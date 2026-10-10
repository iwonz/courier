import * as React from "react";
import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { resetBrowserPreferenceController } from "@courier/ui";
import { DataApp } from "./app";
import { mountData } from "./main";

const response = (body: unknown, status = 200): Response => new Response(typeof body === "string" ? body : JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
const directory = (path = "", entries: unknown[] = []): object => ({ name: path ? "folder" : "root", path, size: 5, type: "directory", entries });

class BrowserUploadRequest {
  static instances: BrowserUploadRequest[] = [];
  static nextSendError: unknown;
  readonly upload: { onprogress: ((event: ProgressEvent) => void) | null } = { onprogress: null };
  readonly open = vi.fn();
  readonly setRequestHeader = vi.fn();
  readonly send = vi.fn(() => {
    if (BrowserUploadRequest.nextSendError !== undefined) {
      const error = BrowserUploadRequest.nextSendError;
      BrowserUploadRequest.nextSendError = undefined;
      throw error;
    }
  });
  withCredentials = false;
  status = 0;
  responseText = "";
  onload: (() => void) | null = null;
  onerror: (() => void) | null = null;
  onabort: (() => void) | null = null;
  ontimeout: (() => void) | null = null;

  constructor() { BrowserUploadRequest.instances.push(this); }

  progress(loaded: number, total: number): void {
    this.upload.onprogress?.({ lengthComputable: true, loaded, total } as ProgressEvent);
  }

  respond(status: number, body = "{}"): void {
    this.status = status;
    this.responseText = body;
    this.onload?.();
  }
}

beforeEach(() => {
  history.replaceState({}, "", "/d/token/");
  localStorage.clear();
  resetBrowserPreferenceController();
  BrowserUploadRequest.instances = [];
  BrowserUploadRequest.nextSendError = undefined;
  vi.stubGlobal("XMLHttpRequest", BrowserUploadRequest);
  vi.stubGlobal("matchMedia", vi.fn(() => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() })));
});

