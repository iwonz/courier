import { parseLocale, type Locale } from "@courier/ui";

const en = {
  title: "Courier delivery",
  privateRoute: "Private route",
  loading: "Preparing the delivery route…",
  retry: "Retry connection",
  download: "Download file",
  downloadArchive: "Download as archive",
  downloadAll: "Download directory",
  upload: "Choose files",
  uploadingAction: "Uploading batch…",
  uploadTitle: "Dispatch files",
  uploadHelp: "Select one or more files. Courier sends them in order and commits each file independently.",
  uploadBatchTitle: "Upload queue",
  uploadBatchProgress: "Completed",
  uploadQueuedTitle: "Queued",
  uploadingTitle: "Upload in progress",
  uploadingHelp: "Keep this page open while Courier sends the selected file.",
  uploadSuccessTitle: "File delivered",
  uploadSuccessHelp: "Courier committed this file to the destination.",
  uploadFailureTitle: "Upload failed",
  uploadFailureHelp: "Courier did not commit this file. You can select it again to retry.",
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
  title: "Доставка Courier",
  privateRoute: "Приватный маршрут",
  loading: "Подготовка маршрута доставки…",
  retry: "Повторить подключение",
  download: "Скачать файл",
  downloadArchive: "Скачать архивом",
  downloadAll: "Скачать каталог",
  upload: "Выбрать файлы",
  uploadingAction: "Загрузка файлов…",
  uploadTitle: "Отправить файлы",
  uploadHelp: "Выберите один или несколько файлов. Courier отправит их по порядку и сохранит каждый независимо.",
  uploadBatchTitle: "Очередь загрузки",
  uploadBatchProgress: "Завершено",
  uploadQueuedTitle: "В очереди",
  uploadingTitle: "Файл загружается",
  uploadingHelp: "Не закрывайте страницу, пока Courier отправляет выбранный файл.",
  uploadSuccessTitle: "Файл доставлен",
  uploadSuccessHelp: "Courier сохранил файл в каталоге назначения.",
  uploadFailureTitle: "Не удалось загрузить файл",
  uploadFailureHelp: "Courier не сохранил этот файл. Его можно выбрать повторно.",
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
