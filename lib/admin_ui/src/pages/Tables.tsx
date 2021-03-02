import React from "react";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { H2 } from "baseui/typography";
import { Table as BaseTable } from "baseui/table";
import { Table, useRemoveTableMutation } from "../graph";
import { Button } from "baseui/button";
import {
  Modal,
  ModalBody,
  ModalButton,
  ModalFooter,
  ModalHeader,
} from "baseui/modal";

const Tables: React.FC<{ tables: Array<Table>; venueId: string }> = ({
  tables,
  venueId,
}) => {
  const [deleteIsOpen, setDeleteIsOpen] = React.useState<boolean>(false);
  const [selectedTable, setSelectedTable] = React.useState<Table | null>(null);
  const [currentTables, setCurrentTables] = React.useState<Array<Table>>(
    tables
  );

  const removeTable = (id: string) => {
    setCurrentTables(currentTables.filter((table) => table.id !== id));
  };

  return (
    <div>
      <FlexGrid>
        <FlexGridItem>
          <H2>Tables</H2>
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
