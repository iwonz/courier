import type { Locale } from "./i18n";

const byteUnits = ["B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"] as const;
const bitUnits = ["bit", "Kibit", "Mibit", "Gibit", "Tibit", "Pibit", "Eibit"] as const;

function number(value: number, locale: Locale, maximumFractionDigits: number): string {
  return new Intl.NumberFormat(locale, { maximumFractionDigits, minimumFractionDigits: 0 })
    .format(value)
    .replace(/[\u00a0\u202f]/g, " ");
}

function quantity(value: number, units: readonly string[], locale: Locale): string {
  if (!Number.isFinite(value) || value < 0) return "—";
  let scaled = value;
  let unit = 0;
  while (scaled >= 1024 && unit < units.length - 1) {
    scaled /= 1024;
    unit += 1;
  }
  return `${number(scaled, locale, unit === 0 ? 0 : 1)} ${units[unit]}`;
}

export function formatBytes(value: number, locale: Locale = "en"): string {
  return quantity(value, byteUnits, locale);
}

export function formatBits(value: number, locale: Locale = "en"): string {
  return quantity(value, bitUnits, locale);
}

export function formatExactBytes(value: number, locale: Locale = "en"): string {
  if (!Number.isFinite(value) || value < 0) return "—";
  return `${number(Math.round(value), locale, 0)} B`;
}

export function formatBytesDetailed(value: number, locale: Locale = "en"): string {
  const compact = formatBytes(value, locale);
  const exact = formatExactBytes(value, locale);
  return compact === exact ? exact : `${compact} (${exact})`;
}

export function formatByteRate(value: number, locale: Locale = "en"): string {
  const formatted = formatBytes(value, locale);
  return formatted === "—" ? formatted : `${formatted}/s`;
}

export function formatBitRate(value: number, locale: Locale = "en"): string {
  const formatted = formatBits(value, locale);
  return formatted === "—" ? formatted : `${formatted}/s`;
}
