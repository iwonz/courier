import { expect, test, type Page } from "@playwright/test";

const landingOrigin = `http://127.0.0.1:${process.env.COURIER_LANDING_PORT ?? "4173"}`;
const dataOrigin = `http://127.0.0.1:${process.env.COURIER_DATA_PORT ?? "4174"}`;
const adminOrigin = `http://127.0.0.1:${process.env.COURIER_ADMIN_PORT ?? "4175"}`;
const landingURL = `${landingOrigin}/courier/`;
const dataURL = `${dataOrigin}/`;
const adminURL = `${adminOrigin}/`;
const secretMarker = "COURIER_SECRET_MUST_NOT_RENDER";

async function exerciseThemes(page: Page, root: string): Promise<void> {
  await page.emulateMedia({ colorScheme: "dark", reducedMotion: "reduce" });
  const button = page.locator(`${root} courier-theme-selector button`);
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "system");
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme", "dark");
  await button.click();
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "light");
  await button.click();
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "dark");
  await button.click();
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", "system");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.theme"))).toBe("system");
}

async function selectRussian(page: Page, root: string): Promise<void> {
  const button = page.locator(`${root} courier-locale-selector button`);
  if (await button.textContent() === "🇬🇧") await button.click();
  await expect(button).toContainText("🇷🇺");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.locale"))).toBe("ru");
}

