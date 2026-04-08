import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { resolve } from 'node:path'

const webpagesDir = resolve(__dirname, 'Webpages')

// https://vite.dev/config/
export default defineConfig({
  // Treat the Webpages folder as the project root so all HTML entry files,
  // /src, and svelte.config.js stay inside Webpages/ for organisation.
  root: webpagesDir,
  plugins: [svelte()],
  build: {
    rollupOptions: {
      input: {
        home: resolve(webpagesDir, 'index.html'),
        demographicQuiz: resolve(webpagesDir, 'demographic-quiz.html'),
        swiping: resolve(webpagesDir, 'swiping.html'),
        results: resolve(webpagesDir, 'results.html'),
        settings: resolve(webpagesDir, 'settings.html'),
        adminHome: resolve(webpagesDir, 'admin-home.html'),
      },
    },
  },
})
