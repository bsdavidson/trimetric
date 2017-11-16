import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import {TrimetricPropTypes} from "./prop_types";
import {getVehicleType} from "../data";
import Marker from "./marker";

export class ArrivalMarkers extends React.Component {
  constructor(props) {
    super(props);
    this.markers = [];
  }

  getStopOpts(stop) {
    return {
      position: {
        lat: stop.lat,
        lng: stop.lng
      },
      name: stop.name,
      animation: this.props.google.maps.Animation.DROP
    };
  }

  getArrivalOpts(arrival) {
    return {
      position: {
        lat: arrival.vehicle_position.lat,
        lng: arrival.vehicle_position.lng
      },
      icon: {
        url: `./assets/${getVehicleType(arrival.route_type)}.png`,
        scaledSize: new this.props.google.maps.Size(25, 25)
      },
      opacity: 0.8,
      title: arrival.headsign
    };
  }

  render() {
    let {stop, google, map} = this.props;
    if (!stop || !stop.arrivals || !google || !map) {
      return null;
    }
    this.markers = stop.arrivals.filter(a => a.vehicle_id).map(a => {
      return (
        <Marker
          key={a.vehicle_id}
          google={google}
          map={map}
          opts={this.getArrivalOpts(a)}
        />
      );
    });
    this.markers.push(
      <Marker
        key={"stop-" + stop.id}
        google={this.props.google}
        map={this.props.map}
        opts={this.getStopOpts(this.props.stop)}
      />
    );
    return <div>{this.markers}</div>;
  }
}

ArrivalMarkers.propTypes = {
  google: TrimetricPropTypes.google,
  map: TrimetricPropTypes.map,
  stop: TrimetricPropTypes.stop
};

function mapStateToProps(state) {
  return {
    locationClicked: state.locationClicked
  };
}

export default withRouter(connect(mapStateToProps)(ArrivalMarkers));
