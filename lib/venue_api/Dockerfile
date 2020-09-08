FROM rust:1.46 as builder
WORKDIR /usr/src/venue_api
COPY . .
RUN cargo install --path .

FROM debian:buster-slim
COPY --from=builder /usr/local/cargo/bin/venue_api /usr/local/bin/venue_api
CMD ["venue_api"]