#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetVenueRequest {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub slug: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CreateVenueRequest {
    #[prost(string, tag = "1")]
    pub name: ::prost::alloc::string::String,
    #[prost(message, repeated, tag = "2")]
    pub opening_hours: ::prost::alloc::vec::Vec<super::models::OpeningHoursSpecification>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTablesRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetTablesResponse {
    #[prost(message, repeated, tag = "1")]
    pub tables: ::prost::alloc::vec::Vec<super::models::Table>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddTableRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub name: ::prost::alloc::string::String,
    #[prost(uint32, tag = "3")]
    pub capacity: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemoveTableRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub table_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct IsAdminRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub email: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct IsAdminResponse {
    #[prost(bool, tag = "1")]
    pub is_admin: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddAdminRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub email: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddAdminResponse {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub email: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemoveAdminRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub email: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemoveAdminResponse {
    #[prost(string, tag = "1")]
    pub email: ::prost::alloc::string::String,
}
#[doc = r" Generated client implementations."]
pub mod venue_api_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    pub struct VenueApiClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl VenueApiClient<tonic::transport::Channel> {
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
    impl<T> VenueApiClient<T>
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
        pub async fn get_venue(
            &mut self,
            request: impl tonic::IntoRequest<super::GetVenueRequest>,
        ) -> Result<tonic::Response<super::super::models::Venue>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/GetVenue");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn create_venue(
            &mut self,
            request: impl tonic::IntoRequest<super::CreateVenueRequest>,
        ) -> Result<tonic::Response<super::super::models::Venue>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/CreateVenue");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn get_tables(
            &mut self,
            request: impl tonic::IntoRequest<super::GetTablesRequest>,
        ) -> Result<tonic::Response<super::GetTablesResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/GetTables");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn add_table(
            &mut self,
            request: impl tonic::IntoRequest<super::AddTableRequest>,
        ) -> Result<tonic::Response<super::super::models::Table>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/AddTable");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn remove_table(
            &mut self,
            request: impl tonic::IntoRequest<super::RemoveTableRequest>,
        ) -> Result<tonic::Response<super::super::models::Table>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/RemoveTable");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn is_admin(
            &mut self,
            request: impl tonic::IntoRequest<super::IsAdminRequest>,
        ) -> Result<tonic::Response<super::IsAdminResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/IsAdmin");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn add_admin(
            &mut self,
            request: impl tonic::IntoRequest<super::AddAdminRequest>,
        ) -> Result<tonic::Response<super::AddAdminResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/AddAdmin");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn remove_admin(
            &mut self,
            request: impl tonic::IntoRequest<super::RemoveAdminRequest>,
        ) -> Result<tonic::Response<super::RemoveAdminResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/venue.api.VenueAPI/RemoveAdmin");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
    impl<T: Clone> Clone for VenueApiClient<T> {
        fn clone(&self) -> Self {
            Self {
                inner: self.inner.clone(),
            }
        }
    }
    impl<T> std::fmt::Debug for VenueApiClient<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "VenueApiClient {{ ... }}")
        }
    }
}
#[doc = r" Generated server implementations."]
pub mod venue_api_server {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = "Generated trait containing gRPC methods that should be implemented for use with VenueApiServer."]
    #[async_trait]
    pub trait VenueApi: Send + Sync + 'static {
        async fn get_venue(
            &self,
            request: tonic::Request<super::GetVenueRequest>,
        ) -> Result<tonic::Response<super::super::models::Venue>, tonic::Status>;
        async fn create_venue(
            &self,
            request: tonic::Request<super::CreateVenueRequest>,
        ) -> Result<tonic::Response<super::super::models::Venue>, tonic::Status>;
        async fn get_tables(
            &self,
            request: tonic::Request<super::GetTablesRequest>,
        ) -> Result<tonic::Response<super::GetTablesResponse>, tonic::Status>;
        async fn add_table(
            &self,
            request: tonic::Request<super::AddTableRequest>,
        ) -> Result<tonic::Response<super::super::models::Table>, tonic::Status>;
        async fn remove_table(
            &self,
            request: tonic::Request<super::RemoveTableRequest>,
        ) -> Result<tonic::Response<super::super::models::Table>, tonic::Status>;
        async fn is_admin(
            &self,
            request: tonic::Request<super::IsAdminRequest>,
        ) -> Result<tonic::Response<super::IsAdminResponse>, tonic::Status>;
        async fn add_admin(
            &self,
            request: tonic::Request<super::AddAdminRequest>,
        ) -> Result<tonic::Response<super::AddAdminResponse>, tonic::Status>;
        async fn remove_admin(
            &self,
            request: tonic::Request<super::RemoveAdminRequest>,
        ) -> Result<tonic::Response<super::RemoveAdminResponse>, tonic::Status>;
    }
    #[derive(Debug)]
    pub struct VenueApiServer<T: VenueApi> {
        inner: _Inner<T>,
    }
    struct _Inner<T>(Arc<T>, Option<tonic::Interceptor>);
    impl<T: VenueApi> VenueApiServer<T> {
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
    impl<T, B> Service<http::Request<B>> for VenueApiServer<T>
    where
        T: VenueApi,
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
                "/venue.api.VenueAPI/GetVenue" => {
                    #[allow(non_camel_case_types)]
                    struct GetVenueSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::GetVenueRequest> for GetVenueSvc<T> {
                        type Response = super::super::models::Venue;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetVenueRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).get_venue(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = GetVenueSvc(inner);
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
                "/venue.api.VenueAPI/CreateVenue" => {
                    #[allow(non_camel_case_types)]
                    struct CreateVenueSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::CreateVenueRequest> for CreateVenueSvc<T> {
                        type Response = super::super::models::Venue;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::CreateVenueRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).create_venue(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = CreateVenueSvc(inner);
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
                "/venue.api.VenueAPI/GetTables" => {
                    #[allow(non_camel_case_types)]
                    struct GetTablesSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::GetTablesRequest> for GetTablesSvc<T> {
                        type Response = super::GetTablesResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetTablesRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).get_tables(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = GetTablesSvc(inner);
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
                "/venue.api.VenueAPI/AddTable" => {
                    #[allow(non_camel_case_types)]
                    struct AddTableSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::AddTableRequest> for AddTableSvc<T> {
                        type Response = super::super::models::Table;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::AddTableRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).add_table(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = AddTableSvc(inner);
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
                "/venue.api.VenueAPI/RemoveTable" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveTableSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::RemoveTableRequest> for RemoveTableSvc<T> {
                        type Response = super::super::models::Table;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::RemoveTableRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).remove_table(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = RemoveTableSvc(inner);
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
                "/venue.api.VenueAPI/IsAdmin" => {
                    #[allow(non_camel_case_types)]
                    struct IsAdminSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::IsAdminRequest> for IsAdminSvc<T> {
                        type Response = super::IsAdminResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::IsAdminRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).is_admin(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = IsAdminSvc(inner);
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
                "/venue.api.VenueAPI/AddAdmin" => {
                    #[allow(non_camel_case_types)]
                    struct AddAdminSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::AddAdminRequest> for AddAdminSvc<T> {
                        type Response = super::AddAdminResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::AddAdminRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).add_admin(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = AddAdminSvc(inner);
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
                "/venue.api.VenueAPI/RemoveAdmin" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveAdminSvc<T: VenueApi>(pub Arc<T>);
                    impl<T: VenueApi> tonic::server::UnaryService<super::RemoveAdminRequest> for RemoveAdminSvc<T> {
                        type Response = super::RemoveAdminResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::RemoveAdminRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).remove_admin(request).await };
                            Box::pin(fut)
                        }
                    }
                    let inner = self.inner.clone();
                    let fut = async move {
                        let interceptor = inner.1.clone();
                        let inner = inner.0;
                        let method = RemoveAdminSvc(inner);
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
    impl<T: VenueApi> Clone for VenueApiServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self { inner }
        }
    }
    impl<T: VenueApi> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(self.0.clone(), self.1.clone())
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: VenueApi> tonic::transport::NamedService for VenueApiServer<T> {
        const NAME: &'static str = "venue.api.VenueAPI";
    }
}
