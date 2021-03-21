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
  Booking,
  GetVenueQuery,
  GetVenueQueryVariables,
  OpeningHoursSpecification,
  Table,
  useCancelBookingMutation,
  useCreateBookingMutation,
} from "../graph";
import { ApolloQueryResult } from "@apollo/client";
import { DatePicker } from "baseui/datepicker";
import { TableBuilder, TableBuilderColumn } from "baseui/table-semantic";
import { Pagination } from "baseui/pagination";

let durations = new Map<string, number>([
  ["30 mins", 30],
  ["1 hour", 60],
  ["1 hour 30 mins", 90],
  ["2 hours", 120],
  ["2 hours 30 mins", 150],
  ["3 hours", 180],
]);

const PAGE_LIMIT = 20;

const Bookings: React.FC<{
  tables: Array<Table>;
  bookings: Array<Booking>;
  pages: number;
  venueId: string | null | undefined;
  openHours: OpeningHoursSpecification | null | undefined;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ tables, bookings, pages, venueId, openHours, refetch }) => {
  const overrides = {
    TableBodyRow: {
      style: ({ $theme, $rowIndex }: any) => ({
        backgroundColor:
          $rowIndex % 2
            ? $theme.colors.backgroundPrimary
            : $theme.colors.backgroundSecondary,
        ":hover": {
          backgroundColor: $theme.colors.backgroundTertiary,
        },
      }),
    },
  };

  const [selectedBooking, setSelectedBooking] = React.useState<Booking | null>(
    null
  );
  const [createIsOpen, setCreateIsOpen] = React.useState<boolean>(false);
  const [cancelIsOpen, setCancelIsOpen] = React.useState<boolean>(false);
  const [date, setDate] = React.useState([new Date(Date.now())]);
  const [currentPage, setCurrentPage] = React.useState(1);

  if (!venueId) return <div>error</div>;

  const tableIdMap = tables.reduce((map, table) => {
    map.set(table.id, table.name);
    return map;
  }, new Map());

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
            openHours={openHours}
            refetch={refetch}
          />
        </FlexGridItem>
        <FlexGridItem>
          <DatePicker
            value={date}
            onChange={({ date }) => {
              refetch({
                venueID: venueId,
                filter: {
                  date: Array.isArray(date)
                    ? date[0].toISOString()
                    : date.toISOString(),
                },
                pageInfo: { page: 0, limit: PAGE_LIMIT },
              }).catch((e) => console.log(e));
              setDate(Array.isArray(date) ? date : [date]);
              setCurrentPage(1);
            }}
          />
        </FlexGridItem>
        <FlexGridItem>
          <Pagination
            numPages={pages}
            currentPage={currentPage}
            onPageChange={({ nextPage }) => {
              refetch({
                venueID: venueId,
                filter: {
                  date: date[0].toISOString(),
                },
                pageInfo: { page: nextPage - 1, limit: PAGE_LIMIT },
              }).catch((e) => console.log(e));
              setCurrentPage(Math.min(Math.max(nextPage, 1), pages));
            }}
          />
        </FlexGridItem>
        <FlexGridItem>
          <TableBuilder data={bookings} overrides={overrides}>
            <TableBuilderColumn header="Email">
              {(row) => row.email}
            </TableBuilderColumn>
            <TableBuilderColumn header="People">
              {(row) => row.people}
            </TableBuilderColumn>
            <TableBuilderColumn header="Starts At">
              {(row) => new Date(Date.parse(row.startsAt)).toLocaleString()}
            </TableBuilderColumn>
            <TableBuilderColumn header="Ends At">
              {(row) => new Date(Date.parse(row.endsAt)).toLocaleTimeString()}
            </TableBuilderColumn>
            <TableBuilderColumn header="Duration">
              {(row) => row.duration}
            </TableBuilderColumn>
            <TableBuilderColumn header="Table">
              {(row) => tableIdMap.get(row.tableId)}
            </TableBuilderColumn>
            <TableBuilderColumn>
              {(row) => (
                <Button
                  onClick={() => {
                    setSelectedBooking(row);
                    setCancelIsOpen(true);
                  }}
                >
                  Cancel
                </Button>
              )}
            </TableBuilderColumn>
          </TableBuilder>
        </FlexGridItem>
      </FlexGrid>
      <CancelBooking
        cancelIsOpen={cancelIsOpen}
        selectedBooking={selectedBooking}
        setCancelIsOpen={setCancelIsOpen}
        venueId={venueId}
        refetch={refetch}
      />
    </div>
  );
};

export default Bookings;

