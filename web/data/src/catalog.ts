import { parseLocale, type Locale } from "@courier/ui";

const en = {
  brandProduct: "Delivery terminal",
  title: "Courier delivery",
  privateRoute: "Private route",
  loading: "Preparing the delivery route…",
  retry: "Retry connection",
  download: "Download file",
  downloadArchive: "Download as archive",
  downloadAll: "Download directory",
  upload: "Choose a file",
  uploadTitle: "Dispatch a file",
  uploadHelp: "Select one file. Courier checks policy and reserves the final path before committing it.",
  up: "Parent directory",
  password: "Delivery password",
  signIn: "Verify access",
  accessTitle: "Identity check",
  accessHelp: "This route is protected. Verify access to reveal its delivery metadata.",
  empty: "No entries are available at this path.",
  failed: "The route is unavailable or authorization is required. No delivery metadata was revealed.",
  ready: "Route ready",
  manifest: "Delivery manifest",
  confirmed: "Verified handoff",
  entryTypeFile: "File",
  entryTypeDirectory: "Directory",
  itemSize: "Bytes",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  brandProduct: "Терминал доставки",
  title: "Доставка Courier",
  privateRoute: "Приватный маршрут",
  loading: "Подготовка маршрута доставки…",
  retry: "Повторить подключение",
  download: "Скачать файл",
  downloadArchive: "Скачать архивом",
  downloadAll: "Скачать каталог",
  upload: "Выбрать файл",
  uploadTitle: "Отправить файл",
  uploadHelp: "Выберите один файл. Courier проверит правила и зарезервирует конечный путь до фиксации.",
  up: "Родительский каталог",
  password: "Пароль доставки",
  signIn: "Подтвердить доступ",
  accessTitle: "Проверка доступа",
  accessHelp: "Маршрут защищён. Подтвердите доступ, чтобы увидеть данные доставки.",
  empty: "По этому пути нет доступных объектов.",
  failed: "Маршрут недоступен или требуется авторизация. Данные доставки не были раскрыты.",
  ready: "Маршрут готов",
  manifest: "Манифест доставки",
  confirmed: "Подтверждённая передача",
  entryTypeFile: "Файл",
  entryTypeDirectory: "Каталог",
  itemSize: "Байт",
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
