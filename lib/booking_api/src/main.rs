#[macro_use]
extern crate diesel;

use alcoholic_jwt::{token_kid, validate, Validation, JWKS};
use protobuf::booking::api::booking_api_server::BookingApiServer;
use serde::{Deserialize, Serialize};
use tonic::transport::{Certificate, Channel, ClientTlsConfig, Identity, Server, ServerTlsConfig};
use tonic::{metadata::MetadataValue, Request, Status};

mod postgres;
mod service;

pub mod models;
pub mod schema;
mod table;
mod venue;

use crate::postgres::Postgres;
use protobuf::venue::api::table_api_client::TableApiClient;
use protobuf::venue::api::venue_api_client::VenueApiClient;
use service::BookingService;
use std::collections::HashMap;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    dotenv::dotenv().ok();

    let cert = tokio::fs::read("localhost.crt").await?;
    let key = tokio::fs::read("localhost.key").await?;
    let server_root_ca_cert = Certificate::from_pem(&cert);

    let identity = Identity::from_pem(cert, key);

    let tls = ClientTlsConfig::new()
        .domain_name("localhost")
        .ca_certificate(server_root_ca_cert)
        .identity(identity.clone());

    let addr = "[::1]:6969".parse()?;

    let mut map = HashMap::new();
    map.insert(
        "client_id",
        std::env::var("AUTH0_CLIENT_ID").expect("client id not set"),
    );
    map.insert(
        "client_secret",
        std::env::var("AUTH0_CLIENT_SECRET").expect("client secret not set"),
    );
    map.insert(
        "audience",
        std::env::var("AUTH0_VENUE_API_IDENTIFIER").expect("venue api identifier not set"),
    );
    map.insert("grant_type", "client_credentials".to_string());
    let client = reqwest::Client::new();
    let mut resp = client
        .post(&format!(
            "{}oauth/token",
            std::env::var("AUTHORITY").expect("AUTHORITY must be set")
        ))
        .json(&map)
        .send()?;
    let token = resp.json::<Token>().map(|t| t.access_token)?;

    let channel = Channel::from_static("http://[::1]:8888");

    let interceptor = Box::new(move |mut req: Request<()>| {
        req.metadata_mut().insert(
            "authorization",
            MetadataValue::from_str(&*format!("Bearer {}", token)).unwrap(),
        );
        Ok(req)
    });

    let venue_client = VenueApiClient::with_interceptor(
        channel.clone().tls_config(tls.clone())?.connect().await?,
        interceptor.clone(),
    );

    let table_client =
        TableApiClient::with_interceptor(channel.tls_config(tls)?.connect().await?, interceptor);

    let service = BookingService::new(
        Box::new(Postgres::new()?),
        Box::new(venue::VenueClient::new(venue_client)),
        Box::new(table::TableClient::new(table_client)),
        None,
    )?;

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();

    tracing::info!(message = "Starting server.", %addr);

    Server::builder()
        .tls_config(ServerTlsConfig::new().identity(identity))?
        .trace_fn(|_| tracing::info_span!("booking_api"))
        .add_service(BookingApiServer::with_interceptor(service, check_auth))
        .serve(addr)
        .await?;

    Ok(())
}

fn check_auth(req: Request<()>) -> Result<Request<()>, Status> {
    let md = req
        .metadata()
        .get("authorization")
        .ok_or_else(|| Status::unauthenticated("no valid auth token"))?;

    let token = md
        .to_str()
        .map_err(|_| Status::invalid_argument("could not parse token"))?;

    if validate_token(token)? {
        Ok(req)
    } else {
        Err(Status::unauthenticated("could not valididate auth token"))
    }
}

pub fn validate_token(token: &str) -> Result<bool, Status> {
    let token = token.trim_start_matches("Bearer ");
    let authority = std::env::var("AUTHORITY").expect("AUTHORITY must be set");
    let jwks = fetch_jwks(&format!(
        "{}{}",
        authority.as_str(),
        ".well-known/jwks.json"
    ))?;
    let validations = vec![Validation::Issuer(authority), Validation::SubjectPresent];
    let kid = token_kid(&token)
        .map_err(|_| Status::internal("failed to fetch jwts"))?
        .ok_or_else(|| Status::invalid_argument("failed to decode kid"))?;
    let jwk = jwks
        .find(&kid)
        .ok_or_else(|| Status::invalid_argument("specified key not found in set"))?;
    let res = validate(token, jwk, validations);

    Ok(res.is_ok())
}

fn fetch_jwks(uri: &str) -> Result<JWKS, Status> {
    let mut res = reqwest::get(uri).map_err(|_| Status::internal("could not get jwks"))?;
    let val = res
        .json::<JWKS>()
        .map_err(|_| Status::internal("could not unmarshall jwks"))?;

    Ok(val)
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Token {
    #[serde(rename = "access_token")]
    pub access_token: String,
    #[serde(rename = "token_type")]
    pub token_type: String,
}
