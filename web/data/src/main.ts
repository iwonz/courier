import { CourierDataApp } from "./app";

if (!customElements.get("courier-data-app")) {
  customElements.define("courier-data-app", CourierDataApp);
}
