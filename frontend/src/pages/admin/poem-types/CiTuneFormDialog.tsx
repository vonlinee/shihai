import { type FormEvent, useState } from 'react'
import { toast } from 'sonner'
import { X } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Combobox, type ComboboxOption } from '@/components/ui/combobox'
import { Input } from '@/components/ui/input'
import type { CiTune } from '@/types'

export interface CiTuneFormPayload {
  name: string
  aliases: string[]
  poemTypeId?: string
  description: string
}

interface CiTuneFormDialogProps {
  editingCiTune?: CiTune | null
  isSubmitting: boolean
  poemTypeOptions: ComboboxOption[]
  onClose: () => void
  onSubmit: (payload: CiTuneFormPayload) => void
}

function formatAliasesForInput(aliases: string[] | undefined): string {
  return (aliases ?? []).join('\n')
}

function parseAliasesInput(value: string): string[] {
  const aliases = value
    .split(/[\n,，、]/)
    .map((item) => item.trim())
    .filter(Boolean)
  return Array.from(new Set(aliases))
}

export function CiTuneFormDialog({
  editingCiTune,
  isSubmitting,
  poemTypeOptions,
  onClose,
  onSubmit,
}: CiTuneFormDialogProps) {
  const [formData, setFormData] = useState({
    name: editingCiTune?.name ?? '',
    aliases: formatAliasesForInput(editingCiTune?.aliases),
    poemTypeId: editingCiTune?.poemTypeId,
    description: editingCiTune?.description ?? '',
  })

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault()
    const name = formData.name.trim()
    if (!name) {
      toast.error('词牌名不能为空')
      return
    }
    onSubmit({
      name,
      aliases: parseAliasesInput(formData.aliases),
      poemTypeId: formData.poemTypeId,
      description: formData.description,
    })
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-background rounded-lg w-full max-w-lg max-h-[90vh] flex flex-col">
        <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
          <h2 className="text-xl font-bold font-serif">
            {editingCiTune ? '编辑词牌' : '添加词牌'}
          </h2>
          <button
            type="button"
            onClick={onClose}
            className="p-1 hover:bg-muted rounded"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4 p-6 pt-4 overflow-y-auto">
          <div className="space-y-2">
            <label className="text-sm font-medium">词牌名 *</label>
            <Input
              value={formData.name}
              onChange={(event) => setFormData({ ...formData, name: event.target.value })}
              placeholder="如：水调歌头"
              required
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">所属词体裁</label>
            <Combobox
              options={poemTypeOptions}
              value={formData.poemTypeId}
              onChange={(value) =>
                setFormData({ ...formData, poemTypeId: value ? String(value) : undefined })
              }
              placeholder="选择小令 / 中调 / 长调等"
              allowCustom={false}
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">别名</label>
            <textarea
              value={formData.aliases}
              onChange={(event) => setFormData({ ...formData, aliases: event.target.value })}
              placeholder="每行一个，也支持逗号分隔"
              rows={4}
              className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">说明</label>
            <textarea
              value={formData.description}
              onChange={(event) => setFormData({ ...formData, description: event.target.value })}
              placeholder="词牌说明"
              rows={3}
              className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button
              type="button"
              className="bg-secondary text-secondary-foreground hover:bg-secondary/90"
              onClick={onClose}
            >
              取消
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? '提交中...' : editingCiTune ? '保存' : '添加'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
