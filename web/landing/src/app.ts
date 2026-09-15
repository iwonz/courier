import { LitElement, css, html } from "lit";
import {
  browserLocale,
  browserThemeState,
  defineCourierElements,
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
  { name: "curl", command: primaryInstall, icon: "terminal" },
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
    activeInstall: { state: true },
    heroActive: { state: true },
  };

  static styles = css`
    :host {
      --masthead-height: 5rem;
      display: block;
      width: 100%;
      min-height: 100svh;
      overflow-x: clip;
      color: var(--courier-color-text);
      background-color: var(--courier-color-canvas);
      background-image: linear-gradient(var(--courier-color-grid) 1px, transparent 1px), linear-gradient(90deg, var(--courier-color-grid) 1px, transparent 1px);
      background-position: 1.25rem 1.25rem;
      background-size: 2.5rem 2.5rem;
      font-family: var(--courier-font-sans);
    }
    * { box-sizing: border-box; }
    a { color: inherit; }
    button { font: inherit; }
    code, pre { font-family: var(--courier-font-mono); }
    h1, h2, h3, p, pre { margin: 0; }
    .shell { width: min(90rem, calc(100% - clamp(2rem, 6vw, 6rem))); margin: 0 auto; }
    .masthead-wrap {
      position: fixed;
      z-index: 20;
      top: 0;
      right: 0;
      left: 0;
      isolation: isolate;
      pointer-events: none;
    }
    .masthead-wrap::before {
      content: "";
      position: absolute;
      z-index: -1;
      top: 0;
      right: 0;
      left: 0;
      height: calc(100% + 2rem);
      background: linear-gradient(180deg, color-mix(in srgb, var(--courier-color-canvas) 94%, transparent) 0%, color-mix(in srgb, var(--courier-color-canvas) 70%, transparent) 58%, transparent 100%);
      backdrop-filter: blur(18px) saturate(0.86);
      -webkit-backdrop-filter: blur(18px) saturate(0.86);
      mask-image: linear-gradient(180deg, #000 0 58%, transparent 100%);
      pointer-events: none;
    }
    .masthead {
      display: flex;
      min-height: var(--masthead-height);
      align-items: center;
      gap: 1.5rem;
      padding: 0.55rem 0;
      pointer-events: auto;
    }
    .header-left, nav, .header-actions, .preferences, .section-heading, .install-actions { display: flex; align-items: center; }
    .header-left { min-width: 0; gap: clamp(1rem, 3vw, 2.5rem); }
    .brand-link { min-width: 0; text-decoration: none; }
    nav { gap: clamp(0.8rem, 2vw, 1.5rem); }
    nav > a { color: var(--courier-color-muted); font-size: 0.8125rem; font-weight: 780; text-decoration: none; }
    nav > a:hover { color: var(--courier-color-text); }
    .header-actions { margin-left: auto; gap: 0.6rem; }
    .preferences { gap: 0.45rem; }
    .github-link { display: inline-flex; width: 2.45rem; height: 2.45rem; align-items: center; justify-content: center; border: 1px solid transparent; border-radius: var(--courier-radius-md); background: color-mix(in srgb, var(--courier-color-surface) 70%, transparent); }
    .github-link:hover { border-color: var(--courier-color-border-strong); }
    .github-link courier-icon { width: 1.25rem; height: 1.25rem; }
    .slide {
      --scene-x: 0px;
      --scene-y: 0px;
      --scene-spot-x: 72%;
      --scene-spot-y: 42%;
      position: relative;
      min-height: 100svh;
      padding: calc(var(--masthead-height) + clamp(1.25rem, 3vh, 2.5rem)) 0 clamp(1.25rem, 3vh, 2.5rem);
      overflow: hidden;
      isolation: isolate;
      scroll-margin: 0;
      scroll-snap-align: start;
      scroll-snap-stop: always;
    }
    .slide-shell { position: relative; z-index: 2; display: grid; height: calc(100svh - var(--masthead-height) - clamp(2.5rem, 6vh, 5rem)); min-height: 0; gap: clamp(1rem, 2.2vh, 1.75rem); }
    .slide-art { position: absolute; z-index: 0; inset: 0; overflow: hidden; background: var(--courier-graphite-900); pointer-events: none; }
    .slide-art courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; opacity: 0.94; transform: translate3d(var(--scene-x), var(--scene-y), 0) scale(1.035); transition: transform 260ms var(--courier-ease); }
    .slide-art courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; filter: saturate(0.86) contrast(1.02); }
    .slide-art::before { content: ""; position: absolute; z-index: 1; inset: 0; background: radial-gradient(circle at var(--scene-spot-x) var(--scene-spot-y), color-mix(in srgb, var(--courier-color-accent) 12%, transparent), transparent 23%); transition: background-position 240ms var(--courier-ease); }
    .slide-art::after { content: ""; position: absolute; z-index: 2; inset: 0; background: linear-gradient(90deg, rgb(10 12 10 / 0.16), transparent 28% 70%, rgb(10 12 10 / 0.28)), linear-gradient(180deg, color-mix(in srgb, var(--courier-color-canvas) 24%, transparent), transparent 24% 78%, rgb(10 12 10 / 0.28)); }
    .hero-slide { position: relative; display: grid; align-items: stretch; }
    .hero { position: relative; z-index: 3; display: grid; height: 100%; min-height: 0; align-items: center; pointer-events: none; }
    .hero-copy { position: relative; z-index: 3; display: grid; width: min(50%, 44rem); min-width: 0; gap: clamp(0.9rem, 2.3vh, 1.45rem); pointer-events: none; }
    h1 { max-width: 44rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(4rem, 7.3vw, 7.25rem); font-weight: 850; letter-spacing: -0.078em; line-height: 0.84; }
    h2 { max-width: 52rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(2rem, 4.6vw, 4.25rem); font-weight: 820; letter-spacing: -0.06em; line-height: 0.94; }
    h3 { font-size: 1rem; letter-spacing: -0.02em; }
    .tagline { max-width: 42rem; font-size: clamp(1.3rem, 2.4vw, 2.2rem); font-weight: 650; letter-spacing: -0.035em; line-height: 1.12; }
    .subline, .intro { max-width: 50rem; color: var(--courier-color-muted); font-size: clamp(0.9rem, 1.2vw, 1.025rem); line-height: 1.58; }
    .hero-visual {
      --relay-x: 0px;
      --relay-y: 0px;
      --spot-x: 68%;
      --spot-y: 32%;
      position: absolute;
      z-index: 1;
      inset: 0;
      display: block;
      width: 100%;
      height: 100%;
      min-height: 0;
      padding: 0;
      overflow: hidden;
      border: 0;
      color: var(--courier-color-text);
      background: transparent;
      text-align: left;
      cursor: crosshair;
      isolation: isolate;
    }
    .hero-visual::before {
      content: "";
      position: absolute;
      z-index: 2;
      inset: 0;
      background: linear-gradient(90deg, var(--courier-color-canvas) 0 27%, color-mix(in srgb, var(--courier-color-canvas) 78%, transparent) 38%, transparent 59%), radial-gradient(circle at var(--spot-x) var(--spot-y), color-mix(in srgb, var(--courier-color-accent) 18%, transparent) 0 6%, transparent 25%), linear-gradient(180deg, color-mix(in srgb, var(--courier-color-canvas) 32%, transparent), transparent 27% 82%, color-mix(in srgb, var(--courier-color-canvas) 30%, transparent));
      pointer-events: none;
    }
    .hero-visual courier-mascot { position: absolute; z-index: 1; inset: 0; transform: translate3d(var(--relay-x), var(--relay-y), 0) scale(1.025); transition: transform 220ms var(--courier-ease); }
    .hero-visual courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; object-position: center; filter: saturate(0.9) contrast(0.98); }
    .route-trace { position: absolute; z-index: 3; right: 8%; bottom: 17%; left: 7%; height: 0.18rem; border-radius: 999px; background: linear-gradient(90deg, transparent, var(--courier-color-accent) 18% 82%, transparent); box-shadow: 0 0 1.4rem color-mix(in srgb, var(--courier-color-accent) 45%, transparent); transform: rotate(-8deg); transform-origin: center; }
    .route-trace::before, .route-trace::after { content: ""; position: absolute; top: 50%; width: 0.8rem; height: 0.8rem; border: 2px solid var(--courier-color-accent); border-radius: 50%; background: var(--courier-color-canvas); transform: translateY(-50%); }
    .route-trace::before { left: 7%; }
    .route-trace::after { right: 7%; }
    .route-signal { position: absolute; top: 50%; left: 9%; width: 0.6rem; height: 0.6rem; border-radius: 50%; background: var(--courier-color-accent); box-shadow: 0 0 0 0.55rem color-mix(in srgb, var(--courier-color-accent) 18%, transparent); transform: translateY(-50%); }
    .hero-visual[data-active="true"] .route-signal { animation: dispatch-signal 850ms var(--courier-ease) both; }
    .hero-node { position: absolute; z-index: 4; bottom: 6%; display: inline-flex; align-items: center; gap: 0.45rem; padding: 0.45rem 0.6rem; border-radius: var(--courier-radius-sm); color: var(--courier-color-muted); background: color-mix(in srgb, var(--courier-color-canvas) 78%, transparent); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 750; letter-spacing: 0.07em; text-transform: uppercase; backdrop-filter: blur(8px); }
    .hero-node courier-icon { width: 1rem; height: 1rem; color: var(--courier-color-accent-ink); }
    .hero-node.source { left: 4%; }
    .hero-node.destination { right: 3%; }
    .hero-hint { position: absolute; z-index: 4; top: 8%; right: 2%; color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; letter-spacing: 0.06em; }
    @keyframes dispatch-signal { from { left: 9%; } to { left: calc(91% - 0.6rem); } }
    .hero-visual:focus-visible, .endpoint:focus-visible, .install-channel:focus-visible, .github-link:focus-visible, a:focus-visible { outline: 3px solid var(--courier-beak); outline-offset: 3px; }
    .section-heading { min-width: 0; justify-content: space-between; align-items: end; color: var(--courier-paper-50); text-shadow: 0 0.14rem 1rem rgb(10 12 10 / 0.72); }
    .section-heading > div { display: grid; min-width: 0; gap: 0.45rem; }
    .section-index, .label, .kind { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6625rem; font-weight: 750; letter-spacing: 0.09em; text-transform: uppercase; }
    .section-heading .section-index { color: var(--courier-signal); }
    .section-heading .intro { color: #d4d8ce; }
    .route-explorer, .install-board, .reference { min-height: 0; overflow: hidden; border: 1px solid #545c4e; border-radius: var(--courier-radius-lg); box-shadow: 0 1.4rem 3.2rem rgb(7 9 7 / 0.34), inset 0 1px rgb(255 255 255 / 0.04); backdrop-filter: blur(18px) saturate(0.82); }
    .route-explorer { display: grid; width: min(56rem, 72%); height: 100%; justify-self: end; border-color: #4d5547; color: var(--courier-paper-50); background: rgb(21 23 20 / 0.94); }
    .route-interface { display: grid; min-width: 0; min-height: 0; align-content: center; gap: clamp(0.75rem, 2vh, 1.25rem); padding: clamp(1rem, 3vw, 2.5rem); }
    .endpoint-columns { display: grid; min-height: 0; grid-template-columns: minmax(0, 1fr) minmax(2.25rem, auto) minmax(0, 1fr); align-items: start; gap: clamp(0.5rem, 1.4vw, 1rem); }
    .endpoint-group { display: grid; min-width: 0; gap: clamp(0.35rem, 0.8vh, 0.55rem); }
    .endpoint-group .label { color: #aeb5a7; }
    .endpoint-word { align-self: center; padding-top: 1.8rem; color: var(--courier-signal); font-family: var(--courier-font-mono); font-weight: 850; text-align: center; }
    .endpoint { appearance: none; display: grid; grid-template-columns: auto minmax(0, 1fr); min-width: 0; min-height: 2.65rem; align-items: center; gap: 0.6rem; padding: 0.5rem 0.65rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-sm); color: #c7ccbf; background: #1c1f1b; text-align: left; cursor: pointer; transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), opacity var(--courier-duration) var(--courier-ease); }
    .endpoint courier-icon { width: 1.1rem; height: 1.1rem; }
    .endpoint.valid { border-color: #6d7567; }
    .endpoint.selected { border-color: var(--courier-signal); color: var(--courier-signal); background: #252a21; }
    .endpoint.invalid { opacity: 0.34; cursor: not-allowed; }
    .route-readout { display: grid; min-width: 0; gap: 0.65rem; padding: 0.85rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-md); background: #10120f; }
    .route-readout code { display: block; max-width: 100%; overflow-wrap: anywhere; color: var(--courier-paper-50); font-size: clamp(0.72rem, 1.35vw, 0.9rem); white-space: normal; }
    .command-shape { color: #aeb5a7; }
    .command-shape b { color: var(--courier-signal); }
    .flag-list { display: flex; max-height: 4.5rem; gap: 0.35rem; overflow-y: auto; flex-wrap: wrap; }
    .flag { padding: 0.28rem 0.42rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-sm); color: #c7ccbf; background: #1c1f1b; font-family: var(--courier-font-mono); font-size: 0.67rem; }
    .install-board { position: relative; display: grid; width: min(52rem, 68%); height: 100%; justify-self: end; color: var(--courier-paper-50); background: rgb(18 21 18 / 0.9); }
    .install-interface { display: grid; min-width: 0; min-height: 0; align-content: center; gap: clamp(0.75rem, 1.8vh, 1.1rem); padding: clamp(1rem, 2.6vw, 2rem); }
    .install-channels { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 0.5rem; }
    .install-channel { appearance: none; display: grid; grid-template-columns: auto minmax(0, 1fr); min-width: 0; min-height: 2.65rem; align-items: center; gap: 0.5rem; padding: 0.5rem 0.6rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-sm); color: #c7ccbf; background: rgb(28 31 27 / 0.88); text-align: left; cursor: pointer; }
    .install-channel courier-icon { width: 1.1rem; height: 1.1rem; color: var(--courier-signal); }
    .install-channel.selected { border-color: var(--courier-signal); color: var(--courier-signal); background: rgb(37 42 33 / 0.96); }
    .install-channel span { min-width: 0; overflow: hidden; font-size: 0.75rem; font-weight: 780; text-overflow: ellipsis; white-space: nowrap; }
    .install-readout { display: grid; min-width: 0; gap: 0.65rem; padding: 0.85rem; border: 1px solid #545c4e; border-radius: var(--courier-radius-md); color: var(--courier-paper-50); background: rgb(10 12 10 / 0.9); }
    .install-readout-head { display: flex; align-items: center; gap: 0.55rem; color: var(--courier-signal); }
    .install-readout-head courier-icon { width: 1.1rem; height: 1.1rem; }
    .install-readout pre { width: 100%; min-width: 0; max-width: 100%; overflow-x: auto; }
    .install-readout code { display: block; min-width: 0; overflow-wrap: anywhere; font-size: clamp(0.7rem, 1.1vw, 0.82rem); line-height: 1.5; white-space: pre-wrap; }
    .install-actions { align-items: stretch; gap: 0.65rem; }
    .download-channel { display: grid; min-width: 0; flex: 1; grid-template-columns: auto minmax(0, 1fr); align-items: center; gap: 0.65rem; padding: 0.75rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-md); background: rgb(28 31 27 / 0.88); text-decoration: none; }
    .download-channel:hover { border-color: var(--courier-signal); }
    .download-channel courier-icon { width: 1.35rem; height: 1.35rem; color: var(--courier-signal); }
    .download-channel div { display: grid; min-width: 0; gap: 0.2rem; }
    .download-channel p { overflow: hidden; color: #aeb5a7; font-size: 0.72rem; line-height: 1.3; text-overflow: ellipsis; white-space: nowrap; }
    .reference { display: grid; width: min(76rem, 88%); height: 100%; grid-template-columns: minmax(15rem, 0.68fr) minmax(0, 1.32fr); justify-self: end; color: var(--courier-paper-50); background: rgb(18 21 18 / 0.91); }
    .reference-column { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(0, 1fr); }
    .reference-column + .reference-column { border-left: 1px solid #4d5547; }
    .reference-label { display: block; padding: 0.75rem 1rem; border-bottom: 1px solid #4d5547; color: #aeb5a7; }
    .reference-list { min-width: 0; min-height: 0; overflow-y: auto; overscroll-behavior: contain; }
    .command-row, .option-row { display: grid; min-width: 0; gap: 0.4rem; padding: 0.72rem 1rem; border-bottom: 1px solid #40473b; transition: color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease); }
    .command-row:hover, .option-row:hover { color: var(--courier-signal); background: rgb(40 45 35 / 0.78); }
    .command-row:last-child, .option-row:last-child { border-bottom: 0; }
    .command-row code, .option-row code { overflow-wrap: anywhere; font-size: 0.76rem; }
    .option-head { display: flex; align-items: baseline; justify-content: space-between; gap: 0.75rem; }
    .option-meta { color: #aeb5a7; font-family: var(--courier-font-mono); font-size: 0.65rem; line-height: 1.45; }
    @media (max-width: 62rem) {
      .route-explorer, .install-board { width: 78%; }
      .reference { grid-template-columns: minmax(12rem, 0.58fr) minmax(0, 1.42fr); }
    }
    @media (max-width: 44rem) {
      :host { --masthead-height: 6.75rem; }
      .shell { width: min(100% - 1.25rem, 90rem); }
      .masthead { display: grid; grid-template-columns: minmax(0, 1fr) auto; align-items: center; gap: 0.35rem 0.65rem; padding: 0.45rem 0; }
      .header-left { display: contents; }
      .brand-link { grid-column: 1; grid-row: 1; }
      nav { grid-column: 1 / -1; grid-row: 2; justify-content: flex-start; gap: 1.1rem; }
      .header-actions { grid-column: 2; grid-row: 1; gap: 0.3rem; }
      .preferences { gap: 0.25rem; }
      .hero { align-items: start; }
      .hero-copy { width: 92%; max-width: 92%; padding-top: 1.2rem; gap: 0.65rem; }
      h1 { max-width: 22rem; overflow-wrap: normal; font-size: clamp(2.75rem, 13.5vw, 3.75rem); hyphens: none; word-break: normal; }
      .tagline { font-size: 1.08rem; }
      .subline { display: none; }
      .hero-visual::before { background: linear-gradient(180deg, color-mix(in srgb, var(--courier-color-canvas) 92%, transparent) 0 27%, color-mix(in srgb, var(--courier-color-canvas) 58%, transparent) 39%, transparent 58% 88%, color-mix(in srgb, var(--courier-color-canvas) 24%, transparent)), radial-gradient(circle at var(--spot-x) var(--spot-y), color-mix(in srgb, var(--courier-color-accent) 17%, transparent), transparent 24%); }
      .hero-visual courier-mascot { inset: 0; }
      .hero-visual courier-mascot::part(image) { object-position: center; }
      .hero-hint { top: 2%; }
      .section-heading .intro { display: none; }
      .slide-shell { gap: 0.75rem; }
      .slide-art courier-mascot { opacity: 0.9; transform: translate3d(var(--scene-x), var(--scene-y), 0) scale(1.045); }
      .slide-art::after { background: linear-gradient(180deg, rgb(10 12 10 / 0.3), transparent 22% 78%, rgb(10 12 10 / 0.38)), linear-gradient(90deg, rgb(10 12 10 / 0.12), transparent 22% 78%, rgb(10 12 10 / 0.16)); }
      .route-explorer, .install-board, .reference { width: 100%; }
      .route-interface, .install-interface { position: relative; z-index: 1; align-content: center; padding: 0.75rem; }
      .endpoint-columns { grid-template-columns: minmax(0, 1fr) 1.6rem minmax(0, 1fr); gap: 0.35rem; }
      .endpoint { grid-template-columns: 1fr; min-height: 2.35rem; justify-items: center; padding: 0.35rem; text-align: center; }
      .endpoint-word { padding-top: 1.6rem; }
      .route-readout { padding: 0.65rem; }
      .flag-list { max-height: 3.5rem; }
      .install-channels { gap: 0.35rem; }
      .install-channel { grid-template-columns: 1fr; min-height: 2.35rem; justify-items: center; gap: 0.2rem; padding: 0.3rem; text-align: center; }
      .install-channel span { max-width: 100%; font-size: 0.65rem; }
      .install-readout { padding: 0.65rem; }
      .install-actions { gap: 0.4rem; }
      .download-channel { padding: 0.55rem; }
      .download-channel p { display: none; }
      .reference { grid-template-columns: minmax(0, 0.84fr) minmax(0, 1.16fr); }
      .reference-label, .command-row, .option-row { padding-right: 0.55rem; padding-left: 0.55rem; }
      .option-head { display: grid; gap: 0.25rem; }
      .option-meta { font-size: 0.58rem; }
      courier-brand::part(product) { display: none; }
    }
    @media (max-width: 30rem) {
      courier-brand::part(words) { display: none; }
      nav { gap: 0.85rem; }
      nav > a { font-size: 0.72rem; }
      .github-link { width: 2.25rem; height: 2.25rem; }
    }
    @media (max-height: 44rem) and (min-width: 44.01rem) {
      :host { --masthead-height: 4.25rem; }
      .hero-visual { min-height: 20rem; }
      .tagline { font-size: 1.2rem; }
      .subline { font-size: 0.84rem; }
      .section-heading .intro { display: none; }
      .endpoint { min-height: 2.2rem; padding-block: 0.32rem; }
      .install-channel { min-height: 2.15rem; padding-block: 0.3rem; }
    }
    @media (prefers-reduced-motion: reduce) {
      .hero-visual courier-mascot, .slide-art courier-mascot { transition: none; transform: none; }
      .hero-visual[data-active="true"] .route-signal { animation: none; left: calc(91% - 0.6rem); }
    }
  `;

  private locale: Locale = browserLocale();
  private selectedSource = initialRoute.source;
  private selectedDestination = initialRoute.destination;
  private activeInstall = installs[0]!.name;
  private heroActive = false;
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

  private chooseInstall(event: Event): void {
    this.activeInstall = (event.currentTarget as HTMLElement).dataset.channel!;
  }

  private navigate(event: MouseEvent): void {
    event.preventDefault();
    const hash = (event.currentTarget as HTMLAnchorElement).getAttribute("href")!;
    (this.renderRoot.querySelector(hash) as HTMLElement).scrollIntoView({ block: "start" });
    globalThis.history.pushState(null, "", hash);
  }

  private moveHero(event: PointerEvent): void {
    const target = event.currentTarget as HTMLElement;
    const bounds = target.getBoundingClientRect();
    const x = Math.min(100, Math.max(0, ((event.clientX - bounds.left) / bounds.width) * 100));
    const y = Math.min(100, Math.max(0, ((event.clientY - bounds.top) / bounds.height) * 100));
    target.style.setProperty("--spot-x", `${x}%`);
    target.style.setProperty("--spot-y", `${y}%`);
    target.style.setProperty("--relay-x", `${((x - 50) / 50) * 12}px`);
    target.style.setProperty("--relay-y", `${((y - 50) / 50) * 8}px`);
  }

  private resetHero(event: PointerEvent): void {
    const target = event.currentTarget as HTMLElement;
    target.style.setProperty("--spot-x", "68%");
    target.style.setProperty("--spot-y", "32%");
    target.style.setProperty("--relay-x", "0px");
    target.style.setProperty("--relay-y", "0px");
  }

  private moveScene(event: PointerEvent): void {
    const target = event.currentTarget as HTMLElement;
    const bounds = target.getBoundingClientRect();
    const x = Math.min(100, Math.max(0, ((event.clientX - bounds.left) / bounds.width) * 100));
    const y = Math.min(100, Math.max(0, ((event.clientY - bounds.top) / bounds.height) * 100));
    target.style.setProperty("--scene-spot-x", `${x}%`);
    target.style.setProperty("--scene-spot-y", `${y}%`);
    target.style.setProperty("--scene-x", `${((x - 50) / 50) * -8}px`);
    target.style.setProperty("--scene-y", `${((y - 50) / 50) * -6}px`);
  }

  private resetScene(event: PointerEvent): void {
    const target = event.currentTarget as HTMLElement;
    target.style.setProperty("--scene-spot-x", "72%");
    target.style.setProperty("--scene-spot-y", "42%");
    target.style.setProperty("--scene-x", "0px");
    target.style.setProperty("--scene-y", "0px");
  }

  private activateHero(): void {
    this.heroActive = !this.heroActive;
  }

  private operateHero(event: KeyboardEvent): void {
    if (event.key !== "Enter" && event.key !== " ") return;
    event.preventDefault();
    this.activateHero();
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
    const install = installs.find((channel) => channel.name === this.activeInstall) ?? installs[0]!;
    return html`
      <header class="masthead-wrap">
        <div class="masthead shell">
          <div class="header-left">
            <a class="brand-link" href="#hero" @click=${this.navigate}><courier-brand product=${this.t("brandProduct")}></courier-brand></a>
            <nav aria-label="Courier">
              <a href="#routes" @click=${this.navigate}>${this.t("routeShort")}</a>
              <a href="#install" @click=${this.navigate}>${this.t("installShort")}</a>
              <a href="#cli" @click=${this.navigate}>${this.t("cliShort")}</a>
            </nav>
          </div>
          <div class="header-actions">
            <a class="github-link" href="https://github.com/iwonz/courier" aria-label=${this.t("githubLabel")}><courier-icon name="github"></courier-icon></a>
            <div class="preferences">
              <courier-theme-selector .locale=${this.locale}></courier-theme-selector>
              <courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector>
            </div>
          </div>
        </div>
      </header>

      <main>
        <section id="hero" class="slide hero-slide">
          <button
            type="button"
            class="hero-visual"
            data-active=${String(this.heroActive)}
            aria-pressed=${String(this.heroActive)}
            aria-label=${this.t("heroInteractionLabel")}
            @pointermove=${this.moveHero}
            @pointerleave=${this.resetHero}
            @click=${this.activateHero}
            @keydown=${this.operateHero}
          >
            <courier-mascot eager alt="" .source=${relayHeroSource} .mobileSource=${relayHeroMobileSource}></courier-mascot>
            <span class="route-trace" aria-hidden="true"><span class="route-signal"></span></span>
            <span class="hero-node source"><courier-icon name="folder"></courier-icon>${this.t("endpointLocal")}</span>
            <span class="hero-node destination"><courier-icon name="server"></courier-icon>${this.t("endpointSsh")}</span>
            <span class="hero-hint" aria-hidden="true">${this.t("heroInteractionHint")}</span>
          </button>
          <div class="hero shell">
            <div class="hero-copy">
              <h1>${this.t("title")}</h1>
              <p class="tagline">${this.t("tagline")}</p>
              <p class="subline">${this.t("subline")}</p>
            </div>
          </div>
        </section>

        <section id="routes" class="slide" @pointermove=${this.moveScene} @pointerleave=${this.resetScene}>
          <div class="slide-art route-art"><courier-mascot alt="" .source=${relayRoutingSource} .mobileSource=${relayRoutingMobileSource}></courier-mascot></div>
          <div class="slide-shell shell">
            <div class="section-heading"><div><span class="section-index">01 / Routing</span><h2>${this.t("routeTitle")}</h2><p class="intro">${this.t("routeIntro")}</p></div></div>
            <div class="route-explorer">
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
          </div>
        </section>

        <section id="install" class="slide" @pointermove=${this.moveScene} @pointerleave=${this.resetScene}>
          <div class="slide-art install-art"><courier-mascot alt="" .source=${relayInstallSource} .mobileSource=${relayInstallMobileSource}></courier-mascot></div>
          <div class="slide-shell shell">
            <div class="section-heading"><div><span class="section-index">02 / Distribution</span><h2>${this.t("install")}</h2><p class="intro">${this.t("installIntro")}</p></div></div>
            <div class="install-board">
              <div class="install-interface">
                <span class="label">${this.t("chooseChannel")}</span>
                <div class="install-channels" role="list">
                  ${installs.map((channel) => html`<button
                    type="button"
                    class="install-channel ${channel.name === install.name ? "selected" : ""}"
                    data-channel=${channel.name}
                    aria-pressed=${String(channel.name === install.name)}
                    @pointerenter=${this.chooseInstall}
                    @focus=${this.chooseInstall}
                    @click=${this.chooseInstall}
                  ><courier-icon name=${channel.icon}></courier-icon><span>${channel.name}</span></button>`)}
                </div>
                <div class="install-readout" aria-live="polite">
                  <div class="install-readout-head"><courier-icon name=${install.icon}></courier-icon><h3>${install.name}</h3></div>
                  <pre tabindex="0"><code>${install.command}</code></pre>
                </div>
                <div class="install-actions">
                  <a class="download-channel" href="https://github.com/iwonz/courier/releases/latest"><courier-icon name="linux"></courier-icon><div><h3>${this.t("packages")}</h3><p>${this.t("packagesDetail")}</p></div></a>
                  <a class="download-channel" href="https://github.com/iwonz/courier/releases/latest"><courier-icon name="download"></courier-icon><div><h3>${this.t("direct")}</h3><p>${this.t("directDetail")}</p></div></a>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section id="cli" class="slide" @pointermove=${this.moveScene} @pointerleave=${this.resetScene}>
          <div class="slide-art cli-art"><courier-mascot alt="" .source=${relayCliSource} .mobileSource=${relayCliMobileSource}></courier-mascot></div>
          <div class="slide-shell shell">
            <div class="section-heading"><div><span class="section-index">03 / Interface</span><h2>${this.t("cliTitle")}</h2><p class="intro">${this.t("cliIntro")}</p></div></div>
            <div class="reference">
              <div class="reference-column">
                <span class="reference-label label">${this.t("commandKind")}</span>
                <div class="reference-list">${contractData.commands.map((command) => html`<article class="command-row"><span class="kind">${command.system ? this.t("system") : this.t("product")}</span><code>${command.usage}</code></article>`)}</div>
              </div>
              <div class="reference-column">
                <span class="reference-label label">${this.t("allowed")}</span>
                <div class="reference-list">${contractData.flags.map((flag) => html`<article class="option-row"><div class="option-head"><code>${flag.syntax}</code><span class="kind">${this.t("optionDefault")}: ${flag.default}</span></div><p class="option-meta">${this.t("repeatable")}: ${flag.repeatable ? this.t("yes") : this.t("no")} · ${this.t("applies")}: ${flag.appliesTo.join(", ")}</p></article>`)}</div>
              </div>
            </div>
          </div>
        </section>
      </main>
    `;
  }
}

if (!customElements.get("courier-landing-app")) {
  customElements.define("courier-landing-app", CourierLandingApp);
}
