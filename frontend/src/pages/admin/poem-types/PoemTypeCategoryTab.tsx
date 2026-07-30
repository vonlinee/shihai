import { useCallback, useMemo, useState } from 'react'
import { Edit2, Plus, Trash2 } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
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

import { PoemTypeFormDialog, type PoemTypeFormPayload } from './PoemTypeFormDialog'

interface PoemTypeCategoryTabProps {
  category: string
}

export function PoemTypeCategoryTab({ category }: PoemTypeCategoryTabProps) {
  const { data: poemTypes, isLoading } = useAdminPoemTypes()
  const createMutation = useAdminCreatePoemType()
  const updateMutation = useAdminUpdatePoemType()
  const deleteMutation = useAdminDeletePoemType()
  const batchDeleteMutation = useAdminBatchDeletePoemTypes()
  const { confirm, ConfirmDialog } = useConfirmDialog()
  const [showDialog, setShowDialog] = useState(false)
  const [editingPoemType, setEditingPoemType] = useState<PoemType | null>(null)
  const [selectedPoemTypeIds, setSelectedPoemTypeIds] = useState<string[]>([])

  const poemTypeList = useMemo(
    () => (poemTypes ?? []).filter((poemType) => poemType.category === category),
    [category, poemTypes],
  )
  const visiblePoemTypeIds = poemTypeList.map((poemType) => toEntityId(poemType.id))
  const allVisiblePoemTypesSelected = visiblePoemTypeIds.length > 0 && visiblePoemTypeIds.every((id) => selectedPoemTypeIds.includes(id))

  const openCreate = () => {
    setEditingPoemType(null)
    setShowDialog(true)
  }

  const openEdit = useCallback((poemType: PoemType) => {
    setEditingPoemType(poemType)
    setShowDialog(true)
  }, [])

  const closeDialog = () => {
    setShowDialog(false)
    setEditingPoemType(null)
  }

  const handleSubmit = (payload: PoemTypeFormPayload) => {
    if (editingPoemType) {
      updateMutation.mutate(
        { id: toEntityId(editingPoemType.id), data: payload },
        { onSuccess: closeDialog },
      )
      return
    }
    createMutation.mutate(payload, { onSuccess: closeDialog })
  }

  const handleDelete = useCallback(
    async (id: string) => {
      const confirmed = await confirm({
        title: `删除${category}体裁`,
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
    [category, confirm, deleteMutation],
  )

  const handleBatchDelete = async () => {
    if (selectedPoemTypeIds.length === 0) return
    const confirmed = await confirm({
      title: `批量删除${category}体裁`,
      description: `确定要删除选中的 ${selectedPoemTypeIds.length} 个体裁吗？已有诗词中的体裁文本不会被自动清空。`,
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    batchDeleteMutation.mutate(selectedPoemTypeIds, {
      onSuccess: () => setSelectedPoemTypeIds([]),
    })
  }

  const poemTypeColumns = useMemo<DataTableColumn<PoemType>[]>(
    () => [
      {
        id: 'select',
        header: () => (
          <Checkbox
            aria-label={`选择全部${category}体裁`}
            checked={allVisiblePoemTypesSelected}
            onCheckedChange={() =>
              setSelectedPoemTypeIds((ids) => toggleAllVisibleIds(ids, visiblePoemTypeIds))
            }
          />
        ),
        cell: ({ row }) => {
          const poemTypeId = toEntityId(row.original.id)
          return (
            <Checkbox
              aria-label={`选择体裁 ${row.original.name}`}
              checked={selectedPoemTypeIds.includes(poemTypeId)}
              onCheckedChange={() =>
                setSelectedPoemTypeIds((ids) => toggleSelectedId(ids, poemTypeId))
              }
            />
          )
        },
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'center', draggable: false, fixed: 'left', headerAlign: 'center', minWidth: 40, resizable: false, width: 40 },
      },
      {
        accessorKey: 'name',
        header: '细分类别',
        cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
        meta: { filterPlaceholder: '筛选体裁', fixed: 'left', minWidth: 220 },
      },
      {
        accessorFn: (poemType) => poemType.lines ?? '',
        id: 'lines',
        header: '句数',
        cell: ({ row }) => (
          <span className="text-muted-foreground">{row.original.lines ?? '不限'}</span>
        ),
        meta: { filterPlaceholder: '筛选句数', width: 130 },
      },
      {
        accessorFn: (poemType) => poemType.charsPerLine ?? '',
        id: 'charsPerLine',
        header: '每句字数',
        cell: ({ row }) => (
          <span className="text-muted-foreground">{row.original.charsPerLine ?? '不限'}</span>
        ),
        meta: { filterPlaceholder: '筛选字数', width: 130 },
      },
      {
        accessorFn: (poemType) => poemType.description ?? '',
        id: 'description',
        header: '说明',
        cell: ({ row }) => (
          <span className="block max-w-md truncate text-muted-foreground">
            {row.original.description || '-'}
          </span>
        ),
        meta: { filterPlaceholder: '筛选说明', minWidth: 280 },
      },
      {
        id: 'actions',
        header: '操作',
        cell: ({ row }) => (
          <div className="flex items-center justify-end gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 rounded-sm"
              onClick={() => openEdit(row.original)}
            >
              <Edit2 className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 rounded-sm"
              onClick={() => handleDelete(toEntityId(row.original.id))}
            >
              <Trash2 className="h-4 w-4 text-cinnabar" />
            </Button>
          </div>
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { align: 'right', draggable: false, fixed: 'right', headerAlign: 'center', width: 120 },
      },
    ],
    [
      allVisiblePoemTypesSelected,
      category,
      handleDelete,
      openEdit,
      selectedPoemTypeIds,
      visiblePoemTypeIds,
    ],
  )

  return (
    <>
      <div className="mb-4 flex items-center justify-end gap-2">
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
          添加{category}体裁
        </Button>
      </div>

      <Card className="ink-border">
        <CardContent className="pt-6">
          <DataTable
            columns={poemTypeColumns}
            data={poemTypeList}
            emptyText={`暂无${category}体裁数据`}
            enableColumnDragging
            enableColumnFilters
            getRowId={(poemType) => toEntityId(poemType.id)}
            loading={isLoading}
          />
        </CardContent>
      </Card>

      {showDialog && (
        <PoemTypeFormDialog
          category={category}
          editingPoemType={editingPoemType}
          isSubmitting={createMutation.isPending || updateMutation.isPending}
          onClose={closeDialog}
          onSubmit={handleSubmit}
        />
      )}
      <ConfirmDialog />
    </>
  )
}
