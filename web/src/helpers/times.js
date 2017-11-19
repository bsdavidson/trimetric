import moment from "moment";

export function formatEstimate(estTimeEpoch) {
  let seconds = secondsUntilEpoch(estTimeEpoch);
  if (seconds > 60) {
    return "ETA " + humanTimeUntilEpoch(estTimeEpoch);
  } else if (seconds > 10) {
    return "Now";
  } else {
    return "Arrived";
  }
}

export function humanTimeUntilEpoch(epoch) {
  let now = new Date();
  let estDate = moment(new Date(epoch));
  let estTime = estDate.diff(now);
  return moment.duration(estTime).humanize(true);
}

export function secondsUntilEpoch(epoch) {
  let now = new Date().getTime();
  return Math.floor((epoch - now) / 1000);
}
