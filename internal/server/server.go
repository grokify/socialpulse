// Package server implements the development server with live reload.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/grokify/socialpulse/internal/site"
)

// Config contains server configuration.
type Config struct {
	Host            string
	Port            int
	SiteTitle       string
	SiteDescription string
	BaseURL         string
	SummariesDir    string
	DigestsDir      string
	ThemeName       string
	WatchEnabled    bool
}

// Server is the development server.
type Server struct {
	config     Config
	httpServer *http.Server
	outputDir  string
	mu         sync.RWMutex
	rebuilding bool
}

// New creates a new development server.
func New(config Config) *Server {
	return &Server{
		config:    config,
		outputDir: filepath.Join(os.TempDir(), "socialpulse-dev"),
	}
}

// Start starts the server.
func (s *Server) Start() error {
	// Initial build
	if err := s.rebuild(); err != nil {
		return fmt.Errorf("initial build failed: %w", err)
	}

	// Setup file watcher if enabled
	if s.config.WatchEnabled {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return fmt.Errorf("failed to create watcher: %w", err)
		}
		defer watcher.Close()

		// Watch content directories
		if err := s.addWatchRecursive(watcher, s.config.SummariesDir); err != nil {
			log.Printf("Warning: failed to watch summaries dir: %v", err)
		}
		if err := s.addWatchRecursive(watcher, s.config.DigestsDir); err != nil {
			log.Printf("Warning: failed to watch digests dir: %v", err)
		}

		// Start watcher goroutine
		go s.watchLoop(watcher)
	}

	// Create HTTP server
	mux := http.NewServeMux()

	// Serve static files with live reload injection
	mux.HandleFunc("/", s.handleRequest)

	// Live reload endpoint (SSE)
	mux.HandleFunc("/__reload", s.handleReload)

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("\nShutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.httpServer.Shutdown(ctx)
	}()

	return s.httpServer.ListenAndServe()
}

func (s *Server) rebuild() error {
	s.mu.Lock()
	s.rebuilding = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.rebuilding = false
		s.mu.Unlock()
	}()

	builder := site.NewBuilder(site.BuilderConfig{
		SiteTitle:       s.config.SiteTitle,
		SiteDescription: s.config.SiteDescription,
		BaseURL:         fmt.Sprintf("http://%s:%d", s.config.Host, s.config.Port),
		SummariesDir:    s.config.SummariesDir,
		DigestsDir:      s.config.DigestsDir,
		OutputDir:       s.outputDir,
		ThemeName:       s.config.ThemeName,
	})

	result, err := builder.Build()
	if err != nil {
		return err
	}

	log.Printf("Rebuilt: %d articles, %d digests, %d pages", result.ArticleCount, result.DigestCount, result.PageCount)
	return nil
}

func (s *Server) addWatchRecursive(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}

func (s *Server) watchLoop(watcher *fsnotify.Watcher) {
	// Debounce rebuilds
	var timer *time.Timer
	rebuildDelay := 500 * time.Millisecond

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only trigger on write/create events for content files
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			ext := strings.ToLower(filepath.Ext(event.Name))
			if ext != ".yaml" && ext != ".yml" && ext != ".json" {
				continue
			}

			log.Printf("Change detected: %s", event.Name)

			// Debounce
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(rebuildDelay, func() {
				if err := s.rebuild(); err != nil {
					log.Printf("Rebuild failed: %v", err)
				}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Serve from output directory
	filePath := filepath.Join(s.outputDir, filepath.Clean(path))

	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// If directory, serve index.html
	if info.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
	}

	// For HTML files, inject live reload script
	if strings.HasSuffix(filePath, ".html") {
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Inject live reload script before </body>
		liveReloadScript := `<script>
(function() {
  var es = new EventSource('/__reload');
  es.onmessage = function(e) {
    if (e.data === 'reload') {
      window.location.reload();
    }
  };
})();
</script>`
		modified := strings.Replace(string(content), "</body>", liveReloadScript+"</body>", 1)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(modified))
		return
	}

	// Serve other files normally
	http.ServeFile(w, r, filePath)
}

// SSE connections for live reload
var (
	reloadClients   = make(map[chan string]bool)
	reloadClientsMu sync.Mutex
)

func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create client channel
	client := make(chan string)

	reloadClientsMu.Lock()
	reloadClients[client] = true
	reloadClientsMu.Unlock()

	defer func() {
		reloadClientsMu.Lock()
		delete(reloadClients, client)
		reloadClientsMu.Unlock()
		close(client)
	}()

	// Keep connection open
	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg := <-client:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

// TriggerReload sends a reload signal to all connected clients.
func TriggerReload() {
	reloadClientsMu.Lock()
	defer reloadClientsMu.Unlock()

	for client := range reloadClients {
		select {
		case client <- "reload":
		default:
			// Client not ready, skip
		}
	}
}
