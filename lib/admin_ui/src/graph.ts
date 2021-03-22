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
  /** fullname of the customer */
  name: Scalars['String'];
  /** given name of the customer. in the u.k., the first name of a person. */
  givenName: Scalars['String'];
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
  /** operating hours of the venue for a specific date */
  openingHoursSpecification?: Maybe<OpeningHoursSpecification>;
  /** tables at the venue */
  tables: Array<Table>;
  /** email addresses of venue administrators */
  admins: Array<Scalars['String']>;
  /** human readable identifier of the venue */
  slug: Scalars['ID'];
  /** paginated list of bookings for a venue */
  bookings?: Maybe<BookingsPage>;
};


/** Venue where a booking can take place. */
export type VenueOpeningHoursSpecificationArgs = {
  date?: Maybe<Scalars['Time']>;
};


/** Venue where a booking can take place. */
export type VenueBookingsArgs = {
  filter?: Maybe<BookingsFilter>;
  pageInfo?: Maybe<PageInfo>;
};

/** An individual table at a venue. */
export type TableInput = {
  /** unique venue identifier the table belongs to */
  venueId: Scalars['ID'];
  /** name of the table */
  name: Scalars['String'];
  /** maximum amount of people that can sit at table */
  capacity: Scalars['Int'];
};

/** Input to remove a venue table */
export type RemoveTableInput = {
  /** unique venue identifier the table belongs to */
  venueId: Scalars['ID'];
  /** unique identifier of the table to be removed */
  tableId: Scalars['ID'];
};

/** An individual table at a venue. */
export type Table = {
  __typename?: 'Table';
  /** unique identifier of the table */
  id: Scalars['ID'];
  /** name of the table */
  name: Scalars['String'];
  /** maximum amount of people that can sit at table */
  capacity: Scalars['Int'];
};

/** Day specific operating hours. */
export type OpeningHoursSpecification = {
  __typename?: 'OpeningHoursSpecification';
  /** the day of the week for which these opening hours are valid */
  dayOfWeek: Scalars['DayOfWeek'];
  /** the opening time of the place or service on the given day(s) of the week */
  opens?: Maybe<Scalars['TimeOfDay']>;
  /** the closing time of the place or service on the given day(s) of the week */
  closes?: Maybe<Scalars['TimeOfDay']>;
  /** date the special opening hours starts at. only valid for special opening hours */
  validFrom?: Maybe<Scalars['Time']>;
  /** date the special opening hours ends at. only valid for special opening hours */
  validThrough?: Maybe<Scalars['Time']>;
};

/** Day specific operating hours. */
export type OpeningHoursSpecificationInput = {
  /** the day of the week for which these opening hours are valid */
  dayOfWeek: Scalars['DayOfWeek'];
  /** the opening time of the place or service on the given day(s) of the week */
  opens: Scalars['TimeOfDay'];
  /** the closing time of the place or service on the given day(s) of the week */
  closes: Scalars['TimeOfDay'];
};

/** Day specific special operating hours. */
export type SpecialOpeningHoursSpecificationInput = {
  /** the day of the week for which these opening hours are valid */
  dayOfWeek: Scalars['DayOfWeek'];
  /** the opening time of the place or service on the given day(s) of the week */
  opens?: Maybe<Scalars['TimeOfDay']>;
  /** the closing time of the place or service on the given day(s) of the week */
  closes?: Maybe<Scalars['TimeOfDay']>;
  /** date the special opening hours starts at. only valid for special opening hours */
  validFrom: Scalars['Time'];
  /** date the special opening hours ends at. only valid for special opening hours */
  validThrough: Scalars['Time'];
};

/** Booking Enquiry Response. */
export type GetSlotResponse = {
  __typename?: 'GetSlotResponse';
  /** slot matching the given enquiy */
  match?: Maybe<Slot>;
  /** slots have match the enquiry but have different starting times */
  otherAvailableSlots?: Maybe<Array<Slot>>;
};

/** Input to query if the user is an admin. Fields AND together. */
export type IsAdminInput = {
  /** unique identifier of the venue */
  venueId?: Maybe<Scalars['ID']>;
  /** human readable identifier of the venue */
  slug?: Maybe<Scalars['ID']>;
};

/** Filter get venue queries. Fields AND together. */
export type VenueFilter = {
  /** unique identifier of the venue */
  id?: Maybe<Scalars['ID']>;
  /** human readable identifier of the venue */
  slug?: Maybe<Scalars['ID']>;
};

/** Filter bookings. */
export type BookingsFilter = {
  /** unique identifier of the venue */
  venueId?: Maybe<Scalars['ID']>;
  /** specific date to query bookings for */
  date: Scalars['Time'];
};

