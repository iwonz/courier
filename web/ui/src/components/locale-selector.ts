import { LitElement, css, html } from "lit";
import { translate, type Locale } from "../i18n";
import { browserPreferenceController, nextPreference, preferenceLocaleOrder, type BrowserPreferenceController } from "../preferences";

const symbols: Readonly<Record<Locale, string>> = { en: "🇬🇧", ru: "🇷🇺" };

export class CourierLocaleSelector extends LitElement {
  static properties = { locale: { type: String } };
  static styles = css`
    :host { display: inline-flex; }
    button { appearance: none; display: inline-grid; width: var(--courier-control-frame-size, 2.35rem); height: var(--courier-control-frame-size, 2.35rem); place-items: center; padding: 0; border: 1px solid var(--courier-control-frame-border, #8f8a81); border-radius: var(--courier-control-frame-radius, 999px); color: var(--courier-control-frame-color, #57544f); background: var(--courier-control-frame-surface, transparent); cursor: pointer; transition: border-color var(--courier-duration, 150ms) var(--courier-ease, ease), background var(--courier-duration, 150ms) var(--courier-ease, ease); }
    button:hover { border-color: var(--courier-control-frame-hover-border, #ad431d); background: var(--courier-control-frame-hover-surface, transparent); }
    button:focus-visible { outline: var(--courier-control-frame-focus-width, 3px) solid var(--courier-control-frame-focus, #d95f2b); outline-offset: 2px; }
    span { font: 1.02rem/1 system-ui, sans-serif; }
  `;
  locale: Locale = "en";
  private controller?: BrowserPreferenceController;
  private readonly synchronize = (): void => { if (this.controller) this.locale = this.controller.locale; };
  connectedCallback(): void { super.connectedCallback(); this.controller = browserPreferenceController(); this.synchronize(); this.controller.addEventListener("change", this.synchronize); }
  disconnectedCallback(): void { this.controller?.removeEventListener("change", this.synchronize); super.disconnectedCallback(); }
  private cycle(): void {
    const locale = this.controller?.cycleLocale() ?? nextPreference(preferenceLocaleOrder, this.locale);
    this.locale = locale;
    this.dispatchEvent(new CustomEvent<Locale>("courier-locale-change", { detail: locale, bubbles: true, composed: true }));
  }
  protected render() {
    const next = nextPreference(preferenceLocaleOrder, this.locale);
    const label = `${translate(this.locale, "locale.label")}: ${translate(this.locale, `locale.${this.locale}`)}. ${translate(this.locale, "preference.next")}: ${translate(this.locale, `locale.${next}`)}`;
    return html`<button type="button" aria-label=${label} title=${label} @click=${this.cycle}><span aria-hidden="true">${symbols[this.locale]}</span></button>`;
  }
}
