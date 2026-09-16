import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import { CourierAdminApp } from "./app";
import type { Delivery, Snapshot } from "./api";

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

let eventListener: EventListener | undefined;
const eventClose = vi.fn();

beforeEach(() => {
  history.replaceState({}, "", "/");
  eventListener = undefined;
  eventClose.mockClear();
  vi.stubGlobal("matchMedia", () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() }));
  class MockEventSource {
    close = eventClose;
    addEventListener(_name: string, listener: EventListenerOrEventListenerObject): void { eventListener = listener as EventListener; }
  }
  vi.stubGlobal("EventSource", MockEventSource);
});

afterEach(() => {
  document.body.replaceChildren();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

async function mount(fetcher = vi.fn(async () => new Response(JSON.stringify(snapshot), { status: 200 }))): Promise<CourierAdminApp> {
  vi.stubGlobal("fetch", fetcher);
  const element = new CourierAdminApp();
  document.body.append(element);
  await vi.waitFor(async () => {
    await element.updateComplete;
    expect(element.shadowRoot?.textContent).toContain("server-a");
  });
  return element;
}

it("loads, renders, receives events, localizes, and disconnects", async () => {
  await import("./main");
  const element = await mount();
  expect(element.shadowRoot?.textContent).toContain("server-a");
  expect(element.shadowRoot?.textContent).toContain("Unreachable");
  expect(element.shadowRoot?.textContent).toContain("Unavailable");
  const forms = [...element.shadowRoot!.querySelectorAll("form")];
  expect((forms[0]!.elements.namedItem("auth") as HTMLSelectElement).value).toBe("none");
  expect((forms[0]!.elements.namedItem("failAction") as HTMLSelectElement).value).toBe("ban");
  expect((forms[1]!.elements.namedItem("auth") as HTMLSelectElement).value).toBe("password");
  expect((forms[1]!.elements.namedItem("failAction") as HTMLSelectElement).value).toBe("stop");
  expect((forms[1]!.elements.namedItem("noUi") as HTMLInputElement).checked).toBe(true);
  eventListener!(new MessageEvent("snapshot", { data: JSON.stringify({ servers: [] }) }));
  await element.updateComplete;
  expect(element.shadowRoot?.textContent).toContain("No live Courier servers");
  const operationsArt = element.shadowRoot?.querySelector("courier-scene") as HTMLElement & { source: string; mobileSource: string };
  expect(operationsArt.source).toContain("relay-terminal-admin-wide");
  expect(operationsArt.mobileSource).toContain("relay-terminal-admin-mobile");
  expect(element.shadowRoot?.querySelector("courier-terminal.registry-terminal")).not.toBeNull();
  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(element.shadowRoot?.textContent).toContain("Активные серверы Courier");
  element.remove();
  expect(eventClose).toHaveBeenCalledOnce();
  new CourierAdminApp().disconnectedCallback();
});

it("runs server and delivery stop actions", async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify(snapshot), { status: 200 }))
    .mockResolvedValueOnce(new Response(null, { status: 204 }))
    .mockResolvedValueOnce(new Response(JSON.stringify(snapshot), { status: 200 }))
    .mockResolvedValueOnce(new Response(null, { status: 204 }))
    .mockResolvedValueOnce(new Response(JSON.stringify(snapshot), { status: 200 }));
  const element = await mount(fetcher);
  const buttons = [...element.shadowRoot!.querySelectorAll("courier-button")];
  const serverStop = buttons.find((button) => button.textContent === "Stop server") as HTMLElement;
  serverStop.click();
  await vi.waitFor(() => expect(fetcher).toHaveBeenCalledWith("/api/v1/servers/server-a/stop", expect.anything()));
  const deliveryStop = buttons.find((button) => button.textContent === "Stop delivery") as HTMLElement;
  deliveryStop.click();
  await vi.waitFor(() => expect(fetcher).toHaveBeenCalledWith("/api/v1/deliveries/delivery-a/stop", expect.anything()));
});

it("saves an optimistic policy and refreshes", async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify(snapshot), { status: 200 }))
    .mockResolvedValueOnce(new Response(null, { status: 204 }))
    .mockResolvedValueOnce(new Response(JSON.stringify(snapshot), { status: 200 }));
  const element = await mount(fetcher);
  const form = element.shadowRoot!.querySelector("form") as HTMLFormElement;
  (form.elements.namedItem("auth") as HTMLSelectElement).value = "basic";
  (form.elements.namedItem("attempts") as HTMLInputElement).value = "7";
  (form.elements.namedItem("failAction") as HTMLSelectElement).value = "stop";
  (form.elements.namedItem("noUi") as HTMLInputElement).checked = true;
  form.dispatchEvent(new SubmitEvent("submit", { bubbles: true, cancelable: true }));
  await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(3));
  expect(fetcher.mock.calls[1][1]?.body).toContain('"version":2');
  expect(fetcher.mock.calls[1][1]?.body).toContain('"authAttempts":7');
});

it("shows conflict and ordinary failures, then retries", async () => {
  const element = new CourierAdminApp();
  vi.stubGlobal("fetch", vi.fn(async () => new Response("", { status: 500 })));
  await element.refresh();
  expect(element.render()).toBeTruthy();
  await element.stop("servers", "server-a");
  expect(element.render()).toBeTruthy();

  vi.stubGlobal("fetch", vi.fn(async () => new Response("", { status: 409 })));
  const form = document.createElement("form");
  for (const [name, value] of [["auth", "none"], ["attempts", "5"], ["failAction", "ban"]]) {
    const input = document.createElement("input");
    input.name = name;
    input.value = value;
    form.append(input);
  }
  await element.save({ preventDefault: vi.fn(), currentTarget: form } as unknown as SubmitEvent, firstDelivery);
  expect(element.render()).toBeTruthy();

  vi.stubGlobal("fetch", vi.fn(async () => new Response("", { status: 500 })));
  await element.save({ preventDefault: vi.fn(), currentTarget: form } as unknown as SubmitEvent, firstDelivery);
  expect(element.render()).toBeTruthy();
});

it("renders initial loading state and does not redefine the element", async () => {
  expect(new CourierAdminApp().render()).toBeTruthy();
  vi.resetModules();
  await import("./main");
  expect(customElements.get("courier-admin-app")).toBe(CourierAdminApp);
});
