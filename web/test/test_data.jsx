import {assert} from "chai";
import sinon from "sinon";
import {Trimet} from "../src/data";
import {getVehicleType} from "../src/store";
import {getMockStopsResponse} from "./mock_data";

let sandbox = sinon.createSandbox();

describe("Data", () => {
  afterEach(function() {
    sandbox.restore();
  });

  describe("getVehicleType", () => {
    it("should return a string for a routeType", function() {
      assert.equal(getVehicleType(0), "tram");
      assert.equal(getVehicleType(1), "subway");
      assert.equal(getVehicleType(2), "rail");
      assert.equal(getVehicleType(3), "bus");
      assert.equal(getVehicleType(1337), "bus");
    });
  });
  describe("handleStopChange", () => {
    it("should update location with a stop", function() {
      let dispatch = sinon.spy();
      let trimet = new Trimet({dispatch});

      trimet.handleStopChange(null);
      assert.isFalse(dispatch.called);
      assert.isNull(trimet.selectedStop);

      let stop = getMockStopsResponse()[0];
      trimet.handleStopChange(stop);
      assert.isTrue(dispatch.calledOnce);
      assert.deepEqual(dispatch.getCall(0).args, [
        {
          locationClick: {
            following: false,
            id: stop.id,
            lat: stop.lat,
            lng: stop.lng,
            locationType: "STOP"
          },
          type: "UPDATE_LOCATION"
        }
      ]);
      assert.equal(trimet.selectedStop, stop);
    });
  });

  describe("start", () => {
    it("should start the update loop", function() {
      let trimet = new Trimet();
      sandbox.stub(trimet, "update");

      trimet.start();
      assert.isTrue(trimet.running);
      assert.isTrue(trimet.update.calledOnce);
    });
  });

  describe("stop", () => {
    it("should stop the update loop", function() {
      let trimet = new Trimet();
      sandbox.stub(window, "clearTimeout");
      trimet.running = true;
      trimet.timeoutID = 1;
      trimet.stop();
      assert.isFalse(trimet.running);
      assert.isTrue(clearTimeout.calledOnce);
      assert.equal(clearTimeout.getCall(0).args[0], 1);
    });
  });
});
