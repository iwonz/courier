import { beforeEach, expect, it, vi } from "vitest";
import { apiURL, loadServers, savePolicy, stopTarget, subscribeSnapshots, type Delivery, type Snapshot } from "./api";

const policy = {
  version: 1, auth: "none", authAttempts: 5, authFailAction: "ban",
  deliveryLimit: { unlimited: true, value: 0 }, allowIp: [],
  maxFileSize: { unlimited: false, value: 10 }, maxExtractedSize: { unlimited: false, value: 100 },
  uploadRate: { unlimited: true, value: 0 }, downloadRate: { unlimited: true, value: 0 }, noUi: false,
} as const;

const delivery = { id: "delivery", policy } as unknown as Delivery;

beforeEach(() => history.replaceState({}, "", "/"));

it("builds administration-relative URLs", () => {
  expect(apiURL("servers")).toBe("/api/v1/servers");
  expect(apiURL("events", "/admin")).toBe("/admin/api/v1/events");
});

it("loads snapshots and classifies failures", async () => {
  const fetcher = vi.fn(async () => new Response(JSON.stringify({ servers: [] }), { status: 200 })) as typeof fetch;
  await expect(loadServers(fetcher)).resolves.toEqual({ servers: [] });
  await expect(loadServers(vi.fn(async () => new Response("", { status: 500 })) as typeof fetch)).rejects.toMatchObject({ name: "RequestError" });
  await expect(loadServers(vi.fn(async () => new Response("", { status: 409 })) as typeof fetch)).rejects.toMatchObject({ name: "ConflictError" });
});

it("updates policies and stops both target kinds", async () => {
  const fetcher = vi.fn(async () => new Response(null, { status: 204 })) as typeof fetch;
  await savePolicy(delivery, { ...policy, version: 2 }, fetcher);
  await stopTarget("deliveries", "delivery", fetcher);
  await stopTarget("servers", "server", fetcher);
  expect(fetcher).toHaveBeenNthCalledWith(1, "/api/v1/deliveries/delivery/policy", expect.objectContaining({ method: "PUT", body: expect.stringContaining('"expectedVersion":1') }));
  expect(fetcher).toHaveBeenNthCalledWith(3, "/api/v1/servers/server/stop", expect.objectContaining({ method: "POST" }));
});

it("subscribes, decodes snapshots, and closes event sources", () => {
  let listener: EventListener | undefined;
  const close = vi.fn();
  const create = vi.fn(() => ({ addEventListener: (_name: string, value: EventListenerOrEventListenerObject) => { listener = value as EventListener; }, close }));
  const receive = vi.fn();
  const unsubscribe = subscribeSnapshots(receive, create);
  const snapshot: Snapshot = { servers: [] };
  listener!(new MessageEvent("snapshot", { data: JSON.stringify(snapshot) }));
  expect(receive).toHaveBeenCalledWith(snapshot);
  unsubscribe();
  expect(close).toHaveBeenCalledOnce();

  const addEventListener = vi.fn();
  const defaultClose = vi.fn();
  class MockEventSource {
    addEventListener = addEventListener;
    close = defaultClose;
  }
  vi.stubGlobal("EventSource", MockEventSource);
  subscribeSnapshots(receive)();
  expect(defaultClose).toHaveBeenCalledOnce();
  vi.unstubAllGlobals();
});
