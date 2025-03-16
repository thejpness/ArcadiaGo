import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";
import { VitePWA } from "vite-plugin-pwa";
import tailwindcss from "@tailwindcss/vite"; // ✅ Tailwind v4 optimized plugin

export default defineConfig(({ mode }) => ({
  plugins: [
    vue(),
    vueDevTools(),
    tailwindcss(), // ✅ Optimized Tailwind v4 plugin
    VitePWA({
      registerType: "autoUpdate",
      workbox: {
        navigateFallback: "/index.html",
        navigateFallbackDenylist: [/\/api\//], // Prevent caching API calls
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/your-backend-api\.com\/api\//,
            handler: "NetworkOnly", // ✅ Ensures API calls always hit the server
            options: {
              cacheName: "api-cache",
            }
          }
        ]
      },
      manifest: {
        name: "ArcadiaGo",
        short_name: "ArcadiaGo",
        start_url: "/",
        display: "standalone",
        background_color: "#ffffff",
        theme_color: "#4A90E2",
        icons: [
          { src: "/icons/icon-192x192.png", sizes: "192x192", type: "image/png" },
          { src: "/icons/icon-512x512.png", sizes: "512x512", type: "image/png" }
        ]
      },
      devOptions: {
        enabled: mode !== "development" ? true : false, // ✅ Disable PWA in dev mode
      }
    })
  ],
  resolve: {
    alias: {
      "@": "/src", // ✅ Directly point @ to /src without fileURLToPath
    }
  },
  server: {
    port: 5173, // Default Vite port, change if needed
    strictPort: true, // Ensures it doesn't auto-switch ports
    watch: {
      usePolling: true // ✅ Fixes HMR issues on some systems
    }
  },
  build: {
    sourcemap: false, // Disable sourcemaps in production for performance
    minify: "esbuild", // Use esbuild for faster builds
    chunkSizeWarningLimit: 600 // Prevent warnings on large files
  }
}));
