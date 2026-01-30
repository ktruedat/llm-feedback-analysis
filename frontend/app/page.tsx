'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth-store';

export default function Home() {
  const router = useRouter();
  const { isAuthenticated, isAdmin } = useAuthStore();

  useEffect(() => {
    if (isAuthenticated) {
      // Redirect based on user role
      if (isAdmin) {
        router.push('/admin');
      } else {
        router.push('/feedback');
      }
    } else {
      router.push('/login');
    }
  }, [isAuthenticated, isAdmin, router]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <h1 className="text-2xl font-bold mb-4">Loading...</h1>
      </div>
    </div>
  );
}
