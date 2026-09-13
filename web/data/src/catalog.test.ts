import { expect, it } from "vitest";
import { dataLocale, dataText } from "./catalog";

it("resolves English and Russian data catalogs", () => {
  expect(dataLocale(["de", "ru-RU"])).toBe("ru");
  expect(dataLocale(["de"])).toBe("en");
  expect(dataText("en", "download")).toBe("Download");
  expect(dataText("ru", "download")).toBe("Скачать");
});
