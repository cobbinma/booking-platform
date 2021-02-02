use crate::service::TableClient as Client;
use async_trait::async_trait;
use protobuf::venue::api::table_api_client::TableApiClient;
use protobuf::venue::api::GetTablesRequest;
use tonic::transport::Channel;
use tonic::Status;

pub struct TableClient {
    client: TableApiClient<tonic::transport::Channel>,
}

impl TableClient {
    pub fn new(client: TableApiClient<Channel>) -> Self {
        TableClient { client }
    }
}

#[async_trait]
impl Client for TableClient {
    async fn get_tables_with_capacity(
        &self,
        venue_id: String,
        capacity: u32,
    ) -> Result<Vec<String>, Status> {
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
