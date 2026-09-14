import { LitElement, html } from "lit";
import { browserLocale, localeStorageKey, parseLocale, translate, writeLocale, type Locale } from "../i18n";

export class CourierLocaleSelector extends LitElement {
  static properties = { locale: { type: String } };

  locale: Locale = "en";

  connectedCallback(): void {
    super.connectedCallback();
    this.locale = browserLocale();
  }

  private change(event: CustomEvent<string>): void {
    const locale = parseLocale(event.detail) ?? "en";
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
    return html`<courier-segmented-control
      .label=${translate(this.locale, "locale.label")}
      .value=${this.locale}
      .options=${[
        { value: "en", label: translate(this.locale, "locale.en") },
        { value: "ru", label: translate(this.locale, "locale.ru") },
      ]}
      data-storage-key=${localeStorageKey}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`;
  }
}
