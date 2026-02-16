import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { listAccounts, deleteAccount } from '@/lib/api/accounts';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Skeleton } from '@/components/ui/skeleton';
import { Badge } from '@/components/ui/badge';
import { AddAccountDialog } from './AddAccountDialog';
import { Instagram, Trash2, Loader2, UserPlus } from 'lucide-react';
import type { Account } from '@/types/api';

export function AccountList() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleteLoading, setDeleteLoading] = useState<string | null>(null);
  const [accountToDelete, setAccountToDelete] = useState<Account | null>(null);

  const fetchAccounts = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await listAccounts();
      setAccounts(response.accounts);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load accounts');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAccounts();
  }, []);

  const handleDelete = async () => {
    if (!accountToDelete) return;

    setDeleteLoading(accountToDelete.id);
    try {
      await deleteAccount(accountToDelete.id);
      setAccountToDelete(null);
      toast.success("Account removed", {
        description: `@${accountToDelete.username} is no longer being tracked.`,
      });
      await fetchAccounts();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete account');
      toast.error("Failed to remove account", {
        description: err.response?.data?.message || 'Please try again.',
      });
    } finally {
      setDeleteLoading(null);
    }
  };

  const getStatusBadge = (status: string) => {
    const variants: Record<string, { variant: 'default' | 'secondary' | 'destructive' | 'outline'; className: string }> = {
      active: { variant: 'default', className: 'bg-success-green/20 text-success-green border-success-green/30' },
      pending: { variant: 'secondary', className: 'bg-caution-amber/20 text-caution-amber border-caution-amber/30' },
      error: { variant: 'destructive', className: 'bg-error-red/20 text-error-red border-error-red/30' },
      scanning: { variant: 'outline', className: 'bg-info-blue/20 text-info-blue border-info-blue/30' },
    };

    const config = variants[status] || variants.pending;
    return (
      <Badge variant={config.variant} className={config.className}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    );
  };

  const formatNumber = (num: number) => {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
    return num.toLocaleString();
  };

  const formatDate = (dateString: string | null) => {
    if (!dateString) return 'Never';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="glass-card p-6">
        <div className="flex justify-between items-center mb-6">
          <Skeleton className="h-6 w-32" />
          <Skeleton className="h-10 w-28" />
        </div>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <Skeleton key={i} className="h-12 w-full" />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="glass-card p-6 text-center">
        <p className="text-error-red mb-4">{error}</p>
        <Button onClick={fetchAccounts} variant="outline">
          Retry
        </Button>
      </div>
    );
  }

  return (
    <>
      <div className="glass-card">
        <div className="p-6 border-b border-glass-border flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-xl font-semibold">Connected Accounts</h2>
            <p className="text-sm text-text-secondary mt-1">
              {accounts.length} account{accounts.length !== 1 ? 's' : ''} connected
            </p>
          </div>
          <AddAccountDialog onAccountAdded={fetchAccounts} />
        </div>

        {accounts.length === 0 ? (
          <div className="p-12 text-center">
            <div className="w-16 h-16 bg-violet-soft/10 rounded-full flex items-center justify-center mx-auto mb-4">
              <UserPlus className="w-8 h-8 text-violet-soft" />
            </div>
            <h3 className="text-lg font-medium mb-2">No accounts connected</h3>
            <p className="text-text-secondary mb-6 max-w-sm mx-auto">
              Add your first Instagram account to start tracking unfollowers and analytics.
            </p>
            <AddAccountDialog onAccountAdded={fetchAccounts} />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-glass-border hover:bg-transparent">
                  <TableHead className="text-text-secondary">Username</TableHead>
                  <TableHead className="text-text-secondary">Platform</TableHead>
                  <TableHead className="text-text-secondary">Followers</TableHead>
                  <TableHead className="text-text-secondary">Status</TableHead>
                  <TableHead className="text-text-secondary">Last Scan</TableHead>
                  <TableHead className="text-text-secondary text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {accounts.map((account) => (
                  <TableRow
                    key={account.id}
                    className="border-glass-border hover:bg-white/5"
                  >
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        <Instagram className="w-4 h-4 text-violet-soft" />
                        @{account.username}
                      </div>
                    </TableCell>
                    <TableCell className="capitalize">{account.platform}</TableCell>
                    <TableCell>{formatNumber(account.follower_count || 0)}</TableCell>
                    <TableCell>{getStatusBadge(account.scan_status)}</TableCell>
                    <TableCell className="text-text-secondary">
                      {formatDate(account.last_scan_at)}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="text-error-red hover:text-error-red hover:bg-error-red/10"
                        onClick={() => setAccountToDelete(account)}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </div>

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!accountToDelete} onOpenChange={() => setAccountToDelete(null)}>
        <DialogContent className="glass-card border-glass-border">
          <DialogHeader>
            <DialogTitle>Delete Account</DialogTitle>
            <DialogDescription className="text-text-secondary">
              Are you sure you want to delete @{accountToDelete?.username}? This action
              cannot be undone and all tracking data will be lost.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2">
            <Button
              variant="outline"
              onClick={() => setAccountToDelete(null)}
              disabled={!!deleteLoading}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={!!deleteLoading}
              className="bg-error-red hover:bg-error-red/90"
            >
              {deleteLoading === accountToDelete?.id ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Deleting...
                </>
              ) : (
                'Delete Account'
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
