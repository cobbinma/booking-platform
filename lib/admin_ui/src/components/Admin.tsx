import React from "react";
import { Switch, Route, useHistory } from "react-router-dom";
import Home from "../pages/Home";
import Tables from "../pages/Tables";
import { AppNavBar, NavItemT, setItemActive } from "baseui/app-nav-bar";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { useStyletron } from "baseui";
import Admins from "../pages/Admins";
import Bookings from "../pages/Bookings";

const Pages: NavItemT[] = [
  {
    label: "Home",
    info: {
      link: "/",
    },
  },
  {
    label: "Bookings",
    info: {
      link: "/bookings",
    },
  },
  {
    label: "Tables",
    info: {
      link: "/tables",
    },
  },
  {
    label: "Administrators",
    info: {
      link: "/admins",
    },
  },
];

const Admin: React.FC<{ venueID: string }> = ({ venueID }) => {
  let history = useHistory();
  const [css] = useStyletron();

  return (
    <div>
      <FlexGrid flexGridColumnCount={1}>
        <FlexGridItem>
          <div
            className={css({
              boxSizing: "border-box",
              width: "100vw",
              top: "0",
              left: "0",
            })}
          >
            <AppNavBar
              title="Admin"
              mainItems={Pages}
              onMainItemSelect={(item) => {
                console.log(item);
                setItemActive(Pages, item);
                history.push(item.info.link);
              }}
            />
          </div>
        </FlexGridItem>
        <FlexGridItem>
          <Switch>
            <Route path="/tables">
              <Tables />
            </Route>
            <Route path="/admins">
              <Admins />
            </Route>
            <Route path="/bookings">
              <Bookings />
            </Route>
            <Route path="/">
              <Home venueID={venueID} />
            </Route>
          </Switch>
        </FlexGridItem>
      </FlexGrid>
    </div>
  );
};

export default Admin;