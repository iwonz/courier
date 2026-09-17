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
  RelaySprite,
  RouteDisplay,
  ThemeSelector,
  buttonVariants,
  cn,
  usePreferences,
} from "@courier/ui";
import { childPath, downloadURL, loadMetadata, login, parentPath, upload, type Entry, type Metadata } from "./api";
import { dataText } from "./catalog";

function DataContent(): React.JSX.Element {
  const { locale } = usePreferences();
  const [metadata, setMetadata] = React.useState<Metadata>();
  const [failed, setFailed] = React.useState(false);
  const [loading, setLoading] = React.useState(true);
  const [csrf, setCsrf] = React.useState("");
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

  const sendFile = async (event: React.ChangeEvent<HTMLInputElement>): Promise<void> => {
    const input = event.currentTarget;
    const file = input.files?.item(0);
    if (!file) return;
    try {
      await upload(file, csrf);
      input.value = "";
      await refresh();
    } catch {
      setFailed(true);
    }
  };

  const entry = (item: Entry): React.JSX.Element => {
    const path = childPath(metadata!.path, item.name);
    return <li key={`${item.type}:${item.name}`} className="grid min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-3 px-4 py-3">
      <span className="grid size-9 place-items-center text-primary"><PixelIcon name={item.type === "directory" ? "folder" : "file"} className="size-4" /></span>
      <div className="min-w-0"><button type="button" className="max-w-full truncate text-left text-sm font-semibold hover:text-primary disabled:pointer-events-none" disabled={item.type !== "directory"} onClick={() => void refresh(path)}>{item.name}</button><p className="text-xs text-muted-foreground">{item.size} {t("itemSize")}</p></div>
      <Button asChild variant="ghost" size="sm"><a href={downloadURL(path, item.type === "directory")}><PixelIcon name="download" />{item.type === "directory" ? t("downloadArchive") : t("download")}</a></Button>
    </li>;
  };

  return <div className="min-h-screen px-4 pb-16 sm:px-6 lg:px-8">
    <header className="mx-auto flex min-h-20 w-full max-w-6xl items-center justify-between gap-4"><Brand /><nav className="flex items-center gap-2"><ThemeSelector /><LocaleSelector /></nav></header>
    <main className="mx-auto grid w-full max-w-6xl gap-6 pt-8 sm:pt-14">
      <div className="grid items-end gap-8 lg:grid-cols-[minmax(0,1fr)_18rem]">
        <div className="grid gap-4"><Badge variant="secondary" className="w-fit font-mono">{t("privateRoute")}</Badge><h1 className="text-5xl font-black tracking-[-.06em] sm:text-7xl">{t("title")}</h1></div>
        <RelaySprite role="delivery" className="hidden w-full max-w-72 justify-self-end lg:block" />
      </div>

      {failed && !metadata ? <Card data-courier-auth-region className="mx-auto w-full max-w-xl">
        <CardHeader><div className="mb-2 text-primary"><PixelIcon name="shield" className="size-8" /></div><CardTitle>{t("accessTitle")}</CardTitle><CardDescription>{t("accessHelp")}</CardDescription></CardHeader>
        <CardContent className="grid gap-4"><Alert variant="destructive"><AlertTitle>{t("privateRoute")}</AlertTitle><AlertDescription>{t("failed")}</AlertDescription></Alert><form className="grid gap-3" onSubmit={signIn}><label className="grid gap-2 text-sm font-semibold" htmlFor="delivery-password">{t("password")}</label><Input id="delivery-password" name="password" type="password" autoComplete="current-password" placeholder={t("password")} /><div className="flex flex-wrap gap-2"><Button type="submit"><PixelIcon name="shield" />{t("signIn")}</Button><Button type="button" variant="outline" onClick={() => void refresh()}><PixelIcon name="retry" />{t("retry")}</Button></div></form></CardContent>
      </Card> : null}

      {metadata ? <>
        <RouteDisplay source="sender" destination={metadata.name} />
        <Card data-courier-manifest><CardHeader className="flex-row items-center justify-between gap-4"><div><CardTitle>{metadata.name}</CardTitle><CardDescription>{t("manifest")}</CardDescription></div><Badge variant="success"><PixelIcon name="check" className="size-3.5" />{t("ready")}</Badge></CardHeader>
          <CardContent className="p-0">
            {metadata.type === "upload" ? <div className="grid justify-items-start gap-3 p-5 sm:p-6"><h2 className="text-lg font-semibold">{t("uploadTitle")}</h2><p className="max-w-2xl text-sm text-muted-foreground">{t("uploadHelp")}</p><label className={cn(buttonVariants({ variant: "default" }), "cursor-pointer")}><PixelIcon name="upload" className="size-4" />{t("upload")}<input type="file" className="sr-only" onChange={(event) => void sendFile(event)} /></label></div> : metadata.type === "file" ? <div className="p-5 sm:p-6"><Button asChild><a href={downloadURL(metadata.path)}><PixelIcon name="download" />{t("download")}</a></Button></div> : <>
              <div className="flex flex-wrap gap-2 p-4"><Button asChild><a href={downloadURL(metadata.path, true)}><PixelIcon name="download" />{t("downloadAll")}</a></Button>{metadata.path ? <Button type="button" variant="ghost" onClick={() => void refresh(parentPath(metadata.path))}><PixelIcon name="arrow-up" />{t("up")}</Button> : null}</div>
              {(metadata.entries ?? []).length ? <ul className="divide-y divide-border/70">{metadata.entries!.map(entry)}</ul> : <p className="p-5 text-sm text-muted-foreground">{t("empty")}</p>}
            </>}
          </CardContent>
        </Card>
      </> : !failed && loading ? <Card><CardContent className="flex min-h-40 items-center justify-center text-sm text-muted-foreground">{t("loading")}</CardContent></Card> : null}
    </main>
  </div>;
}

export function DataApp(): React.JSX.Element {
  return <PreferenceProvider><DataContent /></PreferenceProvider>;
}
