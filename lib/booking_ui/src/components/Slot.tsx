import React from "react";
import { BookingStage } from "./Booking";
import {
  Booking,
  Slot as SlotType,
  SlotInput,
  useCreateBookingMutation,
  useGetSlotQuery,
} from "../graph";
import { Button } from "baseui/button";
import { H2 } from "baseui/typography";
import SlotDisplay from "./SlotDisplay";
import { StyledSpinnerNext } from "baseui/spinner";

interface SlotProps {
  enquiry: SlotInput;
  setBookingStage: React.Dispatch<React.SetStateAction<BookingStage>>;
  setBooking: React.Dispatch<React.SetStateAction<Booking | null>>;
}

const Slot: React.FC<SlotProps> = ({
  enquiry,
  setBooking,
  setBookingStage,
}) => {
  const { data, loading, error } = useGetSlotQuery({
    variables: {
      slot: enquiry,
    },
  });

  if (loading)
    return (
      <div>
        <StyledSpinnerNext />
      </div>
    );

  const match = data?.getSlot.match;

  if (error || !match) {
    return (
      <div>
        <H2>sorry we could not find a slot</H2>
        <br />
        <Button
          onClick={(e) => {
            e.preventDefault();
            setBookingStage(BookingStage.Enquiry);
          }}
        >
          start again
        </Button>
      </div>
    );
  }

  return (
    <div>
      <H2>we found a slot!</H2>
      <SlotDisplay {...match} />
      <br />
      <CreateBookingButton
        match={match}
        setBookingStage={setBookingStage}
        setBooking={setBooking}
      />
    </div>
  );
};

const CreateBookingButton: React.FC<{
  match: SlotType;
  setBookingStage: React.Dispatch<React.SetStateAction<BookingStage>>;
  setBooking: React.Dispatch<React.SetStateAction<Booking | null>>;
}> = ({ match, setBooking, setBookingStage }) => {
  const [createBookingMutation] = useCreateBookingMutation({
    variables: {
      slot: {
        venueId: match.venueId,
        email: match.email,
        people: match.people,
        startsAt: match.startsAt,
        duration: match.duration,
      },
    },
  });
  const handleClick = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    createBookingMutation()
      .then((r) => {
        setBooking(r?.data?.createBooking!);
        setBookingStage(BookingStage.Confirmation);
      })
      .catch((e) => {
        console.log(e);
      });
  };
  return (
    <div>
      <Button
        onClick={(e) => {
          e.preventDefault();
          setBookingStage(BookingStage.Enquiry);
        }}
      >
        Back
      </Button>{" "}
      <Button onClick={handleClick}>Confirm</Button>
    </div>
  );
};

export default Slot;
