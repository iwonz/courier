import * as React from "react";
import { PixelIcon } from "./icons-react";
import { browserPreferenceController, preferenceLocaleOrder, preferenceThemeOrder, type BrowserPreferenceController } from "./preferences";
import { translate, type Locale } from "./i18n";
import { type ThemePreference } from "./theme";
import { Button } from "./components/ui/button";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./components/ui/tooltip";

interface PreferenceValue {
  readonly locale: Locale;
  readonly theme: ThemePreference;
  readonly controller: BrowserPreferenceController;
}

const PreferenceContext = React.createContext<PreferenceValue | undefined>(undefined);

export function PreferenceProvider({ children, controller = browserPreferenceController() }: { readonly children: React.ReactNode; readonly controller?: BrowserPreferenceController }): React.JSX.Element {
  const [value, setValue] = React.useState(() => ({ locale: controller.locale, theme: controller.theme }));
  React.useEffect(() => {
    const update = (): void => setValue({ locale: controller.locale, theme: controller.theme });
    controller.addEventListener("change", update);
    return () => controller.removeEventListener("change", update);
  }, [controller]);
  return <PreferenceContext.Provider value={{ ...value, controller }}>{children}</PreferenceContext.Provider>;
}

export function usePreferences(): PreferenceValue {
  const value = React.useContext(PreferenceContext);
  if (!value) throw new Error("PreferenceProvider is required");
  return value;
}

export function ThemeSelector(): React.JSX.Element {
  const { locale, theme, controller } = usePreferences();
  const next = preferenceThemeOrder[(preferenceThemeOrder.indexOf(theme) + 1) % preferenceThemeOrder.length]!;
  const currentLabel = translate(locale, `theme.${theme}`);
  const nextLabel = translate(locale, `theme.${next}`);
  const label = `${translate(locale, "theme.label")}: ${currentLabel}. ${translate(locale, "preference.next")}: ${nextLabel}`;
  const icon = theme === "light" ? "sun" : theme === "dark" ? "moon" : "system";
  return <TooltipProvider><Tooltip><TooltipTrigger asChild><Button type="button" size="icon" variant="ghost" aria-label={label} onClick={() => controller.cycleTheme()}><PixelIcon name={icon} /></Button></TooltipTrigger><TooltipContent>{currentLabel}</TooltipContent></Tooltip></TooltipProvider>;
}

export function LocaleSelector(): React.JSX.Element {
  const { locale, controller } = usePreferences();
  const next = preferenceLocaleOrder[(preferenceLocaleOrder.indexOf(locale) + 1) % preferenceLocaleOrder.length]!;
  const currentLabel = translate(locale, `locale.${locale}`);
  const nextLabel = translate(locale, `locale.${next}`);
  const label = `${translate(locale, "locale.label")}: ${currentLabel}. ${translate(locale, "preference.next")}: ${nextLabel}`;
  return <TooltipProvider><Tooltip><TooltipTrigger asChild><Button type="button" size="icon" variant="ghost" aria-label={label} onClick={() => controller.cycleLocale()}><span aria-hidden="true" data-locale-icon={locale} className="grid h-3 w-4 overflow-hidden border border-foreground/35">{locale === "ru" ? <><i className="block h-1 bg-white" /><i className="block h-1 bg-[#0039A6]" /><i className="block h-1 bg-[#D52B1E]" /></> : <svg viewBox="0 0 16 12" className="size-full" shapeRendering="crispEdges"><path fill="#012169" d="M0 0h16v12H0z"/><path fill="#FFFFFF" d="M0 0h3l13 9v3h-3L0 3zm16 0h-3L0 9v3h3l13-9zM6 0h4v12H6zM0 4h16v4H0z"/><path fill="#C8102E" d="M7 0h2v12H7zM0 5h16v2H0zM0 0h1l6 4v1H6L0 1zm16 0h-1L9 4v1h1l6-4zM0 12h1l6-4V7H6l-6 4zm16 0h-1L9 8V7h1l6 4z"/></svg>}</span><span className="sr-only">{currentLabel}</span></Button></TooltipTrigger><TooltipContent>{currentLabel}</TooltipContent></Tooltip></TooltipProvider>;
}
