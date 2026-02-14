import { create } from 'zustand';
import type { Recipient, CreateRecipientRequest, UpdateRecipientRequest } from '../types/recipient';
import * as recipientService from '../services/recipientService';

interface RecipientState {
  recipients: Recipient[];
  isLoading: boolean;
  error: string | null;
  fetchRecipients: () => Promise<void>;
  createRecipient: (data: CreateRecipientRequest) => Promise<Recipient>;
  updateRecipient: (id: string, data: UpdateRecipientRequest) => Promise<Recipient>;
  deleteRecipient: (id: string) => Promise<void>;
  bulkDelete: (ids: string[]) => Promise<void>;
  clearError: () => void;
}

export const useRecipientStore = create<RecipientState>((set, get) => ({
  recipients: [],
  isLoading: false,
  error: null,

  fetchRecipients: async () => {
    set({ isLoading: true, error: null });
    try {
      const recipients = await recipientService.listRecipients();
      set({ recipients: recipients || [], isLoading: false });
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to fetch recipients';
      set({ error: message, isLoading: false });
    }
  },

  createRecipient: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const rec = await recipientService.createRecipient(data);
      set({ recipients: [...get().recipients, rec], isLoading: false });
      return rec;
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to create recipient';
      set({ error: message, isLoading: false });
      throw err;
    }
  },

  updateRecipient: async (id, data) => {
    set({ isLoading: true, error: null });
    try {
      const rec = await recipientService.updateRecipient(id, data);
      set({
        recipients: get().recipients.map((r) => (r.id === id ? rec : r)),
        isLoading: false,
      });
      return rec;
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to update recipient';
      set({ error: message, isLoading: false });
      throw err;
    }
  },

  deleteRecipient: async (id) => {
    set({ isLoading: true, error: null });
    try {
      await recipientService.deleteRecipient(id);
      set({
        recipients: get().recipients.filter((r) => r.id !== id),
        isLoading: false,
      });
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to delete';
      set({ error: message, isLoading: false });
      throw err;
    }
  },

  bulkDelete: async (ids) => {
    set({ isLoading: true, error: null });
    try {
      await recipientService.bulkDeleteRecipients(ids);
      set({
        recipients: get().recipients.filter((r) => !ids.includes(r.id)),
        isLoading: false,
      });
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to delete';
      set({ error: message, isLoading: false });
      throw err;
    }
  },

  clearError: () => set({ error: null }),
}));
