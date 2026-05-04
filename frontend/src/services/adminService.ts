import { api } from './api';
import type { User, Poem, Announcement, Dynasty, Poet } from '@/types';

// ─── User Admin ──────────────────────────────────────────────────────────────

export interface UserListResponse {
  list: User[];
  total: number;
  page: number;
  pageSize: number;
}

export interface UserListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  role?: string;
}

export interface AdminCreateUserRequest {
  username: string;
  password: string;
  name?: string;
  roleIds?: number[];
}

// ─── Poem Admin ──────────────────────────────────────────────────────────────

export interface PoemCreateRequest {
  title: string;
  content: string;
  authorId?: number;
  authorName?: string;
  dynastyId?: number;
  dynastyName?: string;
  genre?: string;
  translation?: string;
  appreciation?: string;
  annotation?: string;
  audioUrl?: string;
  coverImage?: string;
}

export interface PoemUpdateRequest {
  title?: string;
  content?: string;
  authorId?: number;
  dynastyId?: number;
  genre?: string;
  translation?: string;
  appreciation?: string;
  annotation?: string;
  audioUrl?: string;
  coverImage?: string;
}

// ─── Announcement Admin ──────────────────────────────────────────────────────

export interface AnnouncementCreateRequest {
  title: string;
  content: string;
  isPinned?: boolean;
}

export interface AnnouncementUpdateRequest {
  title?: string;
  content?: string;
  isPinned?: boolean;
}

// ─── Dynasty Admin ──────────────────────────────────────────────────────────────

export interface DynastyCreateRequest {
  name: string;
  period?: string;
  description?: string;
}

export interface DynastyUpdateRequest {
  name?: string;
  period?: string;
  description?: string;
}

// ─── Poet Admin ──────────────────────────────────────────────────────────────

export interface PoetCreateRequest {
  name: string;
  dynastyId?: number;
  biography?: string;
  avatar?: string;
  birthYear?: number;
  deathYear?: number;
}

export interface PoetUpdateRequest {
  name?: string;
  dynastyId?: number;
  biography?: string;
  avatar?: string;
  birthYear?: number;
  deathYear?: number;
}

export const adminService = {
  // Users
  getUsers(params?: UserListParams): Promise<UserListResponse> {
    return api.get<UserListResponse>('/admin/users', params as Record<string, unknown>);
  },

  getUserById(id: number): Promise<User> {
    return api.get<User>(`/admin/users/${id}`);
  },

  deleteUser(id: number): Promise<void> {
    return api.delete<void>(`/admin/users/${id}`);
  },

  createUser(data: AdminCreateUserRequest): Promise<User> {
    return api.post<User>('/admin/users', data);
  },

  // Poems
  createPoem(data: PoemCreateRequest): Promise<Poem> {
    return api.post<Poem>('/admin/poems', data);
  },

  updatePoem(id: number, data: PoemUpdateRequest): Promise<Poem> {
    return api.put<Poem>(`/admin/poems/${id}`, data);
  },

  deletePoem(id: number): Promise<void> {
    return api.delete<void>(`/admin/poems/${id}`);
  },

  // Announcements
  createAnnouncement(data: AnnouncementCreateRequest): Promise<Announcement> {
    return api.post<Announcement>('/admin/announcements', data);
  },

  updateAnnouncement(id: number, data: AnnouncementUpdateRequest): Promise<Announcement> {
    return api.put<Announcement>(`/admin/announcements/${id}`, data);
  },

  deleteAnnouncement(id: number): Promise<void> {
    return api.delete<void>(`/admin/announcements/${id}`);
  },

  // Admin comments (all comments)
  getComments(page = 1, pageSize = 10): Promise<{
    list: unknown[];
    total: number;
    page: number;
    pageSize: number;
  }> {
    return api.get('/admin/comments/all', { page, pageSize });
  },

  // Dynasties
  createDynasty(data: DynastyCreateRequest): Promise<Dynasty> {
    return api.post<Dynasty>('/admin/dynasties', data);
  },

  updateDynasty(id: number, data: DynastyUpdateRequest): Promise<Dynasty> {
    return api.put<Dynasty>(`/admin/dynasties/${id}`, data);
  },

  deleteDynasty(id: number): Promise<void> {
    return api.delete<void>(`/admin/dynasties/${id}`);
  },

  // Poets
  createPoet(data: PoetCreateRequest): Promise<Poet> {
    return api.post<Poet>('/admin/poets', data);
  },

  updatePoet(id: number, data: PoetUpdateRequest): Promise<Poet> {
    return api.put<Poet>(`/admin/poets/${id}`, data);
  },

  deletePoet(id: number): Promise<void> {
    return api.delete<void>(`/admin/poets/${id}`);
  },
};
