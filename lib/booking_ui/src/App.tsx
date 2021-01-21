import React from "react";
import LoginButton from "./components/LoginButton";
import LogoutButton from "./components/LogoutButton";
import Profile from "./components/Profile";
import { useGetVenueQuery } from "./graph";

function App() {
  const { data, error, loading } = useGetVenueQuery({
    variables: {
      venueID: "12345",
    },
  });
  if (loading) return <p>Loading ...</p>;
  if (error) {
    console.log(error);
    return <p>error</p>;
  }
  console.log(data);
  return (
    <div className="App">
      <header className="App-header">
        <LoginButton />
        <LogoutButton />
        <Profile />
      </header>
    </div>
  );
}

export default App;
