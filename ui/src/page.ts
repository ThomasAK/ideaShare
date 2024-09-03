import { JSX } from 'react'

export type Pages = Page[]

export interface Page {
  name: string
  icon: JSX.Element
  path: string
}
