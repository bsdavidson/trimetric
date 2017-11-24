// @ts-check

import "whatwg-fetch";
import fetchPonyfill from "fetch-ponyfill";
import moment from "moment";

import {LocationTypes, updateData, updateLocation} from "./actions";
import {buildQuery} from "./helpers/http";

const {fetch} = fetchPonyfill(); // eslint-disable-line no-unused-vars

const UPDATE_TIMEOUT = 1000;
const API_ENDPOINTS = {
  arrivals: "/api/v1/arrivals",
  shapes: "/api/v1/shapes",
  stops: "/api/v1/stops",
  routes: "/api/v1/routes",
  vehiclesGTFS: "/api/v1/vehicles"
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

export class Trimet {
  constructor(store, _fetch = fetch) {
    this.appId = process.env.TRIMET_API_KEY;
    this.stopsCache = {
      lat: null,
      lng: null,
      bbox: {
        sw: {},
        ne: {}
      },

      stops: null
    };
    this.fetchStopsID = 0;
    this.routesCache = {
      dir0: {},
      dir1: {}
    };
    this.shapesCache = {};
    this.routesCache = {};
    this.store = store;
    this.fetch = _fetch;
    this.running = false;
    this.timeoutID = null;
    this.selectedStop = null;

    this.handleStopChange = this.handleStopChange.bind(this);
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

  fetchStops(lat, lng, bbox) {
    //
    if (this.stopsCache.stops) {
      return Promise.resolve(this.stopsCache.stops);
    }

    return this.fetch(API_ENDPOINTS.stops)
      .then(response => response.json())
      .then(data => {
        this.stopsCache = {
          stops: data.stops
        };

        return data.stops;
      });
  }

  fetchRoutes() {
    if (this.routesCache.routes) {
      return Promise.resolve(this.routesCache.routes);
    }

    return this.fetch(API_ENDPOINTS.routes)
      .then(response => response.json())
      .then(data => {
        this.routesCache = {
          routes: data
        };

        return data;
      });
  }

  // FIXME: Add protection against no returned Stops.
  fetchArrivals(stop) {
    if (!stop) {
      throw new Error("stop is required");
    }

    let arrivalsAPIURL =
      API_ENDPOINTS.arrivals +
      "?" +
      buildQuery({
        locIDs: stop.id
      });
    return this.fetch(arrivalsAPIURL).then(response => response.json());
  }

  fetchVehicles() {
    return this.fetch(API_ENDPOINTS.vehiclesGTFS).then(response =>
      response.json()
    );
  }

  fetchShapes() {
    if (this.shapesCache.shapes) {
      return Promise.resolve(this.shapesCache.shapes);
    }

    let shapesAPIURL =
      API_ENDPOINTS.shapes +
      "?" +
      buildQuery({
        route_ids: ["100", "90", "190", "290", "193", "194", "195", "200"]
      });

    return this.fetch(shapesAPIURL)
      .then(response => response.json())
      .then(data => {
        this.shapesCache = {
          shapes: data
        };

        return data;
      });
  }

  fetchData(lat, lng, bbox) {
    let promises = [
      this.fetchStops(),
      this.fetchVehicles(),
      this.fetchShapes(),
      this.fetchRoutes()
    ];
    if (this.selectedStop) {
      promises.push(this.fetchArrivals(this.selectedStop));
    }

    return Promise.all(promises).then(results => {
      let [stops, vehicles, shapes, routes, arrivals] = results;

      return {stops, vehicles, shapes, routes, arrivals: arrivals || []};
    });
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
    let {lat, lng} = this.store.getState().location;
    let bbox = this.store.getState().boundingBox;

    let newData = this.fetchData(lat, lng, bbox);
    newData
      .then(data => {
        let geoJsonData = data.stops
          .map(s => ({
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
          }))
          .concat(
            data.vehicles.map(v => ({
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
            }))
          );

        let iconData = data.stops
          .map(s => ({
            position: [s.lng, s.lat, 0],
            icon: "stop",
            size: 1
          }))
          .concat(
            data.vehicles.map(v => ({
              position: [v.position.lng, v.position.lat, 10],
              icon: getVehicleType(v.route_type),
              size: 1.4
            }))
          );
        let routeData = {};
        data.routes.forEach(r => {
          routeData[r.id] = r;
        });

        let routeLineIndexes = {};
        let routeLines = [];
        data.shapes.forEach(s => {
          let idx = routeLineIndexes[s.route_id];
          if (idx === undefined) {
            idx = routeLines.length;
            routeLineIndexes[s.route_id] = idx;
            routeLines.push({
              type: "MultiLineString",
              color: routeData[s.route_id].color,
              coordinates: [],
              routeID: s.route_id
            });
          }
          routeLines[idx].coordinates.push(
            s.point.map(p => {
              return [p.lng, p.lat];
            })
          );
        });

        // console.log(data.shapes.length, routeLines);

        this.store.dispatch(
          updateData(
            data.stops,
            data.vehicles,
            data.arrivals,
            geoJsonData,
            iconData,
            routeLines,
            routeData
          )
        );
        let lc = this.store.getState().locationClicked;
        if (lc) {
          data.vehicles.forEach(v => {
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

        this.timeout();
      })
      .catch(err => {
        this.timeout();
        throw err;
      });
  }
}
