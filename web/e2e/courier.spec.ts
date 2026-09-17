import { expect, test, type Page } from "@playwright/test";

const landingOrigin = `http://127.0.0.1:${process.env.COURIER_LANDING_PORT ?? "4173"}`;
const dataOrigin = `http://127.0.0.1:${process.env.COURIER_DATA_PORT ?? "4174"}`;
const adminOrigin = `http://127.0.0.1:${process.env.COURIER_ADMIN_PORT ?? "4175"}`;
const landingURL = `${landingOrigin}/courier/`;
const dataURL = `${dataOrigin}/`;
const adminURL = `${adminOrigin}/`;
const secretMarker = "COURIER_SECRET_MUST_NOT_RENDER";

async function exerciseThemes(page: Page): Promise<void> {
  await page.emulateMedia({ colorScheme: "dark", reducedMotion: "reduce" });
  const button = page.getByRole("button", { name: /Theme:|Тема:/ });
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "system");
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme", "dark");
  await button.click();
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "light");
  await button.click();
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "dark");
  await button.click();
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "system");
}

async function selectRussian(page: Page): Promise<void> {
  const button = page.getByRole("button", { name: /Language:|Язык:/ });
  if ((await button.getAttribute("aria-label"))?.startsWith("Language:")) await button.click();
  await expect(button).toHaveAttribute("aria-label", /Русский/);
  await expect(button.locator('[data-locale-icon="ru"]')).toHaveText("🇷🇺");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.locale"))).toBe("ru");
}

test("landing uses a compact square Relay and three natural shadcn sections", async ({ page, context }) => {
  const externalRequests: string[] = [];
  page.on("request", (request) => {
    const url = request.url();
    if (url.startsWith("http") && !url.startsWith(landingOrigin)) externalRequests.push(url);
  });
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto(landingURL);

  await expect(page.getByRole("heading", { level: 1 })).toHaveText("From here to anywhere.");
  await expect(page.getByText("COURIER CLI")).toBeVisible();
  await expect(page.locator("main > section")).toHaveCount(3);
  expect(await page.locator("main > section").evaluateAll((sections) => sections.map((section) => section.id))).toEqual(["route", "install", "cli"]);
  await expect(page.getByRole("button", { name: /Run|Replay/ })).toHaveCount(0);
  await expect(page.locator('img[src*="courier-relay-mark-v2"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-tech-v1"]')).toHaveCount(1);
  await expect(page.locator("[data-courier-route-composition]")).toHaveCount(1);
  await expect(page.locator("[data-courier-cli-registry]")).toHaveCount(1);
  expect(await page.locator("header").evaluate((node) => getComputedStyle(node).borderBottomWidth)).toBe("0px");
  expect(await page.locator("#install").evaluate((node) => [getComputedStyle(node).borderTopWidth, getComputedStyle(node).borderBottomWidth])).toEqual(["0px", "0px"]);
  const heroMascot = page.locator('#route img[width="768"][height="768"]');
  await expect(heroMascot).toBeVisible();
  expect(await heroMascot.evaluate((image) => Math.abs(image.getBoundingClientRect().width - image.getBoundingClientRect().height))).toBeLessThanOrEqual(1);
  await expect(page.locator('path[vector-effect="non-scaling-stroke"]')).toHaveAttribute("d", / C /);

  const github = page.getByRole("link", { name: "Courier on GitHub" });
  await expect(github).toHaveAttribute("target", "_blank");
  await expect(github).toHaveAttribute("rel", "noopener noreferrer");
  expect(await page.locator('a[href^="https://"]').evaluateAll((links) => links.every((link) => link.getAttribute("target") === "_blank" && link.getAttribute("rel") === "noopener noreferrer"))).toBe(true);

  const route = page.locator("#route");
  const routeCommand = route.locator("code").last();
  const webSource = route.getByRole("button", { name: "Web" }).first();
  const initial = await routeCommand.textContent();
  await webSource.hover();
  await expect(routeCommand).toHaveText(initial!);
  await webSource.press("Enter");
  await expect(routeCommand).toContainText("web://");
  await expect(route.getByRole("button", { name: "Web" }).last()).toBeDisabled();

  await context.grantPermissions(["clipboard-read", "clipboard-write"], { origin: landingOrigin });
  await route.getByRole("button", { name: "Copy command" }).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toContain("courier from web:// to");
  const install = page.locator("#install");
  await install.getByRole("button", { name: "npm npm", exact: true }).click();
  await expect(install.locator("code")).toHaveText("npm install --global @iwonz/courier");
  await install.getByRole("button", { name: "Copy command" }).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe("npm install --global @iwonz/courier");

  const cli = page.locator("#cli");
  const compatibility = cli.getByRole("checkbox");
  await expect(compatibility).toBeDisabled();
  await cli.getByRole("button", { name: "courier ui start [options]" }).click();
  await expect(compatibility).toBeEnabled();
  await expect(cli.getByText("--listen <host:port>")).toBeVisible();
  await cli.getByRole("button", { name: "Copy command" }).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe("courier ui start [options]");

  await exerciseThemes(page);
  await selectRussian(page);
  await expect(page.getByRole("heading", { level: 1 })).toHaveText("Отсюда — куда угодно.");
  await page.reload();
  await expect(page.getByRole("heading", { level: 1 })).toHaveText("Отсюда — куда угодно.");

  for (const viewport of [{ width: 320, height: 568 }, { width: 390, height: 844 }, { width: 1024, height: 600 }, { width: 1440, height: 900 }]) {
    await page.setViewportSize(viewport);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
    expect(await page.locator("main > section").evaluateAll((sections) => sections.every((section) => getComputedStyle(section).minHeight === "0px"))).toBe(true);
  }
  expect(externalRequests).toEqual([]);
});

