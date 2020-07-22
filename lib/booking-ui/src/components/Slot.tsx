import React from "react";
import { AxiosRequestConfig } from "axios";
import { BookingQuery, BookingSlot } from "./pages/Book";
import useAxios from "axios-hooks";

const Slot = ({ bookingQuery }: { bookingQuery: BookingQuery }) => {
  const [bookingSlot, setBookingSlot] = React.useState<BookingSlot | null>();
  console.log(bookingQuery);
  const request: AxiosRequestConfig = {
    url: "http://localhost:6969/slot",
    method: "POST",
    data: bookingQuery,
  };
  const [{ data, loading, error }] = useAxios(request);

  if (loading) return <p>Loading...</p>;
  if (error) {
    console.log(error.response);
    return <p>Error: {error.response?.data.message}</p>;
  }

  return <div>Slot: {data}</div>;
};

export default Slot;
