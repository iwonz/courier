import { beforeEach, describe, expect, it, vi } from "vitest";
import { CourierRequestError, apiURL, childPath, downloadURL, loadMetadata, login, parentPath, upload } from "./api";

class UploadRequest {
  readonly upload: { onprogress: ((event: ProgressEvent) => void) | null } = { onprogress: null };
  readonly open = vi.fn();
  readonly setRequestHeader = vi.fn();
  readonly send = vi.fn();
  withCredentials = false;
  status = 0;
  responseText = "";
  onload: (() => void) | null = null;
  onerror: (() => void) | null = null;
  onabort: (() => void) | null = null;
  ontimeout: (() => void) | null = null;

  factory = (): XMLHttpRequest => this as unknown as XMLHttpRequest;
}

describe("data API", () => {
  beforeEach(() => history.replaceState({}, "", "/d/token/"));

  it("builds resource-relative URLs", () => {
    expect(apiURL("meta")).toBe("/d/token/api/v1/meta");
    expect(apiURL("meta", "folder", "/d/token")).toBe("/d/token/api/v1/meta?path=folder");
    expect(downloadURL("file name")).toContain("path=file+name");
    expect(downloadURL("folder", true)).toContain("archive=tar.gz");
    expect(childPath("folder", "file")).toBe("folder/file");
    expect(childPath("", "file")).toBe("file");
    expect(parentPath("folder/file")).toBe("folder");
    expect(parentPath("file")).toBe("");
  });

  it("loads metadata and rejects HTTP failures", async () => {
    const fetcher = vi.fn(async () => new Response(JSON.stringify({ name: "root", path: "", size: 0, type: "directory" }), { status: 200 })) as typeof fetch;
    await expect(loadMetadata("", fetcher)).resolves.toMatchObject({ name: "root" });
    await expect(loadMetadata("", vi.fn(async () => new Response(JSON.stringify({ error: "request denied" }), { status: 403 })) as typeof fetch)).rejects.toEqual(new CourierRequestError(403, "request denied"));
  });

  it("bounds malformed and unsafe HTTP failure details", async () => {
    const bodies = [
      "not json",
      JSON.stringify({ error: 42 }),
      JSON.stringify({ error: "x".repeat(300) }),
      JSON.stringify({ error: "unsafe\nreason" }),
      JSON.stringify({ error: "x".repeat(1100) }),
    ];
    for (const body of bodies) {
      const fetcher = vi.fn(async () => new Response(body, { status: 500 })) as typeof fetch;
      await expect(loadMetadata("", fetcher)).rejects.toEqual(new CourierRequestError(500));
    }
    const unreadable = vi.fn(async () => ({ ok: false, status: 502, text: async () => Promise.reject(new Error("read")) }) as Response) as typeof fetch;
    await expect(loadMetadata("", unreadable)).rejects.toEqual(new CourierRequestError(502));
  });

  it("logs in and uploads with bounded request metadata and real progress", async () => {
    const fetcher = vi.fn(async () => new Response(JSON.stringify({ csrf: "token" }), { status: 200 })) as typeof fetch;
    await expect(login("secret", fetcher)).resolves.toEqual({ csrf: "token" });
    expect(fetcher).toHaveBeenCalledWith("/d/token/api/v1/session", expect.objectContaining({ method: "POST" }));

    const request = new UploadRequest();
    const progress = vi.fn();
    const result = upload(new File(["abc"], "a.txt"), "csrf", progress, request.factory);
    expect(request.open).toHaveBeenCalledWith("POST", "/d/token/api/v1/upload");
    expect(request.withCredentials).toBe(true);
    expect(request.setRequestHeader).toHaveBeenCalledWith("X-Courier-CSRF", "csrf");
    expect(request.setRequestHeader).toHaveBeenCalledWith("X-Courier-File-Size", "3");
    expect(request.send.mock.calls[0]?.[0]).toBeInstanceOf(FormData);
    expect(progress).toHaveBeenLastCalledWith({ loaded: 0, total: 3, percent: 0 });
    request.upload.onprogress?.({ lengthComputable: false, loaded: 2, total: 0 } as ProgressEvent);
    request.upload.onprogress?.({ lengthComputable: true, loaded: 2, total: 0 } as ProgressEvent);
    request.upload.onprogress?.({ lengthComputable: true, loaded: 1, total: 2 } as ProgressEvent);
    expect(progress).toHaveBeenLastCalledWith({ loaded: 2, total: 3, percent: 50 });
    request.upload.onprogress?.({ lengthComputable: true, loaded: 3, total: 2 } as ProgressEvent);
    expect(progress).toHaveBeenLastCalledWith({ loaded: 3, total: 3, percent: 100 });
    request.status = 201;
    request.onload?.();
    await expect(result).resolves.toBeUndefined();
    expect(progress).toHaveBeenLastCalledWith({ loaded: 3, total: 3, percent: 100 });
  });

  it("rejects guarded HTTP, network, abort, and timeout upload failures", async () => {
    const rejected = new UploadRequest();
    const failure = upload(new File(["x"], "x.txt"), "csrf", undefined, rejected.factory);
    rejected.status = 409;
    rejected.responseText = JSON.stringify({ error: "destination already exists" });
    rejected.onload?.();
    await expect(failure).rejects.toEqual(new CourierRequestError(409, "destination already exists"));

    for (const [event, reason] of [["onerror", "network error"], ["onabort", "upload aborted"], ["ontimeout", "upload timed out"]] as const) {
      const request = new UploadRequest();
      const result = upload(new File([], "empty.bin"), "csrf", undefined, request.factory);
      request[event]?.();
      await expect(result).rejects.toEqual(new CourierRequestError(0, reason));
    }
  });
});
