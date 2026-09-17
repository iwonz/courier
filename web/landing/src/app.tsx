import * as React from "react";
import {
  Badge,
  Brand,
  BrandIcon,
  Button,
  Checkbox,
  CommandReadout,
  Github,
  Icon,
  LocaleSelector,
  PixelIcon,
  PreferenceProvider,
  RelaySprite,
  ScrollArea,
  Separator,
  ThemeSelector,
  pixelBezierPath,
  type BrandIconName,
  type IconName,
  usePreferences,
} from "@courier/ui";
import { landingText, type LandingMessage } from "./catalog";
import { contractData, type LandingCommand, type LandingFlag, type LandingRoute } from "./contract";

interface InstallChannel { readonly name: string; readonly command: string; readonly icon: BrandIconName; }
export interface RoutePair { readonly source: string; readonly destination: string; readonly routeName: string; readonly allowedFlags: readonly string[]; }
export interface EndpointPresentation { readonly label: LandingMessage; readonly description: LandingMessage; readonly icon: IconName; }

export const installs: readonly InstallChannel[] = [
  { name: "curl", command: "curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh", icon: "curl" },
  { name: "wget", command: "wget -qO- https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh", icon: "wget" },
  { name: "PowerShell", command: "irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex", icon: "powershell" },
  { name: "npm", command: "npm install --global @iwonz/courier", icon: "npm" },
  { name: "npx", command: "npx @iwonz/courier --help", icon: "npx" },
  { name: "Yarn", command: "yarn dlx @iwonz/courier --help", icon: "yarn" },
  { name: "pnpm", command: "pnpm dlx @iwonz/courier --help", icon: "pnpm" },
  { name: "Homebrew", command: "brew tap iwonz/courier https://github.com/iwonz/courier && brew install --cask iwonz/courier/courier", icon: "homebrew" },
  { name: "Scoop", command: "scoop bucket add courier https://github.com/iwonz/courier && scoop install courier/courier", icon: "scoop" },
];

const endpointSamples: Readonly<Record<string, Readonly<Record<"source" | "destination", string>>>> = {
  local: { source: "./project", destination: "./backup/" },
  ssh: { source: "courier@host:/srv/source", destination: "courier@host:/srv/destination/" },
  web: { source: "web://", destination: "web://" },
  webhook: { source: "webhook://", destination: "webhook://" },
  http: { source: "https://api.example.test/source", destination: "https://api.example.test/upload" },
};

const endpointPresentations: Readonly<Record<"source" | "destination", Readonly<Record<string, EndpointPresentation>>>> = {
  source: {
    local: { label: "endpointLocal", description: "sourceLocalDescription", icon: "folder-out" },
    ssh: { label: "endpointRemote", description: "sourceRemoteDescription", icon: "server-out" },
    web: { label: "endpointWeb", description: "sourceWebDescription", icon: "browser-upload" },
    webhook: { label: "endpointWebHook", description: "sourceWebHookDescription", icon: "webhook-in" },
  },
  destination: {
    local: { label: "endpointLocal", description: "destinationLocalDescription", icon: "folder-in" },
    ssh: { label: "endpointRemote", description: "destinationRemoteDescription", icon: "server-in" },
    web: { label: "endpointWeb", description: "destinationWebDescription", icon: "browser-share" },
    http: { label: "endpointWebHook", description: "destinationWebHookDescription", icon: "webhook-out" },
  },
};

export function expandRoutePairs(routes: readonly LandingRoute[]): RoutePair[] {
  return routes.flatMap((route) => route.source.flatMap((source) => route.destination.map((destination) => ({ source, destination, routeName: route.name, allowedFlags: route.allowedFlags }))));
}
export function endpointExample(name: string, side: "source" | "destination"): string { return endpointSamples[name]?.[side] ?? name; }
export function endpointPresentation(name: string, side: "source" | "destination"): EndpointPresentation | undefined { return endpointPresentations[side][name]; }
export function endpointMessage(name: string, side: "source" | "destination"): LandingMessage | undefined { return endpointPresentation(name, side)?.label; }
export function endpointIcon(name: string, side: "source" | "destination"): IconName { return endpointPresentation(name, side)?.icon ?? "route"; }
export function localizedEndpointLabel(name: string, side: "source" | "destination", translate: (message: LandingMessage) => string): string {
  const message = endpointMessage(name, side);
  return message ? translate(message) : name;
}
export function localizedEndpointDescription(name: string, side: "source" | "destination", translate: (message: LandingMessage) => string): string {
  const message = endpointPresentation(name, side)?.description;
  return message ? translate(message) : name;
}
export function applicableFlags(pair: RoutePair, flags: readonly LandingFlag[]): LandingFlag[] {
  const exactRoute = `${pair.source}-to-${pair.destination}`;
  return flags.filter((flag) => pair.allowedFlags.includes(flag.name) && (flag.appliesTo.includes(pair.routeName) || flag.appliesTo.includes(exactRoute)));
}
export function commandFlags(commandName: string, compatibleOnly: boolean, commands: readonly LandingCommand[], flags: readonly LandingFlag[]): LandingFlag[] {
  if (!commandName || !compatibleOnly) return [...flags];
  const names = commands.find((command) => command.name === commandName)?.flags ?? [];
  return flags.filter((flag) => names.includes(flag.name));
}

