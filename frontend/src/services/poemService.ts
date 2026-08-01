import { api } from './api';
import type { Poem, Dynasty, Poet, PoemSearchParams, GenreCategory } from '@/types';

export type PoemID = string | number;

export interface PoemListResponse {
  list: Poem[];
  total: number;
  page: number;
  pageSize: number;
}

export interface PoetListParams {
  keyword?: string;
  dynastyId?: string | number;
  page?: number;
  pageSize?: number;
}

export interface PoetListResponse {
  list: Poet[];
  total: number;
  page: number;
  pageSize: number;
}

export type CorrectionType = 'title' | 'author' | 'dynasty' | 'content' | 'translation' | 'appreciation' | 'annotation' | 'other';

export interface CorrectionCreateRequest {
  poemId: PoemID;
  type: CorrectionType;
  originalText: string;
  suggestedText: string;
  reason: string;
}

export const poemService = {
  getPoems(params?: PoemSearchParams): Promise<PoemListResponse> {
    return api.get<PoemListResponse>('/poems', params as Record<string, unknown>);
  },

  getPoemById(id: PoemID): Promise<Poem> {
    return api.get<Poem>(`/poems/${id}`);
  },

  getRandomPoems(limit = 5): Promise<Poem[]> {
    return api.get<Poem[]>('/poems/random', { limit });
  },

  likePoem(id: PoemID): Promise<void> {
    return api.post<void>(`/poems/${id}/like`);
  },

  createCorrection(data: CorrectionCreateRequest): Promise<void> {
    return api.post<void>('/corrections', data);
  },

  getDynasties(): Promise<Dynasty[]> {
    return api.get<Dynasty[]>('/dynasties');
  },

  getPoets(keyword?: string): Promise<Poet[]> {
    return api.get<Poet[]>('/poets', keyword ? { keyword } : undefined);
  },

  getPoetList(params: PoetListParams): Promise<PoetListResponse> {
    return api.get<PoetListResponse>('/poets', params as Record<string, unknown>);
  },

  getGenres(): Promise<string[]> {
    return api.get<string[]>('/genres');
  },

  getGenreCategories(): Promise<GenreCategory[]> {
    return api.get<GenreCategory[]>('/genre-categories');
  },
};
