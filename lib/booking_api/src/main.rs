use protobuf::booking::api::booking_api_server::{BookingApi, BookingApiServer};
use async_trait::async_trait;
use tonic::{Response, Status, Request};
use protobuf::booking::models::{SlotInput, Booking};
use protobuf::booking::api::GetSlotResponse;
use tonic::transport::{Server, Identity, ServerTlsConfig};

#[derive(Debug, Default)]
pub struct BookingService {}

#[async_trait]
impl BookingApi for BookingService {
    async fn get_slot(&self, _: Request<SlotInput>) -> Result<Response<GetSlotResponse>, Status> {
        unimplemented!()
    }

    async fn create_booking(&self, _: Request<SlotInput>) -> Result<Response<Booking>, Status> {
        unimplemented!()
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
