import React from "react";
import { H2 } from "baseui/typography";
import { Button } from "baseui/button";
import { BookingStage } from "./Booking";

const Error: React.FC<{
  setBookingStage: React.Dispatch<React.SetStateAction<BookingStage>>;
}> = ({ setBookingStage }) => {
  return (
    <div>
      <H2>oops! an error happened...</H2>
      <br />
      <Button
        onClick={(e) => {
          e.preventDefault();
          setBookingStage(BookingStage.Enquiry);
        }}
      >
        try again
      </Button>
    </div>
  );
};

export default Error;
