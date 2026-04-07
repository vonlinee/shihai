import { api } from './api';
import type { User, Poem, Announcement } from '@/types';

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

// ─── Poem Admin ──────────────────────────────────────────────────────────────

export interface PoemCreateRequest {
  title: string;
  content: string;
  authorId: number;
  dynastyId: number;
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
};
