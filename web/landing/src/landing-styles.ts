import { css } from "lit";

export const landingStyles = css`
  :host {
    --masthead-height: 4.6rem;
    display: block;
    min-width: 20rem;
    color: var(--courier-color-text);
    background: var(--courier-color-canvas);
    font-family: var(--courier-font-sans);
  }
  .shell { width: min(100% - clamp(1.25rem, 5vw, 5rem), 94rem); margin-inline: auto; }
  .masthead-wrap {
    --courier-color-text: var(--courier-carbon-950);
    --courier-color-muted: #5f5a52;
    --courier-color-accent: var(--courier-route-light);
    --courier-control-frame-border: rgb(74 70 63 / 0.28);
    --courier-control-frame-color: var(--courier-carbon-950);
    --courier-control-frame-surface: rgb(250 248 242 / 0.44);
    --courier-control-frame-hover-border: var(--courier-route-light);
    --courier-control-frame-hover-color: var(--courier-carbon-950);
    --courier-control-frame-hover-surface: rgb(250 248 242 / 0.82);
    --courier-control-frame-focus: var(--courier-route);
    position: fixed;
    z-index: 30;
    top: 0;
    right: 0;
    left: 0;
    padding-top: env(safe-area-inset-top);
    color: var(--courier-carbon-950);
    background: linear-gradient(180deg, rgb(242 239 230 / 0.88) 0%, rgb(242 239 230 / 0.58) 58%, transparent 100%);
    -webkit-backdrop-filter: blur(14px) saturate(0.82);
    backdrop-filter: blur(14px) saturate(0.82);
    -webkit-mask-image: linear-gradient(#000 0 64%, rgb(0 0 0 / 0.72) 78%, transparent 100%);
    mask-image: linear-gradient(#000 0 64%, rgb(0 0 0 / 0.72) 78%, transparent 100%);
  }
  .masthead { display: grid; min-height: 4.6rem; grid-template-columns: auto 1fr auto; align-items: start; gap: 1.5rem; padding-top: 0.7rem; }
  .brand-link { color: inherit; text-decoration: none; }
  nav, .header-actions, .install-actions { display: flex; align-items: center; }
  nav { justify-content: flex-start; gap: 1.2rem; padding-top: 0.7rem; }
  nav a { color: var(--courier-color-muted); font-size: 0.76rem; font-weight: 720; text-decoration: none; }
  nav a:hover, nav a[aria-current="page"] { color: var(--courier-color-accent); }
  .header-actions { gap: 0.38rem; }
  a:focus-visible, button:focus-visible { outline: 3px solid var(--courier-color-accent); outline-offset: 3px; }

  main { position: relative; overflow: clip; isolation: isolate; }
  .page-panorama { z-index: 0; background: var(--courier-color-canvas); }
  .page-panorama::part(image) { object-position: center top; }
  .page-veil {
    position: absolute;
    z-index: 1;
    inset: 0;
    background:
      linear-gradient(90deg, color-mix(in srgb, var(--courier-color-canvas) 68%, transparent) 0%, color-mix(in srgb, var(--courier-color-canvas) 35%, transparent) 38%, transparent 72%),
      color-mix(in srgb, var(--courier-color-canvas) 18%, transparent);
    pointer-events: none;
  }
  .stage {
    position: relative;
    z-index: 2;
    min-height: max(42rem, 78svh);
    padding: calc(var(--masthead-height) + 3rem) 0 5rem;
    scroll-margin-top: var(--masthead-height);
    content-visibility: auto;
    contain-intrinsic-size: auto 54rem;
  }
  .route-stage { min-height: 100svh; padding-bottom: 3rem; content-visibility: visible; }
  #cli { min-height: max(52rem, 88svh); }
  .stage-shell { position: relative; display: grid; min-height: 0; align-content: start; gap: clamp(1.5rem, 4vh, 3rem); }
  h1, h2, p { margin: 0; }
  h1, h2 {
    position: relative;
    z-index: 0;
    width: fit-content;
    color: var(--courier-color-text);
    font-family: var(--courier-font-display);
    text-wrap: balance;
    text-shadow:
      0 1px color-mix(in srgb, var(--courier-color-canvas) 88%, transparent),
      0 0 1.1rem color-mix(in srgb, var(--courier-color-canvas) 88%, transparent),
      0 0 2.8rem color-mix(in srgb, var(--courier-color-canvas) 72%, transparent);
    isolation: isolate;
  }
  h1::before, h2::before {
    content: "";
    position: absolute;
    z-index: -1;
    inset: -0.22em -0.36em;
    border-radius: 999px;
    background: color-mix(in srgb, var(--courier-color-canvas) 54%, transparent);
    filter: blur(1.15rem);
    pointer-events: none;
  }
  h1 { max-width: 11ch; font-size: clamp(4.2rem, 10vw, 9rem); font-weight: 840; letter-spacing: -0.072em; line-height: 0.88; }
  h2 { font-size: clamp(2.6rem, 6vw, 5.4rem); font-weight: 810; letter-spacing: -0.06em; line-height: 0.95; }

  .route-shell { min-height: calc(100svh - var(--masthead-height) - 6rem); grid-template-rows: minmax(12rem, 1fr) auto; }
  .route-instrument { display: grid; gap: 0.9rem; align-self: end; }
  .route-controls { position: relative; display: grid; grid-template-columns: 1fr 1fr; gap: clamp(2rem, 9vw, 10rem); padding: 0.5rem 0; }
  .route-connector { position: absolute; z-index: 0; inset: 0; width: 100%; height: 100%; overflow: visible; pointer-events: none; }
  .route-connector path { fill: none; vector-effect: non-scaling-stroke; }
  .route-path { stroke: color-mix(in srgb, var(--courier-color-text) 38%, transparent); stroke-width: 1.25; }
  .route-signal { stroke: var(--courier-color-accent-solid); stroke-width: 2; stroke-dasharray: 1 10; stroke-linecap: round; animation: route-pulse 3s linear infinite; }
  @keyframes route-pulse { to { stroke-dashoffset: -22; } }
  .endpoint-group { position: relative; z-index: 1; display: grid; min-width: 0; grid-template-columns: auto repeat(4, minmax(0, 1fr)); align-items: center; gap: 0.35rem; }
  .endpoint-group > strong { color: var(--courier-color-muted); font-size: 0.68rem; font-weight: 720; }
  .endpoint {
    appearance: none;
    display: inline-grid;
    min-width: 0;
    min-height: 2.35rem;
    grid-template-columns: 1.35rem minmax(0, 1fr);
    align-items: center;
    gap: 0.38rem;
    padding: 0.25rem 0.55rem 0.25rem 0.28rem;
    border: 1px solid color-mix(in srgb, var(--courier-color-border-strong) 56%, transparent);
    border-radius: 999px;
    color: var(--courier-color-text);
    background: color-mix(in srgb, var(--courier-color-canvas) 72%, transparent);
    box-shadow: 0 0.35rem 1.2rem rgb(14 15 13 / 0.06);
    text-align: left;
    cursor: pointer;
    -webkit-backdrop-filter: blur(14px) saturate(0.8);
    backdrop-filter: blur(14px) saturate(0.8);
    transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease);
  }
  .endpoint-terminal { display: inline-grid; width: 1.35rem; height: 1.35rem; aspect-ratio: 1; place-items: center; border: 1px solid currentColor; border-radius: 50%; }
  .endpoint courier-icon { width: 0.72rem; height: 0.72rem; }
  .endpoint span:last-child { overflow: hidden; font-size: 0.65rem; font-weight: 760; text-overflow: ellipsis; white-space: nowrap; }
  .endpoint:hover:not(:disabled), .endpoint:focus-visible { border-color: var(--courier-color-text); transform: translateY(-1px); }
  .endpoint.selected { border-color: var(--courier-color-accent); color: var(--courier-color-accent); background: color-mix(in srgb, var(--courier-color-canvas) 90%, transparent); }
  .endpoint:disabled { opacity: 0.32; cursor: not-allowed; }
  .route-readout { height: 10.5rem; }
  .readout-details { display: flex; min-height: 0; align-items: center; gap: 0.6rem; padding: 0.55rem 0.8rem; border-top: 1px solid var(--courier-workbench-line); }
  .flag-list { display: flex; min-width: 0; gap: 0.3rem; overflow: auto; flex-wrap: wrap; }
  .flag { padding: 0.2rem 0.42rem; border: 1px solid var(--courier-workbench-line); border-radius: 999px; color: var(--courier-color-muted); background: color-mix(in srgb, var(--courier-color-canvas) 34%, transparent); font-family: var(--courier-font-mono); font-size: 0.56rem; white-space: nowrap; }

  .install-interface { display: grid; gap: 1rem; }
  .install-channels { display: flex; gap: 0.4rem; flex-wrap: wrap; }
  .install-channel {
    appearance: none;
    display: inline-flex;
    min-height: 2rem;
    align-items: center;
    gap: 0.38rem;
    padding: 0.28rem 0.62rem;
    border: 1px solid color-mix(in srgb, var(--courier-color-border-strong) 56%, transparent);
    border-radius: 999px;
    color: var(--courier-color-text);
    background: color-mix(in srgb, var(--courier-color-canvas) 68%, transparent);
    box-shadow: 0 0.3rem 1rem rgb(14 15 13 / 0.05);
    cursor: pointer;
    -webkit-backdrop-filter: blur(12px) saturate(0.8);
    backdrop-filter: blur(12px) saturate(0.8);
    transition: border-color var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease);
  }
  .install-channel:hover, .install-channel:focus-visible { border-color: var(--courier-color-text); transform: translateY(-1px); }
  .install-channel.selected { border-color: var(--courier-color-accent); color: var(--courier-color-accent); }
  .install-channel courier-brand-icon { width: 0.95rem; height: 0.95rem; }
  .install-channel span { font-size: 0.68rem; font-weight: 740; }
  .install-readout { height: 11.5rem; }
  .install-identity { color: var(--courier-color-accent); font-size: 0.7rem; }
  .install-identity courier-brand-icon { width: 1rem; height: 1rem; }
  .install-actions { gap: 0.9rem; }
  .install-actions a { display: inline-flex; align-items: center; gap: 0.35rem; color: var(--courier-color-muted); font-weight: 720; text-decoration: none; }
  .install-actions a:hover { color: var(--courier-color-accent); }
  .install-actions courier-icon { width: 0.9rem; height: 0.9rem; }

  .cli-workspace { display: grid; height: 34rem; min-height: 34rem; grid-template-columns: minmax(0, 1.25fr) minmax(20rem, 0.75fr); gap: 1rem; }
  .reference-grid { display: grid; height: 100%; min-height: 0; grid-template-columns: minmax(13rem, 0.7fr) minmax(0, 1.3fr); }
  .reference-column { min-width: 0; min-height: 0; overflow: hidden; }
  .reference-column + .reference-column { border-left: 1px solid var(--courier-workbench-line); }
  .reference-list { height: 100%; min-height: 0; overflow-y: auto; overscroll-behavior: contain; }
  .reference-filter { display: flex; min-height: 2.8rem; align-items: center; justify-content: space-between; gap: 0.7rem; padding: 0.35rem 0.65rem; border-bottom: 1px solid var(--courier-workbench-line); color: var(--courier-color-muted); font-size: 0.68rem; font-weight: 720; }
  .command-row, .option-row { display: grid; width: 100%; min-width: 0; gap: 0.34rem; padding: 0.68rem; border: 0; border-bottom: 1px solid color-mix(in srgb, var(--courier-workbench-line) 65%, transparent); color: inherit; background: transparent; text-align: left; }
  .command-row { appearance: none; cursor: pointer; }
  .command-row:hover, .command-row[aria-pressed="true"] { color: var(--courier-color-accent); background: color-mix(in srgb, var(--courier-color-accent-solid) 6%, transparent); }
  .command-row code, .option-row code { overflow-wrap: anywhere; font: 700 0.68rem/1.4 var(--courier-font-mono); }
  .option-row > div { display: flex; align-items: baseline; justify-content: space-between; gap: 0.6rem; }
  .option-row span, .option-row p { color: var(--courier-color-muted); font-size: 0.58rem; line-height: 1.4; }
  .empty-state { display: grid; min-height: 100%; place-items: center; padding: 1rem; color: var(--courier-color-muted); font-size: 0.72rem; text-align: center; }
  .cli-readout { min-height: 0; }

  @media (max-width: 70rem) {
    .endpoint-group { grid-template-columns: auto repeat(2, minmax(0, 1fr)); }
    .endpoint-group > strong { grid-row: 1 / 3; }
    .cli-workspace { grid-template-columns: minmax(0, 1fr) minmax(17rem, 0.72fr); }
  }
  @media (max-width: 44rem) {
    :host { --masthead-height: 4.2rem; }
    .shell { width: min(100% - 1rem, 94rem); }
    .masthead { min-height: 4.2rem; grid-template-columns: auto 1fr auto; gap: 0.55rem; padding-top: 0.5rem; }
    nav { justify-content: flex-end; gap: 0.7rem; }
    nav a { font-size: 0.66rem; }
    .header-actions { gap: 0.2rem; --courier-control-frame-size: 2.1rem; }
    .page-veil { background: linear-gradient(90deg, color-mix(in srgb, var(--courier-color-canvas) 48%, transparent), color-mix(in srgb, var(--courier-color-canvas) 16%, transparent)); }
    .stage { min-height: 44rem; padding: calc(var(--masthead-height) + 1.5rem) 0 3.5rem; }
    .route-stage { min-height: 100svh; }
    h1 { max-width: 8ch; font-size: clamp(3.4rem, 18vw, 5.5rem); letter-spacing: -0.064em; line-height: 0.9; }
    h2 { font-size: clamp(2.4rem, 12vw, 4rem); }
    .route-shell { min-height: 0; grid-template-rows: auto auto; gap: 2rem; }
    .route-controls { gap: 0.7rem; }
    .endpoint-group { grid-template-columns: minmax(0, 1fr); grid-template-rows: auto repeat(4, 1.75rem); gap: 0.2rem; }
    .endpoint-group > strong { grid-row: auto; font-size: 0.58rem; }
    .endpoint { min-height: 1.75rem; grid-template-columns: 1.2rem minmax(0, 1fr); padding: 0.12rem 0.3rem 0.12rem 0.14rem; }
    .endpoint-terminal { width: 1.2rem; height: 1.2rem; }
    .endpoint span:last-child { font-size: 0.56rem; }
    .route-readout, .install-readout, .cli-readout { height: 12.5rem; }
    .install-channel { min-height: 1.75rem; padding: 0.2rem 0.42rem; }
    .install-channel courier-brand-icon { width: 0.76rem; height: 0.76rem; }
    .install-channel span { font-size: 0.57rem; }
    .install-actions { gap: 0.5rem; }
    .install-actions a { font-size: 0.56rem; }
    .cli-workspace { height: 54rem; min-height: 54rem; grid-template-columns: 1fr; grid-template-rows: minmax(31rem, 1fr) auto; gap: 0.7rem; }
    .reference-grid { grid-template-columns: minmax(7.5rem, 0.76fr) minmax(0, 1.24fr); }
    .reference-filter { display: grid; min-height: 4rem; padding: 0.32rem 0.42rem; }
    .command-row, .option-row { padding: 0.48rem 0.42rem; }
    .command-row code, .option-row code { font-size: 0.56rem; }
    .option-row > div { display: grid; gap: 0.15rem; }
    .option-row span, .option-row p { font-size: 0.5rem; }
    courier-brand::part(wordmark) { display: none; }
  }
  @media (max-width: 25rem) { nav { display: none; } }
  @media (prefers-reduced-motion: reduce) { .route-signal { animation: none; stroke-dasharray: none; opacity: 0.78; } }
`;
