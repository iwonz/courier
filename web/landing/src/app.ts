import { LitElement, html, nothing, type PropertyValues } from "lit";
import { browserLocale, cubicBezierPath, defineCourierElements, type BrandIconName, type IconName, type Locale } from "@courier/ui";
import { landingInstallMobileSource, landingInstallSource, landingReferenceMobileSource, landingReferenceSource, landingRouteMobileSource, landingRouteSource } from "@courier/ui/landing-scenes";
import { landingText, type LandingMessage } from "./catalog";
import { contractData, type LandingCommand, type LandingFlag, type LandingRoute } from "./contract";
import { landingStyles } from "./landing-styles";

defineCourierElements();

interface InstallChannel { readonly name: string; readonly command: string; readonly icon: BrandIconName; }
export interface RoutePair { readonly source: string; readonly destination: string; readonly routeName: string; readonly allowedFlags: readonly string[]; }
export interface EndpointPresentation { readonly label: LandingMessage; readonly description: LandingMessage; readonly icon: IconName; }

export const installs: readonly InstallChannel[] = [
  { name: "curl", command: "curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh", icon: "curl" },
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
  ssh: { source: "courier@host:/srv/source", destination: "courier@host:/srv/destination/" },
  web: { source: "web://", destination: "web://" },
  webhook: { source: "webhook://", destination: "webhook://" },
  http: { source: "https://api.example.test/source", destination: "https://api.example.test/upload" },
};

const endpointPresentations: Readonly<Record<"source" | "destination", Readonly<Record<string, EndpointPresentation>>>> = {
  source: {
    local: { label: "endpointLocal", description: "sourceLocalDescription", icon: "folder-out" },
    ssh: { label: "endpointRemote", description: "sourceRemoteDescription", icon: "server-out" },
    web: { label: "endpointWeb", description: "sourceWebDescription", icon: "browser-upload" },
    webhook: { label: "endpointWebHook", description: "sourceWebHookDescription", icon: "webhook-in" },
  },
  destination: {
    local: { label: "endpointLocal", description: "destinationLocalDescription", icon: "folder-in" },
    ssh: { label: "endpointRemote", description: "destinationRemoteDescription", icon: "server-in" },
    web: { label: "endpointWeb", description: "destinationWebDescription", icon: "browser-share" },
    http: { label: "endpointWebHook", description: "destinationWebHookDescription", icon: "webhook-out" },
  },
};

export function expandRoutePairs(routes: readonly LandingRoute[]): RoutePair[] {
  return routes.flatMap((route) => route.source.flatMap((source) => route.destination.map((destination) => ({ source, destination, routeName: route.name, allowedFlags: route.allowedFlags }))));
}
export function endpointExample(name: string, side: "source" | "destination"): string { return endpointSamples[name]?.[side] ?? name; }
export function endpointPresentation(name: string, side: "source" | "destination"): EndpointPresentation | undefined { return endpointPresentations[side][name]; }
export function endpointMessage(name: string, side: "source" | "destination"): LandingMessage | undefined { return endpointPresentation(name, side)?.label; }
export function endpointIcon(name: string, side: "source" | "destination"): IconName { return endpointPresentation(name, side)?.icon ?? "route"; }
export function applicableFlags(pair: RoutePair, flags: readonly LandingFlag[]): LandingFlag[] {
  const exactRoute = `${pair.source}-to-${pair.destination}`;
  return flags.filter((flag) => pair.allowedFlags.includes(flag.name) && (flag.appliesTo.includes(pair.routeName) || flag.appliesTo.includes(exactRoute)));
}
export function commandFlags(commandName: string, compatibleOnly: boolean, commands: readonly LandingCommand[], flags: readonly LandingFlag[]): LandingFlag[] {
  if (!commandName || !compatibleOnly) return [...flags];
  const names = commands.find((command) => command.name === commandName)?.flags ?? [];
  return flags.filter((flag) => names.includes(flag.name));
}
export async function copyText(text: string, clipboard: Pick<Clipboard, "writeText"> | undefined = globalThis.navigator.clipboard): Promise<boolean> {
  if (!clipboard) return false;
  try { await clipboard.writeText(text); return true; } catch { return false; }
}

const routePairs = expandRoutePairs(contractData.routes);
const initialRoute = routePairs.find((pair) => pair.source === "local" && pair.destination === "ssh")!;
const sourceEndpoints = [...new Set(routePairs.map((pair) => pair.source))];
const destinationEndpoints = [...new Set(routePairs.map((pair) => pair.destination))];

