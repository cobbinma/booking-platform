import React, { useState } from "react";
import { Params } from "../App";
import LogoutButton from "./LogoutButton";
import Venue from "./Venue";
import Enquiry from "./Enquiry";
import Slot from "./Slot";
import Confirmation from "./Confirmation";
import { Booking as BookingType, SlotInput } from "../graph";
import { Button } from "baseui/button";
import Error from "./Error";

export enum BookingStage {
  Enquiry = 1,
  Slot,
  Confirmation,
  Error,
}

interface BookingProps {
  params: Params;
  email: string;
}

const Booking: React.FC<BookingProps> = ({ params, email }) => {
  const [stage, setStage] = useState<BookingStage>(BookingStage.Enquiry);
  const [enquiry, setEnquiry] = useState<SlotInput | null>(null);
  const [booking, setBooking] = useState<BookingType | null>(null);
  const { venueId, returnURL } = params;

  const getStageComponent = (stage: BookingStage): React.ReactElement => {
    switch (stage) {
      case BookingStage.Enquiry:
        return (
          <Enquiry
            setBookingStage={setStage}
            setEnquiry={setEnquiry}
            venueId={venueId}
            email={email}
          />
        );
      case BookingStage.Slot:
        if (enquiry == null) return <Error setBookingStage={setStage} />;
        return (
          <Slot
            enquiry={enquiry}
            setBooking={setBooking}
            setBookingStage={setStage}
          />
        );
      case BookingStage.Confirmation:
        return (
          <Confirmation
            booking={booking}
            setBookingStage={setStage}
            returnURL={returnURL}
          />
        );
      case BookingStage.Error:
        return <Error setBookingStage={setStage} />;
      default:
        return <div>error: unknown stage</div>;
    }
  };

  return (
    <div>
      {getStageComponent(stage)}
      <Venue venueId={venueId} />
      <br />
      <a href={decodeURIComponent(returnURL)}>
        <Button type="button">Exit</Button>
      </a>{" "}
      <LogoutButton />
    </div>
  );
};

export default Booking;
