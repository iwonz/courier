import { beforeEach, describe, expect, it, vi } from "vitest";
import { CourierRequestError, apiURL, childPath, downloadURL, loadMetadata, login, parentPath, upload } from "./api";

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

  it("logs in and uploads with bounded request metadata", async () => {
    const fetcher = vi.fn(async () => new Response(JSON.stringify({ csrf: "token" }), { status: 200 })) as typeof fetch;
    await expect(login("secret", fetcher)).resolves.toEqual({ csrf: "token" });
    expect(fetcher).toHaveBeenCalledWith("/d/token/api/v1/session", expect.objectContaining({ method: "POST" }));
    const uploadFetch = vi.fn(async () => new Response("{}", { status: 200 })) as typeof fetch;
    await upload(new File(["abc"], "a.txt"), "csrf", uploadFetch);
    expect(uploadFetch).toHaveBeenCalledWith("/d/token/api/v1/upload", expect.objectContaining({ headers: expect.objectContaining({ "X-Courier-CSRF": "csrf", "X-Courier-File-Size": "3" }) }));
  });
});
