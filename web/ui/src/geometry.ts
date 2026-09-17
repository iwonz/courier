export interface Point {
  readonly x: number;
  readonly y: number;
}

function coordinate(value: number): string {
  return Number(value.toFixed(2)).toString();
}

function assertFinite(source: Point, destination: Point): void {
  if (![source.x, source.y, destination.x, destination.y].every(Number.isFinite)) {
    throw new TypeError("Bezier coordinates must be finite");
  }
}

function controls(source: Point, destination: Point): readonly [Point, Point] {
  const distance = Math.abs(destination.x - source.x);
  const handle = Math.max(32, distance * 0.42);
  const direction = destination.x >= source.x ? 1 : -1;
  return [
    { x: source.x + handle * direction, y: source.y },
    { x: destination.x - handle * direction, y: destination.y },
  ];
}

export function cubicBezierPath(source: Point, destination: Point): string {
  assertFinite(source, destination);
  const [first, second] = controls(source, destination);
  return `M ${coordinate(source.x)} ${coordinate(source.y)} C ${coordinate(first.x)} ${coordinate(first.y)}, ${coordinate(second.x)} ${coordinate(second.y)}, ${coordinate(destination.x)} ${coordinate(destination.y)}`;
}

export function pixelBezierPath(source: Point, destination: Point, grid = 4, samples = 24): string {
  assertFinite(source, destination);
  if (!Number.isFinite(grid) || grid <= 0 || !Number.isInteger(samples) || samples < 2) {
    throw new TypeError("Pixel route grid and sample count must be valid");
  }
  const [first, second] = controls(source, destination);
  const snap = (value: number): number => Math.round(value / grid) * grid;
  const points: Point[] = [];
  for (let index = 0; index <= samples; index += 1) {
    const t = index / samples;
    const inverse = 1 - t;
    const next = {
      x: snap((inverse ** 3 * source.x) + (3 * inverse ** 2 * t * first.x) + (3 * inverse * t ** 2 * second.x) + (t ** 3 * destination.x)),
      y: snap((inverse ** 3 * source.y) + (3 * inverse ** 2 * t * first.y) + (3 * inverse * t ** 2 * second.y) + (t ** 3 * destination.y)),
    };
    const previous = points.at(-1);
    if (!previous || previous.x !== next.x || previous.y !== next.y) points.push(next);
  }
  const [firstPoint, ...rest] = points;
  let path = `M ${coordinate(firstPoint!.x)} ${coordinate(firstPoint!.y)}`;
  for (const point of rest) path += ` H ${coordinate(point.x)} V ${coordinate(point.y)}`;
  return path;
}
