import React from "react";
import { useGetVenueQuery } from "../graph";
import { Spinner } from "baseui/spinner";
import { H2 } from "baseui/typography";
import { Table } from "baseui/table";

const weekday = new Array(7);
weekday[0] = "Unknown";
weekday[1] = "Monday";
weekday[2] = "Tuesday";
weekday[3] = "Wednesday";
weekday[4] = "Thursday";
weekday[5] = "Friday";
weekday[6] = "Saturday";
weekday[7] = "Sunday";

const Venue: React.FC<{ venueId: string }> = ({ venueId }) => {
  const { data, loading, error } = useGetVenueQuery({
    variables: {
      venueID: venueId,
    },
  });

  if (loading)
    return (
      <div>
        <Spinner />
      </div>
    );
  if (error) return <p>error</p>;

  const opening: React.ReactNode[][] | undefined = data?.getVenue?.openingHours
    .slice()
    .sort((h1, h2) => h1.dayOfWeek - h2.dayOfWeek)
    .map((hours) => {
      return [weekday[hours.dayOfWeek], hours.opens, hours.closes];
    });

  return (
    <div>
      <H2>{data?.getVenue.name.toLowerCase()}</H2>
      {opening ? (
        <Table columns={["Day", "Opens", "Closes"]} data={opening} />
      ) : (
        <div />
      )}
    </div>
  );
};

export default Venue;
