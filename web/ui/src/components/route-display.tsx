import * as React from "react";
import { ArrowRight } from "lucide-react";
import { cn } from "../lib/utils";

export interface RouteDisplayProps extends React.HTMLAttributes<HTMLDivElement> {
  source: string;
  destination: string;
}

export function RouteDisplay({ source, destination, className, ...props }: RouteDisplayProps): React.JSX.Element {
  return <div className={cn("grid min-w-0 grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] items-center gap-3 py-3", className)} {...props}>
    <code className="truncate font-mono text-sm font-semibold">{source}</code>
    <span className="grid size-8 place-items-center text-primary"><ArrowRight className="size-4" /></span>
    <code className="truncate text-right font-mono text-sm font-semibold">{destination}</code>
  </div>;
}
