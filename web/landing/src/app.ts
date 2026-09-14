import { LitElement, css, html } from "lit";
import {
  browserLocale,
  browserThemeState,
  defineCourierElements,
  type IconName,
  type Locale,
  type ThemeState,
} from "@courier/ui";
import { relayDispatchSource, relayInstallSource, relayRoutingSource } from "@courier/ui/relay-landing";
import { landingText, type LandingMessage } from "./catalog";
import { contractData, type LandingFlag, type LandingRoute } from "./contract";

defineCourierElements();

interface InstallChannel {
  readonly name: string;
  readonly command: string;
  readonly icon: IconName;
}

export interface RoutePair {
  readonly source: string;
  readonly destination: string;
  readonly routeName: string;
  readonly allowedFlags: readonly string[];
}

const primaryInstall = "curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh";

const installs: readonly InstallChannel[] = [
  { name: "wget", command: "wget -qO- https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh", icon: "download" },
  { name: "PowerShell", command: "irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex", icon: "windows" },
  { name: "npm", command: "npm install --global @iwonz/courier", icon: "npm" },
  { name: "npx", command: "npx @iwonz/courier --help", icon: "terminal" },
  { name: "Yarn", command: "yarn dlx @iwonz/courier --help", icon: "yarn" },
  { name: "pnpm", command: "pnpm dlx @iwonz/courier --help", icon: "pnpm" },
  { name: "Homebrew", command: "brew tap iwonz/courier https://github.com/iwonz/courier && brew install --cask iwonz/courier/courier", icon: "homebrew" },
  { name: "Scoop", command: "scoop bucket add courier https://github.com/iwonz/courier && scoop install courier/courier", icon: "scoop" },
];

const endpointSamples: Readonly<Record<string, Readonly<Record<"source" | "destination", string>>>> = {
  local: { source: "./project", destination: "./backup/" },
  ssh: { source: "relay@host:/srv/source", destination: "relay@host:/srv/destination/" },
  web: { source: "web://", destination: "web://" },
  webhook: { source: "webhook://", destination: "webhook://" },
  http: { source: "https://api.example.test/source", destination: "https://api.example.test/upload" },
};

const endpointMessages: Readonly<Record<string, LandingMessage>> = {
  local: "endpointLocal",
  ssh: "endpointSsh",
  web: "endpointWeb",
  webhook: "endpointWebhook",
  http: "endpointHttp",
};

const endpointIcons: Readonly<Record<string, IconName>> = {
  local: "folder",
  ssh: "server",
  web: "download",
  webhook: "upload",
  http: "upload",
};

export function expandRoutePairs(routes: readonly LandingRoute[]): RoutePair[] {
  return routes.flatMap((route) => route.source.flatMap((source) => route.destination.map((destination) => ({
    source,
    destination,
    routeName: route.name,
    allowedFlags: route.allowedFlags,
  }))));
}

export function endpointExample(name: string, side: "source" | "destination"): string {
  return endpointSamples[name]?.[side] ?? name;
}

export function endpointMessage(name: string): LandingMessage | undefined {
  return endpointMessages[name];
}

export function endpointIcon(name: string): IconName {
  return endpointIcons[name] ?? "route";
}

export function applicableFlags(pair: RoutePair, flags: readonly LandingFlag[]): LandingFlag[] {
  const exactRoute = `${pair.source}-to-${pair.destination}`;
  return flags.filter((flag) => pair.allowedFlags.includes(flag.name) && (flag.appliesTo.includes(pair.routeName) || flag.appliesTo.includes(exactRoute)));
}

const routePairs = expandRoutePairs(contractData.routes);
const initialRoute = routePairs.find((pair) => pair.source === "local" && pair.destination === "ssh")!;
const sourceEndpoints = [...new Set(routePairs.map((pair) => pair.source))];
const destinationEndpoints = [...new Set(routePairs.map((pair) => pair.destination))];

export class CourierLandingApp extends LitElement {
  static properties = {
    locale: { state: true },
    selectedSource: { state: true },
    selectedDestination: { state: true },
  };

