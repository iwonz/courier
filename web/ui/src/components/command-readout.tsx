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
  heading: string;
  command: string;
  copyLabel: string;
  copiedLabel: string;
  copyFailedLabel: string;
  description?: string;
  details?: React.ReactNode;
  footerActions?: React.ReactNode;
  sessionKey?: string;
  copy?: typeof copyText;
}

export function CommandReadout({ heading, command, copyLabel, copiedLabel, copyFailedLabel, description, details, footerActions, sessionKey, copy = copyText, className, ...props }: CommandReadoutProps): React.JSX.Element {
  const [status, setStatus] = React.useState<"idle" | "copied" | "failed">("idle");
  React.useEffect(() => setStatus("idle"), [sessionKey]);
  const copyCommand = async (): Promise<void> => setStatus(await copy(command) ? "copied" : "failed");
  const live = status === "copied" ? copiedLabel : status === "failed" ? copyFailedLabel : "";
  return <div className={cn("grid gap-3 py-2", className)} {...props}>
    <div className="flex min-h-12 items-center justify-between gap-3">
      <span className="text-xs font-semibold text-muted-foreground">{heading}</span>
      <Button type="button" variant="ghost" size="sm" onClick={copyCommand} disabled={!command}>
        <PixelIcon name={status === "copied" ? "check" : "copy"} />{copyLabel}
      </Button>
    </div>
    <div className="grid gap-3 py-2">
      {command ? <code className="overflow-x-auto whitespace-pre-wrap break-words font-mono text-sm font-semibold leading-relaxed text-foreground">{command}</code> : null}
      {description ? <p className="text-sm text-muted-foreground">{description}</p> : null}
      {details}
      <span className="min-h-5 text-xs font-medium text-primary" aria-live="polite">{live}</span>
    </div>
    {footerActions ? <div className="flex flex-wrap items-center gap-3 text-sm">{footerActions}</div> : null}
  </div>;
}
