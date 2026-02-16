import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Plus, Loader2 } from 'lucide-react';
import { createAccount } from '@/lib/api/accounts';

interface AddAccountDialogProps {
  onAccountAdded: () => void;
}

export function AddAccountDialog({ onAccountAdded }: AddAccountDialogProps) {
  const [open, setOpen] = useState(false);
  const [username, setUsername] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!username.trim()) {
      setError('Username is required');
      return;
    }

    // Remove @ if user included it
    const cleanUsername = username.trim().replace(/^@/, '');

    setIsLoading(true);
    try {
      await createAccount(cleanUsername, 'instagram');
      setUsername('');
      setOpen(false);
      toast.success("Account added", {
        description: `@${cleanUsername} is now being tracked.`,
      });
      onAccountAdded();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to add account');
      toast.error("Failed to add account", {
        description: err.response?.data?.message || 'Please try again.',
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="bg-violet-soft hover:bg-violet-600">
          <Plus className="w-4 h-4 mr-2" />
          Add Account
        </Button>
      </DialogTrigger>
      <DialogContent className="glass-card border-glass-border">
        <DialogHeader>
          <DialogTitle>Add Instagram Account</DialogTitle>
          <DialogDescription className="text-text-secondary">
            Connect an Instagram account to track unfollowers.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="space-y-4 py-4">
            {error && (
              <div className="p-3 text-sm text-error-red bg-error-red/10 rounded-lg">
                {error}
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="username">Instagram Username</Label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary">
                  @
                </span>
                <Input
                  id="username"
                  placeholder="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="pl-8"
                  disabled={isLoading}
                />
              </div>
              <p className="text-xs text-text-secondary">
                Enter the Instagram username without the @ symbol
              </p>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isLoading || !username.trim()}
              className="bg-violet-soft hover:bg-violet-600"
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Adding...
                </>
              ) : (
                'Add Account'
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
