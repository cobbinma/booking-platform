import React, { useState } from "react";
import { Params } from "../App";
import LogoutButton from "./LogoutButton";
import Venue from "./Venue";
import Enquiry from "./Enquiry";
import Slot from "./Slot";
import Confirmation from "./Confirmation";
import { Booking as BookingType, Slot as SlotType } from "../graph";

export enum BookingStage {
  Enquiry = 1,
  Slot,
  Confirmation,
}

interface BookingProps {
  params: Params;
  email: string;
}

const Booking: React.FC<BookingProps> = ({ params, email }) => {
  const [stage, setStage] = useState<BookingStage>(BookingStage.Enquiry);
  const [slot, setSlot] = useState<SlotType | null>(null);
  const [booking, setBooking] = useState<BookingType | null>(null);
  const { venueId, returnURL } = params;

  const getStageComponent = (stage: BookingStage): React.ReactElement => {
    switch (stage) {
      case BookingStage.Enquiry:
        return (
          <Enquiry
            setBookingStage={setStage}
            setSlot={setSlot}
            venueId={venueId}
            email={email}
          />
        );
      case BookingStage.Slot:
        return (
          <Slot
            slot={slot}
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
      default:
        return <div>error: unknown stage</div>;
    }
  };

  return (
    <div>
      <LogoutButton />
      {getStageComponent(stage)}
      <Venue venueId={venueId} />
    </div>
  );
};

export default Booking;
