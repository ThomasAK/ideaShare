import React from 'react'

export default function ListPage ({ currentUser }: { currentUser?: boolean }): React.ReactNode {
  return currentUser ? 'myIdeas' : 'allIdeas'
}
