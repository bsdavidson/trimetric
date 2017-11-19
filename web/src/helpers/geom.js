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
