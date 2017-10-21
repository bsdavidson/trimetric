import deepEqual from "deep-equal";
import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";
import PropTypes from "prop-types";

import {updateLocation, LocationTypes} from "../actions";
import Marker from "./marker";
import {store} from "../store";

const VEHICLE_ICON_MAP = {
  bus: "bus",
  rail: "tram"
};

export class ArrivalMarkers extends React.Component {
  constructor(props) {
    super(props);
  }

  getStopOpts(stop) {
    return {
      position: {
        lat: stop.lat,
        lng: stop.lng
      },
      name: stop.desc,
      animation: this.props.google.maps.Animation.DROP
    };
  }

  isFollowing(arrival) {
    return (
      this.props.locationClicked &&
      this.props.locationClicked.id === arrival.vehicleID
    );
  }

  getVehicleOpts(arrival) {
    return {
      position: {
        lat: arrival.latitude,
        lng: arrival.longitude
      },
      icon: {
        url: `./assets/${VEHICLE_ICON_MAP[arrival.type]}.png`,
        scaledSize: new this.props.google.maps.Size(25, 25)
      },
      opacity: 0.8,
      title: arrival.shortSign
    };
  }

  render() {
    let {stop, google, map} = this.props;
    if (!stop || !stop.arrivals || !google || !map) {
      return null;
    }
    let markers = stop.arrivals.map(a => {
      let location = {
        locationType: LocationTypes.VEHICLE,
        id: a.vehicleID,
        lat: a.latitude,
        lng: a.longitude,
        following: true
      };
      if (this.isFollowing(a)) {
        if (!deepEqual(this.props.locationClicked, location)) {
          store.dispatch(
            updateLocation(
              location.locationType,
              location.id,
              location.lat,
              location.lng,
              location.following
            )
          );
        }
      }
      return (
        <Marker
          key={a.vehicleID}
          google={google}
          map={map}
          opts={this.getVehicleOpts(a)}
        />
      );
    });

    markers.push(
      <Marker
        key={"stop-" + stop.id}
        google={google}
        map={map}
        opts={this.getStopOpts(stop)}
      />
    );

    return <div>{markers}</div>;
  }
}

ArrivalMarkers.propTypes = {
  google: PropTypes.object,
  map: PropTypes.object,
  stop: PropTypes.shape({
    arrivals: PropTypes.arrayOf(
      PropTypes.shape({
        latitude: PropTypes.number.isRequired,
        longitude: PropTypes.number.isRequired,
        // FIXME: this is either a string or number depending on where it is
        // used. It should really be only one of these.
        vehicleID: PropTypes.oneOfType([
          PropTypes.number.isRequired,
          PropTypes.string.isRequired
        ])
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
