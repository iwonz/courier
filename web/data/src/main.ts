import { courierMarkSource } from "@courier/ui";
import { CourierDataApp } from "./app";

const favicon = document.createElement("link");
favicon.rel = "icon";
favicon.type = "image/webp";
favicon.href = courierMarkSource;
document.head.append(favicon);

if (!customElements.get("courier-data-app")) {
  customElements.define("courier-data-app", CourierDataApp);
}
