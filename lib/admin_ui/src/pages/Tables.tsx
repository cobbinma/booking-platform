import React, { useState } from "react";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { H2 } from "baseui/typography";
import { Table as BaseTable } from "baseui/table";
import { Table, useAddTableMutation, useRemoveTableMutation } from "../graph";
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

const Tables: React.FC<{ tables: Array<Table>; venueId: string }> = ({
  tables,
  venueId,
}) => {
  const [deleteIsOpen, setDeleteIsOpen] = React.useState<boolean>(false);
  const [addIsOpen, setAddIsOpen] = React.useState<boolean>(false);
  const [selectedTable, setSelectedTable] = React.useState<Table | null>(null);
  const [currentTables, setCurrentTables] = React.useState<Array<Table>>(
    tables
  );

  const removeTable = (id: string) => {
    setCurrentTables(currentTables.filter((table) => table.id !== id));
  };

  const addTable = (table: Table | null | undefined) => {
    if (table) setCurrentTables(currentTables.concat(table));
  };

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
          <BaseTable
            columns={["Name", "Capacity", ""]}
            data={currentTables.slice().map((table) => {
              return [
                table.name,
                table.capacity,
                <Button
                  onClick={() => {
                    setSelectedTable(table);
                    setDeleteIsOpen(true);
                  }}
                >
                  Delete
                </Button>,
              ];
            })}
          />
        </FlexGridItem>
        <DeleteTableModal
          deleteIsOpen={deleteIsOpen}
          setDeleteIsOpen={setDeleteIsOpen}
          selectedTable={selectedTable}
          removeTable={removeTable}
          venueId={venueId}
        />
        <AddTableModal
          addIsOpen={addIsOpen}
          setAddIsOpen={setAddIsOpen}
          venueId={venueId}
          addTable={addTable}
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
  removeTable: (id: string) => void;
  venueId: string;
}> = ({
  deleteIsOpen,
  setDeleteIsOpen,
  selectedTable,
  venueId,
  removeTable,
}) => {
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
          Deleting table will cancel all bookings associated with the booking.
        </ModalBody>
        <ModalFooter>
          <ModalButton kind="tertiary" onClick={() => setDeleteIsOpen(false)}>
            Cancel
          </ModalButton>
          <ModalButton
            onClick={() => {
              removeTableMutation()
                .then((table) => {
                  removeTable(table.data?.removeTable?.id || "");
                  setDeleteIsOpen(false);
                })
                .catch((e) => console.log(e));
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
  addTable: (table: Table | null | undefined) => void;
}> = ({ addIsOpen, setAddIsOpen, venueId, addTable }) => {
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
        <ModalButton
          onClick={() => {
            addTableMutation()
              .then((table) => {
                addTable(table?.data?.addTable);
                setName("");
                setCapacity("");
                setAddIsOpen(false);
              })
              .catch((e) => console.log(e));
          }}
        >
          Okay
        </ModalButton>
      </ModalFooter>
    </Modal>
  );
};
