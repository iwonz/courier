import { test, expect, type Page } from "@playwright/test";
import { spawn, execFile, execFileSync, type ChildProcessWithoutNullStreams } from "node:child_process";
import { createServer, type Server } from "node:http";
import { mkdtemp, mkdir, readFile, readdir, writeFile, lstat } from "node:fs/promises";
import { basename, join } from "node:path";
import { promisify } from "node:util";

const executeFile = promisify(execFile);
const binary = process.env.COURIER_RUNTIME_BINARY ?? "";
const runtimeRoot = process.env.COURIER_RUNTIME_ROOT ?? "";
const portBase = Number(process.env.COURIER_RUNTIME_PORT_BASE ?? "26000");
let fixtureIndex = 0;

test.describe.configure({ mode: "serial" });
test.skip(!binary || !runtimeRoot, "compiled Courier runtime fixture is unavailable");
test.setTimeout(90_000);

type CommandResult = { code: number; stdout: string; stderr: string };

class CourierProcess {
  readonly child: ChildProcessWithoutNullStreams;
  stdout = "";
  stderr = "";
  private exit?: Promise<CommandResult>;

  constructor(readonly home: string, args: string[]) {
    this.child = spawn(binary, args, { env: courierEnvironment(home), stdio: ["pipe", "pipe", "pipe"] });
    this.child.stdout.setEncoding("utf8");
    this.child.stderr.setEncoding("utf8");
    this.child.stdout.on("data", (chunk: string) => { this.stdout += chunk; });
    this.child.stderr.on("data", (chunk: string) => { this.stderr += chunk; });
    this.exit = new Promise((resolve, reject) => {
      this.child.once("error", reject);
      this.child.once("exit", (code, signal) => resolve({ code: code ?? (signal ? 128 : 1), stdout: this.stdout, stderr: this.stderr }));
    });
  }

  async match(pattern: RegExp, timeout = 15_000): Promise<RegExpMatchArray> {
    const deadline = Date.now() + timeout;
    while (Date.now() < deadline) {
      const found = this.stdout.match(pattern);
      if (found) return found;
      if (this.child.exitCode !== null) throw new Error(`Courier exited before readiness (${this.child.exitCode})\nstdout: ${this.stdout}\nstderr: ${this.stderr}`);
      await new Promise((resolve) => setTimeout(resolve, 25));
    }
    throw new Error(`Timed out waiting for ${pattern}\nstdout: ${this.stdout}\nstderr: ${this.stderr}`);
  }

  async wait(): Promise<CommandResult> { return this.exit!; }

  async dispose(): Promise<void> {
    if (this.child.exitCode !== null) return;
    this.child.kill("SIGINT");
    const completed = await Promise.race([this.wait().then(() => true), new Promise<boolean>((resolve) => setTimeout(() => resolve(false), 3_000))]);
    if (!completed && this.child.exitCode === null) this.child.kill("SIGKILL");
    await this.wait();
  }
}

function courierEnvironment(home: string): NodeJS.ProcessEnv {
  return { ...process.env, HOME: home, XDG_CONFIG_HOME: join(home, ".config"), NO_COLOR: "1", LANG: "en_US.UTF-8" };
}

async function fixture(name: string): Promise<{ root: string; home: string }> {
  const root = await mkdtemp(join(runtimeRoot, `${name}-`));
  const home = join(runtimeRoot, `h${fixtureIndex++}`);
  await mkdir(home, { recursive: true, mode: 0o700 });
  return { root, home };
}

async function courier(home: string, args: string[]): Promise<CommandResult> {
  try {
    const result = await executeFile(binary, args, { env: courierEnvironment(home), maxBuffer: 8 << 20 });
    return { code: 0, stdout: result.stdout, stderr: result.stderr };
  } catch (error) {
    const failure = error as { code?: number; stdout?: string; stderr?: string };
    return { code: typeof failure.code === "number" ? failure.code : 1, stdout: failure.stdout ?? "", stderr: failure.stderr ?? String(error) };
  }
}

