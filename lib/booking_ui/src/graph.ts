import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  /** Date (dd-mm-yyyy) */
  Date: any;
  /** Time (hh:mm) */
  TimeOfDay: any;
  /** Day of Week (Monday = 1, Sunday = 7) */
  DayOfWeek: any;
};




/** Slot Input is a booking enquiry. */
export type SlotInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** unique identifier of the customer */
  customerId: Scalars['ID'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** desired date of the booking (dd-mm-yyyy) */
  date: Scalars['Date'];
  /** desired start time of the booking (hh:mm) */
  startsAt: Scalars['TimeOfDay'];
  /** desired duration of the booking in minutes */
  duration: Scalars['Int'];
};

/** Slot is a possible booking that has yet to be confirmed. */
export type Slot = {
  __typename?: 'Slot';
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** unique identifier of the customer */
  customerId: Scalars['ID'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** potential date of the booking (dd-mm-yyyy) */
  date: Scalars['Date'];
  /** potential start time of the booking (hh:mm) */
  startsAt: Scalars['TimeOfDay'];
  /** potential ending time of the booking (hh:mm) */
  endsAt: Scalars['TimeOfDay'];
  /** potential duration of the booking in minutes */
  duration: Scalars['Int'];
};

/** Slot is a possible booking that has yet to be confirmed. */
export type BookingInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** unique identifier of the customer */
  customerId: Scalars['ID'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** date of the booking (dd-mm-yyyy) */
  date: Scalars['Date'];
  /** start time of the booking (hh:mm) */
  startsAt: Scalars['TimeOfDay'];
  /** duration of the booking in minutes */
  duration: Scalars['Int'];
};

/** Booking has now been confirmed. */
export type Booking = {
  __typename?: 'Booking';
  /** unique identifier of the booking */
  id: Scalars['ID'];
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** unique identifier of the customer */
  customerId: Scalars['ID'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** date of the booking (dd-mm-yyyy) */
  date: Scalars['Date'];
  /** start time of the booking (hh:mm) */
  startsAt: Scalars['TimeOfDay'];
  /** end time of the booking (hh:mm) */
  endsAt: Scalars['TimeOfDay'];
  /** duration of the booking in minutes */
  duration: Scalars['Int'];
  /** unique identifier of the booking table */
  tableId: Scalars['ID'];
};

/** Venue where a booking can take place. */
export type Venue = {
  __typename?: 'Venue';
  /** unique identifier of the venue */
  id: Scalars['ID'];
  /** name of the venue */
  name: Scalars['String'];
  /** operating hours of the venue */
  openingHours: Array<OpeningHoursSpecification>;
  /** special operating hours of the venue */
  specialOpeningHours: Array<OpeningHoursSpecification>;
};

/** Day specific operating hours. */
export type OpeningHoursSpecification = {
  __typename?: 'OpeningHoursSpecification';
  /** the day of the week for which these opening hours are valid */
  dayOfWeek: Scalars['DayOfWeek'];
  /** the opening time of the place or service on the given day(s) of the week */
  opens: Scalars['TimeOfDay'];
  /** the closing time of the place or service on the given day(s) of the week */
  closes: Scalars['TimeOfDay'];
  /** date the special opening hours starts at. only valid for special opening hours */
  validFrom?: Maybe<Scalars['Date']>;
  /** date the special opening hours ends at. only valid for special opening hours */
  validThrough?: Maybe<Scalars['Date']>;
};

/** Booking queries. */
export type Query = {
  __typename?: 'Query';
  /** get venue information from an venue identifier */
  getVenue: Venue;
};


/** Booking queries. */
export type QueryGetVenueArgs = {
  id: Scalars['ID'];
};

/** Booking mutations. */
export type Mutation = {
  __typename?: 'Mutation';
  /** create slot is a booking enquiry */
  createSlot: Slot;
  /** create booking is a confirming a booking slot */
  createBooking: Booking;
};


/** Booking mutations. */
export type MutationCreateSlotArgs = {
  input: SlotInput;
};


/** Booking mutations. */
export type MutationCreateBookingArgs = {
  input: BookingInput;
};

export type CreateBookingMutationVariables = Exact<{
  slot: BookingInput;
}>;


export type CreateBookingMutation = (
  { __typename?: 'Mutation' }
  & { createBooking: (
    { __typename?: 'Booking' }
    & Pick<Booking, 'id' | 'venueId' | 'customerId' | 'people' | 'date' | 'startsAt' | 'endsAt' | 'duration' | 'tableId'>
  ) }
);

export type CreateSlotMutationVariables = Exact<{
  slot: SlotInput;
}>;


export type CreateSlotMutation = (
  { __typename?: 'Mutation' }
  & { createSlot: (
    { __typename?: 'Slot' }
    & Pick<Slot, 'venueId' | 'customerId' | 'people' | 'date' | 'startsAt' | 'endsAt' | 'duration'>
  ) }
);

export type GetVenueQueryVariables = Exact<{
  venueID: Scalars['ID'];
}>;


export type GetVenueQuery = (
  { __typename?: 'Query' }
  & { getVenue: (
    { __typename?: 'Venue' }
    & Pick<Venue, 'id' | 'name'>
    & { openingHours: Array<(
      { __typename?: 'OpeningHoursSpecification' }
      & Pick<OpeningHoursSpecification, 'dayOfWeek' | 'opens' | 'closes' | 'validFrom' | 'validThrough'>
    )>, specialOpeningHours: Array<(
      { __typename?: 'OpeningHoursSpecification' }
      & Pick<OpeningHoursSpecification, 'dayOfWeek' | 'opens' | 'closes' | 'validFrom' | 'validThrough'>
    )> }
  ) }
);


export const CreateBookingDocument = gql`
    mutation CreateBooking($slot: BookingInput!) {
  createBooking(input: $slot) {
    id
    venueId
    customerId
    people
    date
    startsAt
    endsAt
    duration
    tableId
  }
}
    `;
export type CreateBookingMutationFn = Apollo.MutationFunction<CreateBookingMutation, CreateBookingMutationVariables>;

/**
 * __useCreateBookingMutation__
 *
 * To run a mutation, you first call `useCreateBookingMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateBookingMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createBookingMutation, { data, loading, error }] = useCreateBookingMutation({
 *   variables: {
 *      slot: // value for 'slot'
 *   },
 * });
 */
export function useCreateBookingMutation(baseOptions?: Apollo.MutationHookOptions<CreateBookingMutation, CreateBookingMutationVariables>) {
        return Apollo.useMutation<CreateBookingMutation, CreateBookingMutationVariables>(CreateBookingDocument, baseOptions);
      }
export type CreateBookingMutationHookResult = ReturnType<typeof useCreateBookingMutation>;
export type CreateBookingMutationResult = Apollo.MutationResult<CreateBookingMutation>;
export type CreateBookingMutationOptions = Apollo.BaseMutationOptions<CreateBookingMutation, CreateBookingMutationVariables>;
export const CreateSlotDocument = gql`
    mutation CreateSlot($slot: SlotInput!) {
  createSlot(input: $slot) {
    venueId
    customerId
    people
    date
    startsAt
    endsAt
    duration
  }
}
    `;
export type CreateSlotMutationFn = Apollo.MutationFunction<CreateSlotMutation, CreateSlotMutationVariables>;

/**
 * __useCreateSlotMutation__
 *
 * To run a mutation, you first call `useCreateSlotMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateSlotMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createSlotMutation, { data, loading, error }] = useCreateSlotMutation({
 *   variables: {
 *      slot: // value for 'slot'
 *   },
 * });
 */
export function useCreateSlotMutation(baseOptions?: Apollo.MutationHookOptions<CreateSlotMutation, CreateSlotMutationVariables>) {
        return Apollo.useMutation<CreateSlotMutation, CreateSlotMutationVariables>(CreateSlotDocument, baseOptions);
      }
export type CreateSlotMutationHookResult = ReturnType<typeof useCreateSlotMutation>;
export type CreateSlotMutationResult = Apollo.MutationResult<CreateSlotMutation>;
export type CreateSlotMutationOptions = Apollo.BaseMutationOptions<CreateSlotMutation, CreateSlotMutationVariables>;
export const GetVenueDocument = gql`
    query GetVenue($venueID: ID!) {
  getVenue(id: $venueID) {
    id
    name
    openingHours {
      dayOfWeek
      opens
      closes
      validFrom
      validThrough
    }
    specialOpeningHours {
      dayOfWeek
      opens
      closes
      validFrom
      validThrough
    }
  }
}
    `;

/**
 * __useGetVenueQuery__
 *
 * To run a query within a React component, call `useGetVenueQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetVenueQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetVenueQuery({
 *   variables: {
 *      venueID: // value for 'venueID'
 *   },
 * });
 */
export function useGetVenueQuery(baseOptions: Apollo.QueryHookOptions<GetVenueQuery, GetVenueQueryVariables>) {
        return Apollo.useQuery<GetVenueQuery, GetVenueQueryVariables>(GetVenueDocument, baseOptions);
      }
export function useGetVenueLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetVenueQuery, GetVenueQueryVariables>) {
          return Apollo.useLazyQuery<GetVenueQuery, GetVenueQueryVariables>(GetVenueDocument, baseOptions);
        }
export type GetVenueQueryHookResult = ReturnType<typeof useGetVenueQuery>;
export type GetVenueLazyQueryHookResult = ReturnType<typeof useGetVenueLazyQuery>;
export type GetVenueQueryResult = Apollo.QueryResult<GetVenueQuery, GetVenueQueryVariables>;