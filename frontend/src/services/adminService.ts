import { api } from './api';
import type { User, Poem, Announcement, Dynasty, Poet, WorkCollection, WorkCollectionItem, PoemAnnotation, CorrectionRequest, ForumPost, PoemType, CiTune } from '@/types';

type AdminID = string | number;

export interface DashboardStat {
  /** 统计项稳定标识。 */
  key: 'users' | 'poems' | 'comments' | 'pendingCorrections';
  /** 统计项展示标题。 */
  title: string;
  /** 当前总数。 */
  value: number;
  /** 最近 7 天相对前 7 天的变化百分比。 */
  trendPercent: number;
}

export interface DashboardActivity {
  /** 活动稳定标识。 */
  id: string;
  /** 活动类型标题。 */
  action: string;
  /** 活动详情文案。 */
  detail: string;
  /** 活动发生时间。 */
  createdAt: string;
}

export interface DashboardResponse {
  /** 仪表盘统计卡片。 */
  stats: DashboardStat[];
  /** 最近活动列表。 */
  recentActivities: DashboardActivity[];
}

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
  /** 与 content 每个文本元素对应的平仄标记，只包含平/仄/? 或空字符串。 */
  pingze?: string[];
  authorId?: number | string;
  authorName?: string;
  dynastyId?: number | string;
  dynastyName?: string;
  /** 体裁一级分类，例如诗、词、曲、文。 */
  genreCategory?: string;
  genre?: string;
  /** 词牌 ID，仅词使用；未传时后端会按标题解析。 */
  ciTuneId?: number | string;
  translation?: string;
  appreciation?: string;
  annotation?: string;
  annotations?: PoemAnnotationUpsertRequest[];
  audioUrl?: string;
  coverImage?: string;
}

export interface PoemUpdateRequest {
  title?: string;
  content?: string[];
  /** 与 content 每个文本元素对应的平仄标记，只包含平/仄/? 或空字符串。 */
  pingze?: string[];
  authorId?: number | string;
  dynastyId?: number | string;
  /** 体裁一级分类，例如诗、词、曲、文。 */
  genreCategory?: string;
  genre?: string;
  /** 词牌 ID，仅词使用；传 0 可清空。 */
  ciTuneId?: number | string;
  translation?: string;
  appreciation?: string;
  annotation?: string;
  annotations?: PoemAnnotationUpsertRequest[];
  audioUrl?: string;
  coverImage?: string;
}

export interface PoemAnnotationUpsertRequest {
  id?: string;
  targetField?: 'content';
  startLine: number;
  startOffset: number;
  endLine: number;
  endOffset: number;
  selectedText: string;
  title?: string;
  content: string;
  type?: 'note';
  displayOrder?: number;
}

export type TextConversionMode =
  | 's2t'
  | 't2s'
  | 's2tw'
  | 'tw2s'
  | 's2hk'
  | 'hk2s'
  | 't2tw'
  | 'tw2t'
  | 't2hk'
  | 'hk2t'
  | 't2jp'
  | 'jp2t'
  | 'tw2sp';

export interface TextConversionRequest {
  mode: TextConversionMode;
  texts: string[];
}

export interface TextConversionResponse {
  mode: TextConversionMode;
  texts: string[];
}

export interface PingzeRecognitionRequest {
  /** 按正文行拆分的待识别文本列表。 */
  texts: string[];
}

export interface PingzeRecognitionResponse {
  /** 与 texts 每个文本元素对应的平仄标记。 */
  pingze: string[];
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

export interface PoemTypeCreateRequest {
  /** 细分类别名称，也是保存到诗词 genre 字段的值。 */
  name: string;
  /** 一级分类，例如诗、词、曲、文。 */
  category: string;
  /** 常见句数；未限定时传 null 或省略。 */
  lines?: number | null;
  /** 常见每句字数；未限定时传 null 或省略。 */
  charsPerLine?: number | null;
  /** 体裁说明。 */
  description?: string;
}

export interface PoemTypeUpdateRequest {
  /** 细分类别名称，也是保存到诗词 genre 字段的值。 */
  name: string;
  /** 一级分类，例如诗、词、曲、文。 */
  category: string;
  /** 常见句数；未限定时传 null 或省略。 */
  lines?: number | null;
  /** 常见每句字数；未限定时传 null 或省略。 */
  charsPerLine?: number | null;
  /** 体裁说明。 */
  description?: string;
}

export interface CiTuneCreateRequest {
  /** 词牌名，例如水调歌头、念奴娇。 */
  name: string;
  /** 词牌别名，用于同调异名和标题解析。 */
  aliases?: string[];
  /** 所属词体裁 ID，例如小令、中调、长调；未归类时省略。 */
  poemTypeId?: number | string;
  /** 词牌说明。 */
  description?: string;
}

export interface CiTuneUpdateRequest {
  /** 词牌名，例如水调歌头、念奴娇。 */
  name: string;
  /** 词牌别名，用于同调异名和标题解析。 */
  aliases?: string[];
  /** 所属词体裁 ID，例如小令、中调、长调；未归类时省略。 */
  poemTypeId?: number | string;
  /** 词牌说明。 */
  description?: string;
}

export interface CiTuneTitleParseResponse {
  /** 匹配到的词牌；未匹配时省略。 */
  ciTune?: CiTune;
  /** 实际命中的词牌名或别名。 */
  matchedName: string;
}

export interface CorrectionListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
}

