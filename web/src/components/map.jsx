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
      (this.props.locationClicked === nextProps.locationClicked ||
        !nextProps.locationClicked) &&
      this.props.location === nextProps.location
    ) {
      return;
    }

    let newPos = null;
    if (this.props.location !== nextProps.location) {
      newPos = {
        lat: nextProps.location.lat,
        lng: nextProps.location.lng,
        locationType: nextProps.location.locationType
      };
    }

    if (newPos === null) {
      newPos = {
        lat: nextProps.locationClicked.lat,
        lng: nextProps.locationClicked.lng,
        locationType: nextProps.locationClicked.locationType
      };
    }

    this.handleViewportChange({
      latitude: newPos.lat,
      longitude: newPos.lng,
      zoom: newPos.locationType === LocationTypes.HOME ? 17 : 18,
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

    let layers = [];

    if (this.props.location) {
      let offset = 0.0001;
      let location = {
        type: "Feature",
        geometry: {
          type: "LineString",
          coordinates: [
            [this.props.location.lng - offset, this.props.location.lat],
            [this.props.location.lng, this.props.location.lat + offset],
            [this.props.location.lng + offset, this.props.location.lat],
            [this.props.location.lng, this.props.location.lat - offset],
            [this.props.location.lng - offset, this.props.location.lat]
          ]
        },

        properties: {
          lineColor: [241, 196, 15, 255],
          fillColor: [241, 196, 15, 255],
          radius: 1
        }
      };

      layers.push(
        new GeoJsonLayer({
          id: "location-home-point-layer",
          data: [location],
          opacity: 0.7,
          stroked: true,
          extruded: true,
          filled: true,
          lineWidthMinPixels: 4,
          visible: 1,
          fp64: true
        })
      );
    }

    if (this.props.locationClicked) {
      let location = {
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [
            this.props.locationClicked.lng,

            this.props.locationClicked.lat
          ]
        },
        properties: {
          lineColor: [0, 0, 0, 255],
          fillColor: [0, 255, 0, 255],
          radius: 1
        }
      };

      layers.push(
        new GeoJsonLayer({
          id: "location-clicked-point-layer",
          data: [location],
          opacity: 0.4,
          stroked: true,
          filled: true,
          lineWidthMinPixels: 2,
          pointRadiusScale: 30,
          visible: 1,
          fp64: true
        })
      );
    }

    if (this.props.stopsPointData) {
      layers.push(
        new GeoJsonLayer({
          id: "stops-point-layer",
          data: this.props.stopsPointData,
          opacity: 1 - clamp(tween),
          stroked: true,
          filled: true,
          pointRadiusScale:
            696.0864 - 106.8473 * zoom + 4.205566 * Math.pow(zoom, 2),
          visible: tween < 1,
          fp64: true
        })
      );
    }

    if (this.props.stopsIconData) {
      layers.push(
        new IconLayer({
          id: "stops-icon-layer",
          data: this.props.stopsIconData,
          iconAtlas: "/assets/sprites.png",
          iconMapping: IconMapping,
          visible: tween > 0,
          opacity: 1,
          fp64: true,
          sizeScale:
            -40.28287 + 0.1462691 * zoom + 0.3593278 * Math.pow(zoom, 2)
        })
      );
    }

    if (this.props.lineData) {
      layers = layers.concat(
        this.props.lineData.map((l, i) => {
          return new GeoJsonLayer({
            id: "geojson-line-layer" + i,
            data: l,
            getLineColor: () => l.color,
            lineWidthMinPixels: l.width,
            fp64: true
          });
        })
      );
    }

    if (this.props.vehiclesPointData) {
      layers.push(
        new GeoJsonLayer({
          id: "vehciles-point-layer",
          data: this.props.vehiclesPointData,
          opacity: 1 - clamp(tween),
          stroked: true,
          filled: true,
          pointRadiusScale:
            696.0864 - 106.8473 * zoom + 4.205566 * Math.pow(zoom, 2),
          visible: tween < 1,
          fp64: true
        })
      );
    }

    if (this.props.vehiclesIconData) {
      layers.push(
        new IconLayer({
          id: "vehicle-icon-layer",
          data: this.props.vehiclesIconData,
          iconAtlas: "/assets/sprites.png",
          iconMapping: IconMapping,
          visible: tween > 0,
          opacity: 1,
          fp64: true,
          sizeScale:
            -40.28287 + 0.1462691 * zoom + 0.3593278 * Math.pow(zoom, 2)
        })
      );
    }

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
            layers={layers}
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
    stopsPointData: state.stopsPointData,
    stopsIconData: state.stopsIconData,
    vehiclesPointData: state.vehiclesPointData,
    vehiclesIconData: state.vehiclesIconData,
    lineData: state.lineData
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(MapBox);
