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

export default function RequestSlotForm() {
  const classes = useStyles();
  const [email, setEmail] = React.useState<string>("");
  const [selectedDate, setSelectedDate] = React.useState<Date | null>(
    new Date()
  );
  const [selectedStartTime, setSelectedStartTime] = React.useState<Date | null>(
    new Date()
  );
  const [durationHours, setDurationHours] = React.useState<number>(0);
  const [people, setPeople] = React.useState<number>(0);

  const handleEmailChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    setEmail(event.target.value as string);
  };
  const handleDateChange = (date: Date | null) => {
    setSelectedDate(date);
  };

  const handleStartTimeChange = (date: Date | null) => {
    setSelectedStartTime(date);
  };

  const handleDurationChange = (
    event: React.ChangeEvent<{ value: unknown }>
  ) => {
    setDurationHours(event.target.value as number);
  };
  const handlePeopleChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    setPeople(event.target.value as number);
  };

  return (
    <form className={classes.root} noValidate autoComplete="off">
      <TextField
        id="email"
        label="Email"
        value={email}
        onChange={handleEmailChange}
      />
      <MuiPickersUtilsProvider utils={DateFnsUtils}>
        <KeyboardDatePicker
          margin="normal"
          id="date"
          label="Date"
          format="MM/dd/yyyy"
          value={selectedDate}
          onChange={handleDateChange}
          KeyboardButtonProps={{
            "aria-label": "change date",
          }}
        />
        <KeyboardTimePicker
          margin="normal"
          id="starts-at"
          label="Starting Time"
          value={selectedStartTime}
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
        value={durationHours}
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
        value={people}
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
