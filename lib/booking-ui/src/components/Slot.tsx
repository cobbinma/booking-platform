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
  if (error) return <p>Error!</p>;

  return (
    <div>
      <pre>{JSON.stringify(bookingSlot, null, 2)}</pre>
    </div>
  );
};

export default Slot;
