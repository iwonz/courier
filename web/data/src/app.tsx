import * as React from "react";
import {
  Alert,
  AlertDescription,
  AlertTitle,
  Badge,
  Brand,
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Input,
  LocaleSelector,
  PixelIcon,
  PreferenceProvider,
  Progress,
  RelaySprite,
  RouteDisplay,
  ThemeSelector,
  buttonVariants,
  cn,
  usePreferences,
} from "@courier/ui";
import { childPath, downloadURL, loadMetadata, login, parentPath, upload, type Entry, type Metadata } from "./api";
import { dataText } from "./catalog";

type UploadPhase = "queued" | "uploading" | "success" | "failure";

interface UploadItem {
  id: string;
  name: string;
  size: number;
  progress: number;
  phase: UploadPhase;
  reason?: string;
}

function DataContent(): React.JSX.Element {
  const { locale } = usePreferences();
  const [metadata, setMetadata] = React.useState<Metadata>();
  const [failed, setFailed] = React.useState(false);
  const [loading, setLoading] = React.useState(true);
  const [csrf, setCsrf] = React.useState("");
  const [uploadItems, setUploadItems] = React.useState<UploadItem[]>([]);
  const [uploading, setUploading] = React.useState(false);
  const uploadBatch = React.useRef(0);
  const t = React.useCallback((message: Parameters<typeof dataText>[1]) => dataText(locale, message), [locale]);

  const refresh = React.useCallback(async (path = ""): Promise<void> => {
    setFailed(false);
    setLoading(true);
    try {
      setMetadata(await loadMetadata(path));
    } catch {
      setFailed(true);
      setMetadata(undefined);
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => { void refresh(); }, [refresh]);

  const signIn = async (event: React.FormEvent<HTMLFormElement>): Promise<void> => {
    event.preventDefault();
    const form = event.currentTarget;
    const password = new FormData(form).get("password")?.toString() ?? "";
    try {
      setCsrf((await login(password)).csrf);
      form.reset();
      await refresh();
    } catch {
      setFailed(true);
    }
  };

  const updateUpload = React.useCallback((id: string, update: Partial<UploadItem>): void => {
    setUploadItems((items) => items.map((item) => item.id === id ? { ...item, ...update } : item));
  }, []);

  const sendFiles = async (event: React.ChangeEvent<HTMLInputElement>): Promise<void> => {
    const input = event.currentTarget;
    const files = Array.from(input.files ?? []);
    if (files.length === 0 || uploading) return;
    const batch = ++uploadBatch.current;
    const selected = files.map((file, index) => ({
      file,
      item: { id: `${batch}:${index}`, name: file.name, size: file.size, progress: 0, phase: "queued" as const },
    }));
    input.value = "";
    setUploadItems(selected.map(({ item }) => item));
    setUploading(true);
    try {
      for (const { file, item } of selected) {
        updateUpload(item.id, { phase: "uploading", progress: 0, reason: undefined });
        try {
          await upload(file, csrf, ({ percent }) => updateUpload(item.id, { progress: percent }));
          updateUpload(item.id, { phase: "success", progress: 100 });
        } catch (error) {
          const reason = error instanceof Error ? error.message : String(error);
          updateUpload(item.id, { phase: "failure", reason });
        }
      }
    } finally {
      setUploading(false);
    }
  };

  const uploadPresentation = (item: UploadItem): { label: string; variant: "secondary" | "default" | "success" | "destructive"; icon: string } => {
    switch (item.phase) {
      case "queued": return { label: t("uploadQueuedTitle"), variant: "secondary", icon: "package" };
      case "uploading": return { label: t("uploadingTitle"), variant: "default", icon: "activity" };
      case "success": return { label: t("uploadSuccessTitle"), variant: "success", icon: "check" };
      case "failure": return { label: t("uploadFailureTitle"), variant: "destructive", icon: "stop" };
    }
  };

  const completedUploads = uploadItems.filter((item) => item.phase === "success" || item.phase === "failure").length;

  const entry = (item: Entry): React.JSX.Element => {
    const path = childPath(metadata!.path, item.name);
    return <li key={`${item.type}:${item.name}`} className="grid min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-3 border-b border-border px-4 py-3 last:border-b-0">
      <span className="grid size-9 place-items-center text-primary"><PixelIcon name={item.type === "directory" ? "folder" : "file"} className="size-4" /></span>
      <div className="min-w-0"><button type="button" className="max-w-full truncate text-left text-sm font-semibold hover:text-primary disabled:pointer-events-none" disabled={item.type !== "directory"} onClick={() => void refresh(path)}>{item.name}</button><p className="text-xs text-muted-foreground">{item.size} {t("itemSize")}</p></div>
      <Button asChild variant="ghost" size="sm"><a href={downloadURL(path, item.type === "directory")}><PixelIcon name="download" />{item.type === "directory" ? t("downloadArchive") : t("download")}</a></Button>
    </li>;
  };

  return <div className="min-h-screen pb-16">
    <header className="border-b border-border bg-background"><div className="mx-auto flex h-16 w-full max-w-6xl items-center justify-between gap-4 px-4 sm:px-6 lg:px-8"><Brand /><div className="flex items-center gap-2"><ThemeSelector /><LocaleSelector /></div></div></header>
    <main className="mx-auto grid w-full max-w-6xl gap-0 px-4 pt-8 sm:px-6 sm:pt-12 lg:px-8">
      <div className="grid items-center border-y border-border bg-card sm:grid-cols-[minmax(0,1fr)_12rem] lg:grid-cols-[minmax(0,1fr)_18rem]">
        <div className="grid gap-4 p-5 sm:p-8"><Badge variant="secondary" className="w-fit font-mono uppercase tracking-[.08em]">{t("privateRoute")}</Badge><h1 className="text-5xl font-black tracking-[-.045em] sm:text-7xl">{t("title")}</h1></div>
        <RelaySprite role="delivery" className="hidden w-full max-w-64 justify-self-end p-3 sm:block" />
      </div>

      {failed && !metadata ? <Card data-courier-auth-region className="mt-8 w-full border-y border-border bg-card">
        <CardHeader><div className="mb-2 text-primary"><PixelIcon name="shield" className="size-8" /></div><CardTitle>{t("accessTitle")}</CardTitle><CardDescription>{t("accessHelp")}</CardDescription></CardHeader>
        <CardContent className="grid gap-4"><Alert variant="destructive"><AlertTitle>{t("privateRoute")}</AlertTitle><AlertDescription>{t("failed")}</AlertDescription></Alert><form className="grid gap-3" onSubmit={signIn}><label className="grid gap-2 text-sm font-semibold" htmlFor="delivery-password">{t("password")}</label><Input id="delivery-password" name="password" type="password" autoComplete="current-password" placeholder={t("password")} /><div className="flex flex-wrap gap-2"><Button type="submit"><PixelIcon name="shield" />{t("signIn")}</Button><Button type="button" variant="outline" onClick={() => void refresh()}><PixelIcon name="retry" />{t("retry")}</Button></div></form></CardContent>
      </Card> : null}

      {metadata ? <>
        <RouteDisplay className="mt-8 border-y border-border bg-card px-4" source="sender" destination={metadata.name} />
        <Card data-courier-manifest className="border-b border-border bg-card"><CardHeader className="flex-row items-center justify-between gap-4 border-b border-border"><div><CardTitle>{metadata.name}</CardTitle><CardDescription>{t("manifest")}</CardDescription></div><Badge variant="success"><PixelIcon name="check" className="size-3.5" />{t("ready")}</Badge></CardHeader>
          <CardContent className="p-0">
            {metadata.type === "upload" ? <div className="grid justify-items-start gap-3 p-5 sm:p-6"><h2 className="text-lg font-semibold">{t("uploadTitle")}</h2><p id="courier-upload-help" className="max-w-2xl text-sm text-muted-foreground">{t("uploadHelp")}</p><label aria-disabled={uploading} className={cn(buttonVariants({ variant: "default" }), uploading ? "cursor-wait opacity-70" : "cursor-pointer")}><PixelIcon name="upload" className="size-4" />{uploading ? t("uploadingAction") : t("upload")}<input type="file" multiple aria-describedby="courier-upload-help" className="sr-only" disabled={uploading} onChange={(event) => void sendFiles(event)} /></label>
              {uploadItems.length ? <section data-courier-upload-status className="mt-2 w-full border-y border-border bg-background" aria-labelledby="courier-upload-batch-title">
                <div className="flex flex-wrap items-baseline justify-between gap-2 border-b border-border px-4 py-3"><h3 id="courier-upload-batch-title" className="font-semibold">{t("uploadBatchTitle")}</h3><p data-courier-upload-summary role="status" aria-live="polite" className="text-xs font-semibold text-muted-foreground">{t("uploadBatchProgress")}: {completedUploads}/{uploadItems.length}</p></div>
                <ul data-courier-upload-list>{uploadItems.map((item) => {
                  const presentation = uploadPresentation(item);
                  return <li key={item.id} data-courier-upload-item data-state={item.phase} className="grid min-w-0 grid-cols-[auto_minmax(0,1fr)] gap-x-3 gap-y-2 border-b border-border px-4 py-3 last:border-b-0 sm:grid-cols-[auto_minmax(0,1fr)_auto]">
                    <span className="row-span-2 grid size-9 place-items-center bg-secondary text-primary"><PixelIcon name="file" className="size-4" /></span>
                    <div className="min-w-0"><p className="truncate text-sm font-semibold">{item.name}</p><p className="text-xs text-muted-foreground">{item.size} {t("itemSize")}</p>{item.reason ? <p data-courier-upload-error className="mt-1 break-words text-xs text-destructive">{item.reason}<br />{t("uploadFailureHelp")}</p> : null}</div>
                    <Badge variant={presentation.variant} className="col-start-2 sm:col-start-3"><PixelIcon name={presentation.icon} className="size-3.5" />{presentation.label}</Badge>
                    <div className="col-span-2 grid min-w-0 grid-cols-[minmax(0,1fr)_3rem] items-center gap-3 sm:col-span-2 sm:col-start-2"><Progress value={item.progress} aria-label={`${item.name}: ${presentation.label}`} /><span className="text-right font-mono text-xs font-semibold tabular-nums">{item.progress}%</span></div>
                  </li>;
                })}</ul>
              </section> : null}
            </div> : metadata.type === "file" ? <div className="p-5 sm:p-6"><Button asChild><a href={downloadURL(metadata.path)}><PixelIcon name="download" />{t("download")}</a></Button></div> : <>
              <div className="flex flex-wrap gap-2 p-4"><Button asChild><a href={downloadURL(metadata.path, true)}><PixelIcon name="download" />{t("downloadAll")}</a></Button>{metadata.path ? <Button type="button" variant="ghost" onClick={() => void refresh(parentPath(metadata.path))}><PixelIcon name="arrow-up" />{t("up")}</Button> : null}</div>
              {(metadata.entries ?? []).length ? <ul>{metadata.entries!.map(entry)}</ul> : <p className="p-5 text-sm text-muted-foreground">{t("empty")}</p>}
            </>}
          </CardContent>
        </Card>
      </> : !failed && loading ? <Card className="mt-8 border-y border-border bg-card"><CardContent className="flex min-h-40 items-center justify-center text-sm text-muted-foreground">{t("loading")}</CardContent></Card> : null}
    </main>
  </div>;
}

export function DataApp(): React.JSX.Element {
  return <PreferenceProvider><DataContent /></PreferenceProvider>;
}
