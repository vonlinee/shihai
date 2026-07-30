import { useMemo, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { useCallback } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Plus, Edit2, Trash2, X } from 'lucide-react'
import { useDynasties } from '@/hooks/usePoems'
import {
  
  useAdminCreateDynasty, useAdminUpdateDynasty, useAdminDeleteDynasty,
  useAdminBatchDeleteDynasties
    ,
} from '@/hooks/useAdmin'
import { toEntityId, toggleAllVisibleIds, toggleSelectedId } from '@/utils/uitools'

export function DynastiesTab() {
  const { data: dynasties, isLoading } = useDynasties()
  const createMutation = useAdminCreateDynasty()
  const updateMutation = useAdminUpdateDynasty()
  const deleteMutation = useAdminDeleteDynasty()
  const batchDeleteMutation = useAdminBatchDeleteDynasties()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [formData, setFormData] = useState({ name: '', period: '', description: '' })
  const [selectedDynastyIds, setSelectedDynastyIds] = useState<string[]>([])
  const dynastiesList = dynasties ?? []
  const visibleDynastyIds = dynastiesList.map((d) => toEntityId(d.id))
  const allVisibleDynastiesSelected = visibleDynastyIds.length > 0 && visibleDynastyIds.every((id) => selectedDynastyIds.includes(id))

  const openCreate = () => { setEditingId(null); setFormData({ name: '', period: '', description: '' }); setShowDialog(true) }
  const openEdit = useCallback((d: { id: string | number; name: string; period?: string; description?: string }) => {
    setEditingId(toEntityId(d.id)); setFormData({ name: d.name, period: d.period || '', description: d.description || '' }); setShowDialog(true)
  }, [])
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingId) updateMutation.mutate({ id: editingId, data: formData }, { onSuccess: () => setShowDialog(false) })
    else createMutation.mutate(formData, { onSuccess: () => setShowDialog(false) })
  }
  const handleDelete = useCallback(async (id: string) => {
    const confirmed = await confirm({
      title: '删除朝代',
      description: '确定要删除该朝代吗？此操作不可撤销。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteMutation.mutate(id, {
      onSuccess: () => setSelectedDynastyIds((ids) => ids.filter((selectedId) => selectedId !== id)),
    })
  }, [confirm, deleteMutation])

  const handleBatchDelete = async () => {
    if (selectedDynastyIds.length === 0) return
    const confirmed = await confirm({
      title: '批量删除朝代',
      description: `确定要删除选中的 ${selectedDynastyIds.length} 个朝代吗？此操作不可撤销。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeleteMutation.mutate(selectedDynastyIds, { onSuccess: () => setSelectedDynastyIds([]) })
  }

  type DynastyRow = (typeof dynastiesList)[number]

  const dynastyColumns = useMemo<DataTableColumn<DynastyRow>[]>(
    () => [
      {
        id: 'select',
        header: () => (
          <Checkbox
            aria-label="选择全部朝代"
            checked={allVisibleDynastiesSelected}
            onCheckedChange={() => setSelectedDynastyIds((ids) => toggleAllVisibleIds(ids, visibleDynastyIds))}
          />
        ),
        cell: ({ row }) => (
          <Checkbox
            aria-label={`选择朝代 ${row.original.name}`}
            checked={selectedDynastyIds.includes(toEntityId(row.original.id))}
            onCheckedChange={() => setSelectedDynastyIds((ids) => toggleSelectedId(ids, toEntityId(row.original.id)))}
          />
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { draggable: false, fixed: 'left', width: 56 },
      },
      {
        accessorKey: 'name',
        header: '朝代名称',
        cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
        meta: { filterPlaceholder: '筛选朝代', fixed: 'left', width: 180 },
      },
      {
        accessorFn: (dynasty) => dynasty.period ?? '',
        id: 'period',
        header: '时期',
        cell: ({ row }) => <span className="text-muted-foreground">{row.original.period || '-'}</span>,
        meta: { filterPlaceholder: '筛选时期', width: 160 },
      },
      {
        accessorFn: (dynasty) => dynasty.description ?? '',
        id: 'description',
        header: '描述',
        cell: ({ row }) => <span className="block max-w-xs truncate text-muted-foreground">{row.original.description || '-'}</span>,
        meta: { filterPlaceholder: '筛选描述', minWidth: 260 },
      },
      {
        id: 'actions',
        header: '操作',
        cell: ({ row }) => (
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="sm" onClick={() => openEdit(row.original)}><Edit2 className="h-4 w-4" /></Button>
            <Button variant="ghost" size="sm" onClick={() => handleDelete(toEntityId(row.original.id))}><Trash2 className="h-4 w-4 text-cinnabar" /></Button>
          </div>
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { draggable: false, fixed: 'right', width: 120 },
      },
    ],
    [allVisibleDynastiesSelected, handleDelete, openEdit, selectedDynastyIds, visibleDynastyIds],
  )

  return (
    <>
      <div className="flex items-center justify-end gap-2">
        {selectedDynastyIds.length > 0 && (
          <Button
            variant="destructive"
            onClick={handleBatchDelete}
            disabled={batchDeleteMutation.isPending}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            删除选中({selectedDynastyIds.length})
          </Button>
        )}
        <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />添加朝代</Button>
      </div>
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? <div className="text-center py-8 text-muted-foreground">加载中...</div>
            : dynastiesList.length === 0 ? <div className="text-center py-8 text-muted-foreground">暂无朝代数据</div>
            : (
              <DataTable
                columns={dynastyColumns}
                data={dynastiesList}
                emptyText="暂无朝代数据"
                enableColumnDragging
                enableColumnFilters
                getRowId={(dynasty) => toEntityId(dynasty.id)}
              />
            )}
        </CardContent>
      </Card>

      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-md max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingId ? '编辑朝代' : '添加朝代'}</h2>
              <button onClick={() => setShowDialog(false)} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="space-y-2">
                <label className="text-sm font-medium">朝代名称 *</label>
                <Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} placeholder="如：唐" required />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">时期</label>
                <Input value={formData.period} onChange={(e) => setFormData({ ...formData, period: e.target.value })} placeholder="如：618-907" />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">描述</label>
                <textarea value={formData.description} onChange={(e) => setFormData({ ...formData, description: e.target.value })} placeholder="朝代简介" rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm" />
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90" onClick={() => setShowDialog(false)}>取消</Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {(createMutation.isPending || updateMutation.isPending) ? '提交中...' : editingId ? '保存' : '添加'}
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
