import { LitElement, css, html, nothing, type PropertyValues } from "lit";

export type TerminalTone = "neutral" | "signal" | "success" | "warning" | "danger";

export interface TerminalStep {
  readonly label: string;
  readonly detail?: string;
  readonly tone?: TerminalTone;
}

export type DemoPhase = "idle" | "running" | "complete";

export function terminalTimestamp(index: number): string {
  const seconds = Math.max(0, index) * 2;
  return `00:${String(seconds).padStart(2, "0")}`;
}

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

export class CourierCommandDemo extends LitElement {
  static properties = {
    command: { type: String },
    description: { type: String },
    steps: { attribute: false },
    copyLabel: { type: String, attribute: "copy-label" },
    runLabel: { type: String, attribute: "run-label" },
    replayLabel: { type: String, attribute: "replay-label" },
    copiedLabel: { type: String, attribute: "copied-label" },
    copyFailedLabel: { type: String, attribute: "copy-failed-label" },
    previewLabel: { type: String, attribute: "preview-label" },
    noEffectLabel: { type: String, attribute: "no-effect-label" },
    sessionKey: { type: String, attribute: "session-key" },
    interval: { type: Number },
    phase: { state: true },
    visibleCount: { state: true },
    copyState: { state: true },
  };

  static styles = css`
    :host { display: block; min-width: 0; min-height: 0; }
    courier-terminal { height: 100%; }
    .toolbar { display: flex; align-items: center; gap: 0.4rem; }
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
    .session { display: grid; height: 100%; min-height: 0; grid-template-rows: auto auto auto minmax(0, 1fr); }
    .prompt, .description, li, .empty { min-width: 0; padding: 0.58rem 0.8rem; }
    .prompt { display: grid; grid-template-columns: auto minmax(0, 1fr); gap: 0.65rem; color: var(--courier-terminal-text, #f3f4e9); }
    .prompt::before { content: "$"; color: var(--courier-terminal-prompt, #d4ff45); font-weight: 800; }
    code { overflow-wrap: anywhere; font: inherit; line-height: 1.45; white-space: pre-wrap; }
    .description { border-top: 1px solid color-mix(in srgb, var(--courier-terminal-border, #596253) 58%, transparent); color: var(--courier-terminal-muted, #b9c0b1); font-size: 0.68rem; line-height: 1.45; }
    ol { min-height: 0; max-height: 11rem; margin: 0; padding: 0; overflow: auto; list-style: none; }
    li { display: grid; grid-template-columns: 2.6rem minmax(7rem, 0.35fr) minmax(0, 1fr); gap: 0.7rem; border-top: 1px solid color-mix(in srgb, var(--courier-terminal-border, #596253) 45%, transparent); font-size: 0.65rem; line-height: 1.4; }
    time { color: var(--courier-terminal-muted, #b9c0b1); font-variant-numeric: tabular-nums; }
    .step-label { color: var(--courier-terminal-text, #f3f4e9); font-weight: 760; }
    .step-detail { color: var(--courier-terminal-muted, #b9c0b1); overflow-wrap: anywhere; }
    li[data-tone="signal"] .step-label, li[data-tone="success"] .step-label { color: var(--courier-terminal-prompt, #d4ff45); }
    li[data-tone="warning"] .step-label { color: var(--courier-warning, #f0b849); }
    li[data-tone="danger"] .step-label { color: var(--courier-danger, #ff6b5f); }
    .empty { color: var(--courier-terminal-muted, #b9c0b1); font-size: 0.65rem; }
    .details { min-width: 0; }
    .footer { display: flex; min-width: 0; align-items: center; justify-content: space-between; gap: 0.75rem; padding: 0.48rem 0.8rem; color: var(--courier-terminal-muted, #b9c0b1); font-size: 0.58rem; line-height: 1.4; }
    @media (max-width: 44rem) {
      .toolbar button span { display: none; }
      li { grid-template-columns: 2.25rem minmax(5.5rem, 0.42fr) minmax(0, 1fr); gap: 0.4rem; padding: 0.45rem 0.55rem; font-size: 0.56rem; }
      .prompt, .description, .empty { padding: 0.48rem 0.55rem; font-size: 0.58rem; }
      ol { max-height: 8rem; }
    }
  `;

