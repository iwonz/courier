import { LitElement, css, html } from "lit";

export interface PointerPosition {
  readonly x: number;
  readonly y: number;
}

export interface SmoothedPointer {
  readonly position: PointerPosition;
  readonly settled: boolean;
}

export const sceneAmbientPosition: PointerPosition = { x: 72, y: 42 };

function clamp(value: number): number {
  return Math.min(100, Math.max(0, value));
}

export function pointerPosition(bounds: Pick<DOMRect, "left" | "top" | "width" | "height">, clientX: number, clientY: number): PointerPosition {
  if (bounds.width <= 0 || bounds.height <= 0) return { x: 50, y: 50 };
  return {
    x: clamp(((clientX - bounds.left) / bounds.width) * 100),
    y: clamp(((clientY - bounds.top) / bounds.height) * 100),
  };
}

export function smoothPointerPosition(current: PointerPosition, target: PointerPosition, elapsedMilliseconds: number): SmoothedPointer {
  const elapsed = Math.min(48, Math.max(0, elapsedMilliseconds));
  const factor = 1 - Math.exp(-elapsed / 42);
  const position = {
    x: current.x + (target.x - current.x) * factor,
    y: current.y + (target.y - current.y) * factor,
  };
  const settled = Math.hypot(target.x - position.x, target.y - position.y) < 0.04;
  return { position: settled ? target : position, settled };
}

export class CourierScene extends LitElement {
  static properties = {
    eager: { type: Boolean },
    mobileSource: { type: String, attribute: "mobile-source" },
    source: { type: String },
  };

  static styles = css`
    :host { --scene-pointer-x: 72%; --scene-pointer-y: 42%; position: absolute; display: block; inset: 0; overflow: hidden; background: var(--courier-graphite-900, #151714); pointer-events: auto; isolation: isolate; }
    .tracking { position: absolute; inset: 0; }
    courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; }
    courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; filter: saturate(0.88) contrast(1.02); }
    .base { z-index: 0; transform: none; }
    .refracted { z-index: 1; opacity: 0.56; transform: scale(1.009); transform-origin: var(--scene-pointer-x) var(--scene-pointer-y); filter: saturate(1.055) contrast(1.02); mask-image: radial-gradient(circle clamp(3.8rem, 8vw, 7.2rem) at var(--scene-pointer-x) var(--scene-pointer-y), #000 36%, transparent 72%), radial-gradient(circle clamp(2.8rem, 6vw, 5.4rem) at calc(var(--scene-pointer-x) - 4.5%) calc(var(--scene-pointer-y) + 2.5%), #000 32%, transparent 74%), radial-gradient(circle clamp(2.4rem, 5vw, 4.8rem) at calc(var(--scene-pointer-x) + 4%) calc(var(--scene-pointer-y) - 3.5%), #000 30%, transparent 72%); mask-repeat: no-repeat; will-change: mask-position, transform; animation: scene-lobes 7s ease-in-out infinite alternate; }
    .glow { position: absolute; z-index: 2; inset: 0; background: radial-gradient(circle clamp(8rem, 20vw, 18rem) at var(--scene-pointer-x) var(--scene-pointer-y), color-mix(in srgb, var(--courier-signal, #d4ff45) 13%, transparent), transparent 68%), radial-gradient(circle clamp(4rem, 9vw, 8rem) at calc(var(--scene-pointer-x) - 5%) calc(var(--scene-pointer-y) + 4%), rgb(255 255 255 / 0.065), transparent 72%), radial-gradient(circle clamp(3rem, 7vw, 6rem) at calc(var(--scene-pointer-x) + 6%) calc(var(--scene-pointer-y) - 5%), color-mix(in srgb, var(--courier-beak, #ff8758) 5%, transparent), transparent 76%); filter: blur(0.55rem); mix-blend-mode: screen; pointer-events: none; animation: scene-glow 8s ease-in-out infinite alternate; }
    .veil { position: absolute; z-index: 3; inset: 0; background: linear-gradient(90deg, rgb(8 10 8 / 0.1), transparent 28% 72%, rgb(8 10 8 / 0.16)), linear-gradient(180deg, rgb(8 10 8 / 0.08), transparent 23% 82%, rgb(8 10 8 / 0.22)); pointer-events: none; }
    @keyframes scene-lobes { to { filter: saturate(1.06) contrast(1.025) blur(0.08rem); } }
    @keyframes scene-glow { to { opacity: 0.86; } }
    @media (pointer: coarse), (hover: none) {
      :host { --scene-pointer-x: 68%; --scene-pointer-y: 40%; }
      .refracted { display: none; }
      .glow { opacity: 0.74; animation: none; }
    }
    @media (prefers-reduced-motion: reduce) {
      .refracted { display: none; }
      .glow { opacity: 0.68; animation: none; }
    }
  `;

