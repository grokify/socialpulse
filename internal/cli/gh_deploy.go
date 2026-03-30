package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/socialpulse/internal/site"
)

var (
	ghDeployConfigPath string
	ghDeployBranch     string
	ghDeployMessage    string
	ghDeployForce      bool
)

var ghDeployCmd = &cobra.Command{
	Use:   "gh-deploy",
	Short: "Deploy to GitHub Pages",
	Long: `Builds the site and deploys it to GitHub Pages. This command:
1. Builds the static site
2. Commits the output to the gh-pages branch
3. Pushes to the remote repository`,
	RunE: runGHDeploy,
}

func init() {
	ghDeployCmd.Flags().StringVarP(&ghDeployConfigPath, "config", "c", "socialpulse.yaml", "Path to configuration file")
	ghDeployCmd.Flags().StringVarP(&ghDeployBranch, "branch", "b", "gh-pages", "Branch to deploy to")
	ghDeployCmd.Flags().StringVarP(&ghDeployMessage, "message", "m", "Deploy SocialPulse site", "Commit message")
	ghDeployCmd.Flags().BoolVarP(&ghDeployForce, "force", "f", false, "Force push to remote")
}

func runGHDeploy(cmd *cobra.Command, args []string) error {
	config, err := loadConfig(ghDeployConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if we're in a git repository
	if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
		return fmt.Errorf("not a git repository")
	}

	// Get current branch to return to later
	currentBranchOutput, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	currentBranch := strings.TrimSpace(string(currentBranchOutput))

	// Create temporary directory for build
	tmpDir, err := os.MkdirTemp("", "socialpulse-deploy-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Build the site
	fmt.Printf("Building site...\n")
	builder := site.NewBuilder(site.BuilderConfig{
		SiteTitle:       config.Site.Title,
		SiteDescription: config.Site.Description,
		BaseURL:         config.Site.BaseURL,
		SummariesDir:    config.Content.SummariesDir,
		DigestsDir:      config.Content.DigestsDir,
		OutputDir:       tmpDir,
		ThemeName:       config.Theme.Name,
	})

	result, err := builder.Build()
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("  Built %d articles, %d digests, %d pages\n", result.ArticleCount, result.DigestCount, result.PageCount)

	// Check if gh-pages branch exists
	branchExists := false
	branches, err := exec.Command("git", "branch", "-a").Output()
	if err == nil {
		branchExists = strings.Contains(string(branches), ghDeployBranch)
	}

	// Create or checkout gh-pages branch
	fmt.Printf("Preparing %s branch...\n", ghDeployBranch)
	if branchExists {
		// Checkout existing branch
		if err := runGit("checkout", ghDeployBranch); err != nil {
			return fmt.Errorf("failed to checkout %s: %w", ghDeployBranch, err)
		}
	} else {
		// Create orphan branch
		if err := runGit("checkout", "--orphan", ghDeployBranch); err != nil {
			return fmt.Errorf("failed to create %s: %w", ghDeployBranch, err)
		}
		// Remove all files from staging
		_ = runGit("rm", "-rf", ".")
	}

	// Get repository root
	repoRootOutput, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("failed to get repository root: %w", err)
	}
	repoRoot := strings.TrimSpace(string(repoRootOutput))

	// Clean existing files (except .git)
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to read repository: %w", err)
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		path := filepath.Join(repoRoot, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	// Copy built files to repository root
	fmt.Printf("Copying built files...\n")
	if err := copyDir(tmpDir, repoRoot); err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	// Add .nojekyll file for GitHub Pages
	nojekyllPath := filepath.Join(repoRoot, ".nojekyll")
	if err := os.WriteFile(nojekyllPath, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create .nojekyll: %w", err)
	}

	// Stage all changes
	if err := runGit("add", "."); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Check if there are changes to commit
	statusOutput, _ := exec.Command("git", "status", "--porcelain").Output()
	if len(statusOutput) == 0 {
		fmt.Printf("No changes to deploy.\n")
		_ = runGit("checkout", currentBranch)
		return nil
	}

	// Commit
	fmt.Printf("Committing changes...\n")
	if err := runGit("commit", "-m", ghDeployMessage); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	// Push
	fmt.Printf("Pushing to remote...\n")
	pushArgs := []string{"push", "origin", ghDeployBranch}
	if ghDeployForce {
		pushArgs = append(pushArgs, "--force")
	}
	if err := runGit(pushArgs...); err != nil {
		_ = runGit("checkout", currentBranch)
		return fmt.Errorf("failed to push: %w", err)
	}

	// Return to original branch
	if err := runGit("checkout", currentBranch); err != nil {
		return fmt.Errorf("failed to return to %s: %w", currentBranch, err)
	}

	fmt.Printf("\nDeployment complete!\n")
	fmt.Printf("Your site should be available shortly at your GitHub Pages URL.\n")

	return nil
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