test("landing presents one continuous illustrated journey without execution surfaces", async ({ page, context }) => {
  const externalRequests: string[] = [];
  const sceneRequests: string[] = [];
  page.on("request", (request) => {
    const url = request.url();
    if (url.startsWith("http") && !url.startsWith(landingOrigin)) externalRequests.push(url);
    if (request.resourceType() === "image" && url.includes("landing-")) sceneRequests.push(url);
  });
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto(landingURL);
  const root = "courier-landing-app";

  await expect(page.locator(`${root} h1`)).toHaveText("From here to anywhere.");
  await expect(page.locator(`${root} section`)).toHaveCount(3);
  expect(await page.locator(`${root} section`).evaluateAll((sections) => sections.map((section) => section.id))).toEqual(["route", "install", "cli"]);
  await expect(page.locator(`${root} button[title*="Run"], ${root} button[title*="Replay"], ${root} textarea, ${root} [contenteditable="true"]`)).toHaveCount(0);
  await expect(page.locator(`${root} courier-scene`)).toHaveCount(1);
  await expect(page.locator(`${root} main > courier-scene.page-panorama .base img`)).toHaveCount(1);
  await expect(page.locator(`${root} section courier-scene`)).toHaveCount(0);
  await page.waitForLoadState("networkidle");
  expect(sceneRequests.filter((url) => url.includes("landing-panorama-"))).toHaveLength(1);
  expect(sceneRequests.some((url) => url.includes("landing-panorama-wide-v2"))).toBe(true);

  await expect(page.locator(`${root} .source-endpoints > strong`)).toHaveText("Source");
  await expect(page.locator(`${root} .destination-endpoints > strong`)).toHaveText("Destination");
  await expect(page.locator(`${root} .source-endpoints button > span:last-child`)).toHaveText(["Local", "Remote", "Web", "Web Hook"]);
  await expect(page.locator(`${root} .destination-endpoints button > span:last-child`)).toHaveText(["Local", "Remote", "Web", "Web Hook"]);
  expect(await page.locator(`${root} .source-endpoints courier-icon`).evaluateAll((icons) => icons.map((icon) => icon.getAttribute("name")))).toEqual(["folder-out", "server-out", "browser-upload", "webhook-in"]);
  expect(await page.locator(`${root} .destination-endpoints courier-icon`).evaluateAll((icons) => icons.map((icon) => icon.getAttribute("name")))).toEqual(["folder-in", "server-in", "browser-share", "webhook-out"]);
  await expect(page.locator(`${root} .route-connector path`).first()).toHaveAttribute("d", / C /);

  const github = page.locator(`${root} courier-icon-link a`);
  await expect(github).toHaveAttribute("target", "_blank");
  await expect(github).toHaveAttribute("rel", "noopener noreferrer");
  expect(await page.locator(`${root} a[href^="https://"]`).evaluateAll((links) => links.every((link) => link.getAttribute("target") === "_blank" && link.getAttribute("rel") === "noopener noreferrer"))).toBe(true);
  const controlChrome = await page.locator(`${root} courier-icon-link a, ${root} courier-theme-selector button`).evaluateAll((nodes) => nodes.map((node) => {
    const style = getComputedStyle(node);
    return { background: style.backgroundColor, border: style.borderColor, radius: style.borderRadius, height: node.getBoundingClientRect().height };
  }));
  expect(controlChrome).toHaveLength(2);
  expect(controlChrome[0]!.background).toBe(controlChrome[1]!.background);
  expect(controlChrome[0]!.border).toBe(controlChrome[1]!.border);
  expect(controlChrome[0]!.radius).toBe(controlChrome[1]!.radius);
  expect(Math.abs(controlChrome[0]!.height - controlChrome[1]!.height)).toBeLessThanOrEqual(1);

  const routeCommand = page.locator(`${root} .route-readout .command code`);
  const webSource = page.locator(`${root} .source-endpoints button[data-endpoint="web"]`);
  const initialRoute = await routeCommand.textContent();
  await webSource.hover();
  await expect(routeCommand).toHaveText(initialRoute!);
  await webSource.focus();
  await webSource.press("Enter");
  await expect(routeCommand).toContainText("web://");
  await expect(page.locator(`${root} .destination-endpoints button[data-endpoint="web"]`)).toBeDisabled();
  await context.grantPermissions(["clipboard-read", "clipboard-write"], { origin: landingOrigin });
  await page.locator(`${root} .route-readout button[title="Copy command"]`).click();
  await expect(page.locator(`${root} .route-readout [role="status"]`)).toContainText("Copied");
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toContain("courier from web:// to");

  const installCommand = page.locator(`${root} .install-readout code`);
  const npmChannel = page.locator(`${root} .install-channel[data-channel="npm"]`);
  const initialInstall = await installCommand.textContent();
  await npmChannel.hover();
  await expect(installCommand).toHaveText(initialInstall!);
  await npmChannel.press("Enter");
  await expect(installCommand).toHaveText("npm install --global @iwonz/courier");
  await page.locator(`${root} .install-readout button[title="Copy command"]`).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe("npm install --global @iwonz/courier");

  const checkbox = page.locator(`${root} courier-checkbox input`);
  await expect(checkbox).toBeChecked();
  await expect(checkbox).toBeDisabled();
  const allOptionCount = await page.locator(`${root} .option-row`).count();
  const uiStart = page.locator(`${root} .command-row[data-command="ui-start"]`);
  await uiStart.press("Enter");
  await expect(checkbox).toBeEnabled();
  await expect(page.locator(`${root} .option-row code`)).toHaveText(["--listen <host:port>", "--background"]);
  await page.locator(`${root} courier-checkbox label`).click();
  await expect(page.locator(`${root} .option-row`)).toHaveCount(allOptionCount);
  await page.locator(`${root} .cli-readout button[title="Copy command"]`).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe("courier ui start [options]");

  for (const id of ["install", "cli"] as const) {
    await page.locator(`${root} #${id}`).evaluate((section) => section.scrollIntoView({ block: "start" }));
    await expect(page.locator(`${root} main > courier-scene.page-panorama .base img`)).toHaveCount(1);
    await expect(page.locator(`${root} nav a[href="#${id}"]`)).toHaveAttribute("aria-current", "page");
  }
  await expect.poll(() => page.locator(`${root} courier-scene .base img`).count()).toBe(1);
  expect(sceneRequests.filter((url) => url.includes("landing-panorama-"))).toHaveLength(1);

  const routeInstrument = page.locator(`${root} .route-instrument`);
  const routeReadout = page.locator(`${root} .route-readout`);
  const routeBox = (await routeInstrument.boundingBox())!;
  const routeReadoutBox = (await routeReadout.boundingBox())!;
  for (const source of ["local", "ssh", "web", "webhook"] as const) {
    await page.locator(`${root} .source-endpoints button[data-endpoint="${source}"]`).click();
    for (const destination of ["local", "ssh", "web", "http"] as const) {
      const button = page.locator(`${root} .destination-endpoints button[data-endpoint="${destination}"]`);
      if (await button.isEnabled()) await button.click();
      const instrument = (await routeInstrument.boundingBox())!;
      const readout = (await routeReadout.boundingBox())!;
      expect(Math.abs(instrument.width - routeBox.width)).toBeLessThanOrEqual(1);
      expect(Math.abs(readout.height - routeReadoutBox.height)).toBeLessThanOrEqual(1);
    }
  }

  const installInterface = page.locator(`${root} .install-interface`);
  const installReadout = page.locator(`${root} .install-readout`);
  const installBox = (await installInterface.boundingBox())!;
  const installReadoutBox = (await installReadout.boundingBox())!;
  for (const channel of ["curl", "wget", "PowerShell", "npm", "npx", "Yarn", "pnpm", "Homebrew", "Scoop"] as const) {
    await page.locator(`${root} .install-channel[data-channel="${channel}"]`).click();
    const panel = (await installInterface.boundingBox())!;
    const readout = (await installReadout.boundingBox())!;
    expect(Math.abs(panel.width - installBox.width)).toBeLessThanOrEqual(1);
    expect(Math.abs(readout.height - installReadoutBox.height)).toBeLessThanOrEqual(1);
  }

  await exerciseThemes(page, root);
  await selectRussian(page, root);
  await expect(page.locator(`${root} h1`)).toHaveText("Отсюда — куда угодно.");
  await expect(page.locator(`${root} .source-endpoints > strong`)).toHaveText("Source");
  await page.reload();
  await expect(page.locator(`${root} h1`)).toHaveText("Отсюда — куда угодно.");

  for (const viewport of [
    { width: 320, height: 568 }, { width: 360, height: 740 }, { width: 390, height: 844 }, { width: 1024, height: 600 }, { width: 1440, height: 900 },
  ]) {
    await page.setViewportSize(viewport);
    for (const target of ["route", "install", "cli"] as const) {
      await page.locator(`${root} #${target}`).evaluate((section) => section.scrollIntoView({ block: "start" }));
      const geometry = await page.locator(`${root} #${target}`).evaluate((section) => {
        const bounds = section.getBoundingClientRect();
        const scene = (section.getRootNode() as ShadowRoot).querySelector(".page-panorama")!.getBoundingClientRect();
        const header = (section.getRootNode() as ShadowRoot).querySelector(".masthead-wrap")!.getBoundingClientRect();
        return { bounds, scene, header, visualHeight: visualViewport?.height ?? innerHeight, overflow: document.documentElement.scrollWidth - innerWidth };
      });
      expect(geometry.bounds.height).toBeGreaterThan(0);
      expect(Math.abs(geometry.scene.left)).toBeLessThan(2);
      expect(Math.abs(geometry.scene.right - viewport.width)).toBeLessThan(2);
      expect(geometry.header.top).toBeGreaterThanOrEqual(-0.5);
      expect(geometry.header.bottom).toBeLessThanOrEqual(geometry.visualHeight + 0.5);
      expect(geometry.overflow).toBeLessThanOrEqual(0);
    }
    expect(await page.locator(`${root} .endpoint-terminal`).evaluateAll((nodes) => nodes.every((node) => Math.abs(node.getBoundingClientRect().width - node.getBoundingClientRect().height) <= 1))).toBe(true);
  }

  expect(await page.locator(`${root} courier-scene`).evaluateAll((scenes) => scenes.length === 1 && scenes.every((scene) => scene.shadowRoot!.querySelectorAll("courier-mascot").length === 1))).toBe(true);
  expect(externalRequests).toEqual([]);
});

