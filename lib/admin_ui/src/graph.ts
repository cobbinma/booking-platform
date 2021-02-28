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
  /** Time (YYYY-MM-DDThh:mm:ssZ) */
  Time: any;
  /** Time Of Day (hh:mm) */
  TimeOfDay: any;
  /** Day of Week (Monday = 1, Sunday = 7) */
  DayOfWeek: any;
};




/** Slot Input is a booking enquiry. */
export type SlotInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** email of the customer */
  email: Scalars['String'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** desired start time of the booking (YYYY-MM-DDThh:mm:ssZ) */
  startsAt: Scalars['Time'];
  /** desired duration of the booking in minutes */
  duration: Scalars['Int'];
};

/** Slot is a possible booking that has yet to be confirmed. */
export type Slot = {
  __typename?: 'Slot';
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** email of the customer */
  email: Scalars['String'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** desired start time of the booking (YYYY-MM-DDThh:mm:ssZ) */
  startsAt: Scalars['Time'];
  /** potential ending time of the booking (YYYY-MM-DDThh:mm:ssZ) */
  endsAt: Scalars['Time'];
  /** potential duration of the booking in minutes */
  duration: Scalars['Int'];
};

/** Slot is a possible booking that has yet to be confirmed. */
export type BookingInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** email of the customer */
  email: Scalars['String'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** start time of the booking (YYYY-MM-DDThh:mm:ssZ) */
  startsAt: Scalars['Time'];
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
  /** email of the customer */
  email: Scalars['String'];
  /** amount of people attending the booking */
  people: Scalars['Int'];
  /** start time of the booking (hh:mm) */
  startsAt: Scalars['Time'];
  /** end time of the booking (hh:mm) */
  endsAt: Scalars['Time'];
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
  validFrom?: Maybe<Scalars['Time']>;
  /** date the special opening hours ends at. only valid for special opening hours */
  validThrough?: Maybe<Scalars['Time']>;
};

/** Booking Enquiry Response. */
export type GetSlotResponse = {
  __typename?: 'GetSlotResponse';
  /** slot matching the given enquiy */
  match?: Maybe<Slot>;
  /** slots have match the enquiry but have different starting times */
  otherAvailableSlots?: Maybe<Array<Slot>>;
};

export type IsAdminInput = {
  venueId: Scalars['String'];
};

/** Booking queries. */
export type Query = {
  __typename?: 'Query';
  /** get venue information from an venue identifier */
  getVenue: Venue;
  /** get slot is a booking enquiry */
  getSlot: GetSlotResponse;
  /** get slot is a booking enquiry */
  isAdmin: Scalars['Boolean'];
};


/** Booking queries. */
export type QueryGetVenueArgs = {
  id: Scalars['ID'];
};


/** Booking queries. */
export type QueryGetSlotArgs = {
  input: SlotInput;
};


/** Booking queries. */
export type QueryIsAdminArgs = {
  input: IsAdminInput;
};

/** Booking mutations. */
export type Mutation = {
  __typename?: 'Mutation';
  /** create booking is a confirming a booking slot */
  createBooking: Booking;
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
    & Pick<Booking, 'id' | 'venueId' | 'email' | 'people' | 'startsAt' | 'endsAt' | 'duration' | 'tableId'>
  ) }
);

export type GetSlotQueryVariables = Exact<{
  slot: SlotInput;
}>;


export type GetSlotQuery = (
  { __typename?: 'Query' }
  & { getSlot: (
    { __typename?: 'GetSlotResponse' }
    & { match?: Maybe<(
      { __typename?: 'Slot' }
      & Pick<Slot, 'venueId' | 'email' | 'people' | 'startsAt' | 'endsAt' | 'duration'>
    )>, otherAvailableSlots?: Maybe<Array<(
      { __typename?: 'Slot' }
      & Pick<Slot, 'venueId' | 'email' | 'people' | 'startsAt' | 'endsAt' | 'duration'>
    )>> }
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

export type IsAdminQueryVariables = Exact<{
  venueId: Scalars['String'];
}>;


export type IsAdminQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'isAdmin'>
);


export const CreateBookingDocument = gql`
    mutation CreateBooking($slot: BookingInput!) {
  createBooking(input: $slot) {
    id
    venueId
    email
    people
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
export const GetSlotDocument = gql`
    query GetSlot($slot: SlotInput!) {
  getSlot(input: $slot) {
    match {
      venueId
      email
      people
      startsAt
      endsAt
      duration
    }
    otherAvailableSlots {
      venueId
      email
      people
      startsAt
      endsAt
      duration
    }
  }
}
    `;

/**
 * __useGetSlotQuery__
 *
 * To run a query within a React component, call `useGetSlotQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSlotQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSlotQuery({
 *   variables: {
 *      slot: // value for 'slot'
 *   },
 * });
 */
export function useGetSlotQuery(baseOptions: Apollo.QueryHookOptions<GetSlotQuery, GetSlotQueryVariables>) {
        return Apollo.useQuery<GetSlotQuery, GetSlotQueryVariables>(GetSlotDocument, baseOptions);
      }
export function useGetSlotLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSlotQuery, GetSlotQueryVariables>) {
          return Apollo.useLazyQuery<GetSlotQuery, GetSlotQueryVariables>(GetSlotDocument, baseOptions);
        }
export type GetSlotQueryHookResult = ReturnType<typeof useGetSlotQuery>;
export type GetSlotLazyQueryHookResult = ReturnType<typeof useGetSlotLazyQuery>;
export type GetSlotQueryResult = Apollo.QueryResult<GetSlotQuery, GetSlotQueryVariables>;
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
export const IsAdminDocument = gql`
    query IsAdmin($venueId: String!) {
  isAdmin(input: {venueId: $venueId})
}
    `;

/**
 * __useIsAdminQuery__
 *
 * To run a query within a React component, call `useIsAdminQuery` and pass it any options that fit your needs.
 * When your component renders, `useIsAdminQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useIsAdminQuery({
 *   variables: {
 *      venueId: // value for 'venueId'
 *   },
 * });
 */
export function useIsAdminQuery(baseOptions: Apollo.QueryHookOptions<IsAdminQuery, IsAdminQueryVariables>) {
        return Apollo.useQuery<IsAdminQuery, IsAdminQueryVariables>(IsAdminDocument, baseOptions);
      }
export function useIsAdminLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<IsAdminQuery, IsAdminQueryVariables>) {
          return Apollo.useLazyQuery<IsAdminQuery, IsAdminQueryVariables>(IsAdminDocument, baseOptions);
        }
export type IsAdminQueryHookResult = ReturnType<typeof useIsAdminQuery>;
export type IsAdminLazyQueryHookResult = ReturnType<typeof useIsAdminLazyQuery>;
export type IsAdminQueryResult = Apollo.QueryResult<IsAdminQuery, IsAdminQueryVariables>;