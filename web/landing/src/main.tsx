import * as React from "react";
import { createRoot } from "react-dom/client";
import { relayMarkSource } from "@courier/ui";
import { LandingApp } from "./app";

export function mountLanding(root: HTMLElement | null): void {
  if (!root) return;
  const favicon = document.createElement("link");
  favicon.rel = "icon";
  favicon.type = "image/webp";
  favicon.href = relayMarkSource;
  document.head.append(favicon);
  createRoot(root).render(<React.StrictMode><LandingApp /></React.StrictMode>);
}

mountLanding(document.getElementById("root"));
