import relayAdminSource from "../assets/courier-relay-pixel-admin-v3.webp";
import relayDeliverySource from "../assets/courier-relay-pixel-delivery-v3.webp";
import relayMarkSource from "../assets/courier-relay-pixel-mark-v3.webp";
import relayNeutralSource from "../assets/courier-relay-pixel-neutral-v3.webp";
import relayRouteSource from "../assets/courier-relay-pixel-route-v3.webp";

export const relayRoles = ["neutral", "route", "delivery", "admin"] as const;
export type RelayRole = typeof relayRoles[number];

const relaySources: Readonly<Record<RelayRole, string>> = {
  neutral: relayNeutralSource,
  route: relayRouteSource,
  delivery: relayDeliverySource,
  admin: relayAdminSource,
};

export function relaySourceForRole(role: RelayRole): string { return relaySources[role]; }

export { relayAdminSource, relayDeliverySource, relayMarkSource, relayNeutralSource, relayRouteSource };
