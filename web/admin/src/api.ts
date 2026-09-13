export interface Limit {
  unlimited: boolean;
  value: number;
}

export interface Policy {
  version: number;
  auth: "none" | "basic" | "password";
  authAttempts: number;
  authFailAction: "ban" | "stop";
  deliveryLimit: Limit;
  allowIp: string[];
  maxFileSize: Limit;
  maxExtractedSize: Limit;
  uploadRate: Limit;
  downloadRate: Limit;
  noUi: boolean;
}

export interface Delivery {
  id: string;
  route: string;
  source: string;
  destination: string;
  state: string;
  policy: Policy;
  counters: { read: number; sent: number; confirmed: number };
  createdAt: string;
  updatedAt: string;
}

export interface Server {
  id: string;
  bind: string;
  processId: number;
  state: string;
  status: "live" | "unreachable";
  startedAt: string;
  updatedAt: string;
  deliveries: Delivery[];
}

export interface Snapshot {
  servers: Server[];
}

export type Fetcher = typeof fetch;
export type EventSourceFactory = (url: string) => Pick<EventSource, "addEventListener" | "close">;

export function apiURL(path: string, pathname = globalThis.location.pathname): string {
  const root = pathname.endsWith("/") ? pathname : `${pathname}/`;
  return `${root}api/v1/${path}`;
}

async function decode<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = new Error(`Courier administration request failed (${response.status})`);
    error.name = response.status === 409 ? "ConflictError" : "RequestError";
    throw error;
  }
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}

export async function loadServers(fetcher: Fetcher = globalThis.fetch): Promise<Snapshot> {
  return decode<Snapshot>(await fetcher(apiURL("servers"), { credentials: "same-origin" }));
}

export async function savePolicy(delivery: Delivery, policy: Policy, fetcher: Fetcher = globalThis.fetch): Promise<void> {
  await decode(await fetcher(apiURL(`deliveries/${delivery.id}/policy`), {
    method: "PUT",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ expectedVersion: delivery.policy.version, policy }),
  }));
}

export async function stopTarget(kind: "servers" | "deliveries", id: string, fetcher: Fetcher = globalThis.fetch): Promise<void> {
  await decode(await fetcher(apiURL(`${kind}/${id}/stop`), { method: "POST", credentials: "same-origin" }));
}

export function subscribeSnapshots(onSnapshot: (snapshot: Snapshot) => void, create: EventSourceFactory = (url) => new EventSource(url)): () => void {
  const source = create(apiURL("events"));
  source.addEventListener("snapshot", (event) => onSnapshot(JSON.parse((event as MessageEvent<string>).data) as Snapshot));
  return () => source.close();
}
