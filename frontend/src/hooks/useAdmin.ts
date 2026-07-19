import { useQuery, useMutation, useQueryClient, type QueryClient } from '@tanstack/react-query';
import {
  adminService,
  type UserListParams,
  type PoemCreateRequest,
  type PoemUpdateRequest,
  type PoemAnnotationUpsertRequest,
  type TextConversionRequest,
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
  type CorrectionListParams,
} from '@/services/adminService';
import { toast } from 'sonner';

type AdminID = string | number;

function syncSavedPoemQueryCache(
  queryClient: QueryClient,
  poem: { id?: AdminID } | null | undefined,
  fallbackId?: AdminID,
) {
  if (!poem) return;

  const ids = Array.from(new Set([poem.id, fallbackId].filter((value): value is AdminID => value !== undefined && value !== null)));
  ids.forEach((id) => {
    queryClient.setQueryData(['poem', String(id)], poem);
  });
}

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
    onSuccess: (poem) => {
      toast.success('诗词添加成功');
      queryClient.invalidateQueries({ queryKey: ['poems'] });
      queryClient.invalidateQueries({ queryKey: ['poem'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'poemAnnotations'] });
      syncSavedPoemQueryCache(queryClient, poem);
      if (poem?.id) {
        queryClient.setQueryData(['admin', 'poemAnnotations', String(poem.id)], poem.annotations ?? []);
      }
    },
    onError: (error: Error) => {
      toast.error(error.message || '添加失败');
    },
  });
}

export function useAdminUpdatePoem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: AdminID; data: PoemUpdateRequest }) =>
      adminService.updatePoem(id, data),
    onSuccess: (poem, variables) => {
      toast.success('诗词更新成功');
      queryClient.invalidateQueries({ queryKey: ['poems'] });
      queryClient.invalidateQueries({ queryKey: ['poem'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'poemAnnotations', variables.id] });
      syncSavedPoemQueryCache(queryClient, poem, variables.id);
      queryClient.setQueryData(['admin', 'poemAnnotations', String(variables.id)], poem.annotations ?? []);
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新失败');
    },
  });
}

export function useAdminDeletePoem() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: AdminID) => adminService.deletePoem(id),
    onSuccess: () => {
      toast.success('诗词已删除');
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

export function useAdminConvertTexts() {
  return useMutation({
    mutationFn: (data: TextConversionRequest) => adminService.convertTexts(data),
    onError: (error: Error) => {
      toast.error(error.message || '文本转换失败');
    },
  });
}

export function useAdminComments(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['admin', 'comments', page, pageSize],
    queryFn: () => adminService.getComments(page, pageSize),
  });
}

export function useAdminCorrections(params?: CorrectionListParams) {
  return useQuery({
    queryKey: ['admin', 'corrections', params],
    queryFn: () => adminService.getCorrections(params),
  });
}

// ─── Dynasty Admin ──────────────────────────────────────────────────────

export function useAdminWorkCollections(params?: WorkCollectionListParams) {
  return useQuery({
    queryKey: ['admin', 'workCollections', params],
    queryFn: () => adminService.getWorkCollections(params),
  });
}

export function useAdminBatchDeletePoems() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ids: AdminID[]) => adminService.batchDeletePoems(ids),
    onSuccess: (_, ids) => {
      toast.success(`已删除 ${ids.length} 首诗词`);
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '批量删除失败');
    },
  });
}

export function useAdminPoemAnnotations(poemId: AdminID | null) {
  return useQuery({
    queryKey: ['admin', 'poemAnnotations', poemId],
    queryFn: () => adminService.getPoemAnnotations(poemId!),
    enabled: !!poemId,
  });
}

export function useAdminCreatePoemAnnotation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ poemId, data }: { poemId: AdminID; data: PoemAnnotationUpsertRequest }) =>
      adminService.createPoemAnnotation(poemId, data),
    onSuccess: (_, variables) => {
      toast.success('标注已添加');
      queryClient.invalidateQueries({ queryKey: ['admin', 'poemAnnotations', variables.poemId] });
      queryClient.invalidateQueries({ queryKey: ['poem'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '添加标注失败');
    },
  });
}

export function useAdminUpdatePoemAnnotation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (variables: { id: AdminID; poemId: AdminID; data: PoemAnnotationUpsertRequest }) =>
      adminService.updatePoemAnnotation(variables.id, variables.data),
    onSuccess: (_, variables) => {
      toast.success('标注已更新');
      queryClient.invalidateQueries({ queryKey: ['admin', 'poemAnnotations', variables.poemId] });
      queryClient.invalidateQueries({ queryKey: ['poem'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新标注失败');
    },
  });
}

export function useAdminDeletePoemAnnotation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (variables: { id: AdminID; poemId: AdminID }) => adminService.deletePoemAnnotation(variables.id),
    onSuccess: (_, variables) => {
      toast.success('标注已删除');
      queryClient.invalidateQueries({ queryKey: ['admin', 'poemAnnotations', variables.poemId] });
      queryClient.invalidateQueries({ queryKey: ['poem'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除标注失败');
    },
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
    mutationFn: ({ id, data }: { id: AdminID; data: DynastyUpdateRequest }) =>
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
    mutationFn: (id: AdminID) => adminService.deleteDynasty(id),
    onSuccess: () => {
      toast.success('朝代已删除');
      queryClient.invalidateQueries({ queryKey: ['dynasties'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

export function useAdminBatchDeleteDynasties() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (ids: AdminID[]) => adminService.batchDeleteDynasties(ids),
    onSuccess: (_, ids) => {
      toast.success(`已删除 ${ids.length} 个朝代`);
      queryClient.invalidateQueries({ queryKey: ['dynasties'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '批量删除失败');
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
    mutationFn: ({ id, data }: { id: AdminID; data: PoetUpdateRequest }) =>
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
    mutationFn: (id: AdminID) => adminService.deletePoet(id),
    onSuccess: () => {
      toast.success('诗人已删除');
      queryClient.invalidateQueries({ queryKey: ['poets'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除失败');
    },
  });
}

export function useAdminBatchDeletePoets() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (ids: AdminID[]) => adminService.batchDeletePoets(ids),
    onSuccess: (_, ids) => {
      toast.success(`已删除 ${ids.length} 位诗人`);
      queryClient.invalidateQueries({ queryKey: ['poets'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '批量删除失败');
    },
  });
}
