import React from "react";
import { H2, H4 } from "baseui/typography";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";

const Home: React.FC<{
  name: string | null | undefined;
  slug: string | null | undefined;
}> = ({ name, slug }) => {
  return (
    <div>
      <FlexGrid>
        <FlexGridItem>
          <H2>Home</H2>
        </FlexGridItem>
        <FlexGridItem>
          <H4>Name: {name || ""}</H4>
        </FlexGridItem>
        <FlexGridItem>
          <H4>Slug: {slug || ""}</H4>
        </FlexGridItem>
      </FlexGrid>
    </div>
  );
};

export default Home;
