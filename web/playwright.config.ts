import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 1 : 0,
  workers: process.env.CI ? 2 : undefined,
  reporter: "line",
  outputDir: process.env.COURIER_PLAYWRIGHT_OUTPUT_DIR ?? "./test-results",
  use: {
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
  },
  projects: [
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
  ],
  webServer: [
    {
      command: "npm run dev --workspace @courier/landing -- --host 127.0.0.1 --port 4173 --strictPort",
      url: "http://127.0.0.1:4173/courier/",
      reuseExistingServer: !process.env.CI,
    },
    {
      command: "npm run dev --workspace @courier/data -- --host 127.0.0.1 --port 4174 --strictPort",
      url: "http://127.0.0.1:4174/",
      reuseExistingServer: !process.env.CI,
    },
    {
      command: "npm run dev --workspace @courier/admin -- --host 127.0.0.1 --port 4175 --strictPort",
      url: "http://127.0.0.1:4175/",
      reuseExistingServer: !process.env.CI,
    },
  ],
});
