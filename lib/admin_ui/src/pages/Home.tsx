import React, { useState } from "react";
import { H2, H5 } from "baseui/typography";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { Input } from "baseui/input";
import { FormControl } from "baseui/form-control";
import {
  GetVenueQuery,
  GetVenueQueryVariables,
  OpeningHoursSpecification,
  useUpdateOpeningHoursMutation,
} from "../graph";
import { ApolloQueryResult } from "@apollo/client";
import { TableBuilder, TableBuilderColumn } from "baseui/table-semantic";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
} from "baseui/modal";
import { Checkbox } from "baseui/checkbox";
import { Button } from "baseui/button";
import { Combobox } from "baseui/combobox";
import { TimePicker } from "baseui/timepicker";

var weekday = new Map<number, string>([
  [1, "Monday"],
  [2, "Tuesday"],
  [3, "Wednesday"],
  [4, "Thursday"],
  [5, "Friday"],
  [6, "Saturday"],
  [7, "Sunday"],
]);

const days = new Map<string, number>([
  ["Monday", 1],
  ["Tuesday", 2],
  ["Wednesday", 3],
  ["Thursday", 4],
  ["Friday", 5],
  ["Saturday", 6],
  ["Sunday", 7],
]);

const Home: React.FC<{
  name: string | null | undefined;
  slug: string | null | undefined;
  openingHours: OpeningHoursSpecification[];
  specialOpeningHours: OpeningHoursSpecification[];
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
  venueId: string | undefined;
}> = ({ name, slug, openingHours, specialOpeningHours, refetch, venueId }) => {
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

  const [addOpeningHoursIsOpen, setAddOpeningHoursIsOpen] = useState<boolean>(
    false
  );

  if (!venueId) return <div>error</div>;

  return (
    <div>
      <FlexGrid>
        <FlexGridItem>
          <H2>Home</H2>
        </FlexGridItem>
        <FlexGridItem>
          <FormControl label="Name">
            <Input value={name || ""} />
          </FormControl>
        </FlexGridItem>
        <FlexGridItem>
          <FormControl label="Slug">
            <Input value={slug || ""} />
          </FormControl>
        </FlexGridItem>
        <FlexGridItem>
          <H5>Opening Hours</H5>
        </FlexGridItem>
        {openingHours.length < 7 && (
          <FlexGridItem>
            <Button onClick={() => setAddOpeningHoursIsOpen(true)}>
              Add Opening Hours
            </Button>
          </FlexGridItem>
        )}
        <FlexGridItem>
          <TableBuilder data={openingHours} overrides={overrides}>
            <TableBuilderColumn header="Day">
              {(row) => weekday.get(row.dayOfWeek)}
            </TableBuilderColumn>
            <TableBuilderColumn header="Opens">
              {(row) => row.opens}
            </TableBuilderColumn>
            <TableBuilderColumn header="Closes">
              {(row) => row.closes}
            </TableBuilderColumn>
            <TableBuilderColumn header="Remove">
              {(row) => (
                <EditOpeningHours
                  venueId={venueId}
                  openingHours={openingHours.filter(
                    (o) => o.dayOfWeek !== row.dayOfWeek
                  )}
                  text="Remove"
                  refetch={refetch}
                />
              )}
            </TableBuilderColumn>
          </TableBuilder>
        </FlexGridItem>
        <FlexGridItem>
          <H5>Special Opening Hours</H5>
        </FlexGridItem>
        <FlexGridItem>
          <TableBuilder data={specialOpeningHours} overrides={overrides}>
            <TableBuilderColumn header="Day">
              {(row) => weekday.get(row.dayOfWeek)}
            </TableBuilderColumn>
            <TableBuilderColumn header="Opens">
              {(row) => row.opens}
            </TableBuilderColumn>
            <TableBuilderColumn header="Closes">
              {(row) => row.closes}
            </TableBuilderColumn>
            <TableBuilderColumn header="Valid From">
              {(row) => row.validFrom}
            </TableBuilderColumn>
            <TableBuilderColumn header="valid Through">
              {(row) => row.validThrough}
            </TableBuilderColumn>
          </TableBuilder>
        </FlexGridItem>
      </FlexGrid>
      <AddOpeningHours
        venueId={venueId}
        openingHours={openingHours}
        setAddOpeningHoursIsOpen={setAddOpeningHoursIsOpen}
        addOpeningHoursIsOpen={addOpeningHoursIsOpen}
        refetch={refetch}
      />
    </div>
  );
};

