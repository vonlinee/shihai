import { useCallback, useState } from 'react'
import * as AlertDialog from '@radix-ui/react-alert-dialog'

import { Button } from '@/components/ui/button'

interface ConfirmDialogOptions {
  title: string
  description?: string
  confirmText?: string
  cancelText?: string
  destructive?: boolean
}

interface ConfirmDialogState {
  options: ConfirmDialogOptions
  resolve: (confirmed: boolean) => void
}

export function useConfirmDialog() {
  const [state, setState] = useState<ConfirmDialogState | null>(null)

  const confirm = useCallback((options: ConfirmDialogOptions) => {
    return new Promise<boolean>((resolve) => {
      setState({ options, resolve })
    })
  }, [])

  const close = useCallback(
    (confirmed: boolean) => {
      state?.resolve(confirmed)
      setState(null)
    },
    [state],
  )

  const ConfirmDialog = useCallback(() => {
    const options = state?.options

    return (
      <AlertDialog.Root open={Boolean(state)} onOpenChange={(open) => !open && close(false)}>
        <AlertDialog.Portal>
          <AlertDialog.Overlay className="fixed inset-0 z-50 bg-black/50" />
          <AlertDialog.Content className="fixed left-1/2 top-1/2 z-50 w-[calc(100vw-2rem)] max-w-md -translate-x-1/2 -translate-y-1/2 rounded-lg border bg-background p-6 shadow-lg">
            <AlertDialog.Title className="text-lg font-semibold text-foreground">
              {options?.title}
            </AlertDialog.Title>
            {options?.description && (
              <AlertDialog.Description className="mt-2 text-sm leading-6 text-muted-foreground">
                {options.description}
              </AlertDialog.Description>
            )}
            <div className="mt-6 flex justify-end gap-2">
              <Button
                type="button"
                variant="secondary"
                onClick={() => close(false)}
              >
                {options?.cancelText ?? '取消'}
              </Button>
              <Button
                type="button"
                variant={options?.destructive ? 'destructive' : 'default'}
                onClick={() => close(true)}
              >
                {options?.confirmText ?? '确定'}
              </Button>
            </div>
          </AlertDialog.Content>
        </AlertDialog.Portal>
      </AlertDialog.Root>
    )
  }, [close, state])

  return { confirm, ConfirmDialog }
}
