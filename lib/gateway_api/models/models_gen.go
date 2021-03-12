// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

// Input to add an administrator to a venue.
type AdminInput struct {
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// email address of the administrator
	Email string `json:"email"`
}

// Booking has now been confirmed.
type Booking struct {
	// unique identifier of the booking
	ID string `json:"id"`
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// email of the customer
	Email string `json:"email"`
	// amount of people attending the booking
	People int `json:"people"`
	// start time of the booking (hh:mm)
	StartsAt time.Time `json:"startsAt"`
	// end time of the booking (hh:mm)
	EndsAt time.Time `json:"endsAt"`
	// duration of the booking in minutes
	Duration int `json:"duration"`
	// unique identifier of the booking table
	TableID string `json:"tableId"`
}

// Slot is a possible booking that has yet to be confirmed.
type BookingInput struct {
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// email of the customer
	Email string `json:"email"`
	// amount of people attending the booking
	People int `json:"people"`
	// start time of the booking (YYYY-MM-DDThh:mm:ssZ)
	StartsAt time.Time `json:"startsAt"`
	// duration of the booking in minutes
	Duration int `json:"duration"`
}

// Filter bookings.
type BookingsFilter struct {
	// unique identifier of the venue
	VenueID *string `json:"venueId"`
	// specific date to query bookings for
	Date time.Time `json:"date"`
}

// A page with a list of bookings.
type BookingsPage struct {
	// list of bookings
	Bookings []*Booking `json:"bookings"`
	// is there a next page
	HasNextPage bool `json:"hasNextPage"`
	// total number of pages
	Pages int `json:"pages"`
}

// Input to cancel an individual booking.
type CancelBookingInput struct {
	// unique identifier of the venue
	VenueID *string `json:"venueId"`
	// unique identifier of the booking
	ID string `json:"id"`
}

// Booking Enquiry Response.
type GetSlotResponse struct {
	// slot matching the given enquiy
	Match *Slot `json:"match"`
	// slots have match the enquiry but have different starting times
	OtherAvailableSlots []*Slot `json:"otherAvailableSlots"`
}

// Input to query if the user is an admin. Fields AND together.
type IsAdminInput struct {
	// unique identifier of the venue
	VenueID *string `json:"venueId"`
	// human readable identifier of the venue
	Slug *string `json:"slug"`
}

// Day specific operating hours.
type OpeningHoursSpecification struct {
	// the day of the week for which these opening hours are valid
	DayOfWeek DayOfWeek `json:"dayOfWeek"`
	// the opening time of the place or service on the given day(s) of the week
	Opens *TimeOfDay `json:"opens"`
	// the closing time of the place or service on the given day(s) of the week
	Closes *TimeOfDay `json:"closes"`
	// date the special opening hours starts at. only valid for special opening hours
	ValidFrom *time.Time `json:"validFrom"`
	// date the special opening hours ends at. only valid for special opening hours
	ValidThrough *time.Time `json:"validThrough"`
}

// Day specific operating hours.
type OpeningHoursSpecificationInput struct {
	// the day of the week for which these opening hours are valid
	DayOfWeek DayOfWeek `json:"dayOfWeek"`
	// the opening time of the place or service on the given day(s) of the week
	Opens TimeOfDay `json:"opens"`
	// the closing time of the place or service on the given day(s) of the week
	Closes TimeOfDay `json:"closes"`
}

// Information about the page being requested. Maximum page limit of 50.
type PageInfo struct {
	// page number
	Page int `json:"page"`
	// maximum amount of results per page
	Limit *int `json:"limit"`
}

// Input to remove an administrator from a venue.
type RemoveAdminInput struct {
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// email address of the administrator
	Email string `json:"email"`
}

// Input to remove a venue table
type RemoveTableInput struct {
	// unique venue identifier the table belongs to
	VenueID string `json:"venueId"`
	// unique identifier of the table to be removed
	TableID string `json:"tableId"`
}

// Slot is a possible booking that has yet to be confirmed.
type Slot struct {
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// email of the customer
	Email string `json:"email"`
	// amount of people attending the booking
	People int `json:"people"`
	// desired start time of the booking (YYYY-MM-DDThh:mm:ssZ)
	StartsAt time.Time `json:"startsAt"`
	// potential ending time of the booking (YYYY-MM-DDThh:mm:ssZ)
	EndsAt time.Time `json:"endsAt"`
	// potential duration of the booking in minutes
	Duration int `json:"duration"`
}

// Slot Input is a booking enquiry.
type SlotInput struct {
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// email of the customer
	Email string `json:"email"`
	// amount of people attending the booking
	People int `json:"people"`
	// desired start time of the booking (YYYY-MM-DDThh:mm:ssZ)
	StartsAt time.Time `json:"startsAt"`
	// desired duration of the booking in minutes
	Duration int `json:"duration"`
}

// Day specific special operating hours.
type SpecialOpeningHoursSpecificationInput struct {
	// the day of the week for which these opening hours are valid
	DayOfWeek DayOfWeek `json:"dayOfWeek"`
	// the opening time of the place or service on the given day(s) of the week
	Opens *TimeOfDay `json:"opens"`
	// the closing time of the place or service on the given day(s) of the week
	Closes *TimeOfDay `json:"closes"`
	// date the special opening hours starts at. only valid for special opening hours
	ValidFrom time.Time `json:"validFrom"`
	// date the special opening hours ends at. only valid for special opening hours
	ValidThrough time.Time `json:"validThrough"`
}

// An individual table at a venue.
type Table struct {
	// unique identifier of the table
	ID string `json:"id"`
	// name of the table
	Name string `json:"name"`
	// maximum amount of people that can sit at table
	Capacity int `json:"capacity"`
}

// An individual table at a venue.
type TableInput struct {
	// unique venue identifier the table belongs to
	VenueID string `json:"venueId"`
	// name of the table
	Name string `json:"name"`
	// maximum amount of people that can sit at table
	Capacity int `json:"capacity"`
}

// Input to update a venue's operating hours.
type UpdateOpeningHoursInput struct {
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// operating hours of the venue
	OpeningHours []*OpeningHoursSpecificationInput `json:"openingHours"`
}

// Input to update a venue's special operating hours.
type UpdateSpecialOpeningHoursInput struct {
	// unique identifier of the venue
	VenueID string `json:"venueId"`
	// special operating hours of the venue
	SpecialOpeningHours []*SpecialOpeningHoursSpecificationInput `json:"specialOpeningHours"`
}

// Venue where a booking can take place.
type Venue struct {
	// unique identifier of the venue
	ID string `json:"id"`
	// name of the venue
	Name string `json:"name"`
	// operating hours of the venue
	OpeningHours []*OpeningHoursSpecification `json:"openingHours"`
	// special operating hours of the venue
	SpecialOpeningHours []*OpeningHoursSpecification `json:"specialOpeningHours"`
	// operating hours of the venue for a specific date
	OpeningHoursSpecification *OpeningHoursSpecification `json:"openingHoursSpecification"`
	// tables at the venue
	Tables []*Table `json:"tables"`
	// email addresses of venue administrators
	Admins []string `json:"admins"`
	// human readable identifier of the venue
	Slug string `json:"slug"`
	// paginated list of bookings for a venue
	Bookings *BookingsPage `json:"bookings"`
}

// Filter get venue queries. Fields AND together.
type VenueFilter struct {
	// unique identifier of the venue
	ID *string `json:"id"`
	// human readable identifier of the venue
	Slug *string `json:"slug"`
}
