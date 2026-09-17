import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";

export default defineConfig({
  base: "/",
  resolve: {
    alias: {
      "@courier/ui/delivery-scenes": fileURLToPath(new URL("../ui/src/delivery-scenes.ts", import.meta.url)),
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
