import { api } from './api';

// ─── Types ───────────────────────────────────────────────────────────────────

export interface Role {
  id: number;
  name: string;
  description: string;
  isActive: boolean;
  permissions?: string[];
}

export interface Permission {
  id: number;
  code: string;
  name: string;
  description: string;
  module: string;
  isActive: boolean;
}

export interface RoleListResponse {
  list: Role[];
  total: number;
  page: number;
  pageSize: number;
}

export interface PermissionListResponse {
  list: Permission[];
  total: number;
  page: number;
  pageSize: number;
}

export interface RoleCreateRequest {
  name: string;
  description?: string;
}

export interface RoleUpdateRequest {
  name?: string;
  description?: string;
  isActive?: boolean;
}

export interface PermissionCreateRequest {
  code: string;
  name: string;
  description?: string;
  module?: string;
}

export interface PermissionUpdateRequest {
  name?: string;
  description?: string;
  module?: string;
  isActive?: boolean;
}

export interface UserRoleResponse {
  userId: number;
  roles: Role[];
}

export interface RolePermissionResponse {
  roleId: number;
  permissions: Permission[];
}

// ─── Service ─────────────────────────────────────────────────────────────────

export const rbacService = {
  // Roles
  getRoles(page = 1, pageSize = 20): Promise<RoleListResponse> {
    return api.get<RoleListResponse>('/rbac/roles', { page, pageSize });
  },

  getRoleById(id: number): Promise<Role> {
    return api.get<Role>(`/rbac/roles/${id}`);
  },

  createRole(data: RoleCreateRequest): Promise<Role> {
    return api.post<Role>('/rbac/roles', data);
  },

  updateRole(id: number, data: RoleUpdateRequest): Promise<Role> {
    return api.put<Role>(`/rbac/roles/${id}`, data);
  },

  deleteRole(id: number): Promise<void> {
    return api.delete<void>(`/rbac/roles/${id}`);
  },

  getRolePermissions(roleId: number): Promise<RolePermissionResponse> {
    return api.get<RolePermissionResponse>(`/rbac/roles/${roleId}/permissions`);
  },

  assignPermissionsToRole(roleId: number, permissionIds: number[]): Promise<void> {
    return api.put<void>(`/rbac/roles/${roleId}/permissions`, { permissionIds });
  },

  // Permissions
  getPermissions(
    page = 1,
    pageSize = 50,
    module?: string,
  ): Promise<PermissionListResponse> {
    return api.get<PermissionListResponse>('/rbac/permissions', {
      page,
      pageSize,
      ...(module ? { module } : {}),
    });
  },

  getAllPermissions(): Promise<Permission[]> {
    return api.get<Permission[]>('/rbac/permissions/all');
  },

  getPermissionById(id: number): Promise<Permission> {
    return api.get<Permission>(`/rbac/permissions/${id}`);
  },

  createPermission(data: PermissionCreateRequest): Promise<Permission> {
    return api.post<Permission>('/rbac/permissions', data);
  },

  updatePermission(id: number, data: PermissionUpdateRequest): Promise<Permission> {
    return api.put<Permission>(`/rbac/permissions/${id}`, data);
  },

  deletePermission(id: number): Promise<void> {
    return api.delete<void>(`/rbac/permissions/${id}`);
  },

  // User Roles
  getUserRoles(userId: number): Promise<UserRoleResponse> {
    return api.get<UserRoleResponse>(`/rbac/users/${userId}/roles`);
  },

  assignRolesToUser(userId: number, roleIds: number[]): Promise<void> {
    return api.put<void>(`/rbac/users/${userId}/roles`, { roleIds });
  },

  getUserPermissions(userId: number): Promise<Permission[]> {
    return api.get<Permission[]>(`/rbac/users/${userId}/permissions`);
  },

  // Current user
  getMyPermissions(): Promise<Permission[]> {
    return api.get<Permission[]>('/rbac/my/permissions');
  },

  checkPermission(permission: string): Promise<{ hasPermission: boolean }> {
    return api.post<{ hasPermission: boolean }>('/rbac/check', { permission });
  },
};