/** Information about the page being requested. Maximum page limit of 50. */
export type PageInfo = {
  /** page number */
  page: Scalars['Int'];
  /** maximum amount of results per page */
  limit?: Maybe<Scalars['Int']>;
};

/** A page with a list of bookings. */
export type BookingsPage = {
  __typename?: 'BookingsPage';
  /** list of bookings */
  bookings: Array<Booking>;
  /** is there a next page */
  hasNextPage: Scalars['Boolean'];
  /** total number of pages */
  pages: Scalars['Int'];
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
  filter: VenueFilter;
};


/** Booking queries. */
export type QueryGetSlotArgs = {
  input: SlotInput;
};


/** Booking queries. */
export type QueryIsAdminArgs = {
  input: IsAdminInput;
};

/** Input to add an administrator to a venue. */
export type AdminInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** email address of the administrator */
  email: Scalars['String'];
};

/** Input to remove an administrator from a venue. */
export type RemoveAdminInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** email address of the administrator */
  email: Scalars['String'];
};

/** Input to cancel an individual booking. */
export type CancelBookingInput = {
  /** unique identifier of the venue */
  venueId?: Maybe<Scalars['ID']>;
  /** unique identifier of the booking */
  id: Scalars['ID'];
};

/** Input to update a venue's operating hours. */
export type UpdateOpeningHoursInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** operating hours of the venue */
  openingHours: Array<OpeningHoursSpecificationInput>;
};

/** Input to update a venue's special operating hours. */
export type UpdateSpecialOpeningHoursInput = {
  /** unique identifier of the venue */
  venueId: Scalars['ID'];
  /** special operating hours of the venue */
  specialOpeningHours: Array<SpecialOpeningHoursSpecificationInput>;
};

/** Booking mutations. */
export type Mutation = {
  __typename?: 'Mutation';
  /** create booking is a confirming a booking slot */
  createBooking: Booking;
  /** add a table to a venue */
  addTable: Table;
  /** remove a table from a venue */
  removeTable: Table;
  /** add an admin to a venue */
  addAdmin: Scalars['String'];
  /** remove an admin from a venue */
  removeAdmin: Scalars['String'];
  /** cancel an individual booking */
  cancelBooking: Booking;
  /** update the venue's opening hours */
  updateOpeningHours: Array<OpeningHoursSpecification>;
  /** update the venue's special opening hours */
  updateSpecialOpeningHours: Array<OpeningHoursSpecification>;
};


/** Booking mutations. */
export type MutationCreateBookingArgs = {
  input: BookingInput;
};


/** Booking mutations. */
export type MutationAddTableArgs = {
  input: TableInput;
};


/** Booking mutations. */
export type MutationRemoveTableArgs = {
  input: RemoveTableInput;
};


/** Booking mutations. */
export type MutationAddAdminArgs = {
  input: AdminInput;
};


/** Booking mutations. */
export type MutationRemoveAdminArgs = {
  input: RemoveAdminInput;
};


/** Booking mutations. */
export type MutationCancelBookingArgs = {
  input: CancelBookingInput;
};


/** Booking mutations. */
export type MutationUpdateOpeningHoursArgs = {
  input: UpdateOpeningHoursInput;
};


/** Booking mutations. */
export type MutationUpdateSpecialOpeningHoursArgs = {
  input: UpdateSpecialOpeningHoursInput;
};

export type AddAdminMutationVariables = Exact<{
  admin: AdminInput;
}>;


export type AddAdminMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'addAdmin'>
);

export type AddTableMutationVariables = Exact<{
  table: TableInput;
}>;


export type AddTableMutation = (
  { __typename?: 'Mutation' }
  & { addTable: (
    { __typename?: 'Table' }
    & Pick<Table, 'id' | 'name' | 'capacity'>
  ) }
);

export type CancelBookingMutationVariables = Exact<{
  input: CancelBookingInput;
}>;


export type CancelBookingMutation = (
  { __typename?: 'Mutation' }
  & { cancelBooking: (
    { __typename?: 'Booking' }
    & Pick<Booking, 'id'>
  ) }
);

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
  slug?: Maybe<Scalars['ID']>;
  venueID?: Maybe<Scalars['ID']>;
  filter?: Maybe<BookingsFilter>;
  pageInfo?: Maybe<PageInfo>;
  date?: Maybe<Scalars['Time']>;
}>;


