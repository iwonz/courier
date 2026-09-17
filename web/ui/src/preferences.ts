import { readLocale, writeLocale, type Locale, type LocaleStorage } from "./i18n";
import { ThemeState, type ThemeMedia, type ThemePreference, type ThemeRoot, type ThemeStorage } from "./theme";

export const preferenceThemeOrder = ["system", "light", "dark"] as const;
export const preferenceLocaleOrder = ["en", "ru"] as const;

export function nextPreference<Value extends string>(order: readonly Value[], value: Value): Value {
  const index = order.indexOf(value);
  return order[(index + 1) % order.length]!;
}

export class BrowserPreferenceController extends EventTarget {
  private readonly themeState: ThemeState;
  locale: Locale;

  constructor(
    root: ThemeRoot,
    private readonly storage: (ThemeStorage & LocaleStorage) | undefined,
    media: ThemeMedia | undefined,
    languages: readonly string[],
    initialTheme?: ThemePreference,
  ) {
    super();
    this.themeState = new ThemeState(root, storage, media, initialTheme);
    this.locale = readLocale(storage, languages);
  }

  get theme(): ThemePreference {
    return this.themeState.preference;
  }

  setTheme(preference: ThemePreference): ThemePreference {
    this.themeState.set(preference);
    this.dispatchEvent(new Event("change"));
    return preference;
  }

  cycleTheme(): ThemePreference {
    return this.setTheme(nextPreference(preferenceThemeOrder, this.theme));
  }

  setLocale(locale: Locale): Locale {
    this.locale = locale;
    writeLocale(this.storage, locale);
    this.dispatchEvent(new Event("change"));
    return locale;
  }

  cycleLocale(): Locale {
    return this.setLocale(nextPreference(preferenceLocaleOrder, this.locale));
  }

  destroy(): void {
    this.themeState.destroy();
  }
}

let browserController: BrowserPreferenceController | undefined;

export function browserPreferenceController(): BrowserPreferenceController {
  if (browserController) return browserController;
  let storage: Storage | undefined;
  try {
    storage = globalThis.localStorage;
  } catch {
    storage = undefined;
  }
  const media = globalThis.matchMedia?.("(prefers-color-scheme: dark)");
  const languages = globalThis.navigator?.languages ?? [];
  browserController = new BrowserPreferenceController(document.documentElement, storage, media, languages);
  return browserController;
}

export function resetBrowserPreferenceController(): void {
  browserController?.destroy();
  browserController = undefined;
}
