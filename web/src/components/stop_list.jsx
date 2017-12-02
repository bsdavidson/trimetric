import React from "react";
import {connect} from "react-redux";
import {withRouter} from "react-router-dom";

import StopListItem from "./stop_list_item";
import Header from "./header";
import {TrimetricPropTypes} from "./prop_types";
import {updateHomeLocation} from "../actions";
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

  render() {
    let stopItems;
    if (this.state.stops.length > 100) {
      stopItems = (
        <div className="stop-list-item">
          Too many stops in view. Please zoom in.
        </div>
      );
    } else {
      stopItems = this.state.stops.map((s, i) => (
        <StopListItem key={s.id + "-" + i} stop={s} />
      ));
    }

    return (
      <div className="stop-list">
        <Header />

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
        updateHomeLocation(
          position.coords.latitude,
          position.coords.longitude,
          true
        )
      );
    }
  };
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(StopList)
);