test("touch landing keeps Relay static and avoids horizontal overflow", async ({ browser }) => {
  const context = await browser.newContext({ viewport: { width: 320, height: 568 }, isMobile: true, hasTouch: true });
  const page = await context.newPage();
  await page.goto(landingURL);
  const mascot = page.locator('#route img[width="768"][height="768"]');
  const before = await mascot.evaluate((node) => getComputedStyle(node).transform);
  await page.touchscreen.tap(250, 420);
  expect(await mascot.evaluate((node) => getComputedStyle(node).transform)).toBe(before);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  await context.close();
});

test("protected delivery reveals no metadata before authentication", async ({ page }) => {
  await page.route("**/api/v1/meta*", (route) => route.fulfill({ status: 401, contentType: "application/json", body: JSON.stringify({ name: secretMarker, path: secretMarker, type: "file" }) }));
  await page.goto(dataURL);
  await expect(page.locator('input[type="password"]')).toBeVisible();
  await expect(page.locator('img[src*="courier-relay-mark-v2"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-tech-v1"]')).toHaveCount(1);
  await expect(page.locator("body")).not.toContainText(secretMarker);
  await exerciseThemes(page);
  await selectRussian(page);
  await expect(page.locator("body")).toContainText("требуется авторизация");
  await page.reload();
  await expect(page.locator("body")).not.toContainText(secretMarker);
});

test("admin keeps secrets out of the shadcn control plane and mutates through APIs", async ({ page }) => {
  const snapshot = { servers: [{
    id: "00000000-0000-4000-8000-000000000001", bind: "127.0.0.1:8080", processId: 42, state: "active", status: "live", startedAt: "", updatedAt: "",
    deliveries: [{
      id: "00000000-0000-4000-8000-000000000002", route: "path-to-web", source: "./data", destination: "web://", state: "active", password: secretMarker,
      policy: { version: 1, auth: "password", authAttempts: 5, authFailAction: "ban", deliveryLimit: { unlimited: true, value: 0 }, allowIp: [], maxFileSize: { unlimited: false, value: 10 }, maxExtractedSize: { unlimited: false, value: 100 }, uploadRate: { unlimited: true, value: 0 }, downloadRate: { unlimited: true, value: 0 }, noUi: false },
      counters: { read: 3, sent: 2, confirmed: 1 }, createdAt: "", updatedAt: "",
    }],
  }] };
  const mutations: string[] = [];
  await page.route("**/api/v1/servers", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(snapshot) }));
  await page.route("**/api/v1/events", (route) => route.fulfill({ status: 200, contentType: "text/event-stream", body: "" }));
  await page.route("**/stop", (route) => { mutations.push(route.request().url()); return route.fulfill({ status: 204 }); });
  await page.route("**/policy", (route) => { mutations.push(route.request().url()); return route.fulfill({ status: 204 }); });
  await page.goto(adminURL);
  await expect(page.locator("body")).toContainText("127.0.0.1:8080");
  await expect(page.locator("[data-courier-metrics]")).toHaveCount(1);
  await expect(page.locator("body")).not.toContainText(secretMarker);
  await page.getByRole("button", { name: "Apply policy" }).click();
  await expect.poll(() => mutations.some((url) => url.endsWith("/policy"))).toBe(true);
  await page.getByRole("button", { name: "Stop delivery" }).click();
  await expect.poll(() => mutations.some((url) => url.endsWith("/stop"))).toBe(true);
  await exerciseThemes(page);
  await selectRussian(page);
  await expect(page.locator("body")).toContainText("Остановить сервер");
});

test.describe("browser language negotiation", () => {
  test.use({ locale: "ru-RU" });
  test("starts every surface in Russian without a saved preference", async ({ page }) => {
    await page.route("**/api/v1/meta*", (route) => route.fulfill({ status: 401, body: "{}" }));
    await page.route("**/api/v1/servers", (route) => route.fulfill({ status: 200, contentType: "application/json", body: '{"servers":[]}' }));
    await page.route("**/api/v1/events", (route) => route.fulfill({ status: 200, contentType: "text/event-stream", body: "" }));
    for (const [url, text] of [[landingURL, "Отсюда — куда угодно."], [dataURL, "требуется авторизация"], [adminURL, "Активные серверы Courier"]] as const) {
      await page.goto(url);
      await expect(page.locator("body")).toContainText(text);
    }
  });
});
