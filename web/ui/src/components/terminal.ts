import { LitElement, css, html, nothing, type PropertyValues } from "lit";

export class CourierTerminal extends LitElement {
  static properties = {
    heading: { type: String },
    status: { type: String },
  };

  static styles = css`
    :host {
      display: grid;
      min-width: 0;
      min-height: 0;
      grid-template-rows: auto minmax(0, 1fr) auto;
      overflow: hidden;
      border: 1px solid var(--courier-terminal-border, rgb(137 147 129 / 0.48));
      border-radius: var(--courier-radius-md, 0.625rem);
      color: var(--courier-terminal-text, var(--courier-color-text, #f3f4e9));
      background: var(--courier-terminal-surface, rgb(12 15 12 / 0.68));
      box-shadow: inset 0 1px rgb(255 255 255 / 0.035), 0 1rem 3rem rgb(0 0 0 / 0.12);
      font-family: var(--courier-font-mono, monospace);
      backdrop-filter: blur(18px) saturate(0.8);
    }
    header {
      display: flex;
      min-height: 2.7rem;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
      padding: 0.55rem 0.75rem;
      border-bottom: 1px solid var(--courier-terminal-border, rgb(137 147 129 / 0.48));
      color: var(--courier-terminal-muted, #b9c0b1);
      background: linear-gradient(90deg, rgb(255 255 255 / 0.035), transparent 62%);
      font-size: 0.66rem;
      font-weight: 760;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .status { color: var(--courier-terminal-prompt, #d4ff45); }
    .body { min-width: 0; min-height: 0; overflow: auto; }
    footer {
      border-top: 1px solid var(--courier-terminal-border, rgb(137 147 129 / 0.48));
    }
    ::slotted([slot="toolbar"]), ::slotted([slot="footer"]) { min-width: 0; }
  `;

  heading = "";
  status = "";

  protected render() {
    return html`
      <header>
        <span>${this.heading}</span>
        <slot name="toolbar"></slot>
        ${this.status ? html`<span class="status">${this.status}</span>` : nothing}
      </header>
      <div class="body"><slot></slot></div>
      <footer><slot name="footer"></slot></footer>
    `;
  }
}

export class CourierCommandReadout extends LitElement {
  static properties = {
    command: { type: String },
    description: { type: String },
    heading: { type: String },
    copyLabel: { type: String, attribute: "copy-label" },
    copiedLabel: { type: String, attribute: "copied-label" },
    copyFailedLabel: { type: String, attribute: "copy-failed-label" },
    sessionKey: { type: String, attribute: "session-key" },
    copyState: { state: true },
  };

  static styles = css`
    :host { display: block; min-width: 0; min-height: 0; }
    courier-terminal { height: 100%; }
    .toolbar { display: flex; align-items: center; }
    button {
      appearance: none;
      display: inline-flex;
      min-height: 1.9rem;
      align-items: center;
      gap: 0.38rem;
      padding: 0.28rem 0.55rem;
      border: 1px solid var(--courier-terminal-border, #596253);
      border-radius: 999px;
      color: var(--courier-terminal-text, #f3f4e9);
      background: rgb(255 255 255 / 0.035);
      font: 700 0.65rem/1 var(--courier-font-sans, sans-serif);
      cursor: pointer;
    }
    button:hover:not(:disabled), button:focus-visible { border-color: var(--courier-terminal-prompt, #d4ff45); color: var(--courier-terminal-prompt, #d4ff45); }
    button:focus-visible { outline: 3px solid var(--courier-beak, #ff8758); outline-offset: 2px; }
    button:disabled { cursor: not-allowed; opacity: 0.42; }
    courier-icon { width: 0.9rem; height: 0.9rem; }
    .readout { display: grid; height: 100%; min-height: 0; grid-template-rows: auto auto minmax(0, 1fr); }
    .prompt, .description { min-width: 0; padding: 0.58rem 0.8rem; }
    .prompt { display: grid; grid-template-columns: auto minmax(0, 1fr); gap: 0.65rem; color: var(--courier-terminal-text, #f3f4e9); }
    .prompt::before { content: "$"; color: var(--courier-terminal-prompt, #d4ff45); font-weight: 800; }
    code { overflow-wrap: anywhere; font: inherit; line-height: 1.45; white-space: pre-wrap; }
    .description { border-top: 1px solid color-mix(in srgb, var(--courier-terminal-border, #596253) 58%, transparent); color: var(--courier-terminal-muted, #b9c0b1); font-size: 0.68rem; line-height: 1.45; }
    .details { min-width: 0; min-height: 0; overflow: auto; }
    .footer { display: flex; min-width: 0; align-items: center; justify-content: space-between; gap: 0.75rem; padding: 0.48rem 0.8rem; color: var(--courier-terminal-muted, #b9c0b1); font-size: 0.58rem; line-height: 1.4; }
    .copy-status { min-width: 7rem; min-height: 1.4em; }
    @media (max-width: 44rem) {
      .toolbar button span { display: none; }
      .prompt, .description { padding: 0.48rem 0.55rem; font-size: 0.58rem; }
      .footer { padding: 0.42rem 0.55rem; }
    }
  `;

  command = "";
  description = "";
  heading = "Command";
  copyLabel = "Copy";
  copiedLabel = "Copied";
  copyFailedLabel = "Copy failed";
  sessionKey = "";
  private copyState: "idle" | "copied" | "failed" = "idle";
  clipboard: Pick<Clipboard, "writeText"> | undefined;

  protected willUpdate(changed: PropertyValues<this>): void {
    if (changed.has("sessionKey") || changed.has("command")) this.reset();
  }

  reset(): void {
    this.copyState = "idle";
  }

  async copyCommand(): Promise<void> {
    const clipboard = this.clipboard ?? globalThis.navigator.clipboard;
    if (!clipboard || !this.command) {
      this.copyState = "failed";
      return;
    }
    try {
      await clipboard.writeText(this.command);
      this.copyState = "copied";
    } catch {
      this.copyState = "failed";
    }
  }

  protected render() {
    const copyStatus = this.copyState === "copied" ? this.copiedLabel : this.copyState === "failed" ? this.copyFailedLabel : "";
    return html`
      <courier-terminal .heading=${this.heading}>
        <div slot="toolbar" class="toolbar">
          <button type="button" ?disabled=${!this.command} aria-label=${this.copyLabel} title=${this.copyLabel} @click=${this.copyCommand}><courier-icon name=${this.copyState === "copied" ? "check" : "copy"}></courier-icon><span>${this.copyLabel}</span></button>
        </div>
        <div class="readout">
          <div class="prompt"><code>${this.command}</code></div>
          ${this.description ? html`<p class="description">${this.description}</p>` : nothing}
          <div class="details"><slot name="details"></slot></div>
        </div>
        <div slot="footer" class="footer"><span class="copy-status" role="status" aria-live="polite">${copyStatus}</span><slot name="footer-actions"></slot></div>
      </courier-terminal>
    `;
  }
}
