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
  vehicles: "/api/v1/vehicles",
  vehiclesGTFS: "/api/v2/vehicles"
};

export class Trimet {
  constructor(store, _fetch = fetch) {
    this.appId = process.env.TRIMET_API_KEY;
    this.stopsCache = {
      location: {
        lat: null,
        lng: null
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

  fetchStops(lat, lng) {
    this.fetchStopsID++;
    let id = this.fetchStopsID;
    if (
      this.stopsCache.location.lat === lat &&
      this.stopsCache.location.lng === lng
    ) {
      return Promise.resolve(this.stopsCache.stops);
    }

    let stopsAPIURL =
      API_ENDPOINTS.stops +
      "?" +
      buildQuery({
        lat: lat,
        lng: lng,
        distance: 200
      });
    return this.fetch(stopsAPIURL)
      .then(response => response.json())
      .then(data => {
        if (id === this.fetchStopsID) {
          this.stopsCache = {
            location: {
              lat: lat,
              lng: lng
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
    // Trimet limits locations in arrivals API to 10

    if (locations.length > 8) {
      locations.length = 8;
    }

    let arrivalsAPIURL =
      API_ENDPOINTS.arrivals +
      "?" +
      buildQuery({
        locIDs: locations.join(",")
      });
    return this.fetch(arrivalsAPIURL).then(response => response.json());
  }

  fetchVehicles(arrivals, fetchAll = false) {
    let params = {};
    let vehicles = arrivals.resultSet.arrival
      .map(arrival => arrival.vehicleID)
      .filter(v => {
        if (v) {
          return v;
        }
      });
    if (!vehicles.length) {
      return [];
    }

    if (!fetchAll) {
      params.ids = vehicles;
    }
    let vehiclesAPIURL = API_ENDPOINTS.vehiclesGTFS + "?" + buildQuery(params);
    return this.fetch(vehiclesAPIURL).then(response => response.json());
  }

  fetchData(lat, lng) {
    let stops, arrivals;
    return this.fetchStops(lat, lng)
      .then(s => {
        stops = s;
        return this.fetchArrivals(s);
      })
      .then(a => {
        arrivals = a;
        return this.fetchVehicles(a, true);
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

    let newData = this.fetchData(lat, lng);
    newData
      .then(data => {
        this.store.dispatch(
          updateData(data.stops, data.vehicles, data.queryTime)
        );
        let lc = this.store.getState().locationClicked;
        if (lc) {
          data.vehicles.arrivals.forEach(v => {
            if (lc.id !== v.vehicleID) {
              return;
            }
            if (lc.lat === v.latitude && lc.lng === v.longitude) {
              return;
            }
            this.store.dispatch(
              updateLocation(
                lc.locationType,
                lc.id,
                v.latitude,
                v.longitude,
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
    throw new Error("stops argument cannot be undefined");
  }
  if (!arrivals) {
    throw new Error("arrivals argument cannot be undefined");
  }
  if (!vehicles) {
    throw new Error("vehicles argument cannot be undefined");
  }

  let newStops = stops.map(stop => {
    let newStop = {
      lng: stop.lon,
      lat: stop.lat,
      locid: parseInt(stop.id, 10),
      desc: stop.name
    };
    newStop.arrivals = arrivals.resultSet.arrival
      .filter(a => +a.locid === +newStop.locid && a.feet)
      .map(arrival => {
        let vehicle = vehicles.find(v => +v.vehicle.id === +arrival.vehicleID);
        if (!vehicle) {
          vehicle = {
            position: {
              latitude: 0,
              longitude: 0
            },
            vehicle: {},
            trip: {}
          };
        }
        return {
          bearing: vehicle.position.bearing,
          estimated: arrival.estimated,
          feet: arrival.feet,
          id: arrival.id,
          latitude: vehicle.position.latitude,
          longitude: vehicle.position.longitude,
          route: arrival.route,
          scheduled: arrival.scheduled,
          shortSign: arrival.shortSign,
          signMessage: vehicle.vehicle.label,
          status: arrival.status,
          type: "bus",
          vehicleID: arrival.vehicleID
        };
      });
    return newStop;
  });
  return {
    queryTime: moment(arrivals.resultSet.queryTime).valueOf(),
    stops: newStops,
    vehicles: {
      arrivals: vehicles.map(v => {
        return {
          latitude: v.position.latitude,
          longitude: v.position.longitude,
          vehicleID: v.vehicle.id,
          type: "bus"
        };
      })
    }
  };
}
