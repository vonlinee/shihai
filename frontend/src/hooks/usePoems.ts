import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { poemService } from '@/services/poemService';
import type { PoemSearchParams } from '@/types';
import { toast } from 'sonner';

export function usePoems(params?: PoemSearchParams) {
  return useQuery({
    queryKey: ['poems', params],
    queryFn: () => poemService.getPoems(params),
  });
}

export function usePoem(id: number) {
  return useQuery({
    queryKey: ['poem', id],
    queryFn: () => poemService.getPoemById(id),
    enabled: !!id,
  });
}

export function useRandomPoems(limit = 5) {
  return useQuery({
    queryKey: ['poems', 'random', limit],
    queryFn: () => poemService.getRandomPoems(limit),
  });
}

export function useLikePoem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => poemService.likePoem(id),
    onSuccess: () => {
      toast.success('点赞成功');
      queryClient.invalidateQueries({ queryKey: ['poems'] });
      queryClient.invalidateQueries({ queryKey: ['poem'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '操作失败');
    },
  });
}

export function useDynasties() {
  return useQuery({
    queryKey: ['dynasties'],
    queryFn: () => poemService.getDynasties(),
    staleTime: 1000 * 60 * 30, // dynasties rarely change
  });
}

export function usePoets(keyword?: string) {
  return useQuery({
    queryKey: ['poets', keyword],
    queryFn: () => poemService.getPoets(keyword),
    staleTime: 1000 * 60 * 5,
  });
}

export function useGenres() {
  return useQuery({
    queryKey: ['genres'],
    queryFn: () => poemService.getGenres(),
    staleTime: 1000 * 60 * 30,
  });
}
