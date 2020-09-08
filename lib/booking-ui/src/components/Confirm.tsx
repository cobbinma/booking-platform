import React, { useEffect } from "react";
import { BookingSlot } from "./pages/Book";
import useAxios from "axios-hooks";

const Confirm = ({
  bookingSlot,
}: {
  bookingSlot: BookingSlot | null | undefined;
}) => {
  const [{ loading, error }] = useAxios({
    url: "http://localhost:6969/venues/1/bookings",
    method: "POST",
    data: bookingSlot,
  });

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.response?.data.message}</p>;

  return (
    <div>
      <h3>Booked!</h3>
      <p>Your table is confirmed.</p>
    </div>
  );
};

export default Confirm;
