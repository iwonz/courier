import { LitElement, css, html, nothing } from "lit";
import { browserLocale, browserThemeState, defineCourierElements, type Locale, type ThemeState } from "@courier/ui";
import { loadServers, savePolicy, stopTarget, subscribeSnapshots, type Delivery, type Policy, type Server, type Snapshot } from "./api";
import { adminText } from "./catalog";

defineCourierElements();

export class CourierAdminApp extends LitElement {
  static properties = {
    locale: { state: true },
    snapshot: { state: true },
    failed: { state: true },
    conflict: { state: true },
  };

  static styles = css`
    :host { box-sizing: border-box; display: block; min-height: 100vh; padding: clamp(1rem, 4vw, 3rem); background: var(--courier-color-canvas); color: var(--courier-color-text); font-family: var(--courier-font-sans); }
    main, courier-panel { min-width: 0; }
    main { width: min(70rem, 100%); margin: 0 auto; display: grid; gap: 1rem; }
    header, nav, .row, .actions { display: flex; gap: .75rem; align-items: center; justify-content: space-between; flex-wrap: wrap; }
    section { display: grid; gap: .75rem; }
    article { padding: 1rem; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-medium); display: grid; gap: .75rem; }
    dl { display: grid; grid-template-columns: max-content 1fr; gap: .35rem .75rem; margin: 0; }
    dt { color: var(--courier-color-muted); }
    dd { margin: 0; overflow-wrap: anywhere; }
    form { display: grid; grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr)); gap: .75rem; align-items: end; }
    label { display: grid; gap: .25rem; }
    input, select { box-sizing: border-box; min-width: 0; max-width: 100%; min-height: 2.75rem; padding: 0 .5rem; }
    .status { text-transform: uppercase; letter-spacing: .08em; font-size: .75rem; }
    .status.live { color: var(--courier-color-accent); }
    h1, h2, p, strong { overflow-wrap: anywhere; }
    @media (max-width: 38rem) { dl { grid-template-columns: 1fr; } dt { margin-top: .35rem; } }
  `;

  private locale: Locale = browserLocale();
  private snapshot?: Snapshot;
  private failed = false;
  private conflict = false;
  private theme?: ThemeState;
  private unsubscribe?: () => void;

  connectedCallback(): void {
    super.connectedCallback();
    this.theme = browserThemeState();
    this.unsubscribe = subscribeSnapshots((snapshot) => {
      this.snapshot = snapshot;
      this.failed = false;
    });
    void this.refresh();
  }

  disconnectedCallback(): void {
    this.unsubscribe?.();
    this.theme?.destroy();
    super.disconnectedCallback();
  }

  async refresh(): Promise<void> {
    this.failed = false;
    this.conflict = false;
    try {
      this.snapshot = await loadServers();
    } catch {
      this.failed = true;
    }
  }

  setLocale(event: CustomEvent<Locale>): void {
    this.locale = event.detail;
  }

  async stop(kind: "servers" | "deliveries", id: string): Promise<void> {
    try {
      await stopTarget(kind, id);
      await this.refresh();
    } catch {
      this.failed = true;
    }
  }

  async save(event: SubmitEvent, delivery: Delivery): Promise<void> {
    event.preventDefault();
    const data = new FormData(event.currentTarget as HTMLFormElement);
    const policy: Policy = {
      ...delivery.policy,
      version: delivery.policy.version + 1,
      auth: String(data.get("auth")) as Policy["auth"],
      authAttempts: Number(data.get("attempts")),
      authFailAction: String(data.get("failAction")) as Policy["authFailAction"],
      noUi: data.get("noUi") === "on",
    };
    this.failed = false;
    this.conflict = false;
    try {
      await savePolicy(delivery, policy);
      await this.refresh();
    } catch (error) {
      this.conflict = error instanceof Error && error.name === "ConflictError";
      this.failed = !this.conflict;
    }
  }

  private t(message: Parameters<typeof adminText>[1]): string {
    return adminText(this.locale, message);
  }

  private delivery(item: Delivery) {
    return html`
      <article>
        <div class="row"><strong>${item.id}</strong><courier-button @click=${() => this.stop("deliveries", item.id)}>${this.t("stopDelivery")}</courier-button></div>
        <dl>
          <dt>${this.t("source")}</dt><dd>${item.source || "<unavailable>"}</dd>
          <dt>${this.t("destination")}</dt><dd>${item.destination || "<unavailable>"}</dd>
          <dt>${this.t("transferred")}</dt><dd>${item.counters.confirmed}</dd>
        </dl>
        <form @submit=${(event: SubmitEvent) => this.save(event, item)}>
          <label>${this.t("authentication")}<select name="auth"><option selected=${item.policy.auth === "none"}>none</option><option selected=${item.policy.auth === "basic"}>basic</option><option selected=${item.policy.auth === "password"}>password</option></select></label>
          <label>${this.t("attempts")}<input name="attempts" type="number" min="1" .value=${String(item.policy.authAttempts)}></label>
          <label>${this.t("failAction")}<select name="failAction"><option selected=${item.policy.authFailAction === "ban"}>ban</option><option selected=${item.policy.authFailAction === "stop"}>stop</option></select></label>
          <label><input name="noUi" type="checkbox" ?checked=${item.policy.noUi}> ${this.t("noUi")}</label>
          <courier-button type="submit">${this.t("save")}</courier-button>
        </form>
      </article>
    `;
  }

  private server(item: Server) {
    return html`
      <courier-panel>
        <section>
          <div class="row"><h2>${item.id}</h2><span class="status ${item.status}">${this.t(item.status)}</span></div>
          <div class="actions"><span>${item.bind}</span><courier-button @click=${() => this.stop("servers", item.id)}>${this.t("stopServer")}</courier-button></div>
          ${item.deliveries.map((delivery) => this.delivery(delivery))}
        </section>
      </courier-panel>
    `;
  }

  render() {
    const servers = this.snapshot?.servers ?? [];
    return html`
      <main>
        <header>
          <h1>${this.t("title")}</h1>
          <nav><courier-button @click=${this.refresh}>${this.t("refresh")}</courier-button><courier-theme-selector></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        ${this.failed ? html`<courier-panel><p role="alert">${this.t("failed")}</p><courier-button @click=${this.refresh}>${this.t("retry")}</courier-button></courier-panel>` : nothing}
        ${this.conflict ? html`<courier-panel><p role="alert">${this.t("conflict")}</p><courier-button @click=${this.refresh}>${this.t("refresh")}</courier-button></courier-panel>` : nothing}
        ${this.snapshot ? servers.length === 0 ? html`<courier-panel>${this.t("empty")}</courier-panel>` : servers.map((server) => this.server(server)) : !this.failed ? html`<courier-panel>${this.t("loading")}</courier-panel>` : nothing}
      </main>
    `;
  }
}

if (!customElements.get("courier-admin-app")) {
  customElements.define("courier-admin-app", CourierAdminApp);
}
