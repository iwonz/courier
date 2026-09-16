import { css } from "lit";

export const landingStyles = css`
  :host { --masthead-height: 5rem; display: block; width: 100%; min-height: 100svh; overflow-x: clip; color: var(--courier-color-text); background: var(--courier-color-canvas); font-family: var(--courier-font-sans); }
  * { box-sizing: border-box; }
  a { color: inherit; }
  button { font: inherit; }
  code { font-family: var(--courier-font-mono); }
  h1, h2, h3, p { margin: 0; }
  .shell { width: min(92rem, calc(100% - clamp(1.25rem, 5vw, 5rem))); margin: 0 auto; }

  .masthead-wrap { position: fixed; z-index: 100; top: 0; right: 0; left: 0; padding-top: env(safe-area-inset-top, 0px); padding-right: env(safe-area-inset-right, 0px); padding-left: env(safe-area-inset-left, 0px); isolation: isolate; pointer-events: none; }
  .masthead-wrap::before { content: ""; position: absolute; z-index: -1; inset: 0 0 -2.7rem; background: linear-gradient(180deg, rgb(12 15 12 / 0.58) 0%, rgb(12 15 12 / 0.34) 38%, rgb(12 15 12 / 0.12) 68%, transparent 100%); backdrop-filter: blur(15px) saturate(0.8); mask-image: linear-gradient(180deg, #000 0 54%, rgb(0 0 0 / 0.72) 70%, transparent 100%); pointer-events: none; }
  .masthead { display: flex; min-height: 4.75rem; align-items: center; gap: 1.5rem; color: var(--courier-paper-50); pointer-events: auto; }
  .header-left, nav, .header-actions, .preferences, .section-heading, .install-actions { display: flex; align-items: center; }
  .header-left { min-width: 0; gap: clamp(1rem, 3vw, 2.5rem); }
  .brand-link, nav a, .github-link, .download-channel { text-decoration: none; }
  .brand-link courier-brand { --courier-color-text: var(--courier-paper-50); --courier-color-muted: #b9c0b1; }
  nav { gap: 1.45rem; }
  nav a { color: #cbd0c3; font-size: 0.76rem; font-weight: 760; }
  nav a:hover, nav a:focus-visible, nav a[aria-current="page"] { color: var(--courier-signal); }
  .header-actions { margin-left: auto; gap: 0.5rem; }
  .preferences { gap: 0.35rem; }
  .github-link { display: inline-grid; width: 2.35rem; height: 2.35rem; place-items: center; border: 1px solid rgb(243 244 233 / 0.24); border-radius: var(--courier-radius-sm); color: var(--courier-paper-50); background: rgb(12 15 12 / 0.18); }
  .github-link:hover { border-color: var(--courier-signal); color: var(--courier-signal); }
  .github-link courier-icon { width: 1.05rem; height: 1.05rem; }

  main { position: relative; }
  .slide { position: relative; height: 100svh; min-height: 32rem; overflow: hidden; padding-top: calc(var(--masthead-height) + 0.65rem); padding-bottom: clamp(0.7rem, 2vh, 1.3rem); scroll-snap-align: start; scroll-snap-stop: always; isolation: isolate; }
  courier-scene { z-index: 0; }
  .slide::after { content: ""; position: absolute; z-index: 1; inset: 0; background: linear-gradient(180deg, rgb(9 12 9 / 0.04), rgb(9 12 9 / 0.16)); pointer-events: none; }
  .slide-shell { position: relative; z-index: 3; display: grid; height: 100%; min-height: 0; grid-template-rows: auto minmax(0, 1fr); gap: clamp(0.55rem, 1.8vh, 1rem); }
  .section-heading { min-width: 0; align-items: end; justify-content: space-between; color: var(--courier-paper-50); }
  .section-heading > div { display: grid; gap: 0.28rem; }
  .section-index, .label { color: var(--courier-signal); font-family: var(--courier-font-mono); font-size: 0.62rem; font-weight: 780; letter-spacing: 0.11em; text-transform: uppercase; }
  h1, h2 { font-family: var(--courier-font-display); font-weight: 860; letter-spacing: -0.065em; line-height: 0.92; }
  h1 { max-width: 48rem; font-size: clamp(4rem, 8.2vw, 8.2rem); }
  h2 { font-size: clamp(2rem, 4.6vw, 4.4rem); }
  .intro { max-width: 48rem; color: #d7dbd1; font-size: clamp(0.76rem, 1.25vw, 0.95rem); line-height: 1.45; }

  .hero-slide { padding: 0; }
  .hero-slide::after { background: none; }
  .hero { position: relative; z-index: 4; display: grid; height: 100%; min-height: 0; align-items: center; pointer-events: none; }
  .hero-copy { display: grid; width: min(62rem, 72%); gap: 1rem; color: var(--courier-graphite-950); }
  .tagline { max-width: 38rem; font-size: clamp(1rem, 1.75vw, 1.45rem); font-weight: 760; line-height: 1.25; }
  .hero-shade { position: absolute; z-index: 1; inset: 0; background: linear-gradient(90deg, rgb(243 244 233 / 0.92) 0 35%, rgb(243 244 233 / 0.55) 52%, transparent 72%), linear-gradient(180deg, rgb(243 244 233 / 0.22), transparent 32% 74%, rgb(16 18 15 / 0.1)); pointer-events: none; }
  .hero-route { position: absolute; z-index: 3; inset: 0; width: 100%; height: 100%; pointer-events: none; }
  .hero-route.mobile-route { display: none; }
  .hero-route path { fill: none; stroke: var(--courier-signal); stroke-width: 0.2; vector-effect: non-scaling-stroke; }
  .hero-route .signal { stroke-width: 0.42; stroke-dasharray: 2 7; animation: route-signal 3s linear infinite; }
  .hero-terminal { position: absolute; z-index: 4; width: clamp(0.82rem, 1.15vw, 1.15rem); height: clamp(0.82rem, 1.15vw, 1.15rem); aspect-ratio: 1; border: 0.2rem solid var(--courier-signal); border-radius: 50%; background: rgb(12 15 12 / 0.88); box-shadow: 0 0 0 0.28rem rgb(212 255 69 / 0.14), 0 0 1.3rem rgb(212 255 69 / 0.2); transform: translate(-50%, -50%); pointer-events: none; }
  .hero-terminal::after { content: ""; position: absolute; inset: 28%; border-radius: 50%; background: var(--courier-signal); }
  .hero-terminal.source { left: 54%; top: 70%; }
  .hero-terminal.destination { left: 91%; top: 38%; }
  .hero-node { position: absolute; z-index: 4; display: grid; gap: 0.12rem; color: var(--courier-paper-50); font-family: var(--courier-font-mono); text-transform: uppercase; pointer-events: none; }
  .hero-node small { color: var(--courier-signal); font-size: 0.55rem; letter-spacing: 0.12em; }
  .hero-node strong { font-size: 0.72rem; letter-spacing: 0.08em; }
  .hero-node.source { left: 54%; top: 70%; margin: 1.35rem 0 0 -1.25rem; }
  .hero-node.destination { right: 7%; top: 38%; margin-top: -2.25rem; text-align: right; }
  @keyframes route-signal { to { stroke-dashoffset: -18; } }

  .route-explorer, .install-board, .cli-workspace { width: 100%; height: 100%; min-width: 0; min-height: 0; }
  .route-interface { display: grid; height: 100%; min-height: 0; grid-template-rows: auto minmax(0, 1fr); gap: clamp(0.55rem, 1.5vh, 0.85rem); }
  .route-controls { position: relative; display: grid; min-width: 0; grid-template-columns: 1fr 1fr; gap: clamp(2.5rem, 7vw, 8rem); padding: 0.35rem 0; }
  .route-connector { position: absolute; z-index: 0; inset: 0; width: 100%; height: 100%; overflow: visible; pointer-events: none; }
  .route-connector path { fill: none; stroke: var(--courier-signal); stroke-width: 1.5; stroke-dasharray: 4 6; vector-effect: non-scaling-stroke; }
  .endpoint-group { position: relative; z-index: 1; display: grid; min-width: 0; grid-template-columns: auto repeat(4, minmax(0, 1fr)); align-items: center; gap: 0.42rem; }
  .endpoint { appearance: none; display: grid; min-width: 0; min-height: 2.45rem; grid-template-columns: 1.55rem minmax(0, 1fr); align-items: center; gap: 0.42rem; padding: 0.28rem 0.58rem 0.28rem 0.32rem; border: 1px solid rgb(203 208 195 / 0.42); border-radius: 999px; color: #d7dbd1; background: rgb(12 15 12 / 0.28); text-align: left; cursor: pointer; backdrop-filter: blur(10px); transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease); }
  .endpoint-terminal { display: inline-grid; width: 1.55rem; height: 1.55rem; aspect-ratio: 1; place-items: center; border: 1px solid rgb(203 208 195 / 0.45); border-radius: 50%; color: #d0d5c9; background: rgb(9 12 9 / 0.64); }
  .endpoint courier-icon { width: 0.82rem; height: 0.82rem; }
  .endpoint-name { overflow: hidden; font-size: 0.67rem; font-weight: 780; text-overflow: ellipsis; white-space: nowrap; }
  .endpoint:hover:not(:disabled), .endpoint:focus-visible { border-color: #d7dbd1; background: rgb(28 33 27 / 0.66); transform: translateY(-1px); }
  .endpoint.selected { border-color: var(--courier-signal); color: var(--courier-paper-50); background: rgb(31 36 28 / 0.68); }
  .endpoint.selected .endpoint-terminal { border-color: var(--courier-signal); color: var(--courier-signal); box-shadow: 0 0 0 0.18rem rgb(212 255 69 / 0.1); }
  .endpoint:disabled { opacity: 0.28; cursor: not-allowed; }
  .endpoint:focus-visible, .install-channel:focus-visible, .command-row:focus-visible { outline: 3px solid var(--courier-beak); outline-offset: 2px; }
  .route-readout { min-height: 0; }
  .demo-details { display: grid; gap: 0.38rem; padding: 0.5rem 0.8rem; border-top: 1px solid var(--courier-terminal-border); }
  .demo-details .label { color: var(--courier-terminal-prompt); }
  .flag-list { display: flex; min-height: 0; align-content: flex-start; gap: 0.3rem; overflow: auto; flex-wrap: wrap; }
  .flag { height: max-content; padding: 0.22rem 0.38rem; border: 1px solid var(--courier-terminal-border); border-radius: 999px; color: var(--courier-terminal-muted); background: color-mix(in srgb, var(--courier-terminal-text) 6%, transparent); font-family: var(--courier-font-mono); font-size: 0.58rem; }

  .install-interface { display: grid; height: 100%; min-height: 0; grid-template-rows: auto auto minmax(0, 1fr); gap: clamp(0.5rem, 1.4vh, 0.8rem); }
  .install-channels { display: flex; min-width: 0; align-content: flex-start; gap: 0.38rem; flex-wrap: wrap; }
  .install-channel { appearance: none; display: inline-grid; min-height: 2rem; grid-template-columns: auto auto; align-items: center; gap: 0.38rem; padding: 0.3rem 0.58rem; border: 1px solid rgb(203 208 195 / 0.42); border-radius: 999px; color: #d7dbd1; background: rgb(12 15 12 / 0.28); cursor: pointer; backdrop-filter: blur(10px); }
  .install-channel:hover, .install-channel:focus-visible { border-color: #d7dbd1; }
  .install-channel.selected { border-color: var(--courier-signal); color: var(--courier-signal); background: rgb(31 36 28 / 0.68); }
  .install-channel courier-brand-icon { width: 1rem; height: 1rem; }
  .install-channel span { font-size: 0.68rem; font-weight: 780; }
  .install-readout { min-height: 0; }
  .install-identity { display: flex; align-items: center; gap: 0.5rem; padding: 0.45rem 0.8rem; border-top: 1px solid var(--courier-terminal-border); color: var(--courier-terminal-prompt); font-family: var(--courier-font-mono); font-size: 0.68rem; }
  .install-identity courier-brand-icon { width: 1rem; height: 1rem; }
  .install-actions { min-width: 0; gap: 0.4rem; }
  .download-channel { display: inline-flex; min-width: 0; align-items: center; gap: 0.42rem; color: var(--courier-terminal-text); }
  .download-channel:hover { color: var(--courier-terminal-prompt); }
  .download-channel > courier-icon { width: 1rem; height: 1rem; color: var(--courier-terminal-prompt); }
  .download-channel div { display: flex; min-width: 0; align-items: center; gap: 0.35rem; }
  .download-channel h3 { font-size: 0.62rem; }
  .download-channel p { display: none; }
  .brand-cloud { display: inline-flex; gap: 0.18rem; color: var(--courier-terminal-prompt); }
  .brand-cloud courier-brand-icon { width: 0.66rem; height: 0.66rem; }

  .cli-workspace { display: grid; grid-template-columns: minmax(0, 1.22fr) minmax(20rem, 0.78fr); gap: 0.75rem; }
  .reference { display: grid; min-height: 0; grid-template-columns: minmax(13rem, 0.72fr) minmax(0, 1.28fr); overflow: hidden; border: 1px solid var(--courier-terminal-border); border-radius: var(--courier-radius-md); color: var(--courier-terminal-text); background: var(--courier-terminal-surface); backdrop-filter: blur(18px) saturate(0.8); }
  .reference-column { display: grid; min-width: 0; min-height: 0; grid-template-rows: auto minmax(0, 1fr); }
  .reference-column + .reference-column { border-left: 1px solid var(--courier-terminal-border); }
  .reference-head { display: flex; min-height: 2.7rem; align-items: center; justify-content: space-between; gap: 0.75rem; padding: 0.52rem 0.72rem; border-bottom: 1px solid var(--courier-terminal-border); }
  .reference-list { min-width: 0; min-height: 0; overflow-y: auto; overscroll-behavior: contain; }
  .command-row, .option-row { display: grid; width: 100%; min-width: 0; gap: 0.35rem; padding: 0.58rem 0.72rem; border: 0; border-bottom: 1px solid color-mix(in srgb, var(--courier-terminal-border) 65%, transparent); color: inherit; background: transparent; text-align: left; }
  .command-row { appearance: none; cursor: pointer; }
  .command-row:hover, .command-row[aria-pressed="true"], .option-row:hover { color: var(--courier-terminal-prompt); background: rgb(212 255 69 / 0.055); }
  .command-row code, .option-row code { overflow-wrap: anywhere; font-size: 0.67rem; }
  .option-head { display: flex; align-items: baseline; justify-content: space-between; gap: 0.65rem; }
  .option-meta { color: var(--courier-terminal-muted); font-family: var(--courier-font-mono); font-size: 0.55rem; line-height: 1.35; }
  .empty-state { display: grid; min-height: 100%; place-items: center; padding: 1rem; color: var(--courier-terminal-muted); font-family: var(--courier-font-mono); font-size: 0.65rem; text-align: center; }
  .cli-demo { min-height: 0; }

  @media (max-width: 70rem) {
    .endpoint-group { grid-template-columns: auto repeat(2, minmax(0, 1fr)); }
    .endpoint-group .label { grid-row: 1 / 3; }
    .cli-workspace { grid-template-columns: minmax(0, 1.1fr) minmax(17rem, 0.9fr); }
  }
  @media (max-width: 44rem) {
    .shell { width: min(100% - 1rem, 92rem); }
    .masthead { display: grid; min-height: 0; grid-template-columns: minmax(0, 1fr) auto; gap: 0.25rem 0.5rem; padding: 0.35rem 0 0.4rem; }
    .header-left { display: contents; }
    .brand-link { grid-column: 1; grid-row: 1; }
    nav { grid-column: 1 / -1; grid-row: 2; justify-content: flex-start; gap: 1rem; }
    .header-actions { grid-column: 2; grid-row: 1; gap: 0.25rem; }
    .preferences { gap: 0.2rem; }
    .github-link { width: 2.2rem; height: 2.2rem; }
    .slide { padding-top: calc(var(--masthead-height) + 0.4rem); padding-bottom: 0.45rem; }
    .slide-shell { gap: 0.4rem; }
    .hero { align-items: start; padding-top: calc(var(--masthead-height) + 1.2rem); }
    .hero-copy { width: 96%; gap: 0.55rem; }
    h1 { max-width: 22rem; font-size: clamp(2.6rem, 13vw, 3.65rem); }
    h2 { font-size: clamp(1.6rem, 7.2vw, 2.15rem); }
    .tagline { max-width: 22rem; font-size: 1rem; }
    .intro { font-size: 0.62rem; }
    .hero-shade { background: linear-gradient(180deg, rgb(243 244 233 / 0.94) 0 30%, rgb(243 244 233 / 0.58) 43%, transparent 63% 88%, rgb(16 18 15 / 0.12)); }
    .hero-route.desktop-route { display: none; }
    .hero-route.mobile-route { display: block; }
    .hero-terminal.source { left: 54%; top: 70%; }
    .hero-terminal.destination { left: 88%; top: 52%; }
    .hero-node small { font-size: 0.46rem; }
    .hero-node strong { font-size: 0.58rem; }
    .hero-node.source { left: 54%; top: 70%; margin: 1.2rem 0 0 -0.75rem; }
    .hero-node.destination { right: 12%; top: 52%; margin-top: -2rem; }
    .route-interface { grid-template-rows: 8.6rem minmax(0, 1fr); gap: 0.36rem; }
    .route-controls { gap: 0.65rem; padding: 0; }
    .endpoint-group { grid-template-columns: minmax(0, 1fr); grid-template-rows: auto repeat(4, 1.65rem); gap: 0.16rem; }
    .endpoint-group .label { grid-row: auto; }
    .endpoint { min-height: 1.65rem; grid-template-columns: 1.25rem minmax(0, 1fr); gap: 0.25rem; padding: 0.12rem 0.28rem 0.12rem 0.12rem; }
    .endpoint-terminal { width: 1.25rem; height: 1.25rem; }
    .endpoint courier-icon { width: 0.68rem; height: 0.68rem; }
    .endpoint-name { font-size: 0.55rem; }
    .demo-details { gap: 0.25rem; padding: 0.34rem 0.55rem; }
    .flag { padding: 0.16rem 0.26rem; font-size: 0.5rem; }
    .install-interface { gap: 0.34rem; }
    .install-channels { gap: 0.24rem; }
    .install-channel { min-height: 1.65rem; gap: 0.25rem; padding: 0.18rem 0.38rem; }
    .install-channel courier-brand-icon { width: 0.75rem; height: 0.75rem; }
    .install-channel span { font-size: 0.55rem; }
    .install-identity { padding: 0.3rem 0.55rem; font-size: 0.56rem; }
    .install-actions { gap: 0.32rem; }
    .download-channel h3 { display: none; }
    .brand-cloud courier-brand-icon { width: 0.58rem; height: 0.58rem; }
    .cli-workspace { grid-template-columns: minmax(0, 1fr); grid-template-rows: minmax(9rem, 0.82fr) minmax(0, 1.18fr); gap: 0.38rem; }
    .reference { grid-template-columns: minmax(7.5rem, 0.78fr) minmax(0, 1.22fr); }
    .reference-head { min-height: 2.5rem; align-items: flex-start; padding: 0.38rem 0.44rem; }
    .reference-column.options .reference-head { display: grid; gap: 0.18rem; }
    .reference-head courier-checkbox { font-size: 0.52rem; }
    .command-row, .option-row { padding: 0.42rem 0.44rem; }
    .command-row code, .option-row code { font-size: 0.56rem; }
    .option-head { display: grid; gap: 0.14rem; }
    .option-meta { font-size: 0.48rem; }
    courier-brand::part(product) { display: none; }
  }
  @media (max-width: 30rem) {
    courier-brand::part(words) { display: none; }
    nav { gap: 0.72rem; }
    nav > a { font-size: 0.66rem; }
  }
  @media (max-height: 44rem) and (min-width: 44.01rem) {
    .masthead { min-height: 4rem; }
    .slide { padding-top: calc(var(--masthead-height) + 0.35rem); padding-bottom: 0.4rem; }
    .slide-shell { gap: 0.35rem; }
    h2 { font-size: 2rem; }
    .route-interface { grid-template-rows: auto minmax(0, 1fr); }
  }
  @media (prefers-reduced-motion: reduce) {
    .hero-route .signal { animation: none; stroke-dashoffset: 0; }
  }
`;
