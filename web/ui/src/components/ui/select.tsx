import * as React from "react";
import * as SelectPrimitive from "@radix-ui/react-select";
import { PixelIcon } from "../../icons-react";
import { cn } from "../../lib/utils";

export const Select = SelectPrimitive.Root;
export const SelectGroup = SelectPrimitive.Group;
export const SelectValue = SelectPrimitive.Value;

export const SelectTrigger = React.forwardRef<React.ElementRef<typeof SelectPrimitive.Trigger>, React.ComponentPropsWithoutRef<typeof SelectPrimitive.Trigger>>(({ className, children, ...props }, ref) => (
  <SelectPrimitive.Trigger ref={ref} className={cn("courier-pixel-control courier-pixel-focus flex h-10 w-full items-center justify-between gap-2 border border-input bg-background px-3 py-2 text-sm outline-none data-[placeholder]:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-45 [&>span]:truncate", className)} {...props}>
    {children}<SelectPrimitive.Icon asChild><PixelIcon name="chevron-down" className="size-4 opacity-70" /></SelectPrimitive.Icon>
  </SelectPrimitive.Trigger>
));
SelectTrigger.displayName = SelectPrimitive.Trigger.displayName;

export const SelectContent = React.forwardRef<React.ElementRef<typeof SelectPrimitive.Content>, React.ComponentPropsWithoutRef<typeof SelectPrimitive.Content>>(({ className, children, position = "popper", ...props }, ref) => (
  <SelectPrimitive.Portal>
    <SelectPrimitive.Content ref={ref} position={position} className={cn("courier-pixel-control relative z-50 max-h-80 min-w-[8rem] overflow-hidden border border-border bg-popover text-popover-foreground", position === "popper" && "data-[side=bottom]:translate-y-1 data-[side=top]:-translate-y-1", className)} {...props}>
      <SelectPrimitive.ScrollUpButton className="flex h-6 items-center justify-center"><PixelIcon name="chevron-up" className="size-4" /></SelectPrimitive.ScrollUpButton>
      <SelectPrimitive.Viewport className="p-1">{children}</SelectPrimitive.Viewport>
      <SelectPrimitive.ScrollDownButton className="flex h-6 items-center justify-center"><PixelIcon name="chevron-down" className="size-4" /></SelectPrimitive.ScrollDownButton>
    </SelectPrimitive.Content>
  </SelectPrimitive.Portal>
));
SelectContent.displayName = SelectPrimitive.Content.displayName;

export const SelectItem = React.forwardRef<React.ElementRef<typeof SelectPrimitive.Item>, React.ComponentPropsWithoutRef<typeof SelectPrimitive.Item>>(({ className, children, ...props }, ref) => (
  <SelectPrimitive.Item ref={ref} className={cn("relative flex w-full cursor-default select-none items-center py-2 pl-8 pr-2 text-sm outline-none focus:bg-primary focus:text-primary-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-45", className)} {...props}>
    <span className="absolute left-2 grid size-4 place-items-center"><SelectPrimitive.ItemIndicator><PixelIcon name="check" className="size-4" /></SelectPrimitive.ItemIndicator></span>
    <SelectPrimitive.ItemText>{children}</SelectPrimitive.ItemText>
  </SelectPrimitive.Item>
));
SelectItem.displayName = SelectPrimitive.Item.displayName;
