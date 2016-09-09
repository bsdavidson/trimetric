import { assert } from "chai"
import { createStore } from "redux"

import { reducer } from "../src/store"
import { updateData } from "../src/actions"

describe("Reducers", function() {
  describe("reducer", function() {
    it("takes action UPDATE_DATA and updates state", function() {
      let store = createStore(reducer, {
        stops: "OldData",
        vehicles: "OldData"
      })
      assert.equal(store.getState().stops, "OldData")
      store.dispatch(updateData("NewStops", "NewVehicles", 123))
      assert.equal(store.getState().stops, "NewStops")
      assert.equal(store.getState().vehicles, "NewVehicles")
      assert.equal(store.getState().queryTime, 123)

    })
  })
})
