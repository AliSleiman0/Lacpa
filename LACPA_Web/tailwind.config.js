/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,js}",
    "./src/components/**/*.{html,js}",
    "./src/pages/**/*.{html,js}"
  ],
  theme: {
    extend: {
      colors: {
        // Custom colors for your admin theme
        'admin-red': 'rgba(214,0,0,0.7)',
        'admin-green': 'rgba(0,121,0,0.7)',
      },
      width: {
        '14/15': '93.333333%',
      }
    },
  },
  plugins: [],
}