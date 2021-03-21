import React, { useState } from "react";
import { Params } from "../App";
import LogoutButton from "./LogoutButton";
import Enquiry from "./Enquiry";
import Slot from "./Slot";
import Confirmation from "./Confirmation";
import { Booking as BookingType, SlotInput, useGetVenueQuery } from "../graph";
import { Button } from "baseui/button";
import { StyledSpinnerNext } from "baseui/spinner";

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
  const [enquiry, setEnquiry] = useState<SlotInput | null>(null);
  const [booking, setBooking] = useState<BookingType | null>(null);
  const { slug, returnURL } = params;

  const { data, loading, error, refetch } = useGetVenueQuery({
    variables: {
      slug: slug,
    },
  });

  if (loading) return <StyledSpinnerNext />;

  if (error) {
    console.log(error);
    return <div>error</div>;
  }

  const getStageComponent = (stage: BookingStage): React.ReactElement => {
    switch (stage) {
      case BookingStage.Enquiry:
        return (
          <Enquiry
            setBookingStage={setStage}
            setEnquiry={setEnquiry}
            venueId={data?.getVenue?.id || ""}
            email={email}
            openingHours={data?.getVenue?.openingHoursSpecification}
            refetch={refetch}
          />
        );
      case BookingStage.Slot:
        if (enquiry == null) return <div>error</div>;
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
      default:
        return <div>error</div>;
    }
  };

  return (
    <div>
      {getStageComponent(stage)}
      <br />
      <a href={decodeURIComponent(returnURL)}>
        <Button type="button">Exit</Button>
      </a>{" "}
      <LogoutButton />
    </div>
  );
};

export default Booking;