const routePairs = expandRoutePairs(contractData.routes);
const initialRoute = routePairs.find((pair) => pair.source === "local" && pair.destination === "ssh")!;
const sourceEndpoints = [...new Set(routePairs.map((pair) => pair.source))];
const destinationEndpoints = [...new Set(routePairs.map((pair) => pair.destination))];

function useRouteConnector(source: string, destination: string): {
  readonly containerRef: React.RefObject<HTMLDivElement | null>;
  readonly sourceRefs: React.MutableRefObject<Map<string, HTMLButtonElement>>;
  readonly destinationRefs: React.MutableRefObject<Map<string, HTMLButtonElement>>;
  readonly viewBox: string;
  readonly path: string;
} {
  const containerRef = React.useRef<HTMLDivElement>(null);
  const sourceRefs = React.useRef(new Map<string, HTMLButtonElement>());
  const destinationRefs = React.useRef(new Map<string, HTMLButtonElement>());
  const [geometry, setGeometry] = React.useState({ viewBox: "0 0 1 1", path: "" });
  React.useLayoutEffect(() => {
    const container = containerRef.current!;
    const sourceNode = sourceRefs.current.get(source)!;
    const destinationNode = destinationRefs.current.get(destination)!;
    const measure = (): void => {
      const bounds = container.getBoundingClientRect();
      const sourceBounds = sourceNode.getBoundingClientRect();
      const destinationBounds = destinationNode.getBoundingClientRect();
      if (bounds.width <= 0 || bounds.height <= 0) return;
      setGeometry({
        viewBox: `0 0 ${bounds.width} ${bounds.height}`,
        path: pixelBezierPath(
          { x: sourceBounds.right - bounds.left, y: sourceBounds.top + sourceBounds.height / 2 - bounds.top },
          { x: destinationBounds.left - bounds.left, y: destinationBounds.top + destinationBounds.height / 2 - bounds.top },
        ),
      });
    };
    measure();
    const observer = globalThis.ResizeObserver ? new ResizeObserver(measure) : undefined;
    observer?.observe(container);
    globalThis.addEventListener("resize", measure);
    return () => { observer?.disconnect(); globalThis.removeEventListener("resize", measure); };
  }, [source, destination]);
  return { containerRef, sourceRefs, destinationRefs, ...geometry };
}

