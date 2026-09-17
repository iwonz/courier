export const en = {
  "theme.label": "Theme",
  "theme.system": "System",
  "theme.light": "Light",
  "theme.dark": "Dark",
  "locale.label": "Language",
  "locale.en": "English",
  "locale.ru": "Russian",
  "preference.next": "Next",
  "progress.label": "Delivery progress",
  "action.cancel": "Cancel",
  "action.close": "Close",
} as const;

export type MessageKey = keyof typeof en;
export type Catalog = { readonly [Key in MessageKey]: string };
