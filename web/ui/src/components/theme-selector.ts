import { LitElement, html } from "lit";
import { translate, type Locale } from "../i18n";
import { browserThemeState, parseTheme, type ThemePreference, type ThemeState } from "../theme";

export class CourierThemeSelector extends LitElement {
  static properties = {
    preference: { type: String },
    locale: { type: String },
  };

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

  private change(event: CustomEvent<string>): void {
    const preference = parseTheme(event.detail);
    this.preference = preference;
    this.state?.set(preference);
    this.dispatchEvent(new CustomEvent<ThemePreference>("courier-theme-change", { detail: preference, bubbles: true, composed: true }));
  }

  protected render() {
    return html`<courier-segmented-control
      icon-only
      .label=${translate(this.locale, "theme.label")}
      .value=${this.preference}
      .options=${[
        { value: "system", label: translate(this.locale, "theme.system"), icon: "system" },
        { value: "light", label: translate(this.locale, "theme.light"), icon: "sun" },
        { value: "dark", label: translate(this.locale, "theme.dark"), icon: "moon" },
      ]}
      @courier-segment-change=${this.change}
    ></courier-segmented-control>`;
  }
}
