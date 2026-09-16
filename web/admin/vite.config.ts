import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";

export default defineConfig({
  base: "/",
  resolve: {
    alias: {
      "@courier/ui/relay-admin": fileURLToPath(new URL("../ui/src/relay-admin.ts", import.meta.url)),
      "@courier/ui": fileURLToPath(new URL("../ui/src/index.ts", import.meta.url)),
    },
  },
  build: {
    outDir: "../../internal/admin/assets",
    emptyOutDir: true,
    rollupOptions: {
      output: {
        entryFileNames: "assets/admin.js",
        chunkFileNames: "assets/[name].js",
        assetFileNames: ({ names }) => names.some((name) => name.endsWith(".svg"))
          ? "assets/admin-[name][extname]"
          : "assets/admin[extname]",
      },
    },
  },
});
