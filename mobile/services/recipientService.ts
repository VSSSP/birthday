import api from "./api";
import {
  Recipient,
  CreateRecipientRequest,
  UpdateRecipientRequest,
} from "../types/recipient";

export const recipientService = {
  create: async (data: CreateRecipientRequest): Promise<Recipient> => {
    const { data: recipient } = await api.post<Recipient>(
      "/api/recipients",
      data
    );
    return recipient;
  },

  list: async (): Promise<Recipient[]> => {
    const { data } = await api.get<Recipient[]>("/api/recipients");
    return data;
  },

  getById: async (id: string): Promise<Recipient> => {
    const { data } = await api.get<Recipient>(`/api/recipients/${id}`);
    return data;
  },

  update: async (
    id: string,
    data: UpdateRecipientRequest
  ): Promise<Recipient> => {
    const { data: recipient } = await api.put<Recipient>(
      `/api/recipients/${id}`,
      data
    );
    return recipient;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/api/recipients/${id}`);
  },

  bulkDelete: async (ids: string[]): Promise<void> => {
    await api.delete("/api/recipients", { data: { ids } });
  },
};
