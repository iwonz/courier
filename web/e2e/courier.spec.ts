import { expect, test, type Page } from "@playwright/test";

const landingOrigin = `http://127.0.0.1:${process.env.COURIER_LANDING_PORT ?? "4173"}`;
const dataOrigin = `http://127.0.0.1:${process.env.COURIER_DATA_PORT ?? "4174"}`;
const adminOrigin = `http://127.0.0.1:${process.env.COURIER_ADMIN_PORT ?? "4175"}`;
const landingURL = `${landingOrigin}/courier/`;
const dataURL = `${dataOrigin}/`;
const adminURL = `${adminOrigin}/`;
const secretMarker = "COURIER_SECRET_MUST_NOT_RENDER";

async function selectTheme(page: Page, root: string, preference: "system" | "light" | "dark"): Promise<void> {
  await page.locator(`${root} courier-theme-selector courier-segmented-control button[data-value="${preference}"]`).click();
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme-preference", preference);
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.theme"))).toBe(preference);
}

async function selectRussian(page: Page, root: string): Promise<void> {
  await page.locator(`${root} courier-locale-selector courier-segmented-control button[data-value="ru"]`).click();
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.locale"))).toBe("ru");
}

async function exerciseThemes(page: Page, root: string): Promise<void> {
  await page.emulateMedia({ colorScheme: "dark", reducedMotion: "reduce" });
  for (const preference of ["light", "dark"] as const) {
    await selectTheme(page, root, preference);
    await expect(page.locator("html")).toHaveAttribute("data-courier-theme", preference);
  }
  await selectTheme(page, root, "system");
  await expect(page.locator("html")).toHaveAttribute("data-courier-theme", "dark");
}

