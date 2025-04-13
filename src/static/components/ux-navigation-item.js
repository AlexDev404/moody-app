class UxNavigationItem extends HTMLElement {
  constructor() {
	super();
	this.attachShadow({ mode: "open" });
	this.onClick = this.onClick.bind(this);
  }

  connectedCallback() {
	this.render();
	this.shadowRoot.querySelector("a").addEventListener("click", this.onClick);
  }

  disconnectedCallback() {
	this.shadowRoot
	  .querySelector("a")
	  .removeEventListener("click", this.onClick);
  }

  static get observedAttributes() {
	return ["href", "target", "rel"];
  }

  attributeChangedCallback() {
	this.render();
  }

  onClick(event) {
	const href = this.getAttribute("href");
	const target = this.getAttribute("target");

	// Only intercept same-origin links and no target="_blank"
	if (
	  href &&
	  !href.startsWith("http") && // Skip absolute URLs
	  !href.startsWith("//") &&
	  target !== "_blank"
	) {
	  event.preventDefault();
	  history.pushState({}, "", href);
	  window.dispatchEvent(new Event("popstate")); // Let your router know
	}
  }

  render() {
	const href = this.getAttribute("href") || "#";
	const target = this.getAttribute("target") || "";
	const rel = this.getAttribute("rel") || "";

	this.shadowRoot.innerHTML = `
		<link rel="stylesheet" href="/static/style.css" />
		<a class="ux-link" href="${href}" target="${target}" rel="${rel}">
		  	<slot></slot>
		</a>
	  `;
  }
}

export default UxNavigationItem;
