#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Booking {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub email: ::prost::alloc::string::String,
    #[prost(uint32, tag = "4")]
    pub people: u32,
    #[prost(string, tag = "5")]
    pub starts_at: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub ends_at: ::prost::alloc::string::String,
    #[prost(uint32, tag = "7")]
    pub duration: u32,
    #[prost(string, tag = "8")]
    pub table_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SlotInput {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub email: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub people: u32,
    #[prost(string, tag = "4")]
    pub starts_at: ::prost::alloc::string::String,
    #[prost(uint32, tag = "5")]
    pub duration: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Slot {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub email: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub people: u32,
    #[prost(string, tag = "4")]
    pub starts_at: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub ends_at: ::prost::alloc::string::String,
    #[prost(uint32, tag = "6")]
    pub duration: u32,
}
