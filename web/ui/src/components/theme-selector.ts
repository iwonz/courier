import { LitElement, css, html } from "lit";
import { translate, type Locale } from "../i18n";
import { browserPreferenceController, nextPreference, preferenceThemeOrder, type BrowserPreferenceController } from "../preferences";
import type { ThemePreference } from "../theme";

const icons: Readonly<Record<ThemePreference, "system" | "sun" | "moon">> = { system: "system", light: "sun", dark: "moon" };

export class CourierThemeSelector extends LitElement {
  static properties = { preference: { type: String }, locale: { type: String } };
  static styles = css`
    :host { display: inline-flex; }
    button { appearance: none; display: inline-grid; width: var(--courier-control-frame-size, 2.35rem); height: var(--courier-control-frame-size, 2.35rem); place-items: center; padding: 0; border: 1px solid var(--courier-control-frame-border, #8f8a81); border-radius: var(--courier-control-frame-radius, 999px); color: var(--courier-control-frame-color, #57544f); background: var(--courier-control-frame-surface, transparent); cursor: pointer; transition: color var(--courier-duration, 150ms) var(--courier-ease, ease), border-color var(--courier-duration, 150ms) var(--courier-ease, ease), background var(--courier-duration, 150ms) var(--courier-ease, ease); }
    button:hover { color: var(--courier-control-frame-hover-color, #0e0f0d); border-color: var(--courier-control-frame-hover-border, #ad431d); background: var(--courier-control-frame-hover-surface, transparent); }
    button:focus-visible { outline: var(--courier-control-frame-focus-width, 3px) solid var(--courier-control-frame-focus, #d95f2b); outline-offset: 2px; }
    courier-icon { width: 1.05rem; height: 1.05rem; }
  `;
  preference: ThemePreference = "system";
  locale: Locale = "en";
  private controller?: BrowserPreferenceController;
  private readonly synchronize = (): void => { if (this.controller) this.preference = this.controller.theme; };
  connectedCallback(): void { super.connectedCallback(); this.controller = browserPreferenceController(); this.synchronize(); this.controller.addEventListener("change", this.synchronize); }
  disconnectedCallback(): void { this.controller?.removeEventListener("change", this.synchronize); super.disconnectedCallback(); }
  private cycle(): void {
    const preference = this.controller?.cycleTheme() ?? nextPreference(preferenceThemeOrder, this.preference);
    this.preference = preference;
    this.dispatchEvent(new CustomEvent<ThemePreference>("courier-theme-change", { detail: preference, bubbles: true, composed: true }));
  }
  protected render() {
    const next = nextPreference(preferenceThemeOrder, this.preference);
    const currentLabel = translate(this.locale, `theme.${this.preference}`);
    const nextLabel = translate(this.locale, `theme.${next}`);
    const label = `${translate(this.locale, "theme.label")}: ${currentLabel}. ${translate(this.locale, "preference.next")}: ${nextLabel}`;
    return html`<button type="button" aria-label=${label} title=${label} @click=${this.cycle}><courier-icon name=${icons[this.preference]}></courier-icon></button>`;
  }
}
