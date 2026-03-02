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
  Shield,
  Lock,
  Eye,
  EyeOff,
  Sparkles,
  CheckCircle,
  Loader2,
  ArrowRight,
} from 'lucide-react';

export function Login() {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const response = await api.post<LoginResponse>('/api/v1/auth/login', {
        email,
        password,
      });
      login(response.data);
      toast.success('Welcome back!', {
        description: 'You have successfully signed in.',
      });
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Invalid credentials');
    } finally {
      setIsLoading(false);
    }
  };

  const trustFeatures = [
    { icon: Shield, text: 'End-to-end encryption' },
    { icon: Lock, text: 'Your data never leaves your device' },
    { icon: CheckCircle, text: 'GDPR compliant' },
  ];

  return (
    <div className="min-h-screen w-full bg-surface-deep flex">
      {/* Animated gradient background */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-violet-soft/20 rounded-full blur-3xl animate-pulse" />
        <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-indigo-deep/40 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '1s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-violet-soft/10 rounded-full blur-3xl" />
      </div>

      {/* Left side - Brand messaging */}
      <div className="hidden lg:flex lg:w-1/2 xl:w-3/5 relative flex-col justify-between p-12 xl:p-16">
        {/* Logo area */}
        <div className="relative z-10">
          <Link to="/" className="group inline-flex">
            <SiteLogo
              iconClassName="w-10 h-10 rounded-xl shadow-lg shadow-violet-soft/25 transition-shadow duration-300 group-hover:shadow-violet-soft/40"
            />
          </Link>
        </div>

        {/* Main value proposition */}
        <div className="relative z-10 max-w-lg">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-violet-soft/10 border border-violet-soft/20 mb-6">
            <Sparkles className="w-4 h-4 text-violet-soft" />
            <span className="text-sm font-medium text-violet-soft">
              Privacy-first insights
            </span>
          </div>

          <h1 className="text-4xl xl:text-5xl font-heading font-bold text-text-primary leading-tight mb-6">
            Understand your{' '}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-violet-soft to-purple-400">
              connections
            </span>{' '}
            with confidence
          </h1>

          <p className="text-lg text-text-secondary leading-relaxed mb-8">
            Track your Instagram unfollowers with a tool that respects your
            privacy. No data mining, no third-party sharing—just clear insights
            that help you nurture meaningful relationships.
          </p>

          {/* Trust signals */}
          <div className="space-y-4">
            {trustFeatures.map((feature, index) => (
              <div
                key={index}
                className="flex items-center gap-3 text-text-secondary"
              >
                <div className="w-8 h-8 rounded-lg bg-violet-soft/10 flex items-center justify-center">
                  <feature.icon className="w-4 h-4 text-violet-soft" />
                </div>
                <span className="text-sm font-medium">{feature.text}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Bottom quote */}
        <div className="relative z-10">
          <blockquote className="text-text-secondary/80 text-sm">
            <p className="mb-2">
              "The Caregiver archetype guides our mission: protecting your data
              while helping you understand your social landscape."
            </p>
            <footer className="text-text-secondary/60">
              — Built with intention
            </footer>
          </blockquote>
        </div>
      </div>

      {/* Right side - Login form */}
      <div className="w-full lg:w-1/2 xl:w-2/5 flex items-center justify-center p-4 sm:p-6 lg:p-8 relative z-10">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="mb-8 flex justify-center lg:hidden">
            <Link to="/">
              <SiteLogo iconClassName="w-10 h-10 rounded-xl shadow-lg shadow-violet-soft/25" />
            </Link>
          </div>

          {/* Glass card */}
          <div className="glass-card p-6 sm:p-8">
            {/* Header */}
            <div className="text-center mb-8">
              <h2 className="text-2xl font-heading font-bold text-text-primary mb-2">
                Welcome back
              </h2>
              <p className="text-text-secondary text-sm">
                Sign in to continue your journey
              </p>
            </div>

            {/* Error message */}
            {error && (
              <div className="mb-6 p-4 rounded-xl bg-error-red/10 border border-error-red/20 flex items-start gap-3">
                <div className="w-5 h-5 rounded-full bg-error-red/20 flex items-center justify-center flex-shrink-0 mt-0.5">
                  <span className="text-error-red text-xs font-bold">!</span>
                </div>
                <p className="text-sm text-error-red">{error}</p>
              </div>
            )}

            {/* Form */}
            <form onSubmit={handleSubmit} className="space-y-5">
              {/* Email field */}
              <div className="space-y-2">
                <Label
                  htmlFor="email"
                  className="text-sm font-medium text-text-primary"
                >
                  Email address
                </Label>
                <div className="relative">
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
              </div>

              {/* Password field */}
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label
                    htmlFor="password"
                    className="text-sm font-medium text-text-primary"
                  >
                    Password
                  </Label>
                  <Link
                    to="/forgot-password"
                    className="text-xs text-violet-soft hover:text-violet-400 transition-colors"
                  >
                    Forgot password?
                  </Link>
                </div>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    placeholder="Enter your password"
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
                    {showPassword ? (
                      <EyeOff className="w-5 h-5" />
                    ) : (
                      <Eye className="w-5 h-5" />
                    )}
                  </button>
                </div>
              </div>

              {/* Submit button */}
              <Button
                type="submit"
                disabled={isLoading}
                className="w-full h-12 bg-gradient-to-r from-violet-soft to-violet-600 hover:from-violet-600 hover:to-violet-700 text-white font-medium rounded-xl shadow-lg shadow-violet-soft/25 hover:shadow-violet-soft/40 transition-all duration-300 disabled:opacity-70 disabled:cursor-not-allowed group"
              >
                {isLoading ? (
                  <>
                    <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                    Signing in...
                  </>
                ) : (
                  <>
                    Sign In
                    <ArrowRight className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
                  </>
                )}
              </Button>
            </form>

            {/* Divider */}
            <div className="relative my-6">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-glass-border" />
              </div>
              <div className="relative flex justify-center text-xs">
                <span className="px-4 bg-[rgba(30,27,75,0.6)] text-text-secondary">
                  Protected by industry-standard security
                </span>
              </div>
            </div>

            {/* Sign up link */}
            <p className="text-center text-sm text-text-secondary">
              Don't have an account?{' '}
              <Link
                to="/register"
                className="font-medium text-violet-soft hover:text-violet-400 transition-colors inline-flex items-center gap-1"
              >
                Create one
                <ArrowRight className="w-3 h-3" />
              </Link>
            </p>

            {/* Privacy badge */}
            <div className="mt-6 flex items-center justify-center gap-2 text-xs text-text-secondary/70">
              <Lock className="w-3 h-3" />
              <span>Your data is encrypted and secure</span>
            </div>
          </div>

          {/* Mobile trust features */}
          <div className="lg:hidden mt-8 grid grid-cols-3 gap-4">
            {trustFeatures.map((feature, index) => (
              <div key={index} className="text-center">
                <div className="w-10 h-10 mx-auto rounded-lg bg-violet-soft/10 flex items-center justify-center mb-2">
                  <feature.icon className="w-5 h-5 text-violet-soft" />
                </div>
                <span className="text-xs text-text-secondary/70">
                  {feature.text}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
