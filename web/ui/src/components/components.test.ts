import { afterEach, beforeAll, describe, expect, it, vi } from "vitest";
import { CourierButton } from "./button";
import { brandIconNames, CourierBrandIcon, resolveBrandIcon } from "./brand-icon";
import { CourierBrand, CourierMascot, CourierRoute, CourierStatus } from "./brand";
import { CourierCheckbox } from "./checkbox";
import { CourierIconLink } from "./icon-link";
import { CourierLocaleSelector } from "./locale-selector";
import { CourierPanel } from "./panel";
import { CourierProgress, progressRatio } from "./progress";
import { CourierThemeSelector } from "./theme-selector";
import { CourierScene } from "./scene";
import { CourierCommandReadout, CourierWorkbench } from "./workbench";
import { defineCourierElements, type ElementRegistry } from "../define";
import { cubicBezierPath } from "../geometry";
import { CourierIcon, iconNames, resolveIcon } from "../icons";
import { adminOperationsMobileSource, adminOperationsSource } from "../admin-scenes";
import { deliveryAccessMobileSource, deliveryAccessSource } from "../delivery-scenes";
import {
  landingInstallMobileSource,
  landingInstallSource,
  landingReferenceMobileSource,
  landingReferenceSource,
  landingRouteMobileSource,
  landingRouteSource,
} from "../landing-scenes";
import { vectorMascotSource } from "../identity-assets";
import { BrowserPreferenceController, browserPreferenceController, nextPreference, preferenceLocaleOrder, preferenceThemeOrder, resetBrowserPreferenceController } from "../preferences";

beforeAll(() => defineCourierElements());

