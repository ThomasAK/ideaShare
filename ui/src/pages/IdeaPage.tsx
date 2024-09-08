import { useParams } from 'react-router-dom'
import { type ReactNode } from 'react'

export default function IdeaPage (): ReactNode {
  const params = useParams()
  return `idea: ${params.id ?? 'no-id'}`
}