async function hosted(home: string, args: string[]): Promise<{ process: CourierProcess; url: string; id: string }> {
  const process = new CourierProcess(home, args);
  const ready = await process.match(/delivery: (http:\/\/[^\n]+)\nid: ([0-9a-f-]{36})\n/);
  return { process, url: ready[1], id: ready[2] };
}

async function stopHosted(home: string, process: CourierProcess, id: string): Promise<void> {
  const stopped = await courier(home, ["servers", "stop", id]);
  expect(stopped.code, stopped.stderr).toBe(0);
  const result = await process.wait();
  expect(result.code, result.stderr).toBe(0);
  expect(result.stdout).toContain(`stopped: ${id}`);
  const listed = await courier(home, ["servers"]);
  expect(listed.code, listed.stderr).toBe(0);
  expect(listed.stdout).toBe("No Courier data servers found.\n");
  await expectNoControlSockets(home);
}

async function stopHostedByServer(home: string, process: CourierProcess, deliveryID: string): Promise<void> {
  const listed = await courier(home, ["servers"]);
  const server = listed.stdout.match(/server ([0-9a-f-]{36})/);
  expect(server, listed.stdout).not.toBeNull();
  const stopped = await courier(home, ["servers", "stop", server![1]]);
  expect(stopped.code, stopped.stderr).toBe(0);
  const result = await process.wait();
  expect(result.code, result.stderr).toBe(0);
  expect(result.stdout).toContain(`stopped: ${deliveryID}`);
  await expectNoControlSockets(home);
}

async function expectNoControlSockets(home: string): Promise<void> {
  for (const directory of [join(home, ".config", "courier", "state"), join(home, "Library", "Application Support", "courier", "state")]) {
    let names: string[];
    try { names = await readdir(directory); } catch { continue; }
    for (const name of names) {
      const info = await lstat(join(directory, name));
      expect(info.isSocket(), `leftover control socket ${join(directory, name)}`).toBe(false);
    }
  }
}

async function upload(page: Page, name: string, data: Buffer, mimeType = "application/octet-stream"): Promise<void> {
  const response = page.waitForResponse((candidate) => candidate.url().includes("/api/v1/upload") && candidate.request().method() === "POST");
  await page.locator('input[type="file"]').setInputFiles({ name, mimeType, buffer: data });
  expect((await response).status()).toBe(201);
}

async function createArchive(root: string, name: string, files: Record<string, Buffer | string>): Promise<string> {
  const source = join(root, name);
  for (const [relative, contents] of Object.entries(files)) {
    const target = join(source, relative);
    await mkdir(join(target, ".."), { recursive: true });
    await writeFile(target, contents);
  }
  const archive = join(root, `${name}.tar.gz`);
  execFileSync("tar", ["-czf", archive, "-C", root, name]);
  return archive;
}

async function extractArchive(archive: string, root: string): Promise<void> {
  await mkdir(root, { recursive: true });
  execFileSync("tar", ["-xzf", archive, "-C", root]);
}

