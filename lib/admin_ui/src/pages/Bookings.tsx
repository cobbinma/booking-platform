import React from "react";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { H2 } from "baseui/typography";

const Bookings = () => {
  return (
    <div>
      <FlexGrid>
        <FlexGridItem>
          <H2>Bookings</H2>
        </FlexGridItem>
      </FlexGrid>
    </div>
  );
};

export default Bookings;
