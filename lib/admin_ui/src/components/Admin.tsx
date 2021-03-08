import React from "react";
import { Switch, Route, useHistory } from "react-router-dom";
import Home from "../pages/Home";
import Tables from "../pages/Tables";
import { AppNavBar, NavItemT } from "baseui/app-nav-bar";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { useStyletron } from "baseui";
import Admins from "../pages/Admins";
import Bookings from "../pages/Bookings";
import { useGetVenueQuery } from "../graph";
import { Spinner } from "baseui/spinner";

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

  const { data, loading, error, refetch } = useGetVenueQuery({
    variables: {
      slug: venueID,
    },
  });

  if (loading)
    return (
      <div>
        <Spinner />
      </div>
    );

  if (error) {
    console.log(error);
    return <p>error</p>;
  }

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
                history.push(item.info.link);
              }}
            />
          </div>
        </FlexGridItem>
        <FlexGridItem>
          <Switch>
            <Route path="/tables">
              <Tables
                tables={data?.getVenue?.tables || []}
                venueId={data?.getVenue?.id}
                refetch={refetch}
              />
            </Route>
            <Route path="/admins">
              <Admins
                admins={data?.getVenue?.admins || []}
                venueId={data?.getVenue?.id}
                refetch={refetch}
              />
            </Route>
            <Route path="/bookings">
              <Bookings
                tables={data?.getVenue?.tables || []}
                bookings={data?.getVenue?.bookings?.bookings || []}
                slug={venueID}
                venueId={data?.getVenue?.id}
                refetch={refetch}
              />
            </Route>
            <Route path="/">
              <Home name={data?.getVenue?.name} slug={data?.getVenue?.slug} />
            </Route>
          </Switch>
        </FlexGridItem>
      </FlexGrid>
    </div>
  );
};

export default Admin;
