import franken from "franken-ui/shadcn-ui/preset-quick";

/** @type {import('tailwindcss').Config} */
export default {
  presets: [franken()],
  content: [
    "./static/errors/**/*.{html,js}",
    "./templates/**/*.{html,js,mustache,tmpl}",
  ],
  safelist: [
    {
      pattern: /^uk-/,
    },
    "ProseMirror",
    "ProseMirror-focused",
    "tiptap",
    "mr-2",
    "mt-2",
    "opacity-50",
  ],
  theme: {
    extend: {
      backgroundImage: {
        'primary-glow': 'var(--primary-glow)',
      },
    },
  },
  plugins: [],
};
