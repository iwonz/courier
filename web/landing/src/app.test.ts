import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import {
  applicableFlags,
  commandFlags,
  CourierLandingApp,
  endpointExample,
  endpointIcon,
  endpointMessage,
  expandRoutePairs,
  installs,
} from "./app";
import { contractData } from "./contract";

function rect(left: number, top: number, width: number, height: number): DOMRect {
  return { left, top, width, height, right: left + width, bottom: top + height, x: left, y: top, toJSON: () => ({}) };
}

beforeEach(() => {
  localStorage.clear();
  vi.stubGlobal("matchMedia", vi.fn(() => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })));
});

afterEach(() => {
  document.body.replaceChildren();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

it("projects routes, endpoint vocabulary, exact route flags, and command compatibility", () => {
  expect(document.head.querySelector<HTMLLinkElement>('link[rel="icon"]')?.href).toContain("relay-mark");
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

  expect(commandFlags("", true, contractData.commands, contractData.flags)).toHaveLength(contractData.flags.length);
  expect(commandFlags("from", false, contractData.commands, contractData.flags)).toHaveLength(contractData.flags.length);
  expect(commandFlags("ui-start", true, contractData.commands, contractData.flags).map((flag) => flag.name)).toEqual(["listen", "background"]);
  expect(commandFlags("servers-stop", true, contractData.commands, contractData.flags).map((flag) => flag.name)).toEqual(["all"]);
  expect(commandFlags("servers", true, contractData.commands, contractData.flags)).toEqual([]);
  expect(commandFlags("unknown", true, contractData.commands, contractData.flags)).toEqual([]);
  expect(installs.map((install) => install.icon)).toEqual(["curl", "wget", "powershell", "npm", "npx", "yarn", "pnpm", "homebrew", "scoop"]);
});

it("renders stable non-interactive scenes and activation-only route, install, and CLI controls", async () => {
  const resizeInstances: { callback: ResizeObserverCallback; observe: ReturnType<typeof vi.fn>; disconnect: ReturnType<typeof vi.fn> }[] = [];
  const sectionInstances: { callback: IntersectionObserverCallback; observe: ReturnType<typeof vi.fn>; disconnect: ReturnType<typeof vi.fn> }[] = [];
  class ResizeObserverStub {
    observe = vi.fn();
    disconnect = vi.fn();
    constructor(readonly callback: ResizeObserverCallback) { resizeInstances.push(this); }
  }
  class IntersectionObserverStub {
    observe = vi.fn();
    disconnect = vi.fn();
    constructor(readonly callback: IntersectionObserverCallback) { sectionInstances.push(this); }
  }
  vi.stubGlobal("ResizeObserver", ResizeObserverStub);
  vi.stubGlobal("IntersectionObserver", IntersectionObserverStub);
  vi.spyOn(HTMLElement.prototype, "getBoundingClientRect").mockImplementation(function (this: HTMLElement) {
    if (this.classList.contains("masthead-wrap")) return rect(0, 0, 1440, 78);
    if (this.classList.contains("route-controls")) return rect(100, 180, 700, 320);
    if (this.closest(".source-endpoints") && this.matches(".selected")) return rect(120, 250, 180, 46);
    if (this.closest(".destination-endpoints") && this.matches(".selected")) return rect(600, 340, 180, 46);
    return rect(0, 0, 0, 0);
  });

  const element = new CourierLandingApp();
  document.body.append(element);
  await element.updateComplete;
  await element.updateComplete;
  const root = element.shadowRoot!;
  const text = root.textContent ?? "";

  expect(text).toContain("Move files. Keep control.");
  expect(text).toContain("courier from <source> to <destination>");
  expect(text).toContain("path-to-path");
  expect(text).toContain("Source");
  expect(text).toContain("Destination");
  expect(text).not.toContain("Remote");
  expect([...root.querySelectorAll("section")].map((section) => section.id)).toEqual(["hero", "routes", "install", "cli"]);
  expect(root.querySelectorAll("courier-scene")).toHaveLength(4);
  expect(root.querySelector("#hero button")).toBeNull();
  expect(root.querySelector("#hero [aria-pressed]")).toBeNull();
  expect(root.querySelector(".hero-route path")?.getAttribute("d")).toContain(" C ");
  expect(element.style.getPropertyValue("--masthead-height")).toBe("78px");
  expect(root.querySelector(".route-connector path")?.getAttribute("d")).toContain(" C ");
  expect(resizeInstances[0]?.observe).toHaveBeenCalledTimes(2);
  expect(sectionInstances[0]?.observe).toHaveBeenCalledTimes(4);

  const scenes = [...root.querySelectorAll("courier-scene")];
  await Promise.all(scenes.map((scene) => scene.updateComplete));
  const baseSources = scenes.map((scene) => scene.shadowRoot?.querySelector(".base") as HTMLElement);
  await Promise.all(baseSources.map((mascot) => (mascot as unknown as { updateComplete: Promise<unknown> }).updateComplete));
  expect(baseSources.map((mascot) => mascot.shadowRoot?.querySelector("img")?.src)).toEqual(expect.arrayContaining([
    expect.stringContaining("relay-hero-wide"), expect.stringContaining("relay-install-wide"), expect.stringContaining("relay-routing-wide"), expect.stringContaining("relay-cli-wide"),
  ]));

  const routeSection = root.querySelector("#routes") as HTMLElement;
  routeSection.scrollIntoView = vi.fn();
  const pushState = vi.spyOn(history, "pushState");
  (root.querySelector('nav a[href="#routes"]') as HTMLAnchorElement).click();
  expect(routeSection.scrollIntoView).toHaveBeenCalledWith({ block: "start" });
  expect(pushState).toHaveBeenCalledWith(null, "", "#routes");

  const routeCommand = () => root.querySelector(".command-shape")?.textContent ?? "";
  const source = (name: string) => root.querySelector(`.source-endpoints button[data-endpoint="${name}"]`) as HTMLButtonElement;
  const destination = (name: string) => root.querySelector(`.destination-endpoints button[data-endpoint="${name}"]`) as HTMLButtonElement;
  source("web").dispatchEvent(new Event("pointerenter"));
  source("web").dispatchEvent(new FocusEvent("focus"));
  await element.updateComplete;
  expect(routeCommand()).toContain("./project to relay@host:/srv/destination/");
  source("web").click();
  await element.updateComplete;
  expect(routeCommand()).toContain("web:// to relay@host:/srv/destination/");
  expect(destination("web").disabled).toBe(true);
  expect(destination("web").getAttribute("aria-disabled")).toBe("true");
  destination("web").click();
  await element.updateComplete;
  expect(routeCommand()).toContain("web:// to relay@host:/srv/destination/");
  destination("local").click();
  await element.updateComplete;
  expect(routeCommand()).toContain("web:// to ./backup/");
  source("local").click();
  await element.updateComplete;
  destination("web").click();
  await element.updateComplete;
  expect(routeCommand()).toContain("./project to web://");
  source("webhook").click();
  await element.updateComplete;
  expect(routeCommand()).toContain("webhook:// to ./backup/");

  const installCommand = () => root.querySelector(".install-readout code")?.textContent ?? "";
  const install = (name: string) => root.querySelector(`.install-channel[data-channel="${name}"]`) as HTMLButtonElement;
  install("npm").dispatchEvent(new Event("pointerenter"));
  install("npm").dispatchEvent(new FocusEvent("focus"));
  await element.updateComplete;
  expect(installCommand()).toBe("curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh");
  install("Homebrew").click();
  await element.updateComplete;
  expect(installCommand()).toContain("brew tap iwonz/courier");
  (element as unknown as { activeInstall: string }).activeInstall = "missing";
  await element.updateComplete;
  expect(installCommand()).toBe("curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh");

  const checkbox = root.querySelector("courier-checkbox")!;
  await checkbox.updateComplete;
  expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).checked).toBe(true);
  expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).disabled).toBe(true);
  expect(root.querySelectorAll(".option-row")).toHaveLength(contractData.flags.length);

  const command = (name: string) => root.querySelector(`.command-row[data-command="${name}"]`) as HTMLButtonElement;
  command("ui-start").click();
  await element.updateComplete;
  expect([...root.querySelectorAll(".option-row code")].map((node) => node.textContent)).toEqual(["--listen <host:port>", "--background"]);
  expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).disabled).toBe(false);
  checkbox.dispatchEvent(new CustomEvent("courier-checkbox-change", { detail: false }));
  await element.updateComplete;
  expect(root.querySelectorAll(".option-row")).toHaveLength(contractData.flags.length);
  command("ui-start").click();
  await element.updateComplete;
  expect(command("ui-start").getAttribute("aria-pressed")).toBe("false");
  expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).checked).toBe(true);
  expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).disabled).toBe(true);
  command("servers").click();
  await element.updateComplete;
  expect(root.querySelector(".empty-state")?.textContent).toBe("This command has no options.");
  command("servers-stop").click();
  await element.updateComplete;
  expect(root.querySelector(".option-row code")?.textContent).toBe("--all");
  command("from").click();
  await element.updateComplete;
  expect(root.querySelectorAll(".option-row")).toHaveLength(17);

  sectionInstances[0]!.callback([
    { target: root.querySelector("#hero")!, isIntersecting: true, intersectionRatio: 0.1 },
    { target: root.querySelector("#routes")!, isIntersecting: true, intersectionRatio: 0.8 },
  ] as IntersectionObserverEntry[], {} as IntersectionObserver);
  await element.updateComplete;
  expect(root.querySelector('nav a[href="#routes"]')?.getAttribute("aria-current")).toBe("page");
  sectionInstances[0]!.callback([
    { target: root.querySelector("#routes")!, isIntersecting: false, intersectionRatio: 0 },
    { target: root.querySelector("#hero")!, isIntersecting: true, intersectionRatio: 0.9 },
  ] as IntersectionObserverEntry[], {} as IntersectionObserver);
  await element.updateComplete;
  expect(root.querySelector("nav [aria-current]")).toBeNull();
  resizeInstances[0]!.callback([], {} as ResizeObserver);
  await element.updateComplete;

  expect((element as unknown as { displayEndpoint(name: string): string }).displayEndpoint("future")).toBe("future");
  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(root.textContent).toContain("Переносите файлы. Сохраняйте контроль.");
  expect(root.textContent).toContain("Команда");
  expect(root.textContent).toContain("Source");
  expect(root.textContent).toContain("Destination");
  expect(root.textContent).toContain("SSH");
  expect(root.textContent).not.toContain("Remote");

  element.remove();
  expect(resizeInstances[0]?.disconnect).toHaveBeenCalledOnce();
  expect(sectionInstances[0]?.disconnect).toHaveBeenCalledOnce();
});

it("covers zero-layout and detached observer-free lifecycle safely", async () => {
  vi.stubGlobal("ResizeObserver", undefined);
  vi.stubGlobal("IntersectionObserver", undefined);
  const element = new CourierLandingApp();
  const internals = element as unknown as {
    measureHeader(): void;
    measureRouteConnector(): void;
    observeSections(entries: readonly IntersectionObserverEntry[]): void;
    navigate(event: MouseEvent): void;
  };
  document.body.append(element);
  await element.updateComplete;
  internals.measureHeader();
  internals.measureRouteConnector();
  element.shadowRoot?.querySelector(".masthead-wrap")?.remove();
  element.shadowRoot?.querySelector(".route-controls")?.remove();
  internals.measureHeader();
  internals.measureRouteConnector();
  internals.observeSections([]);
  const missing = document.createElement("a");
  missing.href = "#missing";
  internals.navigate({ preventDefault: vi.fn(), currentTarget: missing } as unknown as MouseEvent);
  element.remove();

  new CourierLandingApp().disconnectedCallback();
  vi.resetModules();
  await import("./main");
  expect(customElements.get("courier-landing-app")).toBe(CourierLandingApp);
});
