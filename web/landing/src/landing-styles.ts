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

  .shell { width: min(100% - clamp(1.25rem, 5vw, 5rem), 88rem); margin-inline: auto; }
  h1, h2, p { margin: 0; }
  h1, h2 { font-family: var(--courier-font-display); text-wrap: balance; }
  h1 { max-width: 8.8ch; font-size: clamp(3.7rem, 7.8vw, 7.3rem); font-weight: 860; letter-spacing: -0.072em; line-height: 0.88; }
  h2 { font-size: clamp(2.4rem, 5.2vw, 4.8rem); font-weight: 840; letter-spacing: -0.062em; line-height: 0.94; }
  a:focus-visible, button:focus-visible { outline: 3px solid var(--courier-color-accent); outline-offset: 3px; }

  .masthead-wrap {
    position: fixed;
    z-index: 30;
    top: 0;
    right: 0;
    left: 0;
    padding-top: env(safe-area-inset-top);
    color: var(--courier-color-text);
    background: color-mix(in srgb, var(--courier-color-canvas) 84%, transparent);
    box-shadow: 0 0.7rem 2rem color-mix(in srgb, var(--courier-ink-950) 7%, transparent);
    -webkit-backdrop-filter: blur(20px) saturate(1.2);
    backdrop-filter: blur(20px) saturate(1.2);
  }
  .masthead { display: grid; min-height: 4.6rem; grid-template-columns: auto 1fr auto; align-items: center; gap: 1.5rem; }
  .brand-link { color: inherit; text-decoration: none; }
  nav, .header-actions, .install-actions { display: flex; align-items: center; }
  nav { justify-content: flex-start; gap: 1.25rem; }
  nav a { position: relative; color: var(--courier-color-muted); font-size: 0.78rem; font-weight: 780; text-decoration: none; }
  nav a::after { content: ""; position: absolute; right: 0; bottom: -0.42rem; left: 0; height: 0.16rem; border-radius: 999px; background: var(--courier-color-accent); transform: scaleX(0); transition: transform var(--courier-duration) var(--courier-ease); }
  nav a:hover, nav a[aria-current="page"] { color: var(--courier-color-text); }
  nav a[aria-current="page"]::after { transform: scaleX(1); }
  .header-actions { gap: 0.38rem; }

  main { overflow: clip; }
  .stage { position: relative; scroll-margin-top: var(--masthead-height); }
  .stage-shell { position: relative; display: grid; gap: clamp(1.5rem, 3vw, 2.5rem); }

  .route-stage {
    padding: calc(var(--masthead-height) + clamp(2rem, 5vw, 4rem)) 0 clamp(4rem, 7vw, 7rem);
    background:
      radial-gradient(circle at 7% 12%, color-mix(in srgb, var(--courier-mint) 34%, transparent) 0 10rem, transparent 24rem),
      radial-gradient(circle at 92% 32%, color-mix(in srgb, var(--courier-sky) 60%, transparent) 0 13rem, transparent 29rem),
      linear-gradient(155deg, var(--courier-color-canvas) 0 58%, color-mix(in srgb, var(--courier-sky) 34%, var(--courier-color-canvas)) 100%);
  }
  .route-shell { grid-template-columns: minmax(20rem, 0.78fr) minmax(26rem, 1.22fr); grid-template-areas: "copy art" "instrument instrument"; align-items: center; }
  .hero-copy { position: relative; z-index: 2; display: grid; grid-area: copy; align-content: center; min-width: 0; min-height: 21rem; }
  .hero-copy h1 { color: var(--courier-color-text); }
  .hero-copy h1::first-line { color: var(--courier-brand); }
  .hero-orbit { position: absolute; z-index: -1; top: 50%; left: -4rem; width: min(33rem, 55vw); height: min(33rem, 55vw); border: 1px solid color-mix(in srgb, var(--courier-brand) 15%, transparent); border-radius: 50%; transform: translateY(-50%); pointer-events: none; }
  .hero-orbit::before, .hero-orbit::after { content: ""; position: absolute; border: 1px solid color-mix(in srgb, var(--courier-brand) 12%, transparent); border-radius: inherit; }
  .hero-orbit::before { inset: 12%; }
  .hero-orbit::after { inset: 26%; }
  .hero-orbit span { position: absolute; width: 0.72rem; height: 0.72rem; border-radius: 50%; background: var(--courier-color-accent-solid); box-shadow: 0 0 0 0.4rem color-mix(in srgb, var(--courier-color-accent-solid) 16%, transparent); }
  .hero-orbit span:nth-child(1) { top: 18%; right: 14%; }
  .hero-orbit span:nth-child(2) { right: 28%; bottom: 4%; background: var(--courier-mint); }
  .hero-orbit span:nth-child(3) { top: 44%; left: 8%; width: 0.4rem; height: 0.4rem; background: var(--courier-brand); }
  .hero-art { position: relative; grid-area: art; overflow: hidden; aspect-ratio: 3 / 2; border: 1px solid color-mix(in srgb, var(--courier-brand) 18%, transparent); border-radius: clamp(2rem, 5vw, 4.5rem); background: var(--courier-cream); box-shadow: 0 2.4rem 6rem color-mix(in srgb, var(--courier-blue-deep) 19%, transparent); transform: rotate(1.2deg); }
  .hero-art::after { content: ""; position: absolute; inset: 0; border-radius: inherit; box-shadow: inset 0 0 0 0.55rem rgb(255 255 255 / 0.22); pointer-events: none; }
  .hero-art courier-mascot, .hero-art courier-mascot::part(image) { width: 100%; height: 100%; }
  .hero-art courier-mascot::part(image) { object-fit: cover; }

  .route-instrument {
    --courier-color-text: var(--courier-ink-950);
    --courier-color-muted: #58677e;
    --courier-color-accent: var(--courier-coral-dark);
    --courier-workbench-surface: rgb(255 255 255 / 0.94);
    --courier-workbench-line: rgb(10 28 55 / 0.14);
    --courier-workbench-shadow: 0 1rem 2.5rem rgb(4 17 42 / 0.16);
    display: grid;
    grid-area: instrument;
    gap: 1rem;
    padding: clamp(0.8rem, 2.5vw, 1.5rem);
    border-radius: clamp(1.5rem, 3vw, 2.4rem);
    color: var(--courier-cream);
    background:
      radial-gradient(circle at 88% 8%, rgb(116 171 255 / 0.42), transparent 27%),
      linear-gradient(135deg, var(--courier-blue-deep), var(--courier-blue));
    box-shadow: 0 1.5rem 4rem rgb(12 42 94 / 0.18);
  }
  .route-controls { position: relative; display: grid; grid-template-columns: 1fr 1fr; gap: clamp(2rem, 8vw, 9rem); }
  .route-connector { position: absolute; z-index: 0; inset: 0; width: 100%; height: 100%; overflow: visible; pointer-events: none; }
  .route-connector path { fill: none; vector-effect: non-scaling-stroke; }
  .route-path { stroke: rgb(255 255 255 / 0.46); stroke-width: 1.5; }
  .route-signal { stroke: var(--courier-coral); stroke-width: 3; stroke-dasharray: 1 10; stroke-linecap: round; animation: route-pulse 2.5s linear infinite; }
  @keyframes route-pulse { to { stroke-dashoffset: -22; } }
  .endpoint-group { position: relative; z-index: 1; display: grid; min-width: 0; grid-template-columns: auto repeat(4, minmax(0, 1fr)); align-items: center; gap: 0.35rem; }
  .endpoint-group > strong { color: rgb(255 255 255 / 0.72); font-size: 0.7rem; font-weight: 800; }
  .endpoint {
    appearance: none;
    display: inline-grid;
    min-width: 0;
    min-height: 2.45rem;
    grid-template-columns: 1.4rem minmax(0, 1fr);
    align-items: center;
    gap: 0.38rem;
    padding: 0.28rem 0.6rem 0.28rem 0.3rem;
    border: 1px solid rgb(255 255 255 / 0.32);
    border-radius: 999px;
    color: white;
    background: rgb(255 255 255 / 0.1);
    text-align: left;
    cursor: pointer;
    transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease), box-shadow var(--courier-duration) var(--courier-ease);
  }
  .endpoint-terminal { display: inline-grid; width: 1.4rem; height: 1.4rem; aspect-ratio: 1; place-items: center; border: 1px solid currentColor; border-radius: 50%; background: rgb(255 255 255 / 0.08); }
  .endpoint courier-icon { width: 0.76rem; height: 0.76rem; }
  .endpoint span:last-child { overflow: hidden; font-size: 0.67rem; font-weight: 780; text-overflow: ellipsis; white-space: nowrap; }
  .endpoint:hover:not(:disabled), .endpoint:focus-visible { border-color: white; background: rgb(255 255 255 / 0.18); transform: translateY(-1px); }
  .endpoint.selected { border-color: var(--courier-coral); color: var(--courier-ink-950); background: var(--courier-coral); box-shadow: 0 0.7rem 1.5rem rgb(7 20 47 / 0.22); }
  .endpoint:disabled { opacity: 0.28; cursor: not-allowed; }
  .route-readout { height: 10.25rem; color: var(--courier-ink-950); }
  .readout-details { display: flex; min-height: 0; align-items: center; gap: 0.6rem; padding: 0.55rem 0.8rem; border-top: 1px solid var(--courier-workbench-line); }
  .flag-list { display: flex; min-width: 0; gap: 0.3rem; overflow: auto; flex-wrap: wrap; }
  .flag { padding: 0.22rem 0.5rem; border: 1px solid var(--courier-workbench-line); border-radius: 999px; color: var(--courier-color-muted); background: color-mix(in srgb, var(--courier-color-canvas) 45%, transparent); font-family: var(--courier-font-mono); font-size: 0.57rem; white-space: nowrap; }

  #install { padding: clamp(4.5rem, 8vw, 7.5rem) 0; color: var(--courier-cream); background: linear-gradient(145deg, var(--courier-blue-deep) 0%, #163e82 66%, #245ec2 100%); }
  #install::before { content: ""; position: absolute; inset: 0; background: radial-gradient(circle at 92% 18%, rgb(113 173 255 / 0.28), transparent 28%), radial-gradient(circle at 5% 90%, rgb(126 229 205 / 0.15), transparent 25%); pointer-events: none; }
  #install .stage-shell { z-index: 1; }
  #install h2 { color: white; }
  .install-interface { display: grid; gap: 1rem; }
  .install-channels { display: flex; gap: 0.45rem; flex-wrap: wrap; }
  .install-channel {
    appearance: none;
    display: inline-flex;
    min-height: 2.2rem;
    align-items: center;
    gap: 0.42rem;
    padding: 0.32rem 0.7rem;
    border: 1px solid rgb(255 255 255 / 0.3);
    border-radius: 999px;
    color: white;
    background: rgb(255 255 255 / 0.08);
    cursor: pointer;
    transition: border-color var(--courier-duration) var(--courier-ease), background var(--courier-duration) var(--courier-ease), transform var(--courier-duration) var(--courier-ease);
  }
  .install-channel:hover, .install-channel:focus-visible { border-color: white; background: rgb(255 255 255 / 0.15); transform: translateY(-1px); }
  .install-channel.selected { border-color: var(--courier-coral); color: var(--courier-ink-950); background: var(--courier-coral); }
  .install-channel courier-brand-icon { width: 1rem; height: 1rem; }
  .install-channel span { font-size: 0.7rem; font-weight: 780; }
  .install-readout { --courier-color-text: var(--courier-ink-950); --courier-color-muted: #58677e; --courier-color-accent: var(--courier-coral-dark); --courier-workbench-surface: rgb(255 255 255 / 0.96); --courier-workbench-line: rgb(10 28 55 / 0.14); --courier-workbench-shadow: 0 1.2rem 3rem rgb(3 15 39 / 0.26); height: 11.25rem; color: var(--courier-ink-950); }
  .install-identity { color: var(--courier-blue); font-size: 0.7rem; }
  .install-identity courier-brand-icon { width: 1rem; height: 1rem; }
  .install-actions { gap: 0.9rem; }
  .install-actions a { display: inline-flex; align-items: center; gap: 0.35rem; color: #4c607e; font-weight: 760; text-decoration: none; }
  .install-actions a:hover { color: var(--courier-blue); }
  .install-actions courier-icon { width: 0.9rem; height: 0.9rem; }

  #cli { padding: clamp(4.5rem, 8vw, 7.5rem) 0 clamp(5rem, 9vw, 8rem); background: linear-gradient(180deg, color-mix(in srgb, var(--courier-sky) 34%, var(--courier-color-canvas)), var(--courier-color-canvas) 46%); }
  .cli-workspace { display: grid; height: 34rem; grid-template-columns: minmax(0, 1.28fr) minmax(20rem, 0.72fr); gap: 1rem; }
  .reference-grid { display: grid; height: 100%; min-height: 0; grid-template-columns: minmax(13rem, 0.7fr) minmax(0, 1.3fr); }
  .reference-column { min-width: 0; min-height: 0; overflow: hidden; }
  .reference-column + .reference-column { border-left: 1px solid var(--courier-workbench-line); }
  .reference-list { height: 100%; min-height: 0; overflow-y: auto; overscroll-behavior: contain; }
  .reference-filter { display: flex; min-height: 2.8rem; align-items: center; justify-content: space-between; gap: 0.7rem; padding: 0.35rem 0.65rem; border-bottom: 1px solid var(--courier-workbench-line); color: var(--courier-color-muted); font-size: 0.68rem; font-weight: 760; }
  .command-row, .option-row { display: grid; width: 100%; min-width: 0; gap: 0.34rem; padding: 0.72rem; border: 0; border-bottom: 1px solid color-mix(in srgb, var(--courier-workbench-line) 65%, transparent); color: inherit; background: transparent; text-align: left; }
  .command-row { appearance: none; cursor: pointer; }
  .command-row:hover, .command-row[aria-pressed="true"] { color: var(--courier-brand); background: color-mix(in srgb, var(--courier-brand) 9%, transparent); }
  .command-row[aria-pressed="true"] { box-shadow: inset 0.22rem 0 var(--courier-color-accent-solid); }
  .command-row code, .option-row code { overflow-wrap: anywhere; font: 700 0.68rem/1.4 var(--courier-font-mono); }
  .option-row > div { display: flex; align-items: baseline; justify-content: space-between; gap: 0.6rem; }
  .option-row span, .option-row p { color: var(--courier-color-muted); font-size: 0.58rem; line-height: 1.4; }
  .empty-state { display: grid; min-height: 100%; place-items: center; padding: 1rem; color: var(--courier-color-muted); font-size: 0.72rem; text-align: center; }
  .cli-readout { min-height: 0; }

  @media (max-width: 70rem) {
    .route-shell { grid-template-columns: minmax(16rem, 0.72fr) minmax(22rem, 1.28fr); }
    .endpoint-group { grid-template-columns: auto repeat(2, minmax(0, 1fr)); }
    .endpoint-group > strong { grid-row: 1 / 3; }
    .cli-workspace { grid-template-columns: minmax(0, 1fr) minmax(17rem, 0.72fr); }
  }

  @media (max-width: 48rem) {
    :host { --masthead-height: 4.2rem; }
    .shell { width: min(100% - 1rem, 88rem); }
    .masthead { min-height: 4.2rem; grid-template-columns: auto 1fr auto; gap: 0.55rem; }
    nav { justify-content: flex-end; gap: 0.7rem; }
    nav a { font-size: 0.66rem; }
    .header-actions { gap: 0.2rem; --courier-control-frame-size: 2.1rem; }
    .route-stage { padding-top: calc(var(--masthead-height) + 1.5rem); padding-bottom: 3.5rem; }
    .route-shell { grid-template-columns: 1fr; grid-template-areas: "copy" "art" "instrument"; gap: 1.35rem; }
    .hero-copy { min-height: auto; padding: 1rem 0 0.5rem; }
    h1 { max-width: 9ch; font-size: clamp(3.3rem, 17vw, 5.4rem); }
    h2 { font-size: clamp(2.5rem, 12vw, 4rem); }
    .hero-orbit { left: auto; right: -8rem; width: 24rem; height: 24rem; }
    .hero-art { width: min(100%, 30rem); justify-self: center; aspect-ratio: 4 / 5; border-radius: 2.4rem; transform: rotate(0.7deg); }
    .route-instrument { padding: 0.7rem; border-radius: 1.6rem; }
    .route-controls { gap: 0.7rem; }
    .endpoint-group { grid-template-columns: minmax(0, 1fr); grid-template-rows: auto repeat(4, 2rem); gap: 0.22rem; }
    .endpoint-group > strong { grid-row: auto; font-size: 0.6rem; }
    .endpoint { min-height: 2rem; grid-template-columns: 1.25rem minmax(0, 1fr); padding: 0.14rem 0.35rem 0.14rem 0.16rem; }
    .endpoint-terminal { width: 1.25rem; height: 1.25rem; }
    .endpoint span:last-child { font-size: 0.58rem; }
    .route-readout, .install-readout, .cli-readout { height: 12.25rem; }
    #install, #cli { padding-block: 4rem; }
    .install-channel { min-height: 1.9rem; padding: 0.22rem 0.48rem; }
    .install-channel courier-brand-icon { width: 0.8rem; height: 0.8rem; }
    .install-channel span { font-size: 0.59rem; }
    .install-actions { gap: 0.5rem; }
    .install-actions a { font-size: 0.58rem; }
    .cli-workspace { height: 50rem; grid-template-columns: 1fr; grid-template-rows: minmax(29rem, 1fr) auto; gap: 0.7rem; }
    .reference-grid { grid-template-columns: minmax(7.5rem, 0.76fr) minmax(0, 1.24fr); }
    .reference-filter { display: grid; min-height: 4rem; padding: 0.32rem 0.42rem; }
    .command-row, .option-row { padding: 0.5rem 0.42rem; }
    .command-row code, .option-row code { font-size: 0.56rem; }
    .option-row > div { display: grid; gap: 0.15rem; }
    .option-row span, .option-row p { font-size: 0.5rem; }
    courier-brand::part(wordmark) { display: none; }
  }

  @media (max-width: 25rem) { nav { display: none; } }
  @media (prefers-reduced-motion: reduce) { .route-signal { animation: none; stroke-dasharray: none; opacity: 0.82; } }
`;
