import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { SiteLogo } from '@/components/branding/SiteLogo';
import {
  Shield,
  Users,
  TrendingUp,
  Heart,
  Lock,
  EyeOff,
  CheckCircle2,
  ArrowRight,
  Instagram,
  BarChart3,
  Sparkles,
  MessageCircle,
  ChevronDown,
} from 'lucide-react';

// Animated gradient orb component
function GradientOrb({ className, style }: { className?: string; style?: React.CSSProperties }) {
  return (
    <div
      className={`absolute rounded-full blur-[100px] opacity-60 animate-pulse-slow ${className}`}
      style={{
        background: 'linear-gradient(135deg, #7C3AED 0%, #1E1B4B 50%, #0F0D1A 100%)',
        animationDuration: '8s',
        ...style,
      }}
    />
  );
}

// Glass card component
function GlassCard({ children, className = '', style }: { children: React.ReactNode; className?: string; style?: React.CSSProperties }) {
  return (
    <div
      className={`relative overflow-hidden rounded-2xl border border-white/10 bg-white/[0.03] backdrop-blur-xl transition-all duration-500 hover:bg-white/[0.06] hover:border-white/20 ${className}`}
      style={style}
    >
      {children}
    </div>
  );
}

// Feature card component
function FeatureCard({
  icon: Icon,
  title,
  description,
  delay,
}: {
  icon: React.ElementType;
  title: string;
  description: string;
  delay: number;
}) {
  return (
    <GlassCard
      className="group p-8"
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="mb-6 inline-flex h-14 w-14 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500/20 to-violet-600/5 text-violet-400 transition-transform duration-500 group-hover:scale-110">
        <Icon className="h-7 w-7" />
      </div>
      <h3 className="mb-3 text-xl font-semibold tracking-tight text-white">{title}</h3>
      <p className="leading-relaxed text-slate-400">{description}</p>
    </GlassCard>
  );
}

// Step card for How It Works
function StepCard({
  number,
  title,
  description,
  icon: Icon,
}: {
  number: string;
  title: string;
  description: string;
  icon: React.ElementType;
}) {
  return (
    <div className="relative flex flex-col items-center text-center">
      <div className="mb-6 flex h-20 w-20 items-center justify-center rounded-2xl bg-gradient-to-br from-violet-600 to-indigo-700 text-white shadow-lg shadow-violet-500/25">
        <Icon className="h-10 w-10" />
      </div>
      <div className="absolute -top-2 -right-2 flex h-8 w-8 items-center justify-center rounded-full bg-white/10 text-sm font-bold text-white backdrop-blur-md">
        {number}
      </div>
      <h3 className="mb-2 text-lg font-semibold text-white">{title}</h3>
      <p className="max-w-xs text-slate-400">{description}</p>
    </div>
  );
}

// Testimonial card
function TestimonialCard({
  quote,
  author,
  role,
  avatar,
}: {
  quote: string;
  author: string;
  role: string;
  avatar: string;
}) {
  return (
    <GlassCard className="p-8">
      <div className="mb-4 flex gap-1">
        {[...Array(5)].map((_, i) => (
          <Sparkles key={i} className="h-5 w-5 fill-violet-400 text-violet-400" />
        ))}
      </div>
      <p className="mb-6 text-lg leading-relaxed text-slate-300">&ldquo;{quote}&rdquo;</p>
      <div className="flex items-center gap-4">
        <img
          src={avatar}
          alt={author}
          className="h-12 w-12 rounded-full object-cover ring-2 ring-white/10"
        />
        <div>
          <p className="font-semibold text-white">{author}</p>
          <p className="text-sm text-slate-400">{role}</p>
        </div>
      </div>
    </GlassCard>
  );
}

// Trust badge component
function TrustBadge({ icon: Icon, label }: { icon: React.ElementType; label: string }) {
  return (
    <div className="flex items-center gap-3 rounded-full border border-white/10 bg-white/5 px-5 py-2.5 backdrop-blur-sm">
      <Icon className="h-4 w-4 text-violet-400" />
      <span className="text-sm font-medium text-slate-300">{label}</span>
    </div>
  );
}

