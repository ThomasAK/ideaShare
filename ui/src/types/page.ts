import { type ReactNode } from 'react'

export type Pages = Page[]

export interface Page {
  name: string
  icon?: ReactNode
  path: string
  element: ReactNode | null
  errorElement?: ReactNode | null
}
