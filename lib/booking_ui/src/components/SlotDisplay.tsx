import React from "react";
import { Table } from "baseui/table";

const SlotDisplay: React.FC<{
  people: number;
  startsAt: string;
  endsAt: string;
}> = ({ people, startsAt, endsAt }) => {
  const with_leading = (t: number): string => {
    return ("0" + t).slice(-2);
  };
  const starts: Date = new Date(Date.parse(startsAt));
  const ends: Date = new Date(Date.parse(endsAt));
  return (
    <div>
      <Table
        columns={["Guests", "Date", "Starts", "Ends"]}
        data={[
          [
            people,
            `${starts.getDate()}/${
              starts.getMonth() + 1
            }/${starts.getFullYear()}`,
            `${with_leading(starts.getHours())}:${with_leading(
              starts.getMinutes()
            )}`,
            `${with_leading(ends.getHours())}:${with_leading(
              ends.getMinutes()
            )}`,
          ],
        ]}
      />
    </div>
  );
};

export default SlotDisplay;
