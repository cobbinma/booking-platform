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
  Table,
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
  slug: string;
  pages: number;
  venueId: string | null | undefined;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ tables, bookings, slug, pages, venueId, refetch }) => {
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

  const [createIsOpen, setCreateIsOpen] = React.useState<boolean>(false);
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
            refetch={refetch}
          />
        </FlexGridItem>
        <FlexGridItem>
          <DatePicker
            value={date}
            onChange={({ date }) => {
              refetch({
                slug: slug,
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
                slug: slug,
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
              {() => <Button>Cancel</Button>}
            </TableBuilderColumn>
          </TableBuilder>
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
  };

  const [createBookingMutation, { loading, error }] = useCreateBookingMutation({
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
