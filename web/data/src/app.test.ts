import { afterEach, beforeEach, expect, it, vi } from "vitest";
import "./main";
import { CourierDataApp } from "./app";

const flush = async (element: CourierDataApp): Promise<void> => {
  await vi.waitFor(async () => {
    await element.updateComplete;
    expect(element.shadowRoot?.textContent).not.toContain("Loading delivery…");
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
  expect(element.shadowRoot?.textContent).toContain("Download archive");
  expect([...element.shadowRoot!.querySelectorAll("a")].some((link) => link.getAttribute("href")?.includes("file.txt"))).toBe(true);
  const folder = [...element.shadowRoot!.querySelectorAll("button.link")].find((button) => button.textContent === "folder") as HTMLButtonElement;
  folder.click();
  await vi.waitFor(async () => {
    await element.updateComplete;
    expect(element.shadowRoot?.textContent).toContain("Up");
  });
  const up = [...element.shadowRoot!.querySelectorAll("button.link")].find((button) => button.textContent === "Up") as HTMLButtonElement;
  up.click();
  await vi.waitFor(async () => {
    await element.updateComplete;
    expect(fetchMock).toHaveBeenCalledTimes(3);
    expect(element.shadowRoot?.textContent).toContain("This directory is empty");
  });
  element.setLocale(new CustomEvent("courier-locale", { detail: "ru" }));
  await element.updateComplete;
  expect(element.shadowRoot?.textContent).toContain("Каталог пуст");
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
  const form = element.shadowRoot?.querySelector("form") as HTMLFormElement;
  (form.querySelector("input") as HTMLInputElement).value = "secret";
  await element.signIn({ preventDefault: vi.fn(), currentTarget: form } as unknown as SubmitEvent);
  await element.updateComplete;
  expect(element.shadowRoot?.querySelector("a")?.href).toContain("archive=tar.gz");
  await element.refresh();
  expect(element.shadowRoot?.textContent).toContain("empty");
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
  expect(element.shadowRoot?.textContent).not.toContain("This directory is empty");
});

it("renders a shared file as one download", async () => {
  vi.stubGlobal("fetch", vi.fn(async () => new Response(JSON.stringify({ name: "report.pdf", path: "", size: 10, type: "file" }), { status: 200 })));
  const element = new CourierDataApp();
  document.body.append(element);
  await flush(element);
  expect(element.shadowRoot?.querySelector("a")?.textContent).toBe("Download");
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
