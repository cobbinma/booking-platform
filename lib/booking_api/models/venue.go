package models

type VenueID string

func NewVenueID(id string) VenueID {
	return VenueID(id)
}

func (vid VenueID) String() string {
	return string(vid)
}
