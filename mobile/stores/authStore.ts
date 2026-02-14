import { create } from "zustand";
import * as SecureStore from "expo-secure-store";
import { User, TokenPair } from "../types/auth";
import { authService } from "../services/authService";

const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";

interface AuthStore {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;

  initialize: () => Promise<void>;
  loginWithEmail: (email: string, password: string) => Promise<void>;
  registerWithEmail: (
    email: string,
    password: string,
    name: string
  ) => Promise<void>;
  loginWithGoogle: (idToken: string) => Promise<void>;
  loginWithApple: (identityToken: string) => Promise<void>;
  refreshTokens: () => Promise<boolean>;
  logout: () => Promise<void>;
  setTokens: (tokens: TokenPair) => Promise<void>;
  fetchUser: () => Promise<void>;
}

export const useAuthStore = create<AuthStore>((set, get) => ({
  user: null,
  accessToken: null,
  refreshToken: null,
  isLoading: true,
  isAuthenticated: false,

  initialize: async () => {
    try {
      const accessToken = await SecureStore.getItemAsync(ACCESS_TOKEN_KEY);
      const refreshToken = await SecureStore.getItemAsync(REFRESH_TOKEN_KEY);

      if (accessToken && refreshToken) {
        set({ accessToken, refreshToken, isAuthenticated: true });
        try {
          await get().fetchUser();
        } catch {
          const refreshed = await get().refreshTokens();
          if (refreshed) {
            await get().fetchUser();
          } else {
            await get().logout();
          }
        }
      }
    } finally {
      set({ isLoading: false });
    }
  },

  setTokens: async (tokens: TokenPair) => {
    await SecureStore.setItemAsync(ACCESS_TOKEN_KEY, tokens.access_token);
    await SecureStore.setItemAsync(REFRESH_TOKEN_KEY, tokens.refresh_token);
    set({
      accessToken: tokens.access_token,
      refreshToken: tokens.refresh_token,
      isAuthenticated: true,
    });
  },

  loginWithEmail: async (email, password) => {
    const tokens = await authService.login(email, password);
    await get().setTokens(tokens);
    await get().fetchUser();
  },

  registerWithEmail: async (email, password, name) => {
    const tokens = await authService.register(email, password, name);
    await get().setTokens(tokens);
    await get().fetchUser();
  },

  loginWithGoogle: async (idToken) => {
    const tokens = await authService.googleLogin(idToken);
    await get().setTokens(tokens);
    await get().fetchUser();
  },

  loginWithApple: async (identityToken) => {
    const tokens = await authService.appleLogin(identityToken);
    await get().setTokens(tokens);
    await get().fetchUser();
  },

  refreshTokens: async () => {
    const { refreshToken } = get();
    if (!refreshToken) return false;
    try {
      const tokens = await authService.refresh(refreshToken);
      await get().setTokens(tokens);
      return true;
    } catch {
      return false;
    }
  },

  fetchUser: async () => {
    const user = await authService.getMe();
    set({ user });
  },

  logout: async () => {
    await SecureStore.deleteItemAsync(ACCESS_TOKEN_KEY);
    await SecureStore.deleteItemAsync(REFRESH_TOKEN_KEY);
    set({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
    });
  },
}));
