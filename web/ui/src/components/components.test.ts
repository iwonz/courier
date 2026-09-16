import { afterEach, beforeAll, describe, expect, it, vi } from "vitest";
import { CourierButton } from "./button";
import { brandIconNames, CourierBrandIcon, resolveBrandIcon } from "./brand-icon";
import { CourierBrand, CourierMascot, CourierRoute, CourierStatus } from "./brand";
import { CourierCheckbox } from "./checkbox";
import { CourierIconLink } from "./icon-link";
import { CourierLocaleSelector } from "./locale-selector";
import { CourierPanel } from "./panel";
import { CourierProgress, progressRatio } from "./progress";
import { CourierSegmentedControl, nextSegmentIndex } from "./segmented-control";
import { CourierThemeSelector } from "./theme-selector";
import { CourierScene, pointerPosition, sceneAmbientPosition, smoothPointerPosition } from "./scene";
import { CourierCommandReadout, CourierTerminal } from "./terminal";
import { defineCourierElements, type ElementRegistry } from "../define";
import { cubicBezierPath } from "../geometry";
import { CourierIcon, iconNames, resolveIcon } from "../icons";
import { relayOperationsMobileSource, relayOperationsSource } from "../relay-admin";
import { relayAccessMobileSource, relayAccessSource } from "../relay-delivery";
import {
  relayCliMobileSource,
  relayCliSource,
  relayHeroMobileSource,
  relayHeroSource,
  relayInstallMobileSource,
  relayInstallSource,
  relayRoutingMobileSource,
  relayRoutingSource,
} from "../relay-landing";
import { relayMascotSource } from "../relay-mascot";

beforeAll(() => defineCourierElements());

