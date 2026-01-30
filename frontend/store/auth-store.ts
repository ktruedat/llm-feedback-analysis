import { create } from 'zustand';

interface AuthState {
  token: string | null;
  isAuthenticated: boolean;
  setToken: (token: string) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>((set) => {
  // Initialize from localStorage if available
  if (typeof window !== 'undefined') {
    const storedToken = localStorage.getItem('auth_token');
    if (storedToken) {
      return {
        token: storedToken,
        isAuthenticated: true,
        setToken: (token: string) => {
          localStorage.setItem('auth_token', token);
          set({ token, isAuthenticated: true });
        },
        clearAuth: () => {
          localStorage.removeItem('auth_token');
          set({ token: null, isAuthenticated: false });
        },
      };
    }
  }

  return {
    token: null,
    isAuthenticated: false,
    setToken: (token: string) => {
      if (typeof window !== 'undefined') {
        localStorage.setItem('auth_token', token);
      }
      set({ token, isAuthenticated: true });
    },
    clearAuth: () => {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_token');
      }
      set({ token: null, isAuthenticated: false });
    },
  };
});
