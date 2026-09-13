import { LitElement, css, html } from "lit";
import { browserLocale, browserThemeState, defineCourierElements, type Locale, type ThemeState } from "@courier/ui";
import { landingText, type LandingMessage } from "./catalog";
import { contractData } from "./contract";

defineCourierElements();

const installs = [
  ["curl", "curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh"],
  ["wget", "wget -qO- https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh"],
  ["PowerShell", "irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex"],
  ["npm", "npm install --global @iwonz/courier"],
  ["npx", "npx @iwonz/courier --help"],
  ["Yarn", "yarn dlx @iwonz/courier --help"],
  ["pnpm", "pnpm dlx @iwonz/courier --help"],
  ["Homebrew", "brew tap iwonz/courier https://github.com/iwonz/courier && brew install --cask iwonz/courier/courier"],
  ["Scoop", "scoop bucket add courier https://github.com/iwonz/courier && scoop install courier/courier"],
] as const;

export class CourierLandingApp extends LitElement {
  static properties = { locale: { state: true } };

  static styles = css`
    :host { display: block; min-height: 100vh; background: var(--courier-color-canvas); color: var(--courier-color-text); font-family: var(--courier-font-sans); }
    header, main, footer { width: min(76rem, calc(100% - 2rem)); margin: 0 auto; }
    header { min-height: 72vh; display: grid; align-content: center; gap: 1.25rem; padding: 4rem 0; }
    nav, .controls, .links { display: flex; align-items: center; gap: .75rem; flex-wrap: wrap; }
    nav { justify-content: space-between; position: absolute; top: 1rem; width: min(76rem, calc(100% - 2rem)); }
    nav a, footer a { color: inherit; }
    h1 { margin: 0; font-size: clamp(3rem, 12vw, 8rem); letter-spacing: -.07em; line-height: .9; }
    .tagline { max-width: 50rem; margin: 0; font-size: clamp(1.35rem, 3vw, 2.4rem); line-height: 1.15; }
    .subline, .meta { color: var(--courier-color-muted); }
    .hero-links { display: flex; gap: .75rem; flex-wrap: wrap; }
    .button { min-height: 2.75rem; display: inline-flex; align-items: center; padding: 0 1rem; color: var(--courier-color-canvas); background: var(--courier-color-accent); border: 1px solid var(--courier-color-accent); text-decoration: none; border-radius: var(--courier-radius-small); font-weight: 700; }
    .button.secondary { color: var(--courier-color-text); background: transparent; border-color: var(--courier-color-border); }
    main { display: grid; gap: 5rem; padding: 2rem 0 6rem; }
    section { scroll-margin-top: 1rem; display: grid; gap: 1rem; }
    h2 { margin: 0; font-size: clamp(2rem, 5vw, 3.5rem); }
    h3 { margin: 0; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(min(100%, 19rem), 1fr)); gap: 1rem; }
    article, courier-panel, .grid > * { min-width: 0; }
    pre { box-sizing: border-box; max-width: 100%; margin: .75rem 0 0; padding: .9rem; overflow: auto; background: var(--courier-color-surface); border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-small); }
    code { font-family: var(--courier-font-mono); }
    .table { overflow-x: auto; border: 1px solid var(--courier-color-border); border-radius: var(--courier-radius-medium); }
    table { width: 100%; border-collapse: collapse; min-width: 42rem; }
    th, td { padding: .8rem; text-align: left; vertical-align: top; border-bottom: 1px solid var(--courier-color-border); }
    th { color: var(--courier-color-muted); }
    tr:last-child td { border-bottom: 0; }
    footer { padding: 2rem 0 4rem; border-top: 1px solid var(--courier-color-border); }
    a:focus-visible { outline: 3px solid var(--courier-color-accent); outline-offset: 3px; }
    @media (max-width: 44rem) { header { min-height: 82vh; } nav { align-items: flex-start; } .controls { justify-content: flex-end; } main { gap: 3.5rem; } }
    @media (prefers-reduced-motion: reduce) { * { scroll-behavior: auto !important; } }
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
      <header>
        <nav aria-label=${this.t("documentation")}>
          <a href="#top"><strong>COURIER</strong></a>
          <div class="controls">
            <a href="#install">${this.t("install")}</a>
            <a href="#commands">${this.t("commands")}</a>
            <a href="#docs">${this.t("documentation")}</a>
            <courier-theme-selector></courier-theme-selector>
            <courier-locale-selector @courier-locale-change=${this.setLocale}></courier-locale-selector>
          </div>
        </nav>
        <div id="top" class="meta">${this.t("version")} ${contractData.contractVersion} · ${this.t("target")} ${contractData.targetRelease}</div>
        <h1>${this.t("title")}</h1>
        <p class="tagline">${this.t("tagline")}</p>
        <p class="subline">${this.t("subline")}</p>
        <div class="hero-links">
          <a class="button" href="#install">${this.t("install")}</a>
          <a class="button secondary" href="https://github.com/iwonz/courier/releases">${this.t("releases")}</a>
        </div>
      </header>
      <main>
        <section id="install">
          <h2>${this.t("install")}</h2>
          <p>${this.t("installIntro")}</p>
          <div class="grid">${installs.map(([name, command]) => html`<courier-panel><article><h3>${name}</h3><pre tabindex="0"><code>${command}</code></pre></article></courier-panel>`)}</div>
          <div class="grid">
            <courier-panel><article><h3>${this.t("packages")}</h3><p>${this.t("packagesDetail")}</p></article></courier-panel>
            <courier-panel><article><h3>${this.t("direct")}</h3><p>${this.t("directDetail")}</p><a href="https://github.com/iwonz/courier/releases/latest">${this.t("downloads")}</a></article></courier-panel>
          </div>
        </section>

        <section id="commands">
          <h2>${this.t("commands")}</h2>
          <div class="table"><table><thead><tr><th>${this.t("commandKind")}</th><th>${this.t("commandSyntax")}</th></tr></thead><tbody>
            ${contractData.commands.map((command) => html`<tr><td>${command.system ? this.t("system") : this.t("product")}</td><td><code>${command.usage}</code></td></tr>`)}
          </tbody></table></div>
        </section>

        <section id="routes">
          <h2>${this.t("routes")}</h2>
          <div class="table"><table><thead><tr><th>${this.t("route")}</th><th>${this.t("from")}</th><th>${this.t("to")}</th><th>${this.t("allowed")}</th></tr></thead><tbody>
            ${contractData.routes.map((route) => html`<tr><td><code>${route.name}</code></td><td>${route.source.join(", ")}</td><td>${route.destination.join(", ")}</td><td>${route.allowedFlags.map((flag) => `--${flag}`).join(", ")}</td></tr>`)}
          </tbody></table></div>
        </section>

        <section id="options">
          <h2>${this.t("options")}</h2>
          <div class="table"><table><thead><tr><th>${this.t("option")}</th><th>${this.t("default")}</th><th>${this.t("repeatable")}</th><th>${this.t("applies")}</th></tr></thead><tbody>
            ${contractData.flags.map((flag) => html`<tr><td><code>${flag.syntax}</code></td><td><code>${flag.default}</code></td><td>${flag.repeatable ? this.t("yes") : this.t("no")}</td><td>${flag.appliesTo.join(", ")}</td></tr>`)}
          </tbody></table></div>
        </section>

        <section id="examples">
          <h2>${this.t("examples")}</h2>
          <div class="grid">${contractData.examples.map((example) => html`<pre tabindex="0"><code>${example}</code></pre>`)}</div>
        </section>

        <section id="docs">
          <h2>${this.t("documentation")}</h2>
          <div class="links">
            <a class="button secondary" href="https://github.com/iwonz/courier">${this.t("source")}</a>
            <a class="button secondary" href="https://github.com/iwonz/courier/blob/main/docs/cli-reference.md">${this.t("cliReference")}</a>
            <a class="button secondary" href="https://github.com/iwonz/courier/blob/main/docs/installation.md">${this.t("installGuide")}</a>
            <a class="button secondary" href="https://github.com/iwonz/courier/blob/main/docs/security.md">${this.t("security")}</a>
          </div>
        </section>
      </main>
      <footer><div class="links"><span>Courier · ${contractData.targetRelease}</span><a href="https://github.com/iwonz/courier/blob/main/LICENSE">${this.t("license")}</a></div></footer>
    `;
  }
}

if (!customElements.get("courier-landing-app")) {
  customElements.define("courier-landing-app", CourierLandingApp);
}
