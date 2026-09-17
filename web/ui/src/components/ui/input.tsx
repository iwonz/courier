import * as React from "react";
import { cn } from "../../lib/utils";

export const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(({ className, type, ...props }, ref) => (
  <input ref={ref} type={type} className={cn("courier-pixel-control courier-pixel-focus flex h-10 w-full border border-input bg-background px-3 py-2 text-sm text-foreground outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-45 focus-visible:border-primary", className)} {...props} />
));
Input.displayName = "Input";
