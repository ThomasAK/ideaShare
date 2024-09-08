import EditorJS from '@editorjs/editorjs'
import { type ReactNode, useEffect, useState } from 'react'
import { Paper, useTheme, colors } from '@mui/material'
import { EditorJSTools } from './tools.ts'
// @ts-expect-error not a ts module
import Undo from 'editorjs-undo'

interface props { id: string, readOnly?: boolean, placeHolder: string, style: object, editorCreatedCb: (editorJs: EditorJS) => void, data: EditorJS.OutputData | null }
export default function Editor ({ readOnly, placeHolder, id, style, editorCreatedCb, data }: props): ReactNode {
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
      // @ts-expect-error the type is wrong...
      tools: EditorJSTools,
      // @ts-expect-error null is fine
      data,
      onReady: () => {
        // eslint-disable-next-line no-new
        new Undo({ editor: editorJS })
      }
    })
    setEditor(editorJS)
    editorCreatedCb(editorJS)
  })
  const darkMode = theme.palette.mode === 'dark'
  const editorStyles = `
    .ce-toolbar__actions > *:hover {
      background-color: ${darkMode ? colors.grey['600'] : ''};
    }
    .ce-toolbar__actions > * {
      background-color: ${darkMode ? colors.grey['500'] : 'white'};
    }
    .tc-add-column:hover {
      background-color: ${darkMode ? colors.grey['500'] : '#f9f9fb'}
    }
    
    .tc-add-row:hover {
      background-color: ${darkMode ? colors.grey['500'] : '#f9f9fb'}
    }
    .tc-add-row:hover:before {
      background-color: ${darkMode ? colors.grey['500'] : '#f9f9fb'}
    }
    .tc-cell--selected {
      background-color: ${darkMode ? colors.grey['500'] : '#f9f9fb'}
    }
    .tc-row--selected {
      background-color: ${darkMode ? colors.grey['500'] : '#f9f9fb'}
    }
    
    .tc-wrap {
      --color-background: ${darkMode ? 'unset' : '#f9f9fb'}
    }
    .tc-popover *{
      color: black;
    }
    
    .ce-popover__container {
      color: black;
    }
    .ce-block--selected .ce-block__content{
      background-color: ${darkMode ? colors.grey['500'] : '#e1f2ff'}
    }
    .ce-inline-tool-input {
      color: black;    
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
