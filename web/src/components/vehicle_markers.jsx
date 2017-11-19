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
    let {vehicles, boundingBox} = this.props;
    if (!vehicles) {
      return null;
    }

    return (
      <div>
        {vehicles
          .filter(
            v =>
              v.position.lat <= boundingBox.ne.lat &&
              v.position.lat >= boundingBox.sw.lat &&
              v.position.lng <= boundingBox.ne.lng &&
              v.position.lng >= boundingBox.sw.lng
          )
          .map(v => (
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
          ))}
      </div>
    );
  }
}

VehicleMarkers.propTypes = {
  vehicles: TrimetricPropTypes.vehiclePositions
};

function mapStateToProps(state) {
  return {
    locationClicked: state.locationClicked,
    boundingBox: state.boundingBox
  };
}

export default withRouter(connect(mapStateToProps)(VehicleMarkers));
