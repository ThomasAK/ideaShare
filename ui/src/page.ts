import { JSX } from 'react'
import * as React from 'react'

export type Pages = Page[]

export interface Page {
  name: string
  icon?: JSX.Element
  path: string
  element: React.ReactNode | null
  errorElement?: React.ReactNode | null
}
