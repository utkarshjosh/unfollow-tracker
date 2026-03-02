# Plan: Unfollow Tracker Frontend MVP

## Context
Build a complete, beautiful, modular frontend for Unfollow Tracker - a privacy-first Instagram analytics SaaS. The MVP will include a landing page and dashboard, using shadcn/ui for components, Tailwind CSS for styling (brand colors already configured), and real integration with the existing Go backend API.

## Objectives

### Must Have
- Landing page with hero, features, CTA sections
- Authentication (login/register) with JWT
- Dashboard with stats cards and recent unfollows
- Add/connect Instagram accounts
- Real API integration with Go backend
- Dark theme by default (brand-compliant)
- Responsive design (mobile + desktop)
- No broken CSS - polished, professional appearance

### Must NOT
- Include social OAuth (Instagram login)
- Build account detail/settings pages (out of MVP scope)
- Use mock data
- Compromise on accessibility

## UI Library Decision

**Selected: shadcn/ui**

Rationale:
- Native Tailwind CSS integration (existing setup)
- Modern, non-generic aesthetic (stands out from Material Design)
- Full code ownership (copy-paste components)
- Excellent dark mode support
- Built on accessible Radix UI primitives
- Small bundle size (~20-80KB)

## Backend API Summary

**Authentication:**
- POST /api/v1/auth/register - Email/password registration
- POST /api/v1/auth/login - Returns JWT token + user info
- GET /api/v1/me - Current user profile (requires JWT)

**Accounts:**
- GET /api/v1/accounts - List tracked Instagram accounts
- POST /api/v1/auth/register - Create new account
- DELETE /api/v1/accounts/{id} - Remove account

**Analytics:**
- GET /api/v1/unfollows - List detected unfollows (paginated)
- GET /api/v1/unfollows/summary - Aggregated data with health score
- GET /api/v1/accounts/{id}/stats - Per-account statistics

**Auth Format:** `Authorization: Bearer <jwt_token>`

## Implementation Steps

### Phase 1: Setup shadcn/ui (Step 1)
1. Install shadcn/ui CLI and initialize in web/ directory
2. Configure components.json with brand colors
3. Install base components: button, card, input, label, tabs, dialog, dropdown-menu
4. Install data display components: table, badge, skeleton
5. Install form components: form (with react-hook-form + zod)
6. Verify all components render correctly with brand colors

### Phase 2: Project Structure & Routing (Step 2)
1. Install React Router v6
2. Create route structure:
   - `/` - Landing page (public)
   - `/login` - Login page (public)
   - `/register` - Registration page (public)
   - `/dashboard` - Main dashboard (protected)
3. Create layout components:
   - `RootLayout` - Common layout wrapper
   - `AuthLayout` - Centered auth form layout
   - `DashboardLayout` - Sidebar + main content
4. Create AuthContext for global auth state
5. Create ProtectedRoute component for authenticated routes

