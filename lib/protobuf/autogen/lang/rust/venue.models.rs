#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Venue {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub name: ::prost::alloc::string::String,
    #[prost(message, repeated, tag = "3")]
    pub opening_hours: ::prost::alloc::vec::Vec<OpeningHoursSpecification>,
    #[prost(message, repeated, tag = "4")]
    pub special_opening_hours: ::prost::alloc::vec::Vec<OpeningHoursSpecification>,
    #[prost(string, tag = "5")]
    pub slug: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct OpeningHoursSpecification {
    #[prost(uint32, tag = "1")]
    pub day_of_week: u32,
    #[prost(string, tag = "2")]
    pub opens: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub closes: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub valid_from: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub valid_through: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Table {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub name: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub capacity: u32,
}
