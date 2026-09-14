import { LitElement, css, html } from "lit";

export interface SegmentOption {
  readonly value: string;
  readonly label: string;
}

export function nextSegmentIndex(key: string, index: number, length: number): number | undefined {
  if (length <= 0) {
    return undefined;
  }
  switch (key) {
    case "ArrowLeft":
    case "ArrowUp":
      return (index - 1 + length) % length;
    case "ArrowRight":
    case "ArrowDown":
      return (index + 1) % length;
    case "Home":
      return 0;
    case "End":
      return length - 1;
    default:
      return undefined;
  }
}

export class CourierSegmentedControl extends LitElement {
  static properties = {
    label: { type: String },
    options: { attribute: false },
    value: { type: String },
  };

  static styles = css`
    :host { display: block; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    fieldset { min-width: 0; margin: 0; padding: 0; border: 0; }
    legend { margin: 0 0 0.25rem; padding: 0; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 750; letter-spacing: 0.08em; line-height: 1; text-transform: uppercase; }
    .segments { display: inline-grid; max-width: 100%; grid-auto-columns: minmax(0, auto); grid-auto-flow: column; gap: 2px; padding: 2px; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-md, 0.625rem); background: color-mix(in srgb, var(--courier-color-field, #e7e9dc) 68%, transparent); box-shadow: inset 0 1px 2px rgb(16 18 15 / 0.07); }
    button { appearance: none; min-width: 0; min-height: 2.25rem; padding: 0.45rem 0.68rem; overflow: hidden; border: 1px solid transparent; border-radius: calc(var(--courier-radius-md, 0.625rem) - 3px); color: var(--courier-color-muted, #596054); background: transparent; font: inherit; font-size: 0.75rem; font-weight: 780; line-height: 1; text-overflow: ellipsis; white-space: nowrap; cursor: pointer; transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    button:hover { color: var(--courier-color-text, #151714); background: color-mix(in srgb, var(--courier-color-surface-raised, #fff) 72%, transparent); }
    button.selected { border-color: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 64%, var(--courier-color-border, #c8cdbf)); color: var(--courier-color-accent-ink, #151714); background: var(--courier-color-accent, #d4ff45); box-shadow: 0 1px 0 rgb(16 18 15 / 0.12); }
    button:active { transform: translateY(1px); }
    button:focus-visible { position: relative; z-index: 1; outline: 3px solid var(--courier-beak, #ff8758); outline-offset: 2px; }
  `;

  label = "";
  options: readonly SegmentOption[] = [];
  value = "";

  private activate(event: Event): void {
    const value = (event.currentTarget as HTMLButtonElement).dataset.value ?? "";
    if (!value || value === this.value) {
      return;
    }
    this.value = value;
    this.dispatchEvent(new CustomEvent<string>("courier-segment-change", { detail: value, bubbles: true, composed: true }));
  }

  private move(event: KeyboardEvent): void {
    const buttons = [...this.renderRoot.querySelectorAll<HTMLButtonElement>("button")];
    const current = buttons.indexOf(event.currentTarget as HTMLButtonElement);
    const next = nextSegmentIndex(event.key, current, buttons.length);
    if (next === undefined) {
      return;
    }
    event.preventDefault();
    const target = buttons[next]!;
    target.focus();
    target.click();
  }

  protected render() {
    return html`<fieldset>
      <legend>${this.label}</legend>
      <div class="segments" role="radiogroup" aria-label=${this.label}>
        ${this.options.map((option) => {
          const selected = option.value === this.value;
          return html`<button
            type="button"
            role="radio"
            class=${selected ? "selected" : ""}
            data-value=${option.value}
            aria-checked=${String(selected)}
            tabindex=${selected ? 0 : -1}
            @click=${this.activate}
            @keydown=${this.move}
          >${option.label}</button>`;
        })}
      </div>
    </fieldset>`;
  }
}
