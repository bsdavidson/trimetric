export const UPDATE_DATA = "UPDATE_DATA";
export const UPDATE_VIEWPORT = "UPDATE_VIEWPORT";
export const UPDATE_LOCATION = "UPDATE_LOCATION";
export const UPDATE_HOME_LOCATION = "UPDATE_HOME_LOCATION";

export const LocationTypes = {
  VEHICLE: "VEHICLE",
  STOP: "STOP",
  HOME: "HOME"
};

export function updateData(stops, vehicles, arrivals, geoJsonData, iconData) {
  return {
    type: UPDATE_DATA,
    stops,
    vehicles,
    arrivals,
    geoJsonData,
    iconData
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
