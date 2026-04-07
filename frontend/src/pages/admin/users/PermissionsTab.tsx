import { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Search, Plus, Trash2, X, Pencil } from 'lucide-react'
import { usePermissions, useCreatePermission, useUpdatePermission, useDeletePermission } from '@/hooks/useRbac'

export function PermissionsTab() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedModule, setSelectedModule] = useState('')
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [formData, setFormData] = useState({ code: '', name: '', description: '', module: '', isActive: true })

  const { data: permissionsData, isLoading } = usePermissions(1, 100, selectedModule || undefined)
  const createMutation = useCreatePermission()
  const updateMutation = useUpdatePermission()
  const deleteMutation = useDeletePermission()

  const permissions = permissionsData?.list ?? []
  const filteredPermissions = permissions.filter((p) =>
    p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    p.code.toLowerCase().includes(searchQuery.toLowerCase())
  )
  const modules = Array.from(new Set(permissions.map((p) => p.module)))
  const grouped = filteredPermissions.reduce((acc, perm) => {
    if (!acc[perm.module]) acc[perm.module] = []
    acc[perm.module].push(perm)
    return acc
  }, {} as Record<string, typeof permissions>)

  const resetForm = () => setFormData({ code: '', name: '', description: '', module: '', isActive: true })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingId) {
      updateMutation.mutate(
        { id: editingId, data: { name: formData.name, description: formData.description, module: formData.module, isActive: formData.isActive } },
        { onSuccess: () => { setShowDialog(false); resetForm() } },
      )
    } else {
      createMutation.mutate(
        { code: formData.code, name: formData.name, description: formData.description, module: formData.module },
        { onSuccess: () => { setShowDialog(false); resetForm() } },
      )
    }
  }

  return (
    <>
      <div className="flex justify-between items-center">
        <div className="flex gap-4">
          <div className="relative max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input placeholder="搜索权限..." className="pl-10" value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)} />
          </div>
          <select value={selectedModule} onChange={(e) => setSelectedModule(e.target.value)}
            className="border rounded-md px-3 py-2 bg-background text-sm">
            <option value="">所有模块</option>
            {modules.map((m) => (<option key={m} value={m}>{m}</option>))}
          </select>
        </div>
        <Button onClick={() => { setEditingId(null); resetForm(); setShowDialog(true) }}>
          <Plus className="h-4 w-4 mr-2" />添加权限
        </Button>
      </div>

      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">加载中...</div>
          ) : filteredPermissions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">暂无权限数据</div>
          ) : selectedModule ? (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">编码</th>
                    <th className="text-left py-3 px-4 font-medium">名称</th>
                    <th className="text-left py-3 px-4 font-medium">描述</th>
                    <th className="text-left py-3 px-4 font-medium">状态</th>
                    <th className="text-right py-3 px-4 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredPermissions.map((perm) => (
                    <tr key={perm.id} className="border-b last:border-0 hover:bg-muted/50">
                      <td className="py-3 px-4 font-mono text-sm">{perm.code}</td>
                      <td className="py-3 px-4 font-medium">{perm.name}</td>
                      <td className="py-3 px-4">{perm.description}</td>
                      <td className="py-3 px-4">
                        <span className={`text-xs ${perm.isActive ? 'text-green-600' : 'text-red-600'}`}>
                          {perm.isActive ? '启用' : '禁用'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className="flex justify-end gap-2">
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                            onClick={() => { setEditingId(perm.id); setFormData({ code: perm.code, name: perm.name, description: perm.description, module: perm.module, isActive: perm.isActive }); setShowDialog(true) }}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                            onClick={() => { if (confirm(`确定要删除权限 "${perm.name}" 吗？`)) deleteMutation.mutate(perm.id) }}>
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
            <div className="space-y-6">
              {Object.entries(grouped).map(([module, perms]) => (
                <div key={module}>
                  <h3 className="font-medium text-lg mb-3 capitalize">{module} 模块</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {perms.map((perm) => (
                      <div key={perm.id} className="border rounded-lg p-4 hover:shadow-md transition-shadow">
                        <div className="flex justify-between items-start mb-2">
                          <div>
                            <div className="font-medium">{perm.name}</div>
                            <div className="text-xs text-muted-foreground font-mono">{perm.code}</div>
                          </div>
                          <div className="flex gap-1">
                            <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                              onClick={() => { setEditingId(perm.id); setFormData({ code: perm.code, name: perm.name, description: perm.description, module: perm.module, isActive: perm.isActive }); setShowDialog(true) }}>
                              <Pencil className="h-3 w-3" />
                            </Button>
                            <Button className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                              onClick={() => { if (confirm(`确定要删除权限 "${perm.name}" 吗？`)) deleteMutation.mutate(perm.id) }}>
                              <Trash2 className="h-3 w-3" />
                            </Button>
                          </div>
                        </div>
                        <p className="text-sm text-muted-foreground">{perm.description}</p>
                        <div className="mt-2">
                          <span className={`text-xs ${perm.isActive ? 'text-green-600' : 'text-red-600'}`}>
                            {perm.isActive ? '启用' : '禁用'}
                          </span>
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
          <div className="bg-background rounded-lg w-full max-w-md max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingId ? '编辑权限' : '添加权限'}</h2>
              <button onClick={() => setShowDialog(false)} className="p-1 hover:bg-muted rounded">
                <X className="h-5 w-5" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="space-y-2">
                <label className="text-sm font-medium">权限编码</label>
                <Input value={formData.code} onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                  placeholder="如：user:create" required disabled={!!editingId} />
                <p className="text-xs text-muted-foreground">格式：模块:操作，如 user:create</p>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">权限名称</label>
                <Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="如：创建用户" required />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">所属模块</label>
                <Input value={formData.module} onChange={(e) => setFormData({ ...formData, module: e.target.value })}
                  placeholder="如：user, poem, role" required />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">描述</label>
                <Input value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="权限描述..." />
              </div>
              <div className="flex items-center space-x-2">
                <input type="checkbox" checked={formData.isActive}
                  onChange={(e) => setFormData({ ...formData, isActive: e.target.checked })} className="rounded border-gray-300" />
                <label className="text-sm font-medium">启用权限</label>
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                  onClick={() => setShowDialog(false)}>取消</Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {createMutation.isPending || updateMutation.isPending ? '保存中...' : editingId ? '更新' : '创建'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  )
}
