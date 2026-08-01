import { useEffect, useMemo, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Combobox, type ComboboxOption } from '@/components/ui/combobox'
import { DataTable, type DataTableColumn } from '@/components/ui/data-table'
import { Pagination } from '@/components/ui/Pagination'
import { useConfirmDialog } from '@/components/ui/ConfirmDialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Edit2, Plus, Search, Trash2, X } from 'lucide-react'
import {
  useAdminAddWorkCollectionItem,
  useAdminCreateWorkCollection,
  useAdminDeleteWorkCollection,
  useAdminDeleteWorkCollectionItem,
  useAdminUpdateWorkCollection,
  useAdminUpdateWorkCollectionItem,
  useAdminWorkCollection,
  useAdminWorkCollections,
} from '@/hooks/useAdmin'
import { usePoems } from '@/hooks/usePoems'
import type { WorkCollection } from '@/types'

interface CollectionFormState {
  title: string
  description: string
  coverImage: string
  isPublished: boolean
}

export function AdminWorkCollectionsPage() {
  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [showDialog, setShowDialog] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [formData, setFormData] = useState<CollectionFormState>({
    title: '',
    description: '',
    coverImage: '',
    isPublished: false,
  })
  const [selectedPoemId, setSelectedPoemId] = useState<number | undefined>()
  const [sortOrder, setSortOrder] = useState('0')

  const { data, isLoading } = useAdminWorkCollections({ page, pageSize, keyword: keyword || undefined })
  const { data: detail } = useAdminWorkCollection(editingId)
  const { data: poemData } = usePoems({ page: 1, pageSize: 50 })
  const createMutation = useAdminCreateWorkCollection()
  const updateMutation = useAdminUpdateWorkCollection()
  const deleteMutation = useAdminDeleteWorkCollection()
  const addItemMutation = useAdminAddWorkCollectionItem()
  const updateItemMutation = useAdminUpdateWorkCollectionItem()
  const deleteItemMutation = useAdminDeleteWorkCollectionItem()
  const { confirm, ConfirmDialog } = useConfirmDialog()

  const collections = data?.list ?? []
  const total = data?.total ?? 0
  const poemOptions: ComboboxOption[] = useMemo(
    () => (poemData?.list ?? []).map((poem) => ({ value: poem.id, label: poem.title, description: poem.author?.name })),
    [poemData],
  )
  const poemTitleById = useMemo(() => {
    const result = new Map<number, string>()
    for (const poem of poemData?.list ?? []) {
      result.set(poem.id, poem.title)
    }
    return result
  }, [poemData])

  const collectionColumns = useMemo<DataTableColumn<WorkCollection>[]>(
    () => [
      {
        accessorKey: 'title',
        header: '标题',
        cell: ({ row }) => <span className="font-medium">{row.original.title}</span>,
      },
      {
        accessorKey: 'isPublished',
        header: '状态',
        cell: ({ row }) => (row.original.isPublished ? '已发布' : '草稿'),
      },
      {
        accessorKey: 'itemCount',
        header: '作品数',
        enableColumnFilter: false,
      },
      {
        accessorKey: 'createdAt',
        header: '创建时间',
        enableColumnFilter: false,
        cell: ({ row }) => (
          <span className="text-muted-foreground">
            {new Date(row.original.createdAt).toLocaleDateString()}
          </span>
        ),
      },
      {
        id: 'actions',
        header: '操作',
        enableColumnFilter: false,
        enableSorting: false,
        cell: ({ row }) => (
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="sm" onClick={() => openEdit(row.original)}>
              <Edit2 className="h-4 w-4" />
            </Button>
            <Button variant="ghost" size="sm" onClick={() => handleDelete(row.original.id)}>
              <Trash2 className="h-4 w-4 text-cinnabar" />
            </Button>
          </div>
        ),
      },
    ],
    [],
  )

  useEffect(() => {
    if (!detail || !editingId) return
    setFormData({
      title: detail.title,
      description: detail.description || '',
      coverImage: detail.coverImage || '',
      isPublished: detail.isPublished,
    })
  }, [detail, editingId])

  const openCreate = () => {
    setEditingId(null)
    setFormData({ title: '', description: '', coverImage: '', isPublished: false })
    setSelectedPoemId(undefined)
    setSortOrder('0')
    setShowDialog(true)
  }

  const openEdit = (collection: WorkCollection) => {
    setEditingId(collection.id)
    setFormData({
      title: collection.title,
      description: collection.description || '',
      coverImage: collection.coverImage || '',
      isPublished: collection.isPublished,
    })
    setSelectedPoemId(undefined)
    setSortOrder('0')
    setShowDialog(true)
  }

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    const payload = {
      title: formData.title,
      description: formData.description,
      coverImage: formData.coverImage,
      isPublished: formData.isPublished,
    }
    if (editingId) {
      updateMutation.mutate({ id: editingId, data: payload }, { onSuccess: () => setShowDialog(false) })
    } else {
      createMutation.mutate(payload, { onSuccess: () => setShowDialog(false) })
    }
  }

  const handleDelete = async (id: number) => {
    const confirmed = await confirm({
      title: '删除作品集',
      description: '确定要删除该作品集吗？此操作不可撤销。',
      confirmText: '删除',
      destructive: true,
    })
    if (!confirmed) return
    deleteMutation.mutate(id)
  }

  const handleAddItem = () => {
    if (!editingId || !selectedPoemId) return
    addItemMutation.mutate({
      collectionId: editingId,
      data: { workType: 'poem', workId: selectedPoemId, sortOrder: Number(sortOrder) || 0 },
    }, {
      onSuccess: () => {
        setSelectedPoemId(undefined)
        setSortOrder('0')
      },
    })
  }

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between gap-4">
        <div className="relative max-w-sm flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="搜索作品集标题..."
            className="pl-10"
            value={keyword}
            onChange={(event) => {
              setKeyword(event.target.value)
              setPage(1)
            }}
          />
        </div>
        <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />新增作品集</Button>
      </div>

      <Card className="ink-border">
        <CardContent className="pt-6">
          <DataTable
            columns={collectionColumns}
            data={collections}
            emptyText="暂无作品集数据"
            enableColumnDragging
            enableColumnFilters
            getRowId={(collection) => String(collection.id)}
            loading={isLoading}
          />
          <Pagination
            className="mt-4"
            page={page}
            pageSize={pageSize}
            total={total}
            onPageChange={setPage}
            onPageSizeChange={(nextPageSize) => {
              setPageSize(nextPageSize)
              setPage(1)
            }}
          />
        </CardContent>
      </Card>

      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-3xl max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">{editingId ? '编辑作品集' : '新增作品集'}</h2>
              <button onClick={() => setShowDialog(false)} className="p-1 hover:bg-muted rounded"><X className="h-5 w-5" /></button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">标题 *</label>
                  <Input value={formData.title} onChange={(event) => setFormData({ ...formData, title: event.target.value })} required />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">封面 URL</label>
                  <Input value={formData.coverImage} onChange={(event) => setFormData({ ...formData, coverImage: event.target.value })} />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">描述</label>
                <textarea
                  value={formData.description}
                  onChange={(event) => setFormData({ ...formData, description: event.target.value })}
                  rows={3}
                  className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                />
              </div>
              <label className="flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={formData.isPublished}
                  onChange={(event) => setFormData({ ...formData, isPublished: event.target.checked })}
                  className="h-4 w-4"
                />
                发布作品集
              </label>

              {editingId && (
                <div className="border-t pt-4 space-y-4">
                  <h3 className="font-medium">作品条目</h3>
                  <div className="grid grid-cols-[1fr_120px_auto] gap-3 items-end">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">诗词</label>
                      <Combobox
                        options={poemOptions}
                        value={selectedPoemId}
                        onChange={(value) => setSelectedPoemId(typeof value === 'number' ? value : undefined)}
                        placeholder="选择诗词"
                      />
                    </div>
                    <div className="space-y-2">
                      <label className="text-sm font-medium">排序</label>
                      <Input type="number" value={sortOrder} onChange={(event) => setSortOrder(event.target.value)} />
                    </div>
                    <Button type="button" onClick={handleAddItem} disabled={!selectedPoemId || addItemMutation.isPending}>
                      添加
                    </Button>
                  </div>
                  <div className="border rounded-md overflow-hidden">
                    {(detail?.items?.length ?? 0) === 0 ? (
                      <div className="text-center py-6 text-muted-foreground text-sm">暂无作品条目</div>
                    ) : (
                      <Table>
                        <TableHeader>
                          <TableRow className="bg-muted/40">
                            <TableHead className="h-10 px-3 text-sm">类型</TableHead>
                            <TableHead className="h-10 px-3 text-sm">作品</TableHead>
                            <TableHead className="h-10 px-3 text-sm">排序</TableHead>
                            <TableHead className="h-10 px-3 text-sm">操作</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {detail!.items!.map((item) => (
                            <TableRow key={item.id}>
                              <TableCell className="px-3 py-2 text-sm">{item.workType}</TableCell>
                              <TableCell className="px-3 py-2 text-sm">{poemTitleById.get(item.workId) || `#${item.workId}`}</TableCell>
                              <TableCell className="px-3 py-2">
                                <Input
                                  type="number"
                                  className="h-8 w-24"
                                  defaultValue={item.sortOrder}
                                  onBlur={(event) => {
                                    const nextOrder = Number(event.target.value) || 0
                                    if (nextOrder !== item.sortOrder) {
                                      updateItemMutation.mutate({ collectionId: editingId, itemId: item.id, data: { sortOrder: nextOrder } })
                                    }
                                  }}
                                />
                              </TableCell>
                              <TableCell className="px-3 py-2">
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => deleteItemMutation.mutate({ collectionId: editingId, itemId: item.id })}
                                >
                                  <Trash2 className="h-4 w-4 text-cinnabar" />
                                </Button>
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    )}
                  </div>
                </div>
              )}

              <div className="flex justify-end gap-2 pt-2">
                <Button type="button" className="bg-secondary text-secondary-foreground hover:bg-secondary/90" onClick={() => setShowDialog(false)}>取消</Button>
                <Button type="submit" disabled={createMutation.isPending || updateMutation.isPending}>
                  {(createMutation.isPending || updateMutation.isPending) ? '提交中...' : editingId ? '保存' : '新增'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
      <ConfirmDialog />
    </div>
  )
}
