FROM rust:1.49 as builder

RUN USER=root cargo new --bin booking_api
RUN rustup component add rustfmt
WORKDIR ./booking_api
COPY .protobuf ./.protobuf
COPY ./Cargo.toml ./Cargo.toml
RUN cargo build --release
RUN rm src/*.rs

ADD . ./

RUN rm ./target/release/deps/booking_api*
RUN cargo build --release


FROM debian:buster-slim
ARG APP=/usr/src/app