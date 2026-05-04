import { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, Plus, Trash2, X, Pencil } from 'lucide-react'
import { useRoles, useCreateRole, useUpdateRole, useDeleteRole } from '@/hooks/useRbac'
import { getRoleBadgeColor } from './utils'

export function RolesTab() {
  const [searchQuery, setSearchQuery] = useState('')
  const [showDialog, setShowDialog] = useState(false)
  const [editingRoleId, setEditingRoleId] = useState<number | null>(null)
  const [formData, setFormData] = useState({ name: '', description: '', isActive: true })

  const { data: rolesData, isLoading } = useRoles()
  const createMutation = useCreateRole()
  const updateMutation = useUpdateRole()
  const deleteMutation = useDeleteRole()

  const roles = rolesData?.list ?? []
  const filteredRoles = roles.filter((r) =>
    r.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    r.description.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const resetForm = () => setFormData({ name: '', description: '', isActive: true })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingRoleId) {
      updateMutation.mutate({ id: editingRoleId, data: formData },
        { onSuccess: () => { setShowDialog(false); resetForm() } })
    } else {
      createMutation.mutate({ name: formData.name, description: formData.description },
        { onSuccess: () => { setShowDialog(false); resetForm() } })
    }
  }

  return (
    <>
      <div className="flex justify-between items-center">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input placeholder="搜索角色..." className="pl-10" value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)} />
        </div>
        <Button onClick={() => { setEditingRoleId(null); resetForm(); setShowDialog(true) }}>
          <Plus className="h-4 w-4 mr-2" />添加角色
        </Button>
      </div>

      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">加载中...</div>
          ) : filteredRoles.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">暂无角色数据</div>
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
                    <tr key={role.id} className="border-b last:border-0 hover:bg-muted/50">
                      <td className="py-3 px-4">
                        <span className={`px-2 py-1 rounded text-xs ${getRoleBadgeColor(role.name)}`}>
                          {role.name}
                        </span>
                      </td>
                      <td className="py-3 px-4">{role.description}</td>
                      <td className="py-3 px-4">
                        <span className={`text-xs ${role.isActive ? 'text-green-600' : 'text-red-600'}`}>
                          {role.isActive ? '启用' : '禁用'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className="flex justify-end gap-2">
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                            onClick={() => { setEditingRoleId(role.id); setFormData({ name: role.name, description: role.description, isActive: role.isActive }); setShowDialog(true) }}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                            onClick={() => { if (confirm(`确定要删除角色 "${role.name}" 吗？`)) deleteMutation.mutate(role.id) }}
                            disabled={role.name === 'admin'}>
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
          <div className="bg-background rounded-lg w-full max-w-md max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingRoleId ? '编辑角色' : '添加角色'}</h2>
              <button onClick={() => setShowDialog(false)} className="p-1 hover:bg-muted rounded">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="space-y-2">
                <label className="text-sm font-medium">角色名称</label>
                <Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="如：editor" required />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">描述</label>
                <Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="角色描述..." />
              </div>
              <div className="flex items-center space-x-2">
                <input type="checkbox" checked={formData.isActive}
                  onChange={(e) => setFormData({ ...formData, isActive: e.target.checked })} className="rounded border-gray-300" />
                <label className="text-sm font-medium">启用角色</label>
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                  onClick={() => setShowDialog(false)}>取消</Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {createMutation.isPending || updateMutation.isPending ? '保存中...' : editingRoleId ? '更新' : '创建'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  )
}
