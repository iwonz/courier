import { LitElement, css, html, nothing, type PropertyValues } from "lit";

export class CourierWorkbench extends LitElement {
  static properties = { heading: { type: String }, status: { type: String } };
  static styles = css`
    :host { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(0, 1fr) auto; overflow: hidden; border-top: 1px solid var(--courier-workbench-line, #cbc5b8); border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-text, #0e0f0d); background: var(--courier-workbench-surface, rgb(242 239 230 / 0.78)); backdrop-filter: blur(16px); }
    header { display: flex; min-height: 2.75rem; align-items: center; justify-content: space-between; gap: 1rem; padding: 0.55rem 0.2rem; border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-muted, #68645e); font-size: 0.72rem; font-weight: 720; }
    .status { color: var(--courier-color-accent, #ad431d); font-family: var(--courier-font-mono, monospace); }
    .body { min-width: 0; min-height: 0; overflow: auto; }
    footer { border-top: 1px solid var(--courier-workbench-line, #cbc5b8); }
  `;
  heading = "";
  status = "";
  protected render() {
    return html`<header><span>${this.heading}</span><slot name="actions"></slot>${this.status ? html`<span class="status">${this.status}</span>` : nothing}</header><div class="body"><slot></slot></div><footer><slot name="footer"></slot></footer>`;
  }
}

export class CourierCommandReadout extends LitElement {
  static properties = { command: { type: String }, description: { type: String }, heading: { type: String }, copyLabel: { type: String, attribute: "copy-label" }, copiedLabel: { type: String, attribute: "copied-label" }, copyFailedLabel: { type: String, attribute: "copy-failed-label" }, sessionKey: { type: String, attribute: "session-key" }, copyState: { state: true } };
  static styles = css`
    :host { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(0, 1fr) auto; border-top: 1px solid var(--courier-workbench-line, #cbc5b8); border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-text, #0e0f0d); background: var(--courier-workbench-surface, rgb(242 239 230 / 0.78)); backdrop-filter: blur(16px); }
    header { display: flex; min-height: 2.75rem; align-items: center; justify-content: space-between; gap: 1rem; padding: 0.48rem 0.2rem; border-bottom: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-muted, #68645e); font-size: 0.72rem; font-weight: 720; }
    button { appearance: none; display: inline-flex; min-height: 2rem; align-items: center; gap: 0.38rem; padding: 0.3rem 0.65rem; border: 1px solid var(--courier-control-frame-border, #8f8a81); border-radius: 999px; color: var(--courier-color-text, #0e0f0d); background: transparent; font: 700 0.68rem/1 var(--courier-font-sans, sans-serif); cursor: pointer; }
    button:hover:not(:disabled) { border-color: var(--courier-color-accent, #ad431d); color: var(--courier-color-accent, #ad431d); }
    button:focus-visible { outline: 3px solid var(--courier-color-accent, #ad431d); outline-offset: 2px; }
    button:disabled { cursor: not-allowed; opacity: 0.42; }
    courier-icon { width: 0.9rem; height: 0.9rem; }
    .readout { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto auto minmax(0, 1fr); }
    .command, .description { min-width: 0; padding: 0.75rem 0.2rem; }
    .command { color: var(--courier-color-text, #0e0f0d); }
    code { overflow-wrap: anywhere; font: 700 0.78rem/1.5 var(--courier-font-mono, monospace); white-space: pre-wrap; }
    .description { margin: 0; border-top: 1px solid color-mix(in srgb, var(--courier-workbench-line, #cbc5b8) 65%, transparent); color: var(--courier-color-muted, #68645e); font-size: 0.72rem; line-height: 1.45; }
    .details { min-width: 0; min-height: 0; overflow: auto; }
    footer { display: flex; min-width: 0; min-height: 2.4rem; align-items: center; justify-content: space-between; gap: 0.75rem; padding: 0.4rem 0.2rem; border-top: 1px solid var(--courier-workbench-line, #cbc5b8); color: var(--courier-color-muted, #68645e); font-size: 0.62rem; }
    .copy-status { min-width: 6rem; min-height: 1.4em; }
    @media (max-width: 44rem) { header button span { display: none; } .command, .description { padding: 0.6rem 0.1rem; } footer { padding-inline: 0.1rem; } }
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
  protected willUpdate(changed: PropertyValues<this>): void { if (changed.has("sessionKey") || changed.has("command")) this.reset(); }
  reset(): void { this.copyState = "idle"; }
  async copyCommand(): Promise<void> {
    const clipboard = this.clipboard ?? globalThis.navigator.clipboard;
    if (!clipboard || !this.command) { this.copyState = "failed"; return; }
    try { await clipboard.writeText(this.command); this.copyState = "copied"; } catch { this.copyState = "failed"; }
  }
  protected render() {
    const copyStatus = this.copyState === "copied" ? this.copiedLabel : this.copyState === "failed" ? this.copyFailedLabel : "";
    return html`<header><span>${this.heading}</span><button type="button" ?disabled=${!this.command} aria-label=${this.copyLabel} title=${this.copyLabel} @click=${this.copyCommand}><courier-icon name=${this.copyState === "copied" ? "check" : "copy"}></courier-icon><span>${this.copyLabel}</span></button></header><div class="readout"><div class="command"><code>${this.command}</code></div>${this.description ? html`<p class="description">${this.description}</p>` : nothing}<div class="details"><slot name="details"></slot></div></div><footer><span class="copy-status" role="status" aria-live="polite">${copyStatus}</span><slot name="footer-actions"></slot></footer>`;
  }
}
