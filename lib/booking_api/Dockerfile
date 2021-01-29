FROM rust:1.43 as builder

RUN USER=root cargo new --bin booking_api
WORKDIR ./booking_api
COPY ./Cargo.toml ./Cargo.toml
RUN cargo build --release
RUN rm src/*.rs

ADD . ./

RUN rm ./target/release/deps/booking_api*
RUN cargo build --release


FROM debian:buster-slim
ARG APP=/usr/src/app