import { useMemo, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Combobox, type ComboboxOption } from "@/components/ui/combobox";
import { DataTable, type DataTableColumn } from "@/components/ui/data-table";
import { Pagination } from "@/components/ui/Pagination";
import { useConfirmDialog } from "@/components/ui/ConfirmDialog";
import {
  Search,
  Plus,
  Edit2,
  Trash2,
  X,
} from "lucide-react";
import {
  useDynasties,
  usePoetList,
} from "@/hooks/usePoems";
import {
  useAdminCreatePoet,
  useAdminUpdatePoet,
  useAdminDeletePoet,
  useAdminBatchDeletePoets,
} from "@/hooks/useAdmin";
import { toggleSelectedId } from "@/utils/uitools";

function toEntityId(id: string | number): string {
  return String(id);
}

function toggleAllVisibleIds(selectedIds: string[], visibleIds: string[]): string[] {
  const allVisibleSelected = visibleIds.length > 0 && visibleIds.every((id) => selectedIds.includes(id))
  if (allVisibleSelected) {
    return selectedIds.filter((id) => !visibleIds.includes(id))
  }
  return Array.from(new Set([...selectedIds, ...visibleIds]))
}

export default function PoetsTab() {
  const [searchQuery, setSearchQuery] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const { data: poetData, isLoading } = usePoetList({
    keyword: searchQuery || undefined,
    page,
    pageSize,
  });
  const { data: dynasties } = useDynasties();
  const createMutation = useAdminCreatePoet();
  const updateMutation = useAdminUpdatePoet();
  const deleteMutation = useAdminDeletePoet();
  const batchDeleteMutation = useAdminBatchDeletePoets();
  const { confirm, ConfirmDialog } = useConfirmDialog();
  const [showDialog, setShowDialog] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState<{
    name: string;
    dynastyId: string | undefined;
    biography: string;
    avatar: string;
    birthYear: string;
    deathYear: string;
  }>({
    name: "",
    dynastyId: undefined,
    biography: "",
    avatar: "",
    birthYear: "",
    deathYear: "",
  });
  const [selectedPoetIds, setSelectedPoetIds] = useState<string[]>([]);

  const dynastyOptions: ComboboxOption[] = (dynasties ?? []).map((d) => ({
    value: d.id,
    label: d.name,
    description: d.period,
  }));
  const poets = poetData?.list ?? [];
  const total = poetData?.total ?? 0;
  const visiblePoetIds = poets.map((p) => toEntityId(p.id));
  const allVisiblePoetsSelected =
    visiblePoetIds.length > 0 &&
    visiblePoetIds.every((id) => selectedPoetIds.includes(id));

  const openCreate = () => {
    setEditingId(null);
    setFormData({
      name: "",
      dynastyId: undefined,
      biography: "",
      avatar: "",
      birthYear: "",
      deathYear: "",
    });
    setShowDialog(true);
  };
  const openEdit = useCallback(
    (p: {
      id: string | number;
      name: string;
      dynastyId?: string | number;
      biography?: string;
      avatar?: string;
      birthYear?: number;
      deathYear?: number;
    }) => {
      setEditingId(toEntityId(p.id));
      setFormData({
        name: p.name,
        dynastyId: p.dynastyId ? toEntityId(p.dynastyId) : undefined,
        biography: p.biography || "",
        avatar: p.avatar || "",
        birthYear: p.birthYear?.toString() || "",
        deathYear: p.deathYear?.toString() || "",
      });
      setShowDialog(true);
    },
    [],
  );
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const payload = {
      name: formData.name,
      dynastyId: formData.dynastyId,
      biography: formData.biography,
      avatar: formData.avatar,
      birthYear: formData.birthYear ? parseInt(formData.birthYear) : undefined,
      deathYear: formData.deathYear ? parseInt(formData.deathYear) : undefined,
    };
    if (editingId)
      updateMutation.mutate(
        { id: editingId, data: payload },
        { onSuccess: () => setShowDialog(false) },
      );
    else
      createMutation.mutate(payload, { onSuccess: () => setShowDialog(false) });
  };
  const handleDelete = useCallback(
    async (id: string) => {
      const confirmed = await confirm({
        title: "删除诗人",
        description: "确定要删除该诗人吗？此操作不可撤销。",
        confirmText: "删除",
        destructive: true,
      });
      if (!confirmed) return;
      deleteMutation.mutate(id, {
        onSuccess: () =>
          setSelectedPoetIds((ids) =>
            ids.filter((selectedId) => selectedId !== id),
          ),
      });
    },
    [confirm, deleteMutation],
  );

  const handleBatchDelete = async () => {
    if (selectedPoetIds.length === 0) return;
    const confirmed = await confirm({
      title: "批量删除诗人",
      description: `确定要删除选中的 ${selectedPoetIds.length} 位诗人吗？此操作不可撤销。`,
      confirmText: "删除",
      destructive: true,
    });
    if (!confirmed) return;
    batchDeleteMutation.mutate(selectedPoetIds, {
      onSuccess: () => setSelectedPoetIds([]),
    });
  };

  type PoetRow = (typeof poets)[number];

  const poetColumns = useMemo<DataTableColumn<PoetRow>[]>(
    () => [
      {
        id: "select",
        header: () => (
          <Checkbox
            aria-label="选择当前页诗人"
            checked={allVisiblePoetsSelected}
            onCheckedChange={() =>
              setSelectedPoetIds((ids) =>
                toggleAllVisibleIds(ids, visiblePoetIds),
              )
            }
          />
        ),
        cell: ({ row }) => (
          <Checkbox
            aria-label={`选择诗人 ${row.original.name}`}
            checked={selectedPoetIds.includes(toEntityId(row.original.id))}
            onCheckedChange={() =>
              setSelectedPoetIds((ids) =>
                toggleSelectedId(ids, toEntityId(row.original.id)),
              )
            }
          />
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { draggable: false, fixed: "left", width: 56 },
      },
      {
        accessorKey: "name",
        header: "姓名",
        cell: ({ row }) => (
          <span className="font-medium">{row.original.name}</span>
        ),
        meta: { filterPlaceholder: "筛选姓名", fixed: "left", width: 180 },
      },
      {
        accessorFn: (poet) => poet.dynasty?.name ?? "",
        id: "dynasty",
        header: "朝代",
        cell: ({ row }) => (
          <span className="text-muted-foreground">
            {row.original.dynasty?.name || "-"}
          </span>
        ),
        meta: { filterPlaceholder: "筛选朝代", width: 150 },
      },
      {
        accessorFn: (poet) =>
          poet.birthYear && poet.deathYear
            ? `${poet.birthYear}-${poet.deathYear}`
            : "",
        id: "years",
        header: "生卒年",
        cell: ({ row }) => (
          <span className="text-muted-foreground">
            {row.original.birthYear && row.original.deathYear
              ? `${row.original.birthYear}-${row.original.deathYear}`
              : "-"}
          </span>
        ),
        meta: { filterPlaceholder: "筛选年份", width: 160 },
      },
      {
        accessorFn: (poet) => poet.biography ?? "",
        id: "biography",
        header: "简介",
        cell: ({ row }) => (
          <span className="block max-w-xs truncate text-muted-foreground">
            {row.original.biography || "-"}
          </span>
        ),
        meta: { filterPlaceholder: "筛选简介", minWidth: 260 },
      },
      {
        id: "actions",
        header: "操作",
        cell: ({ row }) => (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => openEdit(row.original)}
            >
              <Edit2 className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleDelete(toEntityId(row.original.id))}
            >
              <Trash2 className="h-4 w-4 text-cinnabar" />
            </Button>
          </div>
        ),
        enableColumnFilter: false,
        enableSorting: false,
        meta: { draggable: false, fixed: "right", width: 120 },
      },
    ],
    [
      allVisiblePoetsSelected,
      handleDelete,
      openEdit,
      selectedPoetIds,
      visiblePoetIds,
    ],
  );

  return (
    <>
      <div className="flex items-center justify-between">
        <div className="relative max-w-sm flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="搜索诗人姓名..."
            className="pl-10"
            value={searchQuery}
            onChange={(e) => {
              setSearchQuery(e.target.value);
              setSelectedPoetIds([]);
              setPage(1);
            }}
          />
        </div>
        <div className="flex items-center gap-2">
          {selectedPoetIds.length > 0 && (
            <Button
              variant="destructive"
              onClick={handleBatchDelete}
              disabled={batchDeleteMutation.isPending}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              删除选中({selectedPoetIds.length})
            </Button>
          )}
          <Button onClick={openCreate}>
            <Plus className="h-4 w-4 mr-2" />
            添加诗人
          </Button>
        </div>
      </div>
      <Card className="ink-border">
        <CardContent className="pt-6">
          {isLoading ? (
            <div className="text-center py-8 text-muted-foreground">
              加载中...
            </div>
          ) : poets.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              暂无诗人数据
            </div>
          ) : (
            <DataTable
              columns={poetColumns}
              data={poets}
              emptyText="暂无诗人数据"
              enableColumnDragging
              enableColumnFilters
              getRowId={(poet) => toEntityId(poet.id)}
            />
          )}
          <Pagination
            className="mt-4"
            page={page}
            pageSize={pageSize}
            total={total}
            onPageChange={(nextPage) => {
              setSelectedPoetIds([]);
              setPage(nextPage);
            }}
            onPageSizeChange={(nextPageSize) => {
              setSelectedPoetIds([]);
              setPageSize(nextPageSize);
              setPage(1);
            }}
          />
        </CardContent>
      </Card>

      {showDialog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-background rounded-lg w-full max-w-lg max-h-[90vh] flex flex-col">
            <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
              <h2 className="text-xl font-bold font-serif">
                {editingId ? "编辑诗人" : "添加诗人"}
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
                  <label className="text-sm font-medium">姓名 *</label>
                  <Input
                    value={formData.name}
                    onChange={(e) =>
                      setFormData({ ...formData, name: e.target.value })
                    }
                    placeholder="诗人姓名"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">朝代</label>
                  <Combobox
                    options={dynastyOptions}
                    value={formData.dynastyId}
                    onChange={(val) =>
                      setFormData({
                        ...formData,
                        dynastyId: val ? String(val) : undefined,
                      })
                    }
                    placeholder="选择朝代"
                    allowCustom={false}
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">出生年份</label>
                  <Input
                    type="number"
                    value={formData.birthYear}
                    onChange={(e) =>
                      setFormData({ ...formData, birthYear: e.target.value })
                    }
                    placeholder="如：701"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">逝世年份</label>
                  <Input
                    type="number"
                    value={formData.deathYear}
                    onChange={(e) =>
                      setFormData({ ...formData, deathYear: e.target.value })
                    }
                    placeholder="如：762"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">头像URL</label>
                <Input
                  value={formData.avatar}
                  onChange={(e) =>
                    setFormData({ ...formData, avatar: e.target.value })
                  }
                  placeholder="诗人头像链接"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">简介</label>
                <textarea
                  value={formData.biography}
                  onChange={(e) =>
                    setFormData({ ...formData, biography: e.target.value })
                  }
                  placeholder="诗人简介"
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
                  disabled={
                    createMutation.isPending || updateMutation.isPending
                  }
                >
                  {createMutation.isPending || updateMutation.isPending
                    ? "提交中..."
                    : editingId
                      ? "保存"
                      : "添加"}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
      <ConfirmDialog />
    </>
  );
}
