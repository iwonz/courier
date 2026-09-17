import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  base: "/",
  resolve: {
    alias: {
      "@courier/ui": fileURLToPath(new URL("../ui/src/index.ts", import.meta.url)),
    },
  },
  build: {
    outDir: "../../internal/webdelivery/assets",
    emptyOutDir: true,
    rollupOptions: {
      output: {
        entryFileNames: "assets/data.js",
        chunkFileNames: "assets/[name].js",
        assetFileNames: "assets/data-[name][extname]",
      },
    },
  },
});
