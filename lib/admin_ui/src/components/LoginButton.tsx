import React from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { Params } from "../App";
import { Button } from "baseui/button";

const LoginButton: React.FC<{ params: Params }> = ({ params }) => {
  const { loginWithRedirect } = useAuth0();

  return (
    <Button onClick={() => loginWithRedirect({ appState: { params: params } })}>
      Log In
    </Button>
  );
};

export default LoginButton;
