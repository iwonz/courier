import { beforeEach, describe, expect, it, vi } from "vitest";
import { apiURL, childPath, downloadURL, loadMetadata, login, parentPath, upload } from "./api";

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
    await expect(loadMetadata("", vi.fn(async () => new Response("", { status: 403 })) as typeof fetch)).rejects.toThrow("403");
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
