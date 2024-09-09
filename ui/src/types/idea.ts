// @ts-expect-error not a ts module
import type EditorJS from '@editorjs/editorjs'
import { apiGet, apiPost, apiPut } from '../lib/api.ts'

export function newIdea (): EditableIdea {
  return { id: null, title: '', description: null }
}

export interface EditableIdea {
  id: string | null
  title: string
  description: EditorJS.OutputData | null
}

export interface ApiIdea {
  id: string | null
  title: string
  description: string | null
}

export interface ListIdea {
  id: string
  title: string
  likes: number
  status: string
}

function toEditableIdea (idea: ApiIdea | null): EditableIdea | null {
  return idea && { ...idea, description: JSON.parse(idea.description ?? '{}') }
}

function toApiIdea (idea: EditableIdea): ApiIdea | null {
  return idea && { ...idea, description: JSON.stringify(idea.description ?? {}) }
}

export async function fetchIdea (id: number): Promise<EditableIdea | null> {
  const data = await apiGet<ApiIdea>(`/api/idea/${id}`)
  return toEditableIdea(data)
}

export async function fetchIdeas (page: number, currentUser: boolean): Promise<ListIdea[]> {
  const res = await apiGet<ListIdea[]>(`/api/idea?page=${page}&size=20&currentUser=${currentUser}`)
  if (!res) {
    return []
  }
  return res
}

export async function saveIdea (idea: EditableIdea): Promise<EditableIdea | null> {
  const isNew = !!idea.id
  const path = isNew ? '/api/idea' : `/api/idea/${idea.id}`
  const apiIdea = toApiIdea(idea)
  return await (isNew ? apiPost(path, apiIdea) : apiPut(path, apiIdea))
}
