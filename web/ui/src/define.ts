import { CourierButton } from "./components/button";
import { CourierBrandIcon } from "./components/brand-icon";
import { CourierBrand, CourierMascot, CourierRoute, CourierStatus } from "./components/brand";
import { CourierCheckbox } from "./components/checkbox";
import { CourierIconLink } from "./components/icon-link";
import { CourierLocaleSelector } from "./components/locale-selector";
import { CourierPanel } from "./components/panel";
import { CourierProgress } from "./components/progress";
import { CourierSegmentedControl } from "./components/segmented-control";
import { CourierThemeSelector } from "./components/theme-selector";
import { CourierCommandReadout, CourierTerminal } from "./components/terminal";
import { CourierScene } from "./components/scene";
import { CourierIcon } from "./icons";

export interface ElementRegistry {
  define(name: string, constructor: CustomElementConstructor): void;
  get(name: string): CustomElementConstructor | undefined;
}

const elements = {
  "courier-brand": CourierBrand,
  "courier-brand-icon": CourierBrandIcon,
  "courier-button": CourierButton,
  "courier-checkbox": CourierCheckbox,
  "courier-icon": CourierIcon,
  "courier-icon-link": CourierIconLink,
  "courier-locale-selector": CourierLocaleSelector,
  "courier-mascot": CourierMascot,
  "courier-panel": CourierPanel,
  "courier-progress": CourierProgress,
  "courier-route": CourierRoute,
  "courier-scene": CourierScene,
  "courier-segmented-control": CourierSegmentedControl,
  "courier-status": CourierStatus,
  "courier-theme-selector": CourierThemeSelector,
  "courier-terminal": CourierTerminal,
  "courier-command-readout": CourierCommandReadout,
} as const;

export function defineCourierElements(registry: ElementRegistry = customElements): void {
  for (const [name, constructor] of Object.entries(elements)) {
    if (!registry.get(name)) {
      registry.define(name, constructor);
    }
  }
}
