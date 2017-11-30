// @ts-check

import "whatwg-fetch";
import fetchPonyfill from "fetch-ponyfill";

import {
  LocationTypes,
  updateArrivals,
  updateLines,
  updateLocation,
  updateRoutes,
  updateStops,
  updateVehicles
} from "./actions";
import {buildQuery} from "./helpers/http";

const {fetch} = fetchPonyfill(); // eslint-disable-line no-unused-vars

const UPDATE_TIMEOUT = 1000;
const API_ENDPOINTS = {
  arrivals: "/api/v1/arrivals",
  shapes: "/api/v1/shapes",
  stops: "/api/v1/stops",
  routes: "/api/v1/routes",
  lines: "/api/v1/routes/lines",
  vehiclesGTFS: "/api/v1/vehicles",
  ws: "ws://localhost:8181/api/ws"
};

const VEHICLE_TYPES = {
  0: "tram",
  1: "subway",
  2: "rail",
  3: "bus"
};

export function getVehicleType(routeType) {
  return VEHICLE_TYPES[routeType] || "bus";
}

function parseColor(hex) {
  if (!hex) {
    return [0, 0, 0, 192];
  }
  let color = parseInt(hex, 16);
  return [(color >> 16) & 255, (color >> 8) & 255, color & 255, 255];
}

export class Trimet {
  constructor(store, _fetch = fetch) {
    this.store = store;
    this.fetch = _fetch;
    this.running = false;
    this.timeoutID = null;
    this.selectedStop = null;

    this.handleStopChange = this.handleStopChange.bind(this);
    this.connect();
  }

  connect() {
    this.connection = new WebSocket(API_ENDPOINTS.ws);

    this.connection.onopen = () => {
      console.log("WebSocket Connected");
    };

    this.connection.onclose = () => {
      console.log("WebSocket Disconnected");
    };

    this.connection.onmessage = message => {
      try {
        var parsedMsg = JSON.parse(message.data);
      } catch (err) {
        console.log("WebSocket JSON Error:", err);
        return;
      }
      this.handleMessage(parsedMsg);
    };
    this.connection.onerror = err => {
      console.log("WebSocket Error:", err);
    };
  }

  handleMessage(message) {
    switch (message.type) {
      case "arrivals":
        this.processArrivals(message.data);
        break;
      case "routes":
        this.processRoutes(message.data);
        break;
      case "route_shapes":
        this.processRouteShapes(message.data);
        break;
      case "stops":
        this.processStops(message.data);
        break;
      case "vehicles":
        this.processVehicles(message.data);
        break;
      default:
        break;
    }
  }

  handleStopChange(stop) {
    if (stop) {
      this.selectedStop = stop;
      this.store.dispatch(
        updateLocation(LocationTypes.STOP, stop.id, stop.lat, stop.lng, false)
      );
    } else {
      this.selectedStop = null;
    }
  }

  fetchArrivals(stop) {
    if (!stop) {
      throw new Error("stop is required");
    }

    let arrivalsAPIURL =
      API_ENDPOINTS.arrivals +
      "?" +
      buildQuery({
        stop_ids: stop.id
      });
    return this.fetch(arrivalsAPIURL)
      .then(response => response.json())
      .then(arrivals => {
        this.processArrivals(arrivals);
        return arrivals;
      });
  }

  fetchRoutes() {
    return this.fetch(API_ENDPOINTS.routes)
      .then(response => response.json())
      .then(data => {
        this.store.dispatch(updateRoutes(data));
        return data;
      });
  }

  fetchRouteShapes() {
    let shapesAPIURL =
      API_ENDPOINTS.lines +
      "?" +
      buildQuery({
        route_ids: ["100", "90", "190", "290", "193", "194", "195", "200"]
      });

    return Promise.all([
      this.fetch(shapesAPIURL).then(response => response.json())
    ]).then(([shapes]) => {
      this.processRouteShapes(shapes);
      return shapes;
    });
  }

