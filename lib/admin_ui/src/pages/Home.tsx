import React from "react";
import { useGetVenueQuery } from "../graph";
import { Spinner } from "baseui/spinner";
import { H2, H4 } from "baseui/typography";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";

const Home: React.FC<{ venueID: string }> = ({ venueID }) => {
  const { data, loading, error } = useGetVenueQuery({
    variables: {
      venueID: venueID,
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
      <FlexGrid>
        <FlexGridItem>
          <H2>Home</H2>
        </FlexGridItem>
        <FlexGridItem>
          <H4>Name: {data?.getVenue?.name}</H4>
        </FlexGridItem>
      </FlexGrid>
    </div>
  );
};

export default Home;
