// @ts-check

import "whatwg-fetch";
import fetchPonyfill from "fetch-ponyfill";
import moment from "moment";

import {updateData, updateLocation} from "./actions";
import {buildQuery} from "./helpers/http";

const {fetch} = fetchPonyfill(); // eslint-disable-line no-unused-vars

const UPDATE_TIMEOUT = 1000;
const API_ENDPOINTS = {
  stops: "/api/v1/stops",
  arrivals: "/api/v1/arrivals",
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
      location: {
        lat: null,
        lng: null,
        bbox: {
          sw: {},
          ne: {}
        }
      },
      stops: null
    };
    this.fetchStopsID = 0;
    this.routesCache = {
      dir0: {},
      dir1: {}
    };
    this.store = store;
    this.fetch = _fetch;
    this.running = false;
    this.timeoutID = null;
  }

  fetchStops(lat, lng, bbox) {
    this.fetchStopsID++;
    let id = this.fetchStopsID;

    if (
      this.stopsCache.location.lat === lat &&
      this.stopsCache.location.lng === lng &&
      (bbox && this.stopsCache.location.bbox.sw.lat === bbox.sw.lat)
    ) {
      return Promise.resolve(this.stopsCache.stops);
    }

    let stopsAPIURL =
      API_ENDPOINTS.stops +
      "?" +
      buildQuery({
        lat: lat,
        lng: lng,
        distance: 200,
        south: bbox.sw.lat,
        west: bbox.sw.lng,
        north: bbox.ne.lat,
        east: bbox.ne.lng
      });
    return this.fetch(stopsAPIURL)
      .then(response => response.json())
      .then(data => {
        if (id === this.fetchStopsID) {
          this.stopsCache = {
            location: {
              lat: lat,
              lng: lng,
              bbox: bbox
            },
            stops: data.stops
          };
        }
        return data.stops;
      });
  }

  // FIXME: Add protection against no returned Stops.
  fetchArrivals(stops) {
    if (!stops) {
      throw new Error("no stops");
    }
    let locations = stops.map(location => +location.id);

    let arrivalsAPIURL =
      API_ENDPOINTS.arrivals +
      "?" +
      buildQuery({
        locIDs: locations.join(",")
      });
    return this.fetch(arrivalsAPIURL).then(response => response.json());
  }

  fetchVehicles() {
    return this.fetch(API_ENDPOINTS.vehiclesGTFS).then(response =>
      response.json()
    );
  }

  fetchData(lat, lng, bbox) {
    let stops, arrivals;
    return this.fetchStops(lat, lng, bbox)
      .then(s => {
        stops = s;
        return this.fetchArrivals(s);
      })
      .then(a => {
        arrivals = a;
        return this.fetchVehicles();
      })
      .then(vehicles => ({
        stops,
        arrivals,
        vehicles
      }))
      .then(data => {
        return combineResponses(data.stops, data.arrivals, data.vehicles);
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
        this.store.dispatch(
          updateData(data.stops, data.vehicles, data.queryTime)
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

export function combineResponses(stops, arrivals, vehicles) {
  if (!stops) {
    stops = [];
  }
  if (!arrivals) {
    arrivals = [];
  }
  if (!vehicles) {
    vehicles = [];
  }

  let newStops = stops.map(stop => {
    let newStop = Object.assign({}, stop);
    newStop.arrivals = arrivals
      .filter(a => a.stop_id === newStop.id) // && a.feet)
      .map(arrival => {
        return Object.assign({}, arrival, {
          estimated: moment(arrival.date, "YYYY-MM-DD")
            .add(moment.duration(arrival.arrival_time))
            .valueOf(),
          feet: 123
        });
      });
    return newStop;
  });
  return {
    queryTime: moment().valueOf(),
    stops: newStops,
    vehicles: vehicles
  };
}
