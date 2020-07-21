import React from "react";
import NavBar from "./components/layout/AppBar";
import { BrowserRouter, Route } from "react-router-dom";
import Book from "./components/pages/Book";
import Home from "./components/pages/Home";

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <NavBar />
        <Route exact={true} path="/" component={Home} />
        <Route path="/book" component={Book} />
      </BrowserRouter>
    </div>
  );
}

export default App;
