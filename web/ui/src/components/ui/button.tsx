import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "../../lib/utils";

export const buttonVariants = cva(
  "courier-pixel-control courier-pixel-focus inline-flex shrink-0 items-center justify-center gap-2 whitespace-nowrap border border-transparent text-sm font-semibold outline-none transition-none disabled:pointer-events-none disabled:opacity-45 [&_svg]:pointer-events-none [&_svg]:size-4",
  {
    variants: {
      variant: {
        default: "border-[var(--terminal-fill-action)] bg-[var(--terminal-fill-action)] text-black hover:border-[var(--terminal-fill-selection)] hover:bg-[var(--terminal-fill-selection)] active:bg-foreground active:text-background",
        secondary: "border-[var(--terminal-fill-selection)] bg-[var(--terminal-fill-selection)] text-black hover:border-[var(--terminal-fill-action)] hover:bg-[var(--terminal-fill-action)]",
        outline: "border-border bg-background text-foreground hover:border-primary hover:bg-secondary hover:text-primary",
        ghost: "text-foreground hover:bg-secondary hover:text-primary active:bg-foreground active:text-background",
        destructive: "border-destructive bg-destructive text-white hover:bg-background hover:text-destructive",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-8 px-3 text-xs",
        lg: "h-12 px-6 text-base",
        icon: "size-10 p-0",
        "icon-sm": "size-8 p-0",
      },
    },
    defaultVariants: { variant: "default", size: "default" },
  },
);

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Component = asChild ? Slot : "button";
    return <Component ref={ref} className={cn(buttonVariants({ variant, size }), className)} {...props} />;
  },
);
Button.displayName = "Button";
