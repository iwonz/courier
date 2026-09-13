import { en, type Catalog, type MessageKey } from "./en";
import { ru } from "./ru";

export { en, ru };
export type { Catalog, MessageKey };

export const catalogs = { en, ru } as const satisfies Record<string, Catalog>;
export const supportedLocales = Object.keys(catalogs) as Array<keyof typeof catalogs>;
export type Locale = keyof typeof catalogs;
export const localeStorageKey = "courier.locale";

export interface LocaleStorage {
  getItem(key: string): string | null;
  setItem(key: string, value: string): void;
}

export function parseLocale(value: string | null | undefined): Locale | undefined {
  if (!value) {
    return undefined;
  }
  const base = value.toLowerCase().split("-")[0] as Locale;
  return supportedLocales.includes(base) ? base : undefined;
}

export function resolveLocale(languages: readonly string[]): Locale {
  for (const language of languages) {
    const locale = parseLocale(language);
    if (locale) {
      return locale;
    }
  }
  return "en";
}

export function readLocale(storage: LocaleStorage | undefined, languages: readonly string[]): Locale {
  if (storage) {
    try {
      const stored = parseLocale(storage.getItem(localeStorageKey));
      if (stored) {
        return stored;
      }
    } catch {
      // Browser language negotiation remains available when storage is blocked.
    }
  }
  return resolveLocale(languages);
}

export function writeLocale(storage: LocaleStorage | undefined, locale: Locale): void {
  if (!storage) {
    return;
  }
  try {
    storage.setItem(localeStorageKey, locale);
  } catch {
    // Persistence is optional in restricted or private browser contexts.
  }
}

export function translate(locale: string | null | undefined, key: MessageKey): string {
  const selected = parseLocale(locale) ?? "en";
  return catalogs[selected][key];
}

export function browserLocale(): Locale {
  let storage: Storage | undefined;
  try {
    storage = globalThis.localStorage;
  } catch {
    storage = undefined;
  }
  const languages = globalThis.navigator?.languages ?? [];
  return readLocale(storage, languages);
}