test("landing preserves interaction and geometry across routes, channels, locales, themes, and viewports", async ({ page, context }) => {
  const externalRequests: string[] = [];
  page.on("request", (request) => {
    const url = request.url();
    if (url.startsWith("http") && !url.startsWith(landingOrigin)) externalRequests.push(url);
  });
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto(landingURL);
  const root = "courier-landing-app";

  await expect(page.locator(`${root} h1`)).toHaveText("Move files. Keep control.");
  await expect(page.locator(`${root} section`)).toHaveCount(4);
  expect(await page.locator(root).evaluate((app) => app.shadowRoot!.querySelectorAll(":scope > footer, main > footer").length)).toBe(0);
  await expect(page.locator(`${root} #hero button, ${root} #hero [aria-pressed]`)).toHaveCount(0);
  await expect(page.locator(`${root} .hero-node.source strong`)).toHaveText("Source");
  await expect(page.locator(`${root} .hero-node.destination strong`)).toHaveText("Destination");
  await expect(page.locator(`${root} .hero-route path`).first()).toHaveAttribute("d", / C /);
  await expect(page.locator(`${root} .route-connector path`)).toHaveAttribute("d", / C /);
  await expect(page.locator(`${root} courier-scene`)).toHaveCount(4);
  await expect(page.locator(`${root} courier-scene .base img`)).toHaveCount(4);
  for (const image of await page.locator(`${root} courier-scene .base img`).all()) {
    await expect.poll(() => image.evaluate((node: HTMLImageElement) => node.complete && node.naturalWidth > 0)).toBe(true);
  }
  await expect(page.locator(`${root} .source-endpoints .endpoint-name`)).toHaveText(["Local", "Remote", "Web", "Web Hook"]);
  await expect(page.locator(`${root} .destination-endpoints .endpoint-name`)).toHaveText(["Local", "Remote", "Web", "Web Hook"]);
  expect(await page.locator(`${root} .source-endpoints courier-icon`).evaluateAll((icons) => icons.map((icon) => icon.getAttribute("name")))).toEqual(["folder-out", "server-out", "browser-upload", "webhook-in"]);
  expect(await page.locator(`${root} .destination-endpoints courier-icon`).evaluateAll((icons) => icons.map((icon) => icon.getAttribute("name")))).toEqual(["folder-in", "server-in", "browser-share", "webhook-out"]);
  expect(await page.locator(`${root} a[href^="https://"]`).evaluateAll((links) => links.every((link) => link.getAttribute("target") === "_blank" && link.getAttribute("rel") === "noopener noreferrer"))).toBe(true);

  const heroScene = page.locator(`${root} #hero courier-scene`);
  const baseTransform = await heroScene.locator(".base").evaluate((node) => getComputedStyle(node).transform);
  const beforePointer = await heroScene.evaluate((node) => getComputedStyle(node).getPropertyValue("--scene-pointer-x"));
  await page.mouse.move(1040, 520);
  await expect.poll(() => heroScene.evaluate((node) => node.style.getPropertyValue("--scene-pointer-x"))).not.toBe("");
  expect(await heroScene.locator(".base").evaluate((node) => getComputedStyle(node).transform)).toBe(baseTransform);
  expect(await heroScene.evaluate((node) => getComputedStyle(node).getPropertyValue("--scene-pointer-x"))).not.toBe(beforePointer);

  const routeCommand = page.locator(`${root} .route-readout .prompt code`);
  const webSource = page.locator(`${root} .source-endpoints button[data-endpoint="web"]`);
  const initialRoute = await routeCommand.textContent();
  await webSource.hover();
  await expect(routeCommand).toHaveText(initialRoute!);
  await webSource.focus();
  await expect(routeCommand).toHaveText(initialRoute!);
  await webSource.press("Enter");
  await expect(routeCommand).toContainText("web://");
  await expect(page.locator(`${root} .destination-endpoints button[data-endpoint="web"]`)).toBeDisabled();
  const routeRun = page.locator(`${root} .route-readout button[title="Run demo"]`);
  await routeRun.click();
  await expect(page.locator(`${root} .route-readout courier-terminal .status`)).toHaveText("exit 0");
  await expect(page.locator(`${root} .route-readout ol li`)).toHaveCount(5);
  await expect(page.locator(`${root} .route-readout`)).toContainText("No data left this page");

  const installCommand = page.locator(`${root} .install-readout code`);
  const npmChannel = page.locator(`${root} .install-channel[data-channel="npm"]`);
  const initialInstall = await installCommand.textContent();
  await npmChannel.hover();
  await expect(installCommand).toHaveText(initialInstall!);
  await npmChannel.focus();
  await expect(installCommand).toHaveText(initialInstall!);
  await npmChannel.press("Enter");
  await expect(installCommand).toHaveText("npm install --global @iwonz/courier");
  await context.grantPermissions(["clipboard-read", "clipboard-write"], { origin: landingOrigin });
  await page.locator(`${root} .install-readout button[title="Copy command"]`).click();
  await expect(page.locator(`${root} .install-readout [role="status"]`)).toContainText("Copied");
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe("npm install --global @iwonz/courier");
  const installRun = page.locator(`${root} .install-readout button[title="Run demo"]`);
  await installRun.focus();
  await installRun.press("Enter");
  await expect(page.locator(`${root} .install-readout courier-terminal .status`)).toHaveText("exit 0");
  await expect(page.locator(`${root} .install-readout ol li`)).toHaveCount(5);

  const checkbox = page.locator(`${root} courier-checkbox input`);
  await expect(checkbox).toBeChecked();
  await expect(checkbox).toBeDisabled();
  const allOptionCount = await page.locator(`${root} .option-row`).count();
  const uiStart = page.locator(`${root} .command-row[data-command="ui-start"]`);
  await uiStart.focus();
  await uiStart.press("Enter");
  await expect(checkbox).toBeEnabled();
  await expect(page.locator(`${root} .option-row code`)).toHaveText(["--listen <host:port>", "--background"]);
  const cliRun = page.locator(`${root} .cli-demo button[title="Run demo"]`);
  await cliRun.focus();
  await cliRun.press("Enter");
  await expect(page.locator(`${root} .cli-demo courier-terminal .status`)).toHaveText("exit 0");
  await expect(page.locator(`${root} .cli-demo ol li`)).toHaveCount(4);
  await page.locator(`${root} .cli-demo button[title="Copy command"]`).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe("courier ui start [options]");
  await page.locator(`${root} courier-checkbox label`).click();
  await expect(checkbox).not.toBeChecked();
  await expect(page.locator(`${root} .option-row`)).toHaveCount(allOptionCount);
  await uiStart.press("Enter");
  await expect(checkbox).toBeChecked();
  await expect(checkbox).toBeDisabled();
  await page.locator(`${root} .command-row[data-command="servers"]`).click();
  await expect(page.locator(`${root} .empty-state`)).toHaveText("This command has no options.");
  await page.locator(`${root} .command-row[data-command="servers-stop"]`).click();
  await expect(page.locator(`${root} .option-row code`)).toHaveText(["--all"]);

  const routePanel = page.locator(`${root} .route-explorer`);
  const routeReadout = page.locator(`${root} .route-readout`);
  const routePanelBox = (await routePanel.boundingBox())!;
  const routeReadoutBox = (await routeReadout.boundingBox())!;
  for (const source of ["local", "ssh", "web", "webhook"] as const) {
    await page.locator(`${root} .source-endpoints button[data-endpoint="${source}"]`).click();
    for (const destination of ["local", "ssh", "web", "http"] as const) {
      const button = page.locator(`${root} .destination-endpoints button[data-endpoint="${destination}"]`);
      if (await button.isEnabled()) {
        await button.click();
        const panel = (await routePanel.boundingBox())!;
        const readout = (await routeReadout.boundingBox())!;
        expect(Math.abs(panel.width - routePanelBox.width)).toBeLessThanOrEqual(1);
        expect(Math.abs(panel.height - routePanelBox.height)).toBeLessThanOrEqual(1);
        expect(Math.abs(readout.width - routeReadoutBox.width)).toBeLessThanOrEqual(1);
        expect(Math.abs(readout.height - routeReadoutBox.height)).toBeLessThanOrEqual(1);
      }
    }
  }

  const installPanel = page.locator(`${root} .install-board`);
  const installReadout = page.locator(`${root} .install-readout`);
  const installPanelBox = (await installPanel.boundingBox())!;
  const installReadoutBox = (await installReadout.boundingBox())!;
  for (const channel of ["curl", "wget", "PowerShell", "npm", "npx", "Yarn", "pnpm", "Homebrew", "Scoop"] as const) {
    await page.locator(`${root} .install-channel[data-channel="${channel}"]`).click();
    const panel = (await installPanel.boundingBox())!;
    const readout = (await installReadout.boundingBox())!;
    expect(Math.abs(panel.width - installPanelBox.width)).toBeLessThanOrEqual(1);
    expect(Math.abs(panel.height - installPanelBox.height)).toBeLessThanOrEqual(1);
    expect(Math.abs(readout.width - installReadoutBox.width)).toBeLessThanOrEqual(1);
    expect(Math.abs(readout.height - installReadoutBox.height)).toBeLessThanOrEqual(1);
  }
  await expect(page.locator(`${root} .install-channel courier-brand-icon`)).toHaveCount(9);
  expect(await page.locator(`${root} .endpoint`).evaluateAll((buttons) => buttons.every((button) => button.getBoundingClientRect().height <= 44))).toBe(true);
  expect(await page.locator(`${root} .install-channel`).evaluateAll((buttons) => buttons.every((button) => button.getBoundingClientRect().height <= 40))).toBe(true);

  await exerciseThemes(page, root);
  await expect(page.locator(`${root} courier-theme-selector select`)).toHaveCount(0);
  const localeSymbols = page.locator(`${root} courier-locale-selector courier-segmented-control button .symbol`);
  await expect(localeSymbols).toHaveText(["🇬🇧", "🇷🇺"]);
  const reducedScene = page.locator(`${root} #hero courier-scene`);
  await page.mouse.move(120, 360);
  expect(await reducedScene.locator(".base").evaluate((node) => getComputedStyle(node).transform)).toBe("none");

  await selectRussian(page, root);
  await expect(page.locator(`${root} h2`).filter({ hasText: "Установить Courier" })).toBeVisible();
  await expect(page.locator(`${root} .source-endpoints .label`)).toHaveText("Source");
  await expect(page.locator(`${root} .destination-endpoints .label`)).toHaveText("Destination");
  await expect(page.locator(`${root} .source-endpoints .endpoint-name`)).toHaveText(["Local", "Remote", "Web", "Web Hook"]);
  await page.reload();
  await expect(page.locator(`${root} h2`).filter({ hasText: "Установить Courier" })).toBeVisible();

  for (const viewport of [
    { width: 320, height: 568 },
    { width: 360, height: 740 },
    { width: 390, height: 844 },
    { width: 1024, height: 600 },
    { width: 1440, height: 900 },
  ]) {
    await page.setViewportSize(viewport);
    const expectedSceneSource = viewport.width <= 704 ? "mobile" : "wide";
    await expect.poll(() => page.locator(`${root} courier-scene .base img`).evaluateAll((images, expected) => images.every((image) => (image as HTMLImageElement).currentSrc.includes(expected)), expectedSceneSource)).toBe(true);
    const terminals = await page.locator(`${root} #hero .hero-terminal`).evaluateAll((nodes) => nodes.map((node) => {
      const bounds = node.getBoundingClientRect();
      return { width: bounds.width, height: bounds.height };
    }));
    expect(terminals).toHaveLength(2);
    for (const terminal of terminals) expect(Math.abs(terminal.width - terminal.height)).toBeLessThanOrEqual(1);
    for (const target of ["hero", "routes", "install", "cli"] as const) {
      await page.locator(`${root} #${target}`).evaluate((section) => section.scrollIntoView({ block: "start" }));
      await expect.poll(() => page.locator(`${root} #${target}`).evaluate((section) => Math.abs(section.getBoundingClientRect().top))).toBeLessThan(2);
      const geometry = await page.locator(`${root} #${target}`).evaluate((section) => {
        const bounds = section.getBoundingClientRect();
        const scene = section.querySelector("courier-scene")!.getBoundingClientRect();
        const header = section.getRootNode() instanceof ShadowRoot ? (section.getRootNode() as ShadowRoot).querySelector(".masthead-wrap")!.getBoundingClientRect() : new DOMRect();
        return { bounds, scene, header, visualHeight: visualViewport?.height ?? innerHeight, overflow: document.documentElement.scrollWidth - innerWidth };
      });
      expect(Math.abs(geometry.bounds.height - viewport.height)).toBeLessThanOrEqual(1);
      expect(Math.abs(geometry.scene.left)).toBeLessThan(2);
      expect(Math.abs(geometry.scene.right - viewport.width)).toBeLessThan(2);
      expect(Math.abs(geometry.scene.top)).toBeLessThan(2);
      expect(Math.abs(geometry.scene.bottom - viewport.height)).toBeLessThan(2);
      expect(geometry.header.top).toBeGreaterThanOrEqual(-0.5);
      expect(geometry.header.bottom).toBeLessThanOrEqual(geometry.visualHeight + 0.5);
      expect(geometry.overflow).toBeLessThanOrEqual(0);
      if (target !== "hero") await expect(page.locator(`${root} nav a[href="#${target}"]`)).toHaveAttribute("aria-current", "page");
    }
  }

  expect(externalRequests).toEqual([]);
});

