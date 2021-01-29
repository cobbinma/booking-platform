use protobuf::booking::api::booking_api_server::{BookingApi, BookingApiServer};
use async_trait::async_trait;
use tonic::{Response, Status, Request};
use protobuf::booking::models::{SlotInput, Booking, Slot};
use protobuf::booking::api::GetSlotResponse;
use tonic::transport::{Server, Identity, ServerTlsConfig};
use chrono::{DateTime, Duration};
use uuid::Uuid;
use std::ops::Add;

#[derive(Debug, Default)]
pub struct BookingService {}

#[async_trait]
impl BookingApi for BookingService {
    async fn get_slot(&self, req: Request<SlotInput>) -> Result<Response<GetSlotResponse>, Status> {
        let slot = req.into_inner();
        let starts_at = slot.starts_at;
        let duration = slot.duration;
        let ends_at= DateTime::parse_from_rfc3339(&starts_at)
            .map(|dt| dt.add(Duration::minutes(duration as i64)).to_rfc3339()).
            map_err(|e| Status::internal(e.to_string()))?;
        let venue_id = slot.venue_id;
        let email = slot.email;
        let people = slot.people;
        Ok(Response::new(GetSlotResponse{ r#match: Some(Slot{
            venue_id,
            email,
            people,
            starts_at,
            ends_at,
            duration,
        }), other_available_slots: vec![] }))
    }

    async fn create_booking(&self, req: Request<SlotInput>) -> Result<Response<Booking>, Status> {
        let slot = req.into_inner();
        let starts_at = slot.starts_at;
        let duration = slot.duration;
        let ends_at= DateTime::parse_from_rfc3339(&starts_at)
            .map(|dt| dt.add(Duration::minutes(duration as i64)).to_rfc3339()).
            map_err(|e| Status::internal(e.to_string()))?;
        let venue_id = slot.venue_id;
        let email = slot.email;
        let people = slot.people;
        Ok(Response::new(Booking{
            id: Uuid::new_v4().to_string(),
            venue_id,
            email,
            people,
            starts_at,
            ends_at,
            duration,
            table_id: Uuid::new_v4().to_string(),
        }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cert = tokio::fs::read("localhost.crt").await?;
    let key = tokio::fs::read("localhost.key").await?;

    let identity = Identity::from_pem(cert, key);

    let addr = "[::1]:6969".parse()?;
    let service = BookingService::default();

    Server::builder()
        .tls_config(ServerTlsConfig::new().identity(identity))?
        .add_service(BookingApiServer::new(service))
        .serve(addr)
        .await?;

    Ok(())
}
