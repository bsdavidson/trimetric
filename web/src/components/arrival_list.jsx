import React from "react";
import PropTypes from "prop-types";
import {hashHistory, withRouter, Link} from "react-router-dom";
import {connect} from "react-redux";

import {TrimetricPropTypes} from "./prop_types";
import ArrivalListItem from "./arrival_list_item";
import {updateLocation, LocationTypes} from "../actions";
import {buildQuery} from "../helpers/http";
import {formatEstimate} from "../helpers/times";
import {ColorMap} from "../helpers/colors";

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
    this.props.onRouteNameClick(this.props.stop);
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
          {this.props.stop.name}
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
                color={colorMap.getColorForKey(a.route_id)}
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
  google: TrimetricPropTypes.google,
  location: TrimetricPropTypes.location,
  onRouteNameClick: PropTypes.func.isRequired,
  stop: TrimetricPropTypes.stop
};

function mapDispatchToProps(dispatch) {
  return {
    onRouteNameClick: stop => {
      dispatch(
        updateLocation(LocationTypes.STOP, stop.id, stop.lat, stop.lng, false)
      );
    }
  };
}

function mapStateToProps(state) {
  return {
    location: state.location
  };
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(ArrivalList)
);
