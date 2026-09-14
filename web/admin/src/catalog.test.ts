import { expect, it } from "vitest";
import { adminLocale, adminText } from "./catalog";

it("resolves English and Russian administration catalogs", () => {
  expect(adminLocale(["de", "ru-RU"])).toBe("ru");
  expect(adminLocale(["de"])).toBe("en");
  expect(adminText("en", "title")).toBe("Delivery control");
  expect(adminText("ru", "title")).toBe("Управление доставками");
});
