import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import reportWebVitals from "./reportWebVitals";
import { Client as Styletron } from "styletron-engine-atomic";
import { Provider as StyletronProvider } from "styletron-react";
import { LightTheme, BaseProvider, styled } from "baseui";
import { ApolloClient, ApolloProvider, InMemoryCache } from "@apollo/client";

const engine = new Styletron();

const Theme = styled("div", {
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  height: "100%",
});

const client = new ApolloClient({
  cache: new InMemoryCache(),
  uri: "http://localhost:9999/query",
});

ReactDOM.render(
  <ApolloProvider client={client}>
    <StyletronProvider value={engine}>
      <BaseProvider theme={LightTheme}>
        <Theme>
          <App />
        </Theme>
      </BaseProvider>
    </StyletronProvider>
  </ApolloProvider>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
