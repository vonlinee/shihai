import { useState } from 'react'
import { Plus, Pencil, Trash2, Shield } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useRoles, useCreateRole, useUpdateRole, useDeleteRole } from '@/hooks/useRbac'

export function RolesPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [showDialog, setShowDialog] = useState(false)
  const [editingRoleId, setEditingRoleId] = useState<number | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    isActive: true,
  })

  const { data: rolesData, isLoading } = useRoles()
  const createRoleMutation = useCreateRole()
  const updateRoleMutation = useUpdateRole()
  const deleteRoleMutation = useDeleteRole()

  const roles = rolesData?.list ?? []

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (editingRoleId) {
      updateRoleMutation.mutate(
        { id: editingRoleId, data: formData },
        { onSuccess: () => { setShowDialog(false); resetForm() } },
      )
    } else {
      createRoleMutation.mutate(
        { name: formData.name, description: formData.description },
        { onSuccess: () => { setShowDialog(false); resetForm() } },
      )
    }
  }

  const handleDelete = async (role: { id: number; name: string }) => {
    if (!confirm(`确定要删除角色 "${role.name}" 吗？`)) return
    deleteRoleMutation.mutate(role.id)
  }

  const handleEdit = (role: { id: number; name: string; description: string; isActive: boolean }) => {
    setEditingRoleId(role.id)
    setFormData({
      name: role.name,
      description: role.description,
      isActive: role.isActive,
    })
    setShowDialog(true)
  }

  const handleAdd = () => {
    setEditingRoleId(null)
    resetForm()
    setShowDialog(true)
  }

  const handleManagePermissions = (role: { name: string }) => {
    // TODO: Open permission management dialog
    alert(`管理角色 "${role.name}" 的权限`)
  }

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      isActive: true,
    })
  }

  const filteredRoles = roles.filter(
    (role) =>
      role.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      role.description.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold font-serif">角色管理</h1>
          <p className="text-muted-foreground mt-1">管理系统角色和权限分配</p>
        </div>
        <Button onClick={handleAdd} className="gap-2">
          <Plus className="h-4 w-4" />
          添加角色
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>角色列表</CardTitle>
          <div className="mt-4">
            <Input
              placeholder="搜索角色名称或描述..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="max-w-sm"
            />
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8">加载中...</div>
          ) : filteredRoles.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              暂无角色数据
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">角色名称</th>
                    <th className="text-left py-3 px-4 font-medium">描述</th>
                    <th className="text-left py-3 px-4 font-medium">状态</th>
                    <th className="text-right py-3 px-4 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredRoles.map((role) => (
                    <tr key={role.id} className="border-b">
                      <td className="py-3 px-4 font-medium">{role.name}</td>
                      <td className="py-3 px-4">{role.description}</td>
                      <td className="py-3 px-4">
                        {role.isActive ? (
                          <span className="inline-flex items-center gap-1 text-green-600">
                            启用
                          </span>
                        ) : (
                          <span className="inline-flex items-center gap-1 text-red-600">
                            禁用
                          </span>
                        )}
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className="flex justify-end gap-2">
                          <Button
                            className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                            onClick={() => handleManagePermissions(role)}
                          >
                            <Shield className="h-4 w-4 mr-1" />
                            权限
                          </Button>
                          <Button
                            className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                            onClick={() => handleEdit(role)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                            onClick={() => handleDelete(role)}
                            disabled={role.name === 'admin'}
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
          )}
        </CardContent>
      </Card>

      {/* Role Dialog */}
      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-bold mb-2">
              {editingRoleId ? '编辑角色' : '添加角色'}
            </h2>
            <p className="text-muted-foreground mb-4">
              {editingRoleId ? '修改角色信息' : '创建新的系统角色'}
            </p>
            <form onSubmit={handleSubmit}>
              <div className="space-y-4 mb-6">
                <div className="space-y-2">
                  <label className="text-sm font-medium">角色名称</label>
                  <Input
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="如：editor"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">描述</label>
                  <Input
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    placeholder="角色描述..."
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
                  <label className="text-sm font-medium">启用角色</label>
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
                  disabled={createRoleMutation.isPending || updateRoleMutation.isPending}
                >
                  {editingRoleId ? '更新' : '创建'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
