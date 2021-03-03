import React, { useState } from "react";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { H2 } from "baseui/typography";
import {
  GetVenueQuery,
  GetVenueQueryVariables,
  Table,
  useAddTableMutation,
  useRemoveTableMutation,
} from "../graph";
import { Button } from "baseui/button";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
} from "baseui/modal";
import { Input } from "baseui/input";
import { FormControl } from "baseui/form-control";
import { ApolloQueryResult } from "@apollo/client";
import { TableBuilder, TableBuilderColumn } from "baseui/table-semantic";

const Tables: React.FC<{
  tables: Array<Table>;
  venueId: string | null | undefined;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ tables, venueId, refetch }) => {
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

  const [deleteIsOpen, setDeleteIsOpen] = React.useState<boolean>(false);
  const [addIsOpen, setAddIsOpen] = React.useState<boolean>(false);
  const [selectedTable, setSelectedTable] = React.useState<Table | null>(null);

  if (!venueId) return <div>error</div>;

  return (
    <div>
      <FlexGrid>
        <FlexGridItem>
          <H2>Tables</H2>
        </FlexGridItem>
        <FlexGridItem>
          <Button
            onClick={() => {
              setAddIsOpen(true);
            }}
          >
            Add Table
          </Button>
        </FlexGridItem>
        <FlexGridItem>
          <TableBuilder data={tables} overrides={overrides}>
            <TableBuilderColumn header="Name">
              {(row) => row.name}
            </TableBuilderColumn>
            <TableBuilderColumn header="Capacity">
              {(row) => row.capacity}
            </TableBuilderColumn>
            <TableBuilderColumn>
              {(row) => (
                <Button
                  onClick={() => {
                    setSelectedTable(row);
                    setDeleteIsOpen(true);
                  }}
                >
                  Delete
                </Button>
              )}
            </TableBuilderColumn>
          </TableBuilder>
        </FlexGridItem>
        <DeleteTableModal
          deleteIsOpen={deleteIsOpen}
          setDeleteIsOpen={setDeleteIsOpen}
          selectedTable={selectedTable}
          refetch={refetch}
          venueId={venueId}
        />
        <AddTableModal
          addIsOpen={addIsOpen}
          setAddIsOpen={setAddIsOpen}
          venueId={venueId}
          refetch={refetch}
        />
      </FlexGrid>
    </div>
  );
};

export default Tables;

const DeleteTableModal: React.FC<{
  deleteIsOpen: boolean;
  setDeleteIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  selectedTable: Table | null;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
  venueId: string;
}> = ({ deleteIsOpen, setDeleteIsOpen, selectedTable, venueId, refetch }) => {
  const [removeTableMutation] = useRemoveTableMutation({
    variables: {
      table: { venueId: venueId, tableId: selectedTable?.id || "" },
    },
  });

  return (
    <div>
      <Modal onClose={() => setDeleteIsOpen(false)} isOpen={deleteIsOpen}>
        <ModalHeader>Delete {selectedTable?.name}</ModalHeader>
        <ModalBody>
          Deleting table will cancel all bookings associated with the table.
        </ModalBody>
        <ModalFooter>
          <ModalButton kind="tertiary" onClick={() => setDeleteIsOpen(false)}>
            Cancel
          </ModalButton>
          <ModalButton
            onClick={() => {
              removeTableMutation()
                .then(() => {
                  refetch().catch((e) => console.log(e));
                })
                .catch((e) => console.log(e));
              setDeleteIsOpen(false);
            }}
          >
            Okay
          </ModalButton>
        </ModalFooter>
      </Modal>
    </div>
  );
};

const AddTableModal: React.FC<{
  addIsOpen: boolean;
  setAddIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  venueId: string;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ addIsOpen, setAddIsOpen, venueId, refetch }) => {
  const [name, setName] = useState<string>("");
  const [capacity, setCapacity] = useState<string>("");
  const [addTableMutation] = useAddTableMutation({
    variables: {
      table: { venueId: venueId, name: name, capacity: parseInt(capacity) },
    },
  });
  const close = () => {
    setName("");
    setCapacity("");
    setAddIsOpen(false);
  };

  return (
    <Modal onClose={close} isOpen={addIsOpen}>
      <ModalHeader>Add Table</ModalHeader>
      <ModalBody>
        <FormControl label={() => "Table Name"}>
          <Input
            value={name}
            onChange={(event) => setName(event.currentTarget.value)}
            placeholder="Name"
          />
        </FormControl>
        <FormControl label={() => "Table Capacity"}>
          <Input
            value={capacity}
            onChange={(event) => {
              setCapacity(event.currentTarget.value);
            }}
            placeholder="Capacity"
            error={isNaN(Number(capacity))}
          />
        </FormControl>
      </ModalBody>
      <ModalFooter>
        <ModalButton onClick={close} kind="tertiary">
          Cancel
        </ModalButton>
        {!isNaN(Number(capacity)) && parseInt(capacity) > 0 && name !== "" ? (
          <ModalButton
            onClick={() => {
              addTableMutation()
                .then(() => {
                  refetch().catch((e) => console.log(e));
                  setName("");
                  setCapacity("");
                  setAddIsOpen(false);
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
