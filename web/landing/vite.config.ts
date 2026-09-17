import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";

export default defineConfig({
  base: "/courier/",
  resolve: {
    alias: {
      "@courier/ui/landing-scenes": fileURLToPath(new URL("../ui/src/landing-scenes.ts", import.meta.url)),
      "@courier/ui": fileURLToPath(new URL("../ui/src/index.ts", import.meta.url)),
    },
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
    rollupOptions: {
      output: {
        entryFileNames: "assets/landing.js",
        chunkFileNames: "assets/[name].js",
        assetFileNames: "assets/landing-[name][extname]",
      },
    },
  },
});
