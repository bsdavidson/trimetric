import {createStore, combineReducers} from "redux";
import {douglasPeucker} from "./helpers/geom.js";
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

const VEHICLE_TYPES = {
  0: "tram",
  1: "subway",
  2: "rail",
  3: "bus"
};

function mergeUpdates(state, updates, isEqualFunc) {
  let newState = state.slice();
  let newCount = 0;
  let expired = 0;
  let updateCount = 0;
  let expiredTimestamp = new Date().getTime() / 1000 - 300;
  updates.forEach(u => {
    for (let i = 0; i < newState.length; i++) {
      if (isEqualFunc(u, newState[i])) {
        newState[i] = u;
        updateCount++;
        let current = new Date().getTime() / 1000;

        return;
      }
      if (newState[i].timestamp < expiredTimestamp) {
        expired++;
        return;
      }
    }
    newCount++;
    newState.push(u);
  });
  // console.log(
  //   "NewCount",
  //   newCount,
  //   "Update",
  //   updateCount,
  //   "expired",
  //   expired,
  //   "Total",
  //   newState.length
  // );
  return newState;
}

export function getVehicleType(routeType) {
  return VEHICLE_TYPES[routeType] || "bus";
}

function parseColor(hex) {
  if (!hex) {
    return [0, 0, 0, 192];
  }
  let color = parseInt(hex, 16);
  return [(color >> 16) & 255, (color >> 8) & 255, color & 255, 255];
}

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
    case UPDATE_LINES: {
      let routeLineIndexes = {};
      let routeLines = [];
      action.lineData.forEach(s => {
        if (!s) {
          return;
        }
        let idx = routeLineIndexes[s.color];
        if (idx === undefined) {
          idx = routeLines.length;
          routeLineIndexes[s.color] = idx;
          routeLines.push({
            type: "MultiLineString",
            color: parseColor(s.color),
            coordinates: [],
            routeID: s.route_id,
            width: s.color ? 4 : 2
          });
        }
        routeLines[idx].coordinates.push(
          s.points.map(p => {
            return [p.lng, p.lat];
          })
        );
      });
      return state.concat(routeLines);
    }
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

    case UPDATE_VEHICLES:
      if (!state) {
        return state;
      }
      for (let i = 0; i < action.vehicles.length; i++) {
        let v = action.vehicles[i];
        if (+state.id !== +v.vehicle.id) {
          continue;
        }

        if (state.lat === v.position.lat && state.lng === v.position.lng) {
          continue;
        }
        return {
          locationType: state.locationType,
          id: state.id,
          lat: v.position.lat,
          lng: v.position.lng,
          following: state.following
        };
      }
      return state;
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
      return state.concat(...action.stops);

    default:
      return state;
  }
}

function stopsPointData(state = [], action) {
  switch (action.type) {
    case UPDATE_STOPS: {
      let stops = action.stops.map(s => ({
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
      }));
      let newState = state.concat(...stops);
      return newState;
    }
    default:
      return state;
  }
}

function stopsIconData(state = [], action) {
  switch (action.type) {
    case UPDATE_STOPS: {
      let stops = action.stops.map(s => ({
        position: [s.lng, s.lat, 0],
        icon: "stop",
        size: 1
      }));

      return state.concat(...stops);
    }
    default:
      return state;
  }
}

function vehicles(state = [], action) {
  switch (action.type) {
    case UPDATE_VEHICLES: {
      let newState = mergeUpdates(state, action.vehicles, (a, b) => {
        return a.vehicle.id === b.vehicle.id;
      });

      return newState;
    }
    default:
      return state;
  }
}

function vehiclesIconData(state = [], action) {
  switch (action.type) {
    case UPDATE_VEHICLES: {
      let vehiclesIconData = action.vehicles.map(v => ({
        position: [v.position.lng, v.position.lat, 10],
        icon: getVehicleType(v.route_type),
        size: 1.4,
        vehicle_id: v.vehicle.id,
        timestamp: v.timestamp
      }));

      let newState = mergeUpdates(state, vehiclesIconData, (a, b) => {
        return a.vehicle_id === b.vehicle_id;
      });

      return newState;
    }
    default:
      return state;
  }
}

function vehiclesPointData(state = [], action) {
  switch (action.type) {
    case UPDATE_VEHICLES: {
      let vehiclesPointData = action.vehicles.map(v => ({
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [v.position.lng, v.position.lat]
        },
        properties: {
          lineColor: [246, 76, 0, 255],
          fillColor: [246, 76, 0, 255],
          radius: 2.5
        },
        vehicle_id: v.vehicle.id,
        timestamp: v.timestamp
      }));

      let newState = mergeUpdates(state, vehiclesPointData, (a, b) => {
        return a.vehicle_id === b.vehicle_id;
      });

      return newState;
    }
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
