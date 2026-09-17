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
    expect(document.head.querySelector<HTMLLinkElement>('link[rel="icon"]')?.href).toContain("courier-relay-tech-v1");
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
    expect(document.body.textContent).not.toContain("folder");
    const password = screen.getByLabelText("Delivery password") as HTMLInputElement;
    fireEvent.change(password, { target: { value: "secret" } });
    fireEvent.submit(screen.getByRole("button", { name: "Verify access" }).closest("form")!);
    await screen.findByText("folder");
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

  it("uploads one selected file, refreshes, and ignores an empty selection", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(response({ name: "Upload", path: "", size: 0, type: "upload" }))
      .mockResolvedValueOnce(response({}))
      .mockResolvedValueOnce(response(directory()));
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
    Object.defineProperty(input, "files", { configurable: true, value: { item: () => new File(["x"], "x.txt") } });
    fireEvent.change(input);
    await screen.findByText("No entries are available at this path.");
    expect(fetchMock).toHaveBeenCalledTimes(3);
  });

  it("reports an upload failure and renders a shared file as one download", async () => {
    const uploadFetch = vi.fn()
      .mockResolvedValueOnce(response({ name: "Upload", path: "", size: 0, type: "upload" }))
      .mockResolvedValueOnce(response("", 500));
    vi.stubGlobal("fetch", uploadFetch);
    const first = render(<DataApp />);
    const input = await waitFor(() => {
      const node = document.querySelector<HTMLInputElement>('input[type="file"]');
      expect(node).not.toBeNull();
      return node!;
    });
    Object.defineProperty(input, "files", { configurable: true, value: { item: () => new File(["x"], "x.txt") } });
    fireEvent.change(input);
    await waitFor(() => expect(uploadFetch).toHaveBeenCalledTimes(2));
    first.unmount();
    resetBrowserPreferenceController();
    vi.stubGlobal("fetch", vi.fn(async () => response({ name: "report.pdf", path: "", size: 10, type: "file" })));
    render(<DataApp />);
    expect(await screen.findByRole("link", { name: "Download file" })).toBeTruthy();
    expect(document.body.textContent).not.toContain("No entries are available at this path.");
  });
});
