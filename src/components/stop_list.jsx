import React from "react"
import { Link } from "react-router"

import { updateHomeLocation } from "../actions"
import { formatDistance } from "../helpers/directions"
import { humanTimeUntilEpoch } from "../helpers/times"
import { store } from "../store"


export class StopList extends React.Component {
  constructor(props) {
    super(props)
  }

  createStops(stops) {
    stops = stops
      .filter(s => s.arrivals.length > 0)
      .map(s => <StopListItem key={s.locid} stop={s} state={this.props.state} />)
    if (!stops.length) {
      return "Sorry, no buses are running near you. Better start walking or call an Uber."
    }
    return stops
  }

  handleCurrentLocationClick() {
    if ("geolocation" in navigator) {
      navigator.geolocation.watchPosition((position) => {
        store.dispatch(updateHomeLocation(position.coords.latitude, position.coords.longitude))
      })
    } else {
      alert("Your browser doesn't support geolocation")
    }
  }

  render() {
    return (
      <div className="stop-list">
        <div className="use-current" title="Use Current GPS position" onClick={this.handleCurrentLocationClick}>
          <span className="fui-location"></span>
        </div>
        <h3>Nearby Stops</h3>
        <div className="stop-list-items">
          {this.createStops(this.props.stops)}
        </div>
      </div>
    )
  }
}

StopList.propTypes = {
  state: React.PropTypes.object,
  stops: React.PropTypes.arrayOf(React.PropTypes.shape({
    arrivals: React.PropTypes.array.isRequired,
    locid: React.PropTypes.number.isRequired
  })).isRequired
}

export class StopListItem extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    let location = {
      latitude: this.props.stop.lat,
      longitude: this.props.stop.lng
    }
    let arrivalInfo = {
      route: this.props.stop.arrivals[0].route,
      estimated: humanTimeUntilEpoch(this.props.stop.arrivals[0].estimated)
    }
    return (
      <Link className="stop-link" to={`/stop/${this.props.stop.locid}`}>
        <div className="stop-list-item">
          <div className="stop-arrow">
            <span className="fui-arrow-right"></span>
          </div>
          <h4 className="stop-description">
            {this.props.stop.desc} -
            #<span className="stop-id">{this.props.stop.locid}</span>
          </h4>
          <p className="stop-metrics">
            <span className="stop-distance">{formatDistance(this.props.state.location, location)}</span> feet away.
            <span className="stop-route"> {arrivalInfo.route} arrives </span>
            <span className="stop-estimate">{arrivalInfo.estimated} </span>
          </p>
        </div>
      </Link>
    )
  }
}

StopListItem.propTypes = {
  state: React.PropTypes.shape({
    location: React.PropTypes.shape({
      lat: React.PropTypes.number.isRequired,
      lng: React.PropTypes.number.isRequired
    }).isRequired
  }).isRequired,
  stop: React.PropTypes.shape({
    arrivals: React.PropTypes.arrayOf(React.PropTypes.shape({
      estimated: React.PropTypes.number.isRequired,
      route: React.PropTypes.number.isRequired
    })).isRequired,
    desc: React.PropTypes.string.isRequired,
    lat: React.PropTypes.number.isRequired,
    lng: React.PropTypes.number.isRequired,
    locid: React.PropTypes.number.isRequired
  }).isRequired
}
