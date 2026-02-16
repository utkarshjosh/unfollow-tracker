import { useEffect, useState } from 'react';
import { getUnfollows } from '@/lib/api/unfollows';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Skeleton } from '@/components/ui/skeleton';
import { Badge } from '@/components/ui/badge';
import { ChevronLeft, ChevronRight, UserX, Calendar } from 'lucide-react';
import type { Unfollow } from '@/types/api';

const ITEMS_PER_PAGE = 10;

export function UnfollowsList() {
  const [unfollows, setUnfollows] = useState<Unfollow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [total, setTotal] = useState(0);

  const fetchUnfollows = async (newOffset: number) => {
    try {
      setLoading(true);
      setError(null);
      const response = await getUnfollows(undefined, ITEMS_PER_PAGE, newOffset);
      setUnfollows(response.unfollows);
      setTotal(response.total);
      setOffset(response.offset);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to load unfollows');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUnfollows(0);
  }, []);

  const handlePrevious = () => {
    const newOffset = Math.max(0, offset - ITEMS_PER_PAGE);
    fetchUnfollows(newOffset);
  };

  const handleNext = () => {
    const newOffset = offset + ITEMS_PER_PAGE;
    if (newOffset < total) {
      fetchUnfollows(newOffset);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffDays = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;

    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const totalPages = Math.ceil(total / ITEMS_PER_PAGE);
  const currentPage = Math.floor(offset / ITEMS_PER_PAGE) + 1;

  if (loading) {
    return (
      <div className="glass-card">
        <div className="p-6 border-b border-glass-border">
          <Skeleton className="h-6 w-40" />
        </div>
        <div className="p-6 space-y-4">
          {[...Array(5)].map((_, i) => (
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
        <Button onClick={() => fetchUnfollows(0)} variant="outline">
          Retry
        </Button>
      </div>
    );
  }

  return (
    <div className="glass-card">
      <div className="p-6 border-b border-glass-border">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-semibold">Recent Unfollowers</h2>
            <p className="text-sm text-text-secondary mt-1">
              {total} unfollower{total !== 1 ? 's' : ''} detected
            </p>
          </div>
          {total > 0 && (
            <Badge variant="outline" className="bg-error-red/10 text-error-red border-error-red/30">
              <UserX className="w-3 h-3 mr-1" />
              {total} total
            </Badge>
          )}
        </div>
      </div>

      {unfollows.length === 0 ? (
        <div className="p-12 text-center">
          <div className="w-16 h-16 bg-success-green/10 rounded-full flex items-center justify-center mx-auto mb-4">
            <Calendar className="w-8 h-8 text-success-green" />
          </div>
          <h3 className="text-lg font-medium mb-2">No unfollowers yet</h3>
          <p className="text-text-secondary max-w-sm mx-auto">
            Great news! No unfollowers have been detected for your accounts. We'll keep monitoring and notify you of any changes.
          </p>
        </div>
      ) : (
        <>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-glass-border hover:bg-transparent">
                  <TableHead className="text-text-secondary">Detected</TableHead>
                  <TableHead className="text-text-secondary">Account ID</TableHead>
                  <TableHead className="text-text-secondary">Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {unfollows.map((unfollow) => (
                  <TableRow
                    key={unfollow.id}
                    className="border-glass-border hover:bg-white/5"
                  >
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Calendar className="w-4 h-4 text-text-secondary" />
                        <span>{formatDate(unfollow.detected_date)}</span>
                      </div>
                    </TableCell>
                    <TableCell className="font-mono text-sm text-text-secondary">
                      {unfollow.account_id.slice(0, 8)}...
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="outline"
                        className="bg-error-red/10 text-error-red border-error-red/30"
                      >
                        <UserX className="w-3 h-3 mr-1" />
                        Unfollowed
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="p-4 border-t border-glass-border flex items-center justify-between">
              <p className="text-sm text-text-secondary">
                Page {currentPage} of {totalPages}
              </p>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handlePrevious}
                  disabled={offset === 0 || loading}
                >
                  <ChevronLeft className="w-4 h-4 mr-1" />
                  Previous
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleNext}
                  disabled={offset + ITEMS_PER_PAGE >= total || loading}
                >
                  Next
                  <ChevronRight className="w-4 h-4 ml-1" />
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
