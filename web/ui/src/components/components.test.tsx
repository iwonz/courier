import * as React from "react";
import * as ScrollAreaPrimitive from "@radix-ui/react-scroll-area";
import { act, cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { BrandIcon, brandIconNames, resolveBrandIcon } from "../brand-icons-react";
import { cubicBezierPath } from "../geometry";
import { Icon, iconNames, resolveIcon } from "../icons-react";
import { cn } from "../lib/utils";
import { Brand, Mascot } from "./brand-react";
import { CommandReadout, copyText } from "./command-readout";
import { RouteDisplay } from "./route-display";
import { Alert, AlertDescription, AlertTitle } from "./ui/alert";
import { Badge } from "./ui/badge";
import { Button } from "./ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "./ui/card";
import { Checkbox } from "./ui/checkbox";
import { Input } from "./ui/input";
import { Progress } from "./ui/progress";
import { ScrollArea, ScrollBar } from "./ui/scroll-area";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "./ui/select";
import { Separator } from "./ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./ui/tabs";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./ui/tooltip";

afterEach(() => cleanup());

describe("shared shadcn primitives", () => {
  it("renders every static primitive and variant", () => {
    const { container } = render(<>
      <Alert variant="destructive" className="custom"><AlertTitle>Failure</AlertTitle><AlertDescription>Details</AlertDescription></Alert>
      <Alert variant="warning">Warning</Alert>
      <Alert>Default</Alert>
      {(["default", "secondary", "outline", "success", "warning", "destructive"] as const).map((variant) => <Badge key={variant} variant={variant}>{variant}</Badge>)}
      <Card className="custom"><CardHeader><CardTitle>Title</CardTitle><CardDescription>Description</CardDescription></CardHeader><CardContent>Content</CardContent><CardFooter>Footer</CardFooter></Card>
      <Input type="password" className="custom" aria-label="password" />
      <RouteDisplay source="Source" destination="Destination" className="custom" />
      <Separator /><Separator orientation="vertical" decorative={false} />
    </>);
    expect(screen.getAllByRole("alert")[0]!.className).toContain("custom");
    expect(screen.getByLabelText("password").getAttribute("type")).toBe("password");
    expect(screen.getByText("Destination").className).toContain("text-right");
    expect(container.querySelectorAll('[data-orientation="vertical"]')).toHaveLength(1);
  });

  it("supports button variants, sizes, refs, and slot composition", () => {
    const ref = React.createRef<HTMLButtonElement>();
    render(<>
      {(["default", "secondary", "outline", "ghost", "destructive"] as const).map((variant) => <Button key={variant} variant={variant}>{variant}</Button>)}
      {(["default", "sm", "lg", "icon", "icon-sm"] as const).map((size) => <Button key={size} size={size}>{size}</Button>)}
      <Button ref={ref} className="custom">ref</Button>
      <Button asChild><a href="#target">slot</a></Button>
    </>);
    expect(ref.current?.className).toContain("custom");
    expect(screen.getByRole("link", { name: "slot" }).className).toContain("inline-flex");
  });

  it("renders checkbox, bounded progress, scroll areas, tabs, and tooltip portals", () => {
    const { container, rerender } = render(<>
      <Checkbox aria-label="check" defaultChecked />
      <Progress value={150} />
      <Progress value={-5} />
      <Progress value={null} />
      <ScrollArea><span>Scrollable</span></ScrollArea>
      <ScrollAreaPrimitive.Root><ScrollBar orientation="horizontal" /></ScrollAreaPrimitive.Root>
      <Tabs defaultValue="a"><TabsList><TabsTrigger value="a">A</TabsTrigger><TabsTrigger value="b">B</TabsTrigger></TabsList><TabsContent value="a">Alpha</TabsContent><TabsContent value="b">Beta</TabsContent></Tabs>
      <TooltipProvider><Tooltip open><TooltipTrigger>Hint</TooltipTrigger><TooltipContent sideOffset={9} className="custom">Tooltip</TooltipContent></Tooltip></TooltipProvider>
    </>);
    const transforms = [...container.querySelectorAll<HTMLElement>('[role="progressbar"] > div')].map((node) => node.style.transform);
    expect(transforms).toEqual(["translateX(-0%)", "translateX(-100%)", "translateX(-100%)"]);
    expect(screen.getByText("Scrollable")).toBeTruthy();
    fireEvent.mouseDown(screen.getByRole("tab", { name: "B" }));
    rerender(<TooltipProvider><Tooltip open><TooltipTrigger>Hint</TooltipTrigger><TooltipContent>Tooltip</TooltipContent></Tooltip></TooltipProvider>);
    expect(screen.getByRole("tooltip")).toBeTruthy();
  });

  it("renders select trigger, content, groups, values, and items", () => {
    Object.defineProperty(Element.prototype, "hasPointerCapture", { configurable: true, value: () => false });
    Object.defineProperty(Element.prototype, "setPointerCapture", { configurable: true, value: () => undefined });
    Object.defineProperty(Element.prototype, "releasePointerCapture", { configurable: true, value: () => undefined });
    Object.defineProperty(Element.prototype, "scrollIntoView", { configurable: true, value: () => undefined });
    render(<Select defaultValue="a" open>
      <SelectTrigger><SelectValue /></SelectTrigger>
      <SelectContent position="item-aligned" className="custom"><SelectGroup><SelectItem value="a">Alpha</SelectItem><SelectItem value="b">Beta</SelectItem></SelectGroup></SelectContent>
    </Select>);
    expect(screen.getAllByText("Alpha").length).toBeGreaterThan(0);
    cleanup();
    render(<Select open><SelectTrigger><SelectValue placeholder="Choose" /></SelectTrigger><SelectContent position="popper"><SelectItem value="c">Charlie</SelectItem></SelectContent></Select>);
    expect(screen.getByText("Charlie")).toBeTruthy();
  });
});

describe("brand, icon, and geometry helpers", () => {
  it("renders compact and full brands plus the square mascot", () => {
    const { container } = render(<><Brand /><Brand compact className="compact" /><Mascot alt="Relay" className="mascot" /></>);
    expect(screen.getAllByText("COURIER CLI")).toHaveLength(1);
    expect(container.querySelector('img[src*="courier-relay-mark-v2"]')?.getAttribute("width")).toBe("512");
    const mascot = screen.getByAltText("Relay");
    expect(mascot.getAttribute("width")).toBe("768");
    expect(mascot.getAttribute("height")).toBe("768");
    expect(container.querySelector(".compact")).toBeTruthy();
  });

  it("resolves and renders every semantic icon and fallback", () => {
    expect(resolveIcon("unknown")).toBe("parcel");
    const { container } = render(<>{iconNames.map((name) => <Icon key={name} name={name} label={name} />)}<Icon name="unknown" /></>);
    expect(container.querySelectorAll("svg")).toHaveLength(iconNames.length + 1);
    expect(screen.getByLabelText("github")).toBeTruthy();
    expect(container.querySelector('svg[aria-hidden="true"]')).toBeTruthy();
  });

  it("resolves and renders official brand masks and wget fallback", () => {
    expect(resolveBrandIcon("missing")).toBe("linux");
    const { container } = render(<>{brandIconNames.map((name) => <BrandIcon key={name} name={name} />)}<BrandIcon name="missing" label="Fallback" className="custom" /></>);
    expect(screen.getByLabelText("GNU Wget").querySelector("svg")).toBeTruthy();
    expect(screen.getByLabelText("Fallback").className).toContain("custom");
    expect(container.querySelectorAll('[role="img"]')).toHaveLength(brandIconNames.length + 1);
  });

  it("merges utility classes and builds rounded finite bezier paths", () => {
    expect(cn("px-2", false && "hidden", "px-4")).toBe("px-4");
    expect(cubicBezierPath({ x: 0, y: 1.234 }, { x: 100, y: 50 })).toBe("M 0 1.23 C 42 1.23, 58 50, 100 50");
    expect(cubicBezierPath({ x: 100, y: 50 }, { x: 80, y: 20 })).toBe("M 100 50 C 68 50, 112 20, 80 20");
    expect(() => cubicBezierPath({ x: Number.NaN, y: 0 }, { x: 0, y: 0 })).toThrow(TypeError);
  });
});

describe("command readout", () => {
  it("copies, reports failure, resets by session, and hides optional regions", async () => {
    expect(await copyText("", undefined)).toBe(false);
    expect(await copyText("value", { writeText: vi.fn().mockRejectedValue(new Error("denied")) })).toBe(false);
    expect(await copyText("value", { writeText: vi.fn().mockResolvedValue(undefined) })).toBe(true);
    const copy = vi.fn().mockResolvedValueOnce(true).mockResolvedValueOnce(false);
    const { rerender } = render(<CommandReadout heading="Command" command="courier version" copyLabel="Copy" copiedLabel="Copied" copyFailedLabel="Failed" description="Description" details={<span>Details</span>} footerActions={<a href="#release">Release</a>} sessionKey="one" copy={copy} />);
    fireEvent.click(screen.getByRole("button", { name: "Copy" }));
    await waitFor(() => expect(screen.getByText("Copied")).toBeTruthy());
    fireEvent.click(screen.getByRole("button", { name: "Copy" }));
    await waitFor(() => expect(screen.getByText("Failed")).toBeTruthy());
    rerender(<CommandReadout heading="Command" command="" copyLabel="Copy" copiedLabel="Copied" copyFailedLabel="Failed" sessionKey="two" copy={copy} />);
    await act(async () => undefined);
    expect(screen.getByRole("button", { name: "Copy" }).hasAttribute("disabled")).toBe(true);
    expect(screen.queryByText("Description")).toBeNull();
    expect(screen.queryByText("Release")).toBeNull();
  });
});
