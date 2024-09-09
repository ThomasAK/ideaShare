import { type ReactNode } from 'react'
import { Route, Routes, useLocation, useNavigate } from 'react-router-dom'
import { Fab, Paper } from '@mui/material'
import { Add } from '@mui/icons-material'
import { type Pages } from '../types/page.ts'

export default function Page ({ pages }: { pages: Pages }): ReactNode {
  const navigate = useNavigate()
  const location = useLocation()
  const newIdeaButton = <Fab
    color='primary' size='large' aria-label='new idea'
    sx={{ position: 'absolute', bottom: 16, right: 16 }}
    onClick={function handleAdd () {
      navigate('/idea/new')
    }}
  >
    <Add />
  </Fab>
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
        { location.pathname !== '/idea/new' ? newIdeaButton : ''}
      </Paper>
    </div>
  )
}
