import api from './api';
import type { User, TokenPair, LoginRequest, RegisterRequest } from '../types/auth';

export async function register(data: RegisterRequest): Promise<TokenPair> {
  const res = await api.post<TokenPair>('/api/auth/register', data);
  return res.data;
}

export async function login(data: LoginRequest): Promise<TokenPair> {
  const res = await api.post<TokenPair>('/api/auth/login', data);
  return res.data;
}

export async function getMe(): Promise<User> {
  const res = await api.get<User>('/api/auth/me');
  return res.data;
}

export async function refreshTokens(refreshToken: string): Promise<TokenPair> {
  const res = await api.post<TokenPair>('/api/auth/refresh', { refresh_token: refreshToken });
  return res.data;
}