  static styles = css`
    :host {
      display: block;
      min-height: 100vh;
      color: var(--courier-color-text);
      background-color: var(--courier-color-canvas);
      background-image: linear-gradient(var(--courier-color-grid) 1px, transparent 1px), linear-gradient(90deg, var(--courier-color-grid) 1px, transparent 1px);
      background-size: 2.5rem 2.5rem;
      font-family: var(--courier-font-sans);
    }
    * { box-sizing: border-box; }
    a { color: inherit; }
    button { font: inherit; }
    code, pre { font-family: var(--courier-font-mono); }
    h1, h2, h3, p, pre { margin: 0; }
    .shell { width: min(78rem, calc(100% - 2rem)); margin: 0 auto; }
    .masthead { position: relative; z-index: 4; display: flex; min-height: 5rem; align-items: center; justify-content: space-between; gap: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .brand-link { min-width: 0; text-decoration: none; }
    nav, .preferences, .actions, .section-heading { display: flex; align-items: center; gap: 0.7rem; }
    nav { flex-wrap: wrap; justify-content: flex-end; }
    nav > a { color: var(--courier-color-muted); font-size: 0.875rem; font-weight: 750; text-decoration: none; }
    nav > a:hover { color: var(--courier-color-text); }
    .github-link { display: inline-flex; width: 2.45rem; height: 2.45rem; align-items: center; justify-content: center; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface); }
    .github-link courier-icon { width: 1.25rem; height: 1.25rem; }
    .preferences { padding-left: 0.25rem; border-left: 1px solid var(--courier-color-border); }
    .hero { display: grid; grid-template-columns: minmax(0, 1.08fr) minmax(20rem, 0.92fr); min-height: 43rem; align-items: center; gap: clamp(2rem, 7vw, 6rem); padding: clamp(4rem, 9vw, 8rem) 0; }
    .hero-copy { position: relative; z-index: 2; display: grid; min-width: 0; gap: 1.25rem; }
    h1 { max-width: 48rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(3.8rem, 9vw, 7.75rem); font-weight: 850; letter-spacing: -0.075em; line-height: 0.87; }
    h2 { max-width: 48rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(2.25rem, 5vw, 4.5rem); font-weight: 820; letter-spacing: -0.06em; line-height: 0.95; }
    h3 { font-size: 1rem; letter-spacing: -0.02em; }
    .tagline { max-width: 42rem; font-size: clamp(1.35rem, 2.6vw, 2.15rem); font-weight: 650; letter-spacing: -0.035em; line-height: 1.12; }
    .subline, .intro { max-width: 46rem; color: var(--courier-color-muted); font-size: 1.025rem; line-height: 1.65; }
    .button { display: inline-flex; min-height: 2.875rem; align-items: center; justify-content: center; padding: 0 1rem; border: 1px solid var(--courier-color-accent); border-radius: var(--courier-radius-sm); color: var(--courier-color-accent-ink); background: var(--courier-color-accent); font-weight: 800; text-decoration: none; }
    .button.secondary { border-color: var(--courier-color-border-strong); color: var(--courier-color-text); background: var(--courier-color-surface); }
    .hero-visual { position: relative; min-height: 34rem; overflow: hidden; border: 1px solid #363b33; border-radius: var(--courier-radius-lg); color: var(--courier-paper-50); background: var(--courier-graphite-900); box-shadow: var(--courier-shadow); }
    .hero-visual::before { content: ""; position: absolute; inset: 0; opacity: 0.18; background-image: linear-gradient(rgb(243 244 233 / 0.18) 1px, transparent 1px), linear-gradient(90deg, rgb(243 244 233 / 0.18) 1px, transparent 1px); background-size: 2rem 2rem; }
    .hero-status { position: absolute; z-index: 2; top: 1.25rem; left: 1.25rem; }
    .hero-visual courier-status { --courier-color-muted: #d4ff45; }
    .hero-visual courier-mascot { position: absolute; z-index: 1; right: -10%; bottom: -7%; width: min(112%, 34rem); }
    .route-card { position: absolute; z-index: 2; right: 1rem; bottom: 1rem; left: 1rem; padding: 1rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-md); background: rgb(21 23 20 / 0.88); backdrop-filter: blur(8px); }
    .route-card .label { display: block; margin-bottom: 0.75rem; color: #aeb5a7; }
    .hero-visual courier-route { --courier-color-text: #f3f4e9; --courier-color-surface: #232720; --courier-color-border: #4d5547; --courier-color-border-strong: #aeb5a7; --courier-color-accent: #d4ff45; }
    .quick-command { display: grid; grid-template-columns: auto minmax(0, 1fr); max-width: 48rem; min-width: 0; gap: 0.75rem; padding: 0.9rem 1rem; border: 1px solid var(--courier-color-border-strong); border-radius: var(--courier-radius-sm); color: var(--courier-paper-50); background: var(--courier-graphite-900); box-shadow: var(--courier-shadow); }
    .quick-command span { color: var(--courier-signal); }
    .quick-command code { min-width: 0; overflow-wrap: anywhere; white-space: normal; }
    main { display: grid; gap: clamp(5rem, 10vw, 9rem); padding: 0 0 clamp(5rem, 10vw, 9rem); }
    section { display: grid; min-width: 0; gap: 1.5rem; scroll-margin-top: 2rem; }
    .section-heading { justify-content: space-between; align-items: end; padding-bottom: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .section-heading > div { display: grid; min-width: 0; gap: 0.6rem; }
    .section-index, .label, .kind { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 750; letter-spacing: 0.09em; text-transform: uppercase; }
    .install-board { display: grid; grid-template-columns: minmax(18rem, 0.7fr) minmax(0, 1.3fr); overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .install-scene { position: relative; min-width: 0; min-height: 35rem; overflow: hidden; color: var(--courier-paper-50); background: var(--courier-graphite-900); }
    .install-scene courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; }
    .install-scene courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; }
    .primary-install { position: absolute; right: 1rem; bottom: 1rem; left: 1rem; display: grid; min-width: 0; gap: 0.55rem; padding: 1rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-md); background: rgb(16 18 15 / 0.9); }
    .primary-install .label { color: var(--courier-signal); }
    .primary-install pre { max-width: 100%; padding: 0.7rem 0; overflow-x: auto; color: var(--courier-paper-50); }
    .primary-install code { overflow-wrap: anywhere; white-space: pre-wrap; }
    .channels { display: grid; min-width: 0; align-content: start; }
    .channels-title { padding: 1rem 1.15rem; border-bottom: 1px solid var(--courier-color-border); }
    .channel { display: grid; grid-template-columns: 2rem minmax(5.5rem, 7.5rem) minmax(0, 1fr); min-width: 0; align-items: center; gap: 0.75rem; padding: 0.72rem 1.15rem; border-bottom: 1px solid var(--courier-color-border); }
    .channel:last-child { border-bottom: 0; }
    .channel > courier-icon { width: 1.25rem; height: 1.25rem; color: var(--courier-color-accent-ink); }
    .channel h3 { min-width: 0; }
    .channel pre { width: 100%; min-width: 0; max-width: 100%; padding: 0.55rem 0.65rem; overflow-x: auto; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-sm); background: var(--courier-color-surface); }
    .channel code { overflow-wrap: anywhere; white-space: pre-wrap; }
    .download-channels { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
    .download-channel { display: grid; grid-template-columns: auto minmax(0, 1fr) auto; align-items: start; gap: 0.85rem; padding: 1rem; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface-raised); }
    .download-channel courier-icon { width: 1.5rem; height: 1.5rem; color: var(--courier-color-accent-ink); }
    .download-channel div { display: grid; gap: 0.4rem; }
    .download-channel p { color: var(--courier-color-muted); line-height: 1.5; }
    .download-channel a { font-weight: 750; }
    .route-explorer { display: grid; grid-template-columns: minmax(15rem, 0.55fr) minmax(0, 1.45fr); overflow: hidden; border: 1px solid #363b33; border-radius: var(--courier-radius-lg); color: var(--courier-paper-50); background: var(--courier-graphite-900); box-shadow: var(--courier-shadow); }
    .route-scene { position: relative; min-width: 0; min-height: 35rem; overflow: hidden; border-right: 1px solid #4d5547; }
    .route-scene courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; }
    .route-scene courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; }
    .route-interface { display: grid; min-width: 0; align-content: center; gap: 1.25rem; padding: clamp(1.25rem, 4vw, 3rem); }
    .endpoint-columns { display: grid; grid-template-columns: minmax(0, 1fr) minmax(2.5rem, auto) minmax(0, 1fr); align-items: start; gap: 1rem; }
    .endpoint-group { display: grid; min-width: 0; gap: 0.55rem; }
    .endpoint-group .label { color: #aeb5a7; }
    .endpoint-word { align-self: center; padding-top: 2.2rem; color: var(--courier-signal); font-family: var(--courier-font-mono); font-weight: 850; text-align: center; }
    .endpoint { appearance: none; display: grid; grid-template-columns: auto minmax(0, 1fr); min-width: 0; min-height: 3rem; align-items: center; gap: 0.65rem; padding: 0.65rem 0.75rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-sm); color: #c7ccbf; background: #1c1f1b; text-align: left; cursor: pointer; transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), opacity var(--courier-duration) var(--courier-ease); }
    .endpoint courier-icon { width: 1.2rem; height: 1.2rem; }
    .endpoint.valid { border-color: #6d7567; }
    .endpoint.selected { border-color: var(--courier-signal); color: var(--courier-signal); background: #252a21; }
    .endpoint.invalid { opacity: 0.34; cursor: not-allowed; }
    .endpoint:focus-visible, .github-link:focus-visible, a:focus-visible { outline: 3px solid var(--courier-beak); outline-offset: 3px; }
    .route-readout { display: grid; min-width: 0; gap: 0.85rem; padding: 1rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-md); background: #10120f; }
    .route-readout code { display: block; max-width: 100%; overflow-wrap: anywhere; color: var(--courier-paper-50); font-size: clamp(0.76rem, 1.5vw, 0.95rem); white-space: normal; }
    .command-shape { color: #aeb5a7; }
    .command-shape b { color: var(--courier-signal); }
    .flag-list { display: flex; gap: 0.4rem; flex-wrap: wrap; }
    .flag { padding: 0.3rem 0.45rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-sm); color: #c7ccbf; background: #1c1f1b; font-family: var(--courier-font-mono); font-size: 0.7rem; }
    .reference { display: grid; grid-template-columns: minmax(15rem, 0.68fr) minmax(0, 1.32fr); overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-color-surface-raised); box-shadow: var(--courier-shadow); }
    .reference-column { min-width: 0; }
    .reference-column + .reference-column { border-left: 1px solid var(--courier-color-border); }
    .reference-label { display: block; padding: 0.8rem 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .command-row, .option-row { display: grid; min-width: 0; gap: 0.45rem; padding: 0.8rem 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .command-row:last-child, .option-row:last-child { border-bottom: 0; }
    .command-row code, .option-row code { overflow-wrap: anywhere; font-size: 0.78rem; }
    .option-head { display: flex; align-items: baseline; justify-content: space-between; gap: 0.75rem; }
    .option-meta { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.68rem; line-height: 1.5; }
    .reference-link { justify-self: start; font-weight: 750; }
    @media (max-width: 62rem) {
      .hero { grid-template-columns: 1fr; padding-top: 4rem; }
      .hero-visual { min-height: 31rem; }
      .install-board, .route-explorer { grid-template-columns: 1fr; }
      .install-scene, .route-scene { min-height: 0; aspect-ratio: 3 / 2; border-right: 0; border-bottom: 1px solid #4d5547; }
      .reference { grid-template-columns: 1fr; }
      .reference-column + .reference-column { border-top: 1px solid var(--courier-color-border); border-left: 0; }
    }
    @media (max-width: 44rem) {
      .masthead { display: grid; grid-template-columns: minmax(0, 1fr); align-items: start; padding: 1rem 0; }
      nav { width: 100%; justify-content: space-between; }
      .nav-text { display: none; }
      .preferences { margin-left: auto; }
      .hero { min-height: auto; padding: 4rem 0; }
      .hero-visual { min-height: 26rem; }
      .actions { flex-wrap: wrap; }
      .channel { grid-template-columns: 2rem minmax(0, 1fr); }
      .channel pre { grid-column: 1 / -1; }
      .download-channels { grid-template-columns: 1fr; }
      .download-channel { grid-template-columns: auto minmax(0, 1fr); }
      .download-channel a { grid-column: 2; }
      .endpoint-columns { grid-template-columns: minmax(0, 1fr) 2rem minmax(0, 1fr); gap: 0.5rem; }
      .endpoint { grid-template-columns: 1fr; justify-items: center; text-align: center; }
      .endpoint-word { padding-top: 2.1rem; }
      .route-interface { padding: 1rem; }
    }
  `;

