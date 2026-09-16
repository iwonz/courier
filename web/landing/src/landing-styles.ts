import { css } from "lit";

export const landingStyles = css`
  :host {
    --masthead-height: 5rem;
    display: block;
    width: 100%;
    min-height: 100svh;
    overflow-x: clip;
    color: var(--courier-color-text);
    background: var(--courier-color-canvas);
    font-family: var(--courier-font-sans);
  }
  * { box-sizing: border-box; }
  a { color: inherit; }
  button { font: inherit; }
  code, pre { font-family: var(--courier-font-mono); }
  h1, h2, h3, p, pre { margin: 0; }
  .shell { width: min(90rem, calc(100% - clamp(1.25rem, 6vw, 6rem))); margin: 0 auto; }

  .masthead-wrap {
    position: fixed;
    z-index: 100;
    top: 0;
    right: 0;
    left: 0;
    padding-top: env(safe-area-inset-top, 0px);
    padding-right: env(safe-area-inset-right, 0px);
    padding-left: env(safe-area-inset-left, 0px);
    isolation: isolate;
    pointer-events: none;
  }
  .masthead-wrap::before {
    content: "";
    position: absolute;
    z-index: -1;
    inset: 0 0 -1.7rem;
    background: linear-gradient(180deg, rgb(13 16 13 / 0.9) 0%, rgb(13 16 13 / 0.72) 58%, transparent 100%);
    backdrop-filter: blur(16px) saturate(0.82);
    -webkit-backdrop-filter: blur(16px) saturate(0.82);
    mask-image: linear-gradient(180deg, #000 0 62%, transparent 100%);
    pointer-events: none;
  }
  .masthead {
    --courier-color-text: #f3f4e9;
    --courier-color-muted: #b9c0b1;
    --courier-color-surface-raised: #272b25;
    --courier-color-field: #20241f;
    --courier-color-border: #4e574a;
    --courier-color-border-strong: #899381;
    --courier-color-accent: #d4ff45;
    display: flex;
    min-height: 4.7rem;
    align-items: center;
    gap: 1.5rem;
    padding: 0.5rem 0;
    color: var(--courier-paper-50);
    pointer-events: auto;
  }
  .header-left, nav, .header-actions, .preferences, .section-heading, .install-actions, .brand-cloud { display: flex; align-items: center; }
  .header-left { min-width: 0; gap: clamp(1rem, 3vw, 2.5rem); }
  .brand-link { min-width: 0; text-decoration: none; }
  nav { gap: clamp(0.75rem, 2vw, 1.45rem); }
  nav > a { position: relative; padding: 0.5rem 0; color: #b9c0b1; font-size: 0.8125rem; font-weight: 780; text-decoration: none; }
  nav > a::after { content: ""; position: absolute; right: 0; bottom: 0.2rem; left: 0; height: 2px; border-radius: 999px; background: var(--courier-signal); opacity: 0; transform: scaleX(0.4); transition: opacity var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease); }
  nav > a:hover, nav > a[aria-current="page"] { color: var(--courier-paper-50); }
  nav > a[aria-current="page"]::after { opacity: 1; transform: scaleX(1); }
  .header-actions { margin-left: auto; gap: 0.55rem; }
  .preferences { gap: 0.4rem; }
  .github-link { display: inline-flex; width: 2.35rem; height: 2.35rem; align-items: center; justify-content: center; border: 1px solid rgb(185 192 177 / 0.28); border-radius: var(--courier-radius-md); color: var(--courier-paper-50); background: rgb(31 35 29 / 0.62); transition: border-color var(--courier-duration) var(--courier-ease), color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease); }
  .github-link:hover { border-color: var(--courier-signal); color: var(--courier-signal); background: rgb(42 47 38 / 0.86); }
  .github-link courier-icon { width: 1.2rem; height: 1.2rem; }

  .slide {
    position: relative;
    display: grid;
    width: 100%;
    height: 100svh;
    min-height: 30rem;
    padding: calc(var(--masthead-height) + clamp(0.55rem, 1.8vh, 1.35rem)) 0 clamp(0.65rem, 2vh, 1.5rem);
    overflow: hidden;
    isolation: isolate;
    scroll-snap-align: start;
    scroll-snap-stop: always;
  }
  courier-scene { z-index: 0; }
  .slide-shell { position: relative; z-index: 2; display: grid; height: 100%; min-height: 0; grid-template-rows: auto minmax(0, 1fr); gap: clamp(0.7rem, 1.7vh, 1.25rem); }
  .hero-slide { padding-bottom: 0; }
  .hero { position: relative; z-index: 2; display: grid; height: 100%; min-height: 0; align-items: center; pointer-events: none; }
  .hero-copy { position: relative; z-index: 3; display: grid; width: min(51%, 44rem); min-width: 0; gap: clamp(0.85rem, 2.2vh, 1.35rem); }
  h1 { max-width: 44rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(4rem, 7.3vw, 7.25rem); font-weight: 850; letter-spacing: -0.078em; line-height: 0.84; }
  h2 { max-width: 52rem; overflow-wrap: anywhere; font-family: var(--courier-font-display); font-size: clamp(1.85rem, 4vw, 3.8rem); font-weight: 820; letter-spacing: -0.055em; line-height: 0.94; }
  h3 { font-size: 1rem; letter-spacing: -0.02em; }
  .tagline { max-width: 42rem; font-size: clamp(1.25rem, 2.35vw, 2.1rem); font-weight: 650; letter-spacing: -0.035em; line-height: 1.12; }
  .subline, .intro { max-width: 50rem; color: var(--courier-color-muted); font-size: clamp(0.84rem, 1.1vw, 1rem); line-height: 1.5; }
  .hero-shade { position: absolute; z-index: 1; inset: 0; background: linear-gradient(90deg, color-mix(in srgb, var(--courier-color-canvas) 96%, transparent) 0 24%, color-mix(in srgb, var(--courier-color-canvas) 72%, transparent) 39%, transparent 61%), linear-gradient(180deg, rgb(9 11 9 / 0.12), transparent 27% 82%, rgb(9 11 9 / 0.24)); pointer-events: none; }
  .hero-route { position: absolute; z-index: 3; inset: 0; width: 100%; height: 100%; overflow: visible; pointer-events: none; }
  .hero-route.mobile-route { display: none; }
  .hero-route path { fill: none; stroke: var(--courier-signal); stroke-width: 0.34; vector-effect: non-scaling-stroke; filter: drop-shadow(0 0 0.65rem rgb(212 255 69 / 0.52)); }
  .hero-route .signal { stroke-width: 0.48; stroke-dasharray: 2.2 8; animation: route-signal 2.8s linear infinite; }
  .hero-route .terminal-ring { fill: rgb(12 15 12 / 0.84); stroke: var(--courier-signal); stroke-width: 1.25; vector-effect: non-scaling-stroke; filter: drop-shadow(0 0 0.85rem rgb(212 255 69 / 0.74)); }
  .hero-route .terminal-core { fill: var(--courier-signal); stroke: var(--courier-paper-50); stroke-width: 0.45; vector-effect: non-scaling-stroke; }
  .hero-node { position: absolute; z-index: 4; display: grid; gap: 0.12rem; color: var(--courier-paper-50); font-family: var(--courier-font-mono); text-shadow: 0 0.15rem 0.9rem rgb(5 7 5 / 0.92); pointer-events: none; }
  .hero-node small { color: var(--courier-signal); font-size: 0.57rem; font-weight: 800; letter-spacing: 0.14em; }
  .hero-node strong { font-size: 0.75rem; font-weight: 820; letter-spacing: 0.08em; text-transform: uppercase; }
  .hero-node.source { left: 54%; top: 70%; margin: 1.7rem 0 0 -1rem; }
  .hero-node.destination { right: 9%; top: 38%; margin: -2.55rem -0.25rem 0 0; text-align: right; }
  @keyframes route-signal { to { stroke-dashoffset: -20.4; } }
  .endpoint:focus-visible, .install-channel:focus-visible, .command-row:focus-visible, .github-link:focus-visible, a:focus-visible { outline: 3px solid var(--courier-beak); outline-offset: 3px; }

  .section-heading { min-width: 0; justify-content: space-between; align-items: end; color: var(--courier-paper-50); text-shadow: 0 0.14rem 1rem rgb(10 12 10 / 0.72); }
  .section-heading > div { display: grid; min-width: 0; gap: 0.35rem; }
  .section-index, .label, .kind { color: #b9c0b1; font-family: var(--courier-font-mono); font-size: 0.65rem; font-weight: 750; letter-spacing: 0.09em; text-transform: uppercase; }
  .section-heading .section-index { color: var(--courier-signal); }
  .section-heading .intro { color: #d8dccf; }
  .route-explorer, .install-board, .reference { width: 100%; height: 100%; min-height: 0; overflow: hidden; border: 1px solid rgb(177 186 166 / 0.34); border-radius: var(--courier-radius-lg); color: var(--courier-paper-50); background: linear-gradient(135deg, rgb(10 13 10 / 0.38), rgb(19 22 18 / 0.52)); box-shadow: 0 1.2rem 3rem rgb(7 9 7 / 0.18), inset 0 1px rgb(255 255 255 / 0.045); backdrop-filter: blur(8px) saturate(0.88); }

  .route-interface { display: grid; height: 100%; min-width: 0; min-height: 0; grid-template-columns: minmax(0, 1.08fr) minmax(18rem, 0.92fr); gap: clamp(0.8rem, 2vw, 1.5rem); padding: clamp(0.8rem, 2vw, 1.5rem); }
  .route-controls { position: relative; display: grid; min-width: 0; min-height: 0; grid-template-columns: minmax(0, 1fr) minmax(0, 1fr); align-items: center; gap: clamp(2rem, 8vw, 7rem); }
  .route-connector { position: absolute; z-index: 0; inset: 0; width: 100%; height: 100%; overflow: visible; pointer-events: none; }
  .route-connector path { fill: none; stroke: var(--courier-signal); stroke-width: 2; filter: drop-shadow(0 0 0.35rem rgb(212 255 69 / 0.4)); }
  .endpoint-group { position: relative; z-index: 1; display: grid; min-width: 0; min-height: 0; grid-template-rows: auto repeat(4, 2.55rem); align-content: center; gap: clamp(0.3rem, 0.7vh, 0.45rem); }
  .endpoint { appearance: none; display: grid; min-width: 0; min-height: 0; grid-template-columns: 1.9rem minmax(0, 1fr); align-items: center; gap: 0.5rem; padding: 0.28rem 0.58rem 0.28rem 0.28rem; border: 1px solid #596253; border-radius: 999px; color: #d0d5c9; background: rgb(25 29 24 / 0.78); text-align: left; cursor: pointer; transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease), opacity var(--courier-duration) var(--courier-ease); }
  .endpoint-terminal { display: inline-grid; width: 1.85rem; height: 1.85rem; place-items: center; border: 1px solid rgb(185 192 177 / 0.42); border-radius: 50%; color: #d0d5c9; background: rgb(9 12 9 / 0.76); }
  .endpoint courier-icon { width: 0.98rem; height: 0.98rem; }
  .endpoint-name { overflow: hidden; font-size: 0.77rem; font-weight: 760; text-overflow: ellipsis; white-space: nowrap; }
  .endpoint:hover:not(:disabled), .endpoint:focus-visible { border-color: #9ea795; background: rgb(42 47 38 / 0.92); transform: translateY(-1px); }
  .endpoint.selected { border-color: var(--courier-signal); color: var(--courier-paper-50); background: rgb(43 49 35 / 0.92); }
  .endpoint.selected .endpoint-terminal { border-color: var(--courier-signal); color: var(--courier-signal); box-shadow: 0 0 0 0.2rem rgb(212 255 69 / 0.1), 0 0 1rem rgb(212 255 69 / 0.16); }
  .endpoint:disabled { opacity: 0.3; cursor: not-allowed; }
  .route-readout { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(2.8rem, auto) auto auto minmax(0, 1fr); align-content: stretch; gap: 0.58rem; padding: clamp(0.8rem, 1.7vw, 1.2rem); border: 1px solid rgb(137 147 129 / 0.42); border-radius: var(--courier-radius-md); background: rgb(8 11 8 / 0.52); }
  .route-readout code { display: block; max-width: 100%; overflow-wrap: anywhere; color: var(--courier-paper-50); font-size: clamp(0.72rem, 1.15vw, 0.88rem); line-height: 1.45; white-space: normal; }
  .command-shape b { color: var(--courier-signal); }
  .route-meaning { display: grid; grid-template-columns: 1fr 1fr; gap: 0.45rem; }
  .route-meaning span { display: grid; min-width: 0; gap: 0.14rem; padding: 0.42rem 0.5rem; border-left: 2px solid var(--courier-signal); color: #cbd0c3; background: rgb(32 36 29 / 0.58); font-size: 0.64rem; line-height: 1.35; }
  .route-meaning b { color: var(--courier-signal); font-family: var(--courier-font-mono); font-size: 0.55rem; letter-spacing: 0.08em; text-transform: uppercase; }
  .flag-list { display: flex; min-height: 0; align-content: flex-start; gap: 0.35rem; overflow-y: auto; flex-wrap: wrap; }
  .flag { height: max-content; padding: 0.28rem 0.42rem; border: 1px solid #4d5547; border-radius: var(--courier-radius-sm); color: #cbd0c3; background: rgb(28 31 27 / 0.82); font-family: var(--courier-font-mono); font-size: 0.66rem; }

  .install-interface { display: grid; height: 100%; min-width: 0; min-height: 0; grid-template-columns: minmax(0, 1.35fr) minmax(19rem, 0.65fr); grid-template-rows: auto auto minmax(0, 1fr); gap: clamp(0.6rem, 1.4vw, 1rem); padding: clamp(0.8rem, 2vw, 1.5rem); }
  .install-label { grid-column: 1 / -1; grid-row: 1; }
  .install-channels { display: flex; min-width: 0; grid-column: 1 / -1; grid-row: 2; align-content: flex-start; gap: 0.42rem; flex-wrap: wrap; }
  .install-channel { appearance: none; display: inline-grid; min-width: 0; min-height: 2.15rem; grid-template-columns: auto auto; align-items: center; gap: 0.42rem; padding: 0.36rem 0.62rem; border: 1px solid #596253; border-radius: 999px; color: #d0d5c9; background: rgb(25 29 24 / 0.76); text-align: left; cursor: pointer; transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease); }
  .install-channel:hover, .install-channel:focus-visible { border-color: #9ea795; background: rgb(42 47 38 / 0.92); transform: translateY(-1px); }
  .install-channel.selected { border-color: var(--courier-signal); color: var(--courier-signal); background: rgb(43 49 35 / 0.96); }
  .install-channel courier-brand-icon { width: 1.1rem; height: 1.1rem; }
  .install-channel span { min-width: 0; overflow: hidden; font-size: 0.72rem; font-weight: 780; text-overflow: ellipsis; white-space: nowrap; }
  .install-readout { display: grid; min-width: 0; min-height: 0; grid-column: 1; grid-row: 3; grid-template-rows: auto minmax(0, 1fr) 1rem; gap: 0.55rem; padding: clamp(0.75rem, 1.5vw, 1rem); border: 1px solid rgb(137 147 129 / 0.42); border-radius: var(--courier-radius-md); background: rgb(8 11 8 / 0.52); }
  .install-readout-head, .install-readout-head > div { display: flex; align-items: center; }
  .install-readout-head { min-width: 0; justify-content: space-between; gap: 0.75rem; color: var(--courier-signal); }
  .install-readout-head > div { gap: 0.55rem; }
  .install-readout-head courier-brand-icon { width: 1.1rem; height: 1.1rem; }
  .copy-command { appearance: none; display: inline-flex; min-width: 0; min-height: 2rem; align-items: center; gap: 0.4rem; padding: 0.32rem 0.52rem; border: 1px solid #596253; border-radius: 999px; color: #d0d5c9; background: rgb(31 35 29 / 0.78); cursor: pointer; }
  .copy-command:hover { border-color: var(--courier-signal); color: var(--courier-signal); }
  .copy-command courier-icon { width: 0.9rem; height: 0.9rem; }
  .copy-command span { font-size: 0.66rem; font-weight: 760; }
  .copy-status { min-height: 1rem; color: #b9c0b1; font-family: var(--courier-font-mono); font-size: 0.58rem; }
  .install-readout pre { width: 100%; min-width: 0; min-height: 0; overflow: auto; align-self: stretch; }
  .install-readout code { display: block; min-width: 0; overflow-wrap: anywhere; font-size: clamp(0.69rem, 1vw, 0.82rem); line-height: 1.5; white-space: pre-wrap; }
  .install-actions { min-width: 0; min-height: 0; grid-column: 2; grid-row: 3; align-items: stretch; gap: 0.55rem; flex-direction: column; }
  .download-channel { display: grid; min-width: 0; min-height: 0; flex: 1; grid-template-columns: auto minmax(0, 1fr); align-items: center; gap: 0.6rem; padding: 0.65rem; border: 1px solid rgb(137 147 129 / 0.42); border-radius: var(--courier-radius-md); background: rgb(24 28 23 / 0.58); text-decoration: none; }
  .download-channel:hover { border-color: var(--courier-signal); }
  .download-channel > courier-icon { width: 1.25rem; height: 1.25rem; color: var(--courier-signal); }
  .download-channel div { display: grid; min-width: 0; gap: 0.18rem; }
  .download-channel h3 { font-size: 0.82rem; }
  .download-channel p { overflow: hidden; color: #b9c0b1; font-size: 0.66rem; line-height: 1.3; text-overflow: ellipsis; white-space: nowrap; }
  .brand-cloud { grid-column: 1 / -1; gap: 0.25rem; color: var(--courier-signal); }
  .brand-cloud courier-brand-icon { width: 0.82rem; height: 0.82rem; }

  .reference { display: grid; grid-template-columns: minmax(15rem, 0.72fr) minmax(0, 1.28fr); }
  .reference-column { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(0, 1fr); }
  .reference-column + .reference-column { border-left: 1px solid #50584b; }
  .reference-head { display: flex; min-height: 3.3rem; align-items: center; justify-content: space-between; gap: 0.75rem; padding: 0.6rem 0.9rem; border-bottom: 1px solid #50584b; }
  .reference-list { min-width: 0; min-height: 0; overflow-y: auto; overscroll-behavior: contain; }
  .command-row, .option-row { display: grid; width: 100%; min-width: 0; gap: 0.4rem; padding: 0.72rem 0.9rem; border: 0; border-bottom: 1px solid #40473b; color: inherit; background: transparent; text-align: left; transition: color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease); }
  .command-row { appearance: none; cursor: pointer; }
  .command-row:hover, .command-row[aria-pressed="true"], .option-row:hover { color: var(--courier-signal); background: rgb(40 45 35 / 0.76); }
  .command-row:last-child, .option-row:last-child { border-bottom: 0; }
  .command-row code, .option-row code { overflow-wrap: anywhere; font-size: 0.74rem; }
  .option-head { display: flex; align-items: baseline; justify-content: space-between; gap: 0.75rem; }
  .option-meta { color: #b9c0b1; font-family: var(--courier-font-mono); font-size: 0.62rem; line-height: 1.4; }
  .empty-state { display: grid; min-height: 100%; place-items: center; padding: 1.5rem; color: #b9c0b1; font-family: var(--courier-font-mono); font-size: 0.72rem; text-align: center; }

  @media (max-width: 52rem) {
    .route-interface { grid-template-columns: minmax(0, 1.05fr) minmax(15rem, 0.95fr); }
    .install-interface { grid-template-columns: minmax(0, 1.2fr) minmax(14rem, 0.8fr); }
    .reference { grid-template-columns: minmax(12rem, 0.68fr) minmax(0, 1.32fr); }
  }
  @media (max-width: 44rem) {
    .shell { width: min(100% - 1rem, 90rem); }
    .masthead { display: grid; min-height: 0; grid-template-columns: minmax(0, 1fr) auto; align-items: center; gap: 0.25rem 0.5rem; padding: 0.35rem 0 0.4rem; }
    .header-left { display: contents; }
    .brand-link { grid-column: 1; grid-row: 1; }
    nav { grid-column: 1 / -1; grid-row: 2; justify-content: flex-start; gap: 1rem; }
    .header-actions { grid-column: 2; grid-row: 1; gap: 0.25rem; }
    .preferences { gap: 0.2rem; }
    .github-link { width: 2.2rem; height: 2.2rem; }
    .slide { padding-top: calc(var(--masthead-height) + 0.45rem); padding-bottom: 0.5rem; }
    .slide-shell { gap: 0.5rem; }
    .hero { align-items: start; padding-top: 1rem; }
    .hero-copy { width: 94%; max-width: 94%; gap: 0.62rem; }
    h1 { max-width: 22rem; overflow-wrap: normal; font-size: clamp(2.55rem, 13vw, 3.6rem); hyphens: none; word-break: normal; }
    h2 { font-size: clamp(1.55rem, 7.3vw, 2.2rem); }
    .tagline { font-size: 1.04rem; }
    .subline, .section-heading .intro { display: none; }
    .section-heading > div { gap: 0.18rem; }
    .hero-shade { background: linear-gradient(180deg, color-mix(in srgb, var(--courier-color-canvas) 94%, transparent) 0 28%, color-mix(in srgb, var(--courier-color-canvas) 56%, transparent) 40%, transparent 61% 88%, color-mix(in srgb, var(--courier-color-canvas) 22%, transparent)); }
    .hero-node small { font-size: 0.48rem; }
    .hero-node strong { font-size: 0.62rem; }
    .hero-node.source { left: 54%; top: 70%; margin: 1.2rem 0 0 -0.75rem; }
    .hero-node.destination { right: 12%; top: 52%; margin-top: -2rem; }
    .hero-route.desktop-route { display: none; }
    .hero-route.mobile-route { display: block; }
    .route-interface { grid-template-columns: minmax(0, 1fr); grid-template-rows: 10.8rem minmax(0, 1fr); gap: 0.45rem; padding: 0.52rem; }
    .route-controls { gap: 1rem; }
    .endpoint-group { grid-template-rows: auto repeat(4, 2rem); gap: 0.2rem; }
    .endpoint { grid-template-columns: 1.42rem minmax(0, 1fr); gap: 0.28rem; padding: 0.2rem 0.35rem 0.2rem 0.2rem; }
    .endpoint-terminal { width: 1.4rem; height: 1.4rem; }
    .endpoint courier-icon { width: 0.76rem; height: 0.76rem; }
    .endpoint-name { font-size: 0.59rem; }
    .route-readout { grid-template-rows: auto minmax(2rem, auto) auto auto minmax(0, 1fr); gap: 0.3rem; padding: 0.46rem; }
    .route-readout code { font-size: 0.65rem; }
    .route-meaning { gap: 0.28rem; }
    .route-meaning span { padding: 0.28rem 0.34rem; font-size: 0.53rem; }
    .route-meaning b { font-size: 0.48rem; }
    .flag { padding: 0.2rem 0.3rem; font-size: 0.56rem; }
    .install-interface { grid-template-columns: minmax(0, 1fr); grid-template-rows: auto auto minmax(4.8rem, 1fr) auto; gap: 0.4rem; padding: 0.52rem; }
    .install-label { grid-column: 1; grid-row: 1; }
    .install-channels { grid-column: 1; grid-row: 2; gap: 0.26rem; }
    .install-channel { min-height: 1.72rem; gap: 0.28rem; padding: 0.22rem 0.42rem; }
    .install-channel courier-brand-icon { width: 0.78rem; height: 0.78rem; }
    .install-channel span { font-size: 0.57rem; }
    .install-readout { grid-column: 1; grid-row: 3; grid-template-rows: auto minmax(0, 1fr) 0.72rem; gap: 0.28rem; padding: 0.44rem; }
    .install-readout code { font-size: 0.6rem; line-height: 1.35; }
    .copy-command { min-height: 1.65rem; padding: 0.22rem 0.38rem; }
    .copy-command span { display: none; }
    .copy-status { min-height: 0.72rem; font-size: 0.5rem; }
    .install-actions { grid-column: 1; grid-row: 4; gap: 0.35rem; flex-direction: row; }
    .download-channel { padding: 0.42rem; }
    .download-channel p { display: none; }
    .download-channel h3 { font-size: 0.7rem; }
    .brand-cloud courier-brand-icon { width: 0.66rem; height: 0.66rem; }
    .reference { grid-template-columns: minmax(7.8rem, 0.78fr) minmax(0, 1.22fr); }
    .reference-head { min-height: 3rem; align-items: flex-start; padding: 0.45rem 0.5rem; }
    .reference-column.options .reference-head { display: grid; gap: 0.25rem; }
    .reference-head courier-checkbox { font-size: 0.58rem; }
    .command-row, .option-row { padding: 0.55rem 0.5rem; }
    .command-row code, .option-row code { font-size: 0.62rem; }
    .option-head { display: grid; gap: 0.18rem; }
    .option-meta { font-size: 0.52rem; }
    courier-brand::part(product) { display: none; }
  }
  @media (max-width: 30rem) {
    courier-brand::part(words) { display: none; }
    nav { gap: 0.75rem; }
    nav > a { font-size: 0.68rem; }
  }
  @media (max-height: 44rem) and (min-width: 44.01rem) {
    .masthead { min-height: 4rem; }
    .slide { padding-top: calc(var(--masthead-height) + 0.4rem); padding-bottom: 0.45rem; }
    .slide-shell { gap: 0.45rem; }
    .section-heading .intro, .subline { display: none; }
    h2 { font-size: 2rem; }
    .route-interface, .install-interface { padding: 0.55rem; }
  }
  @media (prefers-reduced-motion: reduce) {
    .hero-route .signal { animation: none; stroke-dashoffset: 0; }
  }
`;
