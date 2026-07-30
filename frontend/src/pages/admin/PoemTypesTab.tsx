import { useCallback, useMemo, useState } from 'react'
import { toast } from 'sonner'
import { ChevronDown, ChevronRight, Edit2, FolderTree, Plus, Trash2, X } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Combobox, type ComboboxOption } from '@/components/ui/combobox'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
import { Input } from '@/components/ui/input'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import {
  useAdminBatchDeletePoemTypes,
  useAdminCreatePoemType,
  useAdminDeletePoemType,
  useAdminPoemTypes,
  useAdminUpdatePoemType,
} from '@/hooks/useAdmin'
import type { PoemType } from '@/types'
import { toEntityId, toggleAllVisibleIds, toggleSelectedId } from '@/utils/uitools'

interface PoemTypeGroup {
  category: string
  poemTypes: PoemType[]
}

type PoemTypeTreeRow =
  | {
    kind: 'category'
    id: string
    category: string
    poemTypes: PoemType[]
  }
  | {
    kind: 'poemType'
    id: string
    category: string
    poemType: PoemType
  }

function parseOptionalPositiveIntInput(value: string): number | null {
  const trimmed = value.trim()
  if (!trimmed) return null
  const parsed = Number(trimmed)
  if (!Number.isInteger(parsed) || parsed <= 0) return Number.NaN
  return parsed
}

function groupPoemTypes(poemTypes: PoemType[]): PoemTypeGroup[] {
  const groups = new Map<string, PoemType[]>()
  for (const poemType of poemTypes) {
    const category = poemType.category || '未分类'
    const items = groups.get(category) ?? []
    items.push(poemType)
    groups.set(category, items)
  }
  return Array.from(groups.entries()).map(([category, items]) => ({
    category,
    poemTypes: items,
  }))
}

