import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import { CourierLandingApp } from "./app";

beforeEach(() => {
  localStorage.clear();
  vi.stubGlobal("matchMedia", () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() }));
});

afterEach(() => {
  document.body.replaceChildren();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

it("renders generated install, command, route, option, and documentation sections", async () => {
  const element = new CourierLandingApp();
  document.body.append(element);
  await element.updateComplete;
  const text = element.shadowRoot?.textContent ?? "";
  expect(text).toContain("Move files. Keep control.");
  expect(text).toContain("npm install --global @iwonz/courier");
  expect(text).toContain("courier from <source> to <destination>");
  expect(text).toContain("path-to-path");
  expect(text).toContain("--archive");
  expect(element.shadowRoot?.querySelectorAll("section").length).toBe(7);

  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(element.shadowRoot?.textContent).toContain("Переносите файлы. Сохраняйте контроль.");
  element.remove();
});

it("supports detached lifecycle and idempotent element definition", async () => {
  new CourierLandingApp().disconnectedCallback();
  vi.resetModules();
  await import("./main");
  expect(customElements.get("courier-landing-app")).toBe(CourierLandingApp);
});
