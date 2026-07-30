import { type FormEvent, useState } from 'react'
import { toast } from 'sonner'
import { X } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import type { PoemType } from '@/types'

export interface PoemTypeFormPayload {
  name: string
  category: string
  lines: number | null
  charsPerLine: number | null
  description: string
}

interface PoemTypeFormDialogProps {
  category: string
  editingPoemType?: PoemType | null
  isSubmitting: boolean
  onClose: () => void
  onSubmit: (payload: PoemTypeFormPayload) => void
}

function parseOptionalPositiveIntInput(value: string): number | null {
  const trimmed = value.trim()
  if (!trimmed) return null
  const parsed = Number(trimmed)
  if (!Number.isInteger(parsed) || parsed <= 0) return Number.NaN
  return parsed
}

export function PoemTypeFormDialog({
  category,
  editingPoemType,
  isSubmitting,
  onClose,
  onSubmit,
}: PoemTypeFormDialogProps) {
  const [formData, setFormData] = useState({
    name: editingPoemType?.name ?? '',
    lines: editingPoemType?.lines ? String(editingPoemType.lines) : '',
    charsPerLine: editingPoemType?.charsPerLine ? String(editingPoemType.charsPerLine) : '',
    description: editingPoemType?.description || '',
  })

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault()

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

    onSubmit({
      name: formData.name.trim(),
      category,
      lines,
      charsPerLine,
      description: formData.description,
    })
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-background rounded-lg w-full max-w-lg max-h-[90vh] flex flex-col">
        <div className="flex items-center justify-between p-6 pb-4 border-b shrink-0">
          <h2 className="text-xl font-bold font-serif">
            {editingPoemType ? `编辑${category}体裁` : `添加${category}体裁`}
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
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">一级分类</label>
              <Input value={category} disabled />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">细分类别 *</label>
              <Input
                value={formData.name}
                onChange={(event) =>
                  setFormData({ ...formData, name: event.target.value })
                }
                placeholder={category === '词' ? '如：浣溪沙' : '如：五言绝句'}
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
                onChange={(event) =>
                  setFormData({ ...formData, lines: event.target.value })
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
                onChange={(event) =>
                  setFormData({ ...formData, charsPerLine: event.target.value })
                }
                placeholder="不限"
              />
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">说明</label>
            <textarea
              value={formData.description}
              onChange={(event) =>
                setFormData({ ...formData, description: event.target.value })
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
              onClick={onClose}
            >
              取消
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? '提交中...' : editingPoemType ? '保存' : '添加'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
