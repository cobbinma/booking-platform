pub mod venue {
    pub mod models {
        include!("venue.models.rs");
    }
    pub mod api {
        include!("venue.api.rs");
    }
}
pub mod booking {
    pub mod api {
        include!("booking.api.rs");
    }
    pub mod models {
        include!("booking.models.rs");
    }
}
