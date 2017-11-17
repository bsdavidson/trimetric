import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import {TrimetricPropTypes} from "./prop_types";
import {getVehicleType} from "../data";
import {Marker} from "react-map-gl";

export class VehicleMarkers extends React.Component {
  constructor(props) {
    super(props);
    this.markers = [];
  }

  render() {
    let {vehicles} = this.props;
    if (!vehicles) {
      return null;
    }
    this.markers = vehicles.map(v => {
      return (
        <Marker
          key={v.vehicle.id}
          latitude={v.position.lat}
          longitude={v.position.lng}
          offsetLeft={-12}
          offsetTop={-12}>
          <img
            src={`/assets/${getVehicleType(v.route_type)}.png`}
            width={25}
            height={25}
          />
        </Marker>
      );
    });
    return <div>{this.markers}</div>;
  }
}

VehicleMarkers.propTypes = {
  vehicles: TrimetricPropTypes.vehiclePositions
};

function mapStateToProps(state) {
  return {
    locationClicked: state.locationClicked
  };
}

export default withRouter(connect(mapStateToProps)(VehicleMarkers));
