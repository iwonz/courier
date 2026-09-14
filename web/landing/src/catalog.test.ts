import { expect, it } from "vitest";
import { landingLocale, landingText } from "./catalog";

it("resolves English and Russian landing catalogs", () => {
  expect(landingLocale(["de", "ru-RU"])).toBe("ru");
  expect(landingLocale(["de"])).toBe("en");
  expect(landingText("en", "install")).toBe("Install Courier");
  expect(landingText("ru", "install")).toBe("Установить Courier");
});
