import { parseLocale, type Locale } from "@courier/ui";

const en = {
  brandProduct: "File delivery CLI",
  title: "Move files. Keep control.",
  tagline: "A dependable CLI for local, SSH, browser, and webhook deliveries.",
  subline: "One binary plans the route, protects the destination, and reports only the bytes it can confirm.",
  install: "Install Courier",
  installShort: "Install",
  releases: "GitHub Releases",
  installIntro: "Choose a repository-owned channel. Every release resolves to the same dependency-free Courier binary.",
  primaryChannel: "Recommended · macOS and Linux",
  alternateChannels: "More ways to install",
  packages: "Linux packages",
  packagesDetail: "Deb, RPM, APK, and Arch packages for Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux, and Manjaro.",
  direct: "Direct binaries",
  directDetail: "Signed release archives and checksums for macOS, Linux, Windows, and supported BSD helpers.",
  downloads: "Open downloads",
  routeShort: "Routes",
  routeTitle: "Choose the handoff",
  routeIntro: "Point at or focus a source, then a valid destination. Every connection below comes from the shipped CLI contract.",
  sourceLabel: "Source endpoints",
  destinationLabel: "Valid destinations",
  from: "from",
  to: "to",
  routeExample: "Command shape",
  allowed: "Available options",
  cliShort: "CLI",
  cliTitle: "Commands and options",
  cliIntro: "One generated reference for the public command tree and every shipped option.",
  commandKind: "Kind",
  product: "Product",
  system: "System",
  optionDefault: "Default",
  applies: "Applies to",
  repeatable: "Repeatable",
  yes: "yes",
  no: "no",
  inspectContract: "Open the complete CLI contract",
  githubLabel: "Courier on GitHub",
  heroStatus: "Route systems ready",
  heroRouteLabel: "Typical dispatch",
  endpointLocal: "Local",
  endpointSsh: "Remote",
  endpointWeb: "Web",
  endpointWebhook: "Webhook",
  endpointHttp: "Webhook",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  brandProduct: "CLI доставки файлов",
  title: "Переносите файлы. Сохраняйте контроль.",
  tagline: "Надёжный CLI для локальных, SSH, браузерных и webhook-доставок.",
  subline: "Один бинарник планирует маршрут, защищает назначение и сообщает только тот объём, который смог подтвердить.",
  install: "Установить Courier",
  installShort: "Установка",
  releases: "Релизы GitHub",
  installIntro: "Выберите канал из репозитория. Каждый релиз устанавливает один и тот же автономный бинарник Courier.",
  primaryChannel: "Рекомендуется · macOS и Linux",
  alternateChannels: "Другие способы установки",
  packages: "Пакеты Linux",
  packagesDetail: "Пакеты Deb, RPM, APK и Arch для Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux и Manjaro.",
  direct: "Готовые бинарники",
  directDetail: "Архивы релизов и checksums для macOS, Linux, Windows и поддерживаемых BSD helper-сборок.",
  downloads: "Открыть загрузки",
  routeShort: "Маршруты",
  routeTitle: "Выберите передачу",
  routeIntro: "Наведите указатель или сфокусируйте источник, затем допустимое назначение. Все связи взяты из опубликованного контракта CLI.",
  sourceLabel: "Варианты источника",
  destinationLabel: "Допустимые назначения",
  from: "from",
  to: "to",
  routeExample: "Форма команды",
  allowed: "Доступные параметры",
  cliShort: "CLI",
  cliTitle: "Команды и параметры",
  cliIntro: "Единый сгенерированный справочник публичных команд и всех опубликованных параметров.",
  commandKind: "Тип",
  product: "Основная",
  system: "Системная",
  optionDefault: "По умолчанию",
  applies: "Применяется к",
  repeatable: "Повторяемый",
  yes: "да",
  no: "нет",
  inspectContract: "Открыть полный контракт CLI",
  githubLabel: "Courier на GitHub",
  heroStatus: "Маршруты готовы",
  heroRouteLabel: "Типовая доставка",
  endpointLocal: "Локально",
  endpointSsh: "Remote",
  endpointWeb: "Web",
  endpointWebhook: "Webhook",
  endpointHttp: "Webhook",
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
