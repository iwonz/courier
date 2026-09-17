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
  Checkbox,
  Input,
  LocaleSelector,
  PixelIcon,
  PreferenceProvider,
  RelaySprite,
  RouteDisplay,
  ScrollArea,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  ThemeSelector,
  usePreferences,
} from "@courier/ui";
import { loadServers, savePolicy, stopTarget, subscribeSnapshots, type Delivery, type Policy, type Server, type Snapshot } from "./api";
import { adminText } from "./catalog";

export function reconcileDeliverySelection(snapshot: Snapshot, selected: string | undefined): string | undefined {
  const deliveries = snapshot.servers.flatMap((server) => server.deliveries);
  return deliveries.some((delivery) => delivery.id === selected) ? selected : deliveries[0]?.id;
}

function AdminContent(): React.JSX.Element {
  const { locale } = usePreferences();
  const [snapshot, setSnapshot] = React.useState<Snapshot>();
  const [failed, setFailed] = React.useState(false);
  const [conflict, setConflict] = React.useState(false);
  const [selectedDeliveryId, setSelectedDeliveryId] = React.useState<string>();
  const t = React.useCallback((message: Parameters<typeof adminText>[1]) => adminText(locale, message), [locale]);

  const applySnapshot = React.useCallback((next: Snapshot): void => {
    setSnapshot(next);
    setSelectedDeliveryId((current) => reconcileDeliverySelection(next, current));
    setFailed(false);
  }, []);

  const refresh = React.useCallback(async (): Promise<void> => {
    setFailed(false);
    setConflict(false);
    try {
      applySnapshot(await loadServers());
    } catch {
      setFailed(true);
    }
  }, [applySnapshot]);

  React.useEffect(() => {
    const unsubscribe = subscribeSnapshots(applySnapshot);
    void refresh();
    return unsubscribe;
  }, [applySnapshot, refresh]);

  const stop = async (kind: "servers" | "deliveries", id: string): Promise<void> => {
    try {
      await stopTarget(kind, id);
      await refresh();
    } catch {
      setFailed(true);
    }
  };

  const save = async (event: React.FormEvent<HTMLFormElement>, delivery: Delivery): Promise<void> => {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    const policy: Policy = {
      ...delivery.policy,
      version: delivery.policy.version + 1,
      auth: String(data.get("auth")) as Policy["auth"],
      authAttempts: Number(data.get("attempts")),
      authFailAction: String(data.get("failAction")) as Policy["authFailAction"],
      noUi: data.get("noUi") === "on",
    };
    setFailed(false);
    setConflict(false);
    try {
      await savePolicy(delivery, policy);
      await refresh();
    } catch (error) {
      setConflict(error instanceof Error && error.name === "ConflictError");
      setFailed(!(error instanceof Error && error.name === "ConflictError"));
    }
  };

  const servers = snapshot?.servers ?? [];
  const deliveries = servers.flatMap((server) => server.deliveries);
  const confirmed = deliveries.reduce((total, delivery) => total + delivery.counters.confirmed, 0);
  const selected = deliveries.find((delivery) => delivery.id === selectedDeliveryId);

  const serverNavigator = (server: Server): React.JSX.Element => <section key={server.id} className="grid gap-2 border-b border-border/35 py-4 last:border-b-0">
    <div className="grid gap-2"><div className="flex items-start justify-between gap-2"><div className="min-w-0"><p className="truncate font-mono text-xs font-semibold">{server.id}</p><p className="truncate font-mono text-xs text-muted-foreground">{server.bind}</p></div><Badge variant={server.status === "live" ? "success" : "destructive"}>{t(server.status)}</Badge></div><Button type="button" size="sm" variant="outline" onClick={() => void stop("servers", server.id)}><PixelIcon name="stop" />{t("stopServer")}</Button></div>
    <div className="grid gap-1">{server.deliveries.map((delivery) => <Button key={delivery.id} type="button" variant={delivery.id === selectedDeliveryId ? "secondary" : "ghost"} aria-pressed={delivery.id === selectedDeliveryId} onClick={() => setSelectedDeliveryId(delivery.id)} className="h-auto justify-start whitespace-normal px-3 py-3 text-left"><span className="grid min-w-0 gap-1"><strong className="truncate text-sm">{delivery.route}</strong><span className="truncate font-mono text-[.68rem] text-muted-foreground">{delivery.id}</span></span></Button>)}</div>
  </section>;

  return <div className="min-h-screen px-4 pb-16 sm:px-6 lg:px-8">
    <header className="mx-auto flex min-h-20 w-full max-w-7xl items-center justify-between gap-4"><Brand /><nav className="flex items-center gap-2"><Button type="button" variant="ghost" size="icon" aria-label={t("refresh")} onClick={() => void refresh()}><PixelIcon name="refresh" /></Button><ThemeSelector /><LocaleSelector /></nav></header>
    <main className="mx-auto grid w-full max-w-7xl gap-6 pt-8 sm:pt-12">
      <div className="grid items-end gap-8 lg:grid-cols-[minmax(0,1fr)_18rem]"><div className="grid gap-4"><Badge variant="secondary" className="w-fit font-mono">{t("eyebrow")}</Badge><h1 className="text-5xl font-bold tracking-[-.035em] sm:text-7xl">{t("title")}</h1><p className="max-w-2xl text-base leading-relaxed text-muted-foreground">{t("intro")}</p></div><RelaySprite role="admin" className="hidden w-full max-w-72 justify-self-end lg:block" /></div>

      <div data-courier-metrics className="grid divide-y divide-border/35 sm:grid-cols-3 sm:divide-x sm:divide-y-0">
        <div className="flex items-center gap-4 py-5 sm:px-5 sm:first:pl-0"><span className="text-primary"><PixelIcon name="server" /></span><div><strong className="font-mono text-2xl">{servers.length}</strong><p className="text-xs text-muted-foreground">{t("serversMetric")}</p></div></div>
        <div className="flex items-center gap-4 py-5 sm:px-5"><span className="text-primary"><PixelIcon name="workflow" /></span><div><strong className="font-mono text-2xl">{deliveries.length}</strong><p className="text-xs text-muted-foreground">{t("deliveriesMetric")}</p></div></div>
        <div className="flex items-center gap-4 py-5 sm:px-5 sm:last:pr-0"><span className="text-primary"><PixelIcon name="database" /></span><div><strong className="font-mono text-2xl">{confirmed}</strong><p className="text-xs text-muted-foreground">{t("confirmedMetric")}</p></div></div>
      </div>

      {failed ? <Alert variant="destructive"><AlertTitle>{t("unreachable")}</AlertTitle><AlertDescription className="flex flex-wrap items-center justify-between gap-3"><span>{t("failed")}</span><Button type="button" variant="outline" size="sm" onClick={() => void refresh()}>{t("retry")}</Button></AlertDescription></Alert> : null}
      {conflict ? <Alert variant="warning"><AlertTitle>{t("policy")}</AlertTitle><AlertDescription className="flex flex-wrap items-center justify-between gap-3"><span>{t("conflict")}</span><Button type="button" variant="outline" size="sm" onClick={() => void refresh()}>{t("refresh")}</Button></AlertDescription></Alert> : null}

      {snapshot ? servers.length ? <Card data-courier-admin-workspace><div className="grid min-h-[34rem] lg:grid-cols-[20rem_minmax(0,1fr)]"><ScrollArea className="h-[34rem] border-b border-border/35 pr-4 lg:border-b-0 lg:border-r">{servers.map(serverNavigator)}</ScrollArea><div className="min-w-0">{selected ? <article className="grid gap-5 py-5 lg:pl-7"><div className="flex flex-wrap items-start justify-between gap-3"><div><p className="text-xs font-semibold text-muted-foreground">{selected.route}</p><p className="mt-1 break-all font-mono text-sm font-bold">{selected.id}</p></div><Button type="button" variant="destructive" size="sm" onClick={() => void stop("deliveries", selected.id)}><PixelIcon name="stop" />{t("stopDelivery")}</Button></div><RouteDisplay source={selected.source || t("unavailable")} destination={selected.destination || t("unavailable")} /><div className="grid gap-4 py-2 sm:grid-cols-3"><div><span className="text-xs text-muted-foreground">{t("source")}</span><p className="break-all font-mono text-xs font-semibold">{selected.source || t("unavailable")}</p></div><div><span className="text-xs text-muted-foreground">{t("destination")}</span><p className="break-all font-mono text-xs font-semibold">{selected.destination || t("unavailable")}</p></div><div><span className="text-xs text-muted-foreground">{t("transferred")}</span><p className="font-mono text-xs font-semibold">{selected.counters.confirmed}</p></div></div><form className="grid gap-4 sm:grid-cols-2 xl:grid-cols-[1fr_1fr_1fr_auto_auto] xl:items-end" onSubmit={(event) => void save(event, selected)}>
        <label className="grid gap-2 text-xs font-semibold text-muted-foreground">{t("authentication")}<Select name="auth" defaultValue={selected.policy.auth}><SelectTrigger><SelectValue /></SelectTrigger><SelectContent><SelectItem value="none">none</SelectItem><SelectItem value="basic">basic</SelectItem><SelectItem value="password">password</SelectItem></SelectContent></Select></label>
        <label className="grid gap-2 text-xs font-semibold text-muted-foreground">{t("attempts")}<Input name="attempts" type="number" min="1" defaultValue={selected.policy.authAttempts} /></label>
        <label className="grid gap-2 text-xs font-semibold text-muted-foreground">{t("failAction")}<Select name="failAction" defaultValue={selected.policy.authFailAction}><SelectTrigger><SelectValue /></SelectTrigger><SelectContent><SelectItem value="ban">ban</SelectItem><SelectItem value="stop">stop</SelectItem></SelectContent></Select></label>
        <label className="flex h-10 items-center gap-2 text-xs font-semibold text-muted-foreground"><Checkbox name="noUi" defaultChecked={selected.policy.noUi} />{t("noUi")}</label>
        <Button type="submit"><PixelIcon name="save" />{t("save")}</Button>
      </form></article> : <div className="grid min-h-[34rem] place-items-center p-6 text-sm text-muted-foreground">{t("empty")}</div>}</div></div></Card> : <Card><CardHeader><CardTitle>{t("live")}</CardTitle><CardDescription>{t("empty")}</CardDescription></CardHeader></Card> : !failed ? <Card><CardContent className="flex min-h-40 items-center justify-center gap-2 text-sm text-muted-foreground"><PixelIcon name="activity" className="size-4" />{t("loading")}</CardContent></Card> : null}
    </main>
  </div>;
}

export function AdminApp(): React.JSX.Element {
  return <PreferenceProvider><AdminContent /></PreferenceProvider>;
}