test("compiled browser upload and webhook routes commit exact files and stop cleanly", async ({ page }) => {
  const direct = await fixture("incoming");
  const destination = join(direct.root, "browser-destination");
  await mkdir(destination, { recursive: true });
  const browser = await hosted(direct.home, ["from", "web://", "to", destination, "--listen", `127.0.0.1:${portBase}`]);
  try {
    await page.goto(browser.url);
    await expect(page.locator("[data-courier-manifest]")).toBeVisible();
    const unicode = Buffer.from("Привет, Courier 🌍\n", "utf8");
    const binaryData = Buffer.from([0, 1, 2, 3, 127, 128, 254, 255]);
    await upload(page, "unicode-текст.txt", unicode, "text/plain");
    await upload(page, "binary.bin", binaryData);
    expect(await readFile(join(destination, "unicode-текст.txt"))).toEqual(unicode);
    expect(await readFile(join(destination, "binary.bin"))).toEqual(binaryData);
    await stopHosted(direct.home, browser.process, browser.id);
  } finally {
    await browser.process.dispose();
  }

  const extractionDestination = join(direct.root, "browser-extracted");
  await mkdir(extractionDestination, { recursive: true });
  const archive = await createArchive(direct.root, "bundle", { "nested/hello.txt": "hello", "bytes.bin": Buffer.from([9, 8, 7, 0]) });
  const extracting = await hosted(direct.home, ["from", "web://", "to", extractionDestination, "--extract", "--listen", `127.0.0.1:${portBase + 1}`]);
  try {
    await page.goto(extracting.url);
    const response = page.waitForResponse((candidate) => candidate.url().includes("/api/v1/upload") && candidate.request().method() === "POST");
    await page.locator('input[type="file"]').setInputFiles(archive);
    expect((await response).status()).toBe(201);
    expect((await readFile(join(extractionDestination, "bundle", "nested", "hello.txt"))).toString()).toBe("hello");
    expect(await readFile(join(extractionDestination, "bundle", "bytes.bin"))).toEqual(Buffer.from([9, 8, 7, 0]));
    await stopHostedByServer(direct.home, extracting.process, extracting.id);
  } finally {
    await extracting.process.dispose();
  }

  const webhookDestination = join(direct.root, "webhook-destination");
  await mkdir(webhookDestination, { recursive: true });
  const webhook = await hosted(direct.home, ["from", "webhook://", "to", webhookDestination, "--listen", `127.0.0.1:${portBase + 2}`]);
  try {
    const textResponse = await fetch(webhook.url, { method: "POST", body: oneFileForm("note.txt", Buffer.from("webhook text"), "text/plain") });
    expect(textResponse.status).toBe(201);
    const binaryResponse = await fetch(webhook.url, { method: "POST", body: oneFileForm("payload.bin", Buffer.from([0, 255, 4, 5])) });
    expect(binaryResponse.status).toBe(201);
    expect(await readFile(join(webhookDestination, "payload.bin"))).toEqual(Buffer.from([0, 255, 4, 5]));
    const collision = await fetch(webhook.url, { method: "POST", body: oneFileForm("note.txt", Buffer.from("replacement"), "text/plain") });
    expect(collision.status).toBe(409);
    expect((await readFile(join(webhookDestination, "note.txt"))).toString()).toBe("webhook text");
    const multiple = new FormData();
    multiple.append("file", new Blob(["first"]), "first.txt");
    multiple.append("file", new Blob(["second"]), "second.txt");
    const rejected = await fetch(webhook.url, { method: "POST", body: multiple });
    expect(rejected.status).toBe(400);
    await stopHosted(direct.home, webhook.process, webhook.id);
  } finally {
    await webhook.process.dispose();
  }

  const webhookExtracted = join(direct.root, "webhook-extracted");
  await mkdir(webhookExtracted, { recursive: true });
  const webhookExtract = await hosted(direct.home, ["from", "webhook://", "to", webhookExtracted, "--extract", "--listen", `127.0.0.1:${portBase + 3}`]);
  try {
    const extracted = await fetch(webhookExtract.url, { method: "POST", body: oneFileForm(basename(archive), await readFile(archive), "application/gzip") });
    expect(extracted.status).toBe(201);
    expect((await readFile(join(webhookExtracted, "bundle", "nested", "hello.txt"))).toString()).toBe("hello");
    await stopHosted(direct.home, webhookExtract.process, webhookExtract.id);
  } finally {
    await webhookExtract.process.dispose();
  }
});

function oneFileForm(name: string, contents: Buffer, type = "application/octet-stream"): FormData {
  const form = new FormData();
  form.append("file", new Blob([contents], { type }), name);
  return form;
}

