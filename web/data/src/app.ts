import { LitElement, css, html, nothing } from "lit";
import { browserLocale, browserThemeState, defineCourierElements, type Locale, type ThemeState } from "@courier/ui";
import { childPath, downloadURL, loadMetadata, login, parentPath, upload, type Entry, type Metadata } from "./api";
import { dataText } from "./catalog";

defineCourierElements();

export class CourierDataApp extends LitElement {
  static properties = {
    locale: { state: true },
    metadata: { state: true },
    failed: { state: true },
    csrf: { state: true },
  };

  static styles = css`
    :host { box-sizing: border-box; display: block; min-height: 100vh; padding: clamp(1rem, 4vw, 3rem); background: var(--courier-color-canvas); color: var(--courier-color-text); font-family: var(--courier-font-sans); }
    main, courier-panel { min-width: 0; }
    main { width: min(52rem, 100%); margin: 0 auto; display: grid; gap: 1rem; }
    header, nav, form, li { display: flex; gap: .75rem; align-items: center; justify-content: space-between; flex-wrap: wrap; }
    ul { list-style: none; padding: 0; display: grid; gap: .5rem; }
    li { padding: .75rem; border-bottom: 1px solid var(--courier-color-border); }
    input { box-sizing: border-box; min-width: 0; max-width: 100%; min-height: 2.75rem; padding: 0 .75rem; }
    h1, h2, p { overflow-wrap: anywhere; }
    a, button.link { color: var(--courier-color-accent); }
    button.link { appearance: none; border: 0; background: transparent; padding: 0; font: inherit; cursor: pointer; }
  `;

  private locale: Locale = browserLocale();
  private metadata?: Metadata;
  private failed = false;
  private csrf = "";
  private theme?: ThemeState;

  connectedCallback(): void {
    super.connectedCallback();
    this.theme = browserThemeState();
    void this.refresh();
  }

  disconnectedCallback(): void {
    this.theme?.destroy();
    super.disconnectedCallback();
  }

  async refresh(path = this.metadata?.path ?? ""): Promise<void> {
    this.failed = false;
    try {
      this.metadata = await loadMetadata(path);
    } catch {
      this.failed = true;
      this.metadata = undefined;
    }
  }

  async openDirectory(event: Event, path: string): Promise<void> {
    event.preventDefault();
    await this.refresh(path);
  }

  async signIn(event: SubmitEvent): Promise<void> {
    event.preventDefault();
    const form = event.currentTarget as HTMLFormElement;
    const password = new FormData(form).get("password")?.toString() ?? "";
    try {
      this.csrf = (await login(password)).csrf;
      form.reset();
      await this.refresh();
    } catch {
      this.failed = true;
    }
  }

  async sendFile(event: Event): Promise<void> {
    const input = event.currentTarget as HTMLInputElement;
    const file = input.files?.item(0);
    if (!file) {
      return;
    }
    try {
      await upload(file, this.csrf);
      input.value = "";
      await this.refresh();
    } catch {
      this.failed = true;
    }
  }

  setLocale(event: CustomEvent<Locale>): void {
    this.locale = event.detail;
  }

  private t(message: Parameters<typeof dataText>[1]): string {
    return dataText(this.locale, message);
  }

  private entry(entry: Entry) {
    const entryPath = childPath(this.metadata!.path, entry.name);
    if (entry.type === "directory") {
      return html`<li><button class="link" @click=${(event: Event) => this.openDirectory(event, entryPath)}>${entry.name}</button><a href=${downloadURL(entryPath, true)}>${this.t("downloadArchive")}</a></li>`;
    }
    return html`<li><span>${entry.name}</span><a href=${downloadURL(entryPath)}>${this.t("download")}</a></li>`;
  }

  render() {
    const entries = this.metadata?.entries ?? [];
    return html`
      <main>
        <header>
          <h1>${this.t("title")}</h1>
          <nav><courier-theme-selector></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        ${this.failed ? html`
          <courier-panel><p>${this.t("failed")}</p><courier-button @click=${this.refresh}>${this.t("retry")}</courier-button></courier-panel>
          <courier-panel><form @submit=${this.signIn}><input name="password" type="password" autocomplete="current-password" placeholder=${this.t("password")}><courier-button type="submit">${this.t("signIn")}</courier-button></form></courier-panel>
        ` : nothing}
        ${this.metadata ? html`
          <courier-panel>
            <h2>${this.metadata.name}</h2>
            ${this.metadata.type === "upload" ? html`<label>${this.t("upload")} <input type="file" @change=${this.sendFile}></label>` : html`
              ${this.metadata.type === "file" ? html`<a href=${downloadURL(this.metadata.path)}>${this.t("download")}</a>` : html`
                <a href=${downloadURL(this.metadata.path, true)}>${this.t("downloadAll")}</a>
                ${this.metadata.path ? html`<button class="link" @click=${(event: Event) => this.openDirectory(event, parentPath(this.metadata!.path))}>${this.t("up")}</button>` : nothing}
                ${entries.length === 0 ? html`<p>${this.t("empty")}</p>` : html`<ul>${entries.map((entry) => this.entry(entry))}</ul>`}
              `}
            `}
          </courier-panel>
        ` : !this.failed ? html`<courier-panel>${this.t("loading")}</courier-panel>` : nothing}
      </main>
    `;
  }
}
