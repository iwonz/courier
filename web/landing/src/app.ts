import { LitElement, html, nothing, type PropertyValues } from "lit";
import {
  browserLocale,
  browserThemeState,
  cubicBezierPath,
  defineCourierElements,
  type BrandIconName,
  type IconName,
  type Locale,
  type ThemeState,
} from "@courier/ui";
import {
  relayCliMobileSource,
  relayCliSource,
  relayHeroMobileSource,
  relayHeroSource,
  relayInstallMobileSource,
  relayInstallSource,
  relayRoutingMobileSource,
  relayRoutingSource,
} from "@courier/ui/relay-landing";
import { landingText, type LandingMessage } from "./catalog";
import { contractData, type LandingCommand, type LandingFlag, type LandingRoute } from "./contract";
import { landingStyles } from "./landing-styles";

defineCourierElements();

interface InstallChannel {
  readonly name: string;
  readonly command: string;
  readonly icon: BrandIconName;
}

export interface RoutePair {
  readonly source: string;
  readonly destination: string;
  readonly routeName: string;
  readonly allowedFlags: readonly string[];
}

const primaryInstall = "curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh";

export const installs: readonly InstallChannel[] = [
  { name: "curl", command: primaryInstall, icon: "curl" },
  { name: "wget", command: "wget -qO- https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh", icon: "wget" },
  { name: "PowerShell", command: "irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex", icon: "powershell" },
  { name: "npm", command: "npm install --global @iwonz/courier", icon: "npm" },
  { name: "npx", command: "npx @iwonz/courier --help", icon: "npx" },
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

export function commandFlags(commandName: string, compatibleOnly: boolean, commands: readonly LandingCommand[], flags: readonly LandingFlag[]): LandingFlag[] {
  if (!commandName || !compatibleOnly) return [...flags];
  const names = commands.find((command) => command.name === commandName)?.flags ?? [];
  return flags.filter((flag) => names.includes(flag.name));
}

const routePairs = expandRoutePairs(contractData.routes);
const initialRoute = routePairs.find((pair) => pair.source === "local" && pair.destination === "ssh")!;
const sourceEndpoints = [...new Set(routePairs.map((pair) => pair.source))];
const destinationEndpoints = [...new Set(routePairs.map((pair) => pair.destination))];
const heroPath = cubicBezierPath({ x: 54, y: 70 }, { x: 91, y: 38 });

export class CourierLandingApp extends LitElement {
  static properties = {
    locale: { state: true },
    selectedSource: { state: true },
    selectedDestination: { state: true },
    activeInstall: { state: true },
    activeSection: { state: true },
    selectedCommand: { state: true },
    compatibleOnly: { state: true },
  };

  static styles = landingStyles;

  private locale: Locale = browserLocale();
  private selectedSource = initialRoute.source;
  private selectedDestination = initialRoute.destination;
  private activeInstall = installs[0]!.name;
  private activeSection = "";
  private selectedCommand = "";
  private compatibleOnly = true;
  private theme?: ThemeState;
  private resizeObserver?: ResizeObserver;
  private sectionObserver?: IntersectionObserver;
  private readonly sectionRatios = new Map<string, number>();

  connectedCallback(): void {
    super.connectedCallback();
    this.theme = browserThemeState();
  }

  disconnectedCallback(): void {
    this.theme?.destroy();
    this.resizeObserver?.disconnect();
    this.sectionObserver?.disconnect();
    super.disconnectedCallback();
  }

  protected firstUpdated(): void {
    this.measureHeader();
    this.measureRouteConnector();
    const masthead = this.renderRoot.querySelector<HTMLElement>(".masthead-wrap");
    const controls = this.renderRoot.querySelector<HTMLElement>(".route-controls");
    if (globalThis.ResizeObserver && masthead && controls) {
      this.resizeObserver = new ResizeObserver(() => {
        this.measureHeader();
        this.measureRouteConnector();
      });
      this.resizeObserver.observe(masthead);
      this.resizeObserver.observe(controls);
    }
    if (globalThis.IntersectionObserver) {
      this.sectionObserver = new IntersectionObserver((entries) => this.observeSections(entries), { threshold: [0.25, 0.5, 0.75] });
      this.renderRoot.querySelectorAll<HTMLElement>("section").forEach((section) => this.sectionObserver?.observe(section));
    }
  }

  protected updated(changed: PropertyValues<this>): void {
    if (changed.has("selectedSource") || changed.has("selectedDestination")) this.measureRouteConnector();
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

  private measureHeader(): void {
    const masthead = this.renderRoot.querySelector<HTMLElement>(".masthead-wrap");
    if (!masthead) return;
    const height = Math.ceil(masthead.getBoundingClientRect().height);
    if (height > 0) this.style.setProperty("--masthead-height", `${height}px`);
  }

  private observeSections(entries: readonly IntersectionObserverEntry[]): void {
    for (const entry of entries) this.sectionRatios.set(entry.target.id, entry.isIntersecting ? entry.intersectionRatio : 0);
    const visible = [...this.sectionRatios.entries()].sort((left, right) => right[1] - left[1])[0];
    this.activeSection = visible && visible[1] > 0 && visible[0] !== "hero" ? visible[0] : "";
  }

  private measureRouteConnector(): void {
    const controls = this.renderRoot.querySelector<HTMLElement>(".route-controls");
    const source = this.renderRoot.querySelector<HTMLElement>(`.source-endpoints [data-endpoint="${this.selectedSource}"]`);
    const destination = this.renderRoot.querySelector<HTMLElement>(`.destination-endpoints [data-endpoint="${this.selectedDestination}"]`);
    const connector = this.renderRoot.querySelector<SVGElement>(".route-connector");
    const connectorPath = connector?.querySelector("path");
    if (!controls || !source || !destination || !connector || !connectorPath) return;
    const bounds = controls.getBoundingClientRect();
    if (bounds.width <= 0 || bounds.height <= 0) return;
    const sourceBounds = source.getBoundingClientRect();
    const destinationBounds = destination.getBoundingClientRect();
    const path = cubicBezierPath(
      { x: sourceBounds.right - bounds.left, y: sourceBounds.top + sourceBounds.height / 2 - bounds.top },
      { x: destinationBounds.left - bounds.left, y: destinationBounds.top + destinationBounds.height / 2 - bounds.top },
    );
    const viewBox = `0 0 ${bounds.width} ${bounds.height}`;
    if (path === connectorPath.getAttribute("d") && viewBox === connector.getAttribute("viewBox")) return;
    connector.setAttribute("viewBox", viewBox);
    connectorPath.setAttribute("d", path);
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
    const selected = routePairs.find((pair) => pair.source === this.selectedSource && pair.destination === destination)!;
    this.selectedDestination = selected.destination;
  }

  private chooseInstall(event: Event): void {
    this.activeInstall = (event.currentTarget as HTMLElement).dataset.channel!;
  }

  private chooseCommand(event: Event): void {
    const command = (event.currentTarget as HTMLElement).dataset.command!;
    if (this.selectedCommand === command) {
      this.selectedCommand = "";
      this.compatibleOnly = true;
      return;
    }
    this.selectedCommand = command;
  }

  private setCompatibility(event: CustomEvent<boolean>): void {
    this.compatibleOnly = event.detail;
  }

  private navigate(event: MouseEvent): void {
    event.preventDefault();
    const hash = (event.currentTarget as HTMLAnchorElement).getAttribute("href")!;
    this.renderRoot.querySelector<HTMLElement>(hash)?.scrollIntoView({ block: "start" });
    globalThis.history.pushState(null, "", hash);
  }

  private renderEndpoint(name: string, side: "source" | "destination") {
    const valid = side === "source" || routePairs.some((pair) => pair.source === this.selectedSource && pair.destination === name);
    const selected = side === "source" ? name === this.selectedSource : name === this.selectedDestination;
    const handler = side === "source" ? this.chooseSource : this.chooseDestination;
    return html`<button
      type="button"
      class="endpoint ${selected ? "selected" : ""}"
      data-endpoint=${name}
      aria-pressed=${String(selected)}
      aria-disabled=${String(!valid)}
      ?disabled=${!valid}
      @click=${handler}
    ><courier-icon name=${endpointIcon(name)}></courier-icon><span>${this.displayEndpoint(name)}</span></button>`;
  }

  render() {
    const selected = routePairs.find((pair) => pair.source === this.selectedSource && pair.destination === this.selectedDestination)!;
    const flags = applicableFlags(selected, contractData.flags);
    const install = installs.find((channel) => channel.name === this.activeInstall) ?? installs[0]!;
    const visibleFlags = commandFlags(this.selectedCommand, this.compatibleOnly, contractData.commands, contractData.flags);
    return html`
      <header class="masthead-wrap">
        <div class="masthead shell">
          <div class="header-left">
            <a class="brand-link" href="#hero" @click=${this.navigate}><courier-brand product=${this.t("brandProduct")}></courier-brand></a>
            <nav aria-label="Courier">
              ${(["routes", "install", "cli"] as const).map((section) => html`<a href=${`#${section}`} aria-current=${this.activeSection === section ? "page" : nothing} @click=${this.navigate}>${this.t(section === "routes" ? "routeShort" : section === "install" ? "installShort" : "cliShort")}</a>`)}
            </nav>
          </div>
          <div class="header-actions">
            <a class="github-link" href="https://github.com/iwonz/courier" aria-label=${this.t("githubLabel")}><courier-icon name="github"></courier-icon></a>
            <div class="preferences"><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></div>
          </div>
        </div>
      </header>

      <main>
        <section id="hero" class="slide hero-slide">
          <courier-scene eager .source=${relayHeroSource} .mobileSource=${relayHeroMobileSource}></courier-scene>
          <div class="hero-shade" aria-hidden="true"></div>
          <svg class="hero-route" viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true"><path d=${heroPath}></path><path class="signal" d=${heroPath}></path><circle cx="54" cy="70" r="0.8"></circle><circle cx="91" cy="38" r="0.8"></circle></svg>
          <span class="hero-node source">Source</span><span class="hero-node destination">Destination</span>
          <div class="hero shell"><div class="hero-copy"><h1>${this.t("title")}</h1><p class="tagline">${this.t("tagline")}</p><p class="subline">${this.t("subline")}</p></div></div>
        </section>

        <section id="routes" class="slide">
          <courier-scene .source=${relayRoutingSource} .mobileSource=${relayRoutingMobileSource}></courier-scene>
          <div class="slide-shell shell">
            <div class="section-heading"><div><span class="section-index">01 / Routing</span><h2>${this.t("routeTitle")}</h2><p class="intro">${this.t("routeIntro")}</p></div></div>
            <div class="route-explorer"><div class="route-interface">
              <div class="route-controls">
                <svg class="route-connector" viewBox="0 0 1 1" preserveAspectRatio="none" aria-hidden="true"><path d=""></path></svg>
                <div class="endpoint-group source-endpoints"><span class="label">${this.t("sourceLabel")}</span>${sourceEndpoints.map((name) => this.renderEndpoint(name, "source"))}</div>
                <div class="endpoint-group destination-endpoints"><span class="label">${this.t("destinationLabel")}</span>${destinationEndpoints.map((name) => this.renderEndpoint(name, "destination"))}</div>
              </div>
              <div class="route-readout" aria-live="polite"><span class="label">${selected.routeName} · ${this.t("routeExample")}</span><code class="command-shape">courier <b>${this.t("from")}</b> ${endpointExample(selected.source, "source")} <b>${this.t("to")}</b> ${endpointExample(selected.destination, "destination")}</code><span class="label">${this.t("allowed")}</span><div class="flag-list">${flags.map((flag) => html`<span class="flag">${flag.syntax}</span>`)}</div></div>
            </div></div>
          </div>
        </section>

        <section id="install" class="slide">
          <courier-scene .source=${relayInstallSource} .mobileSource=${relayInstallMobileSource}></courier-scene>
          <div class="slide-shell shell">
            <div class="section-heading"><div><span class="section-index">02 / Distribution</span><h2>${this.t("install")}</h2><p class="intro">${this.t("installIntro")}</p></div></div>
            <div class="install-board"><div class="install-interface">
              <span class="label install-label">${this.t("chooseChannel")}</span>
              <div class="install-channels" role="list">${installs.map((channel) => html`<button type="button" class="install-channel ${channel.name === install.name ? "selected" : ""}" data-channel=${channel.name} aria-pressed=${String(channel.name === install.name)} @click=${this.chooseInstall}><courier-brand-icon name=${channel.icon}></courier-brand-icon><span>${channel.name}</span></button>`)}</div>
              <div class="install-readout" aria-live="polite"><div class="install-readout-head"><courier-brand-icon name=${install.icon}></courier-brand-icon><h3>${install.name}</h3></div><pre tabindex="0"><code>${install.command}</code></pre></div>
              <div class="install-actions">
                <a class="download-channel" href="https://github.com/iwonz/courier/releases/latest"><courier-icon name="package"></courier-icon><div><h3>${this.t("packages")}</h3><p>${this.t("packagesDetail")}</p><span class="brand-cloud" aria-hidden="true">${(["linux", "ubuntu", "debian", "arch-linux", "manjaro", "fedora", "red-hat", "alpine-linux"] as BrandIconName[]).map((name) => html`<courier-brand-icon name=${name}></courier-brand-icon>`)}</span></div></a>
                <a class="download-channel" href="https://github.com/iwonz/courier/releases/latest"><courier-icon name="download"></courier-icon><div><h3>${this.t("direct")}</h3><p>${this.t("directDetail")}</p></div></a>
              </div>
            </div></div>
          </div>
        </section>

        <section id="cli" class="slide">
          <courier-scene .source=${relayCliSource} .mobileSource=${relayCliMobileSource}></courier-scene>
          <div class="slide-shell shell">
            <div class="section-heading"><div><span class="section-index">03 / Interface</span><h2>${this.t("cliTitle")}</h2><p class="intro">${this.t("cliIntro")}</p></div></div>
            <div class="reference">
              <div class="reference-column commands"><div class="reference-head"><span class="label">${this.t("commandLabel")}</span></div><div class="reference-list">${contractData.commands.map((command) => html`<button type="button" class="command-row" data-command=${command.name} aria-pressed=${String(command.name === this.selectedCommand)} @click=${this.chooseCommand}><code>${command.usage}</code></button>`)}</div></div>
              <div class="reference-column options"><div class="reference-head"><span class="label">${this.t("optionsLabel")}</span><courier-checkbox .checked=${this.compatibleOnly} ?disabled=${!this.selectedCommand} .label=${this.t("compatibleOnly")} @courier-checkbox-change=${this.setCompatibility}></courier-checkbox></div><div class="reference-list">${visibleFlags.length ? visibleFlags.map((flag) => html`<article class="option-row"><div class="option-head"><code>${flag.syntax}</code><span class="kind">${this.t("optionDefault")}: ${flag.default}</span></div><p class="option-meta">${this.t("repeatable")}: ${flag.repeatable ? this.t("yes") : this.t("no")} · ${this.t("applies")}: ${flag.appliesTo.join(", ")}</p></article>`) : html`<p class="empty-state">${this.t("noCompatibleOptions")}</p>`}</div></div>
            </div>
          </div>
        </section>
      </main>
    `;
  }
}

if (!customElements.get("courier-landing-app")) customElements.define("courier-landing-app", CourierLandingApp);