const CreateBooking: React.FC<{
  setCreateIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  createIsOpen: boolean;
  venueId: string;
  openHours: OpeningHoursSpecification | null | undefined;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ setCreateIsOpen, createIsOpen, venueId, openHours, refetch }) => {
  const [email, setEmail] = useState<string>("");
  const [people, setPeople] = useState<number[]>([4]);
  const [date, setDate] = React.useState<Date[] | null>(null);
  const [time, setTime] = React.useState<Date>(new Date(Date.now()));
  const [duration, setDuration] = React.useState<string | null>(null);
  const close = (): void => {
    setCreateIsOpen(false);
    setEmail("");
  };

  const isTimeOfDayBeforeDate = (
    timeOfDay: string,
    date: Date,
    addMinutes?: number
  ): boolean => {
    const splitTime = timeOfDay.split(":");
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

  console.log(openHours);

  const [createBookingMutation, { loading, error }] = useCreateBookingMutation({
    variables: {
      slot: {
        venueId: venueId,
        email: email,
        people: people[0],
        startsAt:
          date && date[0]
            ? new Date(
                date[0].getFullYear(),
                date[0].getMonth(),
                date[0].getDate(),
                time.getHours(),
                time.getMinutes()
              )
            : undefined,
        duration: durations.get(duration || "") || 60,
      },
    },
  });

  return (
    <Modal onClose={close} closeable isOpen={createIsOpen} animate autoFocus>
      <ModalHeader>Create Booking</ModalHeader>
      {loading && (
        <ModalHeader>
          <StyledLoadingSpinner />
        </ModalHeader>
      )}
      {error && <ModalHeader>could not create booking</ModalHeader>}
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
        <FormControl
          label="Date"
          caption={() =>
            (date && !openHours) ||
            (openHours && !openHours.opens) ||
            (openHours && !openHours.closes)
              ? "venue is closed"
              : ""
          }
        >
          <DatePicker
            value={date}
            onChange={({ date }) => {
              const d = Array.isArray(date) ? date : [date];
              setDate(d);
              if (d && d[0]) {
                console.log(d);
                refetch({
                  date: d[0].toISOString(),
                }).catch((e) => console.log(e));
              }
            }}
            minDate={new Date()}
            error={
              !!(
                (date && !openHours) ||
                (openHours && !openHours.opens) ||
                (openHours && !openHours.closes)
              )
            }
          />
        </FormControl>
        {openHours && (
          <div>
            <FormControl
              label="Start Time"
              caption={() =>
                isTimeOfDayBeforeDate(openHours.opens, time) ||
                !isTimeOfDayBeforeDate(openHours.closes, time)
                  ? "venue is closed"
                  : ""
              }
            >
              <TimePicker
                value={time}
                step={1800}
                onChange={(date) => {
                  setDuration(null);
                  setTime(date);
                }}
                disabled={!openHours.opens}
                error={
                  isTimeOfDayBeforeDate(openHours.opens, time) ||
                  !isTimeOfDayBeforeDate(openHours.closes, time)
                }
              />
            </FormControl>
            <FormControl label="Duration">
              <Combobox
                value={duration || ""}
                onChange={(nextValue) => setDuration(nextValue)}
                options={Array.from(durations.keys()).filter((k) => {
                  const duration = durations.get(k);
                  return (
                    duration &&
                    isTimeOfDayBeforeDate(openHours.closes, time, duration)
                  );
                })}
                mapOptionToString={(option) => option}
                disabled={
                  isTimeOfDayBeforeDate(openHours.opens, time) ||
                  !isTimeOfDayBeforeDate(openHours.closes, time) ||
                  !openHours.closes
                }
              />
            </FormControl>
          </div>
        )}
      </ModalBody>
      <ModalFooter>
        <ModalButton kind="tertiary" onClick={close}>
          Cancel
        </ModalButton>
        {isEmailValid(email) &&
        durations.get(duration || "") &&
        openHours &&
        openHours.opens &&
        openHours.closes &&
        !isTimeOfDayBeforeDate(openHours.opens, time) &&
        isTimeOfDayBeforeDate(openHours.closes, time) ? (
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

const CancelBooking: React.FC<{
  cancelIsOpen: boolean;
  setCancelIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  selectedBooking: Booking | null;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
  venueId: string;
}> = ({ cancelIsOpen, setCancelIsOpen, selectedBooking, venueId, refetch }) => {
  const [cancelBookingMutation] = useCancelBookingMutation({
    variables: {
      input: { venueId: venueId, id: selectedBooking?.id || "" },
    },
  });

  return (
    <div>
      <Modal onClose={() => setCancelIsOpen(false)} isOpen={cancelIsOpen}>
        <ModalHeader>Cancel Booking</ModalHeader>
        <ModalBody>
          Cancel booking for {selectedBooking?.email} at{" "}
          {new Date(Date.parse(selectedBooking?.startsAt)).toLocaleString()}
        </ModalBody>
        <ModalFooter>
          <ModalButton kind="tertiary" onClick={() => setCancelIsOpen(false)}>
            Cancel
          </ModalButton>
          <ModalButton
            onClick={() => {
              cancelBookingMutation()
                .then(() => {
                  refetch().catch((e) => console.log(e));
                })
                .catch((e) => console.log(e));
              setCancelIsOpen(false);
            }}
          >
            Okay
          </ModalButton>
        </ModalFooter>
      </Modal>
    </div>
  );
};
