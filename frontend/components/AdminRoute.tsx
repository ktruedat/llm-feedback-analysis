'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth-store';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { AlertCircle } from 'lucide-react';

export function AdminRoute({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const { isAuthenticated, token, isAdmin } = useAuthStore();

  useEffect(() => {
    if (!isAuthenticated || !token) {
      router.push('/login');
      return;
    }

    if (!isAdmin) {
      // User is authenticated but not an admin, redirect to feedback page
      router.push('/feedback');
    }
  }, [isAuthenticated, token, isAdmin, router]);

  if (!isAuthenticated || !token) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Loading...</h1>
        </div>
      </div>
    );
  }

  if (!isAdmin) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Alert variant="destructive" className="max-w-md">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Access Denied</AlertTitle>
          <AlertDescription>
            You do not have permission to access this page. Redirecting...
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  return <>{children}</>;
}
