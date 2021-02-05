#[derive(Clone, PartialEq, ::prost::Message)]
pub struct IsAdminRequest {
    #[prost(string, tag = "1")]
    pub venue_id: ::prost::alloc::string::String,
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
            let path = http::uri::PathAndQuery::from_static("/customer.api.BookingAPI/IsAdmin");
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
            let path = http::uri::PathAndQuery::from_static("/customer.api.BookingAPI/AddAdmin");
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
            let path = http::uri::PathAndQuery::from_static("/customer.api.BookingAPI/RemoveAdmin");
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
                "/customer.api.BookingAPI/IsAdmin" => {
                    #[allow(non_camel_case_types)]
                    struct IsAdminSvc<T: BookingApi>(pub Arc<T>);
                    impl<T: BookingApi> tonic::server::UnaryService<super::IsAdminRequest> for IsAdminSvc<T> {
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
                "/customer.api.BookingAPI/AddAdmin" => {
                    #[allow(non_camel_case_types)]
                    struct AddAdminSvc<T: BookingApi>(pub Arc<T>);
                    impl<T: BookingApi> tonic::server::UnaryService<super::AddAdminRequest> for AddAdminSvc<T> {
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
                "/customer.api.BookingAPI/RemoveAdmin" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveAdminSvc<T: BookingApi>(pub Arc<T>);
                    impl<T: BookingApi> tonic::server::UnaryService<super::RemoveAdminRequest> for RemoveAdminSvc<T> {
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
        const NAME: &'static str = "customer.api.BookingAPI";
    }
}