test("touch landing keeps static scenes and no horizontal overflow", async ({ browser }) => {
  const context = await browser.newContext({ viewport: { width: 320, height: 568 }, isMobile: true, hasTouch: true });
  const page = await context.newPage();
  await page.goto(landingURL);
  const scene = page.locator("courier-landing-app main > courier-scene.page-panorama");
  const before = await scene.locator(".base").evaluate((node) => getComputedStyle(node).transform);
  await page.touchscreen.tap(250, 420);
  expect(await scene.locator(".base").evaluate((node) => getComputedStyle(node).transform)).toBe(before);
  await expect(scene.locator("courier-mascot")).toHaveCount(1);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  await context.close();
});

test("protected delivery keeps metadata private before authentication", async ({ page }) => {
  await page.route("**/api/v1/meta*", (route) => route.fulfill({ status: 401, contentType: "application/json", body: JSON.stringify({ name: secretMarker, path: secretMarker, type: "file" }) }));
  await page.goto(dataURL);
  const root = "courier-data-app";
  const password = page.locator(`${root} input[type="password"]`);
  await expect(password).toBeVisible();
  await expect(page.locator(`${root} courier-scene .base img`)).toHaveAttribute("src", /delivery-access-wide/);
  await expect(page.locator(`${root} courier-workbench.access`)).toBeVisible();
  await expect(page.locator(root)).not.toContainText(secretMarker);
  await exerciseThemes(page, root);
  await password.focus();
  await password.press("Tab");
  await expect(page.locator(`${root} form courier-button button`)).toBeFocused();
  await selectRussian(page, root);
  await expect(page.locator(root)).toContainText("требуется авторизация");
  await page.reload();
  await expect(page.locator(root)).not.toContainText(secretMarker);
});

