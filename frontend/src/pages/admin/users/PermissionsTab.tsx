import { useMemo, useState, Fragment } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Search, Plus, Trash2, X, Pencil, ChevronDown, ChevronRight, FolderTree } from 'lucide-react'
import { usePermissions, useCreatePermission, useUpdatePermission, useDeletePermission } from '@/hooks/useRbac'
import { getModuleDisplayName } from './utils'

// 单条权限数据结构（与后端 dto.PermissionResponse 保持一致）
interface PermissionRow {
  id: number
  code: string
  name: string
  description: string
  module: string
  isActive: boolean
}

/**
 * 权限管理 Tab
 *
 * 展示形态：单一表格，按模块层级以树形结构呈现。
 *   - 父行：模块（中文名称 + 英文 code + 权限数量），可点击展开/收起
 *   - 子行：模块下的权限点，缩进展示编码、名称、描述、状态、操作
 *
 * 交互说明：
 *   - 默认所有模块展开，方便快速浏览全部权限点
 *   - 顶部支持按权限名称/编码搜索，搜索时自动展开匹配模块
 *   - 模块过滤下拉框用于快速定位单一模块
 */
export function PermissionsTab() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedModule, setSelectedModule] = useState('')
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [formData, setFormData] = useState({ code: '', name: '', description: '', module: '', isActive: true })
  // 折叠状态：记录被折叠的模块名称，未在集合中的模块默认展开
  const [collapsedModules, setCollapsedModules] = useState<Set<string>>(new Set())

  const { data: permissionsData, isLoading } = usePermissions(1, 100, selectedModule || undefined)
  const createMutation = useCreatePermission()
  const updateMutation = useUpdatePermission()
  const deleteMutation = useDeletePermission()
  const { confirm, ConfirmDialog } = useConfirmDialog()

  const permissions: PermissionRow[] = permissionsData?.list ?? []

  // 模块下拉框选项（基于全量权限去重）
  const modules = useMemo(() => Array.from(new Set(permissions.map((p) => p.module))), [permissions])

  // 搜索过滤：按权限名称或编码模糊匹配
  const filteredPermissions = useMemo(
    () =>
      permissions.filter(
        (p) =>
          p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          p.code.toLowerCase().includes(searchQuery.toLowerCase()),
      ),
    [permissions, searchQuery],
  )

  // 按模块分组（保留原始顺序），用于树形渲染
  const groupedByModule = useMemo(() => {
    const map = new Map<string, PermissionRow[]>()
    for (const perm of filteredPermissions) {
      if (!map.has(perm.module)) map.set(perm.module, [])
      map.get(perm.module)!.push(perm)
    }
    return Array.from(map.entries())
  }, [filteredPermissions])

  // 切换模块展开/收起状态
  const toggleModule = (module: string) => {
    setCollapsedModules((prev) => {
      const next = new Set(prev)
      if (next.has(module)) next.delete(module)
      else next.add(module)
      return next
    })
  }

  // 一键展开/收起全部模块
  const allCollapsed = collapsedModules.size > 0 && collapsedModules.size === groupedByModule.length
  const toggleAll = () => {
    if (allCollapsed) setCollapsedModules(new Set())
    else setCollapsedModules(new Set(groupedByModule.map(([m]) => m)))
  }

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

  const startEdit = (perm: PermissionRow) => {
    setEditingId(perm.id)
    setFormData({ code: perm.code, name: perm.name, description: perm.description, module: perm.module, isActive: perm.isActive })
    setShowDialog(true)
  }

  const handleDelete = async (perm: PermissionRow) => {
    const confirmed = await confirm({
      title: '删除权限',
      description: `确定要删除权限 "${perm.name}" 吗？此操作不可撤销。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteMutation.mutate(perm.id)
  }

  return (
    <>
      {/* 顶部工具栏：搜索框 + 模块下拉 + 折叠开关 + 添加按钮 */}
      <div className="flex justify-between items-center">
        <div className="flex gap-4">
          <div className="relative max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="搜索权限名称或编码..."
              className="pl-10"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          <select
            value={selectedModule}
            onChange={(e) => setSelectedModule(e.target.value)}
            className="border rounded-md px-3 py-2 bg-background text-sm"
          >
            <option value="">所有模块</option>
            {modules.map((m) => (
              <option key={m} value={m}>{getModuleDisplayName(m)}</option>
            ))}
          </select>
          <Button
            className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
            onClick={toggleAll}
            disabled={groupedByModule.length === 0}
          >
            <FolderTree className="h-4 w-4 mr-2" />
            {allCollapsed ? '全部展开' : '全部收起'}
          </Button>
        </div>
        <Button onClick={() => { setEditingId(null); resetForm(); setShowDialog(true) }}>
          <Plus className="h-4 w-4 mr-2" />添加权限
        </Button>
      </div>

      {/* 权限树形表格 */}
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">加载中...</div>
          ) : filteredPermissions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">暂无权限数据</div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/30">
                  <TableHead className="w-[28%]">模块 / 权限编码</TableHead>
                  <TableHead className="w-[16%]">权限名称</TableHead>
                  <TableHead>描述</TableHead>
                  <TableHead className="w-[80px]">状态</TableHead>
                  <TableHead className="w-[120px] text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                  {groupedByModule.map(([module, perms]) => {
                    const collapsed = collapsedModules.has(module)
                    return (
                      <Fragment key={module}>
                        {/* 父行：模块层级 */}
                        <TableRow
                          key={`module-${module}`}
                          className="border-b bg-muted/20 hover:bg-muted/40 cursor-pointer"
                          onClick={() => toggleModule(module)}
                        >
                          <TableCell className="font-medium">
                            <div className="flex items-center gap-2">
                              {collapsed ? (
                                <ChevronRight className="h-4 w-4 text-muted-foreground" />
                              ) : (
                                <ChevronDown className="h-4 w-4 text-muted-foreground" />
                              )}
                              <span>{getModuleDisplayName(module)}</span>
                              <span className="text-xs font-mono text-muted-foreground">({module})</span>
                            </div>
                          </TableCell>
                          <TableCell className="text-sm text-muted-foreground" colSpan={3}>
                            共 {perms.length} 个权限点
                          </TableCell>
                          <TableCell />
                        </TableRow>

                        {/* 子行：模块下的权限点 */}
                        {!collapsed && perms.map((perm) => (
                          <TableRow key={perm.id} className="hover:bg-muted/30">
                            <TableCell className="pl-12 font-mono text-sm text-muted-foreground">
                              {perm.code}
                            </TableCell>
                            <TableCell className="font-medium">{perm.name}</TableCell>
                            <TableCell className="text-sm text-muted-foreground">{perm.description}</TableCell>
                            <TableCell>
                              <span className={`text-xs ${perm.isActive ? 'text-green-600' : 'text-red-600'}`}>
                                {perm.isActive ? '启用' : '禁用'}
                              </span>
                            </TableCell>
                            <TableCell className="text-right">
                              <div className="flex justify-end gap-2">
                                <Button
                                  className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                                  onClick={() => startEdit(perm)}
                                >
                                  <Pencil className="h-4 w-4" />
                                </Button>
                                <Button
                                  className="bg-secondary text-secondary-foreground hover:bg-secondary/90 h-8 w-8 p-0"
                                  onClick={() => handleDelete(perm)}
                                >
                                  <Trash2 className="h-4 w-4" />
                                </Button>
                              </div>
                            </TableCell>
                          </TableRow>
                        ))}
                      </Fragment>
                    )
                  })}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* 添加 / 编辑权限弹窗 */}
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
                <Input
                  value={formData.code}
                  onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                  placeholder="如：user:create"
                  required
                  disabled={!!editingId}
                />
                <p className="text-xs text-muted-foreground">格式：模块:操作，如 user:create</p>
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
                  onChange={(e) => setFormData({ ...formData, isActive: e.target.checked })}
                  className="rounded border-gray-300"
                />
                <label className="text-sm font-medium">启用权限</label>
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                  onClick={() => setShowDialog(false)}
                >
                  取消
                </Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {createMutation.isPending || updateMutation.isPending ? '保存中...' : editingId ? '更新' : '创建'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
      <ConfirmDialog />
    </>
  )
}
