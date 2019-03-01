"use strict";

const cacheName = "cache-1";
const toCache = ["index.html", "base.css"];

this.addEventListener("install", event => {
  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll(toCache);
    })
  );
});

self.addEventListener("fetch", e => {
  let cached = false;

  e.respondWith(
    caches.open(cacheName).then(cache => {
      cache.match(e.request).then(response => {
        if (response) {
          cached = true;
        }
        return (
          response ||
          fetch(e.request).then(response => {
            cache.put(e.request, response);
            return response;
          })
        );
      });
    })
  );
  if (cached) {
    e.waitUntil(
      caches.open(cacheName).then(cache => {
        fetch(e.request).then(response => {
          cache.put(e.request, response);
        });
      })
    );
  }
});
