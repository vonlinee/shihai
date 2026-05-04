import { BookOpen, Heart } from 'lucide-react'

export function Footer() {
  return (
    <footer className="border-t bg-muted/50">
      <div className="container py-8">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 font-serif text-xl font-bold text-ink">
            <BookOpen className="h-5 w-5 text-cinnabar" />
            <span>识海</span>
          </div>
          <p className="text-sm text-muted-foreground text-center">
            传承中华文化，品味诗词之美
          </p>
          <p className="text-sm text-muted-foreground flex items-center gap-1">
            Made with <Heart className="h-4 w-4 text-cinnabar fill-cinnabar" /> for poetry lovers
          </p>
        </div>
        <div className="mt-6 pt-6 border-t text-center text-xs text-muted-foreground">
          <p>2024 识海古诗词学习平台. All rights reserved.</p>
        </div>
      </div>
    </footer>
  )
}
