import { JSX } from 'react'
import { CssBaseline, ThemeProvider, createTheme } from '@mui/material'
import Layout from './component/Layout.tsx'
import { Pages } from './page.ts'
import LightbulbTwoToneIcon from '@mui/icons-material/LightbulbTwoTone'
import LightbulbRoundedIcon from '@mui/icons-material/LightbulbRounded'
import BarChartIcon from '@mui/icons-material/BarChart'
import SettingsIcon from '@mui/icons-material/Settings'

const theme = createTheme({ colorSchemes: { dark: true } })

const appPages: Pages = [
  {
    name: 'Ideas',
    icon: <LightbulbTwoToneIcon />,
    path: 'ideas'
  },
  {
    name: 'My Ideas',
    icon: <LightbulbRoundedIcon />,
    path: 'my-ideas'
  },
  {
    name: 'Metrics',
    icon: <BarChartIcon />,
    path: 'metrics'
  },
  {
    name: 'Settings',
    icon: <SettingsIcon />,
    path: 'settings'
  }
]

function App (): JSX.Element {
  return (
    <ThemeProvider theme={theme} disableTransitionOnChange>
      <CssBaseline enableColorScheme />
      <Layout pages={appPages} />
    </ThemeProvider>
  )
}

export default App
