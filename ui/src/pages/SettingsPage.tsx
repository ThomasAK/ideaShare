import React, { type ReactNode, useEffect } from 'react'
import {
  Checkbox,
  List,
  ListItem,
  Typography
} from '@mui/material'

export default function SettingsPage (): ReactNode {
  const [checked, setChecked] = React.useState([0])
  const settings = [
    {
      name: 'Email new ideas',
      id: 1,
      valid: true
    },
    {
      name: 'Send Weekly Summary',
      id: 2,
      valid: false
    },
    {
      name: 'Slack Notifications',
      id: 3,
      valid: true
    }
  ]

  useEffect(() => {
    const newChecked = [...checked]
    settings.forEach(setting => setting.valid ? newChecked.push(setting.id) : null)
    setChecked(newChecked)
  }, [])

  const handleToggle = (setting: { name?: string, id: number, valid?: boolean }) => () => {
    const currentIndex = checked.indexOf(setting.id)
    const newChecked = [...checked]

    if (currentIndex === -1) {
      newChecked.push(setting.id)
    } else {
      newChecked.splice(currentIndex, 1)
    }

    (setting.valid === true) ? setting.valid = false : setting.valid = true

    setChecked(newChecked)
  }
  return (
      <>
        <List className="list" sx={{ width: '40%', margin: 'auto' }}>
            {settings.map(setting =>
                <ListItem key={setting.id} sx={{ width: '100%', display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid grey' }}>
                    <Typography gutterBottom variant="h6" component="div">
                        {setting.name}
                    </Typography>
                    <Checkbox
                        edge="end"
                        onChange={handleToggle(setting)}
                        checked={checked.includes(setting.id)}
                        inputProps={{ 'aria-labelledby': '10' }}/>
                </ListItem>
            )}
        </List>
      </>
  )
}
