import { afterEach, describe, expect, it, vi } from "vitest";
import {
  ThemeState,
  parseTheme,
  readTheme,
  resolveTheme,
  themeStorageKey,
  writeTheme,
  type ThemeMedia,
  type ThemeStorage,
} from "./theme";

describe("theme state", () => {
  afterEach(() => {
    vi.restoreAllMocks();
    document.documentElement.removeAttribute("data-courier-theme");
    document.documentElement.removeAttribute("data-courier-theme-preference");
    localStorage.clear();
  });

  it("parses and resolves every preference", () => {
    expect(parseTheme("system")).toBe("system");
    expect(parseTheme("light")).toBe("light");
    expect(parseTheme("dark")).toBe("dark");
    expect(parseTheme("unknown")).toBe("system");
    expect(parseTheme(undefined)).toBe("system");
    expect(resolveTheme("light", { matches: true })).toBe("light");
    expect(resolveTheme("dark", { matches: false })).toBe("dark");
    expect(resolveTheme("system", { matches: true })).toBe("dark");
    expect(resolveTheme("system")).toBe("light");
  });

  it("reads and writes optional or restricted storage", () => {
    expect(readTheme()).toBe("system");
    const values = new Map<string, string>();
    const storage: ThemeStorage = {
      getItem: (key) => values.get(key) ?? null,
      setItem: (key, value) => { values.set(key, value); },
    };
    writeTheme(storage, "dark");
    expect(values.get(themeStorageKey)).toBe("dark");
    expect(readTheme(storage)).toBe("dark");
    writeTheme(undefined, "light");
    const restricted: ThemeStorage = {
      getItem: () => { throw new Error("blocked"); },
      setItem: () => { throw new Error("blocked"); },
    };
    expect(readTheme(restricted)).toBe("system");
    expect(() => writeTheme(restricted, "light")).not.toThrow();
  });

  it("reacts to system changes and removes its listener", () => {
    const listeners = new Set<() => void>();
    const media: ThemeMedia = {
      matches: false,
      addEventListener: (_type, listener) => { listeners.add(listener); },
      removeEventListener: (_type, listener) => { listeners.delete(listener); },
    };
    const root = { dataset: {} as DOMStringMap };
    const storage = { getItem: () => null, setItem: vi.fn() };
    const state = new ThemeState(root, storage, media);
    expect(root.dataset.courierTheme).toBe("light");
    state.set("dark");
    expect(root.dataset.courierTheme).toBe("dark");
    expect(storage.setItem).toHaveBeenCalledWith(themeStorageKey, "dark");
    state.set("system");
    Object.defineProperty(media, "matches", { value: true });
    listeners.forEach((listener) => listener());
    expect(root.dataset.courierTheme).toBe("dark");
    state.destroy();
    expect(listeners.size).toBe(0);

    const explicit = new ThemeState(root, undefined, undefined, "light");
    expect(root.dataset.courierThemePreference).toBe("light");
    explicit.destroy();
  });

});
