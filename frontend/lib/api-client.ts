import axios, { AxiosInstance, AxiosError } from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export interface ApiError {
  message: string;
  code?: string;
}

export class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add request interceptor to include auth token
    this.client.interceptors.request.use((config) => {
      const token = this.getToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          // Token expired or invalid, clear it
          this.clearToken();
          if (typeof window !== 'undefined') {
            window.location.href = '/login';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  private getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('auth_token');
  }

  private clearToken(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem('auth_token');
  }

  setToken(token: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem('auth_token', token);
  }

  clearAuth(): void {
    this.clearToken();
  }

  // Auth endpoints
  async register(email: string, password: string) {
    const response = await this.client.post('/auth/register', { email, password });
    return response.data;
  }

  async login(email: string, password: string) {
    const response = await this.client.post('/auth/login', { email, password });
    // Response is directly the payload
    return response.data;
  }

  // Feedback endpoints
  async createFeedback(rating: number, comment: string) {
    const response = await this.client.post('/feedbacks', { rating, comment });
    return response.data;
  }

  async listFeedbacks(limit?: number, offset?: number) {
    const params = new URLSearchParams();
    if (limit) params.append('limit', limit.toString());
    if (offset) params.append('offset', offset.toString());
    const queryString = params.toString();
    const url = queryString ? `/feedbacks?${queryString}` : '/feedbacks';
    const response = await this.client.get(url);
    return response.data;
  }

  async getFeedback(id: string) {
    const response = await this.client.get(`/feedbacks/${id}`);
    return response.data;
  }

  async deleteFeedback(id: string) {
    await this.client.delete(`/feedbacks/${id}`);
  }

  // Analysis endpoints
  async getLatestAnalysis() {
    const response = await this.client.get('/analyses/latest');
    return response.data;
  }

  async listAnalyses() {
    const response = await this.client.get('/analyses');
    return response.data;
  }

  async getAnalysis(id: string) {
    const response = await this.client.get(`/analyses/${id}`);
    return response.data;
  }

  // Topic endpoints
  async getTopicsWithStats() {
    const response = await this.client.get('/topics');
    return response?.data || response;
  }

  async getTopicDetails(topicEnum: string) {
    const response = await this.client.get(`/topics/${topicEnum}`);
    return response?.data || response;
  }
}

export const apiClient = new ApiClient();