export type GetVenueQuery = (
  { __typename?: 'Query' }
  & { getVenue: (
    { __typename?: 'Venue' }
    & Pick<Venue, 'id' | 'name' | 'admins' | 'slug'>
    & { openingHours: Array<(
      { __typename?: 'OpeningHoursSpecification' }
      & Pick<OpeningHoursSpecification, 'dayOfWeek' | 'opens' | 'closes'>
    )>, specialOpeningHours: Array<(
      { __typename?: 'OpeningHoursSpecification' }
      & Pick<OpeningHoursSpecification, 'dayOfWeek' | 'opens' | 'closes' | 'validFrom' | 'validThrough'>
    )>, openingHoursSpecification?: Maybe<(
      { __typename?: 'OpeningHoursSpecification' }
      & Pick<OpeningHoursSpecification, 'dayOfWeek' | 'opens' | 'closes'>
    )>, tables: Array<(
      { __typename?: 'Table' }
      & Pick<Table, 'id' | 'name' | 'capacity'>
    )>, bookings?: Maybe<(
      { __typename?: 'BookingsPage' }
      & Pick<BookingsPage, 'hasNextPage' | 'pages'>
      & { bookings: Array<(
        { __typename?: 'Booking' }
        & Pick<Booking, 'id' | 'venueId' | 'email' | 'people' | 'startsAt' | 'endsAt' | 'duration' | 'tableId' | 'name' | 'givenName'>
      )> }
    )> }
  ) }
);

export type IsAdminQueryVariables = Exact<{
  slug: Scalars['ID'];
}>;


export type IsAdminQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'isAdmin'>
);

export type RemoveAdminMutationVariables = Exact<{
  admin: RemoveAdminInput;
}>;


export type RemoveAdminMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'removeAdmin'>
);

export type RemoveTableMutationVariables = Exact<{
  table: RemoveTableInput;
}>;


export type RemoveTableMutation = (
  { __typename?: 'Mutation' }
  & { removeTable: (
    { __typename?: 'Table' }
    & Pick<Table, 'id' | 'name' | 'capacity'>
  ) }
);

export type UpdateOpeningHoursMutationVariables = Exact<{
  input: UpdateOpeningHoursInput;
}>;


export type UpdateOpeningHoursMutation = (
  { __typename?: 'Mutation' }
  & { updateOpeningHours: Array<(
    { __typename?: 'OpeningHoursSpecification' }
    & Pick<OpeningHoursSpecification, 'dayOfWeek' | 'opens' | 'closes' | 'validFrom' | 'validThrough'>
  )> }
);

export type UpdateSpecialOpeningHoursMutationVariables = Exact<{
  input: UpdateSpecialOpeningHoursInput;
}>;


export type UpdateSpecialOpeningHoursMutation = (
  { __typename?: 'Mutation' }
  & { updateSpecialOpeningHours: Array<(
    { __typename?: 'OpeningHoursSpecification' }
    & Pick<OpeningHoursSpecification, 'dayOfWeek' | 'opens' | 'closes' | 'validFrom' | 'validThrough'>
  )> }
);


export const AddAdminDocument = gql`
    mutation AddAdmin($admin: AdminInput!) {
  addAdmin(input: $admin)
}
    `;
export type AddAdminMutationFn = Apollo.MutationFunction<AddAdminMutation, AddAdminMutationVariables>;

/**
 * __useAddAdminMutation__
 *
 * To run a mutation, you first call `useAddAdminMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddAdminMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addAdminMutation, { data, loading, error }] = useAddAdminMutation({
 *   variables: {
 *      admin: // value for 'admin'
 *   },
 * });
 */
export function useAddAdminMutation(baseOptions?: Apollo.MutationHookOptions<AddAdminMutation, AddAdminMutationVariables>) {
        return Apollo.useMutation<AddAdminMutation, AddAdminMutationVariables>(AddAdminDocument, baseOptions);
      }
export type AddAdminMutationHookResult = ReturnType<typeof useAddAdminMutation>;
export type AddAdminMutationResult = Apollo.MutationResult<AddAdminMutation>;
export type AddAdminMutationOptions = Apollo.BaseMutationOptions<AddAdminMutation, AddAdminMutationVariables>;
export const AddTableDocument = gql`
    mutation AddTable($table: TableInput!) {
  addTable(input: $table) {
    id
    name
    capacity
  }
}
    `;
export type AddTableMutationFn = Apollo.MutationFunction<AddTableMutation, AddTableMutationVariables>;

/**
 * __useAddTableMutation__
 *
 * To run a mutation, you first call `useAddTableMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddTableMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addTableMutation, { data, loading, error }] = useAddTableMutation({
 *   variables: {
 *      table: // value for 'table'
 *   },
 * });
 */
