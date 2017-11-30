export const UPDATE_ARRIVALS = "UPDATE_ARRIVALS";
export const UPDATE_VEHICLES = "UPDATE_VEHICLES";
export const UPDATE_VIEWPORT = "UPDATE_VIEWPORT";
export const UPDATE_LOCATION = "UPDATE_LOCATION";
export const UPDATE_LINES = "UPDATE_LINES";
export const UPDATE_STOPS = "UPDATE_STOPS";
export const UPDATE_ROUTES = "UPDATE_ROUTES";
export const UPDATE_HOME_LOCATION = "UPDATE_HOME_LOCATION";

export const LocationTypes = {
  VEHICLE: "VEHICLE",
  STOP: "STOP",
  HOME: "HOME"
};

export function updateVehicles(vehicles, vehiclesPointData, vehiclesIconData) {
  return {
    type: UPDATE_VEHICLES,
    vehicles,
    vehiclesPointData,
    vehiclesIconData
  };
}

export function updateArrivals(arrivals) {
  return {
    type: UPDATE_ARRIVALS,
    arrivals
  };
}

export function updateRoutes(routes) {
  return {
    type: UPDATE_ROUTES,
    routes
  };
}

export function updateLines(lineData) {
  return {
    type: UPDATE_LINES,
    lineData
  };
}

export function updateStops(stops, stopsPointData, stopsIconData) {
  return {
    type: UPDATE_STOPS,
    stops,
    stopsPointData,
    stopsIconData
  };
}

export function updateHomeLocation(lat, lng, gps = true) {
  return {
    type: UPDATE_HOME_LOCATION,
    home: {
      locationType: LocationTypes.HOME,
      lat,
      lng,
      gps
    }
  };
}

export function updateViewport(boundingBox, zoom) {
  return {
    type: UPDATE_VIEWPORT,
    boundingBox,
    zoom
  };
}

export function updateLocation(locationType, id, lat, lng, following) {
  return {
    type: UPDATE_LOCATION,
    locationClick: {
      locationType,
      id,
      lat,
      lng,
      following
    }
  };
}

export function clearLocation() {
  return {
    type: UPDATE_LOCATION,
    locationClick: null,
    following: false
  };
}
