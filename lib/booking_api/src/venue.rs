use crate::service::VenueClient as Client;
use async_trait::async_trait;
use protobuf::venue::api::venue_api_client::VenueApiClient;
use protobuf::venue::api::{GetVenueRequest, GetTablesRequest};
use protobuf::venue::models::Venue;
use tonic::transport::Channel;
use tonic::Status;

pub struct VenueClient {
    client: VenueApiClient<tonic::transport::Channel>,
}

impl VenueClient {
    pub fn new(client: VenueApiClient<Channel>) -> Self {
        VenueClient { client }
    }
}

#[async_trait]
impl Client for VenueClient {
    async fn get_venue(&self, venue_id: String) -> Result<Venue, Status> {
        self.client
            .clone()
            .get_venue(GetVenueRequest { id: venue_id, slug: "".to_string() })
            .await
            .map(|v| v.into_inner())
    }

    async fn get_tables_with_capacity(
        &self,
        venue_id: String,
        capacity: u32,
    ) -> Result<Vec<String>, Status> {
        tracing::debug!(
            "getting tables for venue {} with capacity {}",
            &venue_id,
            capacity
        );

        Ok(self
            .client
            .clone()
            .get_tables(GetTablesRequest { venue_id })
            .await?
            .into_inner()
            .tables
            .iter()
            .filter(|table| table.capacity >= capacity)
            .map(|table| table.id.clone())
            .collect())
    }
}
