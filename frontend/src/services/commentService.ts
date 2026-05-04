import { api } from './api';
import type { Comment } from '@/types';

export interface CommentListResponse {
  list: Comment[];
  total: number;
  page: number;
  pageSize: number;
}

export interface CreateCommentRequest {
  poemId: number;
  content: string;
  parentId?: number;
  visitorId?: string;
  visitorName?: string;
}

export interface VoteCommentRequest {
  commentId: number;
  type: 'like' | 'dislike';
}

export const commentService = {
  getComments(
    poemId: number,
    page = 1,
    pageSize = 10,
  ): Promise<CommentListResponse> {
    return api.get<CommentListResponse>('/comments', {
      poemId,
      page,
      pageSize,
    });
  },

  createComment(data: CreateCommentRequest): Promise<Comment> {
    return api.post<Comment>('/comments', data);
  },

  deleteComment(id: number): Promise<void> {
    return api.delete<void>(`/comments/${id}`);
  },

  voteComment(data: VoteCommentRequest): Promise<void> {
    return api.post<void>('/comments/vote', data);
  },
};
