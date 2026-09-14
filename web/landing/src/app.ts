import { LitElement, css, html } from "lit";
import { browserLocale, browserThemeState, defineCourierElements, type Locale, type ThemeState } from "@courier/ui";
import { relayDispatchSource, relayInstallSource, relayRoutingSource, relayVerifySource } from "@courier/ui/relay-landing";
import { landingText, type LandingMessage } from "./catalog";
import { contractData } from "./contract";

defineCourierElements();

const installs = [
  ["wget", "wget -qO- https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh"],
  ["PowerShell", "irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex"],
  ["npm", "npm install --global @iwonz/courier"],
  ["npx", "npx @iwonz/courier --help"],
  ["Yarn", "yarn dlx @iwonz/courier --help"],
  ["pnpm", "pnpm dlx @iwonz/courier --help"],
  ["Homebrew", "brew tap iwonz/courier https://github.com/iwonz/courier && brew install --cask iwonz/courier/courier"],
  ["Scoop", "scoop bucket add courier https://github.com/iwonz/courier && scoop install courier/courier"],
] as const;

const primaryInstall = "curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh";

export class CourierLandingApp extends LitElement {
  static properties = { locale: { state: true } };

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
    a { color: inherit; }
    .shell { width: min(78rem, calc(100% - 2rem)); margin: 0 auto; }
    .masthead { position: relative; z-index: 4; display: flex; align-items: center; justify-content: space-between; gap: 1rem; min-height: 5rem; border-bottom: 1px solid var(--courier-color-border); }
    .brand-link { text-decoration: none; }
    nav, .controls, .actions, .links, .section-heading { display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; }
    nav a { color: var(--courier-color-muted); font-size: 0.875rem; font-weight: 700; text-decoration: none; }
    nav a:hover { color: var(--courier-color-text); }
    .hero { position: relative; display: grid; grid-template-columns: minmax(0, 1.1fr) minmax(20rem, 0.9fr); gap: clamp(2rem, 7vw, 6rem); align-items: center; min-height: 43rem; padding: clamp(4rem, 9vw, 8rem) 0; }
    .hero-copy { position: relative; z-index: 2; display: grid; min-width: 0; gap: 1.25rem; }
    .eyebrow, .section-index, .label { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; font-weight: 750; letter-spacing: 0.09em; text-transform: uppercase; }
    h1, h2, h3, p { margin: 0; }
    h1 { max-width: 48rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(3.8rem, 9vw, 7.75rem); font-weight: 850; letter-spacing: -0.075em; line-height: 0.87; }
    .tagline { max-width: 42rem; font-size: clamp(1.35rem, 2.6vw, 2.15rem); font-weight: 650; letter-spacing: -0.035em; line-height: 1.12; }
    .subline, .intro { max-width: 43rem; color: var(--courier-color-muted); font-size: 1.025rem; line-height: 1.65; }
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
    .quick-command { display: grid; grid-template-columns: auto minmax(0, 1fr); gap: 0.75rem; max-width: 48rem; padding: 0.9rem 1rem; border: 1px solid var(--courier-color-border-strong); border-radius: var(--courier-radius-sm); color: var(--courier-paper-50); background: var(--courier-graphite-900); box-shadow: var(--courier-shadow); }
    .quick-command span { color: var(--courier-signal); }
    code, pre { font-family: var(--courier-font-mono); }
    .quick-command code { overflow: auto; white-space: nowrap; }
    .facts { display: grid; grid-template-columns: repeat(3, 1fr); border-top: 1px solid var(--courier-color-border); border-bottom: 1px solid var(--courier-color-border); background: color-mix(in srgb, var(--courier-color-canvas) 88%, transparent); }
    .fact { display: grid; gap: 0.35rem; padding: 1.5rem 0; }
    .fact + .fact { padding-left: 1.5rem; border-left: 1px solid var(--courier-color-border); }
    .fact strong { font-family: var(--courier-font-display); font-size: 2rem; letter-spacing: -0.05em; }
    .fact span { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; letter-spacing: 0.08em; text-transform: uppercase; }
    main { display: grid; gap: clamp(5rem, 10vw, 9rem); padding: clamp(5rem, 10vw, 9rem) 0; }
    section { display: grid; gap: 1.5rem; scroll-margin-top: 2rem; }
    .section-heading { justify-content: space-between; align-items: end; padding-bottom: 1rem; border-bottom: 1px solid var(--courier-color-border); }
    .section-heading > div { display: grid; min-width: 0; gap: 0.6rem; }
    h2 { max-width: 44rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(2.25rem, 5vw, 4.5rem); font-weight: 820; letter-spacing: -0.06em; line-height: 0.95; }
    h3 { font-size: 1.05rem; letter-spacing: -0.025em; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(min(100%, 18rem), 1fr)); gap: 1rem; }
    .install-primary { display: grid; grid-template-columns: 0.82fr 1.18fr; overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); background: var(--courier-graphite-900); box-shadow: var(--courier-shadow); }
    .install-copy { position: relative; z-index: 1; display: grid; align-content: center; gap: 0.75rem; padding: clamp(1.5rem, 4vw, 3rem); color: var(--courier-paper-50); background: var(--courier-graphite-900); }
    .install-copy .label { color: var(--courier-signal); }
    .install-copy pre { margin: 0; padding: 1rem; overflow: auto; border: 1px solid #4d5547; border-radius: var(--courier-radius-sm); background: #10120f; }
    .install-copy code { color: var(--courier-paper-50); white-space: nowrap; }
    .scene-art { min-width: 0; overflow: hidden; background: var(--courier-graphite-900); }
    .scene-art courier-mascot { display: block; width: 100%; height: 100%; }
    .scene-art courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; }
    .install-primary .scene-art { min-height: 24rem; border-left: 1px solid #4d5547; }
    .narrative-brief { display: grid; grid-template-columns: minmax(15rem, 0.62fr) minmax(0, 1.38fr); overflow: hidden; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-lg); color: var(--courier-paper-50); background: var(--courier-graphite-900); box-shadow: var(--courier-shadow); }
    .brief-copy { display: grid; align-content: center; gap: 0.8rem; padding: clamp(1.5rem, 4vw, 3rem); }
    .brief-copy .label { color: var(--courier-signal); }
    .brief-copy p { color: #c7ccbf; line-height: 1.65; }
    .narrative-brief .scene-art { min-height: 24rem; border-left: 1px solid #4d5547; }
    courier-panel, article, .grid > * { min-width: 0; }
    article { display: grid; gap: 0.75rem; }
    article p { color: var(--courier-color-muted); line-height: 1.55; }
    article pre, .example { max-width: 100%; margin: 0; padding: 0.9rem; overflow: auto; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-sm); background: var(--courier-color-surface); }
    article code, .example code { white-space: nowrap; }
    .route-grid article { align-content: start; }
    .route-grid courier-route { margin-top: 0.35rem; }
    .safety-grid courier-icon { width: 2rem; height: 2rem; color: var(--courier-color-accent); }
    .table { overflow-x: auto; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-md); background: var(--courier-color-surface-raised); }
    table { width: 100%; min-width: 42rem; border-collapse: collapse; }
    th, td { padding: 0.85rem 1rem; text-align: left; vertical-align: top; border-bottom: 1px solid var(--courier-color-border); }
    th { color: var(--courier-color-muted); font-family: var(--courier-font-mono); font-size: 0.6875rem; letter-spacing: 0.08em; text-transform: uppercase; }
    tr:last-child td { border-bottom: 0; }
    .final-cta { display: grid; grid-template-columns: 1fr auto; gap: 1.5rem; align-items: center; padding: clamp(1.5rem, 5vw, 3.5rem); border-radius: var(--courier-radius-lg); color: var(--courier-paper-50); background: var(--courier-graphite-900); }
    .final-cta h2 { font-size: clamp(2rem, 5vw, 3.75rem); }
    .final-cta p { margin-top: 0.75rem; color: #aeb5a7; }
    footer { padding: 2rem 0 4rem; border-top: 1px solid var(--courier-color-border); }
    footer .shell { display: flex; justify-content: space-between; gap: 1rem; color: var(--courier-color-muted); font-size: 0.85rem; }
    a:focus-visible { outline: 3px solid var(--courier-beak); outline-offset: 3px; }
    @media (max-width: 58rem) {
      .hero { grid-template-columns: 1fr; padding-top: 4rem; }
      .hero-visual { min-height: 31rem; }
      .install-primary { grid-template-columns: 1fr; }
      .install-primary .scene-art { min-height: 0; aspect-ratio: 3 / 2; border-top: 1px solid #4d5547; border-left: 0; }
      .narrative-brief { grid-template-columns: 1fr; }
      .narrative-brief .scene-art { min-height: 0; aspect-ratio: 3 / 2; border-top: 1px solid #4d5547; border-left: 0; }
    }
    @media (max-width: 44rem) {
      .masthead { display: grid; grid-template-columns: minmax(0, 1fr); align-items: start; padding: 1rem 0; }
      nav { width: 100%; min-width: 0; justify-content: space-between; }
      nav > * { min-width: 0; max-width: 100%; }
      nav > a { display: none; }
      .hero { min-height: auto; padding: 4rem 0; }
      .hero-visual { min-height: 26rem; }
      .facts { grid-template-columns: 1fr; }
      .fact + .fact { padding-left: 0; border-top: 1px solid var(--courier-color-border); border-left: 0; }
      .final-cta { grid-template-columns: 1fr; }
      footer .shell { display: grid; }
    }
  `;

  private locale: Locale = browserLocale();
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

  render() {
    return html`
      <header id="top">
        <div class="shell">
          <div class="masthead">
            <a class="brand-link" href="#top"><courier-brand product=${this.t("brandProduct")}></courier-brand></a>
            <nav aria-label=${this.t("documentation")}>
              <a href="#install">${this.t("installShort")}</a>
              <a href="#routes">${this.t("routes")}</a>
              <a href="#commands">${this.t("commands")}</a>
              <courier-theme-selector .locale=${this.locale}></courier-theme-selector>
              <courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector>
            </nav>
          </div>
          <div class="hero">
            <div class="hero-copy">
              <span class="eyebrow">${this.t("eyebrow")}</span>
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
          <div class="facts">
            <div class="fact"><strong>0</strong><span>${this.t("oneBinary")}</span></div>
            <div class="fact"><strong>4</strong><span>${this.t("routeFamilies")}</span></div>
            <div class="fact"><strong>6</strong><span>${this.t("targetSystems")}</span></div>
          </div>
        </div>
      </header>

      <main class="shell">
        <section id="install">
          <div class="section-heading"><div><span class="section-index">01 / Dispatch</span><h2>${this.t("install")}</h2><p class="intro">${this.t("installIntro")}</p></div></div>
          <div class="install-primary">
            <div class="install-copy"><span class="label">${this.t("primaryChannel")}</span><h3>curl</h3><pre tabindex="0"><code>${primaryInstall}</code></pre></div>
            <div class="scene-art"><courier-mascot alt="" .source=${relayInstallSource}></courier-mascot></div>
          </div>
          <span class="label">${this.t("alternateChannels")}</span>
          <div class="grid">${installs.map(([name, command]) => html`<courier-panel><article><h3>${name}</h3><pre tabindex="0"><code>${command}</code></pre></article></courier-panel>`)}</div>
          <div class="grid">
            <courier-panel heading=${this.t("packages")}><article><p>${this.t("packagesDetail")}</p></article></courier-panel>
            <courier-panel heading=${this.t("direct")}><article><p>${this.t("directDetail")}</p><a href="https://github.com/iwonz/courier/releases/latest">${this.t("downloads")}</a></article></courier-panel>
          </div>
        </section>

        <section id="routes">
          <div class="section-heading"><div><span class="section-index">02 / Routing</span><h2>${this.t("routes")}</h2><p class="intro">${this.t("routeIntro")}</p></div></div>
          <div class="narrative-brief"><div class="brief-copy"><span class="label">${this.t("routeSceneLabel")}</span><p>${this.t("routeSceneDetail")}</p></div><div class="scene-art"><courier-mascot alt="" .source=${relayRoutingSource}></courier-mascot></div></div>
          <div class="grid route-grid">${contractData.routes.map((route) => html`
            <courier-panel><article><span class="label">${route.name}</span><courier-route source=${route.source.join(" / ")} destination=${route.destination.join(" / ")}></courier-route><p><strong>${this.t("allowed")}:</strong> ${route.allowedFlags.map((flag) => `--${flag}`).join(", ")}</p></article></courier-panel>
          `)}</div>
        </section>

        <section id="safety">
          <div class="section-heading"><div><span class="section-index">03 / Handoff</span><h2>${this.t("safety")}</h2><p class="intro">${this.t("safetyIntro")}</p></div></div>
          <div class="narrative-brief"><div class="brief-copy"><span class="label">${this.t("safetySceneLabel")}</span><p>${this.t("safetySceneDetail")}</p></div><div class="scene-art"><courier-mascot alt="" .source=${relayVerifySource}></courier-mascot></div></div>
          <div class="grid safety-grid">
            <courier-panel><article><courier-icon name="route"></courier-icon><h3>${this.t("preflightTitle")}</h3><p>${this.t("preflightDetail")}</p></article></courier-panel>
            <courier-panel><article><courier-icon name="parcel"></courier-icon><h3>${this.t("stagedTitle")}</h3><p>${this.t("stagedDetail")}</p></article></courier-panel>
            <courier-panel><article><courier-icon name="shield"></courier-icon><h3>${this.t("trustTitle")}</h3><p>${this.t("trustDetail")}</p></article></courier-panel>
            <courier-panel><article><courier-icon name="receipt"></courier-icon><h3>${this.t("reportTitle")}</h3><p>${this.t("reportDetail")}</p></article></courier-panel>
          </div>
        </section>

        <section id="commands">
          <div class="section-heading"><div><span class="section-index">04 / Interface</span><h2>${this.t("commands")}</h2></div><a href="https://github.com/iwonz/courier/blob/main/docs/cli-reference.md">${this.t("inspectContract")}</a></div>
          <div class="table"><table><thead><tr><th>${this.t("commandKind")}</th><th>${this.t("commandSyntax")}</th></tr></thead><tbody>
            ${contractData.commands.map((command) => html`<tr><td>${command.system ? this.t("system") : this.t("product")}</td><td><code>${command.usage}</code></td></tr>`)}
          </tbody></table></div>
        </section>

        <section id="options">
          <div class="section-heading"><div><span class="section-index">05 / Policy</span><h2>${this.t("options")}</h2></div></div>
          <div class="table"><table><thead><tr><th>${this.t("option")}</th><th>${this.t("default")}</th><th>${this.t("repeatable")}</th><th>${this.t("applies")}</th></tr></thead><tbody>
            ${contractData.flags.map((flag) => html`<tr><td><code>${flag.syntax}</code></td><td><code>${flag.default}</code></td><td>${flag.repeatable ? this.t("yes") : this.t("no")}</td><td>${flag.appliesTo.join(", ")}</td></tr>`)}
          </tbody></table></div>
        </section>

        <section id="examples">
          <div class="section-heading"><div><span class="section-index">06 / Practice</span><h2>${this.t("examples")}</h2></div></div>
          <div class="grid">${contractData.examples.map((example) => html`<pre class="example" tabindex="0"><code>${example}</code></pre>`)}</div>
        </section>

        <section id="docs">
          <div class="section-heading"><div><span class="section-index">07 / Reference</span><h2>${this.t("documentation")}</h2></div></div>
          <div class="links">
            <a class="button secondary" href="https://github.com/iwonz/courier">${this.t("source")}</a>
            <a class="button secondary" href="https://github.com/iwonz/courier/blob/main/docs/cli-reference.md">${this.t("cliReference")}</a>
            <a class="button secondary" href="https://github.com/iwonz/courier/blob/main/docs/installation.md">${this.t("installGuide")}</a>
            <a class="button secondary" href="https://github.com/iwonz/courier/blob/main/docs/security.md">${this.t("security")}</a>
          </div>
        </section>

        <aside class="final-cta">
          <div><h2>${this.t("finalTitle")}</h2><p>${this.t("finalDetail")}</p></div>
          <a class="button" href="#install">${this.t("install")}</a>
        </aside>
      </main>

      <footer><div class="shell"><span>${this.t("footerLine")}</span><div class="links"><span>${this.t("version")} ${contractData.contractVersion} · ${this.t("target")} ${contractData.targetRelease}</span><a href="https://github.com/iwonz/courier/blob/main/LICENSE">${this.t("license")}</a></div></div></footer>
    `;
  }
}

if (!customElements.get("courier-landing-app")) {
  customElements.define("courier-landing-app", CourierLandingApp);
}
