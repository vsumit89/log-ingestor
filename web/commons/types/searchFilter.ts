export interface SearchFilters {}

export interface FilterOptions {
  limit?: number;
  resourceId?: string;
  level?: string;
  message?: string;
  traceId?: string;
  spanId?: string;
  commit?: string;
  parentResourceId?: string;
  from?: string;
  to?: string;
}
