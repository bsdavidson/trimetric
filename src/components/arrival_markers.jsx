import deepEqual from "deep-equal"
import React from "react"

import { updateLocation, LocationTypes } from "../actions"
import { Marker } from "./map"
import { store } from "../store"


const VEHICLE_ICON_MAP = {
  bus: "bus",
  rail: "tram"
}

export class ArrivalMarkers extends React.Component {
  constructor(props) {
    super(props)
  }

  getStopOpts(stop) {
    return {
      position: {
        lat: stop.lat,
        lng: stop.lng
      },
      name: stop.desc,
      animation: this.props.google.maps.Animation.DROP
    }
  }

  isFollowing(arrival) {
    return (this.props.state.locationClicked &&
            this.props.state.locationClicked.id === arrival.vehicleID);
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
    }
  }

  render() {
    let {stop, google, map} = this.props
    if (!stop || !google || !map) {
      return null
    }
    let markers = stop.arrivals.map(a => {
      let location = {
        locationType: LocationTypes.VEHICLE,
        id: a.vehicleID,
        lat: a.latitude,
        lng: a.longitude,
        following: true
      }
      if (this.isFollowing(a)) {
        if (!deepEqual(this.props.state.locationClicked, location)) {
          store.dispatch(updateLocation(
            location.locationType, location.id, location.lat,
            location.lng, location.following))
        }
      }
      return (
        <Marker key={a.vehicleID} google={google} map={map} opts={this.getVehicleOpts(a)} />
      )
    })

    markers.push(<Marker key={"stop-" + stop.id} google={google} map={map} opts={this.getStopOpts(stop)} />)

    return <div>{markers}</div>
  }
}

ArrivalMarkers.propTypes = {
  google: React.PropTypes.object,
  map: React.PropTypes.object,
  state: React.PropTypes.shape({
    locationClicked: React.PropTypes.shape({
      id: React.PropTypes.string
    })
  }).isRequired,
  stop: React.PropTypes.shape({
    arrivals: React.PropTypes.arrayOf(React.PropTypes.shape({
      latitude: React.PropTypes.number.isRequired,
      longitude: React.PropTypes.number.isRequired,
      // FIXME: this is either a string or number depending on where it is
      // used. It should really be only one of these.
      vehicleID: React.PropTypes.oneOfType([
        React.PropTypes.number.isRequired,
        React.PropTypes.string.isRequired
      ])
    })).isRequired,
    // FIXME: these should be required, but is optional because this component
    // is being passed the vehicles object as well as a stop, and the vehicles
    // object doesn't have an id.
    id: React.PropTypes.string,
    lat: React.PropTypes.number,
    lng: React.PropTypes.number,
    desc: React.PropTypes.string
  }).isRequired
}
