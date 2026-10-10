import * as React from "react";
import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { resetBrowserPreferenceController } from "@courier/ui";
import { DataApp } from "./app";
import { mountData } from "./main";

const response = (body: unknown, status = 200): Response => new Response(typeof body === "string" ? body : JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
const directory = (path = "", entries: unknown[] = []): object => ({ name: path ? "folder" : "root", path, size: 5, type: "directory", entries });

beforeEach(() => {
  history.replaceState({}, "", "/d/token/");
  localStorage.clear();
  resetBrowserPreferenceController();
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

  it("shows pending and committed upload states, blocks overlap, and ignores an empty selection", async () => {
    let completeUpload!: (value: Response) => void;
    const pendingUpload = new Promise<Response>((resolve) => { completeUpload = resolve; });
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(response({ name: "Upload", path: "", size: 0, type: "upload" }))
      .mockReturnValueOnce(pendingUpload);
    vi.stubGlobal("fetch", fetchMock);
    render(<DataApp />);
    const input = await waitFor(() => {
      const node = document.querySelector<HTMLInputElement>('input[type="file"]');
      expect(node).not.toBeNull();
      return node!;
    });
    Object.defineProperty(input, "files", { configurable: true, value: { item: () => null } });
    fireEvent.change(input);
    expect(fetchMock).toHaveBeenCalledTimes(1);
    const file = new File(["x"], "x.txt");
    Object.defineProperty(input, "files", { configurable: true, value: { item: () => file } });
    Object.defineProperty(input, "value", { configurable: true, writable: true, value: "selected" });
    fireEvent.change(input);
    const pendingStatus = await screen.findByRole("status");
    expect(pendingStatus.getAttribute("aria-live")).toBe("polite");
    expect(pendingStatus.textContent).toContain("Upload in progress");
    expect(input.disabled).toBe(true);
    fireEvent.change(input);
    expect(fetchMock).toHaveBeenCalledTimes(2);
    completeUpload(response({}, 201));
    await waitFor(() => expect(screen.getByRole("status").textContent).toContain("File delivered"));
    expect(screen.getByText("Courier committed this file to the destination.")).toBeTruthy();
    expect(input.disabled).toBe(false);
    expect(input.value).toBe("");
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("keeps metadata visible, reports upload failures, and permits the same-file retry", async () => {
    const uploadFetch = vi.fn()
      .mockResolvedValueOnce(response({ name: "Upload", path: "", size: 0, type: "upload" }))
      .mockResolvedValueOnce(response({ error: "destination already exists" }, 409))
      .mockResolvedValueOnce(response({}, 201))
      .mockRejectedValueOnce("network unavailable");
    vi.stubGlobal("fetch", uploadFetch);
    render(<DataApp />);
    const input = await waitFor(() => {
      const node = document.querySelector<HTMLInputElement>('input[type="file"]');
      expect(node).not.toBeNull();
      return node!;
    });
    const file = new File(["x"], "x.txt");
    Object.defineProperty(input, "files", { configurable: true, value: { item: () => file } });
    Object.defineProperty(input, "value", { configurable: true, writable: true, value: "selected" });
    fireEvent.change(input);
    await screen.findByText("Upload failed");
    expect(screen.getByText(/Courier request failed \(409\): destination already exists/)).toBeTruthy();
    expect(document.querySelector("[data-courier-manifest]")).not.toBeNull();
    expect(input.value).toBe("");
    input.value = "selected";
    fireEvent.change(input);
    await screen.findByText("File delivered");
    expect(input.value).toBe("");
    input.value = "selected";
    fireEvent.change(input);
    await waitFor(() => expect(screen.getByRole("alert").textContent).toContain("network unavailable"));
    expect(input.value).toBe("");
    expect(uploadFetch).toHaveBeenCalledTimes(4);
  });

  it("renders a shared file as one download", async () => {
    resetBrowserPreferenceController();
    vi.stubGlobal("fetch", vi.fn(async () => response({ name: "report.pdf", path: "", size: 10, type: "file" })));
    render(<DataApp />);
    expect(await screen.findByRole("link", { name: "Download file" })).toBeTruthy();
    expect(document.body.textContent).not.toContain("No entries are available at this path.");
  });
});
