import {assert} from "chai";
import {createStore} from "redux";

import {reducer} from "../src/store";
import {updateVehicles} from "../src/actions";

describe("Reducers", function() {
  describe("vehicles", function() {
    it("takes action UPDATE_VEHICLES and updates state", function() {
      let store = createStore(reducer, {
        vehicles: [
          {
            vehicle: {
              id: 1
            },
            position: {
              lat: 123,
              lng: 456
            }
          }
        ],
        vehiclesPointData: [
          {
            position: {
              lat: 123,
              lng: 456
            }
          }
        ],
        vehiclesIconData: [
          {
            position: {
              lat: 123,
              lng: 456
            }
          }
        ]
      });
      store.dispatch(
        updateVehicles(
          [
            {
              vehicle: {
                id: 2
              },
              position: {
                lat: 222,
                lng: 444
              }
            }
          ],
          [
            {
              position: {
                lat: 555,
                lng: 666
              }
            }
          ],
          [
            {
              position: {
                lat: 555,
                lng: 666
              }
            }
          ]
        )
      );
      let state = store.getState();
      assert.deepEqual(state.vehicles, [
        {
          vehicle: {
            id: 1
          },
          position: {
            lat: 123,
            lng: 456
          }
        },
        {
          vehicle: {
            id: 2
          },
          position: {
            lat: 222,
            lng: 444
          }
        }
      ]);
    });
  });
});
