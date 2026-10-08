import * as React from "react";
import { PixelIcon } from "../icons-react";
import { Button } from "./ui/button";
import { cn } from "../lib/utils";

export async function copyText(text: string, clipboard: Pick<Clipboard, "writeText"> | undefined = globalThis.navigator?.clipboard): Promise<boolean> {
  if (!clipboard || !text) return false;
  try {
    await clipboard.writeText(text);
    return true;
  } catch {
    return false;
  }
}

export interface CommandReadoutProps extends React.HTMLAttributes<HTMLDivElement> {
  heading?: string;
  command: string;
  copyLabel: string;
  copiedLabel: string;
  copyFailedLabel: string;
  description?: string;
  details?: React.ReactNode;
  footerActions?: React.ReactNode;
  sessionKey?: string;
  copyDisabled?: boolean;
  copy?: typeof copyText;
}

export function CommandReadout({ heading, command, copyLabel, copiedLabel, copyFailedLabel, description, details, footerActions, sessionKey, copyDisabled = false, copy = copyText, className, ...props }: CommandReadoutProps): React.JSX.Element {
  const [status, setStatus] = React.useState<"idle" | "copied" | "failed">("idle");
  React.useEffect(() => setStatus("idle"), [sessionKey]);
  const copyCommand = async (): Promise<void> => setStatus(await copy(command) ? "copied" : "failed");
  const live = status === "copied" ? copiedLabel : status === "failed" ? copyFailedLabel : "";
  return <div className={cn("grid min-w-0 gap-3 bg-muted/60 px-4 py-4 sm:px-5 sm:py-5", className)} {...props}>
    {heading ? <span className="text-xs font-semibold text-muted-foreground">{heading}</span> : null}
    {command ? <code className="overflow-x-auto whitespace-pre-wrap break-words font-mono text-sm font-semibold leading-relaxed text-foreground">{command}</code> : null}
    <div className="flex flex-wrap items-center gap-x-3 gap-y-2">
      <Button data-courier-copy-action type="button" variant="ghost" size="sm" className="h-auto min-h-0 justify-start border-0 bg-transparent p-0 text-xs text-muted-foreground hover:bg-transparent hover:text-primary hover:underline hover:underline-offset-4 active:bg-transparent active:text-primary focus-visible:text-primary focus-visible:underline focus-visible:underline-offset-4" onClick={copyCommand} disabled={!command || copyDisabled}>
        <PixelIcon name={status === "copied" ? "check" : "copy"} />{copyLabel}
      </Button>
      <span className={cn("text-xs font-medium text-primary", !live && "sr-only")} aria-live="polite">{live}</span>
    </div>
    {description ? <p className="text-sm text-muted-foreground">{description}</p> : null}
    {details}
    {footerActions ? <div className="flex flex-wrap items-center gap-3 text-sm">{footerActions}</div> : null}
  </div>;
}