export function useAddTableMutation(baseOptions?: Apollo.MutationHookOptions<AddTableMutation, AddTableMutationVariables>) {
        return Apollo.useMutation<AddTableMutation, AddTableMutationVariables>(AddTableDocument, baseOptions);
      }
export type AddTableMutationHookResult = ReturnType<typeof useAddTableMutation>;
export type AddTableMutationResult = Apollo.MutationResult<AddTableMutation>;
export type AddTableMutationOptions = Apollo.BaseMutationOptions<AddTableMutation, AddTableMutationVariables>;
export const CancelBookingDocument = gql`
    mutation CancelBooking($input: CancelBookingInput!) {
  cancelBooking(input: $input) {
    id
  }
}
    `;
export type CancelBookingMutationFn = Apollo.MutationFunction<CancelBookingMutation, CancelBookingMutationVariables>;

/**
 * __useCancelBookingMutation__
 *
 * To run a mutation, you first call `useCancelBookingMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCancelBookingMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [cancelBookingMutation, { data, loading, error }] = useCancelBookingMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useCancelBookingMutation(baseOptions?: Apollo.MutationHookOptions<CancelBookingMutation, CancelBookingMutationVariables>) {
        return Apollo.useMutation<CancelBookingMutation, CancelBookingMutationVariables>(CancelBookingDocument, baseOptions);
      }
export type CancelBookingMutationHookResult = ReturnType<typeof useCancelBookingMutation>;
export type CancelBookingMutationResult = Apollo.MutationResult<CancelBookingMutation>;
export type CancelBookingMutationOptions = Apollo.BaseMutationOptions<CancelBookingMutation, CancelBookingMutationVariables>;
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
    query GetVenue($slug: ID, $venueID: ID, $filter: BookingsFilter, $pageInfo: PageInfo, $date: Time) {
  getVenue(filter: {slug: $slug, id: $venueID}) {
    id
    name
    openingHours {
      dayOfWeek
      opens
      closes
    }
    specialOpeningHours {
      dayOfWeek
      opens
      closes
      validFrom
      validThrough
    }
    openingHoursSpecification(date: $date) {
      dayOfWeek
      opens
      closes
    }
    tables {
      id
      name
      capacity
    }
    admins
    slug
    bookings(filter: $filter, pageInfo: $pageInfo) {
      bookings {
        id
        venueId
        email
        people
        startsAt
        endsAt
        duration
        tableId
        name
        givenName
      }
      hasNextPage
      pages
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
 *      slug: // value for 'slug'
 *      venueID: // value for 'venueID'
 *      filter: // value for 'filter'
 *      pageInfo: // value for 'pageInfo'
 *      date: // value for 'date'
 *   },
 * });
 */
export function useGetVenueQuery(baseOptions?: Apollo.QueryHookOptions<GetVenueQuery, GetVenueQueryVariables>) {
        return Apollo.useQuery<GetVenueQuery, GetVenueQueryVariables>(GetVenueDocument, baseOptions);
      }
export function useGetVenueLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetVenueQuery, GetVenueQueryVariables>) {
          return Apollo.useLazyQuery<GetVenueQuery, GetVenueQueryVariables>(GetVenueDocument, baseOptions);
        }
export type GetVenueQueryHookResult = ReturnType<typeof useGetVenueQuery>;
export type GetVenueLazyQueryHookResult = ReturnType<typeof useGetVenueLazyQuery>;
export type GetVenueQueryResult = Apollo.QueryResult<GetVenueQuery, GetVenueQueryVariables>;
export const IsAdminDocument = gql`
    query IsAdmin($slug: ID!) {
  isAdmin(input: {slug: $slug})
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
 *      slug: // value for 'slug'
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
export const RemoveAdminDocument = gql`
    mutation RemoveAdmin($admin: RemoveAdminInput!) {
  removeAdmin(input: $admin)
}
    `;
export type RemoveAdminMutationFn = Apollo.MutationFunction<RemoveAdminMutation, RemoveAdminMutationVariables>;

/**
 * __useRemoveAdminMutation__
 *
 * To run a mutation, you first call `useRemoveAdminMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveAdminMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeAdminMutation, { data, loading, error }] = useRemoveAdminMutation({
 *   variables: {
 *      admin: // value for 'admin'
 *   },
 * });
 */
export function useRemoveAdminMutation(baseOptions?: Apollo.MutationHookOptions<RemoveAdminMutation, RemoveAdminMutationVariables>) {
        return Apollo.useMutation<RemoveAdminMutation, RemoveAdminMutationVariables>(RemoveAdminDocument, baseOptions);
      }
