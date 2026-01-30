export interface User {
  id: string;
  email: string;
  created_at: string;
}

export interface Feedback {
  id: string;
  rating: number;
  comment: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

export interface FeedbackListResponse {
  feedbacks: Feedback[];
  total: number;
}

export interface LoginResponse {
  token: string;
  expires_in: number;
}

export interface RegisterResponse {
  id: string;
  email: string;
  created_at: string;
}

export interface ApiResponse<T> {
  data?: T;
  message?: string;
  error?: string;
}