  command = "";
  description = "";
  steps: readonly TerminalStep[] = [];
  copyLabel = "Copy";
  runLabel = "Run demo";
  replayLabel = "Replay";
  copiedLabel = "Copied";
  copyFailedLabel = "Copy failed";
  previewLabel = "Preview";
  noEffectLabel = "Preview only. No command or transfer ran in this browser.";
  sessionKey = "";
  interval = 180;
  phase: DemoPhase = "idle";
  private visibleCount = 0;
  private copyState: "idle" | "copied" | "failed" = "idle";
  private timer?: number;
  clipboard: Pick<Clipboard, "writeText"> | undefined;

  disconnectedCallback(): void {
    this.stopTimer();
    super.disconnectedCallback();
  }

  protected willUpdate(changed: PropertyValues<this>): void {
    if (changed.has("sessionKey") && changed.get("sessionKey") !== undefined) this.reset();
  }

  reset(): void {
    this.stopTimer();
    this.phase = "idle";
    this.visibleCount = 0;
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

  run(): void {
    this.stopTimer();
    this.copyState = "idle";
    if (!this.command || this.steps.length === 0) {
      this.phase = "idle";
      this.visibleCount = 0;
      return;
    }
    if (globalThis.matchMedia?.("(prefers-reduced-motion: reduce)").matches) {
      this.visibleCount = this.steps.length;
      this.phase = "complete";
      return;
    }
    this.visibleCount = 1;
    this.phase = this.steps.length === 1 ? "complete" : "running";
    if (this.phase === "running") this.scheduleStep();
  }

  private scheduleStep(): void {
    this.timer = globalThis.setTimeout(() => {
      this.timer = undefined;
      this.visibleCount += 1;
      if (this.visibleCount >= this.steps.length) {
        this.phase = "complete";
        return;
      }
      this.scheduleStep();
    }, Math.max(0, this.interval));
  }

  private stopTimer(): void {
    if (this.timer === undefined) return;
    globalThis.clearTimeout(this.timer);
    this.timer = undefined;
  }

  protected render() {
    const copyStatus = this.copyState === "copied" ? this.copiedLabel : this.copyState === "failed" ? this.copyFailedLabel : "";
    const visible = this.steps.slice(0, this.visibleCount);
    return html`
      <courier-terminal .heading=${this.previewLabel} .status=${this.phase === "complete" ? "exit 0" : this.phase}>
        <div slot="toolbar" class="toolbar">
          <button type="button" ?disabled=${!this.command} aria-label=${this.copyLabel} title=${this.copyLabel} @click=${this.copyCommand}><courier-icon name=${this.copyState === "copied" ? "check" : "copy"}></courier-icon><span>${this.copyLabel}</span></button>
          <button type="button" ?disabled=${!this.command || this.steps.length === 0} aria-label=${this.phase === "idle" ? this.runLabel : this.replayLabel} title=${this.phase === "idle" ? this.runLabel : this.replayLabel} @click=${this.run}><courier-icon name="terminal"></courier-icon><span>${this.phase === "idle" ? this.runLabel : this.replayLabel}</span></button>
        </div>
        <div class="session">
          <div class="prompt"><code>${this.command}</code></div>
          ${this.description ? html`<p class="description">${this.description}</p>` : nothing}
          <div class="details"><slot name="details"></slot></div>
          ${visible.length ? html`<ol aria-live="polite">${visible.map((step, index) => html`<li data-tone=${step.tone ?? "neutral"}><time>${terminalTimestamp(index)}</time><span class="step-label">${step.label}</span><span class="step-detail">${step.detail ?? ""}</span></li>`)}</ol>` : html`<p class="empty" aria-live="polite">${copyStatus}</p>`}
        </div>
        <div slot="footer" class="footer" role="status"><span>${copyStatus || this.noEffectLabel}</span><slot name="footer-actions"></slot></div>
      </courier-terminal>
    `;
  }
}