test("compiled browser downloads preserve direct files and nested directory archives", async ({ page }) => {
  const current = await fixture("downloads");
  const directFile = join(current.root, "direct.bin");
  const directBytes = Buffer.from([0, 42, 128, 255, 10]);
  await writeFile(directFile, directBytes);
  const direct = await hosted(current.home, ["from", directFile, "to", "web://", "--listen", `127.0.0.1:${portBase + 4}`]);
  try {
    await page.goto(direct.url);
    const download = page.waitForEvent("download");
    await page.getByRole("link", { name: "Download file" }).click();
    const saved = join(current.root, "downloaded-direct.bin");
    await (await download).saveAs(saved);
    expect(await readFile(saved)).toEqual(directBytes);
    await stopHosted(current.home, direct.process, direct.id);
  } finally {
    await direct.process.dispose();
  }

  const directory = join(current.root, "shared-tree");
  await mkdir(join(directory, "nested"), { recursive: true });
  const rootBytes = Buffer.from([1, 2, 3, 200, 0]);
  await writeFile(join(directory, "unicode.txt"), "данные\n");
  await writeFile(join(directory, "binary.bin"), rootBytes);
  await writeFile(join(directory, "nested", "inside.txt"), "inside");
  const shared = await hosted(current.home, ["from", directory, "to", "web://", "--listen", `127.0.0.1:${portBase + 5}`]);
  try {
    await page.goto(shared.url);
    await expect(page.getByText("binary.bin", { exact: true })).toBeVisible();
    let download = page.waitForEvent("download");
    await page.locator("li").filter({ hasText: "binary.bin" }).getByRole("link", { name: "Download file" }).click();
    const binaryDownload = join(current.root, "binary-download.bin");
    await (await download).saveAs(binaryDownload);
    expect(await readFile(binaryDownload)).toEqual(rootBytes);

    await page.getByRole("button", { name: "nested", exact: true }).click();
    await expect(page.getByText("inside.txt", { exact: true })).toBeVisible();
    download = page.waitForEvent("download");
    await page.locator("li").filter({ hasText: "inside.txt" }).getByRole("link", { name: "Download file" }).click();
    const insideDownload = join(current.root, "inside-download.txt");
    await (await download).saveAs(insideDownload);
    expect((await readFile(insideDownload)).toString()).toBe("inside");
    await page.getByRole("button", { name: "Parent directory" }).click();

    download = page.waitForEvent("download");
    await page.getByRole("link", { name: "Download directory" }).click();
    const archiveDownload = join(current.root, "shared-tree.tar.gz");
    await (await download).saveAs(archiveDownload);
    const extracted = join(current.root, "downloaded-tree");
    await extractArchive(archiveDownload, extracted);
    expect(await readFile(join(extracted, "shared-tree", "binary.bin"))).toEqual(rootBytes);
    expect((await readFile(join(extracted, "shared-tree", "nested", "inside.txt"))).toString()).toBe("inside");
    await stopHosted(current.home, shared.process, shared.id);
  } finally {
    await shared.process.dispose();
  }
});

type MultipartUpload = { name: string; filename: string; data: Buffer };

function parseMultipart(contentType: string, body: Buffer): MultipartUpload[] {
  const boundary = contentType.match(/boundary=(?:"([^"]+)"|([^;]+))/)?.slice(1).find(Boolean);
  if (!boundary) throw new Error(`missing multipart boundary: ${contentType}`);
  return body.toString("latin1").split(`--${boundary}`).flatMap((part) => {
    const trimmed = part.replace(/^\r\n/, "").replace(/\r\n$/, "");
    if (!trimmed || trimmed === "--") return [];
    const separator = trimmed.indexOf("\r\n\r\n");
    if (separator < 0) return [];
    const headers = trimmed.slice(0, separator);
    const disposition = headers.match(/content-disposition:[^\r\n]+/i)?.[0] ?? "";
    const name = disposition.match(/name="([^"]+)"/)?.[1] ?? "";
    const filename = disposition.match(/filename="([^"]*)"/)?.[1] ?? "";
    let data = trimmed.slice(separator + 4);
    if (data.endsWith("\r\n--")) data = data.slice(0, -4);
    if (data.endsWith("\r\n")) data = data.slice(0, -2);
    return [{ name, filename, data: Buffer.from(data, "latin1") }];
  });
}

