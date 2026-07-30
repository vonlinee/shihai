import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { CiTuneManagementTab } from './CiTuneManagementTab'
import { PoemTypeCategoryTab } from './PoemTypeCategoryTab'

export function CiPoemTypesTab() {
  return (
    <Tabs defaultValue="types" className="space-y-4">
      <TabsList className="grid w-64 grid-cols-2">
        <TabsTrigger value="types">词体裁</TabsTrigger>
        <TabsTrigger value="tunes">词牌</TabsTrigger>
      </TabsList>

      <TabsContent value="types" className="space-y-0">
        <PoemTypeCategoryTab category="词" />
      </TabsContent>
      <TabsContent value="tunes" className="space-y-0">
        <CiTuneManagementTab />
      </TabsContent>
    </Tabs>
  )
}
