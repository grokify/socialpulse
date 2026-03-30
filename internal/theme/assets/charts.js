// SocialPulse Charts JavaScript
// Chart functionality for the data dashboard theme using ECharts

(function() {
    'use strict';

    // Stance colors matching dashboard.css
    var stanceColors = {
        'agrees': '#198754',      // Green
        'disagrees': '#dc3545',   // Red
        'nuanced': '#6c757d',     // Gray
        'tangential': '#6f42c1',  // Purple
        'default': '#0d6efd'      // Blue
    };

    // Sentiment colors matching dashboard.css
    var sentimentColors = {
        'optimistic': '#22c55e',   // Green
        'cautious': '#f59e0b',     // Amber
        'pessimistic': '#ef4444',  // Red
        'neutral': '#6b7280',      // Gray
        'provocative': '#8b5cf6'   // Purple
    };

    // Animate bars on page load
    document.addEventListener('DOMContentLoaded', function() {
        animateBars();
        initializeTooltips();
        initializePersonaDonutChart();
        initializeSentimentDonutChart();
    });

    function initializePersonaDonutChart() {
        var chartDom = document.getElementById('persona-donut-chart');
        var dataScript = document.getElementById('persona-data');

        if (!chartDom || !dataScript || typeof echarts === 'undefined') {
            return;
        }

        var personas;
        try {
            personas = JSON.parse(dataScript.textContent);
        } catch (e) {
            console.error('Failed to parse persona data:', e);
            return;
        }

        if (!personas || personas.length === 0) {
            return;
        }

        var chart = echarts.init(chartDom);

        var chartData = personas.map(function(p) {
            return {
                name: p.name + ' (' + p.percentage + '%)',
                value: p.percentage,
                stance: p.stance,
                itemStyle: {
                    color: stanceColors[p.stance] || stanceColors['default']
                }
            };
        });

        var option = {
            tooltip: {
                trigger: 'item',
                formatter: function(params) {
                    return params.name + '<br/>Percentage: ' + params.value + '%';
                }
            },
            legend: {
                orient: 'vertical',
                right: 10,
                top: 'center',
                formatter: function(name) {
                    return name.length > 25 ? name.substring(0, 25) + '...' : name;
                },
                textStyle: {
                    fontSize: 11
                }
            },
            series: [{
                name: 'Personas',
                type: 'pie',
                radius: ['45%', '70%'],
                center: ['35%', '50%'],
                avoidLabelOverlap: true,
                itemStyle: {
                    borderRadius: 4,
                    borderColor: '#fff',
                    borderWidth: 2
                },
                label: {
                    show: false,
                    position: 'center'
                },
                emphasis: {
                    label: {
                        show: true,
                        fontSize: 14,
                        fontWeight: 'bold',
                        formatter: function(params) {
                            return params.value + '%';
                        }
                    }
                },
                labelLine: {
                    show: false
                },
                data: chartData
            }]
        };

        chart.setOption(option);

        // Responsive resize
        window.addEventListener('resize', function() {
            chart.resize();
        });
    }

    function initializeSentimentDonutChart() {
        var chartDom = document.getElementById('sentiment-donut-chart');
        var dataScript = document.getElementById('sentiment-data');

        if (!chartDom || !dataScript || typeof echarts === 'undefined') {
            return;
        }

        var sentiments;
        try {
            sentiments = JSON.parse(dataScript.textContent);
        } catch (e) {
            console.error('Failed to parse sentiment data:', e);
            return;
        }

        if (!sentiments) {
            return;
        }

        var chart = echarts.init(chartDom);

        var chartData = [
            { name: 'Optimistic', value: sentiments.optimistic, itemStyle: { color: sentimentColors.optimistic } },
            { name: 'Cautious', value: sentiments.cautious, itemStyle: { color: sentimentColors.cautious } },
            { name: 'Pessimistic', value: sentiments.pessimistic, itemStyle: { color: sentimentColors.pessimistic } },
            { name: 'Neutral', value: sentiments.neutral, itemStyle: { color: sentimentColors.neutral } },
            { name: 'Provocative', value: sentiments.provocative, itemStyle: { color: sentimentColors.provocative } }
        ].filter(function(d) { return d.value > 0; });

        var option = {
            tooltip: {
                trigger: 'item',
                formatter: '{a} <br/>{b}: {c} ({d}%)'
            },
            legend: {
                orient: 'horizontal',
                bottom: 0,
                textStyle: {
                    fontSize: 11
                }
            },
            series: [{
                name: 'Sentiment',
                type: 'pie',
                radius: ['40%', '65%'],
                center: ['50%', '45%'],
                avoidLabelOverlap: true,
                itemStyle: {
                    borderRadius: 4,
                    borderColor: '#fff',
                    borderWidth: 2
                },
                label: {
                    show: true,
                    position: 'inside',
                    formatter: '{c}',
                    fontSize: 11,
                    fontWeight: 'bold',
                    color: '#fff'
                },
                emphasis: {
                    label: {
                        show: true,
                        fontSize: 14,
                        fontWeight: 'bold'
                    }
                },
                labelLine: {
                    show: false
                },
                data: chartData
            }]
        };

        chart.setOption(option);

        // Responsive resize
        window.addEventListener('resize', function() {
            chart.resize();
        });
    }

    function animateBars() {
        const bars = document.querySelectorAll('.bar');
        bars.forEach(function(bar, index) {
            const width = bar.style.width;
            bar.style.width = '0';
            bar.style.transition = 'width 0.5s ease-out';

            setTimeout(function() {
                bar.style.width = width;
            }, index * 50);
        });
    }

    function initializeTooltips() {
        // Add tooltips to elements with data-tooltip attribute
        document.querySelectorAll('[data-tooltip]').forEach(function(elem) {
            elem.addEventListener('mouseenter', showTooltip);
            elem.addEventListener('mouseleave', hideTooltip);
        });
    }

    function showTooltip(e) {
        const target = e.target;
        const text = target.getAttribute('data-tooltip');
        if (!text) return;

        const tooltip = document.createElement('div');
        tooltip.className = 'tooltip';
        tooltip.textContent = text;
        tooltip.style.cssText = 'position:absolute;background:#333;color:#fff;padding:0.5rem;border-radius:4px;font-size:0.75rem;z-index:1000;';

        document.body.appendChild(tooltip);

        const rect = target.getBoundingClientRect();
        tooltip.style.left = rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2) + 'px';
        tooltip.style.top = rect.top - tooltip.offsetHeight - 8 + window.scrollY + 'px';

        target._tooltip = tooltip;
    }

    function hideTooltip(e) {
        if (e.target._tooltip) {
            e.target._tooltip.remove();
            delete e.target._tooltip;
        }
    }

    // Expose for use by other scripts
    window.SocialPulse = {
        animateBars: animateBars
    };
})();
