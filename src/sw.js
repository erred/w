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
    caches.match(e.request).then(res => {
      if (res) {
        cache = true;
        return res;
      }
      return fetch(e.request).then(res => {
        if (!response || response.status !== 200 || response.type !== "basic") {
          return response;
        }
        let res2cache = res.clone();
        caches.open(cacheName).then(cache => {
          cache.put(e.request, res2cache);
        });
        return res;
      });
    })
  );
  if (cached) {
    e.waitUntil(
      fetch(e.request).then(res => {
        if (!response || response.status !== 200 || response.type !== "basic") {
          return;
        }
        caches.open(cacheName).then(cache => {
          cache.put(e.request, res);
        });
      })
    );
  }
});
