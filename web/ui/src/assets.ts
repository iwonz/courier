import relayAdminSource from "../assets/courier-relay-pixel-admin-v1.webp";
import relayDeliverySource from "../assets/courier-relay-pixel-delivery-v1.webp";
import relayMarkSource from "../assets/courier-relay-pixel-mark-v1.webp";
import relayNeutralSource from "../assets/courier-relay-pixel-neutral-v1.webp";
import relayRouteSource from "../assets/courier-relay-pixel-route-v1.webp";

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
