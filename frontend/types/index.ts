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

export interface UserInfo {
  id: string;
  email: string;
  roles: string[];
}

export interface LoginResponse {
  token: string;
  expires_in: number;
  user: UserInfo;
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

export interface Analysis {
  id: string;
  previous_analysis_id?: string | null;
  period_start: string;
  period_end: string;
  feedback_count: number;
  new_feedback_count?: number | null;
  overall_summary: string;
  sentiment: 'positive' | 'mixed' | 'negative';
  key_insights: string[];
  model: string;
  tokens: number;
  analysis_duration_ms: number;
  status: 'processing' | 'success' | 'failed';
  failure_reason?: string | null;
  created_at: string;
  completed_at?: string | null;
}

export interface TopicAnalysis {
  id: string;
  topic: string;
  topic_name: string;
  summary: string;
  feedback_count: number;
  sentiment: 'positive' | 'mixed' | 'negative';
  created_at: string;
  updated_at: string;
}

export interface FeedbackWithTopics {
  id: string;
  rating: number;
  comment: string;
  created_at: string;
  topics: string[];
}

export interface AnalysisDetail {
  analysis: Analysis;
  topics: TopicAnalysis[];
  feedbacks: FeedbackWithTopics[];
}

export interface AnalysisListResponse {
  analyses: Analysis[];
  total: number;
}
