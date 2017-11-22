import React, {Component} from "react";
import {connect} from "react-redux";
import PropTypes from "prop-types";
import InteractiveMap, {experimental} from "react-map-gl";
import DeckGL, {GeoJsonLayer, IconLayer} from "deck.gl";
import "mapbox-gl/dist/mapbox-gl.css";

import {DEFAULT_ZOOM} from "../store";
import {LocationTypes, clearLocation} from "../actions";
import {TrimetricPropTypes} from "./prop_types";

const IconMapping = {
  tram: {
    x: 0,
    y: 0,
    width: 80,
    height: 80
  },
  bus: {
    x: 80,
    y: 0,
    width: 80,
    height: 80
  },
  home: {
    x: 160,
    y: 0,
    width: 80,
    height: 80,
    mask: true
  },
  stop: {
    x: 240,
    y: 0,
    width: 80,
    height: 80,
    mask: true
  }
};

function clamp(f) {
  return f < 0 ? 0 : f > 1 ? 1 : f;
}

export class CustomMapControls extends experimental.MapControls {
  constructor(props) {
    super(props);
    this.props = props;
    this.events = ["click", "mousedown"];
  }

  handleEvent(event) {
    if (event.type === "mousedown") {
      this.props.onMouseDown();
    }
    return super.handleEvent(event);
  }
}

class MapBox extends Component {
  constructor(props) {
    super(props);

    this.state = {
      viewport: {
        latitude: this.props.location.lat,
        longitude: this.props.location.lng,
        zoom: DEFAULT_ZOOM,
        pitch: 45,
        bearing: 0
      },
      settings: {
        dragPan: true
      },
      test: false
    };
    this.handleMapRef = this.handleMapRef.bind(this);
    this.handleMapMouseDown = this.handleMapMouseDown.bind(this);
    this.handleViewportChange = this.handleViewportChange.bind(this);
    this.mapControls = new CustomMapControls({
      onMouseDown: this.handleMapMouseDown
    });
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
      zoom:
        nextProps.locationClicked.locationType === LocationTypes.HOME ? 17 : 19,
      transitionInterpolator: experimental.viewportFlyToInterpolator,
      transitionDuration: 1200
    });

    if (this.state.viewport.zoom < 15) {
      this.props.onClearLocation();
    }
  }

  handleViewportChange(viewport) {
    this.setState({viewport: Object.assign({}, this.state.viewport, viewport)});
    if (this.props.onViewportChange) {
      let bounds = this.mapRef.getMap().getBounds();
      let zoom = this.state.viewport.zoom;
      this.props.onViewportChange(
        {
          sw: {lat: bounds.getSouth(), lng: bounds.getWest()},
          ne: {lat: bounds.getNorth(), lng: bounds.getEast()}
        },
        zoom
      );
    }
  }

  handleMapMouseDown() {
    this.props.onClearLocation();
  }

  render() {
    let zoom = this.state.viewport.zoom;
    let tween = zoom - 12;

    const geoJsonLayer = new GeoJsonLayer({
      id: "geojson-layer",
      data: this.props.geoJsonData,
      opacity: 1 - clamp(tween),
      stroked: true,
      filled: true,
      pointRadiusScale:
        696.0864 - 106.8473 * zoom + 4.205566 * Math.pow(zoom, 2),
      visible: tween < 1,
      fp64: true
    });

    const iconLayer = new IconLayer({
      id: "icon-layer",
      data: this.props.iconData,
      iconAtlas: "/assets/sprites.png",
      iconMapping: IconMapping,
      visible: tween > 0,
      opacity: 1,
      fp64: true,
      sizeScale: -40.28287 + 0.1462691 * zoom + 0.3593278 * Math.pow(zoom, 2)
    });

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
          dragPan={this.state.settings.dragPan}
          mapControls={this.mapControls}>
          <DeckGL
            width={this.props.width}
            height={this.props.height}
            latitude={this.state.viewport.latitude}
            longitude={this.state.viewport.longitude}
            zoom={this.state.viewport.zoom}
            layers={[geoJsonLayer, iconLayer]}
          />
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

function mapDispatchToProps(dispatch) {
  return {
    onClearLocation: () => {
      dispatch(clearLocation());
    }
  };
}

function mapStateToProps(state) {
  return {
    location: state.location,
    locationClicked: state.locationClicked,
    geoJsonData: state.geoJsonData,
    iconData: state.iconData
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(MapBox);
