import * as React from "react";
import { relayMarkSource, relaySourceForRole, type RelayRole } from "../assets";
import { cn } from "../lib/utils";

export interface BrandProps extends React.HTMLAttributes<HTMLSpanElement> {
  compact?: boolean;
}

export function Brand({ compact = false, className, ...props }: BrandProps): React.JSX.Element {
  return <span className={cn("inline-flex items-center gap-2 font-display font-bold", className)} {...props}>
    <img src={relayMarkSource} alt="" width="256" height="256" className="courier-pixel-image size-9 object-contain" />
    {compact ? null : <span className="text-sm tracking-[.04em]">COURIER CLI</span>}
  </span>;
}

export interface RelaySpriteProps extends React.ImgHTMLAttributes<HTMLImageElement> { readonly role: RelayRole; }

export function RelaySprite({ role, className, alt = "", ...props }: RelaySpriteProps): React.JSX.Element {
  const dimension = role === "delivery" || role === "admin" ? 384 : 512;
  return <img src={relaySourceForRole(role)} width={dimension} height={dimension} alt={alt} className={cn("courier-pixel-image select-none object-contain", className)} draggable={false} {...props} />;
}
