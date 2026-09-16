import { parseLocale, type Locale } from "@courier/ui";

const en = {
  brandProduct: "File delivery CLI",
  title: "Move files. Keep control.",
  install: "Install Courier",
  installShort: "Install",
  chooseChannel: "Choose an installation channel",
  copyCommand: "Copy command",
  copiedCommand: "Copied",
  copyFailed: "Copy failed",
  packages: "Linux packages",
  packagesDetail: "Deb, RPM, APK, and Arch packages for Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux, and Manjaro.",
  direct: "Direct binaries",
  directDetail: "Signed release archives and checksums for macOS, Linux, Windows, and supported BSD helpers.",
  routeShort: "Routes",
  routeTitle: "Choose the handoff",
  routeIntro: "Choose a Source, then activate a valid Destination.",
  sourceLabel: "Source",
  destinationLabel: "Destination",
  from: "from",
  to: "to",
  allowed: "Available options",
  cliShort: "CLI",
  cliTitle: "Commands and options",
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
  installReadoutDescription: "Copy the exact command for this installation channel.",
  cliReadoutDescription: "Copy the generated usage for the selected command.",
  selectCommand: "Select a command to inspect and copy its generated usage.",
  endpointLocal: "Local",
  endpointRemote: "Remote",
  endpointWeb: "Web",
  endpointWebHook: "Web Hook",
  sourceLocalDescription: "Read from this device",
  destinationLocalDescription: "Write to this device",
  sourceRemoteDescription: "Read through a secure SSH connection",
  destinationRemoteDescription: "Write through a secure SSH connection",
  sourceWebDescription: "Open an upload page and send chosen files onward",
  destinationWebDescription: "Serve Source files through a download page",
  sourceWebHookDescription: "Receive files at a unique URL and deliver them onward",
  destinationWebHookDescription: "Send Source files to the destination URL",
} as const;

type Catalog = { readonly [Key in keyof typeof en]: string };

const ru = {
  brandProduct: "CLI доставки файлов",
  title: "Переносите файлы. Сохраняйте контроль.",
  install: "Установить Courier",
  installShort: "Установка",
  chooseChannel: "Выберите канал установки",
  copyCommand: "Скопировать команду",
  copiedCommand: "Скопировано",
  copyFailed: "Не удалось скопировать",
  packages: "Пакеты Linux",
  packagesDetail: "Пакеты Deb, RPM, APK и Arch для Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux и Manjaro.",
  direct: "Готовые бинарники",
  directDetail: "Архивы релизов и checksums для macOS, Linux, Windows и поддерживаемых BSD helper-сборок.",
  routeShort: "Маршруты",
  routeTitle: "Выберите передачу",
  routeIntro: "Выберите Source, затем активируйте допустимый Destination.",
  sourceLabel: "Source",
  destinationLabel: "Destination",
  from: "from",
  to: "to",
  allowed: "Доступные параметры",
  cliShort: "CLI",
  cliTitle: "Команды и параметры",
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
  installReadoutDescription: "Скопируйте точную команду для выбранного канала установки.",
  cliReadoutDescription: "Скопируйте сгенерированный usage выбранной команды.",
  selectCommand: "Выберите команду, чтобы посмотреть и скопировать её usage.",
  endpointLocal: "Local",
  endpointRemote: "Remote",
  endpointWeb: "Web",
  endpointWebHook: "Web Hook",
  sourceLocalDescription: "Читать с этого устройства",
  destinationLocalDescription: "Записать на это устройство",
  sourceRemoteDescription: "Читать через защищённое SSH-подключение",
  destinationRemoteDescription: "Записать через защищённое SSH-подключение",
  sourceWebDescription: "Открыть страницу загрузки и передать выбранные файлы дальше",
  destinationWebDescription: "Раздавать файлы Source через страницу скачивания",
  sourceWebHookDescription: "Принимать файлы по уникальной ссылке и доставлять дальше",
  destinationWebHookDescription: "Отправить файлы Source на destination URL",
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
