import LifecycleManager from "./lifecycle-manager.js";
const event = new CustomEvent("appLoad", {
  detail: { timestamp: Date.now(), path: location.pathname },
});

document.addEventListener("DOMContentLoaded", () => {
  document.dispatchEvent(event);
});

window.addEventListener("popstate", async () => {
  const route = location.pathname;

  // Handle the new route in your frontend
  console.log("Route changed to:", route);

  // Fetch content from the Go backend (can be an API call or full-page fetch)
  try {
    const response = await fetch(route);
    if (response.ok) {
      const htmlContent = await response.text();
      // Create a temporary DOM element to parse the HTML response
      const parser = new DOMParser();
      const doc = parser.parseFromString(htmlContent, "text/html");

      // Extract the content from the new response's #content element
      const newContent = doc.getElementById("content");
      // If the #content element exists in the new HTML, update the current #content
      if (newContent) {
        LifecycleManager.cleanupAll(); // Clear all intervals if needed
        // Update the #content on the current page
        const contentContainer = document.getElementById("content");
        contentContainer.innerHTML = newContent.innerHTML;

        // Re-run any <script> tags within the new content
        const scripts = newContent.querySelectorAll("script");
        scripts.forEach((script) => {
          // Create a new script element and execute it
          const newScript = document.createElement("script");
          newScript.text = script.textContent;
          document.body.appendChild(newScript);
          document.dispatchEvent(event); // Dispatch the appLoad event again for the new content
          lucide.createIcons(); // Recreate icons if using lucide
          document.body.removeChild(newScript); // Optionally remove after execution
          // Before replacing #content or during route change:
        });
      } else {
        console.error("No #content element in the response.");
      }
    } else {
      console.error("Failed to fetch the route:", response.statusText);
    }
  } catch (error) {
    console.error("Error fetching route:", error);
  }
});
