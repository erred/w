var t,n,e = function () {return "".concat(Date.now(), "-").concat(Math.floor(8999999999999 * Math.random()) + 1e12);},i = function (t) {var n = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : -1;return { name: t, value: n, delta: 0, entries: [], id: e(), isFinal: !1 };},a = function (t, n) {try {if (PerformanceObserver.supportedEntryTypes.includes(t)) {var e = new PerformanceObserver(function (t) {return t.getEntries().map(n);});return e.observe({ type: t, buffered: !0 }), e;}} catch (t) {}},r = !1,o = !1,s = function (t) {r = !t.persisted;},u = function () {addEventListener("pagehide", s), addEventListener("unload", function () {});},c = function (t) {var n = arguments.length > 1 && void 0 !== arguments[1] && arguments[1];o || (u(), o = !0), addEventListener("visibilitychange", function (n) {var e = n.timeStamp;"hidden" === document.visibilityState && t({ timeStamp: e, isUnloading: r });}, { capture: !0, once: n });},l = function (t, n, e, i) {var a;return function () {e && n.isFinal && e.disconnect(), n.value >= 0 && (i || n.isFinal || "hidden" === document.visibilityState) && (n.delta = n.value - (a || 0), (n.delta || n.isFinal || void 0 === a) && (t(n), a = n.value));};},p = function (t) {var n = arguments.length > 1 && void 0 !== arguments[1] && arguments[1],e = i("CLS", 0),r = function (t) {t.hadRecentInput || (e.value += t.value, e.entries.push(t), s());},o = a("layout-shift", r),s = l(t, e, o, n);c(function (t) {var n = t.isUnloading;o && o.takeRecords().map(r), n && (e.isFinal = !0), s();});},d = function () {return void 0 === t && (t = "hidden" === document.visibilityState ? 0 : 1 / 0, c(function (n) {var e = n.timeStamp;return t = e;}, !0)), { get timeStamp() {return t;} };},m = function (t) {var n = i("FCP"),e = d(),r = a("paint", function (t) {"first-contentful-paint" === t.name && t.startTime < e.timeStamp && (n.value = t.startTime, n.isFinal = !0, n.entries.push(t), o());}),o = l(t, n, r);},v = function (t) {var n = i("FID"),e = d(),r = function (t) {t.startTime < e.timeStamp && (n.value = t.processingStart - t.startTime, n.entries.push(t), n.isFinal = !0, s());},o = a("first-input", r),s = l(t, n, o);c(function () {o && (o.takeRecords().map(r), o.disconnect());}, !0), o || window.perfMetrics && window.perfMetrics.onFirstInputDelay && window.perfMetrics.onFirstInputDelay(function (t, i) {i.timeStamp < e.timeStamp && (n.value = t, n.isFinal = !0, n.entries = [{ entryType: "first-input", name: i.type, target: i.target, cancelable: i.cancelable, startTime: i.timeStamp, processingStart: i.timeStamp + t }], s());});},f = function () {return n || (n = new Promise(function (t) {return ["scroll", "keydown", "pointerdown"].map(function (n) {addEventListener(n, t, { once: !0, passive: !0, capture: !0 });});})), n;},g = function (t) {var n = arguments.length > 1 && void 0 !== arguments[1] && arguments[1],e = i("LCP"),r = d(),o = function (t) {var n = t.startTime;n < r.timeStamp ? (e.value = n, e.entries.push(t)) : e.isFinal = !0, u();},s = a("largest-contentful-paint", o),u = l(t, e, s, n),p = function () {e.isFinal || (s && s.takeRecords().map(o), e.isFinal = !0, u());};f().then(p), c(p, !0);},h = function (t) {var n,e = i("TTFB");n = function () {try {var n = performance.getEntriesByType("navigation")[0] || function () {var t = performance.timing,n = { entryType: "navigation", startTime: 0 };for (var e in t) "navigationStart" !== e && "toJSON" !== e && (n[e] = Math.max(t[e] - t.navigationStart, 0));return n;}();e.value = e.delta = n.responseStart, e.entries = [n], e.isFinal = !0, t(e);} catch (t) {}}, "complete" === document.readyState ? setTimeout(n, 0) : addEventListener("pageshow", n);};export { p as getCLS, m as getFCP, v as getFID, g as getLCP, h as getTTFB };