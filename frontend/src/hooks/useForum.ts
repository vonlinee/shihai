import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import {
  forumService,
  type CreateForumPostRequest,
  type CreateForumReplyRequest,
  type ForumPostListParams,
  type UpdateForumPostRequest,
} from '@/services/forumService'

export function useForumPosts(params: ForumPostListParams = {}) {
  return useQuery({
    queryKey: ['forum', 'posts', params],
    queryFn: () => forumService.getPosts(params),
  })
}

export function useForumPost(id: string | undefined) {
  return useQuery({
    queryKey: ['forum', 'post', id],
    queryFn: () => forumService.getPostById(id!),
    enabled: !!id,
  })
}

export function useForumReplies(postId: string | undefined, page = 1, pageSize = 50) {
  return useQuery({
    queryKey: ['forum', 'replies', postId, page, pageSize],
    queryFn: () => forumService.getReplies(postId!, page, pageSize),
    enabled: !!postId,
  })
}

export function useCreateForumPost() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateForumPostRequest) => forumService.createPost(data),
    onSuccess: () => {
      toast.success('发帖成功')
      queryClient.invalidateQueries({ queryKey: ['forum', 'posts'] })
    },
    onError: (error: Error) => {
      toast.error(error.message || '发帖失败')
    },
  })
}

export function useUpdateForumPost() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateForumPostRequest }) => forumService.updatePost(id, data),
    onSuccess: (post) => {
      toast.success('帖子已更新')
      queryClient.invalidateQueries({ queryKey: ['forum', 'posts'] })
      queryClient.invalidateQueries({ queryKey: ['forum', 'post', post.id] })
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败')
    },
  })
}

export function useDeleteForumPost() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => forumService.deletePost(id),
    onSuccess: () => {
      toast.success('帖子已删除')
      queryClient.invalidateQueries({ queryKey: ['forum'] })
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败')
    },
  })
}

export function useCreateForumReply(postId: string | undefined) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateForumReplyRequest) => forumService.createReply(postId!, data),
    onSuccess: () => {
      toast.success('回复成功')
      queryClient.invalidateQueries({ queryKey: ['forum', 'posts'] })
      queryClient.invalidateQueries({ queryKey: ['forum', 'post', postId] })
      queryClient.invalidateQueries({ queryKey: ['forum', 'replies', postId] })
    },
    onError: (error: Error) => {
      toast.error(error.message || '回复失败')
    },
  })
}

export function useDeleteForumReply(postId: string | undefined) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => forumService.deleteReply(id),
    onSuccess: () => {
      toast.success('回复已删除')
      queryClient.invalidateQueries({ queryKey: ['forum', 'posts'] })
      queryClient.invalidateQueries({ queryKey: ['forum', 'post', postId] })
      queryClient.invalidateQueries({ queryKey: ['forum', 'replies', postId] })
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败')
    },
  })
}
