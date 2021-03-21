import React from "react";
import { DatePicker } from "baseui/datepicker";
import { OpeningHoursSpecification, SlotInput } from "../graph";
import { BookingStage } from "./Booking";
import { TimePicker } from "baseui/timepicker";
import { Slider } from "baseui/slider";
import { Combobox } from "baseui/combobox";
import { Button } from "baseui/button";
import { H2 } from "baseui/typography";
import { FormControl } from "baseui/form-control";
import { GetVenueQuery, GetVenueQueryVariables } from "../graph";
import { ApolloQueryResult } from "@apollo/client";

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
  openingHours: OpeningHoursSpecification | null | undefined;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}

const Enquiry: React.FC<EnquiryProps> = ({
  venueId,
  email,
  setEnquiry,
  setBookingStage,
  openingHours,
  refetch,
}) => {
  const [date, setDate] = React.useState<Date[] | null>(null);
  const [time, setTime] = React.useState<Date | null>(null);
  const [people, setPeople] = React.useState([4]);
  const [duration, setDuration] = React.useState<string | null>(null);
  const enquiry: SlotInput = {
    venueId: venueId,
    email: email,
    people: people[0],
    startsAt:
      date && date[0] && time
        ? new Date(
            date[0].getFullYear(),
            date[0].getMonth(),
            date[0].getDate(),
            time.getHours(),
            time.getMinutes()
          )
        : [new Date(Date.now())],
    duration: durations.get(duration || "") || 60,
  };

  const handleClick = (e: any) => {
    e.preventDefault();
    if (enquiry) {
      setEnquiry(enquiry);
    }
    setBookingStage(BookingStage.Slot);
  };

  const venue_is_open = (
    hours: OpeningHoursSpecification | null | undefined
  ): boolean => {
    return !!(hours && hours.opens && hours.closes);
  };

  const time_of_day_is_before_date = (
    time_of_day: string,
    date: Date,
    addMinutes?: number
  ): boolean => {
    const splitTime = time_of_day.split(":");
    if (splitTime.length !== 2) return false;
    const day = new Date(
      date.getFullYear(),
      date.getMonth(),
      date.getDate(),
      parseInt(splitTime[0]),
      parseInt(splitTime[1])
    );
    if (addMinutes) date = new Date(date.getTime() + addMinutes * 60000);
    return day >= date;
  };

  const venue_closed_at_time = (): boolean => {
    return (
      time_of_day_is_before_date(openingHours?.opens, time || new Date()) ||
      !time_of_day_is_before_date(openingHours?.closes, time || new Date())
    );
  };

  const get_available_durations = (time: Date): string[] => {
    return Array.from(durations.keys()).filter((k) => {
      const duration = durations.get(k);
      return (
        duration &&
        time_of_day_is_before_date(openingHours?.closes, time, duration)
      );
    });
  };

  return (
    <div>
      <H2>book a table</H2>
      <FormControl
        label="Date"
        caption={() =>
          date && !venue_is_open(openingHours) ? "venue is closed" : ""
        }
      >
        <DatePicker
          value={date}
          onChange={({ date }) => {
            const d: Date[] = Array.isArray(date) ? date : [date];
            setDate(d);
            if (d && d[0]) {
              refetch({
                date: d[0].toISOString(),
                venueId: venueId,
              }).catch((e) => console.log(e));
            }
          }}
          error={!!(date && !venue_is_open(openingHours))}
        />
      </FormControl>
      {date && venue_is_open(openingHours) && (
        <div>
          <FormControl
            label="Time"
            caption={() =>
              venue_closed_at_time() ||
              !time ||
              get_available_durations(time).length === 0
                ? "venue is closed"
                : ""
            }
          >
            <TimePicker
              value={time}
              step={1800}
              onChange={(date) => setTime(date)}
              error={venue_closed_at_time()}
            />
          </FormControl>
          {time &&
            !venue_closed_at_time() &&
            get_available_durations(time).length > 0 && (
              <div>
                <FormControl label="Guests">
                  <Slider
                    value={people}
                    onChange={({ value }) => value && setPeople(value)}
                    min={1}
                    max={20}
                  />
                </FormControl>
                <FormControl label="Duration">
                  <Combobox
                    value={duration || ""}
                    onChange={(nextValue) => setDuration(nextValue)}
                    options={get_available_durations(time)}
                    mapOptionToString={(option) => option}
                  />
                </FormControl>
              </div>
            )}
        </div>
      )}
      <br />
      <Button
        disabled={
          !(
            time &&
            venue_is_open(openingHours) &&
            !venue_closed_at_time() &&
            duration
          )
        }
        onClick={handleClick}
      >
        Next
      </Button>
    </div>
  );
};

export default Enquiry;
