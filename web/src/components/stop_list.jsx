import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import StopListItem from "./stop_list_item";
import Header from "./header";
import {TrimetricPropTypes} from "./prop_types";
import {updateLocation, LocationTypes} from "../actions";
import {withinBoundingBox} from "../helpers/geom";

export class StopList extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      boundingBox: this.props.boundingBox,
      stops: this.props.stops.filter(s =>
        withinBoundingBox(s, this.props.boundingBox)
      )
    };

    this.timeout = null;
    this.stopsInView = 0;
    this.handleCurrentLocationClick = this.handleCurrentLocationClick.bind(
      this
    );
  }

  componentWillReceiveProps(nextProps) {
    if (
      nextProps.boundingBox === this.state.boundingBox &&
      nextProps.stops === this.props.stops
    ) {
      return;
    }

    clearTimeout(this.timeout);
    this.timeout = setTimeout(() => {
      this.setState({
        boundingBox: nextProps.boundingBox,
        stops: nextProps.stops.filter(s =>
          withinBoundingBox(s, this.state.boundingBox)
        )
      });
    }, 100);
  }

  componentWillUnmount() {
    clearTimeout(this.timeout);
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (
      this.state.boundingBox === nextState.boundingBox &&
      this.state.stops === nextState.stops
    ) {
      return false;
    }
    return true;
  }

  handleCurrentLocationClick() {
    if ("geolocation" in navigator) {
      navigator.geolocation.watchPosition(position => {
        this.props.onLocationClick(position);
      });
    } else {
      alert("Your browser doesn't support geolocation");
    }
  }

  render() {
    let stopItems;
    if (this.state.stops.length > 100) {
      stopItems = (
        <div className="stop-list-item">
          Too many stops in view. Please zoom in.
        </div>
      );
    } else {
      stopItems = this.state.stops.map(s => (
        <StopListItem key={s.id} stop={s} />
      ));
    }

    return (
      <div className="stop-list">
        <Header />
        <div
          className="use-current"
          title="Use Current GPS position"
          onClick={this.handleCurrentLocationClick}>
          <span className="fui-location" />
        </div>
        <h3>Visible Stops ({this.state.stops.length})</h3>

        <div className="stop-list-items">{stopItems}</div>
      </div>
    );
  }
}

StopList.propTypes = {
  stops: TrimetricPropTypes.stops
};

function mapStateToProps(state) {
  return {
    boundingBox: state.boundingBox,
    stops: state.stops
  };
}

function mapDispatchToProps(dispatch) {
  return {
    onLocationClick: position => {
      dispatch(
        updateLocation(
          LocationTypes.HOME,
          "GPS",
          position.coords.latitude,
          position.coords.longitude,
          false
        )
      );
    }
  };
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(StopList)
);
