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
  await expect(button.locator('[data-locale-icon="ru"] i')).toHaveCount(3);
  await expect.poll(() => page.evaluate(() => localStorage.getItem("courier.locale"))).toBe("ru");
}

test("landing keeps the static header, official brands, route signal, and command builder exact", async ({ page, context }) => {
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
  await expect(page.locator('img[src*="courier-relay-pixel-mark-v3"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-pixel-route-v3"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-pixel-delivery-v3"], img[src*="courier-relay-pixel-admin-v3"], img[src*="courier-relay-pixel-neutral-v3"]')).toHaveCount(0);
  await expect(page.locator("[data-courier-route-composition]")).toHaveCount(1);
  await expect(page.locator("[data-courier-contract-metrics]")).toHaveCount(0);
  await expect(page.locator("[data-courier-cli-registry]")).toHaveCount(1);
  for (const selector of ["[data-courier-route-composition]", "[data-courier-cli-registry]"]) {
    expect(await page.locator(selector).evaluate((node) => {
      const style = getComputedStyle(node);
      return [style.backgroundImage, style.borderRadius, style.boxShadow, style.backdropFilter];
    })).toEqual(["none", "0px", "none", "none"]);
  }
  await expect(page.locator("header nav")).toHaveCount(0);
  expect(await page.locator("header").evaluate((node) => {
    const style = getComputedStyle(node);
    return [style.position, style.backgroundColor, style.backgroundImage, style.backdropFilter, style.borderBottomWidth];
  })).toEqual(["static", "rgb(255, 255, 255)", "none", "none", "0px"]);
  const headerCenters = await page.locator("header").evaluate((header) => {
    const brand = header.querySelector(":scope > div > a")!.getBoundingClientRect();
    const controls = Array.from(header.querySelectorAll(":scope > div > div > a, :scope > div > div > button"));
    return [brand.top + brand.height / 2, ...controls.map((control) => {
      const bounds = control.getBoundingClientRect();
      return bounds.top + bounds.height / 2;
    })];
  });
  expect(Math.max(...headerCenters) - Math.min(...headerCenters)).toBeLessThanOrEqual(1);
  await page.evaluate(() => scrollTo(0, 600));
  await expect.poll(() => page.locator("header").evaluate((node) => node.getBoundingClientRect().bottom)).toBeLessThan(0);
  await page.evaluate(() => scrollTo(0, 0));
  await expect.poll(() => page.evaluate(() => scrollY)).toBe(0);
  expect(await page.locator("#install").evaluate((node) => [getComputedStyle(node).borderTopWidth, getComputedStyle(node).borderBottomWidth])).toEqual(["0px", "0px"]);
  expect(await page.locator("[data-courier-route-composition]").evaluate((node) => {
    const style = getComputedStyle(node);
    return [style.borderTopWidth, style.borderRightWidth, style.borderBottomWidth, style.borderLeftWidth];
  })).toEqual(["0px", "0px", "0px", "0px"]);
  const activeSource = page.locator('#route button[aria-label="Local"][aria-pressed="true"]');
  expect(await activeSource.evaluate((node) => {
    const style = getComputedStyle(node);
    return style.backgroundColor !== "rgba(0, 0, 0, 0)" && style.color !== style.backgroundColor;
  })).toBe(true);
  const routeBottom = await page.locator("[data-courier-route-composition]").evaluate((node) => node.getBoundingClientRect().bottom + scrollY);
  const installTop = await page.locator("#install h2").evaluate((node) => node.getBoundingClientRect().top + scrollY);
  const installBottom = await page.locator("#install > div").evaluate((node) => node.getBoundingClientRect().bottom + scrollY);
  const cliTop = await page.locator("#cli h2").evaluate((node) => node.getBoundingClientRect().top + scrollY);
  expect(installTop - routeBottom).toBeLessThanOrEqual(80);
  expect(cliTop - installBottom).toBeLessThanOrEqual(80);
  const heroMascot = page.locator('#route img[width="512"][height="512"]');
  await expect(heroMascot).toBeVisible();
  expect(await heroMascot.evaluate((image) => Math.abs(image.getBoundingClientRect().width - image.getBoundingClientRect().height))).toBeLessThanOrEqual(1);
  await expect(page.locator('path[vector-effect="non-scaling-stroke"]')).toHaveAttribute("d", / H .* V /);

  const brandImages = page.locator('img[src*="/brands/"]');
  expect(await brandImages.count()).toBeGreaterThan(0);
  expect(await brandImages.evaluateAll((images) => images.every((image) => image.getAttribute("src")?.endsWith(".png") && !["pixelated", "crisp-edges"].includes(getComputedStyle(image).imageRendering)))).toBe(true);
  await expect(page.locator("#install button", { hasText: "npx" }).locator("img")).toHaveCount(0);
  await expect(page.locator("#install button", { hasText: "wget" }).locator("img")).toHaveCount(0);

  const signal = page.locator("[data-courier-route-signal]");
  expect(await signal.evaluate((node) => [node.getBoundingClientRect().width, node.getBoundingClientRect().height])).toEqual([4, 4]);
  expect(await signal.evaluate((node) => getComputedStyle(node).animationTimingFunction)).toBe("linear");
  await page.emulateMedia({ reducedMotion: "reduce" });
  const signalAlignment = await page.locator("[data-courier-route-path]").evaluate((pathNode) => {
    const path = pathNode as SVGPathElement;
    const signalNode = path.ownerSVGElement!.parentElement!.querySelector("[data-courier-route-signal]")!;
    const point = path.getPointAtLength(path.getTotalLength() / 2).matrixTransform(path.getScreenCTM()!);
    const bounds = signalNode.getBoundingClientRect();
    return { dx: Math.abs(bounds.left + bounds.width / 2 - point.x), dy: Math.abs(bounds.top + bounds.height / 2 - point.y), width: bounds.width, height: bounds.height, animation: getComputedStyle(signalNode).animationName };
  });
  expect(signalAlignment.width).toBeCloseTo(4, 0);
  expect(signalAlignment.height).toBeCloseTo(4, 0);
  expect(signalAlignment.dx).toBeLessThanOrEqual(1);
  expect(signalAlignment.dy).toBeLessThanOrEqual(1);
  expect(signalAlignment.animation).toBe("none");

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
  await install.getByRole("tab", { name: "npm", exact: true }).click();
  await expect(install.getByText("npm", { exact: true })).toHaveCount(1);
  const activeInstall = install.getByRole("tab", { name: "npm", exact: true });
  const activeFill = await activeInstall.evaluate((node) => getComputedStyle(node).backgroundColor);
  await activeInstall.hover();
  expect(await activeInstall.evaluate((node) => [getComputedStyle(node).backgroundColor, getComputedStyle(node).transitionProperty, getComputedStyle(node).animationName])).toEqual([activeFill, "none", "none"]);
  const hoverInstall = install.getByRole("tab", { name: "Homebrew" });
  const restingFill = await hoverInstall.evaluate((node) => getComputedStyle(node).backgroundColor);
  await hoverInstall.hover();
  expect(await hoverInstall.evaluate((node) => getComputedStyle(node).backgroundColor)).not.toBe(restingFill);
  await expect(install.locator("code")).toHaveText("npm install --global @iwonz/courier");
  await install.getByRole("button", { name: "Copy command" }).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe("npm install --global @iwonz/courier");

  const cli = page.locator("#cli");
  const compatibility = cli.getByRole("checkbox", { name: "Compatible with selected command" });
  await expect(compatibility).toBeDisabled();
  for (const selector of ["[data-courier-cli-columns]", "[data-courier-cli-columns] > section"]) {
    expect(await cli.locator(selector).evaluateAll((nodes) => nodes.every((node) => {
      const style = getComputedStyle(node);
      return style.borderTopWidth === "0px" && style.borderRightWidth === "0px" && style.borderBottomWidth === "0px" && style.borderLeftWidth === "0px";
    }))).toBe(true);
  }
  for (const selector of ["[data-courier-cli-command-header]", "[data-courier-cli-options-header]", "[data-courier-cli-command]", "[data-builder-flag]", "[data-courier-cli-readout]"]) {
    expect(await cli.locator(selector).evaluateAll((nodes) => nodes.every((node) => {
      const style = getComputedStyle(node);
      return style.borderTopWidth === "0px" && style.borderRightWidth === "0px" && style.borderBottomWidth === "0px" && style.borderLeftWidth === "0px";
    }))).toBe(true);
  }
  await cli.getByRole("button", { name: "courier from <source> to <destination> [options]" }).click();
  await expect(compatibility).toBeEnabled();
  await cli.getByRole("textbox", { name: "source" }).fill("folder one");
  await cli.getByRole("textbox", { name: "destination" }).fill("host:/srv/it's");
  await cli.locator('[data-builder-flag="extract"]').getByRole("checkbox").check();
  const exclude = cli.locator('[data-builder-flag="exclude"]');
  await exclude.getByRole("textbox", { name: "--exclude <pattern> 1" }).fill("*.tmp");
  await exclude.getByRole("button", { name: "Add value" }).click();
  await exclude.getByRole("textbox", { name: "--exclude <pattern> 2" }).fill("old files/*");
  const posix = "courier from 'folder one' to 'host:/srv/it'\"'\"'s' --extract --exclude '*.tmp' --exclude 'old files/*'";
  await expect(cli.locator("[data-courier-cli-readout] code")).toHaveText(posix);
  await cli.locator("[data-courier-cli-readout]").getByRole("button", { name: "Copy command" }).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe(posix);
  await cli.getByRole("button", { name: "PowerShell", exact: true }).click();
  const powershell = "courier from 'folder one' to 'host:/srv/it''s' --extract --exclude '*.tmp' --exclude 'old files/*'";
  await expect(cli.locator("[data-courier-cli-readout] code")).toHaveText(powershell);
  await cli.locator("[data-courier-cli-readout]").getByRole("button", { name: "Copy command" }).click();
  await expect.poll(() => page.evaluate(() => navigator.clipboard.readText())).toBe(powershell);

  await exerciseThemes(page);
  await selectRussian(page);
  await expect(page.getByRole("heading", { level: 1 })).toHaveText("Отсюда — куда угодно.");
  await page.reload();
  await expect(page.getByRole("heading", { level: 1 })).toHaveText("Отсюда — куда угодно.");

  for (const theme of ["light", "dark"] as const) {
    await page.evaluate((preference) => localStorage.setItem("courier.theme", preference), theme);
    for (const viewport of [{ width: 320, height: 568 }, { width: 390, height: 844 }, { width: 1024, height: 600 }, { width: 1440, height: 900 }]) {
      await page.setViewportSize(viewport);
      await page.reload();
      await expect(page.locator("html")).toHaveAttribute("data-courier-theme", theme);
      expect(await page.evaluate(() => {
        const root = getComputedStyle(document.documentElement);
        const body = getComputedStyle(document.body);
        const heading = getComputedStyle(document.querySelector("h1")!);
        return {
          background: root.getPropertyValue("--background").trim(),
          surface: root.getPropertyValue("--card").trim(),
          primary: root.getPropertyValue("--primary").trim(),
          actionFill: root.getPropertyValue("--terminal-fill-action").trim(),
          selectionFill: root.getPropertyValue("--terminal-fill-selection").trim(),
          destructive: root.getPropertyValue("--destructive").trim(),
          bodyFont: body.fontFamily,
          headingFont: heading.fontFamily,
        };
      })).toEqual(theme === "dark" ? expect.objectContaining({ background: "#000000", surface: "#0D1015", primary: "#71FFF6", actionFill: "#71FFF6", selectionFill: "#FAD14F", destructive: "#C94A55", bodyFont: expect.stringContaining("Overpass Mono"), headingFont: expect.stringContaining("Pixelify Sans") }) : expect.objectContaining({ background: "#FFFFFF", surface: "rgba(0, 0, 0, .04)", primary: "#006B67", actionFill: "#71FFF6", selectionFill: "#FAD14F", destructive: "#C94A55", bodyFont: expect.stringContaining("Overpass Mono"), headingFont: expect.stringContaining("Pixelify Sans") }));
      expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
      expect(await page.locator("main > section").evaluateAll((sections) => sections.every((section) => getComputedStyle(section).minHeight === "0px"))).toBe(true);
    }
  }
  expect(externalRequests).toEqual([]);
});

