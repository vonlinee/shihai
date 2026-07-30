import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { CiPoemTypesTab } from './poem-types/CiPoemTypesTab'
import { PoemTypeCategoryTab } from './poem-types/PoemTypeCategoryTab'

const poemTypeCategories = [
  { value: 'poetry', label: '诗', category: '诗' },
  { value: 'ci', label: '词', category: '词' },
] as const

export function PoemTypesTab() {
  return (
    <Tabs defaultValue={poemTypeCategories[0].value} className="space-y-4">
      <TabsList className="grid w-64 grid-cols-2">
        {poemTypeCategories.map((item) => (
          <TabsTrigger key={item.value} value={item.value}>
            {item.label}
          </TabsTrigger>
        ))}
      </TabsList>

      {poemTypeCategories.map((item) => (
        <TabsContent key={item.value} value={item.value} className="space-y-0">
          {item.category === '词' ? (
            <CiPoemTypesTab />
          ) : (
            <PoemTypeCategoryTab category={item.category} />
          )}
        </TabsContent>
      ))}
    </Tabs>
  )
}
