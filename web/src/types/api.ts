export interface User {
  id: string;
  email: string;
  plan: string;
  created_at: string;
}

export interface Account {
  id: string;
  username: string;
  platform: string;
  follower_count: number;
  scan_status: string;
  last_scan_at: string | null;
  created_at: string;
}

export interface Unfollow {
  id: string;
  account_id: string;
  detected_date: string;
}

export interface UnfollowSummary {
  account_id: string;
  username: string;
  period: string;
  count: number;
  trend_change: number;
}

export interface HealthScore {
  score: number;
  trend: 'improving' | 'stable' | 'declining';
  message: string;
}

export interface LoginResponse {
  token: string;
  expires_at: number;
  user: User;
}

export interface AccountsResponse {
  accounts: Account[];
  total: number;
}

export interface UnfollowsResponse {
  unfollows: Unfollow[];
  total: number;
  limit: number;
  offset: number;
}

export interface UnfollowsSummaryResponse {
  period: string;
  summaries: UnfollowSummary[];
  overall_health: HealthScore;
}
