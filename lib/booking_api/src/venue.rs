use crate::service::VenueClient as Client;
use async_trait::async_trait;
use chrono::{DateTime, Datelike, NaiveDate, NaiveTime, TimeZone, Timelike, Utc};
use protobuf::venue::api::venue_api_client::VenueApiClient;
use protobuf::venue::api::GetVenueRequest;
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
    async fn get_opening_times(
        &self,
        venue_id: String,
        date: NaiveDate,
    ) -> Result<(DateTime<Utc>, DateTime<Utc>), Status> {
        let venue = &self
            .client
            .clone()
            .get_venue(GetVenueRequest { id: venue_id })
            .await?
            .into_inner();

        let opening_hours_specification = venue
            .opening_hours
            .iter()
            .filter(|&hours| hours.day_of_week == date.weekday().number_from_monday())
            .next()
            .ok_or_else(|| Status::invalid_argument("venue not open on given date"))?;

        let opens = NaiveTime::parse_from_str(&opening_hours_specification.opens, "%H:%M")
            .map_err(|e| {
                log::error!("could not parse opens time : {}", e);
                Status::internal("could not parse opens time")
            })
            .map(|o| {
                Utc.ymd(date.year(), date.month(), date.day()).and_hms(
                    o.hour(),
                    o.minute(),
                    o.second(),
                )
            })?;

        let closes = NaiveTime::parse_from_str(&opening_hours_specification.closes, "%H:%M")
            .map_err(|e| {
                log::error!("could not parse closes time : {}", e);
                Status::internal("could not parse closes time")
            })
            .map(|c| {
                Utc.ymd(date.year(), date.month(), date.day()).and_hms(
                    c.hour(),
                    c.minute(),
                    c.second(),
                )
            })?;

        Ok((opens, closes))
    }
}
