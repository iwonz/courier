import { LitElement, css, html, nothing } from "lit";
import { browserLocale, browserThemeState, defineCourierElements, type Locale, type ThemeState } from "@courier/ui";
import { relayMascotSource } from "@courier/ui/relay";
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
    :host {
      display: block;
      min-height: 100vh;
      padding: 0 1rem 3rem;
      color: var(--courier-color-text);
      background-color: var(--courier-color-canvas);
      background-image: linear-gradient(var(--courier-color-grid) 1px, transparent 1px), linear-gradient(90deg, var(--courier-color-grid) 1px, transparent 1px);
      background-size: 2.5rem 2.5rem;
      font-family: var(--courier-font-sans);
    }
    main, courier-panel { min-width: 0; }
    main { width: min(66rem, 100%); margin: 0 auto; }
    header { display: flex; min-height: 5rem; align-items: center; justify-content: space-between; gap: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    nav, .controls, .actions, .row { display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; }
    .workspace { display: grid; gap: 1rem; padding-top: clamp(2rem, 6vw, 5rem); }
    .operation-head { display: flex; align-items: end; justify-content: space-between; gap: 1rem; padding-bottom: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .operation-head > div { display: grid; gap: 0.45rem; }
    .eyebrow, .label { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 750; letter-spacing: 0.08em; text-transform: uppercase; }
    h1, h2, p { margin: 0; overflow-wrap: anywhere; }
    h1 { font-family: var(--courier-font-display); font-size: clamp(2.4rem, 6vw, 4.75rem); font-weight: 830; letter-spacing: -0.06em; line-height: 0.95; }
    h2 { font-size: clamp(1.35rem, 4vw, 2rem); letter-spacing: -0.035em; }
    p { line-height: 1.6; }
    .muted { color: var(--courier-color-muted); }
    .access { display: grid; grid-template-columns: minmax(0, 1fr) minmax(11rem, 0.45fr); gap: 1rem; overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .access-copy { display: grid; align-content: center; gap: 1rem; padding: clamp(1.5rem, 5vw, 3.5rem); }
    .access-art { position: relative; min-height: 24rem; overflow: hidden; background: var(--courier-graphite-900); }
    .access-art::before { content: ""; position: absolute; inset: 0; opacity: 0.15; background-image: linear-gradient(rgb(243 244 233 / 0.2) 1px, transparent 1px), linear-gradient(90deg, rgb(243 244 233 / 0.2) 1px, transparent 1px); background-size: 2rem 2rem; }
    .access-art courier-mascot { position: absolute; right: -16%; bottom: -7%; width: 125%; }
    .error { padding: 0.85rem 1rem; border-left: 3px solid var(--courier-warning); color: var(--courier-color-text); background: color-mix(in srgb, var(--courier-warning) 12%, transparent); }
    form { display: grid; gap: 0.75rem; }
    .signin { grid-template-columns: minmax(0, 1fr) auto; }
    input { width: 100%; min-width: 0; min-height: 2.75rem; padding: 0 0.8rem; border: 1px solid var(--courier-color-border-strong); border-radius: var(--courier-radius-sm); color: var(--courier-color-text); background: var(--courier-color-surface); font: inherit; }
    input:focus-visible, button.link:focus-visible, a:focus-visible { outline: 3px solid var(--courier-beak); outline-offset: 2px; }
    .route-overview { display: grid; gap: 0.75rem; padding: 1rem; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface); }
    .delivery-panel { display: grid; gap: 1.25rem; padding: clamp(1.25rem, 4vw, 2rem); border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .delivery-title { display: flex; align-items: start; justify-content: space-between; gap: 1rem; }
    .delivery-title > div { display: grid; gap: 0.35rem; }
    .toolbar { display: flex; gap: 0.75rem; align-items: center; flex-wrap: wrap; padding: 0.9rem 0; border-top: 1px solid var(--courier-color-border); border-bottom: 1px solid var(--courier-color-border); }
    a, button.link { color: var(--courier-color-text); font-weight: 750; }
    button.link { appearance: none; padding: 0; border: 0; background: transparent; font: inherit; text-decoration: underline; cursor: pointer; }
    .upload-zone { display: grid; gap: 0.75rem; padding: clamp(1.25rem, 4vw, 2rem); border: 1px dashed var(--courier-color-border-strong); border-radius: var(--courier-radius-md); background: var(--courier-color-surface); }
    .upload-zone label { display: grid; gap: 0.5rem; font-weight: 800; }
    .upload-zone input { min-height: auto; padding: 0.75rem; }
    ul { margin: 0; padding: 0; list-style: none; border-top: 1px solid var(--courier-color-border); }
    li { display: grid; grid-template-columns: auto minmax(0, 1fr) auto auto; gap: 0.8rem; align-items: center; min-height: 3.75rem; padding: 0.75rem 0; border-bottom: 1px solid var(--courier-color-border); }
    li courier-icon { color: var(--courier-color-muted); }
    .entry-name { min-width: 0; overflow-wrap: anywhere; }
    .size { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.75rem; font-variant-numeric: tabular-nums; }
    .loading { display: grid; min-height: 14rem; place-items: center; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); color: var(--courier-color-muted); background: var(--courier-color-surface-raised); font-family: var(--courier-font-mono); }
    @media (max-width: 44rem) {
      header { align-items: flex-start; padding: 1rem 0; }
      nav { justify-content: flex-end; }
      .access { grid-template-columns: 1fr; }
      .access-art { display: none; }
      .signin { grid-template-columns: 1fr; }
      li { grid-template-columns: auto minmax(0, 1fr) auto; }
      li .size { display: none; }
      .delivery-title { display: grid; }
    }
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
      return html`<li><courier-icon name="folder"></courier-icon><button class="link entry-name" @click=${(event: Event) => this.openDirectory(event, entryPath)}>${entry.name}</button><span class="size">${entry.size} ${this.t("itemSize")}</span><a href=${downloadURL(entryPath, true)}>${this.t("downloadArchive")}</a></li>`;
    }
    return html`<li><courier-icon name="parcel"></courier-icon><span class="entry-name">${entry.name}</span><span class="size">${entry.size} ${this.t("itemSize")}</span><a href=${downloadURL(entryPath)}>${this.t("download")}</a></li>`;
  }

  render() {
    const entries = this.metadata?.entries ?? [];
    return html`
      <main>
        <header>
          <courier-brand product=${this.t("brandProduct")}></courier-brand>
          <nav><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        <div class="workspace">
          <div class="operation-head"><div><span class="eyebrow">${this.t("privateRoute")}</span><h1>${this.t("title")}</h1></div>${this.metadata ? html`<courier-status tone="signal">${this.t("ready")}</courier-status>` : nothing}</div>
          ${this.failed ? html`
            <div class="access">
              <div class="access-copy">
                <span class="eyebrow">${this.t("privateRoute")}</span>
                <h2>${this.t("accessTitle")}</h2>
                <p class="muted">${this.t("accessHelp")}</p>
                <p class="error" role="alert">${this.t("failed")}</p>
                <form class="signin" @submit=${this.signIn}><input name="password" type="password" autocomplete="current-password" placeholder=${this.t("password")}><courier-button type="submit" variant="primary">${this.t("signIn")}</courier-button></form>
                <courier-button @click=${this.refresh}>${this.t("retry")}</courier-button>
              </div>
              <div class="access-art"><courier-mascot alt="" .source=${relayMascotSource}></courier-mascot></div>
            </div>
          ` : nothing}
          ${this.metadata ? html`
            <div class="route-overview"><span class="label">${this.t("confirmed")}</span><courier-route source="sender" destination=${this.metadata.name}></courier-route></div>
            <section class="delivery-panel">
              <div class="delivery-title"><div><span class="eyebrow">${this.t("manifest")}</span><h2>${this.metadata.name}</h2></div><courier-status tone="signal">${this.t("ready")}</courier-status></div>
              ${this.metadata.type === "upload" ? html`
                <div class="upload-zone"><h2>${this.t("uploadTitle")}</h2><p class="muted">${this.t("uploadHelp")}</p><label>${this.t("upload")}<input type="file" @change=${this.sendFile}></label></div>
              ` : html`
                ${this.metadata.type === "file" ? html`<div class="toolbar"><courier-icon name="download"></courier-icon><a href=${downloadURL(this.metadata.path)}>${this.t("download")}</a></div>` : html`
                  <div class="toolbar"><a href=${downloadURL(this.metadata.path, true)}>${this.t("downloadAll")}</a>${this.metadata.path ? html`<button class="link" @click=${(event: Event) => this.openDirectory(event, parentPath(this.metadata!.path))}>${this.t("up")}</button>` : nothing}</div>
                  ${entries.length === 0 ? html`<p class="muted">${this.t("empty")}</p>` : html`<ul>${entries.map((entry) => this.entry(entry))}</ul>`}
                `}
              `}
            </section>
          ` : !this.failed ? html`<div class="loading"><courier-status>${this.t("loading")}</courier-status></div>` : nothing}
        </div>
      </main>
    `;
  }
}
