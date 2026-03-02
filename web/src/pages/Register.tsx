import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast } from 'sonner';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { api } from '@/lib/api';
import type { LoginResponse } from '@/types/api';
import { SiteLogo } from '@/components/branding/SiteLogo';
import {
  ArrowRight,
  CheckCircle2,
  Eye,
  EyeOff,
  Loader2,
  Lock,
  Shield,
  Sparkles,
} from 'lucide-react';

export function Register() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const trustFeatures = [
    { icon: Shield, text: 'Privacy-first by design' },
    { icon: Lock, text: 'Encrypted account security' },
    { icon: CheckCircle2, text: 'No anxiety-based nudges' },
  ];

  const meetsLength = password.length >= 8;
  const matchesConfirm = confirmPassword.length > 0 && password === confirmPassword;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setIsLoading(true);

    try {
      await api.post('/api/v1/auth/register', {
        email,
        password,
      });

      // Auto-login after registration
      const loginResponse = await api.post<LoginResponse>('/api/v1/auth/login', {
        email,
        password,
      });
      login(loginResponse.data);
      toast.success("Account created!", {
        description: "Welcome to Unfollow Tracker.",
      });
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Registration failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen w-full bg-surface-deep flex">
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-16 left-1/3 w-[460px] h-[460px] bg-violet-soft/20 rounded-full blur-3xl animate-pulse" />
        <div
          className="absolute bottom-0 right-1/4 w-[380px] h-[380px] bg-indigo-deep/40 rounded-full blur-3xl animate-pulse"
          style={{ animationDelay: '1s' }}
        />
      </div>

      <div className="hidden lg:flex lg:w-1/2 xl:w-3/5 relative flex-col justify-between p-12 xl:p-16">
        <div className="relative z-10">
          <Link to="/" className="group inline-flex">
            <SiteLogo
              iconClassName="w-10 h-10 rounded-xl shadow-lg shadow-violet-soft/25 transition-shadow duration-300 group-hover:shadow-violet-soft/40"
            />
          </Link>
        </div>

        <div className="relative z-10 max-w-lg">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-violet-soft/10 border border-violet-soft/20 mb-6">
            <Sparkles className="w-4 h-4 text-violet-soft" />
            <span className="text-sm font-medium text-violet-soft">
              Insights, not obsessions
            </span>
          </div>

          <h1 className="text-4xl xl:text-5xl font-heading font-bold text-text-primary leading-tight mb-6">
            Build healthier growth habits from day one
          </h1>

          <p className="text-lg text-text-secondary leading-relaxed mb-8">
            Create your account and start tracking patterns over time, with
            transparent analytics that protect your privacy and peace of mind.
          </p>

          <div className="space-y-4">
            {trustFeatures.map((feature, index) => (
              <div key={index} className="flex items-center gap-3 text-text-secondary">
                <div className="w-8 h-8 rounded-lg bg-violet-soft/10 flex items-center justify-center">
                  <feature.icon className="w-4 h-4 text-violet-soft" />
                </div>
                <span className="text-sm font-medium">{feature.text}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="relative z-10">
          <blockquote className="text-text-secondary/80 text-sm">
            <p>
              "Understand your growth without compromising your privacy or
              mental well-being."
            </p>
          </blockquote>
        </div>
      </div>

      <div className="w-full lg:w-1/2 xl:w-2/5 flex items-center justify-center p-4 sm:p-6 lg:p-8 relative z-10">
        <div className="w-full max-w-md">
          <div className="mb-8 flex justify-center lg:hidden">
            <Link to="/">
              <SiteLogo iconClassName="w-10 h-10 rounded-xl shadow-lg shadow-violet-soft/25" />
            </Link>
          </div>

          <div className="glass-card p-6 sm:p-8">
            <div className="text-center mb-8">
              <h2 className="text-2xl font-heading font-bold text-text-primary mb-2">
                Create your account
              </h2>
              <p className="text-text-secondary text-sm">
                Start with privacy-first analytics in under a minute
              </p>
            </div>

            {error && (
              <div className="mb-6 p-4 rounded-xl bg-error-red/10 border border-error-red/20 text-sm text-error-red">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-5">
              <div className="space-y-2">
                <Label htmlFor="email" className="text-sm font-medium text-text-primary">
                  Email address
                </Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="you@example.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="h-12 bg-surface-deep/50 border-glass-border text-text-primary placeholder:text-text-secondary/50 focus:border-violet-soft focus:ring-violet-soft/20 rounded-xl transition-all duration-200"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password" className="text-sm font-medium text-text-primary">
                  Password
                </Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="h-12 bg-surface-deep/50 border-glass-border text-text-primary placeholder:text-text-secondary/50 focus:border-violet-soft focus:ring-violet-soft/20 rounded-xl pr-12 transition-all duration-200"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-text-secondary hover:text-text-primary transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-soft rounded"
                    aria-label={showPassword ? 'Hide password' : 'Show password'}
                  >
                    {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                  </button>
                </div>
              </div>

              <div className="space-y-2">
                <Label
                  htmlFor="confirmPassword"
                  className="text-sm font-medium text-text-primary"
                >
                  Confirm password
                </Label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    type={showConfirmPassword ? 'text' : 'password'}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                    className="h-12 bg-surface-deep/50 border-glass-border text-text-primary placeholder:text-text-secondary/50 focus:border-violet-soft focus:ring-violet-soft/20 rounded-xl pr-12 transition-all duration-200"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-text-secondary hover:text-text-primary transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-soft rounded"
                    aria-label={showConfirmPassword ? 'Hide confirm password' : 'Show confirm password'}
                  >
                    {showConfirmPassword ? (
                      <EyeOff className="w-5 h-5" />
                    ) : (
                      <Eye className="w-5 h-5" />
                    )}
                  </button>
                </div>
              </div>

              <div className="space-y-2 text-xs">
                <div
                  className={`flex items-center gap-2 ${
                    meetsLength ? 'text-success-green' : 'text-text-secondary'
                  }`}
                >
                  <CheckCircle2 className="w-3.5 h-3.5" />
                  <span>At least 8 characters</span>
                </div>
                <div
                  className={`flex items-center gap-2 ${
                    matchesConfirm ? 'text-success-green' : 'text-text-secondary'
                  }`}
                >
                  <CheckCircle2 className="w-3.5 h-3.5" />
                  <span>Passwords match</span>
                </div>
              </div>

              <Button
                type="submit"
                className="w-full h-12 bg-gradient-to-r from-violet-soft to-violet-600 hover:from-violet-600 hover:to-violet-700 text-white font-medium rounded-xl shadow-lg shadow-violet-soft/25 hover:shadow-violet-soft/40 transition-all duration-300 disabled:opacity-70 disabled:cursor-not-allowed group"
                disabled={isLoading}
              >
                {isLoading ? (
                  <>
                    <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                    Creating account...
                  </>
                ) : (
                  <>
                    Create Account
                    <ArrowRight className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
                  </>
                )}
              </Button>
            </form>

            <p className="mt-6 text-center text-sm text-text-secondary">
              Already have an account?{' '}
              <Link
                to="/login"
                className="font-medium text-violet-soft hover:text-violet-400 transition-colors inline-flex items-center gap-1"
              >
                Sign in
                <ArrowRight className="w-3 h-3" />
              </Link>
            </p>
          </div>

          <div className="lg:hidden mt-8 grid grid-cols-3 gap-4">
            {trustFeatures.map((feature, index) => (
              <div key={index} className="text-center">
                <div className="w-10 h-10 mx-auto rounded-lg bg-violet-soft/10 flex items-center justify-center mb-2">
                  <feature.icon className="w-5 h-5 text-violet-soft" />
                </div>
                <span className="text-xs text-text-secondary/70">{feature.text}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
