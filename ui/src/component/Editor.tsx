import EditorJS from '@editorjs/editorjs'
import { type ReactNode, useEffect, useState } from 'react'
import { Paper, useTheme, colors } from '@mui/material'
import { EditorJSTools } from './tools.ts'

interface props { id: string, readOnly?: boolean, placeHolder: string, style: object, editorCreatedCb: (editorJs: EditorJS) => void }
export default function Editor ({ readOnly, placeHolder, id, style, editorCreatedCb }: props): ReactNode {
  const theme = useTheme()
  const [editor, setEditor] = useState<EditorJS | null>(null)
  useEffect(() => {
    const editorDiv = document.getElementById(id)
    // if the element gets re-mounted it will duplicate the EditorJS instance
    if ((editorDiv?.hasAttribute('editorJS') ?? false)) {
      if (editor !== null) {
        editorCreatedCb(editor)
      }
      return
    }
    editorDiv?.setAttribute('editorJS', 'true')
    const editorJS = new EditorJS({
      holder: id,
      autofocus: true,
      readOnly,
      placeholder: placeHolder,
      tools: EditorJSTools
    })
    setEditor(editorJS)
    editorCreatedCb(editorJS)
  })
  const darkMode = theme.palette.mode === 'dark'
  const editorStyles = `
    .ce-toolbar svg{
      color: ${darkMode ? 'white' : ''};
    }
    .ce-toolbar__actions > *:hover {
      background-color: ${darkMode ? colors.grey['600'] : ''};
    }
    .ce-toolbar__actions > * {
      background-color: ${darkMode ? colors.grey['500'] : 'white'};
    }
    
    .codex-editor__redactor {
      padding-bottom: 1rem !important;
    }
  `
  return (
    <Paper elevation={2} sx={{ ...style, backgroundColor: darkMode ? colors.grey[800] : '' }}>
      <style>
        {editorStyles}
      </style>
      <div id={id} />
    </Paper>
  )
}
