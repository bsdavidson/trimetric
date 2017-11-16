import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import {TrimetricPropTypes} from "./prop_types";
import {getVehicleType} from "../data";
import Marker from "./marker";

export class VehicleMarkers extends React.Component {
  constructor(props) {
    super(props);
    this.markers = [];
  }

  getVehicleOpts(vehicle) {
    return {
      position: {
        lat: vehicle.position.lat,
        lng: vehicle.position.lng
      },
      icon: {
        url: `./assets/${getVehicleType(vehicle.route_type)}.png`,
        scaledSize: new this.props.google.maps.Size(25, 25)
      },
      opacity: 0.8,
      title: vehicle.vehicle.label
    };
  }

  render() {
    let {vehicles, google, map} = this.props;
    if (!vehicles || !google || !map) {
      return null;
    }
    this.markers = vehicles.map(v => {
      return (
        <Marker
          key={v.vehicle.id}
          google={google}
          map={map}
          opts={this.getVehicleOpts(v)}
        />
      );
    });
    return <div>{this.markers}</div>;
  }
}

VehicleMarkers.propTypes = {
  google: TrimetricPropTypes.google,
  map: TrimetricPropTypes.map,
  vehicles: TrimetricPropTypes.vehiclePositions
};

function mapStateToProps(state) {
  return {
    locationClicked: state.locationClicked
  };
}

export default withRouter(connect(mapStateToProps)(VehicleMarkers));
