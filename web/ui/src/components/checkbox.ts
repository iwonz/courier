import { LitElement, css, html } from "lit";

export class CourierCheckbox extends LitElement {
  static properties = {
    checked: { type: Boolean, reflect: true },
    disabled: { type: Boolean, reflect: true },
    label: { type: String },
  };

  static styles = css`
    :host { display: inline-flex; min-width: 0; color: inherit; font-family: var(--courier-font-sans, sans-serif); }
    label { display: inline-grid; min-width: 0; grid-template-columns: 1.25rem minmax(0, 1fr); align-items: center; gap: 0.55rem; color: inherit; font-size: 0.75rem; font-weight: 720; line-height: 1.25; cursor: pointer; }
    input { position: absolute; width: 1px; height: 1px; margin: -1px; overflow: hidden; clip-path: inset(50%); white-space: nowrap; }
    .box { position: relative; display: grid; width: 1.25rem; height: 1.25rem; place-items: center; border: 1px solid currentColor; border-radius: 0.25rem; background: color-mix(in srgb, currentColor 6%, transparent); transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    .box::after { content: ""; width: 0.55rem; height: 0.3rem; border-bottom: 2px solid currentColor; border-left: 2px solid currentColor; opacity: 0; transform: translateY(-0.08rem) rotate(-45deg) scale(0.65); transition: opacity var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease); }
    input:checked + .box { border-color: var(--courier-color-accent-solid, #ff704c); color: var(--courier-color-accent-ink, #091a33); background: var(--courier-color-accent-solid, #ff704c); }
    input:checked + .box::after { opacity: 1; transform: translateY(-0.08rem) rotate(-45deg) scale(1); }
    input:focus-visible + .box { outline: 3px solid var(--courier-color-accent, #c83f23); outline-offset: 2px; }
    label:hover .box { transform: translateY(-1px); }
    :host([disabled]) { opacity: 0.52; }
    :host([disabled]) label { cursor: not-allowed; }
    :host([disabled]) label:hover .box { transform: none; }
  `;

  checked = false;
  disabled = false;
  label = "";

  private change(event: Event): void {
    this.checked = (event.currentTarget as HTMLInputElement).checked;
    this.dispatchEvent(new CustomEvent<boolean>("courier-checkbox-change", {
      detail: this.checked,
      bubbles: true,
      composed: true,
    }));
  }

  protected render() {
    return html`<label>
      <input type="checkbox" .checked=${this.checked} .disabled=${this.disabled} @change=${this.change}>
      <span class="box" aria-hidden="true"></span>
      <span>${this.label}</span>
    </label>`;
  }
}
