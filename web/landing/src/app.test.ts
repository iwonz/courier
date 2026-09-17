import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import {
  applicableFlags,
  commandFlags,
  copyText,
  CourierLandingApp,
  endpointExample,
  endpointIcon,
  endpointMessage,
  endpointPresentation,
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
  expect(document.head.querySelector<HTMLLinkElement>('link[rel="icon"]')?.href).toContain("data:image/svg+xml");
  const pairs = expandRoutePairs(contractData.routes);
  expect(pairs).toContainEqual(expect.objectContaining({ source: "local", destination: "ssh", routeName: "path-to-path" }));
  expect(pairs).toContainEqual(expect.objectContaining({ source: "webhook", destination: "local", routeName: "webhook-to-path" }));
  expect(endpointExample("local", "source")).toBe("./project");
  expect(endpointExample("ssh", "destination")).toBe("courier@host:/srv/destination/");
  expect(endpointExample("future", "source")).toBe("future");
  expect(endpointMessage("ssh", "source")).toBe("endpointRemote");
  expect(endpointMessage("http", "destination")).toBe("endpointWebHook");
  expect(endpointMessage("future", "source")).toBeUndefined();
  expect(endpointPresentation("web", "source")).toEqual(expect.objectContaining({ icon: "browser-upload" }));
  expect(endpointIcon("webhook", "source")).toBe("webhook-in");
  expect(endpointIcon("http", "destination")).toBe("webhook-out");
  expect(endpointIcon("future", "destination")).toBe("route");

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

it("copies exact commands and reports unavailable or rejected clipboard access", async () => {
  const writeText = vi.fn().mockResolvedValue(undefined);
  await expect(copyText("courier", { writeText })).resolves.toBe(true);
  expect(writeText).toHaveBeenCalledWith("courier");
  await expect(copyText("courier", { writeText: vi.fn().mockRejectedValue(new Error("denied")) })).resolves.toBe(false);
  await expect(copyText("courier", undefined)).resolves.toBe(false);
});

it("renders stable non-interactive scenes and activation-only route, install, and CLI controls", async () => {
  const resizeInstances: { callback: ResizeObserverCallback; observe: ReturnType<typeof vi.fn>; disconnect: ReturnType<typeof vi.fn> }[] = [];
  const sectionInstances: { callback: IntersectionObserverCallback; observe: ReturnType<typeof vi.fn>; disconnect: ReturnType<typeof vi.fn>; options?: IntersectionObserverInit }[] = [];
  class ResizeObserverStub {
    observe = vi.fn();
    disconnect = vi.fn();
    constructor(readonly callback: ResizeObserverCallback) { resizeInstances.push(this); }
  }
  class IntersectionObserverStub {
    observe = vi.fn();
    disconnect = vi.fn();
    constructor(readonly callback: IntersectionObserverCallback, readonly options?: IntersectionObserverInit) { sectionInstances.push(this); }
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

  expect(text).toContain("From here to anywhere.");
  expect(text).toContain("courier from <source> to <destination>");
  expect(text).toContain("path-to-path");
  expect(text).toContain("Source");
  expect(text).toContain("Destination");
  expect(text).toContain("Remote");
  expect(text).toContain("Web Hook");
  expect(text).not.toContain("A dependable CLI");
  expect(text).not.toContain("Run demo");
  expect(text).not.toContain("Replay demo");
  expect(text).not.toContain("One binary plans the route");
  expect(text).not.toContain("All connections come from the published CLI contract");
  expect(text).not.toContain("PATH-TO-PATH");
  expect(text).not.toContain("One generated reference for the public command tree");
  expect([...root.querySelectorAll("section")].map((section) => section.id)).toEqual(["route", "install", "cli"]);
  expect(root.querySelectorAll("courier-scene")).toHaveLength(3);
  expect(root.querySelector("#route h1")?.textContent).toBe("From here to anywhere.");
  expect(element.style.getPropertyValue("--masthead-height")).toBe("78px");
  expect(root.querySelector(".route-connector path")?.getAttribute("d")).toContain(" C ");
  expect(resizeInstances[0]?.observe).toHaveBeenCalledTimes(2);
  const navigationObserver = sectionInstances.find((instance) => instance.observe.mock.calls.length === 3)!;
  expect(navigationObserver.options?.rootMargin).toContain("-78px");
  const github = root.querySelector("courier-icon-link")!;
  await github.updateComplete;
  const githubAnchor = github.shadowRoot?.querySelector("a");
  expect(githubAnchor?.target).toBe("_blank");
  expect(githubAnchor?.rel).toBe("noopener noreferrer");
  expect([...root.querySelectorAll<HTMLAnchorElement>('a[href^="https://"]')].every((link) => link.target === "_blank" && link.rel === "noopener noreferrer")).toBe(true);

  const scenes = [...root.querySelectorAll("courier-scene")];
  for (const instance of sectionInstances.filter((candidate) => candidate.observe.mock.calls[0]?.[0]?.tagName === "COURIER-SCENE")) {
    instance.callback([{ target: instance.observe.mock.calls[0]![0], isIntersecting: true } as IntersectionObserverEntry], {} as IntersectionObserver);
  }
  await Promise.all(scenes.map((scene) => scene.updateComplete));
  const baseSources = scenes.map((scene) => scene.shadowRoot?.querySelector(".base") as HTMLElement);
  await Promise.all(baseSources.map((mascot) => (mascot as unknown as { updateComplete: Promise<unknown> }).updateComplete));
  expect(baseSources.map((mascot) => mascot.shadowRoot?.querySelector("img")?.src)).toEqual(expect.arrayContaining([
    expect.stringContaining("landing-route-wide"), expect.stringContaining("landing-install-wide"), expect.stringContaining("landing-reference-wide"),
  ]));

  const installSection = root.querySelector("#install") as HTMLElement;
  installSection.scrollIntoView = vi.fn();
  const pushState = vi.spyOn(history, "pushState");
  (root.querySelector('nav a[href="#install"]') as HTMLAnchorElement).click();
  expect(installSection.scrollIntoView).toHaveBeenCalledWith({ block: "start", behavior: "smooth" });
  expect(pushState).toHaveBeenCalledWith(null, "", "#install");

  const routeReadout = root.querySelector(".route-readout") as HTMLElement & { command: string };
  const routeCommand = () => routeReadout.command;
  expect(routeReadout.shadowRoot?.textContent).not.toContain("Run demo");
  const source = (name: string) => root.querySelector(`.source-endpoints button[data-endpoint="${name}"]`) as HTMLButtonElement;
  const destination = (name: string) => root.querySelector(`.destination-endpoints button[data-endpoint="${name}"]`) as HTMLButtonElement;
  source("web").dispatchEvent(new Event("pointerenter"));
  source("web").dispatchEvent(new FocusEvent("focus"));
  await element.updateComplete;
  expect(routeCommand()).toContain("./project to courier@host:/srv/destination/");
  source("web").click();
  await element.updateComplete;
  expect(routeCommand()).toContain("web:// to courier@host:/srv/destination/");
  expect(destination("web").disabled).toBe(true);
  destination("web").click();
  await element.updateComplete;
  expect(routeCommand()).toContain("web:// to courier@host:/srv/destination/");
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

  const installReadout = root.querySelector(".install-readout") as HTMLElement & { command: string; clipboard?: Pick<Clipboard, "writeText">; copyCommand(): Promise<void>; updateComplete: Promise<unknown> };
  const installCommand = () => installReadout.command;
  const install = (name: string) => root.querySelector(`.install-channel[data-channel="${name}"]`) as HTMLButtonElement;
  install("npm").dispatchEvent(new Event("pointerenter"));
  install("npm").dispatchEvent(new FocusEvent("focus"));
  await element.updateComplete;
  expect(installCommand()).toBe("curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh");
  install("Homebrew").click();
  await element.updateComplete;
  expect(installCommand()).toContain("brew tap iwonz/courier");
  const clipboard = vi.fn().mockResolvedValue(undefined);
  installReadout.clipboard = { writeText: clipboard };
  await installReadout.copyCommand();
  await installReadout.updateComplete;
  expect(installReadout.shadowRoot?.textContent).toContain("Copied");
  expect(clipboard).toHaveBeenCalledWith(installCommand());
  installReadout.clipboard = { writeText: vi.fn().mockRejectedValue(new Error("denied")) };
  await installReadout.copyCommand();
  await installReadout.updateComplete;
  expect(installReadout.shadowRoot?.textContent).toContain("Copy failed");
  install("npm").click();
  await element.updateComplete;
  await installReadout.updateComplete;
  expect(installReadout.shadowRoot?.textContent).not.toContain("Copy failed");
  (element as unknown as { activeInstall: string }).activeInstall = "missing";
  await element.updateComplete;
  expect(installCommand()).toBe("curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh");
  installReadout.clipboard = { writeText: clipboard };
  await installReadout.copyCommand();
  await installReadout.updateComplete;
  expect(installReadout.shadowRoot?.textContent).toContain("Copied");
  expect(clipboard).toHaveBeenLastCalledWith(installCommand());

  const checkbox = root.querySelector("courier-checkbox")!;
  await checkbox.updateComplete;
  expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).checked).toBe(true);
  expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).disabled).toBe(true);
  expect(root.querySelectorAll(".option-row")).toHaveLength(contractData.flags.length);

  const command = (name: string) => root.querySelector(`.command-row[data-command="${name}"]`) as HTMLButtonElement;
  const cliReadout = root.querySelector(".cli-readout") as HTMLElement & { command: string };
  expect(cliReadout.command).toBe("");
  command("ui-start").click();
  await element.updateComplete;
  expect(cliReadout.command).toBe("courier ui start [options]");
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

  navigationObserver.callback([
    { target: root.querySelector("#install")!, isIntersecting: true, intersectionRatio: 0.8 },
  ] as IntersectionObserverEntry[], {} as IntersectionObserver);
  await element.updateComplete;
  expect(root.querySelector('nav a[href="#install"]')?.getAttribute("aria-current")).toBe("page");
  navigationObserver.callback([
    { target: root.querySelector("#cli")!, isIntersecting: true, intersectionRatio: 0.8 },
  ] as IntersectionObserverEntry[], {} as IntersectionObserver);
  await element.updateComplete;
  expect(root.querySelector('nav a[href="#cli"]')?.getAttribute("aria-current")).toBe("page");
  navigationObserver.callback([
    { target: root.querySelector("#install")!, isIntersecting: false, intersectionRatio: 0 },
    { target: root.querySelector("#cli")!, isIntersecting: false, intersectionRatio: 0 },
    { target: root.querySelector("#route")!, isIntersecting: true, intersectionRatio: 0.9 },
  ] as IntersectionObserverEntry[], {} as IntersectionObserver);
  await element.updateComplete;
  expect(root.querySelector("nav [aria-current]")).toBeNull();
  resizeInstances[0]!.callback([], {} as ResizeObserver);
  globalThis.dispatchEvent(new Event("resize"));
  await element.updateComplete;

  expect((element as unknown as { displayEndpoint(name: string, side: "source" | "destination"): string }).displayEndpoint("future", "source")).toBe("future");
  expect((element as unknown as { describeEndpoint(name: string, side: "source" | "destination"): string }).describeEndpoint("future", "destination")).toBe("future");
  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(root.textContent).toContain("Отсюда — куда угодно.");
  expect((root.querySelector("courier-command-readout") as HTMLElement & { heading: string }).heading).toBe("Команда");
  expect(root.textContent).toContain("Source");
  expect(root.textContent).toContain("Destination");
  expect(root.textContent).toContain("Remote");
  expect(root.textContent).toContain("Web Hook");
  expect(root.textContent).not.toContain("Надёжный CLI");
  expect(root.textContent).not.toContain("Запустить демо");
  expect(root.textContent).not.toContain("Один бинарник планирует маршрут");
  expect(root.textContent).not.toContain("Все связи взяты из опубликованного контракта CLI");
  expect(root.textContent).not.toContain("Единый сгенерированный справочник");

  element.remove();
  expect(resizeInstances[0]?.disconnect).toHaveBeenCalledOnce();
  expect(navigationObserver.disconnect).toHaveBeenCalled();
});

