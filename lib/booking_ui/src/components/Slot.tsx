import React from "react";
import { BookingStage } from "./Booking";
import { Booking, Slot as SlotType, useCreateBookingMutation } from "../graph";

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
        sorry we could not find a slot
        <button
          onClick={(e) => {
            e.preventDefault();
            setBookingStage(BookingStage.Enquiry);
          }}
        >
          start again
        </button>
      </div>
    );
  }

  const handleClick = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    createBookingMutation()
      .then((b) => {
        console.log(slot);
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
      <h2>we found a slot!</h2>
      <p>
        table for {slot?.people} people, starting at {slot?.startsAt}, ending at{" "}
        {slot?.endsAt}
      </p>
      <button onClick={handleClick}>Confirm</button>
    </div>
  );
};

export default Slot;
