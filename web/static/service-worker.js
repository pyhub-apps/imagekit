// Service Worker for ImageKit PWA
const CACHE_NAME = 'imagekit-v1.0.0';
const urlsToCache = [
  '/',
  '/static/app.js',
  '/static/wasm_exec.js',
  '/static/imagekit.wasm',
  '/manifest.json'
];

// Install event - cache resources
self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => {
        console.log('Opened cache');
        return cache.addAll(urlsToCache);
      })
      .catch(err => {
        console.error('Cache install failed:', err);
      })
  );
  // Force the waiting service worker to become the active service worker
  self.skipWaiting();
});

// Activate event - clean up old caches
self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames.map(cacheName => {
          if (cacheName !== CACHE_NAME && cacheName.startsWith('imagekit-')) {
            console.log('Deleting old cache:', cacheName);
            return caches.delete(cacheName);
          }
        })
      );
    })
  );
  // Take control of all pages immediately
  self.clients.claim();
});

// Fetch event - serve from cache with network fallback
self.addEventListener('fetch', event => {
  // Skip non-GET requests
  if (event.request.method !== 'GET') {
    return;
  }

  // Parse URL
  const url = new URL(event.request.url);
  
  // For WASM and app.js files, use cache-first strategy with version check
  if (url.pathname.includes('.wasm') || url.pathname.includes('app.js')) {
    event.respondWith(
      caches.match(event.request)
        .then(response => {
          // Cache hit - check version parameter
          const fetchPromise = fetch(event.request)
            .then(networkResponse => {
              // Update cache with new version
              if (networkResponse && networkResponse.status === 200) {
                const responseToCache = networkResponse.clone();
                caches.open(CACHE_NAME)
                  .then(cache => {
                    cache.put(event.request, responseToCache);
                  });
              }
              return networkResponse;
            })
            .catch(() => response); // Fallback to cache on network error

          return response || fetchPromise;
        })
    );
    return;
  }

  // For other static assets, use cache-first strategy
  if (url.pathname.startsWith('/static/')) {
    event.respondWith(
      caches.match(event.request)
        .then(response => {
          if (response) {
            return response;
          }
          return fetch(event.request)
            .then(response => {
              // Check if valid response
              if (!response || response.status !== 200 || response.type !== 'basic') {
                return response;
              }
              // Clone the response
              const responseToCache = response.clone();
              // Add to cache
              caches.open(CACHE_NAME)
                .then(cache => {
                  cache.put(event.request, responseToCache);
                });
              return response;
            });
        })
    );
    return;
  }

  // For HTML and other requests, use network-first strategy
  event.respondWith(
    fetch(event.request)
      .then(response => {
        // Check if valid response
        if (!response || response.status !== 200) {
          return response;
        }
        // Clone and cache the response
        const responseToCache = response.clone();
        caches.open(CACHE_NAME)
          .then(cache => {
            cache.put(event.request, responseToCache);
          });
        return response;
      })
      .catch(() => {
        // Network failed, try cache
        return caches.match(event.request);
      })
  );
});

// Handle messages from the client
self.addEventListener('message', event => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
  
  if (event.data && event.data.type === 'CHECK_UPDATE') {
    // Check for updates
    event.waitUntil(
      caches.open(CACHE_NAME)
        .then(cache => {
          return cache.match('/static/app.js')
            .then(response => {
              if (response) {
                return fetch('/static/app.js')
                  .then(networkResponse => {
                    // Compare versions or timestamps
                    event.ports[0].postMessage({
                      updateAvailable: false // Implement version comparison logic
                    });
                  });
              }
            });
        })
    );
  }
});