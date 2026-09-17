import { LitElement, css, html, nothing } from "lit";
import { browserLocale, defineCourierElements, formControlStyles, type Locale } from "@courier/ui";
import { adminOperationsMobileSource, adminOperationsSource } from "@courier/ui/admin-scenes";
import { loadServers, savePolicy, stopTarget, subscribeSnapshots, type Delivery, type Policy, type Server, type Snapshot } from "./api";
import { adminText } from "./catalog";

defineCourierElements();

export function reconcileDeliverySelection(snapshot: Snapshot, selected: string | undefined): string | undefined {
  const deliveries = snapshot.servers.flatMap((server) => server.deliveries);
  return deliveries.some((delivery) => delivery.id === selected) ? selected : deliveries[0]?.id;
}

export class CourierAdminApp extends LitElement {
  static properties = {
    locale: { state: true },
    snapshot: { state: true },
    failed: { state: true },
    conflict: { state: true },
    selectedDeliveryId: { state: true },
  };

  static styles = [formControlStyles, css`
    :host { display: block; min-height: 100vh; padding: 0 1rem 4rem; color: var(--courier-color-text); background: var(--courier-color-canvas); font-family: var(--courier-font-sans); }
    main, article { min-width: 0; } main { position: relative; z-index: 2; width: min(78rem, 100%); margin: 0 auto; }
    .page-scene { position: fixed; z-index: 0; inset: 0; opacity: 0.32; }
    :host::after { content: ""; position: fixed; z-index: 1; inset: 0; background: linear-gradient(90deg, var(--courier-color-canvas) 0 18%, color-mix(in srgb, var(--courier-color-canvas) 72%, transparent) 58%, color-mix(in srgb, var(--courier-color-canvas) 90%, transparent)); pointer-events: none; }
    header, nav, .row, .actions { display: flex; align-items: center; gap: 0.7rem; flex-wrap: wrap; }
    header { min-height: 5rem; justify-content: space-between; border-bottom: 1px solid color-mix(in srgb, var(--courier-color-border) 65%, transparent); }
    nav { justify-content: flex-end; }
    .workspace { display: grid; gap: 1rem; padding-top: clamp(2rem, 6vw, 4.5rem); }
    .page-head { display: grid; grid-template-columns: minmax(0, 1fr) auto; align-items: end; gap: 2rem; }
    .page-head > div { display: grid; gap: 0.6rem; }
    .eyebrow, .label { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 740; letter-spacing: 0.05em; }
    h1, h2, p, strong { margin: 0; overflow-wrap: anywhere; }
    h1 { font-family: var(--courier-font-display); font-size: clamp(2.5rem, 7vw, 5rem); font-weight: 830; letter-spacing: -0.065em; line-height: 0.92; }
    h2 { font-size: 1rem; letter-spacing: -0.02em; }
    p { line-height: 1.55; } .intro { max-width: 43rem; color: var(--courier-color-muted); }
    .metrics { display: grid; grid-template-columns: repeat(3, 1fr); }
    .metric { display: grid; gap: 0.25rem; padding: 1rem; }
    .metric + .metric { border-left: 1px solid var(--courier-workbench-line); }
    .metric strong { color: var(--courier-color-accent); font: 760 1.55rem/1 var(--courier-font-mono); font-variant-numeric: tabular-nums; }
    .metric span, .bind, .delivery-id, dd { font-family: var(--courier-font-mono); font-size: 0.75rem; }
    .metric span, .bind, dt { color: var(--courier-color-muted); }
    .operations { display: grid; grid-template-columns: minmax(15rem, 0.34fr) minmax(0, 1fr); min-height: 31rem; }
    .navigator { min-width: 0; padding: 0.75rem; border-right: 1px solid var(--courier-workbench-line); overflow: auto; }
    .server-group { display: grid; gap: 0.45rem; padding: 0.75rem 0; border-bottom: 1px solid var(--courier-workbench-line); }
    .server-head { display: grid; grid-template-columns: minmax(0, 1fr) auto; align-items: start; gap: 0.5rem; }
    .server-title { display: grid; min-width: 0; gap: 0.25rem; }
    .server-id { font-family: var(--courier-font-mono); font-size: 0.74rem; overflow-wrap: anywhere; }
    .delivery-nav { appearance: none; display: grid; width: 100%; gap: 0.2rem; padding: 0.65rem; border: 0; border-left: 2px solid transparent; color: var(--courier-color-text); background: transparent; text-align: left; cursor: pointer; }
    .delivery-nav:hover { background: color-mix(in srgb, var(--courier-color-accent) 8%, transparent); }
    .delivery-nav[aria-pressed="true"] { border-left-color: var(--courier-color-accent); background: color-mix(in srgb, var(--courier-color-accent) 12%, transparent); }
    .delivery-nav:focus-visible { outline: 3px solid var(--courier-color-accent); outline-offset: -3px; }
    .inspector { min-width: 0; overflow: auto; }
    article { display: grid; gap: 1rem; padding: clamp(1rem, 3vw, 1.5rem); }
    .delivery-head { justify-content: space-between; align-items: start; }
    .delivery-head > div { display: grid; gap: 0.35rem; }
    .route { min-width: 0; }
    dl { display: grid; grid-template-columns: max-content minmax(0, 1fr); gap: 0.45rem 1rem; margin: 0; padding: 0.9rem 0; border-block: 1px solid var(--courier-workbench-line); }
    dt { font-size: 0.8rem; } dd { margin: 0; overflow-wrap: anywhere; }
    form { display: grid; grid-template-columns: repeat(4, minmax(8rem, 1fr)) auto; gap: 0.75rem; align-items: end; }
    label { display: grid; gap: 0.35rem; color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 700; }
    label.checkbox { grid-template-columns: auto 1fr; align-items: center; align-content: center; min-height: 2.75rem; }
    .notice, .state-copy { display: grid; justify-items: start; gap: 0.75rem; padding: 1rem; }
    .notice { border-left: 3px solid var(--courier-warning); background: var(--courier-workbench-surface); }
    .state-brief, .loading { min-height: 9rem; } .loading { display: grid; place-items: center; }
    @media (max-width: 64rem) { form { grid-template-columns: repeat(2, minmax(10rem, 1fr)); } }
    @media (max-width: 44rem) {
      header { align-items: flex-start; padding: 1rem 0; } nav { margin-left: auto; }
      .page-head, .operations { grid-template-columns: 1fr; }
      .metrics { grid-template-columns: 1fr; }
      .metric + .metric { border-top: 1px solid var(--courier-workbench-line); border-left: 0; }
      .navigator { max-height: 18rem; border-right: 0; border-bottom: 1px solid var(--courier-workbench-line); }
      form, dl { grid-template-columns: 1fr; } dt { margin-top: 0.3rem; }
    }
  `];

