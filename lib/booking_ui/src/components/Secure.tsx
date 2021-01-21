import React from "react";
import { Params } from "../App";
import Booking from "./Booking";
import { useAuth0 } from "@auth0/auth0-react";

const Secure: React.FC<{ params: Params }> = ({ params }) => {
  const { user, isAuthenticated, isLoading } = useAuth0();

  if (isLoading) {
    return <div>Loading ...</div>;
  }

  if (!isAuthenticated) {
    return <div>Please Login...</div>;
  }

  return (
    <div>
      <Booking params={params} email={user.email} />
    </div>
  );
};

export default Secure;