test("touch landing keeps its stationary ambient scene", async ({ browser }) => {
  const context = await browser.newContext({ viewport: { width: 320, height: 568 }, isMobile: true, hasTouch: true });
  const page = await context.newPage();
  await page.goto(landingURL);
  const scene = page.locator("courier-landing-app #hero courier-scene");
  const before = await scene.evaluate((node) => ({ x: node.style.getPropertyValue("--scene-pointer-x"), transform: getComputedStyle(node.shadowRoot!.querySelector(".base")!).transform }));
  await page.touchscreen.tap(250, 420);
  const after = await scene.evaluate((node) => ({ x: node.style.getPropertyValue("--scene-pointer-x"), transform: getComputedStyle(node.shadowRoot!.querySelector(".base")!).transform }));
  expect(after).toEqual(before);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  await context.close();
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
  await expect(page.locator(`${root} courier-scene .base img`)).toBeVisible();
  await expect(page.locator(`${root} courier-scene .base img`)).toHaveAttribute("src", /relay-terminal-delivery-wide/);
  await expect(page.locator(`${root} courier-terminal.access`)).toBeVisible();
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
  await expect(page.locator(`${root} courier-scene .base img`)).toHaveAttribute("src", /relay-terminal-admin-wide/);
  await expect(page.locator(`${root} courier-terminal.registry-terminal`)).toBeVisible();
  await expect(page.locator(`${root} form select`).first()).toHaveCSS("appearance", "none");
  await expect(page.locator(`${root} form input[type="checkbox"]`).first()).toHaveCSS("appearance", "none");
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
