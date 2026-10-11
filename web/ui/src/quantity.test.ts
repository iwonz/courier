import { describe, expect, it } from "vitest";
import { formatBitRate, formatBits, formatByteRate, formatBytes, formatBytesDetailed, formatExactBytes } from "./quantity";

describe("data quantity formatting", () => {
  it("keeps bytes and bits explicit and groups base units", () => {
    expect(formatBytes(411_789)).toBe("402.1 KiB");
    expect(formatExactBytes(411_789)).toBe("411,789 B");
    expect(formatExactBytes(411_789, "ru")).toBe("411 789 B");
    expect(formatBits(1023)).toBe("1,023 bit");
    expect(formatBits(1024)).toBe("1 Kibit");
  });

  it("scales through the largest binary units and formats rates", () => {
    expect(formatBytes(1024 ** 7)).toBe("1,024 EiB");
    expect(formatByteRate(1536)).toBe("1.5 KiB/s");
    expect(formatBitRate(1536, "ru")).toBe("1,5 Kibit/s");
  });

  it("adds exact context only for scaled values and rejects invalid quantities", () => {
    expect(formatBytesDetailed(7)).toBe("7 B");
    expect(formatBytesDetailed(1536)).toBe("1.5 KiB (1,536 B)");
    expect(formatBytes(-1)).toBe("—");
    expect(formatExactBytes(Number.NaN)).toBe("—");
    expect(formatByteRate(Number.POSITIVE_INFINITY)).toBe("—");
    expect(formatBitRate(-1)).toBe("—");
  });
});
