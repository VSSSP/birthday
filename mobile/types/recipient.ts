export interface Recipient {
  id: string;
  user_id: string;
  name: string;
  age: number;
  gender: string;
  min_budget: number;
  max_budget: number;
  keywords: string[];
  created_at: string;
  updated_at: string;
}

export interface CreateRecipientRequest {
  name: string;
  age: number;
  gender: string;
  min_budget: number;
  max_budget: number;
  keywords: string[];
}

export interface UpdateRecipientRequest {
  name?: string;
  age?: number;
  gender?: string;
  min_budget?: number;
  max_budget?: number;
  keywords?: string[];
}
