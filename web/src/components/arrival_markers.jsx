import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import {TrimetricPropTypes} from "./prop_types";
import {getVehicleType} from "../data";
import {Marker} from "react-map-gl";

export class ArrivalMarkers extends React.Component {
  constructor(props) {
    super(props);
    this.markers = [];
  }

  render() {
    let {stop} = this.props;
    if (!stop || !stop.arrivals) {
      return null;
    }

    this.markers = stop.arrivals.filter(a => a.vehicle_id).map(a => {
      return (
        <Marker
          key={a.vehicle_id}
          latitude={a.vehicle_position.lat}
          longitude={a.vehicle_position.lng}
          offsetLeft={-12}
          offsetTop={-12}>
          <img
            width={25}
            height={25}
            src={`/assets/${getVehicleType(a.route_type)}.png`}
          />
        </Marker>
      );
    });

    this.markers.push(
      <Marker
        key={"stop-" + stop.id}
        latitude={stop.lat}
        longitude={stop.lng}
        offsetLeft={-12}
        offsetTop={-12}>
        <span className="fui-location" />
      </Marker>
    );
    return <div>{this.markers}</div>;
  }
}

ArrivalMarkers.propTypes = {
  stop: TrimetricPropTypes.stop
};

function mapStateToProps(state) {
  return {
    locationClicked: state.locationClicked
  };
}

export default withRouter(connect(mapStateToProps)(ArrivalMarkers));
