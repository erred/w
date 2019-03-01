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
  e.respondWith(
    caches.open(cacheName).then(cache => {
      cache.match(e.request).then(response => {
        return response || Promise.reject("no-match");
      });
    })
  );
  e.waitUntil(
    caches.open(cacheName).then(cache => {
      fetch(e.request).then(response => {
        cache.put(e.request, response);
      });
    })
  );
});
