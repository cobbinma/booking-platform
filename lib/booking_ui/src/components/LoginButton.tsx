import React from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { Params } from "../App";

const LoginButton: React.FC<{ params: Params }> = ({ params }) => {
  const { loginWithRedirect } = useAuth0();

  return (
    <button onClick={() => loginWithRedirect({ appState: { params: params } })}>
      Log In
    </button>
  );
};

export default LoginButton;
