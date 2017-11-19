import React from "react";
import Header from "./header";

export default function Info() {
  return (
    <div className="info">
      <Header />
      <div className="info-text">
        <p>
          Trimetric is a realtime visualization of the{" "}
          <a href="https://www.trimet.org">Trimet</a> transit system in
          Portland, OR. The view that you are currently looking at is showing
          the location of every vehicle and transit stop. The orange dots are
          vehicles, and the black dots are stops.{" "}
        </p>
        <p>
          If you zoom in, you can see more information about the stops,
          including upcoming arrivals.{" "}
        </p>
        <p>
          The data for the view comes from static and realtime{" "}
          <a href="https://developers.google.com/transit/">GTFS</a> feeds
          provided by Trimet.
        </p>
        <p className="info-text-credits">
          Trimetric was built by me,{" "}
          <a href="https://briand.co">Brian Davidson</a>. It&apos;s open source
          and you can find the code at
          <a href="https://github.com/bsdavidson/trimetric">
            {" "}
            github.com/bsdavidson/trimetric
          </a>
        </p>
      </div>
      <div className="info-icons">
        <div className="info-icons-text">
          Trimetric is powered by these technologies:
        </div>
        <a href="https://reactjs.org/" title="React">
          <img
            alt="React"
            className="info-icon react"
            src="/assets/react.svg"
          />
        </a>
        <a href="https://redux.js.org/" title="Redux">
          <img
            alt="Redux"
            className="info-icon redux"
            src="/assets/redux.svg"
          />
        </a>
        <a href="https://kafka.apache.org/" title="Kafka">
          <img
            alt="Kafka"
            className="info-icon kafka"
            src="/assets/kafka.svg"
          />
        </a>
        <a href="https://golang.org/" title="Go">
          <img alt="Go" className="info-icon go" src="/assets/go.svg" />
        </a>
        <a href="https://www.docker.com/" title="Docker">
          <img
            alt="Docker"
            className="info-icon docker"
            src="/assets/docker.svg"
          />
        </a>
        <a href="https://www.postgresql.org/" title="Postgres">
          <img
            alt="Postgres"
            className="info-icon postgresql"
            src="/assets/postgresql.svg"
          />
        </a>
      </div>
    </div>
  );
}
