"use strict";

const cacheName = "cache-1";
const toCache = ["base.css"];

this.addEventListener("install", event => {
  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll(toCache);
    })
  );
});
this.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(function(response) {
      return response || fetch(event.request);
    })
  );
});