export interface CorrectionListResponse {
  list: CorrectionRequest[];
  total: number;
  page: number;
  pageSize: number;
}

export interface CorrectionStatusUpdateRequest {
  status: CorrectionRequest['status'];
}

// ─── Forum Admin ────────────────────────────────────────────────────────────

export interface AdminForumPostListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  includeDeleted?: boolean;
}

export interface AdminForumPostListResponse {
  list: ForumPost[];
  total: number;
  page: number;
  pageSize: number;
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
  getDashboard(): Promise<DashboardResponse> {
    return api.get<DashboardResponse>('/admin/dashboard');
  },

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

  getPoemAnnotations(poemId: AdminID): Promise<PoemAnnotation[]> {
    return api.get<PoemAnnotation[]>(`/admin/poems/${poemId}/annotations`);
  },

  createPoemAnnotation(poemId: AdminID, data: PoemAnnotationUpsertRequest): Promise<PoemAnnotation> {
    return api.post<PoemAnnotation>(`/admin/poems/${poemId}/annotations`, data);
  },

  updatePoemAnnotation(id: AdminID, data: PoemAnnotationUpsertRequest): Promise<PoemAnnotation> {
    return api.put<PoemAnnotation>(`/admin/poem-annotations/${id}`, data);
  },

  deletePoemAnnotation(id: AdminID): Promise<void> {
    return api.delete<void>(`/admin/poem-annotations/${id}`);
  },

  convertTexts(data: TextConversionRequest): Promise<TextConversionResponse> {
    return api.post<TextConversionResponse>('/admin/text-conversion', data);
  },

  recognizePingze(data: PingzeRecognitionRequest): Promise<PingzeRecognitionResponse> {
    return api.post<PingzeRecognitionResponse>('/admin/poems/pingze-recognition', data);
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

  // Admin corrections
  getCorrections(params?: CorrectionListParams): Promise<CorrectionListResponse> {
    return api.get<CorrectionListResponse>('/admin/corrections', params as Record<string, unknown>);
  },

  getCorrection(id: AdminID): Promise<CorrectionRequest> {
    return api.get<CorrectionRequest>(`/admin/corrections/${id}`);
  },

  updateCorrectionStatus(id: AdminID, data: CorrectionStatusUpdateRequest): Promise<CorrectionRequest> {
    return api.put<CorrectionRequest>(`/admin/corrections/${id}/status`, data);
  },

  // Admin forum
  getForumPosts(params?: AdminForumPostListParams): Promise<AdminForumPostListResponse> {
    return api.get<AdminForumPostListResponse>('/admin/forum/posts', params as Record<string, unknown>);
  },

  setForumPostPinned(id: AdminID, isPinned: boolean): Promise<void> {
    return api.put<void>(`/admin/forum/posts/${id}/pin`, { isPinned });
  },

  deleteForumPost(id: AdminID): Promise<void> {
    return api.delete<void>(`/admin/forum/posts/${id}`);
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

  // Poem types
  getPoemTypes(): Promise<PoemType[]> {
    return api.get<PoemType[]>('/admin/poem-types');
  },

  createPoemType(data: PoemTypeCreateRequest): Promise<PoemType> {
    return api.post<PoemType>('/admin/poem-types', data);
  },

  updatePoemType(id: AdminID, data: PoemTypeUpdateRequest): Promise<PoemType> {
    return api.put<PoemType>(`/admin/poem-types/${id}`, data);
  },

  deletePoemType(id: AdminID): Promise<void> {
    return api.delete<void>(`/admin/poem-types/${id}`);
  },

  batchDeletePoemTypes(ids: AdminID[]): Promise<void> {
    return api.delete<void>('/admin/poem-types', { ids });
  },

  // Ci tunes
  getCiTunes(): Promise<CiTune[]> {
    return api.get<CiTune[]>('/admin/ci-tunes');
  },

  createCiTune(data: CiTuneCreateRequest): Promise<CiTune> {
    return api.post<CiTune>('/admin/ci-tunes', data);
  },

  updateCiTune(id: AdminID, data: CiTuneUpdateRequest): Promise<CiTune> {
    return api.put<CiTune>(`/admin/ci-tunes/${id}`, data);
  },

  deleteCiTune(id: AdminID): Promise<void> {
    return api.delete<void>(`/admin/ci-tunes/${id}`);
  },

  batchDeleteCiTunes(ids: AdminID[]): Promise<void> {
    return api.delete<void>('/admin/ci-tunes', { ids });
  },

  parseCiTuneTitle(title: string): Promise<CiTuneTitleParseResponse> {
    return api.post<CiTuneTitleParseResponse>('/admin/ci-tunes/parse-title', { title });
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
