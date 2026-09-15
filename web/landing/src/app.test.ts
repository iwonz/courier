import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import {
  applicableFlags,
  CourierLandingApp,
  endpointExample,
  endpointIcon,
  endpointMessage,
  expandRoutePairs,
} from "./app";
import { contractData } from "./contract";

beforeEach(() => {
  localStorage.clear();
  vi.stubGlobal("matchMedia", () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() }));
});

afterEach(() => {
  document.body.replaceChildren();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

it("projects route pairs, endpoint presentation, and exact flag applicability from the contract", () => {
  const pairs = expandRoutePairs(contractData.routes);
  expect(pairs).toContainEqual(expect.objectContaining({ source: "local", destination: "ssh", routeName: "path-to-path" }));
  expect(pairs).toContainEqual(expect.objectContaining({ source: "webhook", destination: "local", routeName: "webhook-to-path" }));
  expect(endpointExample("local", "source")).toBe("./project");
  expect(endpointExample("ssh", "destination")).toBe("relay@host:/srv/destination/");
  expect(endpointExample("future", "source")).toBe("future");
  expect(endpointMessage("ssh")).toBe("endpointSsh");
  expect(endpointMessage("future")).toBeUndefined();
  expect(endpointIcon("webhook")).toBe("upload");
  expect(endpointIcon("future")).toBe("route");

  const localPair = pairs.find((pair) => pair.source === "local" && pair.destination === "local")!;
  const localFlags = applicableFlags(localPair, contractData.flags).map((flag) => flag.name);
  expect(localFlags).toContain("archive");
  expect(localFlags).not.toContain("upload-rate");
  const sshPair = pairs.find((pair) => pair.source === "ssh" && pair.destination === "ssh")!;
  expect(applicableFlags(sshPair, contractData.flags).map((flag) => flag.name)).toEqual(expect.arrayContaining(["upload-rate", "download-rate"]));
});

it("renders four immersive sections and operates hero, route, and installation interactions", async () => {
  const element = new CourierLandingApp();
  document.body.append(element);
  await element.updateComplete;
  const root = element.shadowRoot!;
  const text = root.textContent ?? "";

  expect(text).toContain("Move files. Keep control.");
  expect(text).toContain("courier from <source> to <destination>");
  expect(text).toContain("path-to-path");
  expect(text).toContain("--archive");
  expect(text).not.toContain("Open source · MIT · self-contained");
  expect([...root.querySelectorAll("section")].map((section) => section.id)).toEqual(["hero", "routes", "install", "cli"]);
  expect(root.querySelectorAll(".slide")).toHaveLength(4);
  expect(root.querySelector(".facts")).toBeNull();
  expect(root.querySelector("#safety")).toBeNull();
  expect(root.querySelector("#examples")).toBeNull();
  expect(root.querySelector("#docs")).toBeNull();
  expect(root.querySelector("footer")).toBeNull();
  expect(root.querySelector(".quick-command")).toBeNull();
  expect(root.querySelector(".actions")).toBeNull();
  expect(root.querySelector(".reference-link")).toBeNull();
  expect(text).not.toContain("Open the complete CLI contract");
  expect(root.querySelector(".github-link")?.getAttribute("href")).toBe("https://github.com/iwonz/courier");
  expect([...root.querySelectorAll("nav a")].map((link) => link.getAttribute("href"))).toEqual(["#routes", "#install", "#cli"]);
  expect(root.querySelectorAll(".install-channel")).toHaveLength(9);

  const routeSection = root.querySelector("#routes") as HTMLElement;
  routeSection.scrollIntoView = vi.fn();
  const pushState = vi.spyOn(history, "pushState");
  (root.querySelector('nav a[href="#routes"]') as HTMLAnchorElement).click();
  expect(routeSection.scrollIntoView).toHaveBeenCalledWith({ block: "start" });
  expect(pushState).toHaveBeenCalledWith(null, "", "#routes");

  const illustrations = [...root.querySelectorAll("courier-mascot")];
  await Promise.all(illustrations.map((illustration) => illustration.updateComplete));
  expect(illustrations).toHaveLength(4);
  expect(illustrations.map((illustration) => illustration.shadowRoot?.querySelector("img")?.src)).toEqual(expect.arrayContaining([
    expect.stringContaining("relay-hero-wide"), expect.stringContaining("relay-install-wide"), expect.stringContaining("relay-routing-wide"), expect.stringContaining("relay-cli-wide"),
  ]));
  expect(illustrations.map((illustration) => illustration.shadowRoot?.querySelector("source")?.srcset)).toEqual(expect.arrayContaining([
    expect.stringContaining("relay-hero-mobile"), expect.stringContaining("relay-install-mobile"), expect.stringContaining("relay-routing-mobile"), expect.stringContaining("relay-cli-mobile"),
  ]));

  const themeSelector = root.querySelector("courier-theme-selector")!;
  await themeSelector.updateComplete;
  const themeControl = themeSelector.shadowRoot?.querySelector("courier-segmented-control")!;
  await themeControl.updateComplete;
  expect(themeControl.shadowRoot?.querySelector("legend")?.classList.contains("sr-only")).toBe(true);
  expect(themeControl.shadowRoot?.querySelector("button span")).toBeNull();

  const hero = root.querySelector(".hero-visual") as HTMLButtonElement;
  hero.getBoundingClientRect = () => ({ left: 10, top: 20, width: 100, height: 200, right: 110, bottom: 220, x: 10, y: 20, toJSON: () => ({}) });
  hero.dispatchEvent(new MouseEvent("pointermove", { clientX: 85, clientY: 70 }));
  expect(hero.style.getPropertyValue("--spot-x")).toBe("75%");
  expect(hero.style.getPropertyValue("--spot-y")).toBe("25%");
  expect(hero.style.getPropertyValue("--relay-x")).toBe("6px");
  expect(hero.style.getPropertyValue("--relay-y")).toBe("-4px");
  hero.dispatchEvent(new MouseEvent("pointerleave"));
  expect(hero.style.getPropertyValue("--spot-x")).toBe("68%");
  expect(hero.style.getPropertyValue("--relay-x")).toBe("0px");
  hero.click();
  await element.updateComplete;
  expect(hero.getAttribute("aria-pressed")).toBe("true");
  hero.dispatchEvent(new KeyboardEvent("keydown", { key: "Tab", bubbles: true }));
  await element.updateComplete;
  expect(hero.getAttribute("aria-pressed")).toBe("true");
  hero.dispatchEvent(new KeyboardEvent("keydown", { key: "Enter", bubbles: true, cancelable: true }));
  await element.updateComplete;
  expect(hero.getAttribute("aria-pressed")).toBe("false");
  hero.dispatchEvent(new KeyboardEvent("keydown", { key: " ", bubbles: true, cancelable: true }));
  await element.updateComplete;
  expect(hero.getAttribute("aria-pressed")).toBe("true");

  routeSection.getBoundingClientRect = () => ({ left: 0, top: 0, width: 200, height: 100, right: 200, bottom: 100, x: 0, y: 0, toJSON: () => ({}) });
  routeSection.dispatchEvent(new MouseEvent("pointermove", { clientX: 50, clientY: 75 }));
  expect(routeSection.style.getPropertyValue("--scene-spot-x")).toBe("25%");
  expect(routeSection.style.getPropertyValue("--scene-spot-y")).toBe("75%");
  expect(routeSection.style.getPropertyValue("--scene-x")).toBe("4px");
  expect(routeSection.style.getPropertyValue("--scene-y")).toBe("-3px");
  routeSection.dispatchEvent(new MouseEvent("pointerleave"));
  expect(routeSection.style.getPropertyValue("--scene-spot-x")).toBe("72%");
  expect(routeSection.style.getPropertyValue("--scene-y")).toBe("0px");

  const installCommand = () => root.querySelector(".install-readout code")?.textContent ?? "";
  const installChannel = (name: string) => root.querySelector(`.install-channel[data-channel="${name}"]`) as HTMLButtonElement;
  installChannel("Homebrew").dispatchEvent(new Event("pointerenter"));
  await element.updateComplete;
  expect(installCommand()).toContain("brew tap iwonz/courier");
  installChannel("npm").dispatchEvent(new FocusEvent("focus"));
  await element.updateComplete;
  expect(installCommand()).toBe("npm install --global @iwonz/courier");
  installChannel("Scoop").click();
  await element.updateComplete;
  expect(installCommand()).toContain("scoop bucket add courier");
  (element as unknown as { activeInstall: string }).activeInstall = "missing";
  await element.updateComplete;
  expect(installCommand()).toBe("curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh");

  const command = () => root.querySelector(".command-shape")?.textContent ?? "";
  expect(command()).toContain("./project to relay@host:/srv/destination/");

  const source = (name: string) => root.querySelector(`.endpoint-group:first-child button[data-endpoint="${name}"]`) as HTMLButtonElement;
  const destination = (name: string) => root.querySelector(`.endpoint-group:last-child button[data-endpoint="${name}"]`) as HTMLButtonElement;

  source("web").dispatchEvent(new Event("pointerenter"));
  await element.updateComplete;
  expect(command()).toContain("web:// to relay@host:/srv/destination/");
  expect(destination("web").getAttribute("aria-disabled")).toBe("true");

  destination("local").dispatchEvent(new FocusEvent("focus"));
  await element.updateComplete;
  expect(command()).toContain("web:// to ./backup/");
  destination("http").click();
  await element.updateComplete;
  expect(command()).toContain("web:// to ./backup/");

  source("local").click();
  await element.updateComplete;
  destination("web").click();
  await element.updateComplete;
  expect(command()).toContain("./project to web://");
  source("webhook").dispatchEvent(new Event("pointerenter"));
  await element.updateComplete;
  expect(command()).toContain("webhook:// to ./backup/");

  expect((element as unknown as { displayEndpoint(name: string): string }).displayEndpoint("future")).toBe("future");
  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(root.textContent).toContain("Переносите файлы. Сохраняйте контроль.");
  expect(root.textContent).toContain("Команды и параметры");
  expect(root.textContent).not.toContain("Открыть полный контракт CLI");
  element.remove();
});

it("supports detached lifecycle and idempotent element definition", async () => {
  new CourierLandingApp().disconnectedCallback();
  vi.resetModules();
  await import("./main");
  expect(customElements.get("courier-landing-app")).toBe(CourierLandingApp);
});
