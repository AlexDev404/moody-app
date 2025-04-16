const LifecycleManager = (() => {
  return {
    cleanupAll() {
      // Clear intervals
      intervals.forEach(id => originalClearInterval(id));
      intervals.clear();

      // Clear timeouts
      timeouts.forEach(id => originalClearTimeout(id));
      timeouts.clear();

      // Cancel animation frames
      rafs.forEach(id => originalCancelAnimationFrame(id));
      rafs.clear();

      // Remove event listeners
      eventListeners.forEach(({ target, type, handler, options }) => {
        target.removeEventListener(type, handler, options);
      });
      eventListeners.length = 0;

      // Disconnect observers
      observers.forEach(observer => observer.disconnect());
      observers.length = 0;
    }
  };
})();

export default LifecycleManager;