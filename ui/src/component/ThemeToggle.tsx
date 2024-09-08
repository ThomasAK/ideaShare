import { useColorScheme, useTheme } from '@mui/material'
import LightModeIcon from '@mui/icons-material/LightMode'
import DarkModeOutlinedIcon from '@mui/icons-material/DarkModeOutlined'
import { type ReactNode } from 'react'

export default function ThemeToggle (): ReactNode {
  const { setMode } = useColorScheme()
  const mode = useTheme().palette.mode
  function handleToggleMode (): void {
    setMode(mode === 'dark' ? 'light' : 'dark')
  }

  return (
    <div onClick={handleToggleMode}>
      {mode === 'dark' ? <DarkModeOutlinedIcon /> : <LightModeIcon />}
    </div>
  )
}