function LandingContent(): React.JSX.Element {
  const { locale } = usePreferences();
  const [selectedSource, setSelectedSource] = React.useState(initialRoute.source);
  const [selectedDestination, setSelectedDestination] = React.useState(initialRoute.destination);
  const [activeInstall, setActiveInstall] = React.useState(installs[0]!.name);
  const [selectedCommand, setSelectedCommand] = React.useState("");
  const [compatibleOnly, setCompatibleOnly] = React.useState(true);
  const t = React.useCallback((message: LandingMessage) => landingText(locale, message), [locale]);
  const selected = routePairs.find((pair) => pair.source === selectedSource && pair.destination === selectedDestination)!;
  const routeFlags = applicableFlags(selected, contractData.flags);
  const install = installs.find((channel) => channel.name === activeInstall)!;
  const visibleFlags = commandFlags(selectedCommand, compatibleOnly, contractData.commands, contractData.flags);
  const selectedCommandData = contractData.commands.find((command) => command.name === selectedCommand);
  const routeCommand = `courier ${t("from")} ${endpointExample(selected.source, "source")} ${t("to")} ${endpointExample(selected.destination, "destination")}`;
  const connector = useRouteConnector(selectedSource, selectedDestination);

  const chooseSource = (source: string): void => {
    const next = routePairs.find((pair) => pair.source === source && pair.destination === selectedDestination) ?? routePairs.find((pair) => pair.source === source)!;
    setSelectedSource(next.source);
    setSelectedDestination(next.destination);
  };
  const chooseCommand = (command: string): void => {
    if (selectedCommand === command) {
      setSelectedCommand("");
      setCompatibleOnly(true);
      return;
    }
    setSelectedCommand(command);
  };
  const endpointButton = (name: string, side: "source" | "destination"): React.JSX.Element => {
    const valid = side === "source" || routePairs.some((pair) => pair.source === selectedSource && pair.destination === name);
    const active = side === "source" ? name === selectedSource : name === selectedDestination;
    const references = side === "source" ? connector.sourceRefs : connector.destinationRefs;
    return <Button
      key={name}
      ref={(node) => { if (node) references.current.set(name, node); else references.current.delete(name); }}
      type="button"
      variant={active ? "default" : "ghost"}
      size="sm"
      aria-pressed={active}
      aria-label={localizedEndpointLabel(name, side, t)}
      aria-description={localizedEndpointDescription(name, side, t)}
      disabled={!valid}
      onClick={() => side === "source" ? chooseSource(name) : setSelectedDestination(name)}
      className="relative z-10 w-full justify-start sm:w-auto sm:min-w-28"
    ><Icon name={endpointIcon(name, side)} /><span>{localizedEndpointLabel(name, side, t)}</span></Button>;
  };

  return <div className="min-h-screen overflow-x-clip">
    <header className="fixed inset-x-0 top-0 z-50 bg-transparent">
      <div className="mx-auto flex min-h-16 w-full max-w-7xl items-center gap-2 px-4 sm:px-6 lg:px-8">
        <a href="#route" className="courier-pixel-focus shrink-0 outline-none"><Brand /></a>
        <nav aria-label="Courier" className="hidden items-center gap-1 sm:flex">
          <Button asChild variant="ghost" size="sm"><a href="#install">{t("installShort")}</a></Button>
          <Button asChild variant="ghost" size="sm"><a href="#cli">{t("cliShort")}</a></Button>
        </nav>
        <div className="ml-auto flex items-center gap-2">
          <Button asChild variant="ghost" size="icon"><a href="https://github.com/iwonz/courier" target="_blank" rel="noopener noreferrer" aria-label={t("githubLabel")}><Github /></a></Button>
          <ThemeSelector />
          <LocaleSelector />
        </div>
      </div>
    </header>

    <main className="pt-20">
      <section id="route" className="scroll-mt-20 px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
        <div data-courier-route-composition className="relative isolate mx-auto grid w-full max-w-7xl gap-12">
          <div className="grid min-h-[27rem] lg:grid-cols-[minmax(0,.9fr)_minmax(24rem,1.1fr)] lg:items-center">
            <div className="relative z-10 max-w-3xl self-start lg:self-center"><h1 className="text-balance text-[clamp(3.4rem,8vw,7.8rem)] font-bold leading-[.84] tracking-[-.04em] text-foreground">{t("title")}</h1></div>
            <div className="relative grid min-h-72 place-items-center lg:min-h-[25rem]">
              <RelaySprite role="route" className="relative z-10 w-[min(32rem,92%)]" />
            </div>
          </div>
          <div className="grid gap-5">
              <div ref={connector.containerRef} className="relative grid grid-cols-2 gap-5 sm:gap-16">
                {connector.path ? <svg className="pointer-events-none absolute inset-0 z-0 size-full overflow-visible" viewBox={connector.viewBox} preserveAspectRatio="none" aria-hidden="true"><path d={connector.path} fill="none" stroke="currentColor" strokeWidth="2" className="text-primary/60" vectorEffect="non-scaling-stroke" shapeRendering="crispEdges" /><rect width="8" height="8" fill="currentColor" className="text-warning" style={{ offsetPath: `path('${connector.path}')`, animation: "courier-pixel-route 2.8s steps(16,end) infinite" }} /></svg> : null}
                <div className="grid content-start gap-2"><span className="px-1 text-xs font-semibold text-muted-foreground">{t("sourceLabel")}</span>{sourceEndpoints.map((name) => endpointButton(name, "source"))}</div>
                <div className="grid content-start gap-2"><span className="px-1 text-xs font-semibold text-muted-foreground">{t("destinationLabel")}</span>{destinationEndpoints.map((name) => endpointButton(name, "destination"))}</div>
              </div>
              <CommandReadout
                heading={t("commandLabel")}
                command={routeCommand}
                sessionKey={`${locale}:${selected.source}:${selected.destination}`}
                copyLabel={t("copyCommand")}
                copiedLabel={t("copiedCommand")}
                copyFailedLabel={t("copyFailed")}
                details={<div className="flex flex-wrap gap-2">{routeFlags.map((flag) => <Badge variant="secondary" key={flag.name} className="font-mono font-medium">{flag.syntax}</Badge>)}</div>}
              />
          </div>
        </div>
      </section>

      <section id="install" className="scroll-mt-20 px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
        <div className="mx-auto grid w-full max-w-7xl gap-6">
          <h2 className="text-4xl font-black tracking-[-.055em] sm:text-6xl">{t("install")}</h2>
          <div className="flex flex-wrap gap-2">{installs.map((channel) => <Button key={channel.name} type="button" variant={channel.name === install.name ? "default" : "ghost"} className={channel.name === install.name ? undefined : "bg-muted/38"} size="sm" aria-pressed={channel.name === install.name} onClick={() => setActiveInstall(channel.name)}><BrandIcon name={channel.icon} />{channel.name}</Button>)}</div>
          <CommandReadout
            heading={t("commandLabel")}
            command={install.command}
            sessionKey={`${locale}:${install.name}`}
            copyLabel={t("copyCommand")}
            copiedLabel={t("copiedCommand")}
            copyFailedLabel={t("copyFailed")}
            details={<div className="flex items-center gap-2 text-sm font-semibold text-primary"><BrandIcon name={install.icon} />{install.name}</div>}
            footerActions={<>
              <Button asChild variant="ghost" size="sm"><a href="https://github.com/iwonz/courier/releases/latest" target="_blank" rel="noopener noreferrer"><PixelIcon name="package" />{t("packages")}</a></Button>
              <Button asChild variant="ghost" size="sm"><a href="https://github.com/iwonz/courier/releases/latest" target="_blank" rel="noopener noreferrer"><PixelIcon name="download" />{t("direct")}</a></Button>
            </>}
          />
        </div>
      </section>

      <section id="cli" className="scroll-mt-20 px-4 pb-14 pt-6 sm:px-6 sm:pb-16 sm:pt-8 lg:px-8">
        <div className="mx-auto grid w-full max-w-7xl gap-6">
          <h2 className="text-4xl font-black tracking-[-.055em] sm:text-6xl">{t("cliTitle")}</h2>
          <div data-courier-cli-registry>
            <div className="grid divide-y divide-border/35 lg:grid-cols-[minmax(18rem,.72fr)_minmax(0,1.28fr)] lg:divide-x lg:divide-y-0">
              <section><div className="px-5 pt-5 text-sm font-semibold">{t("commandLabel")}</div><ScrollArea className="h-[31rem]"><div className="grid gap-1 p-2">{contractData.commands.map((command) => <Button key={command.name} type="button" variant={command.name === selectedCommand ? "secondary" : "ghost"} aria-pressed={command.name === selectedCommand} onClick={() => chooseCommand(command.name)} className="h-auto justify-start whitespace-normal px-3 py-3 text-left"><code className="font-mono text-xs leading-relaxed">{command.usage}</code></Button>)}</div></ScrollArea></section>
              <section><div className="flex min-h-14 flex-wrap items-center justify-between gap-3 px-5 pt-3"><span className="text-sm font-semibold">{t("optionsLabel")}</span><label className="flex items-center gap-2 text-xs font-medium text-muted-foreground"><Checkbox checked={compatibleOnly} disabled={!selectedCommand} onCheckedChange={(checked) => setCompatibleOnly(checked === true)} />{t("compatibleOnly")}</label></div><ScrollArea className="h-[31rem]"><div className="grid gap-1 p-2">{visibleFlags.length ? visibleFlags.map((flag) => <article key={flag.name} className="grid gap-2 px-3 py-3 hover:bg-muted"><div className="flex flex-wrap items-baseline justify-between gap-2"><code className="font-mono text-sm font-bold text-primary">{flag.syntax}</code><span className="text-xs text-muted-foreground">{t("optionDefault")}: {flag.default}</span></div><p className="text-xs leading-relaxed text-muted-foreground">{t("repeatable")}: {flag.repeatable ? t("yes") : t("no")} · {t("applies")}: {flag.appliesTo.join(", ")}</p></article>) : <p className="p-5 text-sm text-muted-foreground">{t("noCompatibleOptions")}</p>}</div></ScrollArea></section>
            </div>
            <CommandReadout className="border-t border-border/35 pt-4" heading={t("commandLabel")} command={selectedCommandData?.usage ?? ""} description={selectedCommandData ? "" : t("selectCommand")} sessionKey={`${locale}:${selectedCommand}`} copyLabel={t("copyCommand")} copiedLabel={t("copiedCommand")} copyFailedLabel={t("copyFailed")} />
          </div>
        </div>
      </section>
    </main>
  </div>;
}

export function LandingApp(): React.JSX.Element {
  return <PreferenceProvider><LandingContent /></PreferenceProvider>;
}
