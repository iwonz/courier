export const themePreferences = ["system", "light", "dark"] as const;
export type ThemePreference = typeof themePreferences[number];
export type ResolvedTheme = Exclude<ThemePreference, "system">;

export interface ThemeRoot {
  dataset: DOMStringMap;
}

export interface ThemeMedia {
  readonly matches: boolean;
  addEventListener(type: "change", listener: () => void): void;
  removeEventListener(type: "change", listener: () => void): void;
}

export interface ThemeStorage {
  getItem(key: string): string | null;
  setItem(key: string, value: string): void;
}

export const themeStorageKey = "courier.theme";

export function parseTheme(value: string | null | undefined): ThemePreference {
  return themePreferences.includes(value as ThemePreference) ? value as ThemePreference : "system";
}

export function resolveTheme(preference: ThemePreference, media?: Pick<ThemeMedia, "matches">): ResolvedTheme {
  if (preference !== "system") {
    return preference;
  }
  return media?.matches ? "dark" : "light";
}

export function readTheme(storage?: ThemeStorage): ThemePreference {
  if (!storage) {
    return "system";
  }
  try {
    return parseTheme(storage.getItem(themeStorageKey));
  } catch {
    return "system";
  }
}

export function writeTheme(storage: ThemeStorage | undefined, preference: ThemePreference): void {
  if (!storage) {
    return;
  }
  try {
    storage.setItem(themeStorageKey, preference);
  } catch {
    // Persistence is optional in restricted or private browser contexts.
  }
}

export class ThemeState {
  preference: ThemePreference;
  private readonly onSystemChange = (): void => this.apply();

  constructor(
    private readonly root: ThemeRoot,
    private readonly storage?: ThemeStorage,
    private readonly media?: ThemeMedia,
    initial?: ThemePreference,
  ) {
    this.preference = initial ?? readTheme(storage);
    this.media?.addEventListener("change", this.onSystemChange);
    this.apply();
  }

  set(preference: ThemePreference): void {
    this.preference = preference;
    writeTheme(this.storage, preference);
    this.apply();
  }

  destroy(): void {
    this.media?.removeEventListener("change", this.onSystemChange);
  }

  private apply(): void {
    this.root.dataset.courierTheme = resolveTheme(this.preference, this.media);
    this.root.dataset.courierThemePreference = this.preference;
  }
}

export function browserThemeState(): ThemeState {
  let storage: Storage | undefined;
  try {
    storage = globalThis.localStorage;
  } catch {
    storage = undefined;
  }
  const media = globalThis.matchMedia?.("(prefers-color-scheme: dark)");
  return new ThemeState(document.documentElement, storage, media);
}
