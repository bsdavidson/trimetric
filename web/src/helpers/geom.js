export function withinBoundingBox(pos, boundingBox) {
  if (!boundingBox) {
    return true;
  }
  return (
    pos.lat <= boundingBox.ne.lat &&
    pos.lat >= boundingBox.sw.lat &&
    pos.lng <= boundingBox.ne.lng &&
    pos.lng >= boundingBox.sw.lng
  );
}

export function Point(lng, lat) {
  this.lng = lng || 0;
  this.lat = lat || 0;
}

export class Line {
  constructor(start, end) {
    this.start = start;
    this.end = end;

    const delta = new Point(end.lng - start.lng, end.lat - start.lat);
    const deltaLength = Math.sqrt(
      delta.lng * delta.lng + delta.lat * delta.lat
    );
    this.direction = new Point(
      delta.lng / deltaLength,
      delta.lat / deltaLength
    );
    this.normal = new Point(this.direction.lat, this.direction.lng);
    this.distance = -(
      this.normal.lng * this.start.lng +
      this.normal.lat * this.start.lat
    );
  }

  getDistance(point) {
    return Math.abs(
      this.normal.lng * point.lng + this.normal.lat * point.lat + this.distance
    );
  }
}

export function douglasPeucker(points, minDistance) {
  // Find the point with the maximum distance
  let furthestDistance = 0;
  let furthestIndex = -1;
  const endIndex = points.length - 1;
  const line = new Line(points[0], points[endIndex]);
  for (let i = 1; i <= endIndex - 1; i++) {
    let d = line.getDistance(points[i]);
    if (d > furthestDistance) {
      furthestIndex = i;
      furthestDistance = d;
    }
  }

  // If max distance is greater than minDistance, recursively simplify
  if (furthestDistance <= minDistance) {
    return [points[0], points[endIndex]];
  }
  // Recursive call
  let newPointsLeft = douglasPeucker(
    points.slice(0, furthestIndex + 1),
    minDistance
  );
  let newPointsRight = douglasPeucker(points.slice(furthestIndex), minDistance);
  // Build the result list
  return newPointsLeft.slice(0, -1).concat(newPointsRight);
}