  private locale: Locale = browserLocale();
  private selectedSource = initialRoute.source;
  private selectedDestination = initialRoute.destination;
  private theme?: ThemeState;

  connectedCallback(): void {
    super.connectedCallback();
    this.theme = browserThemeState();
  }

  disconnectedCallback(): void {
    this.theme?.destroy();
    super.disconnectedCallback();
  }

  setLocale(event: CustomEvent<Locale>): void {
    this.locale = event.detail;
  }

  private t(message: LandingMessage): string {
    return landingText(this.locale, message);
  }

  private displayEndpoint(name: string): string {
    const message = endpointMessage(name);
    return message ? this.t(message) : name;
  }

  private chooseSource(event: Event): void {
    const source = (event.currentTarget as HTMLElement).dataset.endpoint!;
    const current = routePairs.find((pair) => pair.source === source && pair.destination === this.selectedDestination);
    const selected = current ?? routePairs.find((pair) => pair.source === source)!;
    this.selectedSource = selected.source;
    this.selectedDestination = selected.destination;
  }

  private chooseDestination(event: Event): void {
    const destination = (event.currentTarget as HTMLElement).dataset.endpoint!;
    const selected = routePairs.find((pair) => pair.source === this.selectedSource && pair.destination === destination);
    if (selected) {
      this.selectedDestination = selected.destination;
    }
  }

