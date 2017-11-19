import {assert} from "chai";
import {createStore} from "redux";

import {reducer} from "../src/store";
import {updateData} from "../src/actions";

describe("Reducers", function() {
  describe("reducer", function() {
    it("takes action UPDATE_DATA and updates state", function() {
      let store = createStore(reducer, {
        stops: "OldData",
        vehicles: "OldData"
      });
      assert.equal(store.getState().stops, "OldData");
      store.dispatch(
        updateData("stops", "vehicles", "arrivals", "geoJsonData", "iconData")
      );
      let state = store.getState();
      assert.equal(state.stops, "stops");
      assert.equal(state.vehicles, "vehicles");
      assert.equal(state.arrivals, "arrivals");
      assert.equal(state.geoJsonData, "geoJsonData");
      assert.equal(state.iconData, "iconData");
    });
  });
});
