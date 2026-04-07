import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { commentService, type CreateCommentRequest, type VoteCommentRequest } from '@/services/commentService';
import { toast } from 'sonner';

export function useComments(poemId: number, page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['comments', poemId, page, pageSize],
    queryFn: () => commentService.getComments(poemId, page, pageSize),
    enabled: !!poemId,
  });
}

export function useCreateComment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCommentRequest) => commentService.createComment(data),
    onSuccess: () => {
      toast.success('评论发表成功');
      queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '评论发表失败');
    },
  });
}

export function useDeleteComment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => commentService.deleteComment(id),
    onSuccess: () => {
      toast.success('评论已删除');
      queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

export function useVoteComment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: VoteCommentRequest) => commentService.voteComment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '操作失败');
    },
  });
}
