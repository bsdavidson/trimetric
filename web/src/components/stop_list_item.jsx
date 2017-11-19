import React from "react";
import {connect} from "react-redux";
import {Link, withRouter} from "react-router-dom";

import {TrimetricPropTypes} from "./prop_types";
import {formatDistance} from "../helpers/directions";
// import {humanTimeUntilEpoch} from "../helpers/times";

export class StopListItem extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let location = {
      latitude: this.props.stop.lat,
      longitude: this.props.stop.lng
    };
    // let arrivalInfo = {
    //   route: this.props.stop.arrivals[0].route_id,
    //   estimated: humanTimeUntilEpoch(this.props.stop.arrivals[0].estimated)
    // };
    return (
      <Link className="stop-link" to={`/stop/${this.props.stop.id}`}>
        <div className="stop-list-item">
          <div className="stop-arrow">
            <span className="fui-arrow-right" />
          </div>
          <h4 className="stop-description">
            {this.props.stop.name} - #<span className="stop-id">
              {this.props.stop.id}
            </span>
          </h4>
          <p className="stop-metrics">
            <span className="stop-distance">
              {formatDistance(this.props.location, location)}
            </span>{" "}
            feet away.
            {/* <span className="stop-route"> {arrivalInfo.route_id} arrives </span> */}
            {/* <span className="stop-estimate">{arrivalInfo.estimated} </span> */}
          </p>
        </div>
      </Link>
    );
  }
}

StopListItem.propTypes = {
  location: TrimetricPropTypes.location,
  stop: TrimetricPropTypes.stop
};

function mapStateToProps(state) {
  return {
    location: state.location
  };
}

export default withRouter(connect(mapStateToProps)(StopListItem));
