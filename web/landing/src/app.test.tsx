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
    expect(installs.map((install) => install.icon)).toEqual(["curl", undefined, "powershell", "npm", undefined, "yarn", "pnpm", "homebrew", "scoop"]);
  });
});

describe("React landing", () => {
  it("mounts with a square Relay favicon and ignores an absent root", async () => {
    mountLanding(null);
    const root = document.createElement("div");
    document.body.append(root);
    mountLanding(root);
    await waitFor(() => expect(root.querySelector("h1")?.textContent).toBe("From here to anywhere."));
    expect(document.head.querySelector<HTMLLinkElement>('link[rel="icon"]')?.href).toContain("courier-relay-pixel-mark-v3");
  });

  it("renders three natural sections, a compact square mascot, and secure external links", () => {
    render(<LandingApp />);
    expect(screen.getByRole("heading", { level: 1 }).textContent).toBe("From here to anywhere.");
    expect(document.querySelectorAll("main > section")).toHaveLength(3);
    expect(document.querySelector("[data-courier-contract-metrics]")).toBeNull();
    const mascot = document.querySelector<HTMLImageElement>('img[width="512"][height="512"]')!;
    expect(mascot.src).toContain("courier-relay-pixel-route-v3");
    expect(mascot.className).toContain("courier-pixel-image");
    expect(document.body.textContent).not.toContain("Run demo");
    expect(document.body.textContent).not.toContain("One binary plans the route");
    expect(document.querySelector("header nav")).toBeNull();
    expect(document.querySelector("header")?.className).not.toContain("border");
    expect(document.querySelector("main")?.className).not.toContain("pt-");
    expect(document.querySelector("[data-courier-route-composition]")?.className).not.toContain("border");
    expect(document.querySelectorAll('header img[src*="github-"]')).toHaveLength(2);
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

  it("measures a grid-snapped bezier connector, responds to resize, and cleans its observer", () => {
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
    expect(document.querySelector('path[vector-effect="non-scaling-stroke"]')?.getAttribute("d")).toContain(" H ");
    const signal = document.querySelector<HTMLElement>(".courier-route-signal")!;
    expect(signal.className).toContain("size-1");
    expect(signal.tagName).toBe("SPAN");
    expect(signal.style.offsetAnchor).toBe("center");
    callback?.([], {} as ResizeObserver);
    fireEvent(globalThis, new Event("resize"));
    view.unmount();
    expect(disconnect).toHaveBeenCalled();
  });

  it("selects installation channels and copies the exact command", async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", { configurable: true, value: { writeText } });
    render(<LandingApp />);
    const homebrew = screen.getByRole("tab", { name: /Homebrew/ });
    expect(homebrew.className).toContain("courier-install-tab");
    fireEvent.click(homebrew);
    expect(homebrew.getAttribute("aria-selected")).toBe("true");
    expect(screen.getAllByText("Homebrew", { exact: true })).toHaveLength(1);
    const command = screen.getByText(/brew tap iwonz\/courier/).textContent!;
    expect(command).not.toContain("--cask");
    const install = document.querySelector("#install")!;
    fireEvent.click(Array.from(install.querySelectorAll("button")).find((button) => button.textContent?.includes("Copy command"))!);
    await waitFor(() => expect(writeText).toHaveBeenCalledWith(command));
    expect(screen.getByText("Copied")).toBeTruthy();
    fireEvent.click(screen.getByRole("tab", { name: "wget" }));
    expect(document.querySelector('#install img[data-brand-name="GNU Wget"]')).toBeNull();
    expect(screen.getAllByText("wget", { exact: true })).toHaveLength(1);
  });

  it("filters options by selected command, exposes empty states, and clears selection", () => {
    render(<LandingApp />);
    const checkbox = screen.getByRole("checkbox", { name: "Compatible with selected command" });
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
    expect(document.querySelector("[data-courier-cli-empty]")?.className).not.toContain("border");
    fireEvent.click(screen.getByRole("button", { name: "courier servers stop <uuid>|--all" }));
    expect(screen.getByText("--all")).toBeTruthy();
  });

  it("constructs and copies exact POSIX and PowerShell commands", async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", { configurable: true, value: { writeText } });
    render(<LandingApp />);
    fireEvent.click(screen.getByRole("button", { name: "courier from <source> to <destination> [options]" }));
    const cli = document.querySelector("#cli")!;
    const copy = Array.from(cli.querySelectorAll("button")).find((button) => button.textContent?.includes("Copy command"))!;
    expect(copy.hasAttribute("disabled")).toBe(true);
    fireEvent.change(screen.getByRole("textbox", { name: "source" }), { target: { value: "folder one" } });
    fireEvent.change(screen.getByRole("textbox", { name: "destination" }), { target: { value: "host:/srv/it's" } });
    fireEvent.click(screen.getByRole("checkbox", { name: "--extract" }));
    const maximum = screen.getByRole("textbox", { name: "--max-extracted-size <size|unlimited>" });
    expect(maximum.hasAttribute("disabled")).toBe(false);
    fireEvent.change(maximum, { target: { value: "2GiB" } });
    const firstExclude = screen.getByRole("textbox", { name: "--exclude <pattern> 1" });
    fireEvent.change(firstExclude, { target: { value: "*.tmp" } });
    const excludeControl = document.querySelector('[data-builder-flag="exclude"]')!;
    fireEvent.click(Array.from(excludeControl.querySelectorAll("button")).find((button) => button.getAttribute("aria-label") === "Add value")!);
    fireEvent.change(screen.getByRole("textbox", { name: "--exclude <pattern> 2" }), { target: { value: "old files/*" } });
    fireEvent.click(screen.getByRole("combobox", { name: "--auth <none|basic|password>" }));
    fireEvent.click(screen.getByRole("option", { name: "basic" }));
    const posix = "courier from 'folder one' to 'host:/srv/it'\"'\"'s' --extract --auth basic --exclude '*.tmp' --exclude 'old files/*' --max-extracted-size 2GiB";
    expect(cli.textContent).toContain(posix);
    expect(copy.hasAttribute("disabled")).toBe(false);
    fireEvent.click(copy);
    await waitFor(() => expect(writeText).toHaveBeenLastCalledWith(posix));
    fireEvent.click(Array.from(cli.querySelectorAll("button")).find((button) => button.textContent === "PowerShell")!);
    const powershell = "courier from 'folder one' to 'host:/srv/it''s' --extract --auth basic --exclude '*.tmp' --exclude 'old files/*' --max-extracted-size 2GiB";
    expect(cli.textContent).toContain(powershell);
    fireEvent.click(copy);
    await waitFor(() => expect(writeText).toHaveBeenLastCalledWith(powershell));
    fireEvent.click(Array.from(cli.querySelectorAll("button")).find((button) => button.textContent === "POSIX")!);
    fireEvent.click(Array.from(excludeControl.querySelectorAll("button")).find((button) => button.getAttribute("aria-label") === "Remove value")!);
    expect(cli.textContent).not.toContain("*.tmp");
  });

  it("enforces dependencies, conflicts, stop UUID exclusivity, help paths, and reset", () => {
    render(<LandingApp />);
    const from = screen.getByRole("button", { name: "courier from <source> to <destination> [options]" });
    fireEvent.click(from);
    const maximum = screen.getByRole("textbox", { name: "--max-extracted-size <size|unlimited>" });
    expect(maximum.hasAttribute("disabled")).toBe(true);
    fireEvent.click(screen.getByRole("checkbox", { name: "--extract" }));
    fireEvent.change(maximum, { target: { value: "1GiB" } });
    fireEvent.click(screen.getByRole("checkbox", { name: "--archive" }));
    expect(screen.getByRole("checkbox", { name: "--extract" }).getAttribute("aria-checked")).toBe("false");
    expect(maximum.hasAttribute("disabled")).toBe(true);

    fireEvent.click(screen.getByRole("button", { name: "courier servers stop <uuid>|--all" }));
    const uuid = screen.getByRole("textbox", { name: "uuid" });
    const cliCopy = Array.from(document.querySelector("#cli")!.querySelectorAll("button")).find((button) => button.textContent?.includes("Copy command"))!;
    fireEvent.change(uuid, { target: { value: "bad" } });
    expect(cliCopy.hasAttribute("disabled")).toBe(true);
    fireEvent.click(screen.getByRole("checkbox", { name: "--all" }));
    expect(uuid.hasAttribute("disabled")).toBe(true);
    expect(document.querySelector("#cli")?.textContent).toContain("courier servers stop --all");
    fireEvent.click(screen.getByRole("checkbox", { name: "--all" }));
    fireEvent.change(uuid, { target: { value: "550e8400-e29b-41d4-a716-446655440000" } });
    expect(cliCopy.hasAttribute("disabled")).toBe(false);

    fireEvent.click(screen.getByRole("button", { name: "courier help [command]" }));
    fireEvent.click(screen.getByRole("combobox", { name: "command" }));
    fireEvent.click(screen.getByRole("option", { name: "servers stop" }));
    expect(document.querySelector("#cli")?.textContent).toContain("courier help servers stop");
    fireEvent.click(from);
    expect(screen.getByRole<HTMLInputElement>("textbox", { name: "source" }).value).toBe("");
  });

  it("cycles locale without introducing removed explanatory labels", () => {
    render(<LandingApp />);
    expect(document.querySelector("[data-courier-route-composition]")).toBeTruthy();
    expect(document.querySelector("[data-courier-cli-registry]")).toBeTruthy();
    expect(document.querySelector("[data-courier-route-composition]")?.className).toContain("bg-card");
    expect(document.querySelector("[data-courier-route-composition]")?.className).not.toMatch(/rounded|shadow|backdrop|border/);
    expect(document.querySelector("[data-courier-cli-registry]")?.className).not.toMatch(/rounded|bg-card/);
    for (const selector of ["[data-courier-cli-command-header]", "[data-courier-cli-options-header]", "[data-courier-cli-command]", "[data-builder-flag]", "[data-courier-cli-readout]"]) {
      expect(Array.from(document.querySelectorAll<HTMLElement>(selector)).every((node) => !/border-(?:b|y|border)/.test(node.className))).toBe(true);
    }
    expect(document.querySelector("header")?.className).not.toContain("fixed");
    expect(document.querySelector("header")?.className).toContain("bg-background");
    const activeLocal = document.querySelector<HTMLButtonElement>('[data-courier-route-composition] button[aria-label="Local"][aria-pressed="true"]');
    expect(activeLocal?.className).toContain("bg-[var(--terminal-fill-action)]");
    expect(activeLocal?.className).not.toContain("bg-background");
    expect(screen.getByText("COURIER CLI")).toBeTruthy();
    const localeButton = screen.getByRole("button", { name: /Language:/ });
    expect(localeButton.querySelector('[data-locale-icon="en"] svg')).toBeTruthy();
    fireEvent.click(localeButton);
    expect(localeButton.querySelectorAll('[data-locale-icon="ru"] i')).toHaveLength(3);
    expect(screen.getByRole("heading", { level: 1 }).textContent).toBe("Отсюда — куда угодно.");
    expect(document.body.textContent).not.toContain("Все связи взяты из опубликованного контракта CLI");
  });
});
