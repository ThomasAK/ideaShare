import ThemeToggle from './ThemeToggle.tsx'
import { JSX, useState } from 'react'
import {
  AppBar,
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography
} from '@mui/material'
import TipsAndUpdatesOutlinedIcon from '@mui/icons-material/TipsAndUpdatesOutlined'
import { Pages } from '../page.ts'
import KeyboardDoubleArrowLeftIcon from '@mui/icons-material/KeyboardDoubleArrowLeft'
import KeyboardDoubleArrowRightIcon from '@mui/icons-material/KeyboardDoubleArrowRight'
import { useNavigate } from 'react-router-dom'

export default function Layout ({ pages }: { pages: Pages }): JSX.Element {
  const [collapsed, setCollapsed]: [boolean, Function] = useState(false)
  const navigate = useNavigate()
  function handleCollapse (): void {
    setCollapsed(!collapsed)
    if (!collapsed) {
      // @ts-expect-error
      document.getElementById('root').style.marginLeft = '40px'
    } else {
      // @ts-expect-error
      document.getElementById('root').style.marginLeft = '82px'
    }
  }
  function handleNavigate (path: string): void {
    navigate(path)
  }
  return (
    <div>
      <AppBar sx={{ zIndex: '20000' }}>
        <Toolbar>
          <Box sx={{ flexGrow: 1, flexDirection: 'row', display: 'flex', alignItems: 'center' }}>
            <TipsAndUpdatesOutlinedIcon />
            <Typography
              variant='h6' sx={{
                ml: 2,
                letterSpacing: '.3rem'
              }}
            >
              IdeaShare
            </Typography>
          </Box>
          <ThemeToggle />
        </Toolbar>
      </AppBar>
      <Drawer variant='permanent' anchor='left' sx={{ position: 'relative' }} PaperProps={{ sx: { justifyContent: 'center' } }}>
        <List sx={{ minWidth: '24px' }}>
          {pages.filter(p => !(p.icon == null)).map((page) => (
            <ListItem onClick={() => handleNavigate(page.path)} key={page.name} sx={{ paddingLeft: '0', paddingRight: '0' }}>
              <ListItemButton sx={{ flexDirection: 'column', paddingLeft: '.5rem', paddingRight: '.5rem' }}>
                <ListItemIcon sx={{ minWidth: '0' }}>
                  {page.icon}
                </ListItemIcon>
                {collapsed ? '' : <ListItemText primary={page.name} />}
              </ListItemButton>
            </ListItem>
          ))}
        </List>
        <div onClick={handleCollapse} style={{ position: 'absolute', bottom: 0, cursor: 'pointer' }}>
          {collapsed ? <KeyboardDoubleArrowRightIcon /> : <KeyboardDoubleArrowLeftIcon />}
        </div>
      </Drawer>
    </div>
  )
}
