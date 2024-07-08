/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        buyMeButton: "#478955",
        background: "#F0F2F5",
        becomeAnAuthorButton: "#002D62",
        searchBoxFill: "#D9D9D9",
        subText: "#6E6E6E",
      },
      fontFamily: {
        inika: ["Inika", "sans-serif"],
      },
    },
  },
  plugins: [],
};
