import React, { useState } from "react";
import LoginButton from "./components/LoginButton";
import LogoutButton from "./components/LogoutButton";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  useParams,
} from "react-router-dom";
import Profile from "./components/Profile";
import { AppState, Auth0Provider } from "@auth0/auth0-react";

export interface Params {
  venueId: string;
  returnURL: string;
}

function App() {
  const [params, setParams] = useState<Params | null>(null);
  const onRedirectCallback = (appState: AppState) => {
    setParams(appState && appState.params);
  };
  return (
    <div className="App">
      <header className="App-header">
        <Auth0Provider
          domain={process.env.REACT_APP_DOMAIN ?? ""}
          clientId={process.env.REACT_APP_CLIENT_ID! ?? ""}
          redirectUri={window.location.origin}
          onRedirectCallback={onRedirectCallback}
        >
          <Router>
            <Switch>
              <Route path="/:venueId/:returnURL" children={<Child />} />
            </Switch>
          </Router>
          {params && (
            <div>
              venueId: {params.venueId}, return url: {params.returnURL}
              <LogoutButton />
              <Profile />
            </div>
          )}
        </Auth0Provider>
      </header>
    </div>
  );
}

function Child() {
  let params = useParams<Params>();
  return (
    <div>
      <LoginButton params={params} />
    </div>
  );
}

export default App;