export default Home;

const AddOpeningHours: React.FC<{
  venueId: string;
  openingHours: OpeningHoursSpecification[];
  setAddOpeningHoursIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  addOpeningHoursIsOpen: boolean;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({
  openingHours,
  venueId,
  setAddOpeningHoursIsOpen,
  addOpeningHoursIsOpen,
  refetch,
}) => {
  const [day, setDay] = useState<number | null>(null);

  const [opens, setOpens] = useState<Date | null>(null);
  const [closes, setCloses] = useState<Date | null>(null);

  const close = () => {
    setDay(null);
    setOpens(null);
    setCloses(null);
    setAddOpeningHoursIsOpen(false);
  };

  const append_hours = (): OpeningHoursSpecification[] => [
    ...openingHours,
    {
      dayOfWeek: day,
      opens: opens
        ? ("0" + opens.getHours()).slice(-2) +
          ":" +
          ("0" + opens.getMinutes()).slice(-2)
        : "",
      closes: closes
        ? ("0" + closes.getHours()).slice(-2) +
          ":" +
          ("0" + closes.getMinutes()).slice(-2)
        : "",
    },
  ];

  return (
    <div>
      <Modal onClose={close} isOpen={addOpeningHoursIsOpen}>
        <ModalHeader>Add Opening Hours</ModalHeader>
        <ModalBody>
          <FormControl label="Day">
            <Combobox
              value={weekday.get(day || 0) || ""}
              onChange={(nextValue) => setDay(days.get(nextValue) || null)}
              options={Array.from(days.keys()).filter(
                (d) => !openingHours.find((o) => o.dayOfWeek === days.get(d))
              )}
              mapOptionToString={(option) => option}
            />
          </FormControl>
          <FormControl label="Opens">
            <TimePicker
              value={opens}
              step={1800}
              onChange={(date) => {
                setOpens(date);
              }}
              disabled={!day}
            />
          </FormControl>
          <FormControl label="Closes">
            <TimePicker
              value={closes}
              step={1800}
              onChange={(date) => {
                setCloses(date);
              }}
              disabled={!day}
            />
          </FormControl>
        </ModalBody>
        <ModalFooter>
          <ModalButton kind="tertiary" onClick={close}>
            Cancel
          </ModalButton>
          {day && opens && closes && (
            <EditOpeningHours
              venueId={venueId}
              openingHours={append_hours()}
              text="Add"
              close={close}
              refetch={refetch}
            />
          )}
        </ModalFooter>
      </Modal>
    </div>
  );
};

const EditOpeningHours: React.FC<{
  venueId: string;
  openingHours: OpeningHoursSpecification[];
  text: string;
  close?: () => void;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ venueId, openingHours, text, refetch, close }) => {
  const openHours = openingHours.map((o) => {
    return {
      dayOfWeek: o.dayOfWeek,
      opens: o.opens || "",
      closes: o.closes || "",
    };
  });

  const [updateOpeningHoursMutation] = useUpdateOpeningHoursMutation({
    variables: {
      input: {
        venueId: venueId || "",
        openingHours: openHours,
      },
    },
  });

  return (
    <Button
      onClick={() => {
        updateOpeningHoursMutation()
          .then(() =>
            refetch()
              .then(() => {
                if (close) close();
              })
              .catch((e) => console.log(e))
          )
          .catch((e) => console.log(e));
      }}
    >
      {text}
    </Button>
  );
};
