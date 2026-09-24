export type Address = { address: string; lat: number; lon: number };

export type BusinessType = {
  id: string;
  infra_type_id: string;
  infra_type: { id: string; slug: string; name: string; weight: number; max_radius: number };
};

export type TrackedLocation = Address & {
  id: string;
  name: string;
  business_type_id: string;
  rating: number | null;
  rating_calculated_at: string | null;
};

export type CreatedTrackedLocation = TrackedLocation & {
  rating: number;
  rating_calculated_at: string;
};

export type CreatePointRequest = Address & { name: string; business_type_id: string };
export type RatingHistoryEntry = { value: number; calculated_at: string };
export type PointDraft = {
  name: string;
  business_type_id: string;
  addressQuery: string;
  selectedAddress: Address | null;
};

export type LoadState<T> =
  | { status: 'loading' }
  | { status: 'success'; data: T }
  | { status: 'error'; message: string };

export type ApiErrorCode =
  | 'validation_error' | 'unauthorized' | 'not_found'
  | 'provider_unavailable' | 'internal_error';
