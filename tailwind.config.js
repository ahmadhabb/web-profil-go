/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.html", "./static/js/*.js"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Poppins', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
