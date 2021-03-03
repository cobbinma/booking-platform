import React from "react";
import { H2 } from "baseui/typography";
import { FlexGrid, FlexGridItem } from "baseui/flex-grid";
import { Input } from "baseui/input";
import { FormControl } from "baseui/form-control";

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
          <FormControl label={() => "Name"}>
            <Input value={name || ""} />
          </FormControl>
        </FlexGridItem>
        <FlexGridItem>
          <FormControl label={() => "Slug"}>
            <Input value={slug || ""} />
          </FormControl>
        </FlexGridItem>
      </FlexGrid>
    </div>
  );
};

export default Home;
