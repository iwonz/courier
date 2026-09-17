import * as React from "react";
import { createRoot } from "react-dom/client";
import { relaySource } from "@courier/ui";
import { DataApp } from "./app";

export function mountData(root: HTMLElement | null): void {
  if (!root) return;
  const favicon = document.createElement("link");
  favicon.rel = "icon";
  favicon.type = "image/webp";
  favicon.href = relaySource;
  document.head.append(favicon);
  createRoot(root).render(<React.StrictMode><DataApp /></React.StrictMode>);
}

mountData(document.getElementById("root"));
