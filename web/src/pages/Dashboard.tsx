import { StatsCards } from '@/components/dashboard/StatsCards';
import { AccountList } from '@/components/dashboard/AccountList';
import { UnfollowsList } from '@/components/dashboard/UnfollowsList';

export function Dashboard() {
  return (
    <div>
      <h1 className="text-3xl font-bold mb-2">Dashboard</h1>
      <p className="text-text-secondary mb-8">
        Track your Instagram growth and unfollower activity
      </p>

      <StatsCards />

      <div className="mt-8 grid grid-cols-1 xl:grid-cols-2 gap-8">
        <AccountList />
        <UnfollowsList />
      </div>
    </div>
  );
}