  private locale: Locale = browserLocale();
  private snapshot?: Snapshot;
  private failed = false;
  private conflict = false;
  private selectedDeliveryId?: string;
  private unsubscribe?: () => void;

  connectedCallback(): void {
    super.connectedCallback();
    this.unsubscribe = subscribeSnapshots((snapshot) => {
      this.applySnapshot(snapshot);
    });
    void this.refresh();
  }

  disconnectedCallback(): void {
    this.unsubscribe?.();
    super.disconnectedCallback();
  }

  private applySnapshot(snapshot: Snapshot): void {
    this.snapshot = snapshot;
    this.selectedDeliveryId = reconcileDeliverySelection(snapshot, this.selectedDeliveryId);
    this.failed = false;
  }

  async refresh(): Promise<void> {
    this.failed = false;
    this.conflict = false;
    try {
      this.applySnapshot(await loadServers());
    } catch {
      this.failed = true;
    }
  }

  setLocale(event: CustomEvent<Locale>): void {
    this.locale = event.detail;
  }

  selectDelivery(id: string): void {
    this.selectedDeliveryId = id;
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
    return html`<article>
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
      </article>`;
  }

  private navigatorServer(item: Server) {
    return html`<section class="server-group">
      <div class="server-head">
        <div class="server-title"><span class="server-id">${item.id}</span><span class="bind">${item.bind}</span><courier-status tone=${item.status === "live" ? "signal" : "danger"}>${this.t(item.status)}</courier-status></div>
        <courier-button @click=${() => this.stop("servers", item.id)}>${this.t("stopServer")}</courier-button>
      </div>
      ${item.deliveries.map((delivery) => html`<button class="delivery-nav" type="button" aria-pressed=${delivery.id === this.selectedDeliveryId ? "true" : "false"} @click=${() => this.selectDelivery(delivery.id)}><strong>${delivery.route}</strong><span class="delivery-id">${delivery.id}</span></button>`)}
    </section>`;
  }

