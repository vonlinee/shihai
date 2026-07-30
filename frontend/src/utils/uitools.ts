function toggleSelectedId(selectedIds: string[], id: string): string[] {
  return selectedIds.includes(id)
    ? selectedIds.filter((selectedId) => selectedId !== id)
    : [...selectedIds, id];
}

function toggleAllVisibleIds(selectedIds: string[], visibleIds: string[]): string[] {
  const allVisibleSelected = visibleIds.length > 0 && visibleIds.every((id) => selectedIds.includes(id))
  if (allVisibleSelected) {
    return selectedIds.filter((id) => !visibleIds.includes(id))
  }
  return Array.from(new Set([...selectedIds, ...visibleIds]))
}



function toEntityId(id: string | number): string {
  return String(id)
}

export { toggleSelectedId, toEntityId, toggleAllVisibleIds };
