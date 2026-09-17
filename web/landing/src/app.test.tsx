import * as React from "react";
import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import {
  LandingApp,
  applicableFlags,
  commandFlags,
  endpointExample,
  endpointIcon,
  endpointMessage,
  endpointPresentation,
  expandRoutePairs,
  installs,
  localizedEndpointDescription,
  localizedEndpointLabel,
} from "./app";
import { contractData } from "./contract";
import { type LandingMessage } from "./catalog";
import { mountLanding } from "./main";

function rect(left: number, top: number, width: number, height: number): DOMRect {
  return { left, top, width, height, right: left + width, bottom: top + height, x: left, y: top, toJSON: () => ({}) };
}

beforeEach(() => {
  localStorage.clear();
  vi.stubGlobal("matchMedia", vi.fn(() => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })));
});

afterEach(() => {
  cleanup();
  document.body.replaceChildren();
  document.head.querySelectorAll('link[rel="icon"]').forEach((node) => node.remove());
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe("landing contract projection", () => {
  it("projects routes, endpoint vocabulary, exact flags, and command compatibility", () => {
    const pairs = expandRoutePairs(contractData.routes);
    expect(pairs).toContainEqual(expect.objectContaining({ source: "local", destination: "ssh", routeName: "path-to-path" }));
    expect(pairs).toContainEqual(expect.objectContaining({ source: "webhook", destination: "local", routeName: "webhook-to-path" }));
    expect(endpointExample("local", "source")).toBe("./project");
    expect(endpointExample("ssh", "destination")).toBe("courier@host:/srv/destination/");
    expect(endpointExample("future", "source")).toBe("future");
    expect(endpointMessage("ssh", "source")).toBe("endpointRemote");
    expect(endpointMessage("http", "destination")).toBe("endpointWebHook");
    expect(endpointMessage("future", "source")).toBeUndefined();
    expect(endpointPresentation("web", "source")).toEqual(expect.objectContaining({ icon: "browser-upload" }));
    expect(endpointIcon("webhook", "source")).toBe("webhook-in");
    expect(endpointIcon("http", "destination")).toBe("webhook-out");
    expect(endpointIcon("future", "destination")).toBe("route");
    const translate = (message: LandingMessage): string => `translated:${message}`;
    expect(localizedEndpointLabel("local", "source", translate)).toBe("translated:endpointLocal");
    expect(localizedEndpointLabel("future", "source", translate)).toBe("future");
    expect(localizedEndpointDescription("web", "source", translate)).toBe("translated:sourceWebDescription");
    expect(localizedEndpointDescription("future", "destination", translate)).toBe("future");

    const localPair = pairs.find((pair) => pair.source === "local" && pair.destination === "local")!;
    expect(applicableFlags(localPair, contractData.flags).map((flag) => flag.name)).toContain("archive");
    const sshPair = pairs.find((pair) => pair.source === "ssh" && pair.destination === "ssh")!;
    expect(applicableFlags(sshPair, contractData.flags).map((flag) => flag.name)).toEqual(expect.arrayContaining(["upload-rate", "download-rate"]));
    expect(commandFlags("", true, contractData.commands, contractData.flags)).toHaveLength(contractData.flags.length);
    expect(commandFlags("from", false, contractData.commands, contractData.flags)).toHaveLength(contractData.flags.length);
    expect(commandFlags("ui-start", true, contractData.commands, contractData.flags).map((flag) => flag.name)).toEqual(["listen", "background"]);
    expect(commandFlags("servers-stop", true, contractData.commands, contractData.flags).map((flag) => flag.name)).toEqual(["all"]);
    expect(commandFlags("servers", true, contractData.commands, contractData.flags)).toEqual([]);
    expect(commandFlags("unknown", true, contractData.commands, contractData.flags)).toEqual([]);
    expect(installs.map((install) => install.icon)).toEqual(["curl", "wget", "powershell", "npm", "npx", "yarn", "pnpm", "homebrew", "scoop"]);
  });
});

describe("React landing", () => {
  it("mounts with a square Relay favicon and ignores an absent root", async () => {
    mountLanding(null);
    const root = document.createElement("div");
    document.body.append(root);
    mountLanding(root);
    await waitFor(() => expect(root.querySelector("h1")?.textContent).toBe("From here to anywhere."));
    expect(document.head.querySelector<HTMLLinkElement>('link[rel="icon"]')?.href).toContain("courier-relay-tech-v1");
  });

  it("renders three natural sections, a compact square mascot, and secure external links", () => {
    render(<LandingApp />);
    expect(screen.getByRole("heading", { level: 1 }).textContent).toBe("From here to anywhere.");
    expect(document.querySelectorAll("main > section")).toHaveLength(3);
    const mascot = document.querySelector<HTMLImageElement>('img[width="768"][height="768"]')!;
    expect(mascot.src).toContain("courier-relay-tech-v1");
    expect(document.body.textContent).not.toContain("Run demo");
    expect(document.body.textContent).not.toContain("One binary plans the route");
    expect([...document.querySelectorAll<HTMLAnchorElement>('a[href^="https://"]')].every((link) => link.target === "_blank" && link.rel === "noopener noreferrer")).toBe(true);
  });

  it("selects only valid routes and keeps the generated command contract-backed", () => {
    render(<LandingApp />);
    fireEvent.click(screen.getAllByRole("button", { name: "Web Hook" })[0]!);
    expect(screen.getByText(/courier from webhook:\/\//).textContent).toContain("courier@host:/srv/destination/");
    fireEvent.click(screen.getAllByRole("button", { name: "Local" })[0]!);
    fireEvent.click(screen.getAllByRole("button", { name: "Remote" })[1]!);
    const webButtons = screen.getAllByRole("button", { name: "Web" });
    fireEvent.click(webButtons[0]!);
    expect(screen.getByText(/courier from web:\/\//).textContent).toContain("courier@host:/srv/destination/");
    expect(webButtons[1]!.hasAttribute("disabled")).toBe(true);
    fireEvent.click(screen.getAllByRole("button", { name: "Local" })[1]!);
    expect(screen.getByText(/courier from web:\/\//).textContent).toContain("./backup/");
    fireEvent.click(screen.getAllByRole("button", { name: "Web Hook" })[0]!);
    expect(screen.getByText(/courier from webhook:\/\//).textContent).toContain("./backup/");
    fireEvent.click(screen.getAllByRole("button", { name: "Local" })[0]!);
    fireEvent.click(screen.getAllByRole("button", { name: "Web" })[1]!);
    expect(screen.getByText(/courier from \.\/project/).textContent).toContain("web://");
    fireEvent.click(screen.getAllByRole("button", { name: "Web" })[0]!);
    expect(screen.getByText(/courier from web:\/\//).textContent).toContain("./backup/");
  });

  it("measures a bezier connector, responds to resize, and cleans its observer", () => {
    const disconnect = vi.fn();
    let callback: ResizeObserverCallback | undefined;
    class ResizeObserverStub {
      observe = vi.fn();
      disconnect = disconnect;
      constructor(next: ResizeObserverCallback) { callback = next; }
    }
    vi.stubGlobal("ResizeObserver", ResizeObserverStub);
    vi.spyOn(HTMLElement.prototype, "getBoundingClientRect").mockImplementation(function (this: HTMLElement) {
      if (this.tagName === "BUTTON") {
        const destination = this.parentElement?.firstElementChild?.textContent === "Destination";
        return destination ? rect(620, 120, 120, 32) : rect(80, 40, 120, 32);
      }
      return rect(0, 0, 820, 300);
    });
    const view = render(<LandingApp />);
    expect(document.querySelector('path[vector-effect="non-scaling-stroke"]')?.getAttribute("d")).toContain(" C ");
    callback?.([], {} as ResizeObserver);
    fireEvent(globalThis, new Event("resize"));
    view.unmount();
    expect(disconnect).toHaveBeenCalled();
  });

  it("selects installation channels and copies the exact command", async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", { configurable: true, value: { writeText } });
    render(<LandingApp />);
    fireEvent.click(screen.getByRole("button", { name: /Homebrew/ }));
    const command = screen.getByText(/brew tap iwonz\/courier/).textContent!;
    const install = document.querySelector("#install")!;
    fireEvent.click(Array.from(install.querySelectorAll("button")).find((button) => button.textContent?.includes("Copy command"))!);
    await waitFor(() => expect(writeText).toHaveBeenCalledWith(command));
    expect(screen.getByText("Copied")).toBeTruthy();
  });

  it("filters options by selected command, exposes empty states, and clears selection", () => {
    render(<LandingApp />);
    const checkbox = screen.getByRole("checkbox");
    expect(checkbox.hasAttribute("disabled")).toBe(true);
    const uiStart = screen.getByRole("button", { name: "courier ui start [options]" });
    fireEvent.click(uiStart);
    expect(checkbox.hasAttribute("disabled")).toBe(false);
    expect(screen.getByText("--listen <host:port>")).toBeTruthy();
    expect(screen.getByText("--background")).toBeTruthy();
    fireEvent.click(checkbox);
    expect(screen.getAllByText("--archive").length).toBeGreaterThan(0);
    fireEvent.click(uiStart);
    expect(checkbox.hasAttribute("disabled")).toBe(true);
    fireEvent.click(screen.getByRole("button", { name: "courier servers" }));
    expect(screen.getByText("This command has no options.")).toBeTruthy();
    fireEvent.click(screen.getByRole("button", { name: "courier servers stop <uuid>|--all" }));
    expect(screen.getByText("--all")).toBeTruthy();
  });

  it("cycles locale without introducing removed explanatory labels", () => {
    render(<LandingApp />);
    const localeButton = screen.getByRole("button", { name: /Language:/ });
    fireEvent.click(localeButton);
    expect(screen.getByRole("heading", { level: 1 }).textContent).toBe("Отсюда — куда угодно.");
    expect(document.body.textContent).not.toContain("Все связи взяты из опубликованного контракта CLI");
  });
});