  render() {
    const servers = this.snapshot?.servers ?? [];
    const deliveries = servers.reduce((count, server) => count + server.deliveries.length, 0);
    const confirmed = servers.reduce((serverTotal, server) => serverTotal + server.deliveries.reduce((deliveryTotal, delivery) => deliveryTotal + delivery.counters.confirmed, 0), 0);
    const selected = servers.flatMap((server) => server.deliveries).find((delivery) => delivery.id === this.selectedDeliveryId);
    return html`
      <courier-scene class="page-scene" eager .source=${adminOperationsSource} .mobileSource=${adminOperationsMobileSource}></courier-scene>
      <main>
        <header>
          <courier-brand></courier-brand>
          <nav><courier-button @click=${this.refresh}>${this.t("refresh")}</courier-button><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></nav>
        </header>
        <div class="workspace">
          <div class="page-head"><div><span class="eyebrow">${this.t("eyebrow")}</span><h1>${this.t("title")}</h1><p class="intro">${this.t("intro")}</p></div><courier-status tone=${this.failed ? "danger" : "signal"}>${this.t(this.failed ? "unreachable" : "live")}</courier-status></div>
          <courier-workbench class="registry-workbench" .heading=${this.t("eyebrow")} .status=${this.failed ? this.t("unreachable") : this.t("live")}><div class="metrics">
            <div class="metric"><strong>${servers.length}</strong><span>${this.t("serversMetric")}</span></div>
            <div class="metric"><strong>${deliveries}</strong><span>${this.t("deliveriesMetric")}</span></div>
            <div class="metric"><strong>${confirmed}</strong><span>${this.t("confirmedMetric")}</span></div>
          </div></courier-workbench>
          ${this.failed ? html`<courier-workbench class="state-brief" .heading=${this.t("unreachable")} status="request failed"><div class="state-copy"><p role="alert">${this.t("failed")}</p><courier-button @click=${this.refresh}>${this.t("retry")}</courier-button></div></courier-workbench>` : nothing}
          ${this.conflict ? html`<div class="notice"><p role="alert">${this.t("conflict")}</p><courier-button @click=${this.refresh}>${this.t("refresh")}</courier-button></div>` : nothing}
          ${this.snapshot ? servers.length === 0 ? html`<courier-workbench class="state-brief" .heading=${this.t("live")} status="idle"><div class="state-copy"><p>${this.t("empty")}</p></div></courier-workbench>` : html`<courier-workbench class="operations" .heading=${this.t("deliveries")} .status=${selected?.state ?? this.t("live")}><nav class="navigator" aria-label=${this.t("deliveries")}>${servers.map((server) => this.navigatorServer(server))}</nav><div class="inspector">${selected ? this.delivery(selected) : html`<div class="state-copy"><p>${this.t("empty")}</p></div>`}</div></courier-workbench>` : !this.failed ? html`<courier-workbench class="loading" .heading=${this.t("eyebrow")} status="running"><courier-status>${this.t("loading")}</courier-status></courier-workbench>` : nothing}
        </div>
      </main>
    `;
  }
}

if (!customElements.get("courier-admin-app")) {
  customElements.define("courier-admin-app", CourierAdminApp);
}
