import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import { CourierDataApp } from "./app";

const flush = async (element: CourierDataApp): Promise<void> => {
  await vi.waitFor(async () => {
    await element.updateComplete;
    expect(element.shadowRoot?.textContent).not.toContain("Preparing the delivery route…");
  });
};

beforeEach(() => {
  history.replaceState({}, "", "/d/token/");
  vi.stubGlobal("matchMedia", () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() }));
});

afterEach(() => {
  document.body.replaceChildren();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

it("defines once, loads, localizes, and disconnects", async () => {
  const fetchMock = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ name: "root", path: "", size: 5, type: "directory", entries: [{ name: "file.txt", size: 1, type: "file" }, { name: "folder", size: 4, type: "directory" }] }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ name: "folder", path: "folder", size: 4, type: "directory", entries: [] }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ name: "root", path: "", size: 5, type: "directory", entries: [] }), { status: 200 }));
  await import("./main");
  vi.stubGlobal("fetch", fetchMock);
  const element = new CourierDataApp();
  document.body.append(element);
  await flush(element);
  expect(element.shadowRoot?.textContent).toContain("Download as archive");
  expect([...element.shadowRoot!.querySelectorAll("a")].some((link) => link.getAttribute("href")?.includes("file.txt"))).toBe(true);
  const folder = [...element.shadowRoot!.querySelectorAll("button.link")].find((button) => button.textContent === "folder") as HTMLButtonElement;
  folder.click();
  await vi.waitFor(async () => {
    await element.updateComplete;
    expect(element.shadowRoot?.textContent).toContain("Parent directory");
  });
  const up = [...element.shadowRoot!.querySelectorAll("button.link")].find((button) => button.textContent === "Parent directory") as HTMLButtonElement;
  up.click();
  await vi.waitFor(async () => {
    await element.updateComplete;
    expect(fetchMock).toHaveBeenCalledTimes(3);
    expect(element.shadowRoot?.textContent).toContain("No entries are available");
  });
  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(element.shadowRoot?.textContent).toContain("нет доступных объектов");
  element.remove();
  new CourierDataApp().disconnectedCallback();
});

it("renders login failure, signs in, and retries", async () => {
  const fetchMock = vi.fn()
    .mockResolvedValueOnce(new Response("", { status: 401 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ csrf: "csrf" }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ name: "root", path: "", size: 1, type: "directory", entries: [{ name: "folder", size: 1, type: "directory" }] }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ name: "root", path: "", size: 0, type: "directory", entries: [] }), { status: 200 }));
  vi.stubGlobal("fetch", fetchMock);
  const element = new CourierDataApp();
  document.body.append(element);
  await flush(element);
  expect(element.shadowRoot?.textContent).toContain("authorization is required");
  const accessArt = element.shadowRoot?.querySelector("courier-scene") as HTMLElement & { source: string; mobileSource: string };
  expect(accessArt.source).toContain("delivery-access-wide");
  expect(accessArt.mobileSource).toContain("delivery-access-mobile");
  expect(element.shadowRoot?.querySelector("courier-workbench.access")).not.toBeNull();
  expect(element.shadowRoot?.textContent).not.toContain("report.pdf");
  const form = element.shadowRoot?.querySelector("form") as HTMLFormElement;
  (form.querySelector("input") as HTMLInputElement).value = "secret";
  await element.signIn({ preventDefault: vi.fn(), currentTarget: form } as unknown as SubmitEvent);
  await element.updateComplete;
  expect(element.shadowRoot?.querySelector("a")?.href).toContain("archive=tar.gz");
  await element.refresh();
  expect(element.shadowRoot?.textContent).toContain("No entries are available");
});

it("handles login and upload errors plus empty selections", async () => {
  vi.stubGlobal("fetch", vi.fn(async () => new Response("", { status: 500 })));
  const element = new CourierDataApp();
  await element.signIn({ preventDefault: vi.fn(), currentTarget: document.createElement("form") } as unknown as SubmitEvent);
  const input = document.createElement("input");
  Object.defineProperty(input, "files", { configurable: true, value: { item: () => null } });
  await element.sendFile({ currentTarget: input } as unknown as Event);
  const file = new File(["x"], "x.txt");
  Object.defineProperty(input, "files", { configurable: true, value: { item: () => file } });
  await element.sendFile({ currentTarget: input } as unknown as Event);
  expect(element.render()).toBeTruthy();
});

it("renders an upload-only delivery", async () => {
  vi.stubGlobal("fetch", vi.fn(async () => new Response(JSON.stringify({ name: "Upload", path: "", size: 0, type: "upload" }), { status: 200 })));
  const element = new CourierDataApp();
  document.body.append(element);
  await flush(element);
  expect(element.shadowRoot?.querySelector('input[type="file"]')).not.toBeNull();
  expect(element.shadowRoot?.querySelector(".courier-file-action courier-icon")).not.toBeNull();
  expect(element.shadowRoot?.textContent).not.toContain("This directory is empty");
});

it("renders a shared file as one download", async () => {
  vi.stubGlobal("fetch", vi.fn(async () => new Response(JSON.stringify({ name: "report.pdf", path: "", size: 10, type: "file" }), { status: 200 })));
  const element = new CourierDataApp();
  document.body.append(element);
  await flush(element);
  expect(element.shadowRoot?.querySelector("a")?.textContent).toBe("Download file");
  expect(element.shadowRoot?.textContent).not.toContain("This directory is empty");
});

it("uploads a selected file and refreshes", async () => {
  const fetchMock = vi.fn()
    .mockResolvedValueOnce(new Response("{}", { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ name: "root", path: "", size: 0, type: "directory", entries: [] }), { status: 200 }));
  vi.stubGlobal("fetch", fetchMock);
  const element = new CourierDataApp();
  const input = document.createElement("input");
  input.value = "";
  Object.defineProperty(input, "files", { value: { item: () => new File(["x"], "x.txt") } });
  await element.sendFile({ currentTarget: input } as unknown as Event);
  expect(fetchMock).toHaveBeenCalledTimes(2);
});

it("does not redefine the application element", async () => {
  vi.resetModules();
  await import("./main");
  expect(customElements.get("courier-data-app")).toBe(CourierDataApp);
});
