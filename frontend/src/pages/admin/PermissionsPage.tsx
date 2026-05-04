import { useState } from 'react'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { usePermissions, useCreatePermission, useUpdatePermission, useDeletePermission } from '@/hooks/useRbac'

export function PermissionsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedModule, setSelectedModule] = useState('')
  const [showDialog, setShowDialog] = useState(false)
  const [editingPermissionId, setEditingPermissionId] = useState<number | null>(null)
  const [formData, setFormData] = useState({
    code: '',
    name: '',
    description: '',
    module: '',
    isActive: true,
  })

  const { data: permissionsData, isLoading } = usePermissions(1, 100, selectedModule || undefined)
  const createPermissionMutation = useCreatePermission()
  const updatePermissionMutation = useUpdatePermission()
  const deletePermissionMutation = useDeletePermission()

  const permissions = permissionsData?.list ?? []

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (editingPermissionId) {
      updatePermissionMutation.mutate(
        {
          id: editingPermissionId,
          data: {
            name: formData.name,
            description: formData.description,
            module: formData.module,
            isActive: formData.isActive,
          },
        },
        { onSuccess: () => { setShowDialog(false); resetForm() } },
      )
    } else {
      createPermissionMutation.mutate(
        {
          code: formData.code,
          name: formData.name,
          description: formData.description,
          module: formData.module,
        },
        { onSuccess: () => { setShowDialog(false); resetForm() } },
      )
    }
  }

  const handleDelete = async (permission: { id: number; name: string }) => {
    if (!confirm(`确定要删除权限 "${permission.name}" 吗？`)) return
    deletePermissionMutation.mutate(permission.id)
  }

  const handleEdit = (permission: {
    id: number;
    code: string;
    name: string;
    description: string;
    module: string;
    isActive: boolean;
  }) => {
    setEditingPermissionId(permission.id)
    setFormData({
      code: permission.code,
      name: permission.name,
      description: permission.description,
      module: permission.module,
      isActive: permission.isActive,
    })
    setShowDialog(true)
  }

  const handleAdd = () => {
    setEditingPermissionId(null)
    resetForm()
    setShowDialog(true)
  }

  const resetForm = () => {
    setFormData({
      code: '',
      name: '',
      description: '',
      module: '',
      isActive: true,
    })
  }

  // Get unique modules
  const modules = Array.from(new Set(permissions.map((p) => p.module)))

  // Group permissions by module
  const groupedPermissions = permissions.reduce((acc, perm) => {
    if (!acc[perm.module]) acc[perm.module] = []
    acc[perm.module].push(perm)
    return acc
  }, {} as Record<string, typeof permissions>)

  // Filter permissions
  const filteredPermissions = permissions.filter(
    (perm) =>
      (perm.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        perm.code.toLowerCase().includes(searchQuery.toLowerCase()))
  )

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold font-serif">权限管理</h1>
          <p className="text-muted-foreground mt-1">管理系统权限和模块配置</p>
        </div>
        <Button onClick={handleAdd} className="gap-2">
          <Plus className="h-4 w-4" />
          添加权限
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>权限列表</CardTitle>
          <div className="mt-4 flex gap-4">
            <Input
              placeholder="搜索权限名称或编码..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="max-w-sm"
            />
            <select
              value={selectedModule}
              onChange={(e) => setSelectedModule(e.target.value)}
              className="border rounded-md px-3 py-2 bg-background"
            >
              <option value="">所有模块</option>
              {modules.map((module) => (
                <option key={module} value={module}>
                  {module}
                </option>
              ))}
            </select>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8">加载中...</div>
          ) : filteredPermissions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              暂无权限数据
            </div>
          ) : selectedModule ? (
            // Show table view when module is selected
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">权限编码</th>
                    <th className="text-left py-3 px-4 font-medium">权限名称</th>
                    <th className="text-left py-3 px-4 font-medium">描述</th>
                    <th className="text-left py-3 px-4 font-medium">状态</th>
                    <th className="text-right py-3 px-4 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredPermissions.map((perm) => (
                    <tr key={perm.id} className="border-b">
                      <td className="py-3 px-4 font-mono text-sm">{perm.code}</td>
                      <td className="py-3 px-4 font-medium">{perm.name}</td>
                      <td className="py-3 px-4">{perm.description}</td>
                      <td className="py-3 px-4">
                        {perm.isActive ? (
                          <span className="text-green-600">启用</span>
                        ) : (
                          <span className="text-red-600">禁用</span>
                        )}
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className="flex justify-end gap-2">
                          <Button
                            className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                            onClick={() => handleEdit(perm)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                            onClick={() => handleDelete(perm)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            // Show grouped view when no module selected
            <div className="space-y-6">
              {Object.entries(groupedPermissions).map(([module, perms]) => (
                <div key={module}>
                  <h3 className="font-medium text-lg mb-3 capitalize">{module} 模块</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {perms.map((perm) => (
                      <div
                        key={perm.id}
                        className="border rounded-lg p-4 hover:shadow-md transition-shadow"
                      >
                        <div className="flex justify-between items-start mb-2">
                          <div>
                            <div className="font-medium">{perm.name}</div>
                            <div className="text-xs text-muted-foreground font-mono">
                              {perm.code}
                            </div>
                          </div>
                          <div className="flex gap-1">
                            <Button
                              className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                              onClick={() => handleEdit(perm)}
                            >
                              <Pencil className="h-3 w-3" />
                            </Button>
                            <Button
                              className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                              onClick={() => handleDelete(perm)}
                            >
                              <Trash2 className="h-3 w-3" />
                            </Button>
                          </div>
                        </div>
                        <p className="text-sm text-muted-foreground">{perm.description}</p>
                        <div className="mt-2">
                          {perm.isActive ? (
                            <span className="text-xs text-green-600">启用</span>
                          ) : (
                            <span className="text-xs text-red-600">禁用</span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Permission Dialog */}
      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-2">
              {editingPermissionId ? '编辑权限' : '添加权限'}
            </h2>
            <p className="text-muted-foreground mb-4">
              {editingPermissionId ? '修改权限信息' : '创建新的系统权限'}
            </p>
            <form onSubmit={handleSubmit}>
              <div className="space-y-4 mb-6">
                <div className="space-y-2">
                  <label className="text-sm font-medium">权限编码</label>
                  <Input
                    value={formData.code}
                    onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                    placeholder="如：user:create"
                    required
                    disabled={!!editingPermissionId}
                  />
                  <p className="text-xs text-muted-foreground">
                    格式：模块:操作，如 user:create, poem:delete
                  </p>
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">权限名称</label>
                  <Input
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="如：创建用户"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">所属模块</label>
                  <Input
                    value={formData.module}
                    onChange={(e) => setFormData({ ...formData, module: e.target.value })}
                    placeholder="如：user, poem, role"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">描述</label>
                  <Input
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    placeholder="权限描述..."
                  />
                </div>
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={formData.isActive}
                    onChange={(e) =>
                      setFormData({ ...formData, isActive: e.target.checked })
                    }
                    className="rounded border-gray-300"
                  />
                  <label className="text-sm font-medium">启用权限</label>
                </div>
              </div>
              <div className="flex justify-end gap-2">
                <Button
                  type="button"
                  className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                  onClick={() => setShowDialog(false)}
                >
                  取消
                </Button>
                <Button
                  type="submit"
                  disabled={createPermissionMutation.isPending || updatePermissionMutation.isPending}
                >
                  {editingPermissionId ? '更新' : '创建'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