  fetchStops() {
    return this.fetch(API_ENDPOINTS.stops)
      .then(response => response.json())
      .then(data => {
        return data.stops;
      })
      .then(stops => {
        this.processStops(stops);
      });
  }

  fetchVehicles() {
    return this.fetch(API_ENDPOINTS.vehiclesGTFS)
      .then(response => response.json())
      .then(r => {
        this.vehiclesCache = r;
        return r;
      })
      .then(vehicles => {
        this.processVehicles(vehicles);
        return vehicles;
      });
  }

  processArrivals(arrivals) {
    if (arrivals.length === 0) {
      return;
    }
    this.store.dispatch(updateArrivals(arrivals));
  }

  processRoutes(routes) {
    this.store.dispatch(updateRoutes(routes));
  }

  processRouteShapes(shapes) {
    let routeLineIndexes = {};
    let routeLines = [];
    shapes.forEach(s => {
      let idx = routeLineIndexes[s.color];
      if (idx === undefined) {
        idx = routeLines.length;
        routeLineIndexes[s.color] = idx;
        routeLines.push({
          type: "MultiLineString",
          color: parseColor(s.color),
          coordinates: [],
          routeID: s.route_id,
          width: s.color ? 4 : 2
        });
      }
      routeLines[idx].coordinates.push(
        s.points.map(p => {
          return [p.lng, p.lat];
        })
      );
    });
    this.store.dispatch(updateLines(routeLines));
  }

  processStops(stops) {
    let stopsPointData = stops.map(s => ({
      type: "Feature",
      geometry: {
        type: "Point",
        coordinates: [s.lng, s.lat]
      },
      properties: {
        lineColor: [255, 0, 0, 255],
        fillColor: [0, 0, 0, 255],
        radius: 1
      }
    }));

    let stopsIconData = stops.map(s => ({
      position: [s.lng, s.lat, 0],
      icon: "stop",
      size: 1
    }));

    this.store.dispatch(updateStops(stops, stopsPointData, stopsIconData));
    return stops;
  }

  processVehicles(vehicles) {
    let vehiclesPointData = vehicles.map(v => ({
      type: "Feature",
      geometry: {
        type: "Point",
        coordinates: [v.position.lng, v.position.lat]
      },
      properties: {
        lineColor: [246, 76, 0, 255],
        fillColor: [246, 76, 0, 255],
        radius: 2.5
      }
    }));

    let vehiclesIconData = vehicles.map(v => ({
      position: [v.position.lng, v.position.lat, 10],
      icon: getVehicleType(v.route_type),
      size: 1.4
    }));

    this.store.dispatch(
      updateVehicles(vehicles, vehiclesPointData, vehiclesIconData)
    );
    let lc = this.store.getState().locationClicked;
    if (lc) {
      vehicles.forEach(v => {
        if (lc.id !== v.vehicle.id) {
          return;
        }
        if (lc.lat === v.position.lat && lc.lng === v.position.lng) {
          return;
        }
        this.store.dispatch(
          updateLocation(
            lc.locationType,
            lc.id,
            v.position.lat,
            v.position.lng,
            lc.following
          )
        );
      });
    }
  }

  start() {
    this.running = true;
    this.update();
  }

  stop() {
    this.running = false;
    clearTimeout(this.timeoutID);
  }

  timeout() {
    clearTimeout(this.timeoutID);
    this.timeoutID = setTimeout(() => {
      this.update();
    }, UPDATE_TIMEOUT);
  }

  update() {
    let promises = [];
    if (this.selectedStop) {
      promises.push(this.fetchArrivals(this.selectedStop));
    }
    return Promise.all(promises)
      .then(results => {
        this.timeout();
        let [arrivals] = results;
        return {arrivals: arrivals || []};
      })
      .catch(err => {
        this.timeout();
        throw err;
      });
  }
}
