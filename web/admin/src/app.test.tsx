import * as React from "react";
import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { resetBrowserPreferenceController } from "@courier/ui";
import { AdminApp, reconcileDeliverySelection } from "./app";
import type { Delivery, Snapshot } from "./api";
import { mountAdmin } from "./main";

const policy = {
  version: 1, auth: "none", authAttempts: 5, authFailAction: "ban",
  deliveryLimit: { unlimited: true, value: 0 }, allowIp: [],
  maxFileSize: { unlimited: false, value: 10 }, maxExtractedSize: { unlimited: false, value: 100 },
  uploadRate: { unlimited: true, value: 0 }, downloadRate: { unlimited: true, value: 0 }, noUi: false,
} as const;

const firstDelivery: Delivery = {
  id: "delivery-a", route: "path-to-web", source: "./data", destination: "web://", state: "active", policy,
  counters: { read: 3, sent: 2, confirmed: 1 }, createdAt: "2026-09-14T00:00:00Z", updatedAt: "2026-09-14T00:00:00Z",
};
const secondDelivery: Delivery = {
  ...firstDelivery, id: "delivery-b", source: "", destination: "", policy: { ...policy, auth: "password", authFailAction: "stop", noUi: true },
};
const snapshot: Snapshot = { servers: [
  { id: "server-a", bind: "127.0.0.1:8080", processId: 42, state: "active", status: "live", startedAt: "", updatedAt: "", deliveries: [firstDelivery, secondDelivery] },
  { id: "server-b", bind: "127.0.0.1:8081", processId: 43, state: "failed", status: "unreachable", startedAt: "", updatedAt: "", deliveries: [] },
] };

const response = (body: unknown, status = 200): Response => new Response(status === 204 ? null : JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
let eventListener: EventListener | undefined;
const eventClose = vi.fn();

beforeEach(() => {
  history.replaceState({}, "", "/");
  localStorage.clear();
  resetBrowserPreferenceController();
  eventListener = undefined;
  eventClose.mockClear();
  vi.stubGlobal("matchMedia", vi.fn(() => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })));
  vi.stubGlobal("ResizeObserver", class { observe(): void {} unobserve(): void {} disconnect(): void {} });
  class MockEventSource {
    close = eventClose;
    addEventListener(_name: string, listener: EventListenerOrEventListenerObject): void { eventListener = listener as EventListener; }
  }
  vi.stubGlobal("EventSource", MockEventSource);
});

