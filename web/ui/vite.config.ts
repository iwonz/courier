import { defineConfig } from "vite";

export default defineConfig({
  build: {
    emptyOutDir: true,
    lib: {
      entry: "src/index.ts",
      formats: ["es"],
      fileName: "index",
      cssFileName: "courier-ui",
    },
    rollupOptions: {
      external: [/^lit(?:\/.*)?$/],
    },
  },
});