afterEach(() => {
  cleanup();
  resetBrowserPreferenceController();
  document.body.replaceChildren();
  document.head.querySelectorAll('link[rel="icon"]').forEach((node) => node.remove());
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe("delivery React application", () => {
  it("mounts into a root with the shared square Relay favicon and ignores no root", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => response({ name: "root", path: "", size: 0, type: "directory" })));
    mountData(null);
    const root = document.createElement("div");
    document.body.append(root);
    mountData(root);
    await waitFor(() => expect(root.textContent).toContain("root"));
    expect(document.head.querySelector<HTMLLinkElement>('link[rel="icon"]')?.href).toContain("courier-relay-pixel-mark-v3");
  });

  it("loads a directory, navigates into a folder and back, and localizes", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(response(directory("", [
        { name: "file.txt", size: 1, type: "file" },
        { name: "folder", size: 4, type: "directory" },
      ])))
      .mockResolvedValueOnce(response(directory("folder")))
      .mockResolvedValueOnce(response(directory()));
    vi.stubGlobal("fetch", fetchMock);
    render(<DataApp />);
    expect(document.querySelector<HTMLImageElement>('img[src*="courier-relay-pixel-delivery-v3"]')?.className).toContain("courier-pixel-image");
    await screen.findByText("file.txt");
    expect(screen.getByRole("link", { name: /Download file/ }).getAttribute("href")).toContain("file.txt");
    fireEvent.click(screen.getByRole("button", { name: "folder" }));
    await screen.findByRole("button", { name: "Parent directory" });
    fireEvent.click(screen.getByRole("button", { name: "Parent directory" }));
    await screen.findByText("No entries are available at this path.");
    fireEvent.click(screen.getByRole("button", { name: /Language:/ }));
    expect(screen.getByText(/нет доступных объектов/i)).toBeTruthy();
    expect(fetchMock).toHaveBeenCalledTimes(3);
  });

  it("isolates metadata before authentication, signs in, and retries", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(response("", 401))
      .mockResolvedValueOnce(response({ csrf: "csrf" }))
      .mockResolvedValueOnce(response(directory("", [{ name: "folder", size: 1, type: "directory" }])))
      .mockResolvedValueOnce(response("", 500))
      .mockResolvedValueOnce(response(directory()));
    vi.stubGlobal("fetch", fetchMock);
    render(<DataApp />);
    await screen.findByText("The route is unavailable or authorization is required. No delivery metadata was revealed.");
    expect(document.querySelector("[data-courier-auth-region]")?.className).toContain("bg-card");
    expect(document.querySelector("[data-courier-auth-region]")?.className).not.toMatch(/rounded|shadow|backdrop/);
    expect(document.body.textContent).not.toContain("folder");
    const password = screen.getByLabelText("Delivery password") as HTMLInputElement;
    fireEvent.change(password, { target: { value: "secret" } });
    fireEvent.submit(screen.getByRole("button", { name: "Verify access" }).closest("form")!);
    await screen.findByText("folder");
    expect(document.querySelector("[data-courier-manifest]")?.className).toContain("bg-card");
    expect(document.querySelector("[data-courier-manifest]")?.className).not.toMatch(/rounded|shadow|backdrop/);
    expect(password.value).toBe("");
    fireEvent.click(screen.getByRole("button", { name: "folder" }));
    await screen.findByRole("button", { name: "Retry connection" });
    fireEvent.click(screen.getByRole("button", { name: "Retry connection" }));
    await screen.findByText("No entries are available at this path.");
  });

  it("reports login failures including an empty password", async () => {
    const fetchMock = vi.fn().mockResolvedValue(response("", 500));
    vi.stubGlobal("fetch", fetchMock);
    render(<DataApp />);
    await screen.findByRole("button", { name: "Verify access" });
    screen.getByLabelText("Delivery password").removeAttribute("name");
    fireEvent.submit(screen.getByRole("button", { name: "Verify access" }).closest("form")!);
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(2));
    expect(document.body.textContent).toContain("The route is unavailable or authorization is required.");
  });

  it("queues multiple files, reports progress, continues after failure, and blocks overlap", async () => {
    const fetchMock = vi.fn().mockResolvedValue(response({ name: "Upload", path: "", size: 0, type: "upload" }));
    vi.stubGlobal("fetch", fetchMock);
    render(<DataApp />);
    const input = await waitFor(() => {
      const node = document.querySelector<HTMLInputElement>('input[type="file"]');
      expect(node).not.toBeNull();
      return node!;
    });
    expect(input.multiple).toBe(true);
    Object.defineProperty(input, "files", { configurable: true, value: null });
    fireEvent.change(input);
    expect(fetchMock).toHaveBeenCalledTimes(1);
    const first = new File(["abcd"], "first.txt");
    const second = new File(["xy"], "second.bin");
    Object.defineProperty(input, "files", { configurable: true, value: [first, second] });
    Object.defineProperty(input, "value", { configurable: true, writable: true, value: "selected" });
    fireEvent.change(input);
    await waitFor(() => expect(document.querySelectorAll("[data-courier-upload-item]")).toHaveLength(2));
    expect(document.querySelector('[data-courier-upload-item][data-state="uploading"]')?.textContent).toContain("first.txt");
    expect(document.querySelector('[data-courier-upload-item][data-state="queued"]')?.textContent).toContain("second.bin");
    expect(screen.getByText("4 Bytes")).toBeTruthy();
    expect(screen.getByText("2 Bytes")).toBeTruthy();
    expect(screen.getByRole("status").getAttribute("aria-live")).toBe("polite");
    expect(screen.getByRole("status").textContent).toContain("0/2");
    expect(input.disabled).toBe(true);
    expect(input.value).toBe("");
    expect(BrowserUploadRequest.instances).toHaveLength(1);
    BrowserUploadRequest.instances[0]!.progress(1, 4);
    await screen.findByText("25%");
    expect(screen.getByRole("progressbar", { name: "first.txt: Upload in progress" })).toBeTruthy();
    expect(screen.getByRole("progressbar", { name: "second.bin: Queued" })).toBeTruthy();
    fireEvent.change(input);
    expect(BrowserUploadRequest.instances).toHaveLength(1);
    BrowserUploadRequest.instances[0]!.respond(201);
    await waitFor(() => expect(BrowserUploadRequest.instances).toHaveLength(2));
    expect(document.querySelector('[data-courier-upload-item][data-state="success"]')?.textContent).toContain("first.txt");
    expect(document.querySelector('[data-courier-upload-item][data-state="uploading"]')?.textContent).toContain("second.bin");
    BrowserUploadRequest.instances[1]!.respond(409, JSON.stringify({ error: "destination already exists" }));
    await waitFor(() => expect(screen.getByRole("status").textContent).toContain("2/2"));
    expect(document.querySelector('[data-courier-upload-item][data-state="failure"]')?.textContent).toContain("second.bin");
    expect(screen.getByText(/Courier request failed \(409\): destination already exists/)).toBeTruthy();
    expect(input.disabled).toBe(false);
    expect(fetchMock).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByRole("button", { name: /Language:/ }));
    expect(screen.getByText("Очередь загрузки")).toBeTruthy();
    expect(screen.getByText(/Завершено: 2\/2/)).toBeTruthy();
  });

  it("keeps metadata visible and permits the same-file retry after each completed batch", async () => {
    const fetchMock = vi.fn().mockResolvedValue(response({ name: "Upload", path: "", size: 0, type: "upload" }));
    vi.stubGlobal("fetch", fetchMock);
    render(<DataApp />);
    const input = await waitFor(() => {
      const node = document.querySelector<HTMLInputElement>('input[type="file"]');
      expect(node).not.toBeNull();
      return node!;
    });
    const file = new File(["x"], "x.txt");
    Object.defineProperty(input, "files", { configurable: true, value: [file] });
    Object.defineProperty(input, "value", { configurable: true, writable: true, value: "selected" });
    fireEvent.change(input);
    await waitFor(() => expect(BrowserUploadRequest.instances).toHaveLength(1));
    BrowserUploadRequest.instances[0]!.respond(409, JSON.stringify({ error: "destination already exists" }));
    await screen.findByText("Upload failed");
    expect(screen.getByText(/Courier request failed \(409\): destination already exists/)).toBeTruthy();
    expect(document.querySelector("[data-courier-manifest]")).not.toBeNull();
    expect(input.value).toBe("");
    input.value = "selected";
    fireEvent.change(input);
    await waitFor(() => expect(BrowserUploadRequest.instances).toHaveLength(2));
    BrowserUploadRequest.instances[1]!.respond(201);
    await screen.findByText("File delivered");
    expect(input.value).toBe("");
    BrowserUploadRequest.nextSendError = "network unavailable";
    input.value = "selected";
    fireEvent.change(input);
    await waitFor(() => expect(document.querySelector("[data-courier-upload-error]")?.textContent).toContain("network unavailable"));
    expect(input.value).toBe("");
    expect(BrowserUploadRequest.instances).toHaveLength(3);
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  it("renders a shared file as one download", async () => {
    resetBrowserPreferenceController();
    vi.stubGlobal("fetch", vi.fn(async () => response({ name: "report.pdf", path: "", size: 10, type: "file" })));
    render(<DataApp />);
    expect(await screen.findByRole("link", { name: "Download file" })).toBeTruthy();
    expect(document.body.textContent).not.toContain("No entries are available at this path.");
  });
});
