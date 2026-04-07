import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  rbacService,
  type RoleCreateRequest,
  type RoleUpdateRequest,
  type PermissionCreateRequest,
  type PermissionUpdateRequest,
} from '@/services/rbacService';
import { toast } from 'sonner';

// ─── Roles ───────────────────────────────────────────────────────────────────

export function useRoles(page = 1, pageSize = 20) {
  return useQuery({
    queryKey: ['rbac', 'roles', page, pageSize],
    queryFn: () => rbacService.getRoles(page, pageSize),
  });
}

export function useCreateRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: RoleCreateRequest) => rbacService.createRole(data),
    onSuccess: () => {
      toast.success('角色创建成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'roles'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '创建角色失败');
    },
  });
}

export function useUpdateRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: RoleUpdateRequest }) =>
      rbacService.updateRole(id, data),
    onSuccess: () => {
      toast.success('角色更新成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'roles'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新角色失败');
    },
  });
}

export function useDeleteRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => rbacService.deleteRole(id),
    onSuccess: () => {
      toast.success('角色删除成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'roles'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除角色失败');
    },
  });
}

export function useRolePermissions(roleId: number) {
  return useQuery({
    queryKey: ['rbac', 'roles', roleId, 'permissions'],
    queryFn: () => rbacService.getRolePermissions(roleId),
    enabled: !!roleId,
  });
}

export function useAssignPermissionsToRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ roleId, permissions }: { roleId: number; permissions: string[] }) =>
      rbacService.assignPermissionsToRole(roleId, permissions),
    onSuccess: () => {
      toast.success('权限分配成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'roles'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '权限分配失败');
    },
  });
}

// ─── Permissions ─────────────────────────────────────────────────────────────

export function usePermissions(page = 1, pageSize = 50, module?: string) {
  return useQuery({
    queryKey: ['rbac', 'permissions', page, pageSize, module],
    queryFn: () => rbacService.getPermissions(page, pageSize, module),
  });
}

export function useAllPermissions() {
  return useQuery({
    queryKey: ['rbac', 'permissions', 'all'],
    queryFn: () => rbacService.getAllPermissions(),
  });
}

export function useCreatePermission() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: PermissionCreateRequest) => rbacService.createPermission(data),
    onSuccess: () => {
      toast.success('权限创建成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'permissions'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '创建权限失败');
    },
  });
}

export function useUpdatePermission() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: PermissionUpdateRequest }) =>
      rbacService.updatePermission(id, data),
    onSuccess: () => {
      toast.success('权限更新成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'permissions'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '更新权限失败');
    },
  });
}

export function useDeletePermission() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => rbacService.deletePermission(id),
    onSuccess: () => {
      toast.success('权限删除成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'permissions'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '删除权限失败');
    },
  });
}

// ─── User Roles ──────────────────────────────────────────────────────────────

export function useUserRoles(userId: number) {
  return useQuery({
    queryKey: ['rbac', 'users', userId, 'roles'],
    queryFn: () => rbacService.getUserRoles(userId),
    enabled: !!userId,
  });
}

export function useAssignRolesToUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, roleIds }: { userId: number; roleIds: number[] }) =>
      rbacService.assignRolesToUser(userId, roleIds),
    onSuccess: () => {
      toast.success('角色分配成功');
      queryClient.invalidateQueries({ queryKey: ['rbac', 'users'] });
      queryClient.invalidateQueries({ queryKey: ['admin', 'users'] });
    },
    onError: (error: Error) => {
      toast.error(error.message || '角色分配失败');
    },
  });
}

export function useMyPermissions() {
  return useQuery({
    queryKey: ['rbac', 'my', 'permissions'],
    queryFn: () => rbacService.getMyPermissions(),
  });
}
