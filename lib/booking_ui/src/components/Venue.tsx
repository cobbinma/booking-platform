import React from "react";
import { useGetVenueQuery } from "../graph";

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

  if (loading) return <p>Loading...</p>;
  if (error) return <p>error</p>;

  return (
    <div>
      <h2>{data?.getVenue.name}</h2>
      {data?.getVenue?.openingHours
        .slice()
        .sort((h1, h2) => h1.dayOfWeek - h2.dayOfWeek)
        .map((hours) => {
          return (
            <div key={hours.dayOfWeek}>
              <p>
                {weekday[hours.dayOfWeek]}: opens {hours.opens} until{" "}
                {hours.closes}
              </p>
            </div>
          );
        })}
    </div>
  );
};

export default Venue;
