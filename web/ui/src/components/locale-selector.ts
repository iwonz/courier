import { LitElement, html } from "lit";
import { browserLocale, localeStorageKey, parseLocale, translate, writeLocale, type Locale } from "../i18n";
import { controlStyles, fieldStyles } from "../styles";

export class CourierLocaleSelector extends LitElement {
  static properties = { locale: { type: String } };
  static styles = [controlStyles, fieldStyles];

  locale: Locale = "en";

  connectedCallback(): void {
    super.connectedCallback();
    this.locale = browserLocale();
  }

  private change(event: Event): void {
    const locale = parseLocale((event.currentTarget as HTMLSelectElement).value) ?? "en";
    this.locale = locale;
    let storage: Storage | undefined;
    try {
      storage = globalThis.localStorage;
    } catch {
      storage = undefined;
    }
    writeLocale(storage, locale);
    this.dispatchEvent(new CustomEvent<Locale>("courier-locale-change", { detail: locale, bubbles: true, composed: true }));
  }

  protected render() {
    return html`<label>${translate(this.locale, "locale.label")}
      <select .value=${this.locale} @change=${this.change} data-storage-key=${localeStorageKey}>
        <option value="en">${translate(this.locale, "locale.en")}</option>
        <option value="ru">${translate(this.locale, "locale.ru")}</option>
      </select>
    </label>`;
  }
}
