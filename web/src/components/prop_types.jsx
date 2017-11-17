import PropTypes from "prop-types";

const arrival = PropTypes.shape({
  arrival_time: PropTypes.string,
  departure_time: PropTypes.string,
  headsign: PropTypes.string.isRequired,
  route_color: PropTypes.string.isRequired,
  route_id: PropTypes.string.isRequired,
  route_long_name: PropTypes.string.isRequired,
  route_short_name: PropTypes.string.isRequired,
  route_text_color: PropTypes.string.isRequired,
  route_type: PropTypes.number.isRequired,
  stop_id: PropTypes.string.isRequired,
  trip_id: PropTypes.string.isRequired,
  vehicle_id: PropTypes.string,
  vehicle_label: PropTypes.string,
  vehicle_position: PropTypes.shape({
    lat: PropTypes.number.isRequired,
    lng: PropTypes.number.isRequired,
    bearing: PropTypes.number,
    odometer: PropTypes.number,
    speed: PropTypes.number
  }).isRequired,
  date: PropTypes.string.isRequired,
  estimated: PropTypes.number.isRequired
});

const arrivals = PropTypes.arrayOf(arrival);

const arrivalTime = PropTypes.string.isRequired;

const location = PropTypes.shape({
  locationType: PropTypes.string.isRequired,
  lat: PropTypes.number.isRequired,
  lng: PropTypes.number.isRequired
});

const locationClicked = PropTypes.object;

const stop = PropTypes.shape({
  arrivals: arrivals,
  name: PropTypes.string.isRequired,
  lat: PropTypes.number.isRequired,
  lng: PropTypes.number.isRequired,
  id: PropTypes.string.isRequired
});

const stops = PropTypes.arrayOf(stop);

const vehiclePosition = PropTypes.shape({
  position: PropTypes.shape({
    lat: PropTypes.number.isRequired,
    lng: PropTypes.number.isRequired
  }),
  vehicle: PropTypes.shape({
    id: PropTypes.string.isRequired,
    label: PropTypes.string
  }),
  route_type: PropTypes.number.isRequired
});

const vehiclePositions = PropTypes.arrayOf(vehiclePosition);

export const TrimetricPropTypes = {
  arrival,
  arrivals,
  arrivalTime,
  location,
  locationClicked,
  stop,
  stops,
  vehiclePosition,
  vehiclePositions
};
