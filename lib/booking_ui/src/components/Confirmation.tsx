import React from "react";
import { Booking } from "../graph";
import { BookingStage } from "./Booking";
import { Button } from "baseui/button";
import { H2 } from "baseui/typography";
import { Table } from "baseui/table";

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
        <H2>sorry something went wrong</H2>
        <br />
        <Button
          onClick={(e) => {
            e.preventDefault();
            setBookingStage(BookingStage.Enquiry);
          }}
        >
          Try Again
        </Button>
      </div>
    );
  }

  return (
    <div>
      <H2>confirmed!</H2>
      <Table
        columns={["Guests", "Date", "Starts", "Ends"]}
        data={[
          [booking.people, booking.date, booking.startsAt, booking.endsAt],
        ]}
      />
      <br />
      <a href={decodeURIComponent(returnURL)}>
        <Button type="button">Continue</Button>
      </a>
    </div>
  );
};

export default Confirmation;
