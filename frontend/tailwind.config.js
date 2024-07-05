/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        buyMeButton: "#478955",
        background: "#F0F2F5",
        becomeAnAuthorButton: "#2F82FF",
        searchBoxFill: "#D9D9D9",
      },
      fontFamily: {
        inika: ["Inika", "sans-serif"],
      },
    },
  },
  plugins: [],
};