afterEach(() => {
  document.body.replaceChildren();
  document.documentElement.removeAttribute("data-courier-theme");
  document.documentElement.removeAttribute("data-courier-theme-preference");
  localStorage.clear();
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
      "courier-brand", "courier-brand-icon", "courier-button", "courier-checkbox", "courier-command-readout", "courier-icon", "courier-icon-link", "courier-locale-selector", "courier-mascot", "courier-panel", "courier-progress", "courier-route", "courier-scene", "courier-segmented-control", "courier-status", "courier-terminal", "courier-theme-selector",
    ]);
    defineCourierElements(registry);
    expect(values.size).toBe(17);
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
    brand.product = "Operations";
    document.body.append(brand);
    await brand.updateComplete;
    expect(brand.shadowRoot?.textContent).toContain("Courier");
    expect(brand.shadowRoot?.textContent).toContain("Operations");
    expect(brand.shadowRoot?.querySelector("img")?.src).toContain("relay-mark");
    expect(brand.shadowRoot?.querySelector("img")?.alt).toBe("");

    const mascot = document.createElement("courier-mascot") as CourierMascot;
    mascot.alt = "Relay";
    mascot.eager = true;
    mascot.mobileSource = relayHeroMobileSource;
    mascot.source = relayMascotSource;
    document.body.append(mascot);
    await mascot.updateComplete;
    expect(mascot.shadowRoot?.querySelector("img")?.alt).toBe("Relay");
    expect(mascot.shadowRoot?.querySelector("img")?.src).toContain("relay-mascot");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("loading")).toBe("eager");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("fetchpriority")).toBe("high");
    expect(mascot.shadowRoot?.querySelector("source")?.srcset).toContain("relay-journey-hero-mobile");
    mascot.mobileSource = "";
    mascot.eager = false;
    await mascot.updateComplete;
    expect(mascot.shadowRoot?.querySelector("source")).toBeNull();
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("loading")).toBe("lazy");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("fetchpriority")).toBe("auto");

    expect([relayAccessMobileSource, relayAccessSource, relayOperationsMobileSource, relayOperationsSource].every((source) => source.includes("relay-terminal-"))).toBe(true);
    expect([
      relayCliMobileSource, relayCliSource, relayHeroMobileSource, relayHeroSource,
      relayInstallMobileSource, relayInstallSource, relayRoutingMobileSource, relayRoutingSource,
    ].every((source) => source.includes("relay-journey-"))).toBe(true);

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
    expect((CourierSegmentedControl.styles as { cssText: string }).cssText).toContain("--courier-control-frame-size");
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

  it("tracks a fine pointer without transforming the stationary base scene", async () => {
    expect(pointerPosition({ left: 0, top: 0, width: 0, height: 0 }, 10, 10)).toEqual({ x: 50, y: 50 });
    expect(pointerPosition({ left: 10, top: 20, width: 100, height: 200 }, -20, 300)).toEqual({ x: 0, y: 100 });
    expect(sceneAmbientPosition).toEqual({ x: 72, y: 42 });
    expect(smoothPointerPosition({ x: 1, y: 1 }, { x: 1, y: 1 }, 16)).toEqual({ position: { x: 1, y: 1 }, settled: true });
    expect(smoothPointerPosition({ x: 0, y: 0 }, { x: 100, y: 100 }, -1)).toEqual({ position: { x: 0, y: 0 }, settled: false });
    expect(smoothPointerPosition({ x: 0, y: 0 }, { x: 100, y: 100 }, 100).position.x).toBeGreaterThan(50);

    const media = vi.fn((query: string) => ({ matches: query.includes("pointer: fine"), addEventListener: vi.fn(), removeEventListener: vi.fn() }));
    const callbacks: FrameRequestCallback[] = [];
    const request = vi.fn((handler: FrameRequestCallback) => { callbacks.push(handler); return request.mock.calls.length; });
    const cancel = vi.fn();
    vi.stubGlobal("matchMedia", media);
    vi.stubGlobal("requestAnimationFrame", request);
    vi.stubGlobal("cancelAnimationFrame", cancel);

    const scene = document.createElement("courier-scene") as CourierScene;
    scene.source = relayHeroSource;
    scene.mobileSource = relayHeroMobileSource;
    scene.eager = true;
    scene.getBoundingClientRect = () => ({ left: 10, top: 20, width: 100, height: 200, right: 110, bottom: 220, x: 10, y: 20, toJSON: () => ({}) });
    document.body.append(scene);
    await scene.updateComplete;
    const base = scene.shadowRoot?.querySelector(".base") as HTMLElement;
    expect(base.getAttribute("style")).toBeNull();
    expect(scene.shadowRoot?.querySelectorAll("courier-mascot")).toHaveLength(1);
    expect((CourierScene.styles as { cssText: string }).cssText).toContain("pointer-events: none");
    scene.dispatchEvent(new MouseEvent("pointermove", { clientX: 85, clientY: 70 }));
    scene.dispatchEvent(new MouseEvent("pointermove", { clientX: 90, clientY: 80 }));
    await scene.updateComplete;
    expect(scene.shadowRoot?.querySelectorAll("courier-mascot")).toHaveLength(2);
    expect(request).toHaveBeenCalledTimes(1);
    callbacks.shift()?.(16);
    expect(parseFloat(scene.style.getPropertyValue("--scene-pointer-x"))).toBeGreaterThan(72);
    expect(parseFloat(scene.style.getPropertyValue("--scene-pointer-x"))).toBeLessThan(80);
    expect(parseFloat(scene.style.getPropertyValue("--scene-pointer-y"))).toBeLessThan(42);
    expect(callbacks).toHaveLength(1);
    let timestamp = 16;
    for (let index = 0; index < 200 && callbacks.length; index += 1) {
      timestamp += 16;
      callbacks.shift()?.(timestamp);
    }
    expect(scene.style.getPropertyValue("--scene-pointer-x")).toBe("80%");
    expect(scene.style.getPropertyValue("--scene-pointer-y")).toBe("30%");
    scene.dispatchEvent(new MouseEvent("pointerleave"));
    for (let index = 0; index < 200 && callbacks.length; index += 1) {
      timestamp += 16;
      callbacks.shift()?.(timestamp);
    }
    expect(scene.style.getPropertyValue("--scene-pointer-x")).toBe("");
    expect(scene.style.getPropertyValue("--scene-pointer-y")).toBe("");
    await scene.updateComplete;
    expect(scene.shadowRoot?.querySelectorAll("courier-mascot")).toHaveLength(1);
    expect(cancel).not.toHaveBeenCalled();

    scene.dispatchEvent(new MouseEvent("pointermove", { clientX: 50, clientY: 50 }));
    const activeFrame = request.mock.results.at(-1)?.value;
    scene.remove();
    expect(cancel).toHaveBeenCalledWith(activeFrame);

    const detached = new CourierScene();
    (detached as unknown as { advance(timestamp: number): void }).advance(1);
    detached.disconnectedCallback();
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
    scene.source = relayHeroSource;
    document.body.append(scene);
    await scene.updateComplete;
    expect(scene.shadowRoot?.querySelector("courier-mascot")).toBeNull();
    expect(instances[0]?.options?.rootMargin).toBe("50% 0px");
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
  });

  it("keeps scene pointer state fixed for coarse pointers and reduced motion", async () => {
    const request = vi.fn();
    vi.stubGlobal("requestAnimationFrame", request);
    vi.stubGlobal("cancelAnimationFrame", vi.fn());
    vi.stubGlobal("matchMedia", vi.fn((query: string) => ({ matches: query.includes("prefers-reduced-motion"), addEventListener: vi.fn(), removeEventListener: vi.fn() })));
    const scene = document.createElement("courier-scene") as CourierScene;
    document.body.append(scene);
    await scene.updateComplete;
    scene.dispatchEvent(new MouseEvent("pointermove", { clientX: 1, clientY: 1 }));
    scene.dispatchEvent(new MouseEvent("pointerleave"));
    expect(request).not.toHaveBeenCalled();
  });

  it("renders terminals and immutable command readouts with reserved copy feedback", async () => {
    const terminal = document.createElement("courier-terminal") as CourierTerminal;
    terminal.heading = "Operations";
    terminal.status = "ready";
    terminal.textContent = "Transcript";
    document.body.append(terminal);
    await terminal.updateComplete;
    expect(terminal.shadowRoot?.textContent).toContain("Operations");
    expect(terminal.shadowRoot?.textContent).toContain("ready");
    terminal.status = "";
    await terminal.updateComplete;
    expect(terminal.shadowRoot?.querySelector(".status")).toBeNull();

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
  it("moves through branded segments with radio keyboard behavior", async () => {
    expect(nextSegmentIndex("ArrowLeft", 0, 3)).toBe(2);
    expect(nextSegmentIndex("ArrowUp", 1, 3)).toBe(0);
    expect(nextSegmentIndex("ArrowRight", 2, 3)).toBe(0);
    expect(nextSegmentIndex("ArrowDown", 0, 3)).toBe(1);
    expect(nextSegmentIndex("Home", 2, 3)).toBe(0);
    expect(nextSegmentIndex("End", 0, 3)).toBe(2);
    expect(nextSegmentIndex("Tab", 0, 3)).toBeUndefined();
    expect(nextSegmentIndex("ArrowRight", 0, 0)).toBeUndefined();

    const control = document.createElement("courier-segmented-control") as CourierSegmentedControl;
    control.label = "Mode";
    control.value = "one";
    control.options = [{ value: "one", label: "One" }, { value: "two", label: "Two" }];
    const changes: string[] = [];
    control.addEventListener("courier-segment-change", (event) => changes.push((event as CustomEvent<string>).detail));
    document.body.append(control);
    await control.updateComplete;
    const buttons = [...control.shadowRoot!.querySelectorAll("button")];
    expect(buttons[0].getAttribute("aria-checked")).toBe("true");
    buttons[0].click();
    expect(changes).toEqual([]);
    buttons[0].removeAttribute("data-value");
    buttons[0].click();
    expect(changes).toEqual([]);
    buttons[0].dataset.value = "one";
    buttons[0].dispatchEvent(new KeyboardEvent("keydown", { key: "Tab", bubbles: true }));
    buttons[0].dispatchEvent(new KeyboardEvent("keydown", { key: "ArrowRight", bubbles: true }));
    await control.updateComplete;
    expect(changes).toEqual(["two"]);
    expect(control.value).toBe("two");
  });

  it("persists theme changes and follows system state", async () => {
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
    const control = selector.shadowRoot?.querySelector("courier-segmented-control") as CourierSegmentedControl;
    await control.updateComplete;
    expect(control.iconOnly).toBe(true);
    expect(control.shadowRoot?.querySelector("legend")?.classList.contains("sr-only")).toBe(true);
    expect(control.shadowRoot?.querySelectorAll("button courier-icon")).toHaveLength(3);
    expect(control.shadowRoot?.querySelector("button span")).toBeNull();
    expect(control.shadowRoot?.querySelector('button[data-value="system"]')?.getAttribute("aria-label")).toBe("Системная");
    (control.shadowRoot?.querySelector('button[data-value="light"]') as HTMLButtonElement).click();
    await selector.updateComplete;
    expect(detail).toBe("light");
    expect(localStorage.getItem("courier.theme")).toBe("light");
    expect(control.label).toBe("Тема");
    selector.remove();
    expect(listeners.size).toBe(0);

    const disconnected = new CourierThemeSelector();
    disconnected.disconnectedCallback();
  });

  it("negotiates and persists locale changes", async () => {
    localStorage.setItem("courier.locale", "ru");
    const selector = document.createElement("courier-locale-selector") as CourierLocaleSelector;
    let detail = "";
    selector.addEventListener("courier-locale-change", (event) => { detail = (event as CustomEvent<string>).detail; });
    document.body.append(selector);
    await selector.updateComplete;
    expect(selector.locale).toBe("ru");
    const control = selector.shadowRoot?.querySelector("courier-segmented-control") as CourierSegmentedControl;
    await control.updateComplete;
    expect(control.iconOnly).toBe(true);
    expect(control.shadowRoot?.querySelectorAll("button courier-icon")).toHaveLength(0);
    expect([...control.shadowRoot!.querySelectorAll("button .symbol")].map((symbol) => symbol.textContent)).toEqual(["🇬🇧", "🇷🇺"]);
    expect(control.shadowRoot?.querySelector('button[data-value="en"]')?.getAttribute("aria-label")).toBe("Английский");
    control.dispatchEvent(new CustomEvent("courier-segment-change", { detail: "unsupported", bubbles: true }));
    await selector.updateComplete;
    expect(selector.locale).toBe("en");
    expect(detail).toBe("en");
    expect(localStorage.getItem("courier.locale")).toBe("en");

    const descriptor = Object.getOwnPropertyDescriptor(globalThis, "localStorage");
    Object.defineProperty(globalThis, "localStorage", { configurable: true, get: () => { throw new Error("blocked"); } });
    expect(() => control.dispatchEvent(new CustomEvent("courier-segment-change", { detail: "ru", bubbles: true }))).not.toThrow();
    if (descriptor) {
      Object.defineProperty(globalThis, "localStorage", descriptor);
    }
  });
});