  eager = false;
  mobileSource = "";
  source = "";
  private frame = 0;
  private current: PointerPosition = sceneAmbientPosition;
  private target?: PointerPosition;
  private previousTimestamp?: number;
  private returning = false;
  private readonly handlePointerMove = (event: PointerEvent): void => this.move(event);
  private readonly handlePointerLeave = (): void => this.leave();

  connectedCallback(): void {
    super.connectedCallback();
    this.addEventListener("pointermove", this.handlePointerMove);
    this.addEventListener("pointerleave", this.handlePointerLeave);
  }

  disconnectedCallback(): void {
    this.removeEventListener("pointermove", this.handlePointerMove);
    this.removeEventListener("pointerleave", this.handlePointerLeave);
    if (this.frame) globalThis.cancelAnimationFrame(this.frame);
    this.frame = 0;
    this.target = undefined;
    this.previousTimestamp = undefined;
    super.disconnectedCallback();
  }

  private tracksPointer(): boolean {
    return globalThis.matchMedia?.("(hover: hover) and (pointer: fine)").matches === true && !globalThis.matchMedia?.("(prefers-reduced-motion: reduce)").matches;
  }

  private move(event: PointerEvent): void {
    if (!this.tracksPointer()) return;
    this.target = pointerPosition(this.getBoundingClientRect(), event.clientX, event.clientY);
    this.returning = false;
    this.schedule();
  }

  private leave(): void {
    if (!this.tracksPointer()) return;
    this.target = sceneAmbientPosition;
    this.returning = true;
    this.schedule();
  }

  private schedule(): void {
    if (this.frame) return;
    this.frame = globalThis.requestAnimationFrame((timestamp) => this.advance(timestamp));
  }

  private advance(timestamp: number): void {
    this.frame = 0;
    const target = this.target;
    if (!target) return;
    const elapsed = this.previousTimestamp === undefined ? 16 : timestamp - this.previousTimestamp;
    this.previousTimestamp = timestamp;
    const next = smoothPointerPosition(this.current, target, elapsed);
    this.current = next.position;
    this.style.setProperty("--scene-pointer-x", `${this.current.x}%`);
    this.style.setProperty("--scene-pointer-y", `${this.current.y}%`);
    if (!next.settled) {
      this.schedule();
      return;
    }
    this.previousTimestamp = undefined;
    if (this.returning) {
      this.target = undefined;
      this.returning = false;
      this.style.removeProperty("--scene-pointer-x");
      this.style.removeProperty("--scene-pointer-y");
    }
  }

  protected render() {
    return html`<div class="tracking"><courier-mascot class="base" ?eager=${this.eager} alt="" .source=${this.source} .mobileSource=${this.mobileSource}></courier-mascot><courier-mascot class="refracted" aria-hidden="true" alt="" .source=${this.source} .mobileSource=${this.mobileSource}></courier-mascot><span class="glow" aria-hidden="true"></span><span class="veil" aria-hidden="true"></span></div>`;
  }
}
