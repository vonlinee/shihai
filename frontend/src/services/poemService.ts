import { api } from './api';
import type { Poem, Dynasty, Poet, PoemSearchParams } from '@/types';

export interface PoemListResponse {
  list: Poem[];
  total: number;
  page: number;
  pageSize: number;
}

export const poemService = {
  getPoems(params?: PoemSearchParams): Promise<PoemListResponse> {
    return api.get<PoemListResponse>('/poems', params as Record<string, unknown>);
  },

  getPoemById(id: number): Promise<Poem> {
    return api.get<Poem>(`/poems/${id}`);
  },

  getRandomPoems(limit = 5): Promise<Poem[]> {
    return api.get<Poem[]>('/poems/random', { limit });
  },

  likePoem(id: number): Promise<void> {
    return api.post<void>(`/poems/${id}/like`);
  },

  getDynasties(): Promise<Dynasty[]> {
    return api.get<Dynasty[]>('/dynasties');
  },

  getPoets(keyword?: string): Promise<Poet[]> {
    return api.get<Poet[]>('/poets', keyword ? { keyword } : undefined);
  },

  getGenres(): Promise<string[]> {
    return api.get<string[]>('/genres');
  },
};
