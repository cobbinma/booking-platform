import React, { useEffect } from "react";
import { BookingQuery, BookingSlot } from "./pages/Book";
import useAxios from "axios-hooks";

const Slot = ({
  bookingQuery,
  bookingSlot,
  setBookingSlot,
}: {
  bookingQuery: BookingQuery;
  bookingSlot: BookingSlot | null | undefined;
  setBookingSlot: React.Dispatch<
    React.SetStateAction<BookingSlot | null | undefined>
  >;
}) => {
  const date = bookingQuery.date;
  const starts_at = new Date(
    date.getFullYear(),
    date.getMonth(),
    date.getDate(),
    bookingQuery.starts_at.getHours(),
    bookingQuery.starts_at.getMinutes(),
    bookingQuery.starts_at.getSeconds(),
    bookingQuery.starts_at.getMilliseconds()
  );
  const ends_at = new Date(
    date.getFullYear(),
    date.getMonth(),
    date.getDate(),
    starts_at.getHours() + bookingQuery.duration,
    starts_at.getMinutes(),
    starts_at.getSeconds(),
    starts_at.getMilliseconds()
  );

  const [{ data, loading, error }] = useAxios({
    url: "http://localhost:6969/slot",
    method: "POST",
    data: {
      ...bookingQuery,
      date: date?.toISOString().slice(0, 10),
      starts_at: starts_at.toISOString(),
      ends_at: ends_at.toISOString(),
    },
  });

  useEffect(() => {
    if (data) {
      setBookingSlot(data);
    }
  }, [data]);

  if (loading) return <p>Loading...</p>;
  if (error) {
    if (error.response?.status === 404) {
      return (
        <div>
          <h3>No Table Available</h3>
          <p>Sorry we couldn't find you a table. Click back to try again.</p>
        </div>
      );
    }
    console.log(error.response?.data);
    return (
      <div>
        <h3>Oops</h3>
        <p>Sorry a weird error has occurred. Click back to try again.</p>
      </div>
    );
  }

  return (
    <div>
      <h3>We found a table</h3>
      Starts at: {bookingQuery.starts_at.toUTCString()}
      <br />
      Duration: {bookingQuery.duration.toString()} Hours
      <br />
      Guests: {bookingQuery.people.toString()}
      <br />
    </div>
  );
};

export default Slot;
