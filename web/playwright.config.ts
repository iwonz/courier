import { defineConfig, devices } from "@playwright/test";

const landingPort = Number(process.env.COURIER_LANDING_PORT ?? "4173");
const dataPort = Number(process.env.COURIER_DATA_PORT ?? "4174");
const adminPort = Number(process.env.COURIER_ADMIN_PORT ?? "4175");

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
      command: `npm run dev --workspace @courier/landing -- --host 127.0.0.1 --port ${landingPort} --strictPort`,
      url: `http://127.0.0.1:${landingPort}/courier/`,
      reuseExistingServer: false,
    },
    {
      command: `npm run dev --workspace @courier/data -- --host 127.0.0.1 --port ${dataPort} --strictPort`,
      url: `http://127.0.0.1:${dataPort}/`,
      reuseExistingServer: false,
    },
    {
      command: `npm run dev --workspace @courier/admin -- --host 127.0.0.1 --port ${adminPort} --strictPort`,
      url: `http://127.0.0.1:${adminPort}/`,
      reuseExistingServer: false,
    },
  ],
});
