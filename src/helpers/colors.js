const COLORS = [
  "#16a085",
  "#27ae60",
  "#2980b9",
  "#9b59b6",
  "#e67e22",
  "#c0392b",
  "#f50057",
  "#51458d",
  "#486f15",
  "#725e43",
  "#191970",
  "#191919"
];

export class ColorMap {
  constructor() {
    this.map = {};
    this.mapLength = 0;
  }

  getColorForKey(key) {
    if (!this.map.hasOwnProperty(key)) {
      this.mapLength++;
      this.map[key] = COLORS[this.mapLength % COLORS.length];
    }
    return this.map[key];
  }
}
