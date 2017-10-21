import React from "react";
import PropTypes from "prop-types";
import {hashHistory, withRouter, Link} from "react-router-dom";
import {connect} from "react-redux";

import ArrivalListItem from "./arrival_list_item";
import {updateLocation, LocationTypes} from "../actions";
import {buildQuery} from "../helpers/http";
import {formatEstimate} from "../helpers/times";
import {ColorMap} from "../helpers/colors";
import {store} from "../store";

let colorMap = new ColorMap();

export class ArrivalList extends React.Component {
  constructor(props) {
    super(props);
    this.handleDirectionsClick = this.handleDirectionsClick.bind(this);
    this.handleRouteNameClick = this.handleRouteNameClick.bind(this);
  }

  handleBack() {
    hashHistory.goBack();
  }

  handleDirectionsClick() {
    window.open(
      "http://maps.google.com/?" +
        buildQuery({
          f: "d",
          daddr: `${this.props.stop.lat},${this.props.stop.lng}`,
          saddr: `${this.props.location.lat},${this.props.location.lng}`
        })
    );
  }

  handleRouteNameClick() {
    store.dispatch(
      updateLocation(
        LocationTypes.STOP,
        this.props.stop.locid,
        this.props.stop.lat,
        this.props.stop.lng,
        false
      )
    );
  }

  render() {
    return (
      <div className="arrival-list">
        <Link className="back-button" to="/">
          <span className="fui-triangle-left-large" /> Back
        </Link>
        <h3
          className="arrival-list-description"
          onClick={this.handleRouteNameClick}>
          {this.props.stop.desc}
          <span
            className="arrival-list-directions"
            onClick={this.handleDirectionsClick}>
            Get Directions
          </span>
        </h3>
        <div className="arrival-list-items">
          {this.props.stop.arrivals.map((a, idx) => {
            return (
              <ArrivalListItem
                arrival={a}
                arrivalTime={formatEstimate(a.estimated)}
                color={colorMap.getColorForKey(a.route)}
                google={this.props.google}
                key={idx}
              />
            );
          })}
        </div>
      </div>
    );
  }
}

ArrivalList.propTypes = {
  google: PropTypes.object,
  location: PropTypes.shape({
    lat: PropTypes.number.isRequired,
    lng: PropTypes.number.isRequired
  }),
  stop: PropTypes.shape({
    arrivals: PropTypes.arrayOf(
      PropTypes.shape({
        estimated: PropTypes.number.isRequired,
        route: PropTypes.number.isRequired
      })
    ).isRequired,
    desc: PropTypes.string.isRequired,
    lat: PropTypes.number.isRequired,
    lng: PropTypes.number.isRequired,
    locid: PropTypes.number.isRequired
  }).isRequired
};

function mapStateToProps(state) {
  return {
    location: state.location
  };
}

export default withRouter(connect(mapStateToProps)(ArrivalList));
