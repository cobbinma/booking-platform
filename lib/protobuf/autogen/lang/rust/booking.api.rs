#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetSlotResponse {
    #[prost(message, optional, tag = "1")]
    pub r#match: ::core::option::Option<super::models::Slot>,
    #[prost(message, repeated, tag = "2")]
    pub other_available_slots: ::prost::alloc::vec::Vec<super::models::Slot>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBookingsRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub date: ::prost::alloc::string::String,
    #[prost(int32, tag = "3")]
    pub page: i32,
    #[prost(int32, tag = "4")]
    pub limit: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetBookingsResponse {
    #[prost(message, repeated, tag = "1")]
    pub bookings: ::prost::alloc::vec::Vec<super::models::Booking>,
    #[prost(bool, tag = "2")]
    pub has_next_page: bool,
    #[prost(int32, tag = "3")]
    pub pages: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CancelBookingRequest {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BookingInput {
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
    #[prost(string, tag = "6")]
    pub family_name: ::prost::alloc::string::String,
    #[prost(string, tag = "7")]
    pub given_name: ::prost::alloc::string::String,
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
#[doc = r" Generated client implementations."]
pub mod booking_api_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    pub struct BookingApiClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl BookingApiClient<tonic::transport::Channel> {
        #[doc = r" Attempt to create a new client by connecting to a given endpoint."]
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> BookingApiClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::ResponseBody: Body + HttpBody + Send + 'static,
        T::Error: Into<StdError>,
        <T::ResponseBody as HttpBody>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor(inner: T, interceptor: impl Into<tonic::Interceptor>) -> Self {
            let inner = tonic::client::Grpc::with_interceptor(inner, interceptor);
            Self { inner }
        }
        pub async fn get_slot(
            &mut self,
            request: impl tonic::IntoRequest<super::SlotInput>,
        ) -> Result<tonic::Response<super::GetSlotResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/booking.api.BookingAPI/GetSlot");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn create_booking(
            &mut self,
            request: impl tonic::IntoRequest<super::BookingInput>,
        ) -> Result<tonic::Response<super::super::models::Booking>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/booking.api.BookingAPI/CreateBooking");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn get_bookings(
            &mut self,
            request: impl tonic::IntoRequest<super::GetBookingsRequest>,
        ) -> Result<tonic::Response<super::GetBookingsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/booking.api.BookingAPI/GetBookings");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn cancel_booking(
            &mut self,
            request: impl tonic::IntoRequest<super::CancelBookingRequest>,
        ) -> Result<tonic::Response<super::super::models::Booking>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/booking.api.BookingAPI/CancelBooking");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
    impl<T: Clone> Clone for BookingApiClient<T> {
        fn clone(&self) -> Self {
            Self {
                inner: self.inner.clone(),
            }
        }
    }
    impl<T> std::fmt::Debug for BookingApiClient<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "BookingApiClient {{ ... }}")
        }
    }
}
#[doc = r" Generated server implementations."]
pub mod booking_api_server {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = "Generated trait containing gRPC methods that should be implemented for use with BookingApiServer."]
    #[async_trait]
    pub trait BookingApi: Send + Sync + 'static {
        async fn get_slot(
            &self,
            request: tonic::Request<super::SlotInput>,
        ) -> Result<tonic::Response<super::GetSlotResponse>, tonic::Status>;
        async fn create_booking(
            &self,
            request: tonic::Request<super::BookingInput>,
        ) -> Result<tonic::Response<super::super::models::Booking>, tonic::Status>;
        async fn get_bookings(
            &self,
            request: tonic::Request<super::GetBookingsRequest>,
        ) -> Result<tonic::Response<super::GetBookingsResponse>, tonic::Status>;
        async fn cancel_booking(
            &self,
            request: tonic::Request<super::CancelBookingRequest>,
        ) -> Result<tonic::Response<super::super::models::Booking>, tonic::Status>;
    }
    #[derive(Debug)]
    pub struct BookingApiServer<T: BookingApi> {
        inner: _Inner<T>,
    }
    struct _Inner<T>(Arc<T>, Option<tonic::Interceptor>);
    impl<T: BookingApi> BookingApiServer<T> {
        pub fn new(inner: T) -> Self {
            let inner = Arc::new(inner);
            let inner = _Inner(inner, None);
            Self { inner }
        }
        pub fn with_interceptor(inner: T, interceptor: impl Into<tonic::Interceptor>) -> Self {
            let inner = Arc::new(inner);
            let inner = _Inner(inner, Some(interceptor.into()));
            Self { inner }
        }
    }
    impl<T, B> Service<http::Request<B>> for BookingApiServer<T>
    where
        T: BookingApi,
        B: HttpBody + Send + Sync + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = Never;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(&mut self, _cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            let inner = self.inner.clone();
            match req.uri().path() {
                "/booking.api.BookingAPI/GetSlot" => {
                    #[allow(non_camel_case_types)]
                    struct GetSlotSvc<T: BookingApi>(pub Arc<T>);
                    impl<T: BookingApi> tonic::server::UnaryService<super::SlotInput> for GetSlotSvc<T> {
                        type Response = super::GetSlotResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SlotInput>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).get_slot(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = GetSlotSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = if let Some(interceptor) = interceptor {
                            tonic::server::Grpc::with_interceptor(codec, interceptor)
                        } else {
                            tonic::server::Grpc::new(codec)
                        };
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/booking.api.BookingAPI/CreateBooking" => {
                    #[allow(non_camel_case_types)]
                    struct CreateBookingSvc<T: BookingApi>(pub Arc<T>);
                    impl<T: BookingApi> tonic::server::UnaryService<super::BookingInput> for CreateBookingSvc<T> {
                        type Response = super::super::models::Booking;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::BookingInput>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).create_booking(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = CreateBookingSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = if let Some(interceptor) = interceptor {
                            tonic::server::Grpc::with_interceptor(codec, interceptor)
                        } else {
                            tonic::server::Grpc::new(codec)
                        };
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/booking.api.BookingAPI/GetBookings" => {
                    #[allow(non_camel_case_types)]
                    struct GetBookingsSvc<T: BookingApi>(pub Arc<T>);
                    impl<T: BookingApi> tonic::server::UnaryService<super::GetBookingsRequest> for GetBookingsSvc<T> {
                        type Response = super::GetBookingsResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetBookingsRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).get_bookings(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = GetBookingsSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = if let Some(interceptor) = interceptor {
                            tonic::server::Grpc::with_interceptor(codec, interceptor)
                        } else {
                            tonic::server::Grpc::new(codec)
                        };
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/booking.api.BookingAPI/CancelBooking" => {
                    #[allow(non_camel_case_types)]
                    struct CancelBookingSvc<T: BookingApi>(pub Arc<T>);
                    impl<T: BookingApi> tonic::server::UnaryService<super::CancelBookingRequest>
                        for CancelBookingSvc<T>
                    {
                        type Response = super::super::models::Booking;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::CancelBookingRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).cancel_booking(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = CancelBookingSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = if let Some(interceptor) = interceptor {
                            tonic::server::Grpc::with_interceptor(codec, interceptor)
                        } else {
                            tonic::server::Grpc::new(codec)
                        };
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => Box::pin(async move {
                    Ok(http::Response::builder()
                        .status(200)
                        .header("grpc-status", "12")
                        .header("content-type", "application/grpc")
                        .body(tonic::body::BoxBody::empty())
                        .unwrap())
                }),
            }
        }
    }
    impl<T: BookingApi> Clone for BookingApiServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self { inner }
        }
    }
    impl<T: BookingApi> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(self.0.clone(), self.1.clone())
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: BookingApi> tonic::transport::NamedService for BookingApiServer<T> {
        const NAME: &'static str = "booking.api.BookingAPI";
    }
}
