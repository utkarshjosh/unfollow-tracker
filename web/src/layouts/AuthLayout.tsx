import { Outlet } from 'react-router-dom';

export function AuthLayout() {
  // Outlet renders the child route component (Login, Register, etc.)
  // Each auth page handles its own full-screen layout
  return <Outlet />;
}
