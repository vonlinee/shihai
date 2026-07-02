import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  adminService,
  type UserListParams,
  type PoemCreateRequest,
  type PoemUpdateRequest,
  type WorkCollectionListParams,
  type WorkCollectionCreateRequest,
  type WorkCollectionUpdateRequest,
  type WorkCollectionItemCreateRequest,
  type WorkCollectionItemUpdateRequest,
  type DynastyCreateRequest,
  type DynastyUpdateRequest,
  type PoetCreateRequest,
  type PoetUpdateRequest,
  type AdminCreateUserRequest,
} from '@/services/adminService';
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

export function useAdminCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AdminCreateUserRequest) => adminService.createUser(data),
    onSuccess: () => {
      toast.success('用户创建成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'users'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '创建失败');
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

// ─── Dynasty Admin ──────────────────────────────────────────────────────

export function useAdminWorkCollections(params?: WorkCollectionListParams) {
  return useQuery({
    queryKey: ['admin', 'workCollections', params],
    queryFn: () => adminService.getWorkCollections(params),
  });
}

export function useAdminWorkCollection(id: number | null) {
  return useQuery({
    queryKey: ['admin', 'workCollection', id],
    queryFn: () => adminService.getWorkCollectionById(id!),
    enabled: !!id,
  });
}

export function useAdminCreateWorkCollection() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: WorkCollectionCreateRequest) => adminService.createWorkCollection(data),
    onSuccess: () => {
      toast.success('作品集创建成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollections'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '创建失败');
    },
  });
}

export function useAdminUpdateWorkCollection() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: WorkCollectionUpdateRequest }) =>
      adminService.updateWorkCollection(id, data),
    onSuccess: (_, variables) => {
      toast.success('作品集更新成功');
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollections'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollection', variables.id] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败');
    },
  });
}

export function useAdminDeleteWorkCollection() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => adminService.deleteWorkCollection(id),
    onSuccess: () => {
      toast.success('作品集已删除');
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollections'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

export function useAdminAddWorkCollectionItem() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ collectionId, data }: { collectionId: number; data: WorkCollectionItemCreateRequest }) =>
      adminService.addWorkCollectionItem(collectionId, data),
    onSuccess: (_, variables) => {
      toast.success('作品已添加');
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollections'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollection', variables.collectionId] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '添加失败');
    },
  });
}

export function useAdminUpdateWorkCollectionItem() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ collectionId, itemId, data }: { collectionId: number; itemId: number; data: WorkCollectionItemUpdateRequest }) =>
      adminService.updateWorkCollectionItem(collectionId, itemId, data),
    onSuccess: (_, variables) => {
      toast.success('排序已更新');
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollections'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollection', variables.collectionId] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败');
    },
  });
}

export function useAdminDeleteWorkCollectionItem() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ collectionId, itemId }: { collectionId: number; itemId: number }) =>
      adminService.deleteWorkCollectionItem(collectionId, itemId),
    onSuccess: (_, variables) => {
      toast.success('作品已移除');
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollections'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'workCollection', variables.collectionId] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '移除失败');
    },
  });
}

export function useAdminCreateDynasty() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: DynastyCreateRequest) => adminService.createDynasty(data),
    onSuccess: () => {
      toast.success('朝代添加成功');
      queryClient.invalidateQueries({ queryKey: ['dynasties'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '添加失败');
    },
  });
}

export function useAdminUpdateDynasty() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: DynastyUpdateRequest }) =>
      adminService.updateDynasty(id, data),
    onSuccess: () => {
      toast.success('朝代更新成功');
      queryClient.invalidateQueries({ queryKey: ['dynasties'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败');
    },
  });
}

export function useAdminDeleteDynasty() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => adminService.deleteDynasty(id),
    onSuccess: () => {
      toast.success('朝代已删除');
      queryClient.invalidateQueries({ queryKey: ['dynasties'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

// ─── Poet Admin ──────────────────────────────────────────────────────

export function useAdminCreatePoet() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: PoetCreateRequest) => adminService.createPoet(data),
    onSuccess: () => {
      toast.success('诗人添加成功');
      queryClient.invalidateQueries({ queryKey: ['poets'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '添加失败');
    },
  });
}

export function useAdminUpdatePoet() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: PoetUpdateRequest }) =>
      adminService.updatePoet(id, data),
    onSuccess: () => {
      toast.success('诗人更新成功');
      queryClient.invalidateQueries({ queryKey: ['poets'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败');
    },
  });
}

export function useAdminDeletePoet() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => adminService.deletePoet(id),
    onSuccess: () => {
      toast.success('诗人已删除');
      queryClient.invalidateQueries({ queryKey: ['poets'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}
