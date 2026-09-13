import { css } from "lit";

export const controlStyles = css`
  :host {
    color: var(--courier-text, #151714);
    font-family: var(--courier-font-sans, sans-serif);
  }

  button,
  select {
    min-height: 2.75rem;
    border: 1px solid var(--courier-line, #c7ccc0);
    border-radius: var(--courier-radius-sm, 0.25rem);
    color: inherit;
    background: var(--courier-surface, #fff);
    font: inherit;
  }

  button {
    padding: 0.625rem 1rem;
    cursor: pointer;
    transition: background var(--courier-duration, 160ms) var(--courier-ease, ease);
  }

  button:hover:not(:disabled) {
    background: color-mix(in srgb, var(--courier-signal, #d4ff45) 24%, var(--courier-surface, #fff));
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
    color: var(--courier-muted, #51574d);
    font-size: 0.8125rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
  }
`;
