import React from "react";
import { Slot, SlotInput, useCreateSlotMutation } from "../graph";
import { BookingStage } from "./Booking";

interface EnquiryProps {
  setBookingStage: React.Dispatch<React.SetStateAction<BookingStage>>;
  setSlot: React.Dispatch<React.SetStateAction<Slot | null>>;
  venueId: string;
  email: string;
}

const Enquiry: React.FC<EnquiryProps> = ({
  venueId,
  email,
  setSlot,
  setBookingStage,
}) => {
  let enquiry: SlotInput = {
    venueId: venueId,
    customerId: email,
    people: 4,
    date: "21-01-2021",
    startsAt: "18:00",
    duration: 60,
  };
  const [createSlotMutation] = useCreateSlotMutation({
    variables: {
      slot: enquiry,
    },
  });

  const handleClick = (e: { preventDefault: () => void }) => {
    e.preventDefault();
    createSlotMutation()
      .then((r) => {
        setSlot(r?.data?.createSlot || null);
      })
      .catch((e) => {
        setSlot(null);
        console.log(e);
      });
    setBookingStage(BookingStage.Slot);
  };

  return (
    <div>
      Enquiry
      <button onClick={handleClick}>Next</button>
    </div>
  );
};

export default Enquiry;