afterEach(() => {
  document.body.replaceChildren();
  document.documentElement.removeAttribute("data-courier-theme");
  document.documentElement.removeAttribute("data-courier-theme-preference");
  localStorage.clear();
  resetBrowserPreferenceController();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe("element registry", () => {
  it("defines each element once", () => {
    const values = new Map<string, CustomElementConstructor>();
    const registry: ElementRegistry = {
      define: (name, constructor) => { values.set(name, constructor); },
      get: (name) => values.get(name),
    };
    defineCourierElements(registry);
    expect([...values.keys()].sort()).toEqual([
      "courier-brand", "courier-brand-icon", "courier-button", "courier-checkbox", "courier-command-readout", "courier-icon", "courier-icon-link", "courier-locale-selector", "courier-mascot", "courier-panel", "courier-progress", "courier-route", "courier-scene", "courier-status", "courier-theme-selector", "courier-workbench",
    ]);
    defineCourierElements(registry);
    expect(values.size).toBe(16);
  });
});

describe("shared components", () => {
  it("renders buttons and panels", async () => {
    const button = document.createElement("courier-button") as CourierButton;
    button.variant = "primary";
    button.disabled = true;
    button.type = "submit";
    button.textContent = "Send";
    document.body.append(button);
    await button.updateComplete;
    const nativeButton = button.shadowRoot?.querySelector("button");
    expect(nativeButton?.disabled).toBe(true);
    expect(nativeButton?.type).toBe("submit");
    expect(button.getAttribute("variant")).toBe("primary");

    const form = document.createElement("form");
    const submit = document.createElement("courier-button") as CourierButton;
    submit.type = "submit";
    const requestSubmit = vi.spyOn(form, "requestSubmit").mockImplementation(() => undefined);
    form.append(submit);
    document.body.append(form);
    await submit.updateComplete;
    (submit.shadowRoot?.querySelector("button") as HTMLButtonElement).click();
    expect(requestSubmit).toHaveBeenCalledOnce();
    const detachedSubmit = new CourierButton();
    detachedSubmit.type = "submit";
    expect(() => (detachedSubmit as unknown as { forwardSubmit(): void }).forwardSubmit()).not.toThrow();
    expect(() => (new CourierButton() as unknown as { forwardSubmit(): void }).forwardSubmit()).not.toThrow();

    const panel = document.createElement("courier-panel") as CourierPanel;
    panel.heading = "Delivery";
    panel.textContent = "Cargo";
    document.body.append(panel);
    await panel.updateComplete;
    expect(panel.shadowRoot?.querySelector("h2")?.textContent).toBe("Delivery");
    expect(panel.shadowRoot?.querySelector("section")?.getAttribute("aria-labelledby")).toBe("courier-panel-heading");

    panel.heading = "";
    await panel.updateComplete;
    expect(panel.shadowRoot?.querySelector("h2")).toBeNull();
    expect(panel.shadowRoot?.querySelector("section")?.hasAttribute("aria-labelledby")).toBe(false);
  });

  it("renders the shared identity, mascot, route, and status grammar", async () => {
    const brand = document.createElement("courier-brand") as CourierBrand;
    document.body.append(brand);
    await brand.updateComplete;
    expect(brand.shadowRoot?.textContent).toContain("Courier");
    expect(brand.shadowRoot?.textContent).not.toContain("Operations");
    expect(brand.shadowRoot?.querySelector("img")?.src).toContain("data:image/svg+xml");
    expect(brand.shadowRoot?.querySelector("img")?.alt).toBe("");

    const mascot = document.createElement("courier-mascot") as CourierMascot;
    mascot.alt = "Vector";
    mascot.eager = true;
    mascot.mobileSource = landingRouteMobileSource;
    mascot.source = vectorMascotSource;
    document.body.append(mascot);
    await mascot.updateComplete;
    expect(mascot.shadowRoot?.querySelector("img")?.alt).toBe("Vector");
    expect(mascot.shadowRoot?.querySelector("img")?.src).toContain("vector-mascot");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("loading")).toBe("eager");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("fetchpriority")).toBe("high");
    expect(mascot.shadowRoot?.querySelector("source")?.srcset).toContain("landing-route-mobile");
    mascot.mobileSource = "";
    mascot.eager = false;
    await mascot.updateComplete;
    expect(mascot.shadowRoot?.querySelector("source")).toBeNull();
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("loading")).toBe("lazy");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("fetchpriority")).toBe("auto");

    expect([deliveryAccessMobileSource, deliveryAccessSource].every((source) => source.includes("delivery-access-"))).toBe(true);
    expect([adminOperationsMobileSource, adminOperationsSource].every((source) => source.includes("admin-operations-"))).toBe(true);
    expect([
      landingInstallMobileSource, landingInstallSource, landingReferenceMobileSource,
      landingReferenceSource, landingRouteMobileSource, landingRouteSource,
    ].every((source) => source.includes("landing-"))).toBe(true);

    const route = document.createElement("courier-route") as CourierRoute;
    route.source = "./data";
    route.destination = "server:/data";
    document.body.append(route);
    await route.updateComplete;
    expect(route.shadowRoot?.textContent).toContain("./data");
    expect(route.shadowRoot?.textContent).toContain("server:/data");
    expect(route.shadowRoot?.querySelector("path")?.getAttribute("d")).toContain(" C ");
    expect(route.shadowRoot?.querySelector("svg circle")).toBeNull();
    expect(route.shadowRoot?.querySelectorAll(".terminal")).toHaveLength(2);

    const status = document.createElement("courier-status") as CourierStatus;
    status.tone = "signal";
    status.textContent = "Ready";
    document.body.append(status);
    await status.updateComplete;
    expect(status.getAttribute("tone")).toBe("signal");
    expect(status.shadowRoot?.querySelector("slot")).not.toBeNull();
  });

  it("normalizes progress and renders localized or explicit labels", async () => {
    expect(progressRatio(Number.NaN, 1)).toBe(0);
    expect(progressRatio(1, Number.POSITIVE_INFINITY)).toBe(0);
    expect(progressRatio(1, 0)).toBe(0);
    expect(progressRatio(-1, 10)).toBe(0);
    expect(progressRatio(20, 10)).toBe(1);
    expect(progressRatio(5, 10)).toBe(0.5);

    const progress = document.createElement("courier-progress") as CourierProgress;
    progress.value = 5;
    progress.total = 10;
    progress.locale = "ru";
    document.body.append(progress);
    await progress.updateComplete;
    expect(progress.shadowRoot?.querySelector(".track")?.getAttribute("aria-label")).toBe("Ход доставки");
    expect(progress.shadowRoot?.querySelector("output")?.textContent).toBe("50%");
    progress.label = "Files";
    progress.value = -2;
    progress.total = -1;
    await progress.updateComplete;
    expect(progress.shadowRoot?.querySelector(".track")?.getAttribute("aria-label")).toBe("Files");
    expect(progress.shadowRoot?.querySelector(".track")?.getAttribute("aria-valuemax")).toBe("0");
    progress.value = Number.NaN;
    progress.total = Number.NaN;
    await progress.updateComplete;
    expect(progress.shadowRoot?.querySelector(".track")?.getAttribute("aria-valuenow")).toBe("0");
  });

  it("renders every icon without retaining mirror mode", async () => {
    expect(iconNames).not.toContain("mirror");
    for (const name of iconNames) {
      expect(resolveIcon(name)).toBe(name);
    }
    expect(resolveIcon("unknown")).toBe("parcel");
    const icon = document.createElement("courier-icon") as CourierIcon;
    icon.name = "unknown";
    document.body.append(icon);
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector("svg")?.getAttribute("aria-hidden")).toBe("true");
    icon.name = "shield";
    icon.label = "Verified";
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector("svg")?.getAttribute("role")).toBe("img");
    expect(icon.shadowRoot?.querySelector("svg")?.getAttribute("aria-label")).toBe("Verified");
  });

  it("renders a secure semantic icon link with shared control chrome", async () => {
    const link = document.createElement("courier-icon-link") as CourierIconLink;
    link.href = "https://github.com/iwonz/courier";
    link.label = "Courier on GitHub";
    document.body.append(link);
    await link.updateComplete;
    const anchor = link.shadowRoot?.querySelector("a");
    expect(anchor?.href).toBe("https://github.com/iwonz/courier");
    expect(anchor?.target).toBe("_blank");
    expect(anchor?.rel).toBe("noopener noreferrer");
    expect(anchor?.getAttribute("aria-label")).toBe("Courier on GitHub");
    expect((CourierIconLink.styles as { cssText: string }).cssText).toContain("--courier-control-frame-size");
    expect((CourierThemeSelector.styles as { cssText: string }).cssText).toContain("--courier-control-frame-size");
  });

  it("renders pinned monochrome brand marks and the Wget fallback glyph", async () => {
    for (const name of brandIconNames) expect(resolveBrandIcon(name)).toBe(name);
    expect(resolveBrandIcon("unknown")).toBe("linux");

    const icon = document.createElement("courier-brand-icon") as CourierBrandIcon;
    icon.name = "npm";
    document.body.append(icon);
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector(".mask")?.getAttribute("style")).toContain("url(");
    expect(icon.shadowRoot?.querySelector(".mask")?.getAttribute("aria-hidden")).toBe("true");
    icon.label = "npm";
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector(".mask")?.getAttribute("role")).toBe("img");

    icon.name = "wget";
    icon.label = "";
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector("svg")?.getAttribute("role")).toBe("presentation");
    icon.label = "GNU Wget";
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector("svg")?.getAttribute("aria-label")).toBe("GNU Wget");
  });

  it("provides styled native checkbox semantics and one composed change", async () => {
    const checkbox = document.createElement("courier-checkbox") as CourierCheckbox;
    checkbox.label = "Compatible";
    checkbox.checked = true;
    checkbox.disabled = false;
    const changes: boolean[] = [];
    checkbox.addEventListener("courier-checkbox-change", (event) => changes.push((event as CustomEvent<boolean>).detail));
    document.body.append(checkbox);
    await checkbox.updateComplete;
    const input = checkbox.shadowRoot?.querySelector("input") as HTMLInputElement;
    expect(input.checked).toBe(true);
    expect(input.disabled).toBe(false);
    input.checked = false;
    input.dispatchEvent(new Event("change"));
    await checkbox.updateComplete;
    expect(checkbox.checked).toBe(false);
    expect(changes).toEqual([false]);
    checkbox.disabled = true;
    await checkbox.updateComplete;
    expect((checkbox.shadowRoot?.querySelector("input") as HTMLInputElement).disabled).toBe(true);
  });

  it("generates stable Bezier geometry in either direction and rejects invalid input", () => {
    expect(cubicBezierPath({ x: 0, y: 10 }, { x: 100, y: 20 })).toBe("M 0 10 C 42 10, 58 20, 100 20");
    expect(cubicBezierPath({ x: 100, y: 20 }, { x: 0, y: 10 })).toBe("M 100 20 C 58 20, 42 10, 0 10");
    expect(cubicBezierPath({ x: 0, y: 0 }, { x: 10, y: 0 })).toContain("C 32 0, -22 0");
    expect(() => cubicBezierPath({ x: Number.NaN, y: 0 }, { x: 1, y: 1 })).toThrow("finite");
  });

  it("renders one stationary scene without pointer or refraction machinery", async () => {
    const request = vi.fn();
    vi.stubGlobal("requestAnimationFrame", request);
    const scene = document.createElement("courier-scene") as CourierScene;
    scene.source = landingRouteSource;
    scene.mobileSource = landingRouteMobileSource;
    scene.eager = true;
    document.body.append(scene);
    await scene.updateComplete;
    expect(scene.shadowRoot?.querySelectorAll("courier-mascot")).toHaveLength(1);
    expect(scene.shadowRoot?.querySelector(".refraction, .glow")).toBeNull();
    expect((CourierScene.styles as { cssText: string }).cssText).toContain("pointer-events: none");
    scene.dispatchEvent(new MouseEvent("pointermove", { clientX: 50, clientY: 50 }));
    expect(request).not.toHaveBeenCalled();
  });

  it("activates non-eager scenes once near the viewport and cleans observers", async () => {
    const instances: { callback: IntersectionObserverCallback; observe: ReturnType<typeof vi.fn>; disconnect: ReturnType<typeof vi.fn>; options?: IntersectionObserverInit }[] = [];
    class IntersectionObserverStub {
      observe = vi.fn();
      disconnect = vi.fn();
      constructor(readonly callback: IntersectionObserverCallback, readonly options?: IntersectionObserverInit) { instances.push(this); }
    }
    vi.stubGlobal("IntersectionObserver", IntersectionObserverStub);
    const scene = document.createElement("courier-scene") as CourierScene;
    scene.source = landingRouteSource;
    document.body.append(scene);
    await scene.updateComplete;
    expect(scene.shadowRoot?.querySelector("courier-mascot")).toBeNull();
    expect(instances[0]?.options?.rootMargin).toBe("100% 0px");
    expect(instances[0]?.observe).toHaveBeenCalledWith(scene);
    instances[0]!.callback([{ isIntersecting: false } as IntersectionObserverEntry], {} as IntersectionObserver);
    await scene.updateComplete;
    expect(scene.shadowRoot?.querySelector("courier-mascot")).toBeNull();
    instances[0]!.callback([{ isIntersecting: true } as IntersectionObserverEntry], {} as IntersectionObserver);
    await scene.updateComplete;
    expect(scene.shadowRoot?.querySelector("courier-mascot.base")).not.toBeNull();
    expect(instances[0]?.disconnect).toHaveBeenCalledOnce();
    scene.remove();
    expect(instances[0]?.disconnect).toHaveBeenCalledOnce();

    const eager = new CourierScene();
    eager.eager = true;
    eager.connectedCallback();
    await eager.updateComplete;
    expect(eager.shadowRoot?.querySelector("courier-mascot.base")).not.toBeNull();
    eager.disconnectedCallback();
    eager.connectedCallback();
    eager.disconnectedCallback();

    vi.stubGlobal("IntersectionObserver", undefined);
    const fallback = document.createElement("courier-scene") as CourierScene;
    document.body.append(fallback);
    await fallback.updateComplete;
    expect(fallback.shadowRoot?.querySelector("courier-mascot.base")).not.toBeNull();

    const promoted = document.createElement("courier-scene") as CourierScene;
    promoted.eager = true;
    document.body.append(promoted);
    await promoted.updateComplete;
    expect(promoted.shadowRoot?.querySelector("courier-mascot.base")).not.toBeNull();
  });

  it("renders workbenches and immutable command readouts with reserved copy feedback", async () => {
    const workbench = document.createElement("courier-workbench") as CourierWorkbench;
    workbench.heading = "Operations";
    workbench.status = "ready";
    workbench.textContent = "Registry";
    document.body.append(workbench);
    await workbench.updateComplete;
    expect(workbench.shadowRoot?.textContent).toContain("Operations");
    expect(workbench.shadowRoot?.textContent).toContain("ready");
    workbench.status = "";
    await workbench.updateComplete;
    expect(workbench.shadowRoot?.querySelector(".status")).toBeNull();

    const readout = document.createElement("courier-command-readout") as CourierCommandReadout;
    readout.command = "courier from ./data to ./backup";
    readout.description = "Immutable route command";
    readout.heading = "Command";
    readout.clipboard = { writeText: vi.fn().mockResolvedValue(undefined) };
    document.body.append(readout);
    await readout.updateComplete;
    expect(readout.shadowRoot?.querySelector("input, textarea, [contenteditable]")).toBeNull();
    expect(readout.shadowRoot?.textContent).toContain(readout.command);
    expect(readout.shadowRoot?.textContent).not.toContain("Run");
    expect(readout.shadowRoot?.querySelector(".copy-status")?.textContent).toBe("");
    await readout.copyCommand();
    await readout.updateComplete;
    expect(readout.shadowRoot?.textContent).toContain("Copied");
    readout.clipboard = { writeText: vi.fn().mockRejectedValue(new Error("denied")) };
    await readout.copyCommand();
    await readout.updateComplete;
    expect(readout.shadowRoot?.textContent).toContain("Copy failed");
    readout.sessionKey = "next";
    await readout.updateComplete;
    expect(readout.shadowRoot?.textContent).not.toContain("Copy failed");
    readout.command = "";
    await readout.updateComplete;
    await readout.copyCommand();
    expect(readout.shadowRoot?.querySelector("button")?.disabled).toBe(true);

    const noClipboard = new CourierCommandReadout();
    noClipboard.command = "courier help";
    Object.defineProperty(navigator, "clipboard", { configurable: true, value: undefined });
    await noClipboard.copyCommand();
    expect(noClipboard.render()).toBeTruthy();
  });
});

