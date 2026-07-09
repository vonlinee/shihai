import { useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Search, CheckCircle, XCircle, Eye } from 'lucide-react'

const mockCorrections = [
  { id: 1, poem: '静夜思', type: 'content', original: '床前明月光', suggested: '窗前明月光', user: '张三', status: 'pending', createdAt: '2024-01-15' },
  { id: 2, poem: '春晓', type: 'translation', original: '春眠不觉晓', suggested: '春天睡觉不知不觉天就亮了', user: '李四', status: 'approved', createdAt: '2024-01-14' },
  { id: 3, poem: '登鹳雀楼', type: 'appreciation', original: '赏析内容有误', suggested: '更详细的赏析内容', user: '王五', status: 'rejected', createdAt: '2024-01-13' },
  { id: 4, poem: '水调歌头', type: 'content', original: '明月几时有', suggested: '明月何时有', user: '赵六', status: 'pending', createdAt: '2024-01-12' },
]

const typeLabels: Record<string, string> = {
  content: '原文',
  translation: '译文',
  appreciation: '赏析',
  annotation: '注释',
}

export function AdminCorrectionsPage() {
  const [searchQuery, setSearchQuery] = useState('')

  return (
    <div className="p-8 space-y-6">
      <Card className="ink-border">
        <CardHeader className="pb-4">
          <div className="flex items-center gap-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="搜索诗词或用户..."
                className="pl-10"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>诗词</TableHead>
                <TableHead>纠错类型</TableHead>
                <TableHead>原文</TableHead>
                <TableHead>建议修改</TableHead>
                <TableHead>提交用户</TableHead>
                <TableHead>状态</TableHead>
                <TableHead>操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
                {mockCorrections.map((correction) => (
                  <TableRow key={correction.id}>
                    <TableCell>{correction.id}</TableCell>
                    <TableCell className="font-medium">{correction.poem}</TableCell>
                    <TableCell>{typeLabels[correction.type]}</TableCell>
                    <TableCell className="max-w-xs truncate text-muted-foreground">{correction.original}</TableCell>
                    <TableCell className="max-w-xs truncate">{correction.suggested}</TableCell>
                    <TableCell>{correction.user}</TableCell>
                    <TableCell>
                      <span className={`px-2 py-1 rounded text-xs ${
                        correction.status === 'approved' ? 'bg-green-500/10 text-green-500' :
                        correction.status === 'pending' ? 'bg-yellow-500/10 text-yellow-500' :
                        'bg-red-500/10 text-red-500'
                      }`}>
                        {correction.status === 'approved' ? '已通过' : correction.status === 'pending' ? '待审核' : '已驳回'}
                      </span>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Button variant="ghost" size="sm">
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="sm">
                          <CheckCircle className="h-4 w-4 text-green-500" />
                        </Button>
                        <Button variant="ghost" size="sm">
                          <XCircle className="h-4 w-4 text-cinnabar" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}
