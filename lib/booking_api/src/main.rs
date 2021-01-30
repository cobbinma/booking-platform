#[macro_use]
extern crate diesel;

use alcoholic_jwt::{token_kid, validate, Validation, JWKS};
use protobuf::booking::api::booking_api_server::BookingApiServer;
use tonic::transport::{Identity, Server, ServerTlsConfig};
use tonic::{Request, Status};

mod postgres;
mod service;

pub mod models;
pub mod schema;

use service::BookingService;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    femme::with_level(femme::LevelFilter::Info);

    dotenv::dotenv().ok();

    let cert = tokio::fs::read("localhost.crt").await?;
    let key = tokio::fs::read("localhost.key").await?;

    let identity = Identity::from_pem(cert, key);

    let addr = "[::1]:6969".parse()?;
    let service = BookingService::default();

    log::info!("listening on port {}", &addr);

    Server::builder()
        .tls_config(ServerTlsConfig::new().identity(identity))?
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