describe("preference selectors", () => {
  it("cycles and persists theme through one semantic button", async () => {
    const listeners = new Set<() => void>();
    const media = {
      matches: true,
      addEventListener: (_type: "change", listener: () => void) => { listeners.add(listener); },
      removeEventListener: (_type: "change", listener: () => void) => { listeners.delete(listener); },
    };
    vi.stubGlobal("matchMedia", vi.fn(() => media));
    const selector = document.createElement("courier-theme-selector") as CourierThemeSelector;
    selector.locale = "ru";
    let detail = "";
    selector.addEventListener("courier-theme-change", (event) => { detail = (event as CustomEvent<string>).detail; });
    document.body.append(selector);
    await selector.updateComplete;
    expect(document.documentElement.dataset.courierTheme).toBe("dark");
    const control = selector.shadowRoot?.querySelector("button") as HTMLButtonElement;
    expect(selector.shadowRoot?.querySelector("[role=radiogroup]")).toBeNull();
    expect(control.getAttribute("aria-label")).toContain("Системная");
    expect(control.getAttribute("aria-label")).toContain("Светлая");
    control.click();
    await selector.updateComplete;
    expect(detail).toBe("light");
    expect(localStorage.getItem("courier.theme")).toBe("light");
    control.click();
    await selector.updateComplete;
    expect(detail).toBe("dark");
    control.click();
    await selector.updateComplete;
    expect(detail).toBe("system");
    selector.remove();
    resetBrowserPreferenceController();
    expect(listeners.size).toBe(0);

    const disconnected = new CourierThemeSelector();
    (disconnected as unknown as { synchronize(): void }).synchronize();
    (disconnected as unknown as { cycle(): void }).cycle();
    expect(disconnected.preference).toBe("light");
    disconnected.disconnectedCallback();
  });

  it("cycles and persists locale through one semantic button", async () => {
    localStorage.setItem("courier.locale", "ru");
    const selector = document.createElement("courier-locale-selector") as CourierLocaleSelector;
    let detail = "";
    selector.addEventListener("courier-locale-change", (event) => { detail = (event as CustomEvent<string>).detail; });
    document.body.append(selector);
    await selector.updateComplete;
    expect(selector.locale).toBe("ru");
    const control = selector.shadowRoot?.querySelector("button") as HTMLButtonElement;
    expect(control.textContent).toContain("🇷🇺");
    expect(control.getAttribute("aria-label")).toContain("Русский");
    expect(control.getAttribute("aria-label")).toContain("Английский");
    control.click();
    await selector.updateComplete;
    expect(selector.locale).toBe("en");
    expect(detail).toBe("en");
    expect(localStorage.getItem("courier.locale")).toBe("en");
    control.click();
    await selector.updateComplete;
    expect(detail).toBe("ru");

    const disconnected = new CourierLocaleSelector();
    (disconnected as unknown as { synchronize(): void }).synchronize();
    (disconnected as unknown as { cycle(): void }).cycle();
    expect(disconnected.locale).toBe("ru");
    disconnected.disconnectedCallback();
  });

  it("centralizes browser defaults, explicit overrides, cycling, and storage failures", () => {
    expect(nextPreference(preferenceThemeOrder, "system")).toBe("light");
    expect(nextPreference(preferenceThemeOrder, "dark")).toBe("system");
    expect(nextPreference(preferenceLocaleOrder, "ru")).toBe("en");
    const listeners = new Set<() => void>();
    const media = { matches: true, addEventListener: (_type: "change", listener: () => void) => listeners.add(listener), removeEventListener: (_type: "change", listener: () => void) => listeners.delete(listener) };
    const root = { dataset: {} as DOMStringMap };
    const storage = { getItem: (key: string) => key === "courier.theme" ? "light" : "ru", setItem: vi.fn() };
    const controller = new BrowserPreferenceController(root, storage, media, ["en-US"]);
    const changes = vi.fn();
    controller.addEventListener("change", changes);
    expect(controller.theme).toBe("light");
    expect(controller.locale).toBe("ru");
    expect(controller.cycleTheme()).toBe("dark");
    expect(controller.setLocale("en")).toBe("en");
    expect(controller.cycleLocale()).toBe("ru");
    expect(changes).toHaveBeenCalledTimes(3);
    controller.destroy();
    expect(listeners.size).toBe(0);

    const restricted = { getItem: () => { throw new Error("blocked"); }, setItem: () => { throw new Error("blocked"); } };
    expect(() => new BrowserPreferenceController(root, restricted, undefined, ["ru-RU"], "system").cycleTheme()).not.toThrow();

    resetBrowserPreferenceController();
    const storageDescriptor = Object.getOwnPropertyDescriptor(globalThis, "localStorage");
    const navigatorDescriptor = Object.getOwnPropertyDescriptor(globalThis, "navigator");
    Object.defineProperty(globalThis, "localStorage", { configurable: true, get: () => { throw new Error("blocked"); } });
    Object.defineProperty(globalThis, "navigator", { configurable: true, value: undefined });
    expect(browserPreferenceController().theme).toBe("system");
    expect(browserPreferenceController().locale).toBe("en");
    resetBrowserPreferenceController();
    if (storageDescriptor) Object.defineProperty(globalThis, "localStorage", storageDescriptor);
    if (navigatorDescriptor) Object.defineProperty(globalThis, "navigator", navigatorDescriptor);
  });
});
