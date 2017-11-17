import React, {Component} from "react";
import {connect} from "react-redux";

import {TrimetricPropTypes} from "./prop_types";
import PropTypes from "prop-types";
import InteractiveMap, {experimental} from "react-map-gl";
import "mapbox-gl/dist/mapbox-gl.css";

class MapBox extends Component {
  constructor(props) {
    super(props);

    this.state = {
      viewport: {
        latitude: this.props.location.lat,
        longitude: this.props.location.lng,
        zoom: 16
      },
      settings: {
        dragPan: true
      }
    };
    this.handleMapRef = this.handleMapRef.bind(this);
    this.handleViewportChange = this.handleViewportChange.bind(this);
  }

  handleMapRef(map) {
    this.mapRef = map;
  }

  componentWillReceiveProps(nextProps) {
    if (
      this.props.locationClicked === nextProps.locationClicked ||
      !nextProps.locationClicked
    ) {
      return;
    }
    this.handleViewportChange({
      latitude: nextProps.locationClicked.lat,
      longitude: nextProps.locationClicked.lng,
      zoom: 16,
      transitionInterpolator: experimental.viewportFlyToInterpolator,
      transitionDuration: 1000
    });
  }

  handleViewportChange(viewport) {
    this.setState({viewport: Object.assign({}, this.state.viewport, viewport)});
    if (this.props.onBoundsChanged) {
      let bounds = this.mapRef.getMap().getBounds();
      this.props.onBoundsChanged({
        sw: {lat: bounds.getSouth(), lng: bounds.getWest()},
        ne: {lat: bounds.getNorth(), lng: bounds.getEast()}
      });
    }
  }

  renderChildren() {
    return React.Children.map(this.props.children, child => {
      if (!React.isValidElement(child)) {
        return child;
      }
      return React.cloneElement(child, {
        map: {}
      });
    });
  }

  render() {
    return (
      <div id="mapbox" className="app-map">
        <InteractiveMap
          ref={this.handleMapRef}
          onViewportChange={this.handleViewportChange}
          mapboxApiAccessToken={process.env.MAPBOX_ACCESS_TOKEN}
          width={this.props.width}
          height={this.props.height}
          latitude={this.state.viewport.latitude}
          longitude={this.state.viewport.longitude}
          transitionInterpolator={this.state.viewport.transitionInterpolator}
          transitionDuration={this.state.viewport.transitionDuration}
          zoom={this.state.viewport.zoom}
          dragPan={this.state.settings.dragPan}>
          {this.renderChildren()}
        </InteractiveMap>
      </div>
    );
  }
}

MapBox.propTypes = {
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  onBoundsChanged: PropTypes.func,
  location: TrimetricPropTypes.location,
  locationClicked: TrimetricPropTypes.locationClicked
};

function mapStateToProps(state) {
  return {
    location: state.location,
    locationClicked: state.locationClicked
  };
}

export default connect(mapStateToProps)(MapBox);
