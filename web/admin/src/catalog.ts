import { parseLocale, type Locale } from "@courier/ui";

const en = {
  title: "Courier administration",
  refresh: "Refresh",
  retry: "Retry",
  loading: "Loading servers…",
  empty: "No Courier data servers found.",
  failed: "Administration data is temporarily unavailable.",
  conflict: "The policy changed elsewhere. Refresh before editing again.",
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
  save: "Save policy",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  title: "Управление Courier",
  refresh: "Обновить",
  retry: "Повторить",
  loading: "Загрузка серверов…",
  empty: "Серверы данных Courier не найдены.",
  failed: "Данные управления временно недоступны.",
  conflict: "Политика была изменена. Обновите данные перед повторным редактированием.",
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
  save: "Сохранить политику",
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
