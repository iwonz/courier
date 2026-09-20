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
  Input,
  LocaleSelector,
  PixelIcon,
  PreferenceProvider,
  RelaySprite,
  ScrollArea,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  ThemeSelector,
  pixelBezierPath,
  type BrandIconName,
  type IconName,
  usePreferences,
} from "@courier/ui";
import { landingText, type LandingMessage } from "./catalog";
import { contractData, type LandingCommand, type LandingFlag, type LandingRoute } from "./contract";
import { argumentValid, buildCommand, flagEnabled, updateFlagValues, type FlagValues, type ShellMode } from "./command-builder";

interface InstallChannel { readonly name: string; readonly command: string; readonly icon?: BrandIconName; }
export interface RoutePair { readonly source: string; readonly destination: string; readonly routeName: string; readonly allowedFlags: readonly string[]; }
export interface EndpointPresentation { readonly label: LandingMessage; readonly description: LandingMessage; readonly icon: IconName; }

export const installs: readonly InstallChannel[] = [
  { name: "curl", command: "curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh", icon: "curl" },
  { name: "wget", command: "wget -qO- https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh" },
  { name: "PowerShell", command: "irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex", icon: "powershell" },
  { name: "npm", command: "npm install --global @iwonz/courier", icon: "npm" },
  { name: "npx", command: "npx @iwonz/courier --help" },
  { name: "Yarn", command: "yarn dlx @iwonz/courier --help", icon: "yarn" },
  { name: "pnpm", command: "pnpm dlx @iwonz/courier --help", icon: "pnpm" },
  { name: "Homebrew", command: "brew tap iwonz/courier https://github.com/iwonz/courier && brew install iwonz/courier/courier", icon: "homebrew" },
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
const endpointCount = new Set([...sourceEndpoints, ...destinationEndpoints]).size;

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
  const [argumentValues, setArgumentValues] = React.useState<Record<string, string>>({});
  const [flagValues, setFlagValues] = React.useState<FlagValues>({});
  const [shellMode, setShellMode] = React.useState<ShellMode>("posix");
  const t = React.useCallback((message: LandingMessage) => landingText(locale, message), [locale]);
  const selected = routePairs.find((pair) => pair.source === selectedSource && pair.destination === selectedDestination)!;
  const routeFlags = applicableFlags(selected, contractData.flags);
  const install = installs.find((channel) => channel.name === activeInstall)!;
  const visibleFlags = commandFlags(selectedCommand, compatibleOnly, contractData.commands, contractData.flags);
  const selectedCommandData = contractData.commands.find((command) => command.name === selectedCommand);
  const builtCommand = buildCommand(selectedCommandData, argumentValues, flagValues, contractData.flags, shellMode);
  const routeCommand = `courier ${t("from")} ${endpointExample(selected.source, "source")} ${t("to")} ${endpointExample(selected.destination, "destination")}`;
  const connector = useRouteConnector(selectedSource, selectedDestination);

  const chooseSource = (source: string): void => {
    const next = routePairs.find((pair) => pair.source === source && pair.destination === selectedDestination) ?? routePairs.find((pair) => pair.source === source)!;
    setSelectedSource(next.source);
    setSelectedDestination(next.destination);
  };
  const chooseCommand = (command: string): void => {
    setArgumentValues({});
    setFlagValues({});
    setCompatibleOnly(true);
    if (selectedCommand === command) {
      setSelectedCommand("");
      return;
    }
    setSelectedCommand(command);
  };
  const setArgumentValue = (name: string, value: string, omitWhenFlag?: string): void => {
    setArgumentValues((current) => ({ ...current, [name]: value }));
    if (omitWhenFlag && value.trim()) setFlagValues((current) => ({ ...current, [omitWhenFlag]: [] }));
  };
  const setFlagValue = (flag: LandingFlag, values: readonly string[]): void => {
    setFlagValues((current) => updateFlagValues(current, flag, values, contractData.flags));
    if (flag.name === "all" && values.some((value) => value === "true")) setArgumentValues((current) => ({ ...current, uuid: "" }));
  };
  const flagControl = (flag: LandingFlag): React.JSX.Element => {
    const compatible = Boolean(selectedCommandData?.flags?.includes(flag.name));
    const dependenciesMet = (flag.requires ?? []).every((name) => flagEnabled(flagValues, name));
    const disabled = !compatible || !dependenciesMet;
    const values = flagValues[flag.name] ?? [];
    if (flag.valueKind === "boolean") return <label className="flex min-h-10 items-center justify-end gap-2 text-xs font-medium text-muted-foreground">
      <Checkbox aria-label={flag.syntax} checked={flagEnabled(flagValues, flag.name)} disabled={disabled} onCheckedChange={(checked) => setFlagValue(flag, checked === true ? ["true"] : [])} />
      {t("enabled")}
    </label>;
    if (flag.valueKind === "enum") return <Select value={values[0] ?? ""} disabled={disabled} onValueChange={(value) => setFlagValue(flag, [value])}>
      <SelectTrigger aria-label={flag.syntax}><SelectValue placeholder={`${flag.placeholder} · ${t("optionDefault")}: ${flag.default}`} /></SelectTrigger>
      <SelectContent>{(flag.choices as string[]).map((choice) => <SelectItem key={choice} value={choice}>{choice}</SelectItem>)}</SelectContent>
    </Select>;
    if (flag.repeatable) {
      const rows = values.length ? values : [""];
      return <div className="grid gap-2">{rows.map((value, index) => <div key={`${flag.name}:${index}`} className="flex gap-2">
        <Input aria-label={`${flag.syntax} ${index + 1}`} value={value} disabled={disabled} placeholder={`${flag.placeholder} · ${t("optionDefault")}: ${flag.default}`} onChange={(event) => {
          const next = [...rows];
          next[index] = event.currentTarget.value;
          setFlagValue(flag, next);
        }} />
        <Button type="button" variant="ghost" size="icon" disabled={disabled} aria-label={index === rows.length - 1 ? t("addValue") : t("removeValue")} onClick={() => index === rows.length - 1 ? setFlagValue(flag, [...rows, ""]) : setFlagValue(flag, rows.filter((_, row) => row !== index))}>
          <span aria-hidden="true" className="text-base font-bold">{index === rows.length - 1 ? "+" : "−"}</span>
        </Button>
      </div>)}</div>;
    }
    return <Input aria-label={flag.syntax} value={values[0] ?? ""} disabled={disabled} inputMode={flag.valueKind === "unsigned" ? "numeric" : undefined} placeholder={`${flag.placeholder} · ${t("optionDefault")}: ${flag.default}`} onChange={(event) => setFlagValue(flag, [event.currentTarget.value])} />;
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
    <header className="border-b border-border bg-background">
      <div className="mx-auto flex h-16 w-full max-w-7xl items-center gap-2 px-4 sm:px-6 lg:px-8">
        <a href="#route" className="courier-pixel-focus shrink-0 outline-none"><Brand /></a>
        <div className="ml-auto flex items-center gap-2">
          <Button asChild variant="ghost" size="icon"><a href="https://github.com/iwonz/courier" target="_blank" rel="noopener noreferrer" aria-label={t("githubLabel")}><Github /></a></Button>
          <ThemeSelector />
          <LocaleSelector />
        </div>
      </div>
    </header>

    <main>
      <section id="route" className="px-4 py-8 sm:px-6 sm:py-10 lg:px-8">
        <div data-courier-route-composition className="relative isolate mx-auto grid w-full max-w-7xl gap-0 border-y border-border bg-card">
          <div className="grid border-b border-border lg:grid-cols-[minmax(0,.92fr)_minmax(22rem,1.08fr)]">
            <div className="relative z-10 grid content-between gap-10 p-5 sm:p-8 lg:min-h-[28rem]">
              <div className="grid gap-4"><span className="font-mono text-xs font-semibold uppercase tracking-[.12em] text-primary">{t("routeConsole")}</span><h1 className="text-balance text-[clamp(3rem,7vw,6.8rem)] font-bold leading-[.88] tracking-[-.04em] text-foreground">{t("title")}</h1></div>
              <div data-courier-contract-metrics className="grid grid-cols-3 border-y border-border">
                {[[t("endpointsMetric"), endpointCount], [t("routesMetric"), contractData.routes.length], [t("commandsMetric"), contractData.commands.length]].map(([label, value]) => <div key={label} className="grid gap-1 py-3 pr-3 sm:py-4"><span className="text-[.65rem] uppercase tracking-[.08em] text-muted-foreground">{label}</span><strong className="font-mono text-xl text-warning sm:text-2xl">{value}</strong></div>)}
              </div>
            </div>
            <div className="relative grid min-h-72 place-items-center overflow-hidden bg-[var(--card)] p-4 lg:min-h-[28rem]">
              <RelaySprite role="route" className="relative z-10 w-[min(28rem,92%)]" />
            </div>
          </div>
          <div className="grid gap-5 p-5 sm:p-8">
              <div ref={connector.containerRef} className="relative grid grid-cols-2 gap-5 sm:gap-16">
                {connector.path ? <><svg className="pointer-events-none absolute inset-0 z-0 size-full overflow-visible" viewBox={connector.viewBox} preserveAspectRatio="none" aria-hidden="true"><path data-courier-route-path d={connector.path} fill="none" stroke="currentColor" strokeWidth="2" className="text-primary/60" vectorEffect="non-scaling-stroke" shapeRendering="crispEdges" /></svg><span data-courier-route-signal aria-hidden="true" className="courier-route-signal pointer-events-none absolute left-0 top-0 z-[1] size-1 bg-warning" style={{ offsetPath: `path('${connector.path}')`, offsetAnchor: "center", offsetRotate: "0deg" }} /></> : null}
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

      <section id="install" className="px-4 py-8 sm:px-6 sm:py-10 lg:px-8">
        <div className="mx-auto grid w-full max-w-7xl gap-6">
          <h2 className="break-words text-3xl font-black leading-none tracking-[-.04em] sm:text-6xl">{t("install")}</h2>
          <div className="flex overflow-x-auto border-y border-border" role="tablist">{installs.map((channel) => <Button key={channel.name} type="button" role="tab" variant={channel.name === install.name ? "secondary" : "ghost"} className="h-12 shrink-0 border-0 px-4" size="sm" aria-selected={channel.name === install.name} aria-pressed={channel.name === install.name} onClick={() => setActiveInstall(channel.name)}>{channel.icon ? <BrandIcon name={channel.icon} /> : null}{channel.name}</Button>)}</div>
          <CommandReadout className="courier-terminal-panel px-4 sm:px-6"
            heading={t("commandLabel")}
            command={install.command}
            sessionKey={`${locale}:${install.name}`}
            copyLabel={t("copyCommand")}
            copiedLabel={t("copiedCommand")}
            copyFailedLabel={t("copyFailed")}
            details={<div className="flex items-center gap-2 text-sm font-semibold text-primary">{install.icon ? <BrandIcon name={install.icon} /> : null}{install.name}</div>}
            footerActions={<>
              <Button asChild variant="ghost" size="sm"><a href="https://github.com/iwonz/courier/releases/latest" target="_blank" rel="noopener noreferrer"><PixelIcon name="package" />{t("packages")}</a></Button>
              <Button asChild variant="ghost" size="sm"><a href="https://github.com/iwonz/courier/releases/latest" target="_blank" rel="noopener noreferrer"><PixelIcon name="download" />{t("direct")}</a></Button>
            </>}
          />
        </div>
      </section>

      <section id="cli" className="px-4 pb-14 pt-8 sm:px-6 sm:pb-16 sm:pt-10 lg:px-8">
        <div className="mx-auto grid w-full max-w-7xl gap-6">
          <h2 className="break-words text-3xl font-black leading-none tracking-[-.04em] sm:text-6xl">{t("cliTitle")}</h2>
          <div data-courier-cli-registry>
            <div data-courier-cli-columns className="grid gap-8 lg:grid-cols-[minmax(18rem,.72fr)_minmax(0,1.28fr)]">
              <section className="min-w-0"><div className="border-y border-border px-3 py-4 text-sm font-semibold uppercase tracking-[.08em]">{t("commandLabel")}</div><ScrollArea className="h-[31rem]"><div>{contractData.commands.map((command) => <Button key={command.name} type="button" variant={command.name === selectedCommand ? "secondary" : "ghost"} aria-pressed={command.name === selectedCommand} onClick={() => chooseCommand(command.name)} className="h-auto min-w-0 w-full justify-start whitespace-normal border-b border-border px-3 py-3 text-left"><code className="min-w-0 break-words font-mono text-xs leading-relaxed">{command.usage}</code></Button>)}</div></ScrollArea></section>
              <section className="min-w-0">
                <div className="flex min-h-14 min-w-0 flex-wrap items-center justify-between gap-3 border-y border-border px-3 py-3"><span className="text-sm font-semibold uppercase tracking-[.08em]">{t("optionsLabel")}</span><label className="flex min-w-0 items-center gap-2 text-xs font-medium text-muted-foreground"><Checkbox checked={compatibleOnly} disabled={!selectedCommand} onCheckedChange={(checked) => setCompatibleOnly(checked === true)} /><span className="min-w-0 break-words">{t("compatibleOnly")}</span></label></div>
                <ScrollArea className="h-[31rem]"><div>
                  {selectedCommandData?.arguments.map((argument) => {
                    const omitted = Boolean(argument.omitWhenFlag && flagEnabled(flagValues, argument.omitWhenFlag));
                    const valid = argumentValid(argument, argumentValues[argument.name] ?? "", flagValues);
                    return <article key={argument.name} data-builder-argument={argument.name} className="grid gap-3 border-b border-border px-3 py-3 sm:grid-cols-[minmax(0,1fr)_minmax(12rem,.9fr)] sm:items-center">
                      <div className="grid gap-1"><code className="font-mono text-sm font-bold text-primary">{`<${argument.name}>`}</code><span className="text-xs text-muted-foreground">{argument.required ? t("requiredValue") : t("optionalValue")}{argument.prefix ? ` · ${argument.prefix}` : ""}</span></div>
                      {argument.kind === "command-path" ? <Select value={argumentValues[argument.name] ?? ""} onValueChange={(value) => setArgumentValue(argument.name, value, argument.omitWhenFlag)}><SelectTrigger aria-label={argument.name}><SelectValue placeholder={t("chooseCommandPath")} /></SelectTrigger><SelectContent>{contractData.commands.filter((candidate) => candidate.name !== selectedCommandData.name).map((candidate) => <SelectItem key={candidate.name} value={candidate.path}>{candidate.path}</SelectItem>)}</SelectContent></Select> : <Input aria-label={argument.name} value={argumentValues[argument.name] ?? ""} disabled={omitted} aria-invalid={!valid} placeholder={omitted ? `--${argument.omitWhenFlag}` : argument.name} onChange={(event) => setArgumentValue(argument.name, event.currentTarget.value, argument.omitWhenFlag)} />}
                    </article>;
                  })}
                  {visibleFlags.length ? visibleFlags.map((flag) => <article key={flag.name} data-builder-flag={flag.name} className="grid gap-3 border-b border-border px-3 py-3 hover:bg-muted sm:grid-cols-[minmax(0,1fr)_minmax(12rem,.9fr)] sm:items-center"><div className="grid gap-1"><div className="flex flex-wrap items-baseline justify-between gap-2"><code className="font-mono text-sm font-bold text-primary">{flag.syntax}</code><span className="text-xs text-muted-foreground">{t("optionDefault")}: {flag.default}</span></div><p className="text-xs leading-relaxed text-muted-foreground">{t("repeatable")}: {flag.repeatable ? t("yes") : t("no")} · {t("applies")}: {flag.appliesTo.join(", ")}{(flag.requires ?? []).length ? ` · ${t("requires")}: ${(flag.requires as string[]).map((name) => `--${name}`).join(", ")}` : ""}</p></div>{flagControl(flag)}</article>) : <p className="border-b border-border p-5 text-sm text-muted-foreground">{t("noCompatibleOptions")}</p>}
                </div></ScrollArea>
              </section>
            </div>
            <div className="flex justify-end gap-2 py-2" aria-label={t("shellLabel")}><Button type="button" size="sm" variant={shellMode === "posix" ? "default" : "ghost"} aria-pressed={shellMode === "posix"} onClick={() => setShellMode("posix")}>POSIX</Button><Button type="button" size="sm" variant={shellMode === "powershell" ? "default" : "ghost"} aria-pressed={shellMode === "powershell"} onClick={() => setShellMode("powershell")}>PowerShell</Button></div>
            <CommandReadout className="mt-3 border-y border-border px-3" data-courier-cli-readout heading={t("commandLabel")} command={builtCommand.command} copyDisabled={!builtCommand.valid} description={!selectedCommandData ? t("selectCommand") : !builtCommand.valid ? t("completeCommand") : ""} sessionKey={`${locale}:${selectedCommand}:${shellMode}:${builtCommand.command}`} copyLabel={t("copyCommand")} copiedLabel={t("copiedCommand")} copyFailedLabel={t("copyFailed")} />
          </div>
        </div>
      </section>
    </main>
  </div>;
}

export function LandingApp(): React.JSX.Element {
  return <PreferenceProvider><LandingContent /></PreferenceProvider>;
}
