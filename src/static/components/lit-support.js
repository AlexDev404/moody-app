// Icons: https://iconify.design/

// Directives not included in this distribution: https://cdn.jsdelivr.net/gh/lit/dist@3.3.0/all/lit-all.min.js.map
// Tree-shakable directives: https://www.jsdelivr.com/package/npm/lit?tab=files&path=directives
import * as Lit_ from "https://cdn.jsdelivr.net/gh/lit/dist@3/core/lit-core.min.js";
import * as LitUnsafeHTML from "https://cdn.jsdelivr.net/npm/lit-html/directives/unsafe-html.js";

// Merge into a fresh object (Lit stays immutable)
const Lit = { ...Lit_, ...LitUnsafeHTML };
window.Lit = Lit;
// Only for debug: console.log(Lit);
