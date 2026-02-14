import { create } from "zustand";
import { Recipient, CreateRecipientRequest, UpdateRecipientRequest } from "../types/recipient";
import { recipientService } from "../services/recipientService";

interface RecipientStore {
  recipients: Recipient[];
  isLoading: boolean;

  fetchRecipients: () => Promise<void>;
  createRecipient: (data: CreateRecipientRequest) => Promise<Recipient>;
  updateRecipient: (id: string, data: UpdateRecipientRequest) => Promise<Recipient>;
  deleteRecipient: (id: string) => Promise<void>;
  bulkDeleteRecipients: (ids: string[]) => Promise<void>;
}

export const useRecipientStore = create<RecipientStore>((set, get) => ({
  recipients: [],
  isLoading: false,

  fetchRecipients: async () => {
    set({ isLoading: true });
    try {
      const recipients = await recipientService.list();
      set({ recipients });
    } finally {
      set({ isLoading: false });
    }
  },

  createRecipient: async (data) => {
    const recipient = await recipientService.create(data);
    set({ recipients: [recipient, ...get().recipients] });
    return recipient;
  },

  updateRecipient: async (id, data) => {
    const updated = await recipientService.update(id, data);
    set({
      recipients: get().recipients.map((r) => (r.id === id ? updated : r)),
    });
    return updated;
  },

  deleteRecipient: async (id) => {
    await recipientService.delete(id);
    set({ recipients: get().recipients.filter((r) => r.id !== id) });
  },

  bulkDeleteRecipients: async (ids) => {
    await recipientService.bulkDelete(ids);
    set({
      recipients: get().recipients.filter((r) => !ids.includes(r.id)),
    });
  },
}));
