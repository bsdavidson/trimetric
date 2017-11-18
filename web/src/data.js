// @ts-check

import "whatwg-fetch";
import fetchPonyfill from "fetch-ponyfill";
import moment from "moment";

import {LocationTypes, updateData, updateLocation} from "./actions";
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
    this.store = store;
    this.fetch = _fetch;
    this.running = false;
    this.timeoutID = null;
    this.selectedStopID = null;

    this.handleStopChange = this.handleStopChange.bind(this);
  }

  handleStopChange(stop) {
    if (stop) {
      this.selectedStopID = stop;
      this.store.dispatch(
        updateLocation(LocationTypes.STOP, stop.id, stop.lat, stop.lng, false)
      );
    } else {
      this.selectedStopID = null;
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

  fetchData(lat, lng, bbox) {
    let promises = [this.fetchStops(), this.fetchVehicles()];
    if (this.selectedStopID) {
      promises.push(this.fetchArrivals(this.selectedStopID));
    }

    return Promise.all(promises).then(results => {
      let [stops, vehicles, arrivals] = results;

      return {stops, vehicles, arrivals: arrivals || []};
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
              icon: getVehicleType(v.routeType),
              size: 1.4
            }))
          );

        this.store.dispatch(
          updateData(
            data.stops,
            data.vehicles,
            data.arrivals,
            geoJsonData,
            iconData
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
