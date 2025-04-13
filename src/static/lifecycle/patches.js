window.intervals = new Set();
window.timeouts = new Set();
window.rafs = new Set();
window.eventListeners = [];
window.observers = [];

const originalSetInterval = window.setInterval;
const originalClearInterval = window.clearInterval;
const originalSetTimeout = window.setTimeout;
const originalClearTimeout = window.clearTimeout;
const originalRequestAnimationFrame = window.requestAnimationFrame;
const originalCancelAnimationFrame = window.cancelAnimationFrame;

// Monkey patch
window.setInterval = function (fn, delay, ...args) {
  const id = originalSetInterval(fn, delay, ...args);
  intervals.add(id);
  return id;
};

window.clearInterval = function (id) {
  originalClearInterval(id);
  intervals.delete(id);
};

window.setTimeout = function (fn, delay, ...args) {
  const id = originalSetTimeout(fn, delay, ...args);
  timeouts.add(id);
  return id;
};

window.clearTimeout = function (id) {
  originalClearTimeout(id);
  timeouts.delete(id);
};

window.requestAnimationFrame = function (cb) {
  const id = originalRequestAnimationFrame(cb);
  rafs.add(id);
  return id;
};

window.cancelAnimationFrame = function (id) {
  originalCancelAnimationFrame(id);
  rafs.delete(id);
};

window.trackEventListener = function (target, type, handler, options) {
  target.addEventListener(type, handler, options);
  eventListeners.push({ target, type, handler, options });
};

window.trackObserver = function (observer) {
  observers.push(observer);
};
