import React from "react";
import { useIsAdminQuery } from "../graph";
import { Params } from "../App";
import { Spinner } from "baseui/spinner";
import Admin from "./Admin";
import { BrowserRouter } from "react-router-dom";

const IsAdmin: React.FC<{ params: Params }> = ({ params }) => {
  const { data, loading, error } = useIsAdminQuery({
    variables: {
      venueId: params.venueId,
    },
  });

  if (loading)
    return (
      <div>
        <Spinner />
      </div>
    );

  if (error) {
    console.log(error);
    return <p>error</p>;
  }

  if (data == null || !data.isAdmin) {
    return <p>you must be admin</p>;
  }

  return (
    <div>
      <BrowserRouter>
        <Admin />
      </BrowserRouter>
    </div>
  );
};

export default IsAdmin;
