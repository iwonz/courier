import { css } from "lit";

export const controlStyles = css`
  :host {
    color: var(--courier-color-text, #151714);
    font-family: var(--courier-font-sans, sans-serif);
  }

  button,
  select {
    min-height: 2.75rem;
    border: 1px solid var(--courier-color-border, #c8cdbf);
    border-radius: var(--courier-radius-sm, 0.25rem);
    color: inherit;
    background: var(--courier-color-surface-raised, #fff);
    font: inherit;
  }

  button {
    padding: 0.625rem 1rem;
    box-shadow: 0 1px 0 rgb(16 18 15 / 0.08);
    cursor: pointer;
    font-weight: 750;
    letter-spacing: -0.01em;
    transition: background var(--courier-duration, 160ms) var(--courier-ease, ease), border-color var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease);
  }

  button:hover:not(:disabled) {
    border-color: var(--courier-color-border-strong, #8e9587);
    background: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 18%, var(--courier-color-surface-raised, #fff));
    transform: translateY(-1px);
  }

  button:active:not(:disabled) {
    transform: translateY(0);
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  select {
    padding: 0.5rem 2rem 0.5rem 0.75rem;
  }

  button:focus-visible,
  select:focus-visible {
    outline: 3px solid var(--courier-beak, #ff8758);
    outline-offset: 2px;
  }
`;

export const fieldStyles = css`
  label {
    display: grid;
    gap: var(--courier-space-1, 0.25rem);
    color: var(--courier-color-muted, #596054);
    font-family: var(--courier-font-mono, monospace);
    font-size: 0.6875rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }
`;
