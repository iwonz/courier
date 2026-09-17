import * as React from "react";
import { cn } from "../../lib/utils";

export const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(({ className, type, ...props }, ref) => (
  <input ref={ref} type={type} className={cn("flex h-10 w-full rounded-xl border border-input bg-background/70 px-3 py-2 text-sm text-foreground shadow-sm outline-none transition-colors placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-45 focus-visible:border-primary/60 focus-visible:ring-2 focus-visible:ring-ring", className)} {...props} />
));
Input.displayName = "Input";
