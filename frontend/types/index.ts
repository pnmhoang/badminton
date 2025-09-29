// Type definitions for the badminton application

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  role: 'player' | 'admin';
  ranking?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Tournament {
  id: number;
  name: string;
  description: string;
  tournament_type: 'singles' | 'doubles' | 'mixed_doubles';
  max_participants: number;
  start_date: string;
  end_date: string;
  registration_deadline: string;
  entry_fee: number;
  prize_pool: number;
  status: 'upcoming' | 'ongoing' | 'completed' | 'cancelled';
  created_at: string;
  updated_at: string;
}

export interface Team {
  id: number;
  name: string;
  player1_id: number;
  player2_id?: number;
  created_at: string;
  updated_at: string;
}

export interface Match {
  id: number;
  tournament_id: number;
  team1_id?: number;
  team2_id?: number;
  player1_id?: number;
  player2_id?: number;
  player3_id?: number;
  player4_id?: number;
  score_team1: number;
  score_team2: number;
  status: 'scheduled' | 'in_progress' | 'completed' | 'cancelled';
  scheduled_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  full_name: string;
}

export interface AuthResponse {
  message: string;
  data: {
    user: User;
    token: string;
  };
}

export interface ApiResponse<T = any> {
  message: string;
  data?: T;
  error?: string;
}