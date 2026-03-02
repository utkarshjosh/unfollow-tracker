import { Outlet, Link, useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { SiteLogo } from '@/components/branding/SiteLogo';
import { LayoutDashboard, User, Settings, LogOut, Menu } from 'lucide-react';
import { useState } from 'react';

export function DashboardLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const handleLogout = () => {
    logout();
    toast.info("Logged out", {
      description: "You have been signed out successfully.",
    });
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-surface-deep flex">
      {/* Sidebar */}
      <aside
        className={`fixed lg:static inset-y-0 left-0 z-50 w-64 bg-indigo-deep border-r border-glass-border transform transition-transform duration-200 ease-in-out lg:transform-none ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        <div className="p-6">
          <Link to="/" className="flex items-center gap-2">
            <SiteLogo
              iconClassName="h-9 w-9 rounded-lg"
              textClassName="text-2xl font-bold text-violet-soft"
            />
          </Link>
        </div>

        <nav className="px-4 py-4 space-y-2">
          <Link
            to="/dashboard"
            className="flex items-center gap-3 px-4 py-3 rounded-lg text-text-primary hover:bg-white/5 transition-colors"
          >
            <LayoutDashboard className="w-5 h-5" />
            <span>Dashboard</span>
          </Link>
          <div className="flex items-center gap-3 px-4 py-3 rounded-lg text-text-secondary cursor-not-allowed">
            <User className="w-5 h-5" />
            <span>Accounts</span>
          </div>
          <div className="flex items-center gap-3 px-4 py-3 rounded-lg text-text-secondary cursor-not-allowed">
            <Settings className="w-5 h-5" />
            <span>Settings</span>
          </div>
        </nav>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <header className="h-16 border-b border-glass-border flex items-center justify-between px-4 lg:px-8">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="lg:hidden p-2 text-text-secondary hover:text-text-primary"
          >
            <Menu className="w-6 h-6" />
          </button>

          <div className="flex-1" />

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="text-text-primary">
                {user?.email}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <DropdownMenuItem onClick={handleLogout}>
                <LogOut className="w-4 h-4 mr-2" />
                Logout
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </header>

        {/* Page content */}
        <main className="flex-1 p-4 lg:p-8 overflow-auto">
          <Outlet />
        </main>
      </div>

      {/* Mobile overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
    </div>
  );
}
