import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const isProduction = mode === 'production'

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      port: 3000,
      host: true,
    },
    build: {
      // Output directory for production build
      outDir: 'dist',
      // Generate source maps in development, hidden source maps in production
      sourcemap: isProduction ? 'hidden' : true,
      // Minification options (default is esbuild, can also use 'terser' for more aggressive minification)
      minify: isProduction ? 'esbuild' : false,
      // CSS configuration
      cssMinify: isProduction ? 'esbuild' : false,
      // Target modern browsers for smaller bundles
      target: 'es2020',
      // Rollup options for advanced chunking and output configuration
      rollupOptions: {
        output: {
          // Entry file names with content hash for cache busting
          entryFileNames: 'assets/[name]-[hash].js',
          // Chunk file names with content hash
          chunkFileNames: 'assets/[name]-[hash].js',
          // Asset file names with content hash
          assetFileNames: (assetInfo) => {
            // Keep fonts in a separate folder
            if (/\.(woff2?|ttf|otf|eot)$/.test(assetInfo.name ?? '')) {
              return 'assets/fonts/[name]-[hash][extname]'
            }
            // Keep images in a separate folder
            if (/\.(png|jpe?g|gif|svg|webp|ico)$/.test(assetInfo.name ?? '')) {
              return 'assets/images/[name]-[hash][extname]'
            }
            // Default for other assets (CSS, etc.)
            return 'assets/[name]-[hash][extname]'
          },
          // Manual chunks for better code splitting and caching
          manualChunks: (id) => {
            // Vendor chunk: all node_modules except specific large libraries
            // We split out only the largest dependencies for better caching
            if (id.includes('node_modules')) {
              // Axios - large HTTP client, changes infrequently
              if (id.includes('node_modules/axios')) {
                return 'vendor-axios'
              }
              // Radix UI - large UI library, changes infrequently
              if (id.includes('node_modules/@radix-ui')) {
                return 'vendor-radix'
              }
              // All other dependencies go into a single vendor chunk
              // This avoids circular dependency issues while still providing
              // separation between app code and vendor code
              return 'vendor'
            }
          },
        },
      },
      // Chunk size warning limit (in KB)
      chunkSizeWarningLimit: 500,
      // Empty output directory before building
      emptyOutDir: true,
    },
    // CSS configuration
    css: {
      // Enable CSS source maps in development
      devSourcemap: !isProduction,
      // PostCSS configuration (uses postcss.config.js if present)
      postcss: {
        // Additional PostCSS options if needed
      },
    },
    // Dependency optimization for development
    optimizeDeps: {
      // Pre-bundle these dependencies for faster dev server startup
      include: [
        'react',
        'react-dom',
        'react-router-dom',
        'axios',
        'zod',
        'react-hook-form',
        '@hookform/resolvers',
        'lucide-react',
      ],
      // Exclude these from pre-bundling if they cause issues
      exclude: [],
    },
    // Preview server configuration (for `vite preview`)
    preview: {
      port: 4173,
      host: true,
    },
    // Esbuild options
    esbuild: {
      // Drop console and debugger statements in production
      drop: isProduction ? ['console', 'debugger'] : [],
      // Keep class names for better debugging in development
      keepNames: !isProduction,
    },
  }
})