### Phase 3: API Client Setup (Step 3)
1. Create `src/lib/api.ts` with axios instance
2. Configure base URL (from env: VITE_API_URL=http://localhost:8080)
3. Add request interceptor to inject JWT token
4. Add response interceptor for 401 handling (redirect to login)
5. Create type definitions in `src/types/api.ts`:
   - User, Account, Unfollow, LoginResponse, etc.
6. Create API functions:
   - auth.ts - login, register, getMe
   - accounts.ts - listAccounts, createAccount, deleteAccount
   - unfollows.ts - getUnfollows, getUnfollowsSummary

### Phase 4: Authentication Pages (Step 4)
1. Create Login page (`src/pages/Login.tsx`):
   - Email/password form with validation
   - Error handling (invalid credentials)
   - Link to register page
   - Submit to /api/v1/auth/login
   - Store token in localStorage
2. Create Register page (`src/pages/Register.tsx`):
   - Email/password/confirm password form
   - Validation (password min 8 chars, matching)
   - Error handling (email exists)
   - Link to login page
   - Submit to /api/v1/auth/register
   - Auto-login after registration
3. Create AuthGuard for protected routes
4. Add logout functionality in dashboard

### Phase 5: Landing Page (Step 5)
1. Create `src/pages/Landing.tsx`:
   - Hero section with value proposition
   - Features section (3-4 key benefits)
   - Privacy commitment section
   - CTA buttons (Get Started -> /register)
   - Simple footer
2. Components to build:
   - Hero section with gradient background
   - Feature cards with icons (Lucide React)
   - Trust badges section
3. Navigation header with login/register links
4. Responsive design (mobile hamburger menu)

### Phase 6: Dashboard Layout (Step 6)
1. Create DashboardLayout (`src/layouts/DashboardLayout.tsx`):
   - Fixed sidebar (240px) with navigation
   - Header with user menu (logout)
   - Main content area
2. Sidebar navigation:
   - Dashboard (active)
   - Accounts (for future expansion)
   - Settings (for future expansion)
3. User dropdown menu with logout
4. Mobile responsive (collapsible sidebar)

### Phase 7: Dashboard Stats Cards (Step 7)
1. Create StatsCards component (`src/components/dashboard/StatsCards.tsx`):
   - Total Followers card (from accounts data)
   - Recent Unfollowers card (from unfollows summary)
   - Engagement Rate card (calculated/placeholder)
   - Account Health card (from health score)
2. Use glass-card utility for consistent styling
3. Add loading skeletons
4. Color-code metrics:
   - Green for positive growth
   - Red for unfollowers
   - Amber for warnings
   - Violet for primary CTAs

### Phase 8: Accounts Management (Step 8)
1. Create AccountList component (`src/components/dashboard/AccountList.tsx`):
   - Table of connected Instagram accounts
   - Columns: Username, Platform, Follower Count, Status, Actions
   - Delete account action with confirmation dialog
2. Create AddAccountDialog:
   - Input for Instagram username
   - Platform selector (Instagram only for now)
   - Submit to POST /api/v1/accounts
   - Refresh list after add
3. Empty state for no accounts
4. Loading skeleton

### Phase 9: Unfollows Display (Step 9)
1. Create UnfollowsList component (`src/components/dashboard/UnfollowsList.tsx`):
   - List of recent unfollowers
   - Show: detected date, associated account
   - Pagination support
2. Create UnfollowsSummary widget:
   - Mini chart or trend indicator
   - Period selector (day/week/month)
3. Connect to real API endpoints:
   - GET /api/v1/unfollows
   - GET /api/v1/unfollows/summary
4. Empty state when no unfollows detected

### Phase 10: Polish & Integration (Step 10)
1. Add loading states throughout
2. Add error handling with toast notifications
3. Add empty states for all data views
4. Test responsive breakpoints (mobile, tablet, desktop)
5. Verify dark theme consistency
6. Test real API integration end-to-end:
   - Register new user
   - Login
   - Add Instagram account
   - View dashboard data
7. Fix any CSS issues
8. Verify build succeeds: npm run build

## Files to Modify/Create

| File | Purpose |
|------|---------|
| `web/components.json` | shadcn/ui configuration |
| `web/src/lib/utils.ts` | shadcn/ui utilities (cn function) |
| `web/src/components/ui/*` | shadcn/ui components (auto-generated) |
| `web/src/lib/api.ts` | Axios API client setup |
| `web/src/types/api.ts` | TypeScript API types |
| `web/src/contexts/AuthContext.tsx` | Global auth state |
| `web/src/layouts/RootLayout.tsx` | Root layout wrapper |
| `web/src/layouts/AuthLayout.tsx` | Auth pages layout |
| `web/src/layouts/DashboardLayout.tsx` | Dashboard layout with sidebar |
| `web/src/pages/Landing.tsx` | Landing page |
| `web/src/pages/Login.tsx` | Login page |
| `web/src/pages/Register.tsx` | Registration page |
| `web/src/pages/Dashboard.tsx` | Main dashboard page |
| `web/src/components/dashboard/StatsCards.tsx` | Stats cards grid |
| `web/src/components/dashboard/AccountList.tsx` | Accounts table |
| `web/src/components/dashboard/AddAccountDialog.tsx` | Add account modal |
| `web/src/components/dashboard/UnfollowsList.tsx` | Unfollowers list |
| `web/src/App.tsx` | App with routing |
| `web/.env.example` | Environment variables template |

## Dependencies to Install

```bash
# shadcn/ui base
npx shadcn-ui@latest init

# shadcn components
npx shadcn add button card input label tabs dialog dropdown-menu table badge skeleton form

# Routing & state
npm install react-router-dom

# HTTP client
npm install axios

# Forms & validation
npm install react-hook-form zod @hookform/resolvers

# Icons
npm install lucide-react

# Utilities
npm install clsx tailwind-merge
```

## Environment Variables

```bash
VITE_API_URL=http://localhost:8080
```

## Acceptance Criteria

- [ ] Landing page renders with hero, features, and CTA sections
- [ ] Users can register with email/password
- [ ] Users can login and receive JWT token
- [ ] Authenticated users see dashboard at /dashboard
- [ ] Unauthenticated users are redirected to /login
- [ ] Dashboard displays stats cards with real data from API
- [ ] Users can add Instagram accounts (POST /api/v1/accounts)
- [ ] Users can view connected accounts in a table
- [ ] Users can delete accounts with confirmation
- [ ] Dashboard displays recent unfollows from API
- [ ] Users can logout (token cleared, redirected to /)
- [ ] All pages are responsive (mobile, tablet, desktop)
- [ ] Dark theme is applied consistently
- [ ] No console errors or broken CSS
- [ ] npm run build succeeds without errors
- [ ] Real API integration works end-to-end

## Brand Compliance Checklist

- [ ] Deep Indigo (#1E1B4B) used for surfaces
- [ ] Soft Violet (#7C3AED) used for accents/CTAs
- [ ] Outfit font for headings
- [ ] DM Sans font for body text
- [ ] Glassmorphism cards (glass-card utility)
- [ ] Dark mode default
- [ ] No harsh shadows, layered shadows only
