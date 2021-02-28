import React from "react";
import { Params } from "../App";
import { useAuth0 } from "@auth0/auth0-react";
import {
  ApolloClient,
  ApolloProvider,
  createHttpLink,
  InMemoryCache,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import IsAdmin from "./IsAdmin";
import { Spinner } from "baseui/spinner";

const Secure: React.FC<{ params: Params }> = ({ params }) => {
  const { isAuthenticated, isLoading, getAccessTokenSilently } = useAuth0();

  if (isLoading) {
    return (
      <div>
        <Spinner />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <div>Please Login...</div>;
  }

  const httpLink = createHttpLink({
    uri: "http://localhost:9999/query",
  });

  const authLink = setContext(async (_, { headers }) => {
    // get the authentication token from local storage if it exists
    const token = await getAccessTokenSilently().catch((e) => {
      console.log(e);
    });
    // return the headers to the context so httpLink can read them
    return {
      headers: {
        ...headers,
        authorization: token ? `Bearer ${token}` : "",
      },
    };
  });

  const client = new ApolloClient({
    cache: new InMemoryCache(),
    link: authLink.concat(httpLink),
  });

  return (
    <div>
      <ApolloProvider client={client}>
        <IsAdmin params={params} />
      </ApolloProvider>
    </div>
  );
};

export default Secure;
