import * as React from "react";
import { Monitor, Moon, Sun } from "lucide-react";
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
  const ThemeIcon = theme === "light" ? Sun : theme === "dark" ? Moon : Monitor;
  return <TooltipProvider><Tooltip><TooltipTrigger asChild><Button type="button" size="icon" variant="ghost" className="bg-muted/45" aria-label={label} onClick={() => controller.cycleTheme()}><ThemeIcon /></Button></TooltipTrigger><TooltipContent>{currentLabel}</TooltipContent></Tooltip></TooltipProvider>;
}

export function LocaleSelector(): React.JSX.Element {
  const { locale, controller } = usePreferences();
  const next = preferenceLocaleOrder[(preferenceLocaleOrder.indexOf(locale) + 1) % preferenceLocaleOrder.length]!;
  const currentLabel = translate(locale, `locale.${locale}`);
  const nextLabel = translate(locale, `locale.${next}`);
  const label = `${translate(locale, "locale.label")}: ${currentLabel}. ${translate(locale, "preference.next")}: ${nextLabel}`;
  return <TooltipProvider><Tooltip><TooltipTrigger asChild><Button type="button" size="icon" variant="ghost" className="bg-muted/45" aria-label={label} onClick={() => controller.cycleLocale()}><span aria-hidden="true" data-locale-icon={locale} className="text-base leading-none">{locale === "ru" ? "🇷🇺" : "🇬🇧"}</span><span className="sr-only">{currentLabel}</span></Button></TooltipTrigger><TooltipContent>{currentLabel}</TooltipContent></Tooltip></TooltipProvider>;
}
