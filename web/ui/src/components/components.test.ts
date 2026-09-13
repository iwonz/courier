import { afterEach, beforeAll, describe, expect, it, vi } from "vitest";
import { CourierButton } from "./button";
import { CourierLocaleSelector } from "./locale-selector";
import { CourierPanel } from "./panel";
import { CourierProgress, progressRatio } from "./progress";
import { CourierThemeSelector } from "./theme-selector";
import { defineCourierElements, type ElementRegistry } from "../define";
import { CourierIcon, iconNames, resolveIcon } from "../icons";

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
      "courier-button", "courier-icon", "courier-locale-selector", "courier-panel", "courier-progress", "courier-theme-selector",
    ]);
    defineCourierElements(registry);
    expect(values.size).toBe(6);
  });
});

describe("shared components", () => {
  it("renders buttons and panels", async () => {
    const button = document.createElement("courier-button") as CourierButton;
    button.variant = "primary";
    button.disabled = true;
    button.textContent = "Send";
    document.body.append(button);
    await button.updateComplete;
    const nativeButton = button.shadowRoot?.querySelector("button");
    expect(nativeButton?.disabled).toBe(true);
    expect(button.getAttribute("variant")).toBe("primary");

    const panel = document.createElement("courier-panel") as CourierPanel;
    panel.heading = "Delivery";
    panel.textContent = "Cargo";
    document.body.append(panel);
    await panel.updateComplete;
    expect(panel.shadowRoot?.querySelector("h2")?.textContent).toBe("Delivery");
    expect(panel.shadowRoot?.querySelector("section")?.getAttribute("aria-labelledby")).toBe("courier-panel-heading");
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
});

describe("preference selectors", () => {
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
    const select = selector.shadowRoot?.querySelector("select") as HTMLSelectElement;
    select.value = "light";
    select.dispatchEvent(new Event("change"));
    await selector.updateComplete;
    expect(detail).toBe("light");
    expect(localStorage.getItem("courier.theme")).toBe("light");
    expect(selector.shadowRoot?.querySelector("label")?.textContent).toContain("Тема");
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
    const select = selector.shadowRoot?.querySelector("select") as HTMLSelectElement;
    select.value = "unsupported";
    select.dispatchEvent(new Event("change"));
    await selector.updateComplete;
    expect(selector.locale).toBe("en");
    expect(detail).toBe("en");
    expect(localStorage.getItem("courier.locale")).toBe("en");

    const descriptor = Object.getOwnPropertyDescriptor(globalThis, "localStorage");
    Object.defineProperty(globalThis, "localStorage", { configurable: true, get: () => { throw new Error("blocked"); } });
    select.value = "ru";
    expect(() => select.dispatchEvent(new Event("change"))).not.toThrow();
    if (descriptor) {
      Object.defineProperty(globalThis, "localStorage", descriptor);
    }
  });
});
