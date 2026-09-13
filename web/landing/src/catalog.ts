import { parseLocale, type Locale } from "@courier/ui";

const en = {
  title: "Courier CLI",
  tagline: "Move files safely across local, SSH, browser, and webhook boundaries.",
  subline: "One dependency-free binary. Explicit routes. Transactional destinations.",
  install: "Install",
  commands: "Commands",
  routes: "Routes",
  options: "Options",
  examples: "Examples",
  downloads: "Downloads",
  documentation: "Documentation",
  source: "Source on GitHub",
  releases: "GitHub Releases",
  installIntro: "Choose the repository-owned channel that fits your environment.",
  packages: "Linux packages",
  packagesDetail: "Deb, RPM, APK, and Arch packages cover Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux, and Manjaro.",
  direct: "Direct binaries",
  directDetail: "Signed release metadata and SHA-256 checksums are published for macOS, Linux, Windows, and BSD helpers.",
  commandKind: "Kind",
  commandSyntax: "Command",
  product: "Product",
  system: "System",
  route: "Route",
  from: "Source",
  to: "Destination",
  allowed: "Allowed options",
  option: "Option",
  default: "Default",
  applies: "Applies to",
  repeatable: "Repeatable",
  yes: "yes",
  no: "no",
  version: "Contract",
  target: "Target release",
  security: "Security model",
  cliReference: "CLI reference",
  installGuide: "Installation guide",
  license: "MIT License",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  title: "Courier CLI",
  tagline: "Безопасный перенос файлов между локальными, SSH, браузерными и webhook-средами.",
  subline: "Один автономный бинарник. Явные маршруты. Транзакционная запись.",
  install: "Установка",
  commands: "Команды",
  routes: "Маршруты",
  options: "Параметры",
  examples: "Примеры",
  downloads: "Загрузки",
  documentation: "Документация",
  source: "Исходный код на GitHub",
  releases: "Релизы GitHub",
  installIntro: "Выберите канал из репозитория Courier для своей среды.",
  packages: "Пакеты Linux",
  packagesDetail: "Пакеты Deb, RPM, APK и Arch поддерживают Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux и Manjaro.",
  direct: "Готовые бинарники",
  directDetail: "Метаданные релиза и SHA-256 публикуются для macOS, Linux, Windows и вспомогательных BSD-сборок.",
  commandKind: "Тип",
  commandSyntax: "Команда",
  product: "Основная",
  system: "Системная",
  route: "Маршрут",
  from: "Источник",
  to: "Назначение",
  allowed: "Допустимые параметры",
  option: "Параметр",
  default: "По умолчанию",
  applies: "Применяется к",
  repeatable: "Повторяемый",
  yes: "да",
  no: "нет",
  version: "Контракт",
  target: "Целевой релиз",
  security: "Модель безопасности",
  cliReference: "Справочник CLI",
  installGuide: "Руководство по установке",
  license: "Лицензия MIT",
} as const satisfies Catalog;

export const landingCatalogs = { en, ru } as const;
export type LandingMessage = keyof typeof en;

export function landingLocale(languages: readonly string[] = globalThis.navigator.languages): Locale {
  for (const language of languages) {
    const locale = parseLocale(language);
    if (locale) return locale;
  }
  return "en";
}

export function landingText(locale: Locale, message: LandingMessage): string {
  return landingCatalogs[locale][message];
}
