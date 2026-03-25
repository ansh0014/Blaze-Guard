/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        fiery: {
          DEFAULT: '#ff4d00', // Fiery Orange
          500: '#ff4d00',
          600: '#e64500',
          glow: '#ff4d0020'
        },
        // Emerald is already in Tailwind, but we can define specific semantic aliases if we want
        primary: {
          light: '#34d399', // emerald-400
          DEFAULT: '#10b981', // emerald-500
          dark: '#059669', // emerald-600
        },
        background: {
          DEFAULT: '#0a0a0a', // Premium off-black
          paper: '#121212', // Slightly lighter
        }
      },
      fontFamily: {
        sans: ['Outfit', 'Inter', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
