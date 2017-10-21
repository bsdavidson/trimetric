create database trimetric;

\c trimetric

create table vehicles (
  id serial PRIMARY KEY,
  vehicle_id integer not null UNIQUE,
  type varchar(255) not null,
  sign_message varchar(255) not null
);