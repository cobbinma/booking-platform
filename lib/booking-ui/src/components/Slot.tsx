import React from "react";
import { AxiosRequestConfig } from "axios";
import { BookingQuery, BookingSlot } from "./pages/Book";
import useAxios from "axios-hooks";

const Slot = ({ bookingQuery }: { bookingQuery: BookingQuery }) => {
  const [bookingSlot, setBookingSlot] = React.useState<BookingSlot | null>();
  const request: AxiosRequestConfig = {
    url: "http://localhost:6969/booking",
    method: "GET",
    data: bookingQuery,
  };
  const [{ data, loading, error }] = useAxios(request);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error!</p>;

  return <div>Slot: {data}</div>;
};

export default Slot;
