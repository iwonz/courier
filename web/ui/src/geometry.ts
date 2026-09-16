export interface Point {
  readonly x: number;
  readonly y: number;
}

function coordinate(value: number): string {
  return Number(value.toFixed(2)).toString();
}

export function cubicBezierPath(source: Point, destination: Point): string {
  if (![source.x, source.y, destination.x, destination.y].every(Number.isFinite)) {
    throw new TypeError("Bezier coordinates must be finite");
  }
  const distance = Math.abs(destination.x - source.x);
  const handle = Math.max(32, distance * 0.42);
  const direction = destination.x >= source.x ? 1 : -1;
  const first = source.x + handle * direction;
  const second = destination.x - handle * direction;
  return `M ${coordinate(source.x)} ${coordinate(source.y)} C ${coordinate(first)} ${coordinate(source.y)}, ${coordinate(second)} ${coordinate(destination.y)}, ${coordinate(destination.x)} ${coordinate(destination.y)}`;
}