export class CourierLandingApp extends LitElement {
  static properties = { locale: { state: true }, selectedSource: { state: true }, selectedDestination: { state: true }, activeInstall: { state: true }, activeSection: { state: true }, selectedCommand: { state: true }, compatibleOnly: { state: true } };
  static styles = landingStyles;
  private locale: Locale = browserLocale();
  private selectedSource = initialRoute.source;
  private selectedDestination = initialRoute.destination;
  private activeInstall = installs[0]!.name;
  private activeSection = "";
  private selectedCommand = "";
  private compatibleOnly = true;
  private resizeObserver?: ResizeObserver;
  private sectionObserver?: IntersectionObserver;
  private readonly observedSections = new Set<string>();
  private headerHeight = 0;
  private readonly handleViewportResize = (): void => { this.observeSectionBand(); this.measureRouteConnector(); };

  connectedCallback(): void { super.connectedCallback(); globalThis.addEventListener("resize", this.handleViewportResize); }
  disconnectedCallback(): void { this.resizeObserver?.disconnect(); this.sectionObserver?.disconnect(); globalThis.removeEventListener("resize", this.handleViewportResize); super.disconnectedCallback(); }
  protected firstUpdated(): void {
    this.measureHeader(); this.measureRouteConnector();
    const masthead = this.renderRoot.querySelector<HTMLElement>(".masthead-wrap");
    const controls = this.renderRoot.querySelector<HTMLElement>(".route-controls");
    if (globalThis.ResizeObserver && masthead && controls) {
      this.resizeObserver = new ResizeObserver(() => { this.measureHeader(); this.measureRouteConnector(); });
      this.resizeObserver.observe(masthead); this.resizeObserver.observe(controls);
    }
  }
  protected updated(changed: PropertyValues<this>): void { if (changed.has("selectedSource") || changed.has("selectedDestination")) this.measureRouteConnector(); }
  setLocale(event: CustomEvent<Locale>): void { this.locale = event.detail; }
  private t(message: LandingMessage): string { return landingText(this.locale, message); }
  private displayEndpoint(name: string, side: "source" | "destination"): string { const message = endpointMessage(name, side); return message ? this.t(message) : name; }
  private describeEndpoint(name: string, side: "source" | "destination"): string { const message = endpointPresentation(name, side)?.description; return message ? this.t(message) : name; }
  private measureHeader(): void {
    const masthead = this.renderRoot.querySelector<HTMLElement>(".masthead-wrap");
    if (!masthead) return;
    const height = Math.ceil(masthead.getBoundingClientRect().height);
    if (height <= 0) return;
    this.style.setProperty("--masthead-height", `${height}px`);
    if (height !== this.headerHeight) { this.headerHeight = height; this.observeSectionBand(); }
  }
  private observeSections(entries: readonly IntersectionObserverEntry[]): void {
    for (const entry of entries) entry.isIntersecting && entry.intersectionRatio > 0 ? this.observedSections.add(entry.target.id) : this.observedSections.delete(entry.target.id);
    const visible = [...this.renderRoot.querySelectorAll<HTMLElement>("section")].reverse().find((section) => this.observedSections.has(section.id));
    if (visible) this.activeSection = visible.id === "route" ? "" : visible.id;
  }
  private observeSectionBand(): void {
    this.sectionObserver?.disconnect(); this.sectionObserver = undefined; this.observedSections.clear();
    if (!globalThis.IntersectionObserver) return;
    const viewportHeight = Math.max(this.headerHeight + 2, globalThis.innerHeight);
    const readingBandHeight = Math.max(2, Math.min(160, viewportHeight * 0.25));
    const bottomInset = Math.max(0, viewportHeight - this.headerHeight - readingBandHeight);
    this.sectionObserver = new IntersectionObserver((entries) => this.observeSections(entries), { rootMargin: `-${this.headerHeight}px 0px -${bottomInset}px 0px`, threshold: 0 });
    this.renderRoot.querySelectorAll<HTMLElement>("section").forEach((section) => this.sectionObserver?.observe(section));
  }
  private measureRouteConnector(): void {
    const controls = this.renderRoot.querySelector<HTMLElement>(".route-controls");
    const source = this.renderRoot.querySelector<HTMLElement>(`.source-endpoints [data-endpoint="${this.selectedSource}"]`);
    const destination = this.renderRoot.querySelector<HTMLElement>(`.destination-endpoints [data-endpoint="${this.selectedDestination}"]`);
    const connector = this.renderRoot.querySelector<SVGElement>(".route-connector");
    const pathNode = connector?.querySelector(".route-path");
    if (!controls || !source || !destination || !connector || !pathNode) return;
    const bounds = controls.getBoundingClientRect();
    if (bounds.width <= 0 || bounds.height <= 0) return;
    const sourceBounds = source.getBoundingClientRect();
    const destinationBounds = destination.getBoundingClientRect();
    const path = cubicBezierPath({ x: sourceBounds.right - bounds.left, y: sourceBounds.top + sourceBounds.height / 2 - bounds.top }, { x: destinationBounds.left - bounds.left, y: destinationBounds.top + destinationBounds.height / 2 - bounds.top });
    connector.setAttribute("viewBox", `0 0 ${bounds.width} ${bounds.height}`);
    pathNode.setAttribute("d", path);
    connector.querySelector(".route-signal")?.setAttribute("d", path);
  }
  private chooseSource(event: Event): void {
    const source = (event.currentTarget as HTMLElement).dataset.endpoint!;
    const selected = routePairs.find((pair) => pair.source === source && pair.destination === this.selectedDestination) ?? routePairs.find((pair) => pair.source === source)!;
    this.selectedSource = selected.source; this.selectedDestination = selected.destination;
  }
  private chooseDestination(event: Event): void { this.selectedDestination = (event.currentTarget as HTMLElement).dataset.endpoint!; }
  private chooseInstall(event: Event): void { this.activeInstall = (event.currentTarget as HTMLElement).dataset.channel!; }
  private chooseCommand(event: Event): void {
    const command = (event.currentTarget as HTMLElement).dataset.command!;
    if (this.selectedCommand === command) { this.selectedCommand = ""; this.compatibleOnly = true; return; }
    this.selectedCommand = command;
  }
  private setCompatibility(event: CustomEvent<boolean>): void { this.compatibleOnly = event.detail; }
  private navigate(event: MouseEvent): void {
    event.preventDefault();
    const hash = (event.currentTarget as HTMLAnchorElement).getAttribute("href")!;
    const reducedMotion = globalThis.matchMedia?.("(prefers-reduced-motion: reduce)").matches === true;
    this.renderRoot.querySelector<HTMLElement>(hash)?.scrollIntoView({ block: "start", behavior: reducedMotion ? "auto" : "smooth" });
    globalThis.history.pushState(null, "", hash);
  }
  private renderEndpoint(name: string, side: "source" | "destination") {
    const valid = side === "source" || routePairs.some((pair) => pair.source === this.selectedSource && pair.destination === name);
    const selected = side === "source" ? name === this.selectedSource : name === this.selectedDestination;
    return html`<button type="button" class="endpoint ${selected ? "selected" : ""}" data-endpoint=${name} aria-pressed=${String(selected)} aria-label=${this.displayEndpoint(name, side)} aria-description=${this.describeEndpoint(name, side)} ?disabled=${!valid} @click=${side === "source" ? this.chooseSource : this.chooseDestination}><span class="endpoint-terminal" aria-hidden="true"><courier-icon name=${endpointIcon(name, side)}></courier-icon></span><span>${this.displayEndpoint(name, side)}</span></button>`;
  }

