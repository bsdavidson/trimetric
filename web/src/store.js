import {createStore, combineReducers} from "redux";

import {
  LocationTypes,
  UPDATE_ARRIVALS,
  UPDATE_HOME_LOCATION,
  UPDATE_LINES,
  UPDATE_LOCATION,
  UPDATE_ROUTES,
  UPDATE_STOPS,
  UPDATE_VEHICLES,
  UPDATE_VIEWPORT
} from "./actions";

export const DEFAULT_LOCATION = {
  lat: 45.522236,
  lng: -122.675827,
  gps: false,
  locationType: LocationTypes.HOME
};

export const DEFAULT_BOUNDING_BOX = {
  sw: {
    lat: 45.50889931447199,
    lng: -122.68664166674807
  },
  ne: {
    lat: 45.53556952479618,
    lng: -122.66501233325198
  }
};

export const DEFAULT_ZOOM = 11;

function arrivals(state = [], action) {
  switch (action.type) {
    case UPDATE_ARRIVALS:
      return action.arrivals;
    default:
      return state;
  }
}

function boundingBox(state = DEFAULT_BOUNDING_BOX, action) {
  switch (action.type) {
    case UPDATE_VIEWPORT:
      return action.boundingBox;
    default:
      return state;
  }
}

function lineData(state = [], action) {
  switch (action.type) {
    case UPDATE_LINES:
      return action.lineData;
    default:
      return state;
  }
}

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

function routes(state = [], action) {
  switch (action.type) {
    case UPDATE_ROUTES:
      return action.routes;
    default:
      return state;
  }
}

function stops(state = [], action) {
  switch (action.type) {
    case UPDATE_STOPS:
      return action.stops;
    default:
      return state;
  }
}

function stopsPointData(state = [], action) {
  switch (action.type) {
    case UPDATE_STOPS:
      return action.stopsPointData;
    default:
      return state;
  }
}

function stopsIconData(state = [], action) {
  switch (action.type) {
    case UPDATE_STOPS:
      return action.stopsIconData;
    default:
      return state;
  }
}

function vehicles(state = [], action) {
  switch (action.type) {
    case UPDATE_VEHICLES:
      return action.vehicles;
    default:
      return state;
  }
}

function vehiclesIconData(state = [], action) {
  switch (action.type) {
    case UPDATE_VEHICLES:
      return action.vehiclesIconData;
    default:
      return state;
  }
}

function vehiclesPointData(state = [], action) {
  switch (action.type) {
    case UPDATE_VEHICLES:
      return action.vehiclesPointData;
    default:
      return state;
  }
}

function zoom(state = DEFAULT_ZOOM, action) {
  switch (action.type) {
    case UPDATE_VIEWPORT:
      return action.zoom;
    default:
      return state;
  }
}

export const reducer = combineReducers({
  arrivals,
  boundingBox,
  lineData,
  location,
  locationClicked,
  routes,
  stops,
  stopsIconData,
  stopsPointData,
  vehicles,
  vehiclesIconData,
  vehiclesPointData,
  zoom
});

export const store = createStore(
  reducer,
  window.devToolsExtension && window.devToolsExtension()
);
