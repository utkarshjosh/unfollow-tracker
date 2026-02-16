import { useEffect, useState } from 'react';
import { listAccounts } from '@/lib/api/accounts';
import { getUnfollowsSummary } from '@/lib/api/unfollows';
import { Skeleton } from '@/components/ui/skeleton';
import { Users, UserMinus, Activity, Heart } from 'lucide-react';

interface StatsData {
  totalFollowers: number;
  recentUnfollowers: number;
  engagementRate: number;
  healthScore: number;
  healthTrend: 'improving' | 'stable' | 'declining';
  healthMessage: string;
  followerChange: number;
  unfollowerChange: number;
}

export function StatsCards() {
  const [stats, setStats] = useState<StatsData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchStats() {
      try {
        setLoading(true);
        setError(null);

        // Fetch accounts and unfollows summary in parallel
        const [accountsRes, summaryRes] = await Promise.all([
          listAccounts(),
          getUnfollowsSummary('week'),
        ]);

        // Calculate total followers from all accounts
        const totalFollowers = accountsRes.accounts.reduce(
          (sum, acc) => sum + (acc.follower_count || 0),
          0
        );

        // Calculate recent unfollowers (last 7 days from summaries)
        const recentUnfollowers = summaryRes.summaries.reduce(
          (sum, s) => sum + s.count,
          0
        );

        // Calculate engagement rate (mock calculation for now)
        const engagementRate = totalFollowers > 0
          ? Math.min(4.2 + (totalFollowers / 100000), 8.5)
          : 4.2;

        // Health score from API
        const healthScore = summaryRes.overall_health?.score ?? 100;
        const healthTrend = summaryRes.overall_health?.trend ?? 'stable';
        const healthMessage = summaryRes.overall_health?.message ?? 'Your audience is healthy!';

        // Mock changes for demo (would come from historical comparison)
        const followerChange = Math.floor(Math.random() * 200) - 50;
        const unfollowerChange = Math.floor(Math.random() * 20) - 5;

        setStats({
          totalFollowers,
          recentUnfollowers,
          engagementRate,
          healthScore,
          healthTrend,
          healthMessage,
          followerChange,
          unfollowerChange,
        });
      } catch (err: any) {
        setError(err.response?.data?.message || 'Failed to load stats');
      } finally {
        setLoading(false);
      }
    }

    fetchStats();
  }, []);

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="glass-card p-6">
            <Skeleton className="h-4 w-24 mb-4" />
            <Skeleton className="h-10 w-20 mb-2" />
            <Skeleton className="h-4 w-16" />
          </div>
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="glass-card p-6 text-center">
        <p className="text-error-red">{error}</p>
        <button
          onClick={() => window.location.reload()}
          className="mt-4 text-violet-soft hover:underline"
        >
          Retry
        </button>
      </div>
    );
  }

  if (!stats) return null;

  const formatNumber = (num: number) => {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
    return num.toString();
  };

  const cards = [
    {
      title: 'Total Followers',
      value: formatNumber(stats.totalFollowers),
      change: stats.followerChange,
      changeLabel: 'this week',
      icon: Users,
      color: 'text-violet-soft',
      changeColor: stats.followerChange >= 0 ? 'text-success-green' : 'text-error-red',
    },
    {
      title: 'Recent Unfollowers',
      value: stats.recentUnfollowers.toString(),
      change: stats.unfollowerChange,
      changeLabel: 'from last week',
      icon: UserMinus,
      color: 'text-error-red',
      changeColor: stats.unfollowerChange > 0 ? 'text-error-red' : 'text-success-green',
    },
    {
      title: 'Engagement Rate',
      value: stats.engagementRate.toFixed(1) + '%',
      change: 0.3,
      changeLabel: 'improvement',
      icon: Activity,
      color: 'text-info-blue',
      changeColor: 'text-success-green',
    },
    {
      title: 'Account Health',
      value: Math.round(stats.healthScore) + '/100',
      subtitle: stats.healthMessage,
      icon: Heart,
      color: stats.healthScore >= 80 ? 'text-success-green' : stats.healthScore >= 60 ? 'text-caution-amber' : 'text-error-red',
      showChange: false,
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      {cards.map((card, index) => (
        <div
          key={index}
          className="glass-card glass-card-hover p-6 transition-all duration-200"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-medium text-text-secondary">{card.title}</h3>
            <card.icon className={`w-5 h-5 ${card.color}`} />
          </div>

          <p className={`text-3xl font-bold ${card.color}`}>{card.value}</p>

          {card.showChange !== false && card.change !== undefined && (
            <p className={`text-sm mt-2 flex items-center gap-1 ${card.changeColor}`}>
              <span>{card.change >= 0 ? '+' : ''}{card.change}</span>
              <span className="text-text-secondary">{card.changeLabel}</span>
            </p>
          )}

          {card.subtitle && (
            <p className="text-sm text-text-secondary mt-2">{card.subtitle}</p>
          )}
        </div>
      ))}
    </div>
  );
}
