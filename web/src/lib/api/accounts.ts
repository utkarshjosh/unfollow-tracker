import { api } from './api';
import type { AccountsResponse, Account } from '@/types/api';

export async function listAccounts(): Promise<AccountsResponse> {
  const response = await api.get<AccountsResponse>('/api/v1/accounts');
  return response.data;
}

export async function createAccount(username: string, platform = 'instagram'): Promise<Account> {
  const response = await api.post<Account>('/api/v1/accounts', {
    username,
    platform,
  });
  return response.data;
}

export async function deleteAccount(accountId: string): Promise<void> {
  await api.delete(`/api/v1/accounts/${accountId}`);
}
