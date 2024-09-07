import { JSX } from 'react'
import { createTheme, CssBaseline, ThemeProvider } from '@mui/material'
import Layout from './component/Layout.tsx'
import { Pages } from './page.ts'
import LightbulbTwoToneIcon from '@mui/icons-material/LightbulbTwoTone'
import LightbulbRoundedIcon from '@mui/icons-material/LightbulbRounded'
import BarChartIcon from '@mui/icons-material/BarChart'
import SettingsIcon from '@mui/icons-material/Settings'
import { BrowserRouter } from 'react-router-dom'
import ListPage from './pages/ListPage.tsx'
import IdeaPage from './pages/IdeaPage.tsx'
import MetricsPage from './pages/MetricsPage.tsx'
import SettingsPage from './pages/SettingsPage.tsx'
import EditIdeaPage from "./pages/EditIdeaPage.tsx";
import Page from "./component/Page.tsx";

const theme = createTheme({colorSchemes: {dark: true}})


const appPages: Pages = [
  {
    name: 'Ideas',
    icon: <LightbulbTwoToneIcon/>,
    path: '/',
    element: <ListPage/>,
    errorElement: <div>Page Not Found</div>
  },
  {
    name: 'My Ideas',
    icon: <LightbulbRoundedIcon/>,
    path: '/my-ideas',
    element: <ListPage currentUser/>
  },
  {
    name: 'Metrics',
    icon: <BarChartIcon/>,
    path: '/metrics',
    element: <MetricsPage/>
  },
  {
    name: 'Settings',
    icon: <SettingsIcon/>,
    path: '/settings',
    element: <SettingsPage/>
  },
  {
    name: 'NewIdea',
    path: '/idea/new',
    element: <EditIdeaPage isNew={true}/>
  },
  {
    name: 'EditIdea',
    path: '/idea/:id/edit',
    element: <EditIdeaPage/>
  },
  {
    name: 'Idea',
    path: '/idea/:id',
    element: <IdeaPage/>
  }
]

function App(): JSX.Element {
  return (
    <ThemeProvider theme={theme} disableTransitionOnChange>
      <BrowserRouter>
        <CssBaseline enableColorScheme/>
        <Layout pages={appPages}/>
        <Page pages={appPages}/>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
