import * as React from "react";
import { relayMarkSource, relaySource } from "../assets";
import { cn } from "../lib/utils";

export interface BrandProps extends React.HTMLAttributes<HTMLSpanElement> {
  compact?: boolean;
}

export function Brand({ compact = false, className, ...props }: BrandProps): React.JSX.Element {
  return <span className={cn("inline-flex items-center gap-2.5 font-black tracking-[-0.045em]", className)} {...props}>
    <img src={relayMarkSource} alt="" width="512" height="512" className="size-9 object-contain" />
    {compact ? null : <span className="text-[.82rem] tracking-[.06em]">COURIER CLI</span>}
  </span>;
}

export interface MascotProps extends React.ImgHTMLAttributes<HTMLImageElement> {}

export function Mascot({ className, alt = "", ...props }: MascotProps): React.JSX.Element {
  return <img src={relaySource} width="768" height="768" alt={alt} className={cn("select-none object-contain", className)} draggable={false} {...props} />;
}