async function multipartReceiver(port: number): Promise<{ server: Server; uploads: MultipartUpload[][] }> {
  const uploads: MultipartUpload[][] = [];
  const server = createServer((request, response) => {
    const chunks: Buffer[] = [];
    request.on("data", (chunk: Buffer) => chunks.push(chunk));
    request.on("end", () => {
      try {
        uploads.push(parseMultipart(request.headers["content-type"] ?? "", Buffer.concat(chunks)));
        response.writeHead(204).end();
      } catch (error) {
        response.writeHead(400).end(String(error));
      }
    });
  });
  await new Promise<void>((resolve, reject) => {
    server.once("error", reject);
    server.listen(port, "127.0.0.1", () => resolve());
  });
  return { server, uploads };
}

test("compiled outgoing HTTP route sends one exact multipart file directly or as an archive", async () => {
  const current = await fixture("outgoing-webhook");
  const receiver = await multipartReceiver(portBase + 6);
  try {
    const direct = join(current.root, "payload.bin");
    const bytes = Buffer.from([0, 10, 13, 128, 255]);
    await writeFile(direct, bytes);
    let result = await courier(current.home, ["from", direct, "to", `http://127.0.0.1:${portBase + 6}/receive`]);
    expect(result.code, result.stderr).toBe(0);
    expect(receiver.uploads).toHaveLength(1);
    expect(receiver.uploads[0]).toHaveLength(1);
    expect(receiver.uploads[0][0].name).toBe("file");
    expect(receiver.uploads[0][0].data).toEqual(bytes);

    const directory = join(current.root, "outgoing-tree");
    await mkdir(join(directory, "nested"), { recursive: true });
    await writeFile(join(directory, "nested", "note.txt"), "archive payload");
    result = await courier(current.home, ["from", directory, "to", `http://127.0.0.1:${portBase + 6}/receive`, "--archive"]);
    expect(result.code, result.stderr).toBe(0);
    expect(receiver.uploads).toHaveLength(2);
    expect(receiver.uploads[1]).toHaveLength(1);
    expect(receiver.uploads[1][0].name).toBe("file");
    const archive = join(current.root, "outgoing.tar.gz");
    await writeFile(archive, receiver.uploads[1][0].data);
    const extracted = join(current.root, "outgoing-extracted");
    await extractArchive(archive, extracted);
    expect((await readFile(join(extracted, "outgoing-tree", "nested", "note.txt"))).toString()).toBe("archive payload");
  } finally {
    await new Promise<void>((resolve) => receiver.server.close(() => resolve()));
  }
});

test("compiled foreground administration UI exits cleanly after ui stop", async ({ page }) => {
  const current = await fixture("admin");
  const process = new CourierProcess(current.home, ["ui", "start", "--listen", `127.0.0.1:${portBase + 7}`]);
  try {
    const ready = await process.match(/Courier administration UI: (http:\/\/[^\n]+)\n/);
    await page.goto(ready[1]);
    await expect(page.locator("body")).toContainText("Courier");
    const stopped = await courier(current.home, ["ui", "stop"]);
    expect(stopped.code, stopped.stderr).toBe(0);
    const result = await process.wait();
    expect(result.code, result.stderr).toBe(0);
    expect(result.stdout).toContain("Courier administration UI stopped.");
    const repeated = await courier(current.home, ["ui", "stop"]);
    expect(repeated.code, repeated.stderr).toBe(0);
    expect(repeated.stdout).toContain("already stopped");
  } finally {
    await process.dispose();
  }
});
