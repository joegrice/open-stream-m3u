(function() {
    'use strict';

    const THEME_KEY = 'open-stream-theme';

    function initTheme() {
        const saved = localStorage.getItem(THEME_KEY);
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const theme = saved || (prefersDark ? 'dark' : 'light');
        document.documentElement.setAttribute('data-theme', theme);
    }

    function toggleTheme() {
        const current = document.documentElement.getAttribute('data-theme');
        const next = current === 'dark' ? 'light' : 'dark';
        document.documentElement.setAttribute('data-theme', next);
        localStorage.setItem(THEME_KEY, next);
    }

    document.addEventListener('DOMContentLoaded', function() {
        initTheme();

        const themeBtn = document.getElementById('themeToggle');
        if (themeBtn) {
            themeBtn.addEventListener('click', toggleTheme);
        }

        const tabs = document.querySelectorAll('.tab');
        const contents = document.querySelectorAll('.tab-content');

        tabs.forEach(function(tab) {
            tab.addEventListener('click', function() {
                const target = this.getAttribute('data-tab');

                tabs.forEach(function(t) { t.classList.remove('active'); });
                contents.forEach(function(c) { c.classList.remove('active'); });

                this.classList.add('active');
                var targetContent = document.getElementById(target + 'Tab');
                if (targetContent) {
                    targetContent.classList.add('active');
                }
            });
        });

        const params = new URLSearchParams(window.location.search);
        const mode = params.get('mode');
        if (mode) {
            var targetTab = document.querySelector('.tab[data-tab="' + mode + '"]');
            if (targetTab) {
                targetTab.click();
            }
        }

        var form = document.getElementById('configForm');
        if (form) {
            form.addEventListener('submit', handleInstall);
        }

        var copyBtn = document.getElementById('copyBtn');
        if (copyBtn) {
            copyBtn.addEventListener('click', copyManifestUrl);
        }
    });

    async function handleInstall(e) {
        e.preventDefault();
        console.log('Form submitted, starting install...');

        var activeTab = document.querySelector('.tab.active');
        var mode = activeTab ? activeTab.getAttribute('data-tab') : 'direct';
        console.log('Active mode:', mode);

        var config = {};

        if (mode === 'direct') {
            var m3uUrl = document.getElementById('m3uUrl').value.trim();
            console.log('M3U URL:', m3uUrl);
            if (!m3uUrl) {
                alert('Please enter an M3U URL');
                return;
            }

            config = {
                provider: 'direct',
                m3uUrl: m3uUrl,
                enableEpg: document.getElementById('enableEpg').checked,
                epgUrl: document.getElementById('epgUrl').value.trim(),
                epgOffsetHours: parseFloat(document.getElementById('epgOffset').value) || 0
            };
        } else {
            var xtreamUrl = document.getElementById('xtreamUrl').value.trim();
            var xtreamUsername = document.getElementById('xtreamUsername').value.trim();
            var xtreamPassword = document.getElementById('xtreamPassword').value.trim();

            if (!xtreamUrl || !xtreamUsername || !xtreamPassword) {
                alert('Please fill in all Xtream fields');
                return;
            }

            config = {
                provider: 'xtream',
                xtreamUrl: xtreamUrl,
                xtreamUsername: xtreamUsername,
                xtreamPassword: xtreamPassword,
                xtreamUseM3U: document.getElementById('xtreamUseM3U').checked,
                enableEpg: document.getElementById('xtreamEnableEpg').checked
            };
        }

        console.log('Config:', config);
        showProgress();
        await buildAndInstall(config);
    }

    function showProgress() {
        var overlay = document.getElementById('progressOverlay');
        var result = document.getElementById('resultPanel');
        if (overlay) overlay.classList.remove('hidden');
        if (result) result.classList.add('hidden');
        updateProgress(10, 'Building configuration...');
    }

    function updateProgress(percent, text) {
        var fill = document.getElementById('progressFill');
        var textEl = document.getElementById('progressText');
        if (fill) fill.style.width = percent + '%';
        if (textEl) textEl.textContent = text;
    }

    function updateStats(stats) {
        var el = function(id) { return document.getElementById(id); };
        var ids = ['statChannels', 'statMovies', 'statSeries', 'statEPG'];
        var values = [stats.channels || 0, stats.movies || 0, stats.series || 0, stats.epgChannels || 0];
        
        for (var i = 0; i < ids.length; i++) {
            var element = el(ids[i]);
            if (element) {
                element.textContent = formatNumber(values[i]);
                element.classList.remove('loading');
            }
        }
    }

    function setLoadingStats() {
        var ids = ['statChannels', 'statMovies', 'statSeries', 'statEPG'];
        for (var i = 0; i < ids.length; i++) {
            var element = document.getElementById(ids[i]);
            if (element) {
                element.textContent = '...';
                element.classList.add('loading');
            }
        }
    }

    function formatNumber(n) {
        if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M';
        if (n >= 1000) return (n / 1000).toFixed(1) + 'K';
        return n.toLocaleString();
    }

    async function buildAndInstall(config) {
        updateProgress(10, 'Testing connection...');

        try {
            // Test connection based on provider type
            if (config.provider === 'direct' && config.m3uUrl) {
                const prefetchResp = await fetch('/api/prefetch', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ url: config.m3uUrl, purpose: 'playlist' })
                });

                if (!prefetchResp.ok) {
                    throw new Error('Failed to fetch M3U playlist: ' + prefetchResp.statusText);
                }

                const prefetchData = await prefetchResp.json();
                if (!prefetchData.ok) {
                    throw new Error('M3U playlist not accessible');
                }

                updateProgress(40, `M3U playlist loaded (${(prefetchData.bytes / 1024).toFixed(1)} KB)`);
            } else if (config.provider === 'xtream' && config.xtreamUrl) {
                // Test Xtream API connection
                const apiUrl = config.xtreamUrl.replace(/\/$/, '') + '/player_api.php?username=' + 
                    encodeURIComponent(config.xtreamUsername) + '&password=' + 
                    encodeURIComponent(config.xtreamPassword);
                
                const prefetchResp = await fetch('/api/prefetch', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ url: apiUrl, purpose: 'xtream_auth' })
                });

                if (!prefetchResp.ok) {
                    throw new Error('Failed to connect to Xtream panel: ' + prefetchResp.statusText);
                }

                updateProgress(40, 'Xtream panel connected');
            }

            // Test EPG URL if provided
            if (config.enableEpg && config.epgUrl) {
                updateProgress(50, 'Testing EPG connection...');

                const epgResp = await fetch('/api/prefetch', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ url: config.epgUrl, purpose: 'epg' })
                });

                if (epgResp.ok) {
                    const epgData = await epgResp.json();
                    if (epgData.ok) {
                        updateProgress(70, `EPG loaded (${(epgData.bytes / 1024).toFixed(1)} KB)`);
                    }
                } else {
                    console.warn('EPG fetch failed, continuing without it');
                    updateProgress(70, 'EPG not accessible (continuing without it)');
                }
            } else if (config.provider === 'xtream' && config.enableEpg) {
                updateProgress(70, 'Using panel EPG');
            } else {
                updateProgress(70, 'Skipping EPG test');
            }

            updateProgress(85, 'Generating configuration token...');

            // Encode config to base64url
            const jsonStr = JSON.stringify(config);
            const token = btoa(jsonStr)
                .replace(/\+/g, '-')
                .replace(/\//g, '_')
                .replace(/=+$/, '');

            updateProgress(95, 'Building manifest URL...');

            const baseUrl = window.location.origin;
            const manifestUrl = baseUrl + '/' + token + '/manifest.json';
            const stremioUrl = 'stremio://' + baseUrl.replace(/^https?:\/\//, '') + '/' + token + '/manifest.json';

            updateProgress(100, 'Complete!');

            // Show result panel immediately
            setTimeout(function() {
                var overlay = document.getElementById('progressOverlay');
                var result = document.getElementById('resultPanel');
                var manifestCodeBlock = document.getElementById('manifestCodeBlock');
                var stremioLink = document.getElementById('stremioLink');

                if (overlay) overlay.classList.add('hidden');
                if (result) result.classList.remove('hidden');
                if (manifestCodeBlock) manifestCodeBlock.textContent = manifestUrl;
                if (stremioLink) stremioLink.href = stremioUrl;
                
                // Show loading state for stats
                setLoadingStats();
            }, 500);

            // Fetch stats asynchronously
            try {
                const statsResp = await fetch('/api/info?token=' + encodeURIComponent(token));
                if (statsResp.ok) {
                    const stats = await statsResp.json();
                    updateStats(stats);
                }
            } catch (e) {
                console.warn('Failed to fetch stats:', e);
            }

        } catch (error) {
            console.error('Install error:', error);
            alert('Installation failed: ' + error.message);
            var overlay = document.getElementById('progressOverlay');
            if (overlay) overlay.classList.add('hidden');
        }
    }

    function copyManifestUrl() {
        var codeBlock = document.getElementById('manifestCodeBlock');
        if (!codeBlock) return;

        var text = codeBlock.textContent;
        if (!text) return;

        if (navigator.clipboard) {
            navigator.clipboard.writeText(text).then(function() {
                var btn = document.getElementById('copyBtn');
                if (btn) {
                    var original = btn.textContent;
                    btn.textContent = 'Copied!';
                    setTimeout(function() { btn.textContent = original; }, 2000);
                }
            });
        } else {
            // Fallback
            var range = document.createRange();
            range.selectNodeContents(codeBlock);
            var selection = window.getSelection();
            selection.removeAllRanges();
            selection.addRange(range);
            document.execCommand('copy');
            selection.removeAllRanges();
        }
    }
})();
