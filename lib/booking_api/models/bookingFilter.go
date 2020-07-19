package models

type BookingFilter struct {
	Date     *Date
	TableIDs []TableID
}

func BookingFilterWithDate(date *Date) func(*BookingFilter) *BookingFilter {
	return func(filter *BookingFilter) *BookingFilter {
		filter.Date = date
		return filter
	}
}

func BookingFilterWithTableIDs(tableIDs []TableID) func(*BookingFilter) *BookingFilter {
	return func(filter *BookingFilter) *BookingFilter {
		filter.TableIDs = tableIDs
		return filter
	}
}
