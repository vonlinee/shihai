import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminService, type UserListParams, type PoemCreateRequest, type PoemUpdateRequest } from '@/services/adminService';
import { toast } from 'sonner';

export function useAdminUsers(params?: UserListParams) {
  return useQuery({
    queryKey: ['admin', 'users', params],
    queryFn: () => adminService.getUsers(params),
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => adminService.deleteUser(id),
    onSuccess: () => {
      toast.success('用户已删除');
      queryClient.invalidateQueries({ queryKey: ['admin', 'users'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

export function useAdminCreatePoem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: PoemCreateRequest) => adminService.createPoem(data),
    onSuccess: () => {
      toast.success('诗词添加成功');
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '添加失败');
    },
  });
}

export function useAdminUpdatePoem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: PoemUpdateRequest }) =>
      adminService.updatePoem(id, data),
    onSuccess: () => {
      toast.success('诗词更新成功');
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败');
    },
  });
}

export function useAdminDeletePoem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => adminService.deletePoem(id),
    onSuccess: () => {
      toast.success('诗词已删除');
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

export function useAdminComments(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['admin', 'comments', page, pageSize],
    queryFn: () => adminService.getComments(page, pageSize),
  });
}
