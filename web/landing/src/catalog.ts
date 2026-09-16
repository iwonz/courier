import { parseLocale, type Locale } from "@courier/ui";

const en = {
  brandProduct: "File delivery CLI",
  title: "Move files. Keep control.",
  tagline: "A dependable CLI for local, SSH, browser, and webhook deliveries.",
  subline: "One binary plans the route, protects the destination, and reports only the bytes it can confirm.",
  install: "Install Courier",
  installShort: "Install",
  installIntro: "Choose a repository-owned channel. Every release resolves to the same dependency-free Courier binary.",
  chooseChannel: "Choose an installation channel",
  packages: "Linux packages",
  packagesDetail: "Deb, RPM, APK, and Arch packages for Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux, and Manjaro.",
  direct: "Direct binaries",
  directDetail: "Signed release archives and checksums for macOS, Linux, Windows, and supported BSD helpers.",
  routeShort: "Routes",
  routeTitle: "Choose the handoff",
  routeIntro: "Choose a Source, then activate a valid Destination. Every connection below comes from the shipped CLI contract.",
  sourceLabel: "Source",
  destinationLabel: "Destination",
  from: "from",
  to: "to",
  routeExample: "Command shape",
  allowed: "Available options",
  cliShort: "CLI",
  cliTitle: "Commands and options",
  cliIntro: "One generated reference for the public command tree and every shipped option.",
  commandLabel: "Command",
  optionsLabel: "Options",
  compatibleOnly: "Compatible with selected command",
  noCompatibleOptions: "This command has no options.",
  product: "Product",
  system: "System",
  optionDefault: "Default",
  applies: "Applies to",
  repeatable: "Repeatable",
  yes: "yes",
  no: "no",
  githubLabel: "Courier on GitHub",
  endpointLocal: "Local",
  endpointSsh: "SSH",
  endpointWeb: "Web",
  endpointWebhook: "Webhook",
  endpointHttp: "HTTP(S)",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  brandProduct: "CLI доставки файлов",
  title: "Переносите файлы. Сохраняйте контроль.",
  tagline: "Надёжный CLI для локальных, SSH, браузерных и webhook-доставок.",
  subline: "Один бинарник планирует маршрут, защищает назначение и сообщает только тот объём, который смог подтвердить.",
  install: "Установить Courier",
  installShort: "Установка",
  installIntro: "Выберите канал из репозитория. Каждый релиз устанавливает один и тот же автономный бинарник Courier.",
  chooseChannel: "Выберите канал установки",
  packages: "Пакеты Linux",
  packagesDetail: "Пакеты Deb, RPM, APK и Arch для Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux и Manjaro.",
  direct: "Готовые бинарники",
  directDetail: "Архивы релизов и checksums для macOS, Linux, Windows и поддерживаемых BSD helper-сборок.",
  routeShort: "Маршруты",
  routeTitle: "Выберите передачу",
  routeIntro: "Выберите Source, затем активируйте допустимый Destination. Все связи взяты из опубликованного контракта CLI.",
  sourceLabel: "Source",
  destinationLabel: "Destination",
  from: "from",
  to: "to",
  routeExample: "Форма команды",
  allowed: "Доступные параметры",
  cliShort: "CLI",
  cliTitle: "Команды и параметры",
  cliIntro: "Единый сгенерированный справочник публичных команд и всех опубликованных параметров.",
  commandLabel: "Команда",
  optionsLabel: "Параметры",
  compatibleOnly: "Совместимые с выбранной командой",
  noCompatibleOptions: "У этой команды нет параметров.",
  product: "Основная",
  system: "Системная",
  optionDefault: "По умолчанию",
  applies: "Применяется к",
  repeatable: "Повторяемый",
  yes: "да",
  no: "нет",
  githubLabel: "Courier на GitHub",
  endpointLocal: "Local",
  endpointSsh: "SSH",
  endpointWeb: "Web",
  endpointWebhook: "Webhook",
  endpointHttp: "HTTP(S)",
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
