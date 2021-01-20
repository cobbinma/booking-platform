import React from "react";
import { useAuth0 } from "@auth0/auth0-react";

const LoginButton = () => {
  const { loginWithRedirect } = useAuth0();

  return (
    <button
      onClick={() =>
        loginWithRedirect({ appState: { targetUrl: window.location.pathname } })
      }
    >
      Log In
    </button>
  );
};

export default LoginButton;
