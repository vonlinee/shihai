import { useMemo, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { useCallback } from 'react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Combobox, type ComboboxOption } from '@/components/ui/combobox'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Plus, Edit2, Trash2, X } from 'lucide-react'
import {
  
  useAdminPoemTypes, useAdminCreatePoemType, useAdminUpdatePoemType,
  useAdminDeletePoemType, useAdminBatchDeletePoemTypes
  ,
} from '@/hooks/useAdmin'
import type { PoemType } from '@/types'
import { toEntityId, toggleAllVisibleIds, toggleSelectedId } from '@/utils/uitools'


function parseOptionalPositiveIntInput(value: string): number | null {
  const trimmed = value.trim()
  if (!trimmed) return null
  const parsed = Number(trimmed)
  if (!Number.isInteger(parsed) || parsed <= 0) return Number.NaN
  return parsed
}


export function PoemTypesTab() {
  const { data: poemTypes, isLoading } = useAdminPoemTypes()
  const createMutation = useAdminCreatePoemType()
  const updateMutation = useAdminUpdatePoemType()
  const deleteMutation = useAdminDeletePoemType()
  const batchDeleteMutation = useAdminBatchDeletePoemTypes()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    category: '',
    lines: '',
    charsPerLine: '',
    description: '',
  })
  const [selectedPoemTypeIds, setSelectedPoemTypeIds] = useState<string[]>([])

  const poemTypeList = poemTypes ?? []
  const visiblePoemTypeIds = poemTypeList.map((poemType) => toEntityId(poemType.id))
  const allVisiblePoemTypesSelected = visiblePoemTypeIds.length > 0 && visiblePoemTypeIds.every((id) => selectedPoemTypeIds.includes(id))
  const categoryOptions: ComboboxOption[] = useMemo(() => {
    const categories = new Set(['诗', '词', '曲', '文', '其他'])
    poemTypeList.forEach((poemType) => {
      if (poemType.category) categories.add(poemType.category)
    })
    return Array.from(categories).map((category) => ({ value: category, label: category }))
  }, [poemTypeList])

  const openCreate = () => {
    setEditingId(null)
    setFormData({ name: '', category: '', lines: '', charsPerLine: '', description: '' })
    setShowDialog(true)
  }

  const openEdit = useCallback((poemType: PoemType) => {
    setEditingId(toEntityId(poemType.id))
    setFormData({
      name: poemType.name,
      category: poemType.category,
      lines: poemType.lines ? String(poemType.lines) : '',
      charsPerLine: poemType.charsPerLine ? String(poemType.charsPerLine) : '',
      description: poemType.description || '',
    })
    setShowDialog(true)
  }, [])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const lines = parseOptionalPositiveIntInput(formData.lines)
    if (Number.isNaN(lines)) {
      toast.error('句数必须是大于 0 的整数')
      return
    }
    const charsPerLine = parseOptionalPositiveIntInput(formData.charsPerLine)
    if (Number.isNaN(charsPerLine)) {
      toast.error('每句字数必须是大于 0 的整数')
      return
    }
    const payload = {
      name: formData.name.trim(),
      category: formData.category.trim(),
      lines,
      charsPerLine,
      description: formData.description,
    }
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: payload }, { onSuccess: () => setShowDialog(false) })
      return
    }
    createMutation.mutate(payload, { onSuccess: () => setShowDialog(false) })
  }

  const handleDelete = useCallback(async (id: string) => {
    const confirmed = await confirm({
      title: '删除体裁',
      description: '确定要删除该体裁吗？已有诗词中的体裁文本不会被自动清空。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteMutation.mutate(id, {
      onSuccess: () => setSelectedPoemTypeIds((ids) => ids.filter((selectedId) => selectedId !== id)),
    })
  }, [confirm, deleteMutation])

  const handleBatchDelete = async () => {
    if (selectedPoemTypeIds.length === 0) return
    const confirmed = await confirm({
      title: '批量删除体裁',
      description: `确定要删除选中的 ${selectedPoemTypeIds.length} 个体裁吗？已有诗词中的体裁文本不会被自动清空。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeleteMutation.mutate(selectedPoemTypeIds, { onSuccess: () => setSelectedPoemTypeIds([]) })
  }

  const poemTypeColumns = useMemo<DataTableColumn<PoemType>[]>(
    () => [
      {
        id: 'select',
        header: () => (
          <Checkbox
            aria-label="选择全部体裁"
            checked={allVisiblePoemTypesSelected}
            onCheckedChange={() => setSelectedPoemTypeIds((ids) => toggleAllVisibleIds(ids, visiblePoemTypeIds))}
          />
        ),
        cell: ({ row }) => (
          <Checkbox
            aria-label={`选择体裁 ${row.original.name}`}
            checked={selectedPoemTypeIds.includes(toEntityId(row.original.id))}
            onCheckedChange={() => setSelectedPoemTypeIds((ids) => toggleSelectedId(ids, toEntityId(row.original.id)))}
          />
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { draggable: false, fixed: 'left', width: 56 },
      },
      {
        accessorKey: 'category',
        header: '一级分类',
        cell: ({ row }) => <span className="font-medium">{row.original.category}</span>,
        meta: { filterPlaceholder: '筛选大类', fixed: 'left', width: 140 },
      },
      {
        accessorKey: 'name',
        header: '细分类别',
        cell: ({ row }) => <span>{row.original.name}</span>,
        meta: { filterPlaceholder: '筛选体裁', width: 180 },
      },
      {
        accessorFn: (poemType) => poemType.lines ?? '',
        id: 'lines',
        header: '句数',
        cell: ({ row }) => <span className="text-muted-foreground">{row.original.lines ?? '不限'}</span>,
        meta: { filterPlaceholder: '筛选句数', width: 120 },
      },
      {
        accessorFn: (poemType) => poemType.charsPerLine ?? '',
        id: 'charsPerLine',
        header: '每句字数',
        cell: ({ row }) => <span className="text-muted-foreground">{row.original.charsPerLine ?? '不限'}</span>,
        meta: { filterPlaceholder: '筛选字数', width: 130 },
      },
      {
        accessorFn: (poemType) => poemType.description ?? '',
        id: 'description',
        header: '说明',
        cell: ({ row }) => <span className="block max-w-md truncate text-muted-foreground">{row.original.description || '-'}</span>,
        meta: { filterPlaceholder: '筛选说明', minWidth: 280 },
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
    [allVisiblePoemTypesSelected, handleDelete, openEdit, selectedPoemTypeIds, visiblePoemTypeIds],
  )

  return (
    <>
      <div className="flex items-center justify-end gap-2">
        {selectedPoemTypeIds.length > 0 && (
          <Button
            variant="destructive"
            onClick={handleBatchDelete}
            disabled={batchDeleteMutation.isPending}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            删除选中({selectedPoemTypeIds.length})
          </Button>
        )}
        <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />添加体裁</Button>
      </div>
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? <div className="text-center py-8 text-muted-foreground">加载中...</div>
            : poemTypeList.length === 0 ? <div className="text-center py-8 text-muted-foreground">暂无体裁数据</div>
            : (
              <DataTable
                columns={poemTypeColumns}
                data={poemTypeList}
                emptyText="暂无体裁数据"
                enableColumnDragging
                enableColumnFilters
                getRowId={(poemType) => toEntityId(poemType.id)}
              />
            )}
        </CardContent>
      </Card>

      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-lg max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingId ? '编辑体裁' : '添加体裁'}</h2>
              <button onClick={() => setShowDialog(false)} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">一级分类 *</label>
                  <Combobox
                    options={categoryOptions}
                    value={formData.category}
                    onChange={(val) => setFormData({ ...formData, category: String(val || '') })}
                    placeholder="选择或输入大类"
                    allowCustom
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">细分类别 *</label>
                  <Input value={formData.name} onChange={(e) => setFormData({ ...formData, name: e.target.value })} placeholder="如：五言绝句" required />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">句数</label>
                  <Input type="number" min={1} value={formData.lines} onChange={(e) => setFormData({ ...formData, lines: e.target.value })} placeholder="不限" />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">每句字数</label>
                  <Input type="number" min={1} value={formData.charsPerLine} onChange={(e) => setFormData({ ...formData, charsPerLine: e.target.value })} placeholder="不限" />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">说明</label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="体裁说明"
                  rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                />
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
