import * as React from "react";
import { createRoot } from "react-dom/client";
import { relayMarkSource } from "@courier/ui";
import { AdminApp } from "./app";

export function mountAdmin(root: HTMLElement | null): void {
  if (!root) return;
  const favicon = document.createElement("link");
  favicon.rel = "icon";
  favicon.type = "image/webp";
  favicon.href = relayMarkSource;
  document.head.append(favicon);
  createRoot(root).render(<React.StrictMode><AdminApp /></React.StrictMode>);
}

mountAdmin(document.getElementById("root"));
