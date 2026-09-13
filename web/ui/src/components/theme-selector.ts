import { LitElement, html } from "lit";
import { translate, type Locale } from "../i18n";
import { controlStyles, fieldStyles } from "../styles";
import { browserThemeState, parseTheme, type ThemePreference, type ThemeState } from "../theme";

export class CourierThemeSelector extends LitElement {
  static properties = {
    preference: { type: String },
    locale: { type: String },
  };

  static styles = [controlStyles, fieldStyles];

  preference: ThemePreference = "system";
  locale: Locale = "en";
  private state?: ThemeState;

  connectedCallback(): void {
    super.connectedCallback();
    this.state = browserThemeState();
    this.preference = this.state.preference;
  }

  disconnectedCallback(): void {
    this.state?.destroy();
    super.disconnectedCallback();
  }

  private change(event: Event): void {
    const preference = parseTheme((event.currentTarget as HTMLSelectElement).value);
    this.preference = preference;
    this.state?.set(preference);
    this.dispatchEvent(new CustomEvent<ThemePreference>("courier-theme-change", { detail: preference, bubbles: true, composed: true }));
  }

  protected render() {
    return html`<label>${translate(this.locale, "theme.label")}
      <select .value=${this.preference} @change=${this.change}>
        <option value="system">${translate(this.locale, "theme.system")}</option>
        <option value="light">${translate(this.locale, "theme.light")}</option>
        <option value="dark">${translate(this.locale, "theme.dark")}</option>
      </select>
    </label>`;
  }
}
