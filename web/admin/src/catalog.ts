import { parseLocale, type Locale } from "@courier/ui";

const en = {
  eyebrow: "Local control plane",
  title: "Delivery control",
  intro: "Inspect live routes, confirmed volume, and delivery policy from one private local console.",
  refresh: "Refresh registry",
  retry: "Retry connection",
  loading: "Checking the live registry…",
  empty: "No live Courier servers are registered.",
  failed: "Administration data is temporarily unavailable. Retry the local connection.",
  conflict: "This delivery policy changed elsewhere. Refresh before editing again.",
  live: "Live",
  unreachable: "Unreachable",
  stopServer: "Stop server",
  stopDelivery: "Stop delivery",
  source: "Source",
  destination: "Destination",
  transferred: "Confirmed bytes",
  authentication: "Authentication",
  attempts: "Authentication attempts",
  failAction: "Failure action",
  noUi: "Hide delivery UI",
  save: "Apply policy",
  serversMetric: "Live registry servers",
  deliveriesMetric: "Active deliveries",
  confirmedMetric: "Confirmed bytes",
  uploadsMetric: "Recorded uploads",
  server: "Server",
  bind: "Bound address",
  deliveries: "Delivery routes",
  policy: "Delivery policy",
  unavailable: "Unavailable",
  activity: "Recent uploads",
  activityHelp: "Immutable receipts recorded for browser and webhook uploads. Newest first.",
  noActivity: "No uploads have been recorded yet.",
  originalName: "Original name",
  storedName: "Stored as",
  extractedArchive: "Extracted archive",
  actor: "Actor",
  peer: "Peer",
  received: "Received",
  size: "Size",
  entries: "Entries",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  eyebrow: "Локальный контур управления",
  title: "Управление доставками",
  intro: "Проверяйте активные маршруты, подтверждённый объём и правила доставки в одной приватной локальной консоли.",
  refresh: "Обновить реестр",
  retry: "Повторить подключение",
  loading: "Проверка активного реестра…",
  empty: "Активные серверы Courier не зарегистрированы.",
  failed: "Данные управления временно недоступны. Повторите локальное подключение.",
  conflict: "Правила этой доставки были изменены. Обновите данные перед повторным редактированием.",
  live: "Работает",
  unreachable: "Недоступен",
  stopServer: "Остановить сервер",
  stopDelivery: "Остановить доставку",
  source: "Источник",
  destination: "Назначение",
  transferred: "Подтверждено байт",
  authentication: "Аутентификация",
  attempts: "Попытки аутентификации",
  failAction: "Действие при ошибке",
  noUi: "Скрыть интерфейс доставки",
  save: "Применить правила",
  serversMetric: "Серверы активного реестра",
  deliveriesMetric: "Активные доставки",
  confirmedMetric: "Подтверждено байт",
  uploadsMetric: "Учтённые загрузки",
  server: "Сервер",
  bind: "Адрес привязки",
  deliveries: "Маршруты доставки",
  policy: "Правила доставки",
  unavailable: "Недоступно",
  activity: "Последние загрузки",
  activityHelp: "Неизменяемые записи о загрузках из браузера и через вебхук. Сначала новые.",
  noActivity: "Загрузок пока не зарегистрировано.",
  originalName: "Исходное имя",
  storedName: "Сохранено как",
  extractedArchive: "Архив распакован",
  actor: "Отправитель",
  peer: "Адрес",
  received: "Получено",
  size: "Размер",
  entries: "Объектов",
} as const satisfies Catalog;

export const adminCatalogs = { en, ru } as const;
export type AdminMessage = keyof typeof en;

export function adminLocale(languages: readonly string[] = globalThis.navigator.languages): Locale {
  for (const language of languages) {
    const locale = parseLocale(language);
    if (locale) return locale;
  }
  return "en";
}

export function adminText(locale: Locale, message: AdminMessage): string {
  return adminCatalogs[locale][message];
}
