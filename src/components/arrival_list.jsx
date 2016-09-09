import React from "react"
import { hashHistory } from "react-router"

import { updateLocation, LocationTypes } from "../actions"
import { ColorMap } from "../helpers/colors"
import { degreeToCompass } from "../helpers/directions"
import { buildQuery } from "../helpers/http"
import { formatEstimate } from "../helpers/times"
import { store } from "../store"


let colorMap = new ColorMap()

export class ArrivalList extends React.Component {
  constructor(props) {
    super(props)
    this.handleDirectionsClick = this.handleDirectionsClick.bind(this)
    this.handleRouteNameClick = this.handleRouteNameClick.bind(this)
  }

  handleBack() {
    hashHistory.goBack()
  }

  handleDirectionsClick() {
    window.open("http://maps.google.com/?" + buildQuery({
      f: "d",
      daddr: `${this.props.stop.lat},${this.props.stop.lng}`,
      saddr: `${this.props.state.location.lat},${this.props.state.location.lng}`
    }))
  }

  handleRouteNameClick() {
    store.dispatch(updateLocation(
      LocationTypes.STOP, this.props.stop.locid, this.props.stop.lat,
      this.props.stop.lng, false))
  }

  render() {
    return (
      <div className="arrival-list">
        <div className="back-button" onClick={this.handleBack}>
          <span className="fui-triangle-left-large"></span> Back
        </div>
        <h3 className="arrival-list-description" onClick={this.handleRouteNameClick}>
          {this.props.stop.desc}
          <span className="arrival-list-directions" onClick={this.handleDirectionsClick}>
            Get Directions
          </span>
        </h3>
        <div className="arrival-list-items">
          {this.props.stop.arrivals.map((a, idx) =>
            <ArrivalListItem
              arrival={a}
              arrivalTime={formatEstimate(a.estimated)}
              color={colorMap.getColorForKey(a.route)}
              google={this.props.google}
              key={idx}
              state={this.props.state} />
          )}
        </div>
      </div>
    )
  }
}

ArrivalList.propTypes = {
  google: React.PropTypes.object,
  state: React.PropTypes.shape({
    location: React.PropTypes.shape({
      lat: React.PropTypes.number.isRequired,
      lng: React.PropTypes.number.isRequired
    })
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

export class ArrivalListItem extends React.Component {
  constructor(props) {
    super(props)
    this.handleVehicleClick = this.handleVehicleClick.bind(this)
  }

  componentWillUnmount() {
    if (this.props.state.locationClicked &&
        this.props.arrival.vehicleID === this.props.state.locationClicked.id) {
      store.dispatch(updateLocation(
        LocationTypes.HOME , null, this.props.state.location.lat,
        this.props.state.location.lng, false))
    }
  }

  handleVehicleClick() {
    if (window) {
      window.scrollTo(0, 0);
    }
    store.dispatch(updateLocation(
      LocationTypes.VEHICLE, this.props.arrival.vehicleID, this.props.arrival.latitude,
      this.props.arrival.longitude, true))
  }

  render() {
    let routeClass = ""
    if (this.props.state.locationClicked &&
        this.props.state.locationClicked.following &&
        this.props.state.locationClicked.id === this.props.arrival.vehicleID) {
      routeClass = "active"
    }

    let routeStyle = {
      backgroundColor: colorMap.getColorForKey(this.props.arrival.route)
    }

    let vehicleDistance
    if (this.props.arrival.feet > 500) {
      vehicleDistance = Math.round(this.props.arrival.feet / 5280 * 100) / 100 + " miles away"
    } else {
      vehicleDistance = this.props.arrival.feet + " feet away"
    }
    return (
      <div className={"arrival-list-item " + routeClass} onClick={this.handleVehicleClick} >
        <div className="arrival-id" style={routeStyle}>{this.props.arrival.route}</div>
        <div className="arrival-name">{this.props.arrival.shortSign}</div>
        <div className="arrival-metrics">
          <div className="arrival-metric arrival-bus-distance">{vehicleDistance}</div>
          <div className="arrival-metric arrival-est-time">{this.props.arrivalTime}</div>
          <div className="arrival-metric arrival-direction">Traveling: {degreeToCompass(this.props.arrival.bearing)}</div>
        </div>
      </div>
    )
  }
}

ArrivalListItem.propTypes = {
  arrival: React.PropTypes.shape({
    bearing: React.PropTypes.number,
    feet: React.PropTypes.number.isRequired,
    latitude: React.PropTypes.number.isRequired,
    longitude: React.PropTypes.number.isRequired,
    route: React.PropTypes.number.isRequired,
    shortSign: React.PropTypes.string.isRequired,
    vehicleID: React.PropTypes.string
  }).isRequired,
  arrivalTime: React.PropTypes.string.isRequired,
  state: React.PropTypes.shape({
    location: React.PropTypes.object,
    locationClicked: React.PropTypes.object
  }).isRequired
}
