import { type ReactNode } from 'react'

export default function ListPage ({ currentUser }: { currentUser?: boolean }): ReactNode {
  return currentUser ?? false ? 'myIdeas' : 'allIdeas'
}
