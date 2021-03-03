import React, { useState } from "react";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { H2 } from "baseui/typography";
import { Table as BaseTable } from "baseui/table";
import {
  GetVenueQuery,
  GetVenueQueryVariables,
  useAddAdminMutation,
  useRemoveAdminMutation,
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
import { isBoolean } from "util";

const Admins: React.FC<{
  admins: Array<string>;
  venueId: string | null | undefined;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ admins, venueId, refetch }) => {
  const [deleteIsOpen, setDeleteIsOpen] = React.useState<boolean>(false);
  const [addIsOpen, setAddIsOpen] = React.useState<boolean>(false);
  const [selectedAdmin, setSelectedAdmin] = React.useState<string | null>(null);

  if (!venueId) return <div>error</div>;

  return (
    <div>
      <FlexGrid>
        <FlexGridItem>
          <H2>Administrators</H2>
        </FlexGridItem>
        <FlexGridItem>
          <Button
            onClick={() => {
              setAddIsOpen(true);
            }}
          >
            Add Administrator
          </Button>
        </FlexGridItem>
        <FlexGridItem>
          <BaseTable
            columns={["Name", ""]}
            data={admins.slice().map((admin) => {
              return [
                admin,
                <Button
                  onClick={() => {
                    setSelectedAdmin(admin);
                    setDeleteIsOpen(true);
                  }}
                >
                  Delete
                </Button>,
              ];
            })}
          />
        </FlexGridItem>
        <DeleteAdminModal
          deleteIsOpen={deleteIsOpen}
          setDeleteIsOpen={setDeleteIsOpen}
          selectedAdmin={selectedAdmin}
          refetch={refetch}
          venueId={venueId}
        />
        <AddAdminModal
          addIsOpen={addIsOpen}
          setAddIsOpen={setAddIsOpen}
          venueId={venueId}
          refetch={refetch}
        />
      </FlexGrid>
    </div>
  );
};

export default Admins;

const DeleteAdminModal: React.FC<{
  deleteIsOpen: boolean;
  setDeleteIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  selectedAdmin: string | null;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
  venueId: string;
}> = ({ deleteIsOpen, setDeleteIsOpen, selectedAdmin, venueId, refetch }) => {
  const [removeAdminMutation] = useRemoveAdminMutation({
    variables: {
      admin: { venueId: venueId, email: selectedAdmin || "" },
    },
  });

  return (
    <div>
      <Modal onClose={() => setDeleteIsOpen(false)} isOpen={deleteIsOpen}>
        <ModalHeader>Delete {selectedAdmin}</ModalHeader>
        <ModalBody>Deleting administrator will remove their access.</ModalBody>
        <ModalFooter>
          <ModalButton kind="tertiary" onClick={() => setDeleteIsOpen(false)}>
            Cancel
          </ModalButton>
          <ModalButton
            onClick={() => {
              removeAdminMutation()
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

const AddAdminModal: React.FC<{
  addIsOpen: boolean;
  setAddIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
  venueId: string;
  refetch: (
    variables?: GetVenueQueryVariables
  ) => Promise<ApolloQueryResult<GetVenueQuery>>;
}> = ({ addIsOpen, setAddIsOpen, venueId, refetch }) => {
  const [email, setEmail] = useState<string>("");
  const [addAdminMutation] = useAddAdminMutation({
    variables: {
      admin: { venueId: venueId, email: email },
    },
  });
  const close = () => {
    setEmail("");
    setAddIsOpen(false);
  };
  function isEmailValid(email: string): boolean {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  }

  return (
    <Modal onClose={close} isOpen={addIsOpen}>
      <ModalHeader>Add Admin</ModalHeader>
      <ModalBody>
        <FormControl label={() => "Email"}>
          <Input
            value={email}
            onChange={(event) => setEmail(event.currentTarget.value)}
            placeholder="Email"
            error={email !== "" && !isEmailValid(email)}
          />
        </FormControl>
      </ModalBody>
      <ModalFooter>
        <ModalButton onClick={close} kind="tertiary">
          Cancel
        </ModalButton>
        {isEmailValid(email) ? (
          <ModalButton
            onClick={() => {
              addAdminMutation()
                .then(() => {
                  refetch().catch((e) => console.log(e));
                  setEmail("");
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
