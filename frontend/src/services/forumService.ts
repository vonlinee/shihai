import { api } from './api'
import type { ForumPost, ForumReply } from '@/types'

export interface ForumListResponse<T> {
  /** 当前页数据列表。 */
  list: T[]
  /** 符合筛选条件的总数。 */
  total: number
  /** 当前页码，从 1 开始。 */
  page: number
  /** 每页数量。 */
  pageSize: number
}

export interface CreateForumPostRequest {
  /** 帖子标题，2-200 个字符。 */
  title: string
  /** 帖子正文。 */
  content: string
}

export interface UpdateForumPostRequest {
  /** 帖子标题；不传时保持原值。 */
  title?: string
  /** 帖子正文；不传时保持原值。 */
  content?: string
}

export interface CreateForumReplyRequest {
  /** 回复正文。 */
  content: string
  /** 父回复 ID，用于楼中楼回复。 */
  parentId?: string
}

export interface ForumPostListParams {
  /** 当前页码，从 1 开始。 */
  page?: number
  /** 每页数量。 */
  pageSize?: number
  /** 搜索关键词，匹配标题和正文。 */
  keyword?: string
}

export const forumService = {
  getPosts(params: ForumPostListParams = {}): Promise<ForumListResponse<ForumPost>> {
    return api.get<ForumListResponse<ForumPost>>('/forum/posts', {
      page: params.page,
      pageSize: params.pageSize,
      keyword: params.keyword,
    })
  },

  getPostById(id: string): Promise<ForumPost> {
    return api.get<ForumPost>(`/forum/posts/${id}`)
  },

  createPost(data: CreateForumPostRequest): Promise<ForumPost> {
    return api.post<ForumPost>('/forum/posts', data)
  },

  updatePost(id: string, data: UpdateForumPostRequest): Promise<ForumPost> {
    return api.put<ForumPost>(`/forum/posts/${id}`, data)
  },

  deletePost(id: string): Promise<void> {
    return api.delete<void>(`/forum/posts/${id}`)
  },

  getReplies(postId: string, page = 1, pageSize = 50): Promise<ForumListResponse<ForumReply>> {
    return api.get<ForumListResponse<ForumReply>>(`/forum/posts/${postId}/replies`, {
      page,
      pageSize,
    })
  },

  createReply(postId: string, data: CreateForumReplyRequest): Promise<ForumReply> {
    return api.post<ForumReply>(`/forum/posts/${postId}/replies`, data)
  },

  deleteReply(id: string): Promise<void> {
    return api.delete<void>(`/forum/replies/${id}`)
  },
}
