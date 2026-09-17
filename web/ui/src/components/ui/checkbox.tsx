import * as React from "react";
import * as CheckboxPrimitive from "@radix-ui/react-checkbox";
import { PixelIcon } from "../../icons-react";
import { cn } from "../../lib/utils";

export const Checkbox = React.forwardRef<React.ElementRef<typeof CheckboxPrimitive.Root>, React.ComponentPropsWithoutRef<typeof CheckboxPrimitive.Root>>(
  ({ className, ...props }, ref) => (
    <CheckboxPrimitive.Root ref={ref} className={cn("courier-pixel-control courier-pixel-focus peer size-5 shrink-0 border border-border bg-background outline-none data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground disabled:cursor-not-allowed disabled:opacity-45", className)} {...props}>
      <CheckboxPrimitive.Indicator className="grid place-items-center"><PixelIcon name="check" className="size-3.5" /></CheckboxPrimitive.Indicator>
    </CheckboxPrimitive.Root>
  ),
);
Checkbox.displayName = CheckboxPrimitive.Root.displayName;
