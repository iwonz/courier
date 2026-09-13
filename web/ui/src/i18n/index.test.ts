import { afterEach, describe, expect, it, vi } from "vitest";
import {
  browserLocale,
  catalogs,
  localeStorageKey,
  parseLocale,
  readLocale,
  resolveLocale,
  supportedLocales,
  translate,
  writeLocale,
  type LocaleStorage,
} from "./index";

describe("localization", () => {
  afterEach(() => {
    vi.restoreAllMocks();
    localStorage.clear();
  });

  it("keeps catalogs complete and resolves regional browser languages", () => {
    expect(supportedLocales).toEqual(["en", "ru"]);
    expect(Object.keys(catalogs.en)).toEqual(Object.keys(catalogs.ru));
    expect(parseLocale(null)).toBeUndefined();
    expect(parseLocale("")).toBeUndefined();
    expect(parseLocale("RU-ru")).toBe("ru");
    expect(parseLocale("en-US")).toBe("en");
    expect(parseLocale("de")).toBeUndefined();
    expect(resolveLocale(["de-DE", "ru-RU"])).toBe("ru");
    expect(resolveLocale(["de-DE"])).toBe("en");
    expect(translate("ru", "action.close")).toBe("Закрыть");
    expect(translate("unsupported", "action.close")).toBe("Close");
  });

  it("prefers storage and falls back when storage is absent or blocked", () => {
    const values = new Map<string, string>();
    const storage: LocaleStorage = {
      getItem: (key) => values.get(key) ?? null,
      setItem: (key, value) => { values.set(key, value); },
    };
    expect(readLocale(undefined, ["ru"])).toBe("ru");
    expect(readLocale(storage, ["ru"])).toBe("ru");
    values.set(localeStorageKey, "en");
    expect(readLocale(storage, ["ru"])).toBe("en");
    values.set(localeStorageKey, "unsupported");
    expect(readLocale(storage, ["ru"])).toBe("ru");
    writeLocale(storage, "ru");
    expect(values.get(localeStorageKey)).toBe("ru");
    writeLocale(undefined, "en");
    const restricted: LocaleStorage = {
      getItem: () => { throw new Error("blocked"); },
      setItem: () => { throw new Error("blocked"); },
    };
    expect(readLocale(restricted, ["ru"])).toBe("ru");
    expect(() => writeLocale(restricted, "en")).not.toThrow();
  });

  it("negotiates browser locale when local storage is available or blocked", () => {
    localStorage.setItem(localeStorageKey, "ru");
    expect(browserLocale()).toBe("ru");
    localStorage.clear();
    const navigatorDescriptor = Object.getOwnPropertyDescriptor(globalThis, "navigator");
    Object.defineProperty(globalThis, "navigator", { configurable: true, value: { languages: ["ru-RU"] } });
    expect(browserLocale()).toBe("ru");

    const storageDescriptor = Object.getOwnPropertyDescriptor(globalThis, "localStorage");
    Object.defineProperty(globalThis, "localStorage", { configurable: true, get: () => { throw new Error("blocked"); } });
    expect(browserLocale()).toBe("ru");
    if (storageDescriptor) {
      Object.defineProperty(globalThis, "localStorage", storageDescriptor);
    }
    Object.defineProperty(globalThis, "navigator", { configurable: true, value: undefined });
    expect(browserLocale()).toBe("en");
    if (navigatorDescriptor) {
      Object.defineProperty(globalThis, "navigator", navigatorDescriptor);
    }
  });
});
