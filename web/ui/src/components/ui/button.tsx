import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "../../lib/utils";

export const buttonVariants = cva(
  "courier-pixel-control courier-pixel-focus inline-flex shrink-0 items-center justify-center gap-2 whitespace-nowrap border border-transparent text-sm font-semibold outline-none transition-[color,background-color,border-color,box-shadow,transform] duration-75 ease-[steps(1,end)] disabled:pointer-events-none disabled:opacity-45 active:translate-x-0.5 active:translate-y-0.5 active:shadow-none [&_svg]:pointer-events-none [&_svg]:size-4",
  {
    variants: {
      variant: {
        default: "border-primary bg-primary text-primary-foreground shadow-[2px_2px_0_var(--pixel-shadow)] hover:-translate-x-px hover:-translate-y-px hover:shadow-[3px_3px_0_var(--pixel-shadow)]",
        secondary: "border-primary bg-primary text-primary-foreground shadow-[inset_2px_2px_0_color-mix(in_srgb,var(--pixel-paper)_24%,transparent)]",
        outline: "border-border bg-background text-foreground shadow-[2px_2px_0_var(--pixel-shadow)] hover:border-primary hover:text-primary",
        ghost: "text-foreground hover:bg-secondary hover:text-secondary-foreground",
        destructive: "border-destructive bg-destructive text-white shadow-[2px_2px_0_var(--pixel-shadow)]",
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
