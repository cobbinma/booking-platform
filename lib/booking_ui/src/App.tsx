import React, { useState } from "react";
import LoginButton from "./components/LoginButton";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  useParams,
} from "react-router-dom";
import { AppState, Auth0Provider } from "@auth0/auth0-react";
import Secure from "./components/Secure";
export interface Params {
  slug: string;
  returnURL: string;
}

const App = () => {
  const [params, setParams] = useState<Params | null>(null);
  const onRedirectCallback = (appState: AppState) => {
    setParams(appState && appState.params);
  };
  return (
    <div className="App">
      <header className="App-header">
        <Auth0Provider
          domain={process.env.REACT_APP_DOMAIN ?? ""}
          clientId={process.env.REACT_APP_CLIENT_ID ?? ""}
          redirectUri={window.location.origin}
          audience={process.env.REACT_APP_AUDIENCE ?? ""}
          issuer={process.env.REACT_APP_ISSUER ?? ""}
          onRedirectCallback={onRedirectCallback}
        >
          <Router>
            <Switch>
              <Route path="/:slug/:returnURL" children={<GetParams />} />
            </Switch>
          </Router>
          {params && <Secure params={params} />}
        </Auth0Provider>
      </header>
    </div>
  );
};

const GetParams = () => {
  let params = useParams<Params>();
  return (
    <div>
      <LoginButton params={params} />
    </div>
  );
};

export default App;
