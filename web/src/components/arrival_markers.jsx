import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";
import PropTypes from "prop-types";

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
  google: PropTypes.object,
  map: PropTypes.object,
  stop: PropTypes.shape({
    arrivals: PropTypes.arrayOf(
      PropTypes.shape({
        vehiclePosition: PropTypes.shape({
          vehicle: PropTypes.shape({
            id: PropTypes.string.isRequired
          }),
          position: PropTypes.shape({
            latitude: PropTypes.number.isRequired,
            longitude: PropTypes.number.isRequired
          })
        })
      })
    ).isRequired,
    // FIXME: these should be required, but is optional because this component
    // is being passed the vehicles object as well as a stop, and the vehicles
    // object doesn't have an id.
    id: PropTypes.string,
    lat: PropTypes.number,
    lng: PropTypes.number,
    desc: PropTypes.string
  }).isRequired
};

function mapStateToProps(state) {
  return {
    locationClicked: state.locationClicked
  };
}

export default withRouter(connect(mapStateToProps)(ArrivalMarkers));
