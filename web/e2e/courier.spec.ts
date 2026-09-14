import { expect, test, type Page } from "@playwright/test";

const landingURL = "http://127.0.0.1:4173/courier/";
const dataURL = "http://127.0.0.1:4174/";
const adminURL = "http://127.0.0.1:4175/";
const secretMarker = "COURIER_SECRET_MUST_NOT_RENDER";

async function selectTheme(page: Page, root: string, preference: "system" | "light" | "dark"): Promise<void> {
  await page.locator(`${root} courier-theme-selector select`).selectOption(preference);
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", preference);
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.theme"))).toBe(preference);
}

async function selectRussian(page: Page, root: string): Promise<void> {
  await page.locator(`${root} courier-locale-selector select`).selectOption("ru");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.locale"))).toBe("ru");
}

async function exerciseThemes(page: Page, root: string): Promise<void> {
  await page.emulateMedia({ colorScheme: "dark", reducedMotion: "reduce" });
  await selectTheme(page, root, "system");
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme", "dark");
  for (const preference of ["light", "dark"] as const) {
    await selectTheme(page, root, preference);
    await expect(page.locator("html")).toHaveAttribute("data-courier-theme", preference);
  }
}

test("landing covers locales, themes, keyboard, and responsive layouts", async ({ page }) => {
  await page.goto(landingURL);
  const root = "courier-landing-app";
  await expect(page.locator(`${root} h1`)).toHaveText("Move files. Keep control.");
  await expect(page.locator(`${root} main`)).toBeVisible();
  await expect(page.locator(`${root} footer`)).toBeVisible();
  await expect(page.locator(`${root} courier-mascot img`).first()).toBeVisible();
  await expect(page.locator(`${root} courier-route`).first()).toBeVisible();
  await expect(page.locator(`${root} code`).filter({ hasText: "npm install --global @iwonz/courier" })).toBeVisible();

  await exerciseThemes(page, root);

  await selectRussian(page, root);
  await expect(page.locator(`${root} h2`).filter({ hasText: "Установить Courier" })).toBeVisible();
  await page.reload();
  await expect(page.locator(`${root} h2`).filter({ hasText: "Установить Courier" })).toBeVisible();

  const installLink = page.locator(`${root} a[href="#install"]`).first();
  await installLink.focus();
  await expect(installLink).toBeFocused();
  await installLink.press("Enter");
  await expect(page).toHaveURL(/#install$/);

  for (const viewport of [{ width: 360, height: 740 }, { width: 1440, height: 900 }]) {
    await page.setViewportSize(viewport);
    await expect(page.locator(`${root} nav`)).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  }
});

test("protected data metadata never renders before authentication", async ({ page }) => {
  await page.route("**/api/v1/meta*", (route) => route.fulfill({
    status: 401,
    contentType: "application/json",
    body: JSON.stringify({ name: secretMarker, path: secretMarker, type: "file" }),
  }));
  await page.goto(dataURL);
  const root = "courier-data-app";
  const password = page.locator(`${root} input[type="password"]`);
  await expect(password).toBeVisible();
  await expect(page.locator(`${root} courier-mascot img`)).toBeVisible();
  await expect(page.locator(root)).not.toContainText(secretMarker);
  await exerciseThemes(page, root);
  await password.focus();
  await expect(password).toBeFocused();
  await password.press("Tab");
  await expect(page.locator(`${root} form courier-button button`)).toBeFocused();
  await selectRussian(page, root);
  await expect(page.locator(root)).toContainText("требуется авторизация");
  await page.reload();
  await expect(page.locator(root)).toContainText("требуется авторизация");
  await expect(page.locator(root)).not.toContainText(secretMarker);
  await page.setViewportSize({ width: 360, height: 740 });
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
});

test("admin renders secret-free state with localized controls", async ({ page }) => {
  const snapshot = {
    servers: [{
      id: "00000000-0000-4000-8000-000000000001",
      bind: "127.0.0.1:8080",
      processId: 42,
      state: "active",
      status: "live",
      startedAt: "2026-09-14T00:00:00Z",
      updatedAt: "2026-09-14T00:00:00Z",
      deliveries: [{
        id: "00000000-0000-4000-8000-000000000002",
        route: "path-to-web",
        source: "./data",
        destination: "web://",
        state: "active",
        password: secretMarker,
        policy: {
          version: 1, auth: "password", authAttempts: 5, authFailAction: "ban",
          deliveryLimit: { unlimited: true, value: 0 }, allowIp: [],
          maxFileSize: { unlimited: false, value: 10 }, maxExtractedSize: { unlimited: false, value: 100 },
          uploadRate: { unlimited: true, value: 0 }, downloadRate: { unlimited: true, value: 0 }, noUi: false,
        },
        counters: { read: 3, sent: 2, confirmed: 1 },
        createdAt: "2026-09-14T00:00:00Z", updatedAt: "2026-09-14T00:00:00Z",
      }],
    }],
  };
  await page.route("**/api/v1/servers", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(snapshot) }));
  await page.route("**/api/v1/events", (route) => route.fulfill({ status: 200, contentType: "text/event-stream", body: "" }));
  await page.goto(adminURL);
  const root = "courier-admin-app";
  await expect(page.locator(root)).toContainText("127.0.0.1:8080");
  await expect(page.locator(`${root} .metrics`)).toContainText("Active deliveries");
  await expect(page.locator(root)).not.toContainText(secretMarker);
  await exerciseThemes(page, root);
  const refresh = page.locator(`${root} nav courier-button button`).first();
  await refresh.focus();
  await expect(refresh).toBeFocused();
  await refresh.press("Enter");
  await expect(page.locator(root)).toContainText("127.0.0.1:8080");
  await selectRussian(page, root);
  await expect(page.locator(root)).toContainText("Остановить сервер");
  await page.reload();
  await expect(page.locator(root)).toContainText("Остановить сервер");
  await expect(page.locator(root)).not.toContainText(secretMarker);
  await page.setViewportSize({ width: 360, height: 740 });
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
});

test.describe("browser language negotiation", () => {
  test.use({ locale: "ru-RU" });

  test("starts every surface in Russian when no preference exists", async ({ page }) => {
    await page.route("**/api/v1/meta*", (route) => route.fulfill({ status: 401, body: "{}" }));
    await page.route("**/api/v1/servers", (route) => route.fulfill({ status: 200, contentType: "application/json", body: '{"servers":[]}' }));
    await page.route("**/api/v1/events", (route) => route.fulfill({ status: 200, contentType: "text/event-stream", body: "" }));
    for (const [url, root, text] of [
      [landingURL, "courier-landing-app", "Установить Courier"],
      [dataURL, "courier-data-app", "требуется авторизация"],
      [adminURL, "courier-admin-app", "Активные серверы Courier"],
    ] as const) {
      await page.goto(url);
      await expect(page.locator(root)).toContainText(text);
    }
  });
});
