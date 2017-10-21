import React from "react";
import PropTypes from "prop-types";

import {clearLocation, LocationTypes} from "../actions";
import {store} from "../store";

export class PanTo extends React.Component {
  constructor(props) {
    super(props);
    this.dragListener = null;
    this.handleDragStart = this.handleDragStart.bind(this);
  }

  addMapListeners(map) {
    if (!map) {
      return;
    }
    this.removeMapListeners();
    this.dragListener = map.addListener("dragstart", this.handleDragStart);
  }

  componentDidMount() {
    this.addMapListeners(this.props.map);
  }

  componentDidUpdate() {
    if (!this.props.location || !this.props.map) {
      return;
    }
    this.props.map.panTo({
      lat: this.props.location.lat,
      lng: this.props.location.lng
    });
    if (this.props.location.locationType !== LocationTypes.VEHICLE) {
      store.dispatch(clearLocation());
    }
  }

  componentWillUnmount() {
    this.removeMapListeners();
  }

  componentWillUpdate(nextProps) {
    if (this.props.map !== nextProps.map) {
      this.removeMapListeners();
      this.addMapListeners(nextProps.map);
    }
  }

  handleDragStart() {
    store.dispatch(clearLocation());
  }

  removeMapListeners() {
    if (!this.dragListener || !this.props.google) {
      return;
    }
    this.props.google.maps.event.removeListener(this.dragListener);
  }

  render() {
    return null;
  }
}

PanTo.propTypes = {
  google: PropTypes.object,
  location: PropTypes.shape({
    locationType: PropTypes.string.isRequired,
    lat: PropTypes.number.isRequired,
    lng: PropTypes.number.isRequired
  }),
  map: PropTypes.object
};
