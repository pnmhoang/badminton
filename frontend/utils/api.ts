import axios, { InternalAxiosRequestConfig } from 'axios';
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '../types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_BASE}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

export const authAPI = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/login', data);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/register', data);
    return response.data;
  },

  getProfile: async (): Promise<User> => {
    const response = await api.get<User>('/profile');
    return response.data;
  },

  updateProfile: async (data: Partial<User>): Promise<{ message: string; data: User }> => {
    const response = await api.put('/profile', data);
    return response.data;
  },

  changePassword: async (data: { current_password: string; new_password: string }) => {
    const response = await api.post('/change-password', data);
    return response.data;
  },
};

export const tournamentAPI = {
  getTournaments: async () => {
    const response = await api.get('/tournaments');
    return response.data;
  },

  getTournament: async (id: number) => {
    const response = await api.get(`/tournaments/${id}`);
    return response.data;
  },

  createTournament: async (data: any) => {
    const response = await api.post('/tournaments', data);
    return response.data;
  },

  registerForTournament: async (tournamentId: number) => {
    const response = await api.post(`/tournament-registration/${tournamentId}`);
    return response.data;
  },

  getMyRegistrations: async () => {
    const response = await api.get('/my-registrations');
    return response.data;
  },
};

export default api;