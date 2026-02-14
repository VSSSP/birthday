import api from './api';
import type { Recipient, CreateRecipientRequest, UpdateRecipientRequest } from '../types/recipient';

export async function listRecipients(): Promise<Recipient[]> {
  const res = await api.get<Recipient[]>('/api/recipients');
  return res.data;
}

export async function getRecipient(id: string): Promise<Recipient> {
  const res = await api.get<Recipient>(`/api/recipients/${id}`);
  return res.data;
}

export async function createRecipient(data: CreateRecipientRequest): Promise<Recipient> {
  const res = await api.post<Recipient>('/api/recipients', data);
  return res.data;
}

export async function updateRecipient(id: string, data: UpdateRecipientRequest): Promise<Recipient> {
  const res = await api.put<Recipient>(`/api/recipients/${id}`, data);
  return res.data;
}

export async function deleteRecipient(id: string): Promise<void> {
  await api.delete(`/api/recipients/${id}`);
}

export async function bulkDeleteRecipients(ids: string[]): Promise<void> {
  await api.delete('/api/recipients', { data: { ids } });
}
