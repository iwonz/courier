import { parseLocale, type Locale } from "@courier/ui";

const en = {
  title: "Courier delivery",
  loading: "Loading delivery…",
  retry: "Retry",
  download: "Download",
  downloadArchive: "Download archive",
  downloadAll: "Download this directory",
  upload: "Upload file",
  up: "Up",
  password: "Password",
  signIn: "Sign in",
  empty: "This directory is empty.",
  failed: "The delivery is unavailable or authorization is required.",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  title: "Доставка Courier",
  loading: "Загрузка доставки…",
  retry: "Повторить",
  download: "Скачать",
  downloadArchive: "Скачать архив",
  downloadAll: "Скачать этот каталог",
  upload: "Загрузить файл",
  up: "Наверх",
  password: "Пароль",
  signIn: "Войти",
  empty: "Каталог пуст.",
  failed: "Доставка недоступна или требуется авторизация.",
} as const satisfies Catalog;

export const dataCatalogs = { en, ru } as const;
export type DataMessage = keyof typeof en;

export function dataLocale(languages: readonly string[] = globalThis.navigator.languages): Locale {
  for (const language of languages) {
    const locale = parseLocale(language);
    if (locale) {
      return locale;
    }
  }
  return "en";
}

export function dataText(locale: Locale, message: DataMessage): string {
  return dataCatalogs[locale][message];
}