  private renderEndpoint(name: string, side: "source" | "destination") {
    const valid = side === "source" || routePairs.some((pair) => pair.source === this.selectedSource && pair.destination === name);
    const selected = side === "source" ? name === this.selectedSource : name === this.selectedDestination;
    const handler = side === "source" ? this.chooseSource : this.chooseDestination;
    return html`<button
      type="button"
      class="endpoint ${valid ? "valid" : "invalid"} ${selected ? "selected" : ""}"
      data-endpoint=${name}
      aria-pressed=${String(selected)}
      aria-disabled=${String(!valid)}
      @pointerenter=${handler}
      @focus=${handler}
      @click=${handler}
    ><courier-icon name=${endpointIcon(name)}></courier-icon><span>${this.displayEndpoint(name)}</span></button>`;
  }

  render() {
    const selected = routePairs.find((pair) => pair.source === this.selectedSource && pair.destination === this.selectedDestination)!;
    const flags = applicableFlags(selected, contractData.flags);
    return html`
      <header id="top">
        <div class="shell">
          <div class="masthead">
            <a class="brand-link" href="#top"><courier-brand product=${this.t("brandProduct")}></courier-brand></a>
            <nav aria-label="Courier">
              <a class="nav-text" href="#install">${this.t("installShort")}</a>
              <a class="nav-text" href="#routes">${this.t("routeShort")}</a>
              <a class="nav-text" href="#cli">${this.t("cliShort")}</a>
              <a class="github-link" href="https://github.com/iwonz/courier" aria-label=${this.t("githubLabel")}><courier-icon name="github"></courier-icon></a>
              <div class="preferences">
                <courier-theme-selector .locale=${this.locale}></courier-theme-selector>
                <courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector>
              </div>
            </nav>
          </div>
          <div class="hero">
            <div class="hero-copy">
              <h1>${this.t("title")}</h1>
              <p class="tagline">${this.t("tagline")}</p>
              <p class="subline">${this.t("subline")}</p>
              <div class="quick-command"><span>$</span><code>${primaryInstall}</code></div>
              <div class="actions">
                <a class="button" href="#install">${this.t("install")}</a>
                <a class="button secondary" href="https://github.com/iwonz/courier/releases">${this.t("releases")}</a>
              </div>
            </div>
            <div class="hero-visual" aria-label=${this.t("heroRouteLabel")}>
              <div class="hero-status"><courier-status tone="signal">${this.t("heroStatus")}</courier-status></div>
              <courier-mascot eager alt="" .source=${relayDispatchSource}></courier-mascot>
              <div class="route-card"><span class="label">${this.t("heroRouteLabel")}</span><courier-route source="./project" destination="server:/srv/data/"></courier-route></div>
            </div>
          </div>
        </div>
      </header>

      <main class="shell">
        <section id="install">
          <div class="section-heading"><div><span class="section-index">01 / Dispatch</span><h2>${this.t("install")}</h2><p class="intro">${this.t("installIntro")}</p></div></div>
          <div class="install-board">
            <div class="install-scene">
              <courier-mascot alt="" .source=${relayInstallSource}></courier-mascot>
              <div class="primary-install"><span class="label">${this.t("primaryChannel")}</span><h3>curl</h3><pre tabindex="0"><code>${primaryInstall}</code></pre></div>
            </div>
            <div class="channels">
              <span class="channels-title label">${this.t("alternateChannels")}</span>
              ${installs.map((channel) => html`<article class="channel"><courier-icon name=${channel.icon}></courier-icon><h3>${channel.name}</h3><pre tabindex="0"><code>${channel.command}</code></pre></article>`)}
            </div>
          </div>
          <div class="download-channels">
            <article class="download-channel"><courier-icon name="linux"></courier-icon><div><h3>${this.t("packages")}</h3><p>${this.t("packagesDetail")}</p></div><a href="https://github.com/iwonz/courier/releases/latest">${this.t("downloads")}</a></article>
            <article class="download-channel"><courier-icon name="download"></courier-icon><div><h3>${this.t("direct")}</h3><p>${this.t("directDetail")}</p></div><a href="https://github.com/iwonz/courier/releases/latest">${this.t("downloads")}</a></article>
          </div>
        </section>

        <section id="routes">
          <div class="section-heading"><div><span class="section-index">02 / Routing</span><h2>${this.t("routeTitle")}</h2><p class="intro">${this.t("routeIntro")}</p></div></div>
          <div class="route-explorer">
            <div class="route-scene"><courier-mascot alt="" .source=${relayRoutingSource}></courier-mascot></div>
            <div class="route-interface">
              <div class="endpoint-columns">
                <div class="endpoint-group"><span class="label">${this.t("sourceLabel")}</span>${sourceEndpoints.map((name) => this.renderEndpoint(name, "source"))}</div>
                <span class="endpoint-word">→</span>
                <div class="endpoint-group"><span class="label">${this.t("destinationLabel")}</span>${destinationEndpoints.map((name) => this.renderEndpoint(name, "destination"))}</div>
              </div>
              <div class="route-readout" aria-live="polite">
                <span class="label">${selected.routeName} · ${this.t("routeExample")}</span>
                <code class="command-shape">courier <b>${this.t("from")}</b> ${endpointExample(selected.source, "source")} <b>${this.t("to")}</b> ${endpointExample(selected.destination, "destination")}</code>
                <span class="label">${this.t("allowed")}</span>
                <div class="flag-list">${flags.map((flag) => html`<span class="flag">${flag.syntax}</span>`)}</div>
              </div>
            </div>
          </div>
        </section>

        <section id="cli">
          <div class="section-heading"><div><span class="section-index">03 / Interface</span><h2>${this.t("cliTitle")}</h2><p class="intro">${this.t("cliIntro")}</p></div></div>
          <div class="reference">
            <div class="reference-column">
              <span class="reference-label label">${this.t("commandKind")}</span>
              ${contractData.commands.map((command) => html`<article class="command-row"><span class="kind">${command.system ? this.t("system") : this.t("product")}</span><code>${command.usage}</code></article>`)}
            </div>
            <div class="reference-column">
              <span class="reference-label label">${this.t("allowed")}</span>
              ${contractData.flags.map((flag) => html`<article class="option-row"><div class="option-head"><code>${flag.syntax}</code><span class="kind">${this.t("optionDefault")}: ${flag.default}</span></div><p class="option-meta">${this.t("repeatable")}: ${flag.repeatable ? this.t("yes") : this.t("no")} · ${this.t("applies")}: ${flag.appliesTo.join(", ")}</p></article>`)}
            </div>
          </div>
          <a class="reference-link" href="https://github.com/iwonz/courier/blob/main/docs/cli-reference.md">${this.t("inspectContract")}</a>
        </section>
      </main>
    `;
  }
}

if (!customElements.get("courier-landing-app")) {
  customElements.define("courier-landing-app", CourierLandingApp);
}
