// @ts-expect-error not a ts module
import type EditorJS from '@editorjs/editorjs'

export function newIdea (): EditableIdea {
  return { id: null, title: '', description: null }
}

export interface EditableIdea {
  id: string | null
  title: string
  description: EditorJS.OutputData | null
}