export type RemoveAdminMutationHookResult = ReturnType<typeof useRemoveAdminMutation>;
export type RemoveAdminMutationResult = Apollo.MutationResult<RemoveAdminMutation>;
export type RemoveAdminMutationOptions = Apollo.BaseMutationOptions<RemoveAdminMutation, RemoveAdminMutationVariables>;
export const RemoveTableDocument = gql`
    mutation RemoveTable($table: RemoveTableInput!) {
  removeTable(input: $table) {
    id
    name
    capacity
  }
}
    `;
export type RemoveTableMutationFn = Apollo.MutationFunction<RemoveTableMutation, RemoveTableMutationVariables>;

/**
 * __useRemoveTableMutation__
 *
 * To run a mutation, you first call `useRemoveTableMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveTableMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeTableMutation, { data, loading, error }] = useRemoveTableMutation({
 *   variables: {
 *      table: // value for 'table'
 *   },
 * });
 */
export function useRemoveTableMutation(baseOptions?: Apollo.MutationHookOptions<RemoveTableMutation, RemoveTableMutationVariables>) {
        return Apollo.useMutation<RemoveTableMutation, RemoveTableMutationVariables>(RemoveTableDocument, baseOptions);
      }
export type RemoveTableMutationHookResult = ReturnType<typeof useRemoveTableMutation>;
export type RemoveTableMutationResult = Apollo.MutationResult<RemoveTableMutation>;
export type RemoveTableMutationOptions = Apollo.BaseMutationOptions<RemoveTableMutation, RemoveTableMutationVariables>;
export const UpdateOpeningHoursDocument = gql`
    mutation UpdateOpeningHours($input: UpdateOpeningHoursInput!) {
  updateOpeningHours(input: $input) {
    dayOfWeek
    opens
    closes
    validFrom
    validThrough
  }
}
    `;
export type UpdateOpeningHoursMutationFn = Apollo.MutationFunction<UpdateOpeningHoursMutation, UpdateOpeningHoursMutationVariables>;

/**
 * __useUpdateOpeningHoursMutation__
 *
 * To run a mutation, you first call `useUpdateOpeningHoursMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateOpeningHoursMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateOpeningHoursMutation, { data, loading, error }] = useUpdateOpeningHoursMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateOpeningHoursMutation(baseOptions?: Apollo.MutationHookOptions<UpdateOpeningHoursMutation, UpdateOpeningHoursMutationVariables>) {
        return Apollo.useMutation<UpdateOpeningHoursMutation, UpdateOpeningHoursMutationVariables>(UpdateOpeningHoursDocument, baseOptions);
      }
export type UpdateOpeningHoursMutationHookResult = ReturnType<typeof useUpdateOpeningHoursMutation>;
export type UpdateOpeningHoursMutationResult = Apollo.MutationResult<UpdateOpeningHoursMutation>;
export type UpdateOpeningHoursMutationOptions = Apollo.BaseMutationOptions<UpdateOpeningHoursMutation, UpdateOpeningHoursMutationVariables>;
export const UpdateSpecialOpeningHoursDocument = gql`
    mutation UpdateSpecialOpeningHours($input: UpdateSpecialOpeningHoursInput!) {
  updateSpecialOpeningHours(input: $input) {
    dayOfWeek
    opens
    closes
    validFrom
    validThrough
  }
}
    `;
export type UpdateSpecialOpeningHoursMutationFn = Apollo.MutationFunction<UpdateSpecialOpeningHoursMutation, UpdateSpecialOpeningHoursMutationVariables>;

/**
 * __useUpdateSpecialOpeningHoursMutation__
 *
 * To run a mutation, you first call `useUpdateSpecialOpeningHoursMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSpecialOpeningHoursMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSpecialOpeningHoursMutation, { data, loading, error }] = useUpdateSpecialOpeningHoursMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useUpdateSpecialOpeningHoursMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSpecialOpeningHoursMutation, UpdateSpecialOpeningHoursMutationVariables>) {
        return Apollo.useMutation<UpdateSpecialOpeningHoursMutation, UpdateSpecialOpeningHoursMutationVariables>(UpdateSpecialOpeningHoursDocument, baseOptions);
      }
export type UpdateSpecialOpeningHoursMutationHookResult = ReturnType<typeof useUpdateSpecialOpeningHoursMutation>;
export type UpdateSpecialOpeningHoursMutationResult = Apollo.MutationResult<UpdateSpecialOpeningHoursMutation>;
export type UpdateSpecialOpeningHoursMutationOptions = Apollo.BaseMutationOptions<UpdateSpecialOpeningHoursMutation, UpdateSpecialOpeningHoursMutationVariables>;