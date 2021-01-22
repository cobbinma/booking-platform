import React from "react";
import { BookingStage } from "./Booking";
import { Booking, Slot as SlotType, useCreateBookingMutation } from "../graph";
import { Button } from "baseui/button";
import { H2 } from "baseui/typography";
import { Table } from "baseui/table";

interface SlotProps {
  slot: SlotType | null;
  setBookingStage: React.Dispatch<React.SetStateAction<BookingStage>>;
  setBooking: React.Dispatch<React.SetStateAction<Booking | null>>;
}

const Slot: React.FC<SlotProps> = ({ slot, setBooking, setBookingStage }) => {
  const [createBookingMutation] = useCreateBookingMutation({
    variables: {
      slot: {
        venueId: slot?.venueId!,
        customerId: slot?.customerId!,
        people: slot?.people!,
        date: slot?.date!,
        startsAt: slot?.startsAt!,
        duration: slot?.duration!,
      },
    },
  });

  if (slot == null) {
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

  const handleClick = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    createBookingMutation()
      .then((b) => {
        setBooking(b?.data?.createBooking || null);
      })
      .catch((e) => {
        console.log(e);
        setBooking(null);
      });
    setBookingStage(BookingStage.Confirmation);
  };

  return (
    <div>
      <H2>we found a slot!</H2>
      <Table
        columns={["Guests", "Date", "Starts", "Ends"]}
        data={[[slot.people, slot.date, slot.startsAt, slot.endsAt]]}
      />
      <br />
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
