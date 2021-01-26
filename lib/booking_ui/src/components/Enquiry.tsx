import React from "react";
import { DatePicker } from "baseui/datepicker";
import { SlotInput } from "../graph";
import { BookingStage } from "./Booking";
import { TimePicker } from "baseui/timepicker";
import { Slider } from "baseui/slider";
import { Combobox } from "baseui/combobox";
import { Button } from "baseui/button";
import { H2, Label1 } from "baseui/typography";

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
  setEnquiry: React.Dispatch<React.SetStateAction<SlotInput | null>>;
  venueId: string;
  email: string;
}

const Enquiry: React.FC<EnquiryProps> = ({
  venueId,
  email,
  setEnquiry,
  setBookingStage,
}) => {
  const [date, setDate] = React.useState([new Date(Date.now())]);
  const [time, setTime] = React.useState(new Date(Date.now()));
  const [people, setPeople] = React.useState([4]);
  const [duration, setDuration] = React.useState<string>("1 hour");
  const enquiry: SlotInput = {
    venueId: venueId,
    email: email,
    people: people[0],
    startsAt: new Date(
      date[0].getFullYear(),
      date[0].getMonth(),
      date[0].getDate(),
      time.getHours(),
      time.getMinutes()
    ),
    duration: durations.get(duration) || 60,
  };

  const handleClick = (e: any) => {
    e.preventDefault();
    if (enquiry) {
      setEnquiry(enquiry);
    }
    setBookingStage(BookingStage.Slot);
  };

  return (
    <div>
      <H2>book a table</H2>
      <Label1>Date</Label1>
      <DatePicker
        value={date}
        onChange={({ date }) => setDate(Array.isArray(date) ? date : [date])}
      />
      <Label1>Time</Label1>
      <TimePicker value={time} step={1800} onChange={(date) => setTime(date)} />
      <Label1>Guests</Label1>
      <Slider
        value={people}
        onChange={({ value }) => value && setPeople(value)}
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
      <br />
      <Button onClick={handleClick}>Next</Button>
    </div>
  );
};

export default Enquiry;
