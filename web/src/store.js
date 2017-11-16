import {createStore, combineReducers} from "redux";

import {
  LocationTypes,
  UPDATE_DATA,
  UPDATE_BOUNDING_BOX,
  UPDATE_HOME_LOCATION,
  UPDATE_LOCATION
} from "./actions";

const DEFAULT_LOCATION = {
  lat: 45.522236,
  lng: -122.675827,
  gps: false,
  locationType: LocationTypes.HOME
};

const DEFAULT_BOUNDING_BOX = {
  south: 45.50889931447199,
  west: -122.68664166674807,
  north: 45.53556952479618,
  east: -122.66501233325198
};

function location(state = DEFAULT_LOCATION, action) {
  switch (action.type) {
    case UPDATE_HOME_LOCATION:
      return action.home;
    default:
      return state;
  }
}

function locationClicked(state = null, action) {
  switch (action.type) {
    case UPDATE_LOCATION:
      return action.locationClick;
    default:
      return state;
  }
}

function boundingBox(state = DEFAULT_BOUNDING_BOX, action) {
  switch (action.type) {
    case UPDATE_BOUNDING_BOX:
      return action.boundingBox;
    default:
      return state;
  }
}

function queryTime(state = null, action) {
  switch (action.type) {
    case UPDATE_DATA:
      return action.queryTime;
    default:
      return state;
  }
}

function stops(state = [], action) {
  switch (action.type) {
    case UPDATE_DATA:
      return action.stops;
    default:
      return state;
  }
}

function vehicles(state = [], action) {
  switch (action.type) {
    case UPDATE_DATA:
      return action.vehicles;
    default:
      return state;
  }
}

export const reducer = combineReducers({
  boundingBox,
  location,
  locationClicked,
  queryTime,
  stops,
  vehicles
});

export const store = createStore(
  reducer,
  window.devToolsExtension && window.devToolsExtension()
);
