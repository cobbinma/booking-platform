import React from "react";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";
import DateFnsUtils from "@date-io/date-fns";
import {
  MuiPickersUtilsProvider,
  KeyboardDatePicker,
  KeyboardTimePicker,
} from "@material-ui/pickers";
import { InputLabel, MenuItem, Select } from "@material-ui/core";
import { BookingQuery } from "./pages/Book";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      "& > *": {
        margin: theme.spacing(2),
        width: "25ch",
        display: "block",
      },
    },
  })
);

export default function RequestSlotForm({
  bookingQuery,
  setBookingQuery,
}: {
  bookingQuery: BookingQuery;
  setBookingQuery: React.Dispatch<React.SetStateAction<BookingQuery>>;
}) {
  const classes = useStyles();

  const handleEmailChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    setBookingQuery({
      ...bookingQuery,
      customer_id: event.target.value as string,
    });
  };
  const handleDateChange = (date: Date | null) => {
    setBookingQuery({
      ...bookingQuery,
      date,
    });
  };

  const handleStartTimeChange = (starts_at: Date | null) => {
    setBookingQuery({
      ...bookingQuery,
      starts_at,
    });
  };

  const handleDurationChange = (
    event: React.ChangeEvent<{ value: unknown }>
  ) => {
    setBookingQuery({
      ...bookingQuery,
      duration: event.target.value as number,
    });
  };
  const handlePeopleChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    setBookingQuery({
      ...bookingQuery,
      people: event.target.value as number,
    });
  };

  return (
    <form className={classes.root} noValidate autoComplete="off">
      <TextField
        id="email"
        label="Email"
        value={bookingQuery.customer_id}
        onChange={handleEmailChange}
      />
      <MuiPickersUtilsProvider utils={DateFnsUtils}>
        <KeyboardDatePicker
          margin="normal"
          id="date"
          label="Date"
          format="MM/dd/yyyy"
          value={bookingQuery.date}
          onChange={handleDateChange}
          KeyboardButtonProps={{
            "aria-label": "change date",
          }}
        />
        <KeyboardTimePicker
          margin="normal"
          id="starts-at"
          label="Starting Time"
          value={bookingQuery.starts_at}
          onChange={handleStartTimeChange}
          KeyboardButtonProps={{
            "aria-label": "change time",
          }}
        />
      </MuiPickersUtilsProvider>
      <InputLabel id="duration">Duration</InputLabel>
      <Select
        labelId="duration"
        id="duration"
        value={bookingQuery.duration}
        onChange={handleDurationChange}
      >
        <MenuItem value={1}>1 Hour</MenuItem>
        <MenuItem value={2}>2 hours</MenuItem>
        <MenuItem value={3}>3 hours</MenuItem>
        <MenuItem value={4}>4 hours</MenuItem>
        <MenuItem value={5}>5 hours</MenuItem>
        <MenuItem value={6}>6 hours</MenuItem>
      </Select>
      <InputLabel id="guests">Guests</InputLabel>
      <Select
        labelId="guests"
        id="guests"
        value={bookingQuery.people}
        onChange={handlePeopleChange}
      >
        <MenuItem value={1}>1 guest</MenuItem>
        <MenuItem value={2}>2 guests</MenuItem>
        <MenuItem value={3}>3 guests</MenuItem>
        <MenuItem value={4}>4 guests</MenuItem>
        <MenuItem value={5}>5 guests</MenuItem>
        <MenuItem value={6}>6 guests</MenuItem>
      </Select>
    </form>
  );
}
