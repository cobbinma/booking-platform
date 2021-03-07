import React, { useState } from "react";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { H2 } from "baseui/typography";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
} from "baseui/modal";
import { Button, StyledLoadingSpinner } from "baseui/button";
import { FormControl } from "baseui/form-control";
import { Input } from "baseui/input";
import { isEmailValid } from "./Admins";
import { Slider } from "baseui/slider";
import { TimePicker } from "baseui/timepicker";
import { Combobox } from "baseui/combobox";
import {
  GetVenueQuery,
  GetVenueQueryVariables,
  useCreateBookingMutation,
} from "../graph";
import { ApolloQueryResult } from "@apollo/client";
import { DatePicker } from "baseui/datepicker";

let durations = new Map<string, number>([
  ["30 mins", 30],
  ["1 hour", 60],
  ["1 hour 30 mins", 90],
  ["2 hours", 120],
  ["2 hours 30 mins", 150],
  ["3 hours", 180],
]);

const Bookings: React.FC<{
  venueId: string | null | undefined;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ venueId, refetch }) => {
  const [createIsOpen, setCreateIsOpen] = React.useState<boolean>(false);

  if (!venueId) return <div>error</div>;

  return (
    <div>
      <FlexGrid>
        <FlexGridItem>
          <H2>Bookings</H2>
        </FlexGridItem>
        <FlexGridItem>
          <Button onClick={() => setCreateIsOpen(true)}>Create Booking</Button>
          <CreateBooking
            setCreateIsOpen={setCreateIsOpen}
            createIsOpen={createIsOpen}
            venueId={venueId}
            refetch={refetch}
          />
        </FlexGridItem>
      </FlexGrid>
    </div>
  );
};

export default Bookings;

const CreateBooking: React.FC<{
  setCreateIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  createIsOpen: boolean;
  venueId: string;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ setCreateIsOpen, createIsOpen, venueId, refetch }) => {
  const [email, setEmail] = useState<string>("");
  const [people, setPeople] = useState<number[]>([4]);
  const [date, setDate] = React.useState([new Date(Date.now())]);
  const [time, setTime] = React.useState(new Date(Date.now()));
  const [duration, setDuration] = React.useState<string>("1 hour");
  const close = (): void => {
    setCreateIsOpen(false);
    setEmail("");
    error = undefined;
  };

  let [createBookingMutation, { loading, error }] = useCreateBookingMutation({
    variables: {
      slot: {
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
      },
    },
  });

  if (loading)
    return (
      <Modal onClose={close} closeable isOpen={createIsOpen} animate autoFocus>
        <StyledLoadingSpinner />
      </Modal>
    );

  if (error)
    return (
      <Modal onClose={close} closeable isOpen={createIsOpen} animate autoFocus>
        <H2>could not create booking</H2>
      </Modal>
    );

  return (
    <Modal onClose={close} closeable isOpen={createIsOpen} animate autoFocus>
      <ModalHeader>Create Booking</ModalHeader>
      <ModalBody>
        <FormControl label="Email">
          <Input
            value={email}
            onChange={(event) => setEmail(event.currentTarget.value)}
            placeholder="Email"
            error={email !== "" && !isEmailValid(email)}
          />
        </FormControl>
        <FormControl label="People">
          <Slider
            value={people}
            onChange={({ value }) => value && setPeople(value)}
            min={1}
            max={20}
          />
        </FormControl>
        <FormControl label="Date">
          <DatePicker
            value={date}
            onChange={({ date }) =>
              setDate(Array.isArray(date) ? date : [date])
            }
          />
        </FormControl>
        <FormControl label="Start Time">
          <TimePicker
            value={time}
            step={1800}
            onChange={(date) => setTime(date)}
          />
        </FormControl>
        <FormControl label="Duration">
          <Combobox
            value={duration}
            onChange={(nextValue) => setDuration(nextValue)}
            options={Array.from(durations.keys())}
            mapOptionToString={(option) => option}
          />
        </FormControl>
      </ModalBody>
      <ModalFooter>
        <ModalButton kind="tertiary" onClick={close}>
          Cancel
        </ModalButton>
        {isEmailValid(email) && durations.get(duration) ? (
          <ModalButton
            onClick={() => {
              createBookingMutation()
                .then(() => {
                  refetch().catch((e) => console.log(e));
                  close();
                })
                .catch((e) => console.log(e));
            }}
          >
            Okay
          </ModalButton>
        ) : null}
      </ModalFooter>
    </Modal>
  );
};
