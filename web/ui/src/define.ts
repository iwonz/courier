import { CourierButton } from "./components/button";
import { CourierBrand, CourierMascot, CourierRoute, CourierStatus } from "./components/brand";
import { CourierLocaleSelector } from "./components/locale-selector";
import { CourierPanel } from "./components/panel";
import { CourierProgress } from "./components/progress";
import { CourierSegmentedControl } from "./components/segmented-control";
import { CourierThemeSelector } from "./components/theme-selector";
import { CourierIcon } from "./icons";

export interface ElementRegistry {
  define(name: string, constructor: CustomElementConstructor): void;
  get(name: string): CustomElementConstructor | undefined;
}

const elements = {
  "courier-brand": CourierBrand,
  "courier-button": CourierButton,
  "courier-icon": CourierIcon,
  "courier-locale-selector": CourierLocaleSelector,
  "courier-mascot": CourierMascot,
  "courier-panel": CourierPanel,
  "courier-progress": CourierProgress,
  "courier-route": CourierRoute,
  "courier-segmented-control": CourierSegmentedControl,
  "courier-status": CourierStatus,
  "courier-theme-selector": CourierThemeSelector,
} as const;

export function defineCourierElements(registry: ElementRegistry = customElements): void {
  for (const [name, constructor] of Object.entries(elements)) {
    if (!registry.get(name)) {
      registry.define(name, constructor);
    }
  }
}
