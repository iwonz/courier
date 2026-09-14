import { afterEach, beforeAll, describe, expect, it, vi } from "vitest";
import { CourierButton } from "./button";
import { CourierBrand, CourierMascot, CourierRoute, CourierStatus } from "./brand";
import { CourierLocaleSelector } from "./locale-selector";
import { CourierPanel } from "./panel";
import { CourierProgress, progressRatio } from "./progress";
import { CourierSegmentedControl, nextSegmentIndex } from "./segmented-control";
import { CourierThemeSelector } from "./theme-selector";
import { defineCourierElements, type ElementRegistry } from "../define";
import { CourierIcon, iconNames, resolveIcon } from "../icons";
import { relayOperationsSource } from "../relay-admin";
import { relayAccessSource } from "../relay-delivery";
import { relayDispatchSource, relayInstallSource, relayRoutingSource, relayVerifySource } from "../relay-landing";
import { relayMascotSource } from "../relay-mascot";

beforeAll(() => defineCourierElements());

afterEach(() => {
  document.body.replaceChildren();
  document.documentElement.removeAttribute("data-courier-theme");
  document.documentElement.removeAttribute("data-courier-theme-preference");
  localStorage.clear();
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
      "courier-brand", "courier-button", "courier-icon", "courier-locale-selector", "courier-mascot", "courier-panel", "courier-progress", "courier-route", "courier-segmented-control", "courier-status", "courier-theme-selector",
    ]);
    defineCourierElements(registry);
    expect(values.size).toBe(11);
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
    expect(brand.shadowRoot?.querySelectorAll("svg path").length).toBeGreaterThan(3);

    const mascot = document.createElement("courier-mascot") as CourierMascot;
    mascot.alt = "Relay";
    mascot.eager = true;
    mascot.source = relayMascotSource;
    document.body.append(mascot);
    await mascot.updateComplete;
    expect(mascot.shadowRoot?.querySelector("img")?.alt).toBe("Relay");
    expect(mascot.shadowRoot?.querySelector("img")?.src).toContain("relay-mascot");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("loading")).toBe("eager");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("fetchpriority")).toBe("high");
    mascot.eager = false;
    await mascot.updateComplete;
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("loading")).toBe("lazy");
    expect(mascot.shadowRoot?.querySelector("img")?.getAttribute("fetchpriority")).toBe("auto");

    expect([
      relayAccessSource, relayDispatchSource, relayInstallSource,
      relayOperationsSource, relayRoutingSource, relayVerifySource,
    ].every((source) => source.includes("relay-") || source.includes("readme-route"))).toBe(true);

    const route = document.createElement("courier-route") as CourierRoute;
    route.source = "./data";
    route.destination = "server:/data";
    document.body.append(route);
    await route.updateComplete;
    expect(route.shadowRoot?.textContent).toContain("./data");
    expect(route.shadowRoot?.textContent).toContain("server:/data");

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
    icon.name = "language-en";
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector("text")?.textContent).toBe("EN");
    icon.name = "language-ru";
    await icon.updateComplete;
    expect(icon.shadowRoot?.querySelector("text")?.textContent).toBe("RU");
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
    expect(control.shadowRoot?.querySelectorAll("button courier-icon")).toHaveLength(2);
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
