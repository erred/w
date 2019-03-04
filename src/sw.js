"use strict";

const cacheName = "cache-1";
const toCache = ["/", "./base.css"];

this.addEventListener("install", event => {
  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll(toCache);
    })
  );
});

self.addEventListener("fetch", e => {
  let update = true;
  e.respondWith(
    caches.open(cacheName).then(cache => {
      return cache.match(e.request).then(res => {
        return (
          res ||
          fetch(e.request).then(res => {
            update = false;
            cache.put(e.request, res.clone());
            return res;
          })
        );
      });
    })
  );
  if (update) {
    e.waitUntil(
      caches.open(cacheName).then(cache => {
        fetch(e.request).then(res => {
          cache.put(e.request, res);
        });
      })
    );
  }
});
