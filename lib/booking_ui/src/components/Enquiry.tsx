import React from "react";
import { DatePicker } from "baseui/datepicker";
import { Slot, SlotInput, useCreateSlotMutation } from "../graph";
import { BookingStage } from "./Booking";
import { TimePicker } from "baseui/timepicker";
import { Slider } from "baseui/slider";
import { Combobox } from "baseui/combobox";
import { Button } from "baseui/button";
import { Label1 } from "baseui/typography";

let durations = new Map<string, number>([
  ["30 mins", 30],
  ["1 hour", 60],
  ["1 hour 30 mins", 90],
  ["2 hours", 120],
  ["2 hours 30 mins", 150],
  ["3 hours", 180],
]);

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
  const [date, setDate] = React.useState([
    new Date("2021-01-12T00:00:00.000Z"),
  ]);
  const [time, setTime] = React.useState(new Date("2021-01-22T00:01:42.445Z"));
  const [people, setPeople] = React.useState([4]);
  const [duration, setDuration] = React.useState<string>("1 hour");
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

  const handleClick = (e: any) => {
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
      <Label1>Date</Label1>
      <DatePicker
        value={date}
        onChange={({ date }) => setDate(Array.isArray(date) ? date : [date])}
      />
      <Label1>Time</Label1>
      <TimePicker value={time} onChange={(date) => setTime(date)} />
      <Label1>Guests</Label1>
      <Slider
        value={people}
        onChange={({ value }) => value && setPeople(value)}
        onFinalChange={({ value }) => console.log(value)}
        min={1}
        max={20}
      />
      <Label1>Duration</Label1>
      <Combobox
        value={duration}
        onChange={(nextValue) => setDuration(nextValue)}
        options={Array.from(durations.keys())}
        mapOptionToString={(option) => option}
      />
      <Button onClick={handleClick}>Next</Button>
    </div>
  );
};

export default Enquiry;