it("covers zero-layout and detached observer-free lifecycle safely", async () => {
  vi.stubGlobal("ResizeObserver", undefined);
  vi.stubGlobal("IntersectionObserver", undefined);
  const element = new CourierLandingApp();
  const internals = element as unknown as {
    measureHeader(): void;
    measureRouteConnector(): void;
    observeSections(entries: readonly IntersectionObserverEntry[]): void;
    observeSectionBand(): void;
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
  internals.observeSectionBand();
  const missing = document.createElement("a");
  missing.href = "#missing";
  vi.stubGlobal("matchMedia", vi.fn(() => ({ matches: true, addEventListener: vi.fn(), removeEventListener: vi.fn() })));
  internals.navigate({ preventDefault: vi.fn(), currentTarget: missing } as unknown as MouseEvent);
  const routes = element.shadowRoot?.querySelector("#route") as HTMLElement;
  routes.scrollIntoView = vi.fn();
  const routeLink = document.createElement("a");
  routeLink.href = "#route";
  internals.navigate({ preventDefault: vi.fn(), currentTarget: routeLink } as unknown as MouseEvent);
  expect(routes.scrollIntoView).toHaveBeenCalledWith({ block: "start", behavior: "auto" });
  element.remove();

  new CourierLandingApp().disconnectedCallback();
  vi.resetModules();
  await import("./main");
  expect(customElements.get("courier-landing-app")).toBe(CourierLandingApp);
});