function flattenPoemTypeGroups(groups: PoemTypeGroup[], collapsedCategories: Set<string>): PoemTypeTreeRow[] {
  return groups.flatMap((group) => {
    const categoryRow: PoemTypeTreeRow = {
      kind: 'category',
      id: `category:${group.category}`,
      category: group.category,
      poemTypes: group.poemTypes,
    }

    if (collapsedCategories.has(group.category)) {
      return [categoryRow]
    }

    return [
      categoryRow,
      ...group.poemTypes.map((poemType) => ({
        kind: 'poemType' as const,
        id: toEntityId(poemType.id),
        category: group.category,
        poemType,
      })),
    ]
  })
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
  const [collapsedCategories, setCollapsedCategories] = useState<Set<string>>(new Set())

  const poemTypeList = poemTypes ?? []
  const poemTypeGroups = useMemo(() => groupPoemTypes(poemTypeList), [poemTypeList])
  const poemTypeTreeRows = useMemo(
    () => flattenPoemTypeGroups(poemTypeGroups, collapsedCategories),
    [collapsedCategories, poemTypeGroups],
  )
  const visiblePoemTypeIds = poemTypeList.map((poemType) => toEntityId(poemType.id))
  const allVisiblePoemTypesSelected = visiblePoemTypeIds.length > 0 && visiblePoemTypeIds.every((id) => selectedPoemTypeIds.includes(id))
  const allCategoriesCollapsed = poemTypeGroups.length > 0 && poemTypeGroups.every((group) => collapsedCategories.has(group.category))

  const categoryOptions: ComboboxOption[] = useMemo(() => {
    const categories = new Set(['诗', '词', '曲', '文', '其他'])
    poemTypeList.forEach((poemType) => {
      if (poemType.category) categories.add(poemType.category)
    })
    return Array.from(categories).map((category) => ({
      value: category,
      label: category,
    }))
  }, [poemTypeList])

  const toggleCategory = (category: string) => {
    setCollapsedCategories((current) => {
      const next = new Set(current)
      if (next.has(category)) next.delete(category)
      else next.add(category)
      return next
    })
  }

  const toggleAllCategories = () => {
    if (allCategoriesCollapsed) {
      setCollapsedCategories(new Set())
      return
    }
    setCollapsedCategories(new Set(poemTypeGroups.map((group) => group.category)))
  }

  const toggleCategorySelection = (ids: string[]) => {
    setSelectedPoemTypeIds((current) => toggleAllVisibleIds(current, ids))
  }

  const openCreate = () => {
    setEditingId(null)
    setFormData({
      name: '',
      category: '',
      lines: '',
      charsPerLine: '',
      description: '',
    })
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
      updateMutation.mutate(
        { id: editingId, data: payload },
        { onSuccess: () => setShowDialog(false) },
      )
      return
    }
    createMutation.mutate(payload, { onSuccess: () => setShowDialog(false) })
  }

  const handleDelete = useCallback(
    async (id: string) => {
      const confirmed = await confirm({
        title: '删除体裁',
        description: '确定要删除该体裁吗？已有诗词中的体裁文本不会被自动清空。',
        confirmText: '删除',
        destructive: true,
      })
      if (!confirmed) return
      deleteMutation.mutate(id, {
        onSuccess: () =>
          setSelectedPoemTypeIds((ids) => ids.filter((selectedId) => selectedId !== id)),
      })
    },
    [confirm, deleteMutation],
  )

  const handleBatchDelete = async () => {
    if (selectedPoemTypeIds.length === 0) return
    const confirmed = await confirm({
      title: '批量删除体裁',
      description: `确定要删除选中的 ${selectedPoemTypeIds.length} 个体裁吗？已有诗词中的体裁文本不会被自动清空。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeleteMutation.mutate(selectedPoemTypeIds, {
      onSuccess: () => setSelectedPoemTypeIds([]),
    })
  }

  const poemTypeColumns = useMemo<DataTableColumn<PoemTypeTreeRow>[]>(
    () => [
      {
        id: 'select',
        header: () => (
          <Checkbox
            aria-label="选择全部体裁"
            checked={allVisiblePoemTypesSelected}
            onCheckedChange={() =>
              setSelectedPoemTypeIds((ids) => toggleAllVisibleIds(ids, visiblePoemTypeIds))
            }
          />
        ),
        cell: ({ row }) => {
          if (row.original.kind === 'category') {
            const categoryIds = row.original.poemTypes.map((poemType) => toEntityId(poemType.id))
            const categorySelected = categoryIds.length > 0 && categoryIds.every((id) => selectedPoemTypeIds.includes(id))
            return (
              <Checkbox
                aria-label={`选择 ${row.original.category} 下全部体裁`}
                checked={categorySelected}
                onCheckedChange={() => toggleCategorySelection(categoryIds)}
              />
            )
          }

          return (
            <Checkbox
              aria-label={`选择体裁 ${row.original.poemType.name}`}
              checked={selectedPoemTypeIds.includes(row.original.id)}
              onCheckedChange={() =>
                setSelectedPoemTypeIds((ids) => toggleSelectedId(ids, row.original.id))
              }
            />
          )
        },
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'center', draggable: false, fixed: 'left', headerAlign: 'center', minWidth: 40, resizable: false, width: 40 },
      },
      {
        accessorFn: (row) => row.kind === 'category' ? row.category : row.poemType.name,
        id: 'name',
        header: '分类 / 体裁',
        cell: ({ row }) => {
          if (row.original.kind === 'category') {
            const collapsed = collapsedCategories.has(row.original.category)
            return (
              <button
                type="button"
                className="inline-flex min-w-0 items-center gap-2 font-medium"
                onClick={() => toggleCategory(row.original.category)}
              >
                {collapsed ? (
                  <ChevronRight className="h-4 w-4 shrink-0 text-muted-foreground" />
                ) : (
                  <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground" />
                )}
                <span className="truncate">{row.original.category}</span>
              </button>
            )
          }

          return <span className="pl-8 font-medium">{row.original.poemType.name}</span>
        },
        enableSorting: false,
        meta: { fixed: 'left', minWidth: 220 },
      },
      {
        accessorFn: (row) => row.kind === 'poemType' ? row.poemType.lines ?? '' : '',
        id: 'lines',
        header: '句数',
        cell: ({ row }) => row.original.kind === 'poemType'
          ? <span className="text-muted-foreground">{row.original.poemType.lines ?? '不限'}</span>
          : <span className="text-muted-foreground">共 {row.original.poemTypes.length} 个细分类别</span>,
        enableSorting: false,
        meta: { width: 150 },
      },
      {
        accessorFn: (row) => row.kind === 'poemType' ? row.poemType.charsPerLine ?? '' : '',
        id: 'charsPerLine',
        header: '每句字数',
        cell: ({ row }) => row.original.kind === 'poemType'
          ? <span className="text-muted-foreground">{row.original.poemType.charsPerLine ?? '不限'}</span>
          : null,
        enableSorting: false,
        meta: { width: 130 },
      },
      {
        accessorFn: (row) => row.kind === 'poemType' ? row.poemType.description ?? '' : '',
        id: 'description',
        header: '说明',
        cell: ({ row }) => row.original.kind === 'poemType'
          ? (
            <span className="block max-w-md truncate text-muted-foreground">
              {row.original.poemType.description || '-'}
            </span>
          )
          : null,
        enableSorting: false,
        meta: { minWidth: 280 },
      },
      {
        id: 'actions',
        header: '操作',
        cell: ({ row }) => {
          if (row.original.kind !== 'poemType') return null
          const poemTypeRow = row.original

          return (
            <div className="flex items-center justify-end gap-1">
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 rounded-sm"
                onClick={() => openEdit(poemTypeRow.poemType)}
              >
                <Edit2 className="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 rounded-sm"
                onClick={() => handleDelete(poemTypeRow.id)}
              >
                <Trash2 className="h-4 w-4 text-cinnabar" />
              </Button>
            </div>
          )
        },
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'right', draggable: false, fixed: 'right', headerAlign: 'center', width: 120 },
      },
    ],
    [
      allVisiblePoemTypesSelected,
      collapsedCategories,
      handleDelete,
      openEdit,
      selectedPoemTypeIds,
      visiblePoemTypeIds,
    ],
  )

  return (
    <>
      <div className="flex items-center justify-end gap-2">
        <Button
          className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
          onClick={toggleAllCategories}
          disabled={poemTypeGroups.length === 0}
        >
          <FolderTree className="h-4 w-4 mr-2" />
          {allCategoriesCollapsed ? '全部展开' : '全部收起'}
        </Button>
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
        <Button onClick={openCreate}>
          <Plus className="h-4 w-4 mr-2" />
          添加体裁
        </Button>
      </div>

      <Card className="ink-border">
        <CardContent className="pt-6">
          <DataTable
            columns={poemTypeColumns}
            data={poemTypeTreeRows}
            emptyText="暂无体裁数据"
            enableColumnDragging
            enableSorting={false}
            getRowId={(row) => row.id}
            loading={isLoading}
          />
        </CardContent>
      </Card>

      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-lg max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">
                {editingId ? '编辑体裁' : '添加体裁'}
              </h2>
              <button
                onClick={() => setShowDialog(false)}
                className="p-1 hover:bg-muted rounded"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            <form
              onSubmit={handleSubmit}
              className="space-y-4 p-6 pt-4 overflow-y-auto"
            >
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">一级分类 *</label>
                  <Combobox
                    options={categoryOptions}
                    value={formData.category}
                    onChange={(val) =>
                      setFormData({ ...formData, category: String(val || '') })
                    }
                    placeholder="选择或输入大类"
                    allowCustom
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">细分类别 *</label>
                  <Input
                    value={formData.name}
                    onChange={(e) =>
                      setFormData({ ...formData, name: e.target.value })
                    }
                    placeholder="如：五言绝句"
                    required
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">句数</label>
                  <Input
                    type="number"
                    min={1}
                    value={formData.lines}
                    onChange={(e) =>
                      setFormData({ ...formData, lines: e.target.value })
                    }
                    placeholder="不限"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">每句字数</label>
                  <Input
                    type="number"
                    min={1}
                    value={formData.charsPerLine}
                    onChange={(e) =>
                      setFormData({ ...formData, charsPerLine: e.target.value })
                    }
                    placeholder="不限"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">说明</label>
                <textarea
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                  placeholder="体裁说明"
                  rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                />
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
                  onClick={() => setShowDialog(false)}
                >
                  取消
                </Button>
                <Button
                  type="submit"
                  disabled={createMutation.isPending || updateMutation.isPending}
                >
                  {createMutation.isPending || updateMutation.isPending
                    ? '提交中...'
                    : editingId
                      ? '保存'
                      : '添加'}
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
