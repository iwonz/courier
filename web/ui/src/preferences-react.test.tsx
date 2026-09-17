import * as React from "react";
import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import {
  BrowserPreferenceController,
  browserPreferenceController,
  nextPreference,
  preferenceLocaleOrder,
  preferenceThemeOrder,
  resetBrowserPreferenceController,
} from "./preferences";
import { LocaleSelector, PreferenceProvider, ThemeSelector, usePreferences } from "./preferences-react";
import { type ThemeMedia } from "./theme";

function makeController(initialTheme: "system" | "light" | "dark" = "system"): { controller: BrowserPreferenceController; listeners: Set<() => void>; storage: Map<string, string> } {
  const listeners = new Set<() => void>();
  const storage = new Map<string, string>();
  const media: ThemeMedia = {
    matches: false,
    addEventListener: (_type, listener) => listeners.add(listener),
    removeEventListener: (_type, listener) => listeners.delete(listener),
  };
  const persistence = { getItem: (key: string) => storage.get(key) ?? null, setItem: (key: string, value: string) => { storage.set(key, value); } };
  return { controller: new BrowserPreferenceController({ dataset: {} as DOMStringMap }, persistence, media, ["en-US"], initialTheme), listeners, storage };
}

afterEach(() => {
  cleanup();
  resetBrowserPreferenceController();
  localStorage.clear();
  document.documentElement.removeAttribute("data-courier-theme");
  document.documentElement.removeAttribute("data-courier-theme-preference");
  vi.restoreAllMocks();
});

describe("browser preference controller", () => {
  it("cycles generic, theme, and locale orders and persists changes", () => {
    expect(nextPreference(["one", "two"] as const, "one")).toBe("two");
    expect(nextPreference(["one", "two"] as const, "two")).toBe("one");
    const { controller, listeners, storage } = makeController();
    const changed = vi.fn();
    controller.addEventListener("change", changed);
    expect(controller.theme).toBe("system");
    expect(controller.locale).toBe("en");
    expect(controller.cycleTheme()).toBe("light");
    expect(controller.setTheme("dark")).toBe("dark");
    expect(controller.cycleTheme()).toBe("system");
    expect(controller.cycleLocale()).toBe("ru");
    expect(controller.setLocale("en")).toBe("en");
    expect(storage.get("courier.theme")).toBe("system");
    expect(storage.get("courier.locale")).toBe("en");
    expect(changed).toHaveBeenCalledTimes(5);
    expect(listeners.size).toBe(1);
    controller.destroy();
    expect(listeners.size).toBe(0);
  });

  it("creates and reuses a singleton with browser fallbacks", () => {
    const first = browserPreferenceController();
    expect(browserPreferenceController()).toBe(first);
    resetBrowserPreferenceController();
    const descriptor = Object.getOwnPropertyDescriptor(globalThis, "localStorage");
    const navigatorDescriptor = Object.getOwnPropertyDescriptor(globalThis, "navigator");
    Object.defineProperty(globalThis, "localStorage", { configurable: true, get: () => { throw new Error("blocked"); } });
    Object.defineProperty(globalThis, "navigator", { configurable: true, value: undefined });
    const blocked = browserPreferenceController();
    expect(blocked.theme).toBe("system");
    resetBrowserPreferenceController();
    if (descriptor) Object.defineProperty(globalThis, "localStorage", descriptor);
    if (navigatorDescriptor) Object.defineProperty(globalThis, "navigator", navigatorDescriptor);
  });
});

describe("React preference controls", () => {
  it("requires a provider", () => {
    function Consumer(): React.JSX.Element { usePreferences(); return <span />; }
    expect(() => render(<Consumer />)).toThrow("PreferenceProvider is required");
  });

  it("renders one button per preference and updates every theme and locale", () => {
    const { controller } = makeController();
    render(<PreferenceProvider controller={controller}><ThemeSelector /><LocaleSelector /></PreferenceProvider>);
    const buttons = screen.getAllByRole("button");
    expect(buttons).toHaveLength(2);
    expect(buttons[0]!.getAttribute("aria-label")).toContain("System");
    fireEvent.click(buttons[0]!);
    expect(buttons[0]!.getAttribute("aria-label")).toContain("Light");
    fireEvent.click(buttons[0]!);
    expect(buttons[0]!.getAttribute("aria-label")).toContain("Dark");
    fireEvent.click(buttons[0]!);
    expect(buttons[0]!.getAttribute("aria-label")).toContain("System");
    fireEvent.click(buttons[1]!);
    expect(buttons[1]!.getAttribute("aria-label")).toContain("Русский");
    expect(preferenceThemeOrder).toEqual(["system", "light", "dark"]);
    expect(preferenceLocaleOrder).toEqual(["en", "ru"]);
  });

  it("uses the browser controller when none is injected", () => {
    render(<PreferenceProvider><ThemeSelector /></PreferenceProvider>);
    expect(screen.getByRole("button").getAttribute("aria-label")).toContain("Theme");
  });
});
