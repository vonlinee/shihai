import { api } from './api';
import type { User, Poem, Announcement, Dynasty, Poet, WorkCollection, WorkCollectionItem } from '@/types';

type AdminID = string | number;

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
  content: string[];
  authorId?: number | string;
  authorName?: string;
  dynastyId?: number | string;
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
  content?: string[];
  authorId?: number | string;
  dynastyId?: number | string;
  genre?: string;
  translation?: string;
  appreciation?: string;
  annotation?: string;
  audioUrl?: string;
  coverImage?: string;
}

// ─── Work Collection Admin ───────────────────────────────────────────────────

export interface WorkCollectionListResponse {
  list: WorkCollection[];
  total: number;
  page: number;
  pageSize: number;
}

export interface WorkCollectionListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  published?: boolean;
}

export interface WorkCollectionCreateRequest {
  title: string;
  description?: string;
  coverImage?: string;
  isPublished?: boolean;
}

export interface WorkCollectionUpdateRequest {
  title?: string;
  description?: string;
  coverImage?: string;
  isPublished?: boolean;
}

export interface WorkCollectionItemCreateRequest {
  workType: string;
  workId: number;
  sortOrder?: number;
}

export interface WorkCollectionItemUpdateRequest {
  sortOrder: number;
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
  dynastyId?: number | string;
  biography?: string;
  avatar?: string;
  birthYear?: number;
  deathYear?: number;
}

export interface PoetUpdateRequest {
  name?: string;
  dynastyId?: number | string;
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

  updatePoem(id: AdminID, data: PoemUpdateRequest): Promise<Poem> {
    return api.put<Poem>(`/admin/poems/${id}`, data);
  },

  deletePoem(id: AdminID): Promise<void> {
    return api.delete<void>(`/admin/poems/${id}`);
  },

  batchDeletePoems(ids: AdminID[]): Promise<void> {
    return api.delete<void>('/admin/poems', { ids });
  },

  // Work Collections
  getWorkCollections(params?: WorkCollectionListParams): Promise<WorkCollectionListResponse> {
    return api.get<WorkCollectionListResponse>('/admin/work-collections', params as Record<string, unknown>);
  },

  getWorkCollectionById(id: number): Promise<WorkCollection> {
    return api.get<WorkCollection>(`/admin/work-collections/${id}`);
  },

  createWorkCollection(data: WorkCollectionCreateRequest): Promise<WorkCollection> {
    return api.post<WorkCollection>('/admin/work-collections', data);
  },

  updateWorkCollection(id: number, data: WorkCollectionUpdateRequest): Promise<WorkCollection> {
    return api.put<WorkCollection>(`/admin/work-collections/${id}`, data);
  },

  deleteWorkCollection(id: number): Promise<void> {
    return api.delete<void>(`/admin/work-collections/${id}`);
  },

  addWorkCollectionItem(collectionId: number, data: WorkCollectionItemCreateRequest): Promise<WorkCollectionItem> {
    return api.post<WorkCollectionItem>(`/admin/work-collections/${collectionId}/items`, data);
  },

  updateWorkCollectionItem(collectionId: number, itemId: number, data: WorkCollectionItemUpdateRequest): Promise<WorkCollectionItem> {
    return api.put<WorkCollectionItem>(`/admin/work-collections/${collectionId}/items/${itemId}`, data);
  },

  deleteWorkCollectionItem(collectionId: number, itemId: number): Promise<void> {
    return api.delete<void>(`/admin/work-collections/${collectionId}/items/${itemId}`);
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

  updateDynasty(id: AdminID, data: DynastyUpdateRequest): Promise<Dynasty> {
    return api.put<Dynasty>(`/admin/dynasties/${id}`, data);
  },

  deleteDynasty(id: AdminID): Promise<void> {
    return api.delete<void>(`/admin/dynasties/${id}`);
  },

  batchDeleteDynasties(ids: AdminID[]): Promise<void> {
    return api.delete<void>('/admin/dynasties', { ids });
  },

  // Poets
  createPoet(data: PoetCreateRequest): Promise<Poet> {
    return api.post<Poet>('/admin/poets', data);
  },

  updatePoet(id: AdminID, data: PoetUpdateRequest): Promise<Poet> {
    return api.put<Poet>(`/admin/poets/${id}`, data);
  },

  deletePoet(id: AdminID): Promise<void> {
    return api.delete<void>(`/admin/poets/${id}`);
  },

  batchDeletePoets(ids: AdminID[]): Promise<void> {
    return api.delete<void>('/admin/poets', { ids });
  },
};
