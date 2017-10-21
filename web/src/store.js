import {createStore, combineReducers} from "redux";

import {UPDATE_DATA, UPDATE_HOME_LOCATION, UPDATE_LOCATION} from "./actions";

const DEFAULT_LOCATION = {
  // WeWork
  lat: 45.5247402,
  lng: -122.6787931,
  gps: false
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

function vehicles(state = {}, action) {
  switch (action.type) {
    case UPDATE_DATA:
      return action.vehicles;
    default:
      return state;
  }
}

export const reducer = combineReducers({
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
