import api from "./api";
import { TokenPair, User } from "../types/auth";

export const authService = {
  register: async (
    email: string,
    password: string,
    name: string
  ): Promise<TokenPair> => {
    const { data } = await api.post<TokenPair>("/api/auth/register", {
      email,
      password,
      name,
    });
    return data;
  },

  login: async (email: string, password: string): Promise<TokenPair> => {
    const { data } = await api.post<TokenPair>("/api/auth/login", {
      email,
      password,
    });
    return data;
  },

  googleLogin: async (idToken: string): Promise<TokenPair> => {
    const { data } = await api.post<TokenPair>("/api/auth/google", {
      id_token: idToken,
    });
    return data;
  },

  appleLogin: async (identityToken: string): Promise<TokenPair> => {
    const { data } = await api.post<TokenPair>("/api/auth/apple", {
      id_token: identityToken,
    });
    return data;
  },

  refresh: async (refreshToken: string): Promise<TokenPair> => {
    const { data } = await api.post<TokenPair>("/api/auth/refresh", {
      refresh_token: refreshToken,
    });
    return data;
  },

  getMe: async (): Promise<User> => {
    const { data } = await api.get<User>("/api/auth/me");
    return data;
  },
};
