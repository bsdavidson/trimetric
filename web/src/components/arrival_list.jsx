import React from "react";
import PropTypes from "prop-types";
import {withRouter, Link} from "react-router-dom";
import {connect} from "react-redux";
import moment from "moment";

import {TrimetricPropTypes} from "./prop_types";
import ArrivalListItem from "./arrival_list_item";
import Header from "./header";
import {clearLocation, updateLocation, LocationTypes} from "../actions";
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
    let {arrivals, stop} = this.props;
    if (arrivals === null) {
      return null;
    }
    let items = arrivals.map((a, idx) => {
      let arrivalTime = moment(a.date, "YYYY-MM-DD")
        .add(moment.duration(a.arrival_time))
        .valueOf();
      return (
        <ArrivalListItem
          arrival={a}
          arrivalTime={formatEstimate(arrivalTime)}
          color={colorMap.getColorForKey(a.route_id)}
          key={idx}
        />
      );
    });

    if (items.length === 0) {
      items = (
        <div className="arrival-list-none">
          {" "}
          There are no arrivals within the next hour for this stop. Better start
          walking.
        </div>
      );
    }

    return (
      <div className="arrival-list">
        <Header />
        <Link className="back-button" to="/">
          <span className="fui-triangle-left-large" /> Back
        </Link>
        <span
          className="arrival-list-directions"
          onClick={this.handleDirectionsClick}>
          Get Directions
        </span>

        <h3
          className="arrival-list-description"
          onClick={this.handleRouteNameClick}>
          {stop.name}
        </h3>

        <div className="arrival-list-items">{items}</div>
      </div>
    );
  }
}

ArrivalList.propTypes = {
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
    },
    onClearLocation: () => {
      dispatch(clearLocation());
    }
  };
}

function mapStateToProps(state) {
  return {
    location: state.location,
    arrivals: state.arrivals
  };
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(ArrivalList)
);
