import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import {
  applicableFlags,
  CourierLandingApp,
  endpointExample,
  endpointIcon,
  endpointMessage,
  expandRoutePairs,
} from "./app";
import { contractData } from "./contract";

beforeEach(() => {
  localStorage.clear();
  vi.stubGlobal("matchMedia", () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() }));
});

afterEach(() => {
  document.body.replaceChildren();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

it("projects route pairs, endpoint presentation, and exact flag applicability from the contract", () => {
  const pairs = expandRoutePairs(contractData.routes);
  expect(pairs).toContainEqual(expect.objectContaining({ source: "local", destination: "ssh", routeName: "path-to-path" }));
  expect(pairs).toContainEqual(expect.objectContaining({ source: "webhook", destination: "local", routeName: "webhook-to-path" }));
  expect(endpointExample("local", "source")).toBe("./project");
  expect(endpointExample("ssh", "destination")).toBe("relay@host:/srv/destination/");
  expect(endpointExample("future", "source")).toBe("future");
  expect(endpointMessage("ssh")).toBe("endpointSsh");
  expect(endpointMessage("future")).toBeUndefined();
  expect(endpointIcon("webhook")).toBe("upload");
  expect(endpointIcon("future")).toBe("route");

  const localPair = pairs.find((pair) => pair.source === "local" && pair.destination === "local")!;
  const localFlags = applicableFlags(localPair, contractData.flags).map((flag) => flag.name);
  expect(localFlags).toContain("archive");
  expect(localFlags).not.toContain("upload-rate");
  const sshPair = pairs.find((pair) => pair.source === "ssh" && pair.destination === "ssh")!;
  expect(applicableFlags(sshPair, contractData.flags).map((flag) => flag.name)).toEqual(expect.arrayContaining(["upload-rate", "download-rate"]));
});

it("renders a compact landing and operates the route illustration with pointer, focus, and activation", async () => {
  const element = new CourierLandingApp();
  document.body.append(element);
  await element.updateComplete;
  const root = element.shadowRoot!;
  const text = root.textContent ?? "";

  expect(text).toContain("Move files. Keep control.");
  expect(text).toContain("npm install --global @iwonz/courier");
  expect(text).toContain("courier from <source> to <destination>");
  expect(text).toContain("path-to-path");
  expect(text).toContain("--archive");
  expect(text).not.toContain("Open source · MIT · self-contained");
  expect(root.querySelectorAll("section")).toHaveLength(3);
  expect(root.querySelector(".facts")).toBeNull();
  expect(root.querySelector("#safety")).toBeNull();
  expect(root.querySelector("#examples")).toBeNull();
  expect(root.querySelector("#docs")).toBeNull();
  expect(root.querySelector("footer")).toBeNull();
  expect(root.querySelector(".github-link")?.getAttribute("href")).toBe("https://github.com/iwonz/courier");
  expect(root.querySelectorAll(".channel")).toHaveLength(8);

  const illustrations = [...root.querySelectorAll("courier-mascot")];
  await Promise.all(illustrations.map((illustration) => illustration.updateComplete));
  expect(illustrations).toHaveLength(3);
  expect(illustrations.map((illustration) => illustration.shadowRoot?.querySelector("img")?.src)).toEqual(expect.arrayContaining([
    expect.stringContaining("relay-dispatch"), expect.stringContaining("relay-install"), expect.stringContaining("relay-routing"),
  ]));

  const themeSelector = root.querySelector("courier-theme-selector")!;
  await themeSelector.updateComplete;
  const themeControl = themeSelector.shadowRoot?.querySelector("courier-segmented-control")!;
  await themeControl.updateComplete;
  expect(themeControl.shadowRoot?.querySelector("legend")?.classList.contains("sr-only")).toBe(true);
  expect(themeControl.shadowRoot?.querySelector("button span")).toBeNull();

  const command = () => root.querySelector(".command-shape")?.textContent ?? "";
  expect(command()).toContain("./project to relay@host:/srv/destination/");

  const source = (name: string) => root.querySelector(`.endpoint-group:first-child button[data-endpoint="${name}"]`) as HTMLButtonElement;
  const destination = (name: string) => root.querySelector(`.endpoint-group:last-child button[data-endpoint="${name}"]`) as HTMLButtonElement;

  source("web").dispatchEvent(new Event("pointerenter"));
  await element.updateComplete;
  expect(command()).toContain("web:// to relay@host:/srv/destination/");
  expect(destination("web").getAttribute("aria-disabled")).toBe("true");

  destination("local").dispatchEvent(new FocusEvent("focus"));
  await element.updateComplete;
  expect(command()).toContain("web:// to ./backup/");
  destination("http").click();
  await element.updateComplete;
  expect(command()).toContain("web:// to ./backup/");

  source("local").click();
  await element.updateComplete;
  destination("web").click();
  await element.updateComplete;
  expect(command()).toContain("./project to web://");
  source("webhook").dispatchEvent(new Event("pointerenter"));
  await element.updateComplete;
  expect(command()).toContain("webhook:// to ./backup/");

  expect((element as unknown as { displayEndpoint(name: string): string }).displayEndpoint("future")).toBe("future");
  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(root.textContent).toContain("Переносите файлы. Сохраняйте контроль.");
  expect(root.textContent).toContain("Команды и параметры");
  element.remove();
});

it("supports detached lifecycle and idempotent element definition", async () => {
  new CourierLandingApp().disconnectedCallback();
  vi.resetModules();
  await import("./main");
  expect(customElements.get("courier-landing-app")).toBe(CourierLandingApp);
});
