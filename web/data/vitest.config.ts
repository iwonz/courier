import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vitest/config";

export default defineConfig({
  resolve: {
    alias: {
      "@courier/ui/relay": fileURLToPath(new URL("../ui/src/assets.ts", import.meta.url)),
      "@courier/ui": fileURLToPath(new URL("../ui/src/index.ts", import.meta.url)),
    },
  },
  test: {
    environment: "jsdom",
    coverage: {
      provider: "v8",
      include: ["src/**/*.ts"],
      exclude: ["src/**/*.test.ts"],
      reporter: ["text"],
      thresholds: { statements: 100, branches: 100, functions: 100, lines: 100 },
    },
  },
});