afterEach(() => {
  cleanup();
  resetBrowserPreferenceController();
  document.body.replaceChildren();
  document.head.querySelectorAll('link[rel="icon"]').forEach((node) => node.remove());
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe("selection reconciliation", () => {
  it("preserves valid delivery IDs and falls back deterministically", () => {
    expect(reconcileDeliverySelection(snapshot, "delivery-b")).toBe("delivery-b");
    expect(reconcileDeliverySelection(snapshot, "missing")).toBe("delivery-a");
    expect(reconcileDeliverySelection({ servers: [] }, "delivery-a")).toBeUndefined();
  });
});

describe("administration React application", () => {
  it("mounts with the shared square Relay favicon and ignores no root", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => response({ servers: [] })));
    mountAdmin(null);
    const root = document.createElement("div");
    document.body.append(root);
    mountAdmin(root);
    await waitFor(() => expect(root.textContent).toContain("No live Courier servers"));
    expect(document.head.querySelector<HTMLLinkElement>('link[rel="icon"]')?.href).toContain("courier-relay-tech-v1");
  });

  it("loads, selects, preserves SSE selection, falls back, localizes, and unsubscribes", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => response(snapshot)));
    const view = render(<AdminApp />);
    await screen.findByText("server-a");
    expect(screen.getByText("Unreachable")).toBeTruthy();
    fireEvent.click(screen.getByRole("button", { name: "Refresh registry" }));
    await screen.findByText("server-b");
    fireEvent.click(screen.getByRole("button", { name: /path-to-webdelivery-b/ }));
    expect(screen.getAllByText("Unavailable").length).toBeGreaterThan(0);
    eventListener!(new MessageEvent("snapshot", { data: JSON.stringify(snapshot) }));
    expect(screen.getByRole("button", { name: /path-to-webdelivery-b/ }).getAttribute("aria-pressed")).toBe("true");
    eventListener!(new MessageEvent("snapshot", { data: JSON.stringify({ servers: [] }) }));
    await screen.findByText("No live Courier servers are registered.");
    fireEvent.click(screen.getByRole("button", { name: /Language:/ }));
    expect(screen.getByRole("heading", { level: 1 }).textContent).toBe("Управление доставками");
    view.unmount();
    expect(eventClose).toHaveBeenCalled();
  });

  it("stops a server and delivery and refreshes after each action", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(response(snapshot))
      .mockResolvedValueOnce(response({}, 204))
      .mockResolvedValueOnce(response(snapshot))
      .mockResolvedValueOnce(response({}, 204))
      .mockResolvedValueOnce(response(snapshot));
    vi.stubGlobal("fetch", fetchMock);
    render(<AdminApp />);
    await screen.findByText("server-a");
    fireEvent.click(screen.getAllByRole("button", { name: "Stop server" })[0]!);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith("/api/v1/servers/server-a/stop", expect.anything()));
    fireEvent.click(screen.getByRole("button", { name: "Stop delivery" }));
    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith("/api/v1/deliveries/delivery-a/stop", expect.anything()));
  });

  it("saves an optimistic policy and refreshes", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(response(snapshot))
      .mockResolvedValueOnce(response({}, 204))
      .mockResolvedValueOnce(response(snapshot));
    vi.stubGlobal("fetch", fetchMock);
    render(<AdminApp />);
    const form = await waitFor(() => {
      const node = document.querySelector<HTMLFormElement>("form");
      expect(node).not.toBeNull();
      return node!;
    });
    fireEvent.change(form.elements.namedItem("attempts")!, { target: { value: "7" } });
    fireEvent.click(form.elements.namedItem("noUi") as HTMLElement);
    fireEvent.submit(form);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(3));
    const body = String(fetchMock.mock.calls[1]![1]?.body);
    expect(body).toContain('"version":2');
    expect(body).toContain('"authAttempts":7');
    expect(body).toContain('"noUi":true');
  });

  it("shows a policy conflict and clears it through refresh", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(response(snapshot))
      .mockResolvedValueOnce(response({}, 409))
      .mockResolvedValueOnce(response(snapshot));
    vi.stubGlobal("fetch", fetchMock);
    render(<AdminApp />);
    const form = await waitFor(() => {
      const node = document.querySelector<HTMLFormElement>("form");
      expect(node).not.toBeNull();
      return node!;
    });
    fireEvent.submit(form);
    await screen.findByText("This delivery policy changed elsewhere. Refresh before editing again.");
    fireEvent.click(screen.getAllByRole("button", { name: "Refresh registry" })[1]!);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(3));
    expect(screen.queryByText("This delivery policy changed elsewhere. Refresh before editing again.")).toBeNull();
  });

  it("reports registry, stop, and non-conflict policy failures and retries", async () => {
    const fetchMock = vi.fn()
      .mockRejectedValueOnce(new Error("offline"))
      .mockResolvedValueOnce(response(snapshot))
      .mockRejectedValueOnce("stop failed")
      .mockRejectedValueOnce(new Error("save failed"));
    vi.stubGlobal("fetch", fetchMock);
    render(<AdminApp />);
    await screen.findByText("Administration data is temporarily unavailable. Retry the local connection.");
    fireEvent.click(screen.getByRole("button", { name: "Retry connection" }));
    await screen.findByText("server-a");
    fireEvent.click(screen.getAllByRole("button", { name: "Stop server" })[0]!);
    await screen.findByText("Administration data is temporarily unavailable. Retry the local connection.");
    fireEvent.submit(document.querySelector("form")!);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(4));
    expect(screen.queryByText("This delivery policy changed elsewhere. Refresh before editing again.")).toBeNull();
  });

  it("renders loading and a server with no selectable delivery", async () => {
    let resolveFetch: ((value: Response) => void) | undefined;
    vi.stubGlobal("fetch", vi.fn(() => new Promise<Response>((resolve) => { resolveFetch = resolve; })));
    const view = render(<AdminApp />);
    expect(screen.getByText("Checking the live registry…")).toBeTruthy();
    resolveFetch!(response({ servers: [snapshot.servers[1]!] }));
    await screen.findByText("server-b");
    expect(screen.getByText("No live Courier servers are registered.")).toBeTruthy();
    view.unmount();
  });
});
