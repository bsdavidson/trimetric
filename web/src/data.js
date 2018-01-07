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

const UPDATE_TIMEOUT = 4000;
const API_ENDPOINTS = {
  arrivals: "/api/v1/arrivals",
  shapes: "/api/v1/shapes",
  stops: "/api/v1/stops",
  routes: "/api/v1/routes",
  lines: "/api/v1/routes/lines",
  vehiclesGTFS: "/api/v1/vehicles"
};

const messageTypeToAction = {
  arrivals: updateArrivals,
  routes: updateRoutes,
  route_shapes: updateLines,
  stops: updateStops,
  vehicles: updateVehicles
};

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
    const loc = window.location;
    let origin = loc.origin;
    if (!origin) {
      origin = `${loc.protocol}//${loc.hostname}${
        loc.port ? ":" + loc.port : ""
      }`;
    }

    const url = `${origin}/ws`.replace(/^http(s?)/, "ws$1");
    this.connection = new WebSocket(url);

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
    let action = messageTypeToAction[message.type];
    if (!action) {
      return;
    }
    this.store.dispatch(action(message.data));
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
      this.store.dispatch(updateLines(shapes));

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
        this.store.dispatch(updateStops(stops));
        return stops;
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
    this.store.dispatch(updateArrivals(arrivals));
  }

  processVehicles(vehicles) {
    this.store.dispatch(updateVehicles(vehicles));
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
