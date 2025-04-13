class UxAvatar extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: "open" });
	}

	connectedCallback() {
		this.render();
	}

	static get observedAttributes() {
		return ["src", "alt", "size"];
	}

	attributeChangedCallback() {
		this.render();
	}

	render() {
		const src = this.getAttribute("src") || "";
		const alt = this.getAttribute("alt") || "User avatar";
		const size = this.getAttribute("size") || "40px";

		this.shadowRoot.innerHTML = `
			<link rel="stylesheet" href="/static/style.css" />
			<style>
				.avatar {
					inline-size: ${size};
					block-size: ${size};
					border-radius: 100%;
					object-fit: cover;
				}
			</style>
			<img class="avatar" src="${src}" alt="${alt}" />
		`;
	}
}

export default UxAvatar;
customElements.define("ux-avatar", UxAvatar);
