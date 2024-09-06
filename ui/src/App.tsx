import { JSX } from 'react'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import Layout from './component/Layout.tsx'
import { Pages } from './page.ts'
import LightbulbTwoToneIcon from '@mui/icons-material/LightbulbTwoTone'
import LightbulbRoundedIcon from '@mui/icons-material/LightbulbRounded'
import BarChartIcon from '@mui/icons-material/BarChart'
import SettingsIcon from '@mui/icons-material/Settings'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import ListPage from './pages/ListPage.tsx'
import IdeaPage from './pages/IdeaPage.tsx'
import MetricsPage from './pages/MetricsPage.tsx'
import SettingsPage from './pages/SettingsPage.tsx'

const theme = createTheme({ colorSchemes: { dark: true } })

const appPages: Pages = [
  {
    name: 'Ideas',
    icon: <LightbulbTwoToneIcon />,
    path: '/',
    element: <ListPage />,
    errorElement: <div>Page Not Found</div>
  },
  {
    name: 'My Ideas',
    icon: <LightbulbRoundedIcon />,
    path: '/my-ideas',
    element: <ListPage currentUser />
  },
  {
    name: 'Metrics',
    icon: <BarChartIcon />,
    path: '/metrics',
    element: <MetricsPage />
  },
  {
    name: 'Settings',
    icon: <SettingsIcon />,
    path: '/settings',
    element: <SettingsPage />
  },
  {
    name: 'Idea',
    path: '/idea/:id',
    element: <IdeaPage />
  }
]

function App (): JSX.Element {
  return (
    <ThemeProvider theme={theme} disableTransitionOnChange>
      <BrowserRouter>
        <CssBaseline enableColorScheme />
        <Layout pages={appPages} />
        <Routes>
          {appPages.map(p => <Route key={p.path} path={p.path} element={p.element} errorElement={p.errorElement} />)}
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
