export interface Entry {
  name: string;
  size: number;
  type: "file" | "directory" | "upload";
}

export interface Metadata extends Entry {
  path: string;
  entries?: Entry[];
}

export interface LoginResult {
  csrf: string;
}

export type Fetcher = typeof fetch;

export function apiURL(endpoint: string, path = "", pathname = globalThis.location.pathname): string {
  const root = pathname.endsWith("/") ? pathname : `${pathname}/`;
  const url = new URL(`api/v1/${endpoint}`, globalThis.location.origin);
  url.pathname = `${root}api/v1/${endpoint}`;
  if (path) {
    url.searchParams.set("path", path);
  }
  return `${url.pathname}${url.search}`;
}

async function decode<T>(response: Response): Promise<T> {
  if (!response.ok) {
    throw new Error(`Courier request failed (${response.status})`);
  }
  return response.json() as Promise<T>;
}

export async function loadMetadata(path = "", fetcher: Fetcher = globalThis.fetch): Promise<Metadata> {
  return decode<Metadata>(await fetcher(apiURL("meta", path), { credentials: "same-origin" }));
}

export async function login(password: string, fetcher: Fetcher = globalThis.fetch): Promise<LoginResult> {
  return decode<LoginResult>(await fetcher(apiURL("session"), {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ password }),
  }));
}

export async function upload(file: File, csrf: string, fetcher: Fetcher = globalThis.fetch): Promise<void> {
  const body = new FormData();
  body.append("file", file);
  await decode(await fetcher(apiURL("upload"), {
    method: "POST",
    credentials: "same-origin",
    headers: { "X-Courier-CSRF": csrf, "X-Courier-File-Size": String(file.size) },
    body,
  }));
}

export function downloadURL(path: string, archive = false): string {
  const url = new URL(apiURL("download", path), globalThis.location.origin);
  if (archive) {
    url.searchParams.set("archive", "tar.gz");
  }
  return `${url.pathname}${url.search}`;
}

export function childPath(parent: string, child: string): string {
  return parent ? `${parent}/${child}` : child;
}

export function parentPath(path: string): string {
  const separator = path.lastIndexOf("/");
  return separator < 0 ? "" : path.slice(0, separator);
}
