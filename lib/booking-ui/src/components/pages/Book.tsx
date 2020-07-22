import React from "react";
import RequestSlotForm from "../RequestSlotForm";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import { Step, Stepper, StepLabel, Button } from "@material-ui/core";
import Slot from "../Slot";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      width: "100%",
    },
    button: {
      marginRight: theme.spacing(1),
    },
    instructions: {
      marginTop: theme.spacing(1),
      marginBottom: theme.spacing(1),
    },
  })
);

function getSteps() {
  return ["Enquiry", "Confirm", "Done"];
}

export interface BookingQuery {
  customer_id: string;
  people: number;
  date: Date;
  starts_at: Date;
  duration: number;
}

export interface BookingSlot {
  customer_id: string;
  people: number;
  date: string;
  starts_at: string;
  ends_at: string;
  table_id: number;
}

const Book = () => {
  const classes = useStyles();
  const [bookingQuery, setBookingQuery] = React.useState<BookingQuery>({
    customer_id: "",
    people: 0,
    date: new Date(),
    starts_at: new Date(),
    duration: 0,
  });
  const [bookingSlot, setBookingSlot] = React.useState<BookingSlot | null>();
  const [activeStep, setActiveStep] = React.useState(0);
  const [skipped, setSkipped] = React.useState(new Set<number>());
  const steps = getSteps();

  const getStepContent = (step: number) => {
    switch (step) {
      case 0:
        return (
          <RequestSlotForm
            bookingQuery={bookingQuery}
            setBookingQuery={setBookingQuery}
          />
        );
      case 1:
        return (
          <Slot
            bookingQuery={bookingQuery}
            bookingSlot={bookingSlot}
            setBookingSlot={setBookingSlot}
          />
        );
      case 2:
        return "Your booking is confirmed!";
      default:
        return "Unknown step";
    }
  };

  const getNextButton = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Button
            variant="contained"
            color="primary"
            onClick={handleNext}
            className={classes.button}
          >
            Next
          </Button>
        );
      case 1:
        return (
          <Button
            variant="contained"
            color="primary"
            onClick={handleNext}
            className={classes.button}
          >
            Confirm
          </Button>
        );
      case 2:
        return;
      default:
        return;
    }
  };

  const isStepSkipped = (step: number) => {
    return skipped.has(step);
  };

  const handleNext = () => {
    let newSkipped = skipped;
    if (isStepSkipped(activeStep)) {
      newSkipped = new Set(newSkipped.values());
      newSkipped.delete(activeStep);
    }

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
    setSkipped(newSkipped);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  return (
    <div className={classes.root}>
      <h1>Book</h1>
      <Stepper activeStep={activeStep}>
        {steps.map((label, index) => {
          const stepProps: { completed?: boolean } = {};
          const labelProps: { optional?: React.ReactNode } = {};
          if (isStepSkipped(index)) {
            stepProps.completed = false;
          }
          return (
            <Step key={label} {...stepProps}>
              <StepLabel {...labelProps}>{label}</StepLabel>
            </Step>
          );
        })}
      </Stepper>
      {getStepContent(activeStep)}
      <Button
        disabled={activeStep === 0 || activeStep === 2}
        onClick={handleBack}
        className={classes.button}
      >
        Back
      </Button>
      {getNextButton(activeStep)}
    </div>
  );
};

export default Book;