test("touch landing keeps Relay static and avoids horizontal overflow", async ({ browser }) => {
  const context = await browser.newContext({ viewport: { width: 320, height: 568 }, isMobile: true, hasTouch: true });
  const page = await context.newPage();
  await page.goto(landingURL);
  const mascot = page.locator('#route img[width="512"][height="512"]');
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
  await expect(page.locator('img[src*="courier-relay-pixel-mark-v3"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-pixel-delivery-v3"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-pixel-route-v3"], img[src*="courier-relay-pixel-admin-v3"], img[src*="courier-relay-pixel-neutral-v3"]')).toHaveCount(0);
  expect(await page.locator("[data-courier-auth-region]").evaluate((node) => {
    const style = getComputedStyle(node);
    return [style.backgroundImage, style.borderRadius, style.boxShadow, style.backdropFilter];
  })).toEqual(["none", "0px", "none", "none"]);
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
  await expect(page.locator('img[src*="courier-relay-pixel-mark-v3"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-pixel-admin-v3"]')).toHaveCount(1);
  await expect(page.locator('img[src*="courier-relay-pixel-route-v3"], img[src*="courier-relay-pixel-delivery-v3"], img[src*="courier-relay-pixel-neutral-v3"]')).toHaveCount(0);
  await expect(page.locator("[data-courier-metrics]")).toHaveCount(1);
  for (const selector of ["[data-courier-metrics]", "[data-courier-admin-workspace]"]) {
    expect(await page.locator(selector).evaluate((node) => {
      const style = getComputedStyle(node);
      return [style.backgroundImage, style.borderRadius, style.boxShadow, style.backdropFilter];
    })).toEqual(["none", "0px", "none", "none"]);
  }
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