  render() {
    const selected = routePairs.find((pair) => pair.source === this.selectedSource && pair.destination === this.selectedDestination)!;
    const flags = applicableFlags(selected, contractData.flags);
    const install = installs.find((channel) => channel.name === this.activeInstall) ?? installs[0]!;
    const visibleFlags = commandFlags(this.selectedCommand, this.compatibleOnly, contractData.commands, contractData.flags);
    const selectedCommand = contractData.commands.find((command) => command.name === this.selectedCommand);
    const routeCommand = `courier ${this.t("from")} ${endpointExample(selected.source, "source")} ${this.t("to")} ${endpointExample(selected.destination, "destination")}`;
    return html`
      <header class="masthead-wrap"><div class="masthead shell"><a class="brand-link" href="#route" @click=${this.navigate}><courier-brand></courier-brand></a><nav aria-label="Courier"><a href="#install" aria-current=${this.activeSection === "install" ? "page" : nothing} @click=${this.navigate}>${this.t("installShort")}</a><a href="#cli" aria-current=${this.activeSection === "cli" ? "page" : nothing} @click=${this.navigate}>${this.t("cliShort")}</a></nav><div class="header-actions"><courier-icon-link href="https://github.com/iwonz/courier" icon="github" target="_blank" rel="noopener noreferrer" .label=${this.t("githubLabel")}></courier-icon-link><courier-theme-selector .locale=${this.locale}></courier-theme-selector><courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector></div></div></header>
      <main>
        <section id="route" class="stage route-stage"><courier-scene eager .source=${landingRouteSource} .mobileSource=${landingRouteMobileSource}></courier-scene><div class="stage-shell route-shell shell"><h1>${this.t("title")}</h1><div class="route-instrument"><div class="route-controls"><svg class="route-connector" viewBox="0 0 1 1" preserveAspectRatio="none" aria-hidden="true"><path class="route-path" d=""></path><path class="route-signal" d=""></path></svg><div class="endpoint-group source-endpoints"><strong>${this.t("sourceLabel")}</strong>${sourceEndpoints.map((name) => this.renderEndpoint(name, "source"))}</div><div class="endpoint-group destination-endpoints"><strong>${this.t("destinationLabel")}</strong>${destinationEndpoints.map((name) => this.renderEndpoint(name, "destination"))}</div></div><courier-command-readout class="route-readout" .heading=${this.t("commandLabel")} .command=${routeCommand} .sessionKey=${`${this.locale}:${selected.source}:${selected.destination}`} .copyLabel=${this.t("copyCommand")} .copiedLabel=${this.t("copiedCommand")} .copyFailedLabel=${this.t("copyFailed")}><div slot="details" class="readout-details"><div class="flag-list">${flags.map((flag) => html`<span class="flag">${flag.syntax}</span>`)}</div></div></courier-command-readout></div></div></section>
        <section id="install" class="stage"><courier-scene .source=${landingInstallSource} .mobileSource=${landingInstallMobileSource}></courier-scene><div class="stage-shell shell"><h2>${this.t("install")}</h2><div class="install-interface"><div class="install-channels" role="list">${installs.map((channel) => html`<button type="button" class="install-channel ${channel.name === install.name ? "selected" : ""}" data-channel=${channel.name} aria-pressed=${String(channel.name === install.name)} @click=${this.chooseInstall}><courier-brand-icon name=${channel.icon}></courier-brand-icon><span>${channel.name}</span></button>`)}</div><courier-command-readout class="install-readout" .heading=${this.t("commandLabel")} .command=${install.command} .sessionKey=${`${this.locale}:${install.name}`} .copyLabel=${this.t("copyCommand")} .copiedLabel=${this.t("copiedCommand")} .copyFailedLabel=${this.t("copyFailed")}><div slot="details" class="readout-details install-identity"><courier-brand-icon name=${install.icon}></courier-brand-icon><strong>${install.name}</strong></div><div slot="footer-actions" class="install-actions"><a href="https://github.com/iwonz/courier/releases/latest" target="_blank" rel="noopener noreferrer"><courier-icon name="package"></courier-icon>${this.t("packages")}</a><a href="https://github.com/iwonz/courier/releases/latest" target="_blank" rel="noopener noreferrer"><courier-icon name="download"></courier-icon>${this.t("direct")}</a></div></courier-command-readout></div></div></section>
        <section id="cli" class="stage"><courier-scene .source=${landingReferenceSource} .mobileSource=${landingReferenceMobileSource}></courier-scene><div class="stage-shell cli-shell shell"><h2>${this.t("cliTitle")}</h2><div class="cli-workspace"><courier-workbench class="reference" .heading=${this.t("commandLabel")}><div class="reference-grid"><div class="reference-column commands"><div class="reference-list">${contractData.commands.map((command) => html`<button type="button" class="command-row" data-command=${command.name} aria-pressed=${String(command.name === this.selectedCommand)} @click=${this.chooseCommand}><code>${command.usage}</code></button>`)}</div></div><div class="reference-column options"><div class="reference-filter"><span>${this.t("optionsLabel")}</span><courier-checkbox .checked=${this.compatibleOnly} ?disabled=${!this.selectedCommand} .label=${this.t("compatibleOnly")} @courier-checkbox-change=${this.setCompatibility}></courier-checkbox></div><div class="reference-list">${visibleFlags.length ? visibleFlags.map((flag) => html`<article class="option-row"><div><code>${flag.syntax}</code><span>${this.t("optionDefault")}: ${flag.default}</span></div><p>${this.t("repeatable")}: ${flag.repeatable ? this.t("yes") : this.t("no")} · ${this.t("applies")}: ${flag.appliesTo.join(", ")}</p></article>`) : html`<p class="empty-state">${this.t("noCompatibleOptions")}</p>`}</div></div></div></courier-workbench><courier-command-readout class="cli-readout" .heading=${this.t("commandLabel")} .command=${selectedCommand?.usage ?? ""} .description=${selectedCommand ? "" : this.t("selectCommand")} .sessionKey=${`${this.locale}:${this.selectedCommand}:${this.compatibleOnly}`} .copyLabel=${this.t("copyCommand")} .copiedLabel=${this.t("copiedCommand")} .copyFailedLabel=${this.t("copyFailed")}><div slot="details" class="readout-details"><div class="flag-list">${visibleFlags.map((flag) => html`<span class="flag">${flag.syntax}</span>`)}</div></div></courier-command-readout></div></div></section>
      </main>`;
  }
}

if (!customElements.get("courier-landing-app")) customElements.define("courier-landing-app", CourierLandingApp);