test("admin exposes a secret-free navigator and API-backed controls", async ({ page }) => {
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
  const root = "courier-admin-app";
  await expect(page.locator(root)).toContainText("127.0.0.1:8080");
  await expect(page.locator(`${root} courier-scene .base img`)).toHaveAttribute("src", /admin-operations-wide/);
  await expect(page.locator(`${root} courier-workbench.registry-workbench`)).toBeVisible();
  await expect(page.locator(`${root} .delivery-nav[aria-pressed="true"]`)).toContainText("path-to-web");
  await expect(page.locator(root)).not.toContainText(secretMarker);
  await page.locator(`${root} form courier-button button`).click();
  await expect.poll(() => mutations.some((url) => url.endsWith("/policy"))).toBe(true);
  await page.locator(`${root} .delivery-head courier-button button`).click();
  await expect.poll(() => mutations.some((url) => url.endsWith("/stop"))).toBe(true);
  await exerciseThemes(page, root);
  await selectRussian(page, root);
  await expect(page.locator(root)).toContainText("Остановить сервер");
  await page.reload();
  await expect(page.locator(root)).not.toContainText(secretMarker);
});

test.describe("browser language negotiation", () => {
  test.use({ locale: "ru-RU" });
  test("starts every surface in Russian without a saved preference", async ({ page }) => {
    await page.route("**/api/v1/meta*", (route) => route.fulfill({ status: 401, body: "{}" }));
    await page.route("**/api/v1/servers", (route) => route.fulfill({ status: 200, contentType: "application/json", body: '{"servers":[]}' }));
    await page.route("**/api/v1/events", (route) => route.fulfill({ status: 200, contentType: "text/event-stream", body: "" }));
    for (const [url, root, text] of [[landingURL, "courier-landing-app", "Отсюда — куда угодно."], [dataURL, "courier-data-app", "требуется авторизация"], [adminURL, "courier-admin-app", "Активные серверы Courier"]] as const) {
      await page.goto(url);
      await expect(page.locator(root)).toContainText(text);
    }
  });
});
