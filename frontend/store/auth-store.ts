import { create } from 'zustand';
import { UserInfo } from '@/types';

interface AuthState {
  token: string | null;
  user: UserInfo | null;
  isAuthenticated: boolean;
  isAdmin: boolean;
  setAuth: (token: string, user: UserInfo) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>((set) => {
  // Initialize from localStorage if available
  if (typeof window !== 'undefined') {
    const storedToken = localStorage.getItem('auth_token');
    const storedUser = localStorage.getItem('auth_user');
    if (storedToken && storedUser) {
      try {
        const user: UserInfo = JSON.parse(storedUser);
        const isAdmin = user.roles?.includes('admin') || false;
        return {
          token: storedToken,
          user,
          isAuthenticated: true,
          isAdmin,
          setAuth: (token: string, user: UserInfo) => {
            const admin = user.roles?.includes('admin') || false;
            localStorage.setItem('auth_token', token);
            localStorage.setItem('auth_user', JSON.stringify(user));
            set({ token, user, isAuthenticated: true, isAdmin: admin });
          },
          clearAuth: () => {
            localStorage.removeItem('auth_token');
            localStorage.removeItem('auth_user');
            set({ token: null, user: null, isAuthenticated: false, isAdmin: false });
          },
        };
      } catch {
        // Invalid stored user, clear it
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
      }
    }
  }

  return {
    token: null,
    user: null,
    isAuthenticated: false,
    isAdmin: false,
    setAuth: (token: string, user: UserInfo) => {
      const admin = user.roles?.includes('admin') || false;
      if (typeof window !== 'undefined') {
        localStorage.setItem('auth_token', token);
        localStorage.setItem('auth_user', JSON.stringify(user));
      }
      set({ token, user, isAuthenticated: true, isAdmin: admin });
    },
    clearAuth: () => {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
      }
      set({ token: null, user: null, isAuthenticated: false, isAdmin: false });
    },
  };
});
