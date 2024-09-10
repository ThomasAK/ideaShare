import { type ReactNode } from 'react'
import { Route, Routes } from 'react-router-dom'
import { Paper } from '@mui/material'
import { type Pages } from '../types/page.ts'

export default function Page ({ pages }: { pages: Pages }): ReactNode {
  return (
    <div className='page'>
      <Paper id='edit-idea-page' sx={{ width: '100%', height: '100%', position: 'relative', zIndex: 1300 }} elevation={4}>
        <Routes>
          {pages.map(p => <Route
            key={p.path}
            path={p.path}
            element={p.element}
            errorElement={p.errorElement}
          />)}
        </Routes>
      </Paper>
    </div>
  )
}
