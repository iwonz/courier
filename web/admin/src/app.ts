import { LitElement, css, html, nothing } from "lit";
import { browserLocale, browserThemeState, defineCourierElements, formControlStyles, type Locale, type ThemeState } from "@courier/ui";
import { relayOperationsSource } from "@courier/ui/relay-admin";
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

  static styles = [formControlStyles, css`
    :host {
      display: block;
      min-height: 100vh;
      padding: 0 1rem 4rem;
      color: var(--courier-color-text);
      background-color: var(--courier-color-canvas);
      background-image: linear-gradient(var(--courier-color-grid) 1px, transparent 1px), linear-gradient(90deg, var(--courier-color-grid) 1px, transparent 1px);
      background-size: 2.5rem 2.5rem;
      font-family: var(--courier-font-sans);
    }
    main, courier-panel, article { min-width: 0; }
    main { width: min(78rem, 100%); margin: 0 auto; }
    header { display: flex; min-height: 5rem; align-items: center; justify-content: space-between; gap: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    nav, .row, .actions { display: flex; gap: 0.75rem; align-items: center; justify-content: space-between; flex-wrap: wrap; }
    .workspace { display: grid; gap: 1.25rem; padding-top: clamp(2rem, 6vw, 5rem); }
    .page-head { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 2rem; align-items: end; padding-bottom: 1.25rem; border-bottom: 1px solid var(--courier-color-border); }
    .page-head > div { display: grid; gap: 0.65rem; }
    .eyebrow, .label { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 750; letter-spacing: 0.08em; text-transform: uppercase; }
    h1, h2, h3, p, strong { margin: 0; overflow-wrap: anywhere; }
    h1 { font-family: var(--courier-font-display); font-size: clamp(2.5rem, 7vw, 5rem); font-weight: 830; letter-spacing: -0.065em; line-height: 0.92; }
    h2 { font-size: 1.2rem; letter-spacing: -0.03em; }
    h3 { font-size: 1rem; }
    p { line-height: 1.6; }
    .intro { max-width: 43rem; color: var(--courier-color-muted); }
    .metrics { display: grid; grid-template-columns: repeat(3, 1fr); border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .metric { display: grid; gap: 0.35rem; padding: 1.35rem; }
    .metric + .metric { border-left: 1px solid var(--courier-color-border); }
    .metric strong { font-family: var(--courier-font-display); font-size: 2rem; font-variant-numeric: tabular-nums; letter-spacing: -0.05em; }
    .metric span { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; letter-spacing: 0.075em; text-transform: uppercase; }
    .notice { display: grid; gap: 0.75rem; padding: 1rem; border: 1px solid var(--courier-color-border); border-left: 3px solid var(--courier-warning); border-radius: var(--courier-radius-sm); background: var(--courier-color-surface-raised); }
    .notice.error { border-left-color: var(--courier-danger); }
    .server-list { display: grid; gap: 1rem; }
    .server { display: grid; gap: 1rem; }
    .server-head { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 1rem; padding-bottom: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .server-title { display: grid; gap: 0.4rem; }
    .server-id, .delivery-id, dd { font-family: var(--courier-font-mono); font-size: 0.8rem; font-variant-numeric: tabular-nums; }
    .bind { display: flex; align-items: center; gap: 0.5rem; color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.8rem; }
    .delivery-stack { display: grid; gap: 0.75rem; }
    article { display: grid; gap: 1rem; padding: clamp(1rem, 3vw, 1.5rem); border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface); }
    .delivery-head { align-items: start; }
    .delivery-head > div { display: grid; gap: 0.35rem; }
    .route { display: grid; gap: 0.75rem; }
    dl { display: grid; grid-template-columns: max-content minmax(0, 1fr); gap: 0.45rem 1rem; margin: 0; padding: 0.9rem 0; border-top: 1px solid var(--courier-color-border); border-bottom: 1px solid var(--courier-color-border); }
    dt { color: var(--courier-color-muted); font-size: 0.8rem; }
    dd { margin: 0; overflow-wrap: anywhere; }
    form { display: grid; grid-template-columns: repeat(4, minmax(9rem, 1fr)) auto; gap: 0.75rem; align-items: end; }
    label { display: grid; gap: 0.35rem; color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 700; letter-spacing: 0.04em; }
    label.checkbox { grid-template-columns: auto 1fr; align-items: center; align-content: center; min-height: 2.75rem; }
    .state-brief { display: grid; grid-template-columns: minmax(0, 1fr) minmax(18rem, 0.7fr); min-height: 20rem; overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .state-copy { display: grid; align-content: center; justify-items: start; gap: 1rem; padding: clamp(1.5rem, 5vw, 3rem); }
    .state-art { position: relative; min-height: 20rem; overflow: hidden; background: var(--courier-graphite-900); }
    .state-art courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; }
    .state-art courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; }
    .empty, .loading { display: grid; min-height: 13rem; place-items: center; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); color: var(--courier-color-muted); background: var(--courier-color-surface-raised); font-family: var(--courier-font-mono); }
    @media (max-width: 64rem) { form { grid-template-columns: repeat(2, minmax(10rem, 1fr)); } }
    @media (max-width: 44rem) {
      header { align-items: flex-start; padding: 1rem 0; }
      nav { justify-content: flex-end; }
      .page-head, .server-head { grid-template-columns: 1fr; }
      .metrics { grid-template-columns: 1fr; }
      .metric + .metric { border-top: 1px solid var(--courier-color-border); border-left: 0; }
      form { grid-template-columns: 1fr; }
      dl { grid-template-columns: 1fr; }
      dt { margin-top: 0.3rem; }
      .state-brief { grid-template-columns: 1fr; }
      .state-art { min-height: 15rem; }
    }
  `];

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
    const source = item.source || this.t("unavailable");
    const destination = item.destination || this.t("unavailable");
    return html`
      <article>
        <div class="row delivery-head"><div><span class="label">${this.t("deliveries")} · ${item.route}</span><strong class="delivery-id">${item.id}</strong></div><courier-button @click=${() => this.stop("deliveries", item.id)}>${this.t("stopDelivery")}</courier-button></div>
        <div class="route"><courier-route source=${source} destination=${destination}></courier-route></div>
        <dl>
          <dt>${this.t("source")}</dt><dd>${source}</dd>
          <dt>${this.t("destination")}</dt><dd>${destination}</dd>
          <dt>${this.t("transferred")}</dt><dd>${item.counters.confirmed}</dd>
        </dl>
        <span class="label">${this.t("policy")}</span>
        <form @submit=${(event: SubmitEvent) => this.save(event, item)}>
          <label>${this.t("authentication")}<select name="auth"><option ?selected=${item.policy.auth === "none"}>none</option><option ?selected=${item.policy.auth === "basic"}>basic</option><option ?selected=${item.policy.auth === "password"}>password</option></select></label>
          <label>${this.t("attempts")}<input name="attempts" type="number" min="1" .value=${String(item.policy.authAttempts)}></label>
          <label>${this.t("failAction")}<select name="failAction"><option ?selected=${item.policy.authFailAction === "ban"}>ban</option><option ?selected=${item.policy.authFailAction === "stop"}>stop</option></select></label>
          <label class="checkbox"><input name="noUi" type="checkbox" ?checked=${item.policy.noUi}><span>${this.t("noUi")}</span></label>
          <courier-button type="submit" variant="primary">${this.t("save")}</courier-button>
        </form>
      </article>
    `;
  }

  private server(item: Server) {
    return html`
      <courier-panel>
        <section class="server">
          <div class="server-head">
            <div class="server-title"><span class="label">${this.t("server")}</span><h2 class="server-id">${item.id}</h2><span class="bind"><courier-icon name="server"></courier-icon>${item.bind}</span></div>
            <div class="actions"><courier-status tone=${item.status === "live" ? "signal" : "danger"}>${this.t(item.status)}</courier-status><courier-button @click=${() => this.stop("servers", item.id)}>${this.t("stopServer")}</courier-button></div>
          </div>
          <span class="label">${this.t("deliveries")}</span>
          <div class="delivery-stack">${item.deliveries.map((delivery) => this.delivery(delivery))}</div>
        </section>
      </courier-panel>
    `;
  }

  render() {
    const servers = this.snapshot?.servers ?? [];
    const deliveries = servers.reduce((count, server) => count + server.deliveries.length, 0);
    const confirmed = servers.reduce((serverTotal, server) => serverTotal + server.deliveries.reduce((deliveryTotal, delivery) => deliveryTotal + delivery.counters.confirmed, 0), 0);
    return html`
      <main>
        <header>
          <courier-brand product=${this.t("brandProduct")}></courier-brand>
          <nav><courier-button @click=${this.refresh}>${this.t("refresh")}</courier-button><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        <div class="workspace">
          <div class="page-head"><div><span class="eyebrow">${this.t("eyebrow")}</span><h1>${this.t("title")}</h1><p class="intro">${this.t("intro")}</p></div><courier-status tone=${this.failed ? "danger" : "signal"}>${this.t(this.failed ? "unreachable" : "live")}</courier-status></div>
          <div class="metrics">
            <div class="metric"><strong>${servers.length}</strong><span>${this.t("serversMetric")}</span></div>
            <div class="metric"><strong>${deliveries}</strong><span>${this.t("deliveriesMetric")}</span></div>
            <div class="metric"><strong>${confirmed}</strong><span>${this.t("confirmedMetric")}</span></div>
          </div>
          ${this.failed ? html`<div class="state-brief"><div class="state-copy"><span class="eyebrow">${this.t("unreachable")}</span><p role="alert">${this.t("failed")}</p><courier-button @click=${this.refresh}>${this.t("retry")}</courier-button></div><div class="state-art"><courier-mascot alt="" .source=${relayOperationsSource}></courier-mascot></div></div>` : nothing}
          ${this.conflict ? html`<div class="notice"><p role="alert">${this.t("conflict")}</p><courier-button @click=${this.refresh}>${this.t("refresh")}</courier-button></div>` : nothing}
          ${this.snapshot ? servers.length === 0 ? html`<div class="state-brief"><div class="state-copy"><span class="eyebrow">${this.t("live")}</span><p>${this.t("empty")}</p></div><div class="state-art"><courier-mascot alt="" .source=${relayOperationsSource}></courier-mascot></div></div>` : html`<div class="server-list">${servers.map((server) => this.server(server))}</div>` : !this.failed ? html`<div class="loading"><courier-status>${this.t("loading")}</courier-status></div>` : nothing}
        </div>
      </main>
    `;
  }
}

if (!customElements.get("courier-admin-app")) {
  customElements.define("courier-admin-app", CourierAdminApp);
}
