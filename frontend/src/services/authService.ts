import { api } from './api';
import type { User } from '@/types';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
  name?: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface UpdateProfileRequest {
  name?: string;
  avatar?: string;
  gender?: string;
  age?: number;
  phone?: string;
}

export interface ChangePasswordRequest {
  oldPassword: string;
  newPassword: string;
}

export const authService = {
  login(data: LoginRequest): Promise<LoginResponse> {
    return api.post<LoginResponse>('/auth/login', data);
  },

  register(data: RegisterRequest): Promise<User> {
    return api.post<User>('/auth/register', data);
  },

  getProfile(): Promise<User> {
    return api.get<User>('/user/profile');
  },

  updateProfile(data: UpdateProfileRequest): Promise<User> {
    return api.put<User>('/user/profile', data);
  },

  changePassword(data: ChangePasswordRequest): Promise<void> {
    return api.put<void>('/user/password', data);
  },
};
