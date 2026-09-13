import { CourierButton } from "./components/button";
import { CourierLocaleSelector } from "./components/locale-selector";
import { CourierPanel } from "./components/panel";
import { CourierProgress } from "./components/progress";
import { CourierThemeSelector } from "./components/theme-selector";
import { CourierIcon } from "./icons";

export interface ElementRegistry {
  define(name: string, constructor: CustomElementConstructor): void;
  get(name: string): CustomElementConstructor | undefined;
}

const elements = {
  "courier-button": CourierButton,
  "courier-icon": CourierIcon,
  "courier-locale-selector": CourierLocaleSelector,
  "courier-panel": CourierPanel,
  "courier-progress": CourierProgress,
  "courier-theme-selector": CourierThemeSelector,
} as const;

export function defineCourierElements(registry: ElementRegistry = customElements): void {
  for (const [name, constructor] of Object.entries(elements)) {
    if (!registry.get(name)) {
      registry.define(name, constructor);
    }
  }
}
