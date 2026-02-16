import { api } from './api';
import type { UnfollowsResponse, UnfollowsSummaryResponse } from '@/types/api';

export async function getUnfollows(
  accountId?: string,
  limit = 50,
  offset = 0
): Promise<UnfollowsResponse> {
  const params = new URLSearchParams();
  if (accountId) params.append('account_id', accountId);
  params.append('limit', limit.toString());
  params.append('offset', offset.toString());

  const response = await api.get<UnfollowsResponse>(`/api/v1/unfollows?${params}`);
  return response.data;
}

export async function getUnfollowsSummary(
  period: 'day' | 'week' | 'month' = 'week'
): Promise<UnfollowsSummaryResponse> {
  const response = await api.get<UnfollowsSummaryResponse>(
    `/api/v1/unfollows/summary?period=${period}`
  );
  return response.data;
}
