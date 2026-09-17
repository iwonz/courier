import * as React from "react";
import * as ProgressPrimitive from "@radix-ui/react-progress";
import { cn } from "../../lib/utils";

export const Progress = React.forwardRef<React.ElementRef<typeof ProgressPrimitive.Root>, React.ComponentPropsWithoutRef<typeof ProgressPrimitive.Root>>(({ className, value = 0, ...props }, ref) => (
  <ProgressPrimitive.Root ref={ref} className={cn("relative h-2 w-full overflow-hidden bg-secondary", className)} {...props}>
    <ProgressPrimitive.Indicator className="h-full w-full bg-primary transition-transform duration-100 ease-[steps(8,end)]" style={{ transform: `translateX(-${100 - Math.max(0, Math.min(100, value ?? 0))}%)` }} />
  </ProgressPrimitive.Root>
));
Progress.displayName = ProgressPrimitive.Root.displayName;
