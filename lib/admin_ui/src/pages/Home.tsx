import React, { ChangeEvent, useState } from "react";
import { H2, H5 } from "baseui/typography";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { Input } from "baseui/input";
import { FormControl } from "baseui/form-control";
import {
  GetVenueQuery,
  GetVenueQueryVariables,
  OpeningHoursSpecification,
  useUpdateOpeningHoursMutation,
  useUpdateSpecialOpeningHoursMutation,
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
import { Button } from "baseui/button";
import { Combobox } from "baseui/combobox";
import { TimePicker } from "baseui/timepicker";
import { Checkbox } from "baseui/checkbox";
import { DatePicker } from "baseui/datepicker";

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
  const [
    addSpecialOpeningHoursIsOpen,
    setAddSpecialOpeningHoursIsOpen,
  ] = useState<boolean>(false);

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
            <TableBuilderColumn>
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
          <Button onClick={() => setAddSpecialOpeningHoursIsOpen(true)}>
            Add Special Opening Hours
          </Button>
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
            <TableBuilderColumn header="Valid Through">
              {(row) => row.validThrough}
            </TableBuilderColumn>
            <TableBuilderColumn>
              {(row) => (
                <EditSpecialOpeningHours
                  venueId={venueId}
                  openingHours={specialOpeningHours.filter(
                    (o) => o.dayOfWeek !== row.dayOfWeek
                  )}
                  text="Remove"
                  refetch={refetch}
                />
              )}
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
      <AddSpecialOpeningHours
        venueId={venueId}
        openingHours={specialOpeningHours}
        setAddSpecialOpeningHoursIsOpen={setAddSpecialOpeningHoursIsOpen}
        addSpecialOpeningHoursIsOpen={addSpecialOpeningHoursIsOpen}
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

const AddSpecialOpeningHours: React.FC<{
  venueId: string;
  openingHours: OpeningHoursSpecification[];
  setAddSpecialOpeningHoursIsOpen: React.Dispatch<
    React.SetStateAction<boolean>
  >;
  addSpecialOpeningHoursIsOpen: boolean;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({
  openingHours,
  venueId,
  setAddSpecialOpeningHoursIsOpen,
  addSpecialOpeningHoursIsOpen,
  refetch,
}) => {
  const [day, setDay] = useState<number | null>(null);

  const [open, setOpen] = useState<boolean>(true);

  const [opens, setOpens] = useState<Date | null>(null);
  const [closes, setCloses] = useState<Date | null>(null);

  const [validFrom, setValidFrom] = React.useState<Date[] | null>(null);
  const [validThrough, setValidThrough] = React.useState<Date[] | null>(null);

  const close = () => {
    setDay(null);
    setOpens(null);
    setCloses(null);
    setValidFrom(null);
    setValidThrough(null);
    setAddSpecialOpeningHoursIsOpen(false);
  };

  const append_hours = (): OpeningHoursSpecification[] => [
    ...openingHours,
    {
      dayOfWeek: day,
      opens: opens
        ? ("0" + opens.getHours()).slice(-2) +
          ":" +
          ("0" + opens.getMinutes()).slice(-2)
        : null,
      closes: closes
        ? ("0" + closes.getHours()).slice(-2) +
          ":" +
          ("0" + closes.getMinutes()).slice(-2)
        : null,
      validFrom: validFrom && validFrom[0] ? validFrom[0].toISOString() : "",
      validThrough:
        validThrough && validThrough[0] ? validThrough[0].toISOString() : "",
    },
  ];

  return (
    <div>
      <Modal onClose={close} isOpen={addSpecialOpeningHoursIsOpen}>
        <ModalHeader>Add Special Opening Hours</ModalHeader>
        <ModalBody>
          <FormControl label="Day">
            <Combobox
              value={weekday.get(day || 0) || ""}
              onChange={(nextValue) => setDay(days.get(nextValue) || null)}
              options={Array.from(days.keys())}
              mapOptionToString={(option) => option}
            />
          </FormControl>
          <FormControl label="Open">
            <Checkbox
              checked={open}
              onChange={(e: ChangeEvent<HTMLInputElement>) => {
                setOpens(null);
                setCloses(null);
                setOpen(e.target.checked);
              }}
            />
          </FormControl>
          <FormControl label="Opens">
            <TimePicker
              value={opens}
              step={1800}
              onChange={(date) => {
                setOpens(date);
              }}
              disabled={!day || !open}
            />
          </FormControl>
          <FormControl label="Closes">
            <TimePicker
              value={closes}
              step={1800}
              onChange={(date) => {
                setCloses(date);
              }}
              disabled={!day || !open}
            />
          </FormControl>
          <FormControl label="Valid From">
            <DatePicker
              value={validFrom}
              onChange={({ date }) => {
                setValidFrom(Array.isArray(date) ? date : [date]);
              }}
            />
          </FormControl>
          <FormControl label="Valid Through">
            <DatePicker
              value={validThrough}
              onChange={({ date }) => {
                setValidThrough(Array.isArray(date) ? date : [date]);
              }}
            />
          </FormControl>
        </ModalBody>
        <ModalFooter>
          <ModalButton kind="tertiary" onClick={close}>
            Cancel
          </ModalButton>
          {day && validFrom && validThrough && (
            <EditSpecialOpeningHours
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

const EditSpecialOpeningHours: React.FC<{
  venueId: string;
  openingHours: OpeningHoursSpecification[];
  text: string;
  close?: () => void;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ venueId, openingHours, text, refetch, close }) => {
  const [
    updateSpecialOpeningHoursMutation,
  ] = useUpdateSpecialOpeningHoursMutation({
    variables: {
      input: {
        venueId: venueId || "",
        specialOpeningHours: openingHours.map((o) => {
          return {
            dayOfWeek: o.dayOfWeek,
            opens: o.opens,
            closes: o.closes,
            validFrom: o.validFrom || "",
            validThrough: o.validThrough || "",
          };
        }),
      },
    },
  });

  return (
    <Button
      onClick={() => {
        updateSpecialOpeningHoursMutation()
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
