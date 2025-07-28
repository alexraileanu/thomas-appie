import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import fs from 'fs'
import path from 'path'

// Custom plugin to move favicon files to assets folder
const moveFaviconsPlugin = () => {
  return {
    name: 'move-favicons',
    generateBundle(options, bundle) {
      // List of favicon files to move
      const faviconFiles = [
        'favicon-16x16.png',
        'favicon-32x32.png',
        'apple-touch-icon.png',
        'apple-touch-icon-57x57.png',
        'apple-touch-icon-60x60.png',
        'apple-touch-icon-72x72.png',
        'apple-touch-icon-76x76.png',
        'apple-touch-icon-114x114.png',
        'apple-touch-icon-120x120.png',
        'apple-touch-icon-144x144.png',
        'apple-touch-icon-152x152.png',
        'vite.svg'
      ]

      // Add favicon files to bundle as assets
      faviconFiles.forEach(filename => {
        const filepath = path.resolve(__dirname, 'public', filename)
        if (fs.existsSync(filepath)) {
          const source = fs.readFileSync(filepath)
          this.emitFile({
            type: 'asset',
            fileName: `assets/${filename}`,
            source
          })
        }
      })
    }
  }
}

export default defineConfig({
  plugins: [vue(), moveFaviconsPlugin()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  // Exclude favicon files from public directory copying
  publicDir: false,
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html')
      }
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:7008',
        changeOrigin: true,
      }
    }
  }
})