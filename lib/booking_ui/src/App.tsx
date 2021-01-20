import React from "react";
import LoginButton from "./components/LoginButton";
import LogoutButton from "./components/LogoutButton";
import Profile from "./components/Profile";

function App() {
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
