import "whatwg-fetch";
import fetchPonyfill from "fetch-ponyfill";
import moment from "moment"

import { updateData } from "./actions"
import { buildQuery } from "./helpers/http"


const fetch = fetchPonyfill(); // eslint-disable-line no-unused-vars

const UPDATE_TIMEOUT = 4000
const API_ENDPOINTS = {
  stops: "https://developer.trimet.org/ws/V1/stops",
  arrivals: "https://developer.trimet.org/ws/v2/arrivals",
  vehicles: "https://developer.trimet.org/ws/v2/vehicles",
  routes: "https://developer.trimet.org/ws/V1/routeConfig"
}

export class Trimet {
  constructor(store, fetch = fetch) {
    this.appId = process.env.TRIMET_API_KEY
    this.stopsCache = {
      location: {
        lat: null,
        lng: null
      },
      data: null
    }
    this.routesCache = {dir0: {}, dir1: {}}
    this.vehicleCache = {}
    this.lastVehicleQueryTime = 100
    this.store = store
    this.fetch = fetch
    this.running = false
    this.timeoutID = null
  }

  fetchStops(lat, lng) {
    if (this.stopsCache) {
      if (this.stopsCache.location && this.stopsCache.location.lat === lat && this.stopsCache.location.lng === lng) {
        return Promise.resolve(this.stopsCache.data)
      } else {
        this.stopsCache = {
          data: null,
          location: {lat: lat, lng: lng}
        }
      }
    }

    let stopsAPIURL = API_ENDPOINTS.stops + "?" + buildQuery({
      appID: this.appId,
      json: true,
      ll: lat + "," + lng,
      feet: 5000
    })
    return this.fetch(stopsAPIURL)
      .then(response => response.json())
  }

  // FIXME: Add protection against no returned Stops.
  fetchArrivals(stops) {
    if (!this.stopsCache.data) {
      this.stopsCache.data = {}
      this.stopsCache.data = Object.assign({}, this.stopsCache.data, stops)
    }
    let locations = stops.resultSet.location.map(location => location.locid)
    // Trimet limits locations in arrivals API to 10
    locations.length = 8
    let arrivalsAPIURL = API_ENDPOINTS.arrivals + "?" + buildQuery({
      appID: this.appId,
      json: true,
      locIDs: locations
    })
    return this.fetch(arrivalsAPIURL)
      .then(response => response.json())
  }

  fetchVehicles(arrivals, fetchAll = false) {
    let params = {}
    let vehicles = arrivals.resultSet.arrival.map(arrival => arrival.vehicleID).filter((v) => {
      if (v) {
        return v
      }
    })
    if (!vehicles.length) {
      return []
    }

    if (fetchAll) {
      params.appID = this.appId
      params.json = true
    } else {
      params.appID = this.appId
      params.json = true
      params.ids = vehicles
    }
    let vehiclesAPIURL = API_ENDPOINTS.vehicles + "?" + buildQuery(params)
    return this.fetch(vehiclesAPIURL)
      .then(response => response.json())
  }

  fetchRoute(route) {
    if (!route) {
      throw new Error("route argument cannot be undefined")
    }
    let params = {
      appID: this.appID,
      json: true,
      stops: true,
      route: route
    }

    let routesAPIURL = API_ENDPOINTS.routes + "?" + buildQuery(params)
    return this.fetch(routesAPIURL)
      .then(response => response.json())
  }

  fetchAllRoutes() {
    let params = {
      appID: this.appID,
      json: true,
      stops: true
    }

    let routesAPIURL = API_ENDPOINTS.routes + "?" + buildQuery(params)
    return this.fetch(routesAPIURL)
      .then(response => response.json())
  }

  fetchData(lat, lng) {
    let stops,
      arrivals
    return this.fetchStops(lat, lng)
      .then(s => {
        stops = s
        return this.fetchArrivals(s)
      })
      .then(a => {
        arrivals = a
        return this.fetchVehicles(a, true)
      })
      .then(vehicles => ({ stops, arrivals, vehicles }))
      .then(data => combineResponses(data.stops, data.arrivals, data.vehicles))
  }

  start() {
    this.running = true
    this.update()
  }

  stop() {
    this.running = false
    clearTimeout(this.timeoutID)
  }

  timeout() {
    clearTimeout(this.timeoutID)
    this.timeoutID = setTimeout(() => {
      this.update()
    }, UPDATE_TIMEOUT)
  }

  update() {
    let {lat, lng} = this.store.getState().location
    let newData = this.fetchData(lat, lng)
    newData.then(data => {
      this.store.dispatch(updateData(data.stops, data.vehicles, data.queryTime))
      this.timeout()
    })
  }
}

export function combineResponses(stops, arrivals, vehicles) {
  if (!stops) {
    throw new Error("stops argument cannot be undefined")
  }
  if (!arrivals) {
    throw new Error("arrivals argument cannot be undefined")
  }
  if (!vehicles) {
    throw new Error("vehicles argument cannot be undefined")
  }

  let response = stops.resultSet.location.map(location => {
    let loc = {
      lng: location.lng,
      lat: location.lat,
      locid: location.locid,
      desc: location.desc
    }
    loc.arrivals = arrivals.resultSet.arrival
      .filter(a => a.locid === location.locid && a.feet)
      .map(arrival => {
        let vehicle = vehicles.resultSet.vehicle
          .find(v => v.vehicleID == arrival.vehicleID)
        if (!vehicle) {
          vehicle = {
            latitude: 0,
            longitude: 0
          }
        }

        return {
          bearing: vehicle.bearing,
          estimated: arrival.estimated,
          feet: arrival.feet,
          id: arrival.id,
          latitude: vehicle.latitude,
          longitude: vehicle.longitude,
          route: arrival.route,
          scheduled: arrival.scheduled,
          shortSign: arrival.shortSign,
          signMessage: vehicle.signMessage,
          status: arrival.status,
          type: vehicle.type,
          vehicleID: arrival.vehicleID
        }
      })
    return loc
  })

  return {
    queryTime: moment(arrivals.resultSet.queryTime).valueOf(),
    stops: response,
    vehicles: {
      arrivals: vehicles.resultSet.vehicle
    }
  }
}