// Dashboard preview mockup
function DashboardPreview() {
  return (
    <div className="relative mx-auto max-w-4xl">
      {/* Glow effect behind */}
      <div className="absolute inset-0 -m-4 rounded-3xl bg-gradient-to-r from-violet-600/30 via-indigo-600/30 to-violet-600/30 blur-2xl" />

      {/* Browser chrome */}
      <div className="relative overflow-hidden rounded-2xl border border-white/10 bg-[#0a0814] shadow-2xl">
        {/* Browser header */}
        <div className="flex items-center gap-2 border-b border-white/5 bg-white/5 px-4 py-3">
          <div className="flex gap-1.5">
            <div className="h-3 w-3 rounded-full bg-red-500/80" />
            <div className="h-3 w-3 rounded-full bg-amber-500/80" />
            <div className="h-3 w-3 rounded-full bg-green-500/80" />
          </div>
          <div className="ml-4 flex flex-1 items-center justify-center">
            <div className="flex items-center gap-2 rounded-lg bg-white/5 px-4 py-1.5 text-xs text-slate-400">
              <Lock className="h-3 w-3" />
              app.unfollowtracker.io/dashboard
            </div>
          </div>
        </div>

        {/* Dashboard content */}
        <div className="p-6">
          {/* Stats row */}
          <div className="mb-6 grid grid-cols-3 gap-4">
            {[
              { label: 'Total Followers', value: '12.4K', change: '+2.3%' },
              { label: 'Net Growth', value: '+284', change: '+5.1%' },
              { label: 'Unfollows', value: '12', change: '-8.2%' },
            ].map((stat, i) => (
              <div key={i} className="rounded-xl border border-white/5 bg-white/[0.03] p-4">
                <p className="mb-1 text-xs text-slate-500">{stat.label}</p>
                <div className="flex items-end justify-between">
                  <span className="text-xl font-bold text-white">{stat.value}</span>
                  <span className={`text-xs ${stat.change.startsWith('+') ? 'text-emerald-400' : 'text-rose-400'}`}>
                    {stat.change}
                  </span>
                </div>
              </div>
            ))}
          </div>

          {/* Chart area */}
          <div className="mb-6 rounded-xl border border-white/5 bg-white/[0.03] p-4">
            <div className="mb-4 flex items-center justify-between">
              <span className="text-sm font-medium text-slate-300">30-Day Growth</span>
              <div className="flex gap-2">
                <div className="h-2 w-8 rounded-full bg-violet-500" />
                <div className="h-2 w-8 rounded-full bg-white/20" />
              </div>
            </div>
            {/* Chart visualization */}
            <div className="flex items-end gap-2 h-24">
              {[40, 55, 45, 60, 58, 72, 68, 85, 80, 95, 88, 100].map((h, i) => (
                <div
                  key={i}
                  className="flex-1 rounded-t bg-gradient-to-t from-violet-600 to-violet-400 transition-all duration-500 hover:from-violet-500 hover:to-violet-300"
                  style={{ height: `${h}%`, opacity: 0.6 + (i * 0.03) }}
                />
              ))}
            </div>
          </div>

          {/* Recent unfollows list */}
          <div className="rounded-xl border border-white/5 bg-white/[0.03] p-4">
            <span className="mb-3 block text-sm font-medium text-slate-300">Recent Activity</span>
            <div className="space-y-2">
              {[
                { name: 'sarah_creates', action: 'unfollowed', time: '2h ago', avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=64&h=64&fit=crop&crop=face' },
                { name: 'tech_daily', action: 'followed', time: '5h ago', avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=64&h=64&fit=crop&crop=face' },
                { name: 'design_hub', action: 'unfollowed', time: '1d ago', avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=64&h=64&fit=crop&crop=face' },
              ].map((item, i) => (
                <div key={i} className="flex items-center justify-between rounded-lg bg-white/5 px-3 py-2">
                  <div className="flex items-center gap-3">
                    <img
                      src={item.avatar}
                      alt={item.name}
                      className="h-8 w-8 rounded-full object-cover ring-1 ring-white/10"
                    />
                    <span className="text-sm text-slate-300">@{item.name}</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className={`text-xs ${item.action === 'unfollowed' ? 'text-rose-400' : 'text-emerald-400'}`}>
                      {item.action}
                    </span>
                    <span className="text-xs text-slate-500">{item.time}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export function Landing() {
  const { isAuthenticated } = useAuth();

  return (
    <div className="relative min-h-screen overflow-hidden bg-[#05040a]">
      {/* Animated background orbs */}
      <GradientOrb className="-left-40 -top-40 h-[600px] w-[600px]" />
      <GradientOrb className="right-0 top-1/3 h-[500px] w-[500px]" style={{ animationDelay: '2s' }} />
      <GradientOrb className="bottom-0 left-1/3 h-[700px] w-[700px]" style={{ animationDelay: '4s' }} />

      {/* Noise texture overlay */}
      <div
        className="pointer-events-none fixed inset-0 opacity-[0.015]"
        style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E")`,
        }}
      />

      {/* Header */}
      <header className="fixed left-0 right-0 top-0 z-50 border-b border-white/5 bg-[#05040a]/80 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <Link to="/" className="flex items-center gap-3">
            <SiteLogo
              iconClassName="h-10 w-10 rounded-xl shadow-lg shadow-violet-500/25"
              textClassName="text-xl font-bold tracking-tight text-white"
            />
          </Link>

          <nav className="flex items-center gap-6">
            {isAuthenticated ? (
              <Link to="/dashboard">
                <Button className="bg-violet-600 hover:bg-violet-700">Dashboard</Button>
              </Link>
            ) : (
              <>
                <Link to="/login" className="hidden text-sm font-medium text-slate-400 transition-colors hover:text-white sm:block">
                  Log in
                </Link>
                <Link to="/register">
                  <Button className="bg-violet-600 hover:bg-violet-700">Get Started</Button>
                </Link>
              </>
            )}
          </nav>
        </div>
      </header>

      <main className="relative z-10">
        {/* Hero Section */}
        <section className="relative px-6 pb-32 pt-40">
          <div className="mx-auto max-w-7xl">
            <div className="mb-20 text-center">
              {/* Badge */}
              <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-violet-500/30 bg-violet-500/10 px-4 py-1.5">
                <Shield className="h-4 w-4 text-violet-400" />
                <span className="text-sm font-medium text-violet-300">Privacy-First Analytics</span>
              </div>

              {/* Headline */}
              <h1 className="mx-auto mb-6 max-w-4xl text-5xl font-bold leading-[1.1] tracking-tight text-white sm:text-6xl lg:text-7xl">
                Understand your growth{' '}
                <span className="bg-gradient-to-r from-violet-400 to-indigo-400 bg-clip-text text-transparent">
                  without compromising
                </span>{' '}
                your privacy
              </h1>

              {/* Subheadline */}
              <p className="mx-auto mb-10 max-w-2xl text-lg leading-relaxed text-slate-400 sm:text-xl">
                The ethical Instagram analytics platform that respects your data and your mental well-being.
                Track unfollowers, analyze trends, and grow with confidence.
              </p>

              {/* CTAs */}
              <div className="flex flex-col items-center justify-center gap-4 sm:flex-row">
                <Link to="/register">
                  <Button
                    size="lg"
                    className="group h-14 bg-violet-600 px-8 text-lg font-semibold hover:bg-violet-700"
                  >
                    Get Started Free
                    <ArrowRight className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" />
                  </Button>
                </Link>
                <a href="#features">
                  <Button
                    size="lg"
                    variant="outline"
                    className="h-14 border-white/20 px-8 text-lg font-semibold text-white hover:bg-white/10"
                  >
                    Learn More
                  </Button>
                </a>
              </div>

              {/* Trust badges */}
              <div className="mt-12 flex flex-wrap items-center justify-center gap-3">
                <TrustBadge icon={Lock} label="End-to-End Encrypted" />
                <TrustBadge icon={EyeOff} label="No Data Selling" />
                <TrustBadge icon={Shield} label="GDPR Compliant" />
              </div>
            </div>

            {/* Dashboard Preview */}
            <DashboardPreview />

            {/* Scroll indicator */}
            <div className="mt-16 flex justify-center">
              <a
                href="#features"
                className="flex h-12 w-12 animate-bounce items-center justify-center rounded-full border border-white/10 text-slate-500 transition-colors hover:border-white/20 hover:text-slate-300"
              >
                <ChevronDown className="h-6 w-6" />
              </a>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section id="features" className="relative px-6 py-32">
          <div className="mx-auto max-w-7xl">
            <div className="mb-16 text-center">
              <h2 className="mb-4 text-3xl font-bold tracking-tight text-white sm:text-4xl">
                Everything you need to{' '}
                <span className="text-violet-400">grow confidently</span>
              </h2>
              <p className="mx-auto max-w-2xl text-lg text-slate-400">
                Powerful insights without the anxiety. We believe in transparency, not surveillance.
              </p>
            </div>

            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
              <FeatureCard
                icon={Shield}
                title="Privacy-First Analytics"
                description="Your data is encrypted and never sold. We collect only what's necessary and are transparent about everything."
                delay={0}
              />
              <FeatureCard
                icon={Users}
                title="Unfollower Tracking"
                description="Know who unfollowed you without obsessively checking. Get daily or weekly summaries, not instant anxiety."
                delay={100}
              />
              <FeatureCard
                icon={TrendingUp}
                title="Growth Insights"
                description="Understand your growth patterns with beautiful charts. Focus on trends, not vanity metrics."
                delay={200}
              />
              <FeatureCard
                icon={Heart}
                title="Wellness-Focused Design"
                description="Built to reduce social media anxiety. No push notifications for unfollows. Insights, not obsessions."
                delay={300}
              />
            </div>
          </div>
        </section>

        {/* How It Works */}
        <section className="relative px-6 py-32">
          <div className="mx-auto max-w-5xl">
            <div className="mb-16 text-center">
              <h2 className="mb-4 text-3xl font-bold tracking-tight text-white sm:text-4xl">
                Simple, secure, <span className="text-violet-400">transparent</span>
              </h2>
              <p className="mx-auto max-w-2xl text-lg text-slate-400">
                Get started in minutes. No complicated setup, no invasive permissions.
              </p>
            </div>

            <div className="relative">
              {/* Connection line */}
              <div className="absolute left-1/2 top-20 hidden h-1 w-2/3 -translate-x-1/2 bg-gradient-to-r from-transparent via-violet-500/30 to-transparent lg:block" />

              <div className="grid gap-12 lg:grid-cols-3">
                <StepCard
                  number="1"
                  title="Connect Account"
                  description="Securely link your Instagram account with read-only access. We never post on your behalf."
                  icon={Instagram}
                />
                <StepCard
                  number="2"
                  title="We Track Changes"
                  description="Our system monitors follower changes in the background. You focus on creating."
                  icon={BarChart3}
                />
                <StepCard
                  number="3"
                  title="Get Insights"
                  description="Review digestible reports on your schedule. Daily, weekly, or monthly summaries."
                  icon={Sparkles}
                />
              </div>
            </div>
          </div>
        </section>

        {/* Privacy Commitment */}
        <section className="relative px-6 py-32">
          <div className="mx-auto max-w-4xl">
            <GlassCard className="overflow-hidden">
              <div className="relative p-12 sm:p-16">
                {/* Background gradient */}
                <div className="absolute inset-0 bg-gradient-to-br from-violet-600/10 via-transparent to-indigo-600/10" />

                <div className="relative text-center">
                  <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-violet-600 to-indigo-700 shadow-lg shadow-violet-500/25">
                    <Shield className="h-8 w-8 text-white" />
                  </div>

                  <h2 className="mb-4 text-3xl font-bold tracking-tight text-white sm:text-4xl">
                    Our Privacy Promise
                  </h2>

                  <p className="mb-8 text-lg leading-relaxed text-slate-300">
                    In a world of data harvesting and surveillance capitalism, we choose a different path.
                    We believe you deserve tools that respect your privacy and your mental well-being.
                  </p>

                  <div className="grid gap-4 sm:grid-cols-2">
                    {[
                      'No third-party data sharing',
                      'End-to-end encryption',
                      'Read-only API access',
                      'Delete your data anytime',
                    ].map((item, i) => (
                      <div key={i} className="flex items-center gap-3 rounded-xl bg-white/5 px-4 py-3">
                        <CheckCircle2 className="h-5 w-5 flex-shrink-0 text-violet-400" />
                        <span className="text-sm font-medium text-slate-300">{item}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </GlassCard>
          </div>
        </section>

        {/* Social Proof */}
        <section className="relative px-6 py-32">
          <div className="mx-auto max-w-7xl">
            <div className="mb-16 text-center">
              <h2 className="mb-4 text-3xl font-bold tracking-tight text-white sm:text-4xl">
                Loved by <span className="text-violet-400">creators</span>
              </h2>
              <p className="mx-auto max-w-2xl text-lg text-slate-400">
                Join thousands of creators who chose ethical analytics.
              </p>
            </div>

            <div className="grid gap-6 md:grid-cols-3">
              <TestimonialCard
                quote="Finally, an analytics tool that doesn't make me anxious. The weekly digest is perfect—I stay informed without obsessing over numbers."
                author="Maya Thompson"
                role="Content Creator, 52K followers"
                avatar="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=128&h=128&fit=crop&crop=face"
              />
              <TestimonialCard
                quote="The privacy focus sold me. Knowing my data isn't being sold to advertisers makes this worth every penny. Plus the UI is gorgeous."
                author="Jordan Miller"
                role="Photographer, 118K followers"
                avatar="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=128&h=128&fit=crop&crop=face"
              />
              <TestimonialCard
                quote="I manage 12 brand accounts and Unfollow Tracker saves me hours every week. Clean, simple, and respects our clients' privacy."
                author="Sofia Rodriguez"
                role="Social Media Manager"
                avatar="https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=128&h=128&fit=crop&crop=face"
              />
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="relative px-6 py-32">
          <div className="mx-auto max-w-4xl">
            <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-violet-600 via-indigo-700 to-violet-800 p-12 text-center sm:p-16">
              {/* Background pattern */}
              <div className="absolute inset-0 opacity-20">
                <div
                  className="absolute inset-0"
                  style={{
                    backgroundImage: `radial-gradient(circle at 2px 2px, rgba(255,255,255,0.15) 1px, transparent 0)`,
                    backgroundSize: '32px 32px',
                  }}
                />
              </div>

              <div className="relative">
                <h2 className="mb-4 text-3xl font-bold tracking-tight text-white sm:text-4xl">
                  Ready to grow with confidence?
                </h2>
                <p className="mx-auto mb-8 max-w-xl text-lg text-violet-100">
                  Join thousands of creators who chose ethical, privacy-first analytics.
                  Start your free trial today.
                </p>
                <div className="flex flex-col items-center justify-center gap-4 sm:flex-row">
                  <Link to="/register">
                    <Button
                      size="lg"
                      className="h-14 bg-white px-8 text-lg font-semibold text-violet-700 hover:bg-violet-50"
                    >
                      Get Started Free
                    </Button>
                  </Link>
                  <Link to="/login">
                    <Button
                      size="lg"
                      variant="outline"
                      className="h-14 border-white/30 px-8 text-lg font-semibold text-white hover:bg-white/10"
                    >
                      Sign In
                    </Button>
                  </Link>
                </div>
                <p className="mt-6 text-sm text-violet-200">
                  No credit card required. 14-day free trial.
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="relative z-10 border-t border-white/10 bg-[#05040a]/80 px-6 py-16 backdrop-blur-xl">
        <div className="mx-auto max-w-7xl">
          <div className="grid gap-12 lg:grid-cols-4">
            {/* Brand */}
            <div>
              <Link to="/" className="mb-4 flex items-center gap-3">
                <SiteLogo
                  iconClassName="h-10 w-10 rounded-xl"
                  textClassName="text-xl font-bold text-white"
                />
              </Link>
              <p className="text-sm leading-relaxed text-slate-400">
                Privacy-first Instagram analytics for creators who value their data and their peace of mind.
              </p>
            </div>

            {/* Links */}
            <div>
              <h4 className="mb-4 text-sm font-semibold uppercase tracking-wider text-white">Product</h4>
              <ul className="space-y-3">
                {['Features', 'Pricing', 'Security', 'API'].map((item) => (
                  <li key={item}>
                    <a href="#" className="text-sm text-slate-400 transition-colors hover:text-violet-400">
                      {item}
                    </a>
                  </li>
                ))}
              </ul>
            </div>

            <div>
              <h4 className="mb-4 text-sm font-semibold uppercase tracking-wider text-white">Company</h4>
              <ul className="space-y-3">
                {['About', 'Blog', 'Careers', 'Contact'].map((item) => (
                  <li key={item}>
                    <a href="#" className="text-sm text-slate-400 transition-colors hover:text-violet-400">
                      {item}
                    </a>
                  </li>
                ))}
              </ul>
            </div>

            <div>
              <h4 className="mb-4 text-sm font-semibold uppercase tracking-wider text-white">Legal</h4>
              <ul className="space-y-3">
                {['Privacy Policy', 'Terms of Service', 'Cookie Policy', 'GDPR'].map((item) => (
                  <li key={item}>
                    <a href="#" className="text-sm text-slate-400 transition-colors hover:text-violet-400">
                      {item}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          </div>

          <div className="mt-12 flex flex-col items-center justify-between gap-4 border-t border-white/10 pt-8 sm:flex-row">
            <p className="text-sm text-slate-500">
              © {new Date().getFullYear()} Unfollow Tracker. All rights reserved.
            </p>
            <div className="flex items-center gap-2 text-slate-500">
              <MessageCircle className="h-4 w-4" />
              <span className="text-sm">support@unfollowtracker.io</span>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
