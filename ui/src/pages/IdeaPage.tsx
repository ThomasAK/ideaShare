import { useParams } from 'react-router-dom'
import React from 'react'

export default function IdeaPage (): React.ReactNode {
  const params = useParams()
  return `idea: ${params.id}`
}
