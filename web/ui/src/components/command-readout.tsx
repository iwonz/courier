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
  return <div className={cn("grid min-w-0 gap-2 bg-muted/60 p-3 sm:p-4", className)} {...props}>
    {heading ? <span className="text-xs font-semibold text-muted-foreground">{heading}</span> : null}
    {command ? <code className="overflow-x-auto whitespace-pre-wrap break-words font-mono text-sm font-semibold leading-relaxed text-foreground">{command}</code> : null}
    <div className="flex min-h-8 flex-wrap items-center gap-3">
      <Button type="button" variant="ghost" size="sm" onClick={copyCommand} disabled={!command || copyDisabled}>
        <PixelIcon name={status === "copied" ? "check" : "copy"} />{copyLabel}
      </Button>
      <span className={cn("text-xs font-medium text-primary", !live && "sr-only")} aria-live="polite">{live}</span>
    </div>
    {description ? <p className="text-sm text-muted-foreground">{description}</p> : null}
    {details}
    {footerActions ? <div className="flex flex-wrap items-center gap-3 text-sm">{footerActions}</div> : null}
  </div>;
}
