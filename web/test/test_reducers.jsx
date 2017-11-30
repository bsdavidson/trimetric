import {assert} from "chai";
import {createStore} from "redux";

import {reducer} from "../src/store";
import {updateVehicles} from "../src/actions";

describe("Reducers", function() {
  describe("vehicles", function() {
    it("takes action UPDATE_VEHICLES and updates state", function() {
      let store = createStore(reducer, {
        vehicles: "old vehicles data",
        vehiclesPointData: "old vehicles point data",
        vehiclesIconData: "old vehicles icon data"
      });
      store.dispatch(
        updateVehicles("vehicles", "vehiclesPointData", "vehiclesIconData")
      );
      let state = store.getState();
      assert.equal(state.vehicles, "vehicles");
      assert.equal(state.vehiclesPointData, "vehiclesPointData");
      assert.equal(state.vehiclesIconData, "vehiclesIconData");
    });
  });
});
