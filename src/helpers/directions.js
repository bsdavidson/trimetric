import haversine from "haversine-distance";

const COMPASS_DIRECTIONS = [
  "N",
  "NNE",
  "NE",
  "ENE",
  "E",
  "ESE",
  "SE",
  "SSE",
  "S",
  "SSW",
  "SW",
  "WSW",
  "W",
  "WNW",
  "NW",
  "NNW"
];

export function degreeToCompass(degree) {
  var compassIndex = Math.floor(degree / 22.5 + 0.5);
  return COMPASS_DIRECTIONS[compassIndex % 16] || "";
}

export function formatDistance(start, end) {
  return Math.round(haversine(start, end) * 3.28084);
}
