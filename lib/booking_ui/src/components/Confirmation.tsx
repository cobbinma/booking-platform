import React from "react";
import { Booking } from "../graph";
import { BookingStage } from "./Booking";

interface ConfirmationProps {
  booking: Booking | null;
  setBookingStage: React.Dispatch<React.SetStateAction<BookingStage>>;
  returnURL: string;
}

const Confirmation: React.FC<ConfirmationProps> = ({
  booking,
  setBookingStage,
  returnURL,
}) => {
  if (booking == null) {
    return (
      <div>
        sorry something went wrong
        <button
          onClick={(e) => {
            e.preventDefault();
            setBookingStage(BookingStage.Enquiry);
          }}
        >
          try again
        </button>
      </div>
    );
  }

  return (
    <div>
      <h2>Confirmed!</h2>
      <p>
        table for {booking?.people} people, starting at {booking?.startsAt},
        ending at {booking?.endsAt}
      </p>
      <a href={decodeURIComponent(returnURL)}>
        <button type="button">Continue</button>
      </a>
    </div>
  );
};

export default Confirmation;
