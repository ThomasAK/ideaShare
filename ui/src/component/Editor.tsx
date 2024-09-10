import EditorJS from '@editorjs/editorjs'
import { type ReactNode, useEffect, useRef } from 'react'
import { colors, Paper, useTheme } from '@mui/material'
import { EditorJSTools } from './tools.ts'
// @ts-expect-error not a ts module
import Undo from 'editorjs-undo'

interface props { id: string, readOnly?: boolean, placeHolder?: string, style?: object, editorCreatedCb: (editorJs: EditorJS) => void, data: EditorJS.OutputData | null }
export default function Editor ({ readOnly, placeHolder, id, style, editorCreatedCb, data }: props): ReactNode {
  const theme = useTheme()
  const editor = useRef<EditorJS>()
  if (editor.current) {
    editorCreatedCb(editor.current)
  }
  useEffect(() => {
    if (editor.current) {
      editor.current?.isReady
        .then(async () => {
          const currentBlockCount = editor.current?.blocks.getBlocksCount() ?? 0
          const dataBlockCount = data?.blocks.length ?? 0
          if (dataBlockCount > 0) {
            await editor.current?.blocks.render(data ?? { blocks: [] })
          } else if ((currentBlockCount > 1 && dataBlockCount === 0) ||
            (currentBlockCount === 1 && dataBlockCount === 0 && !editor.current?.blocks.getBlockByIndex(0)?.isEmpty)) {
            // currentBlockCount is always 1 because it inserts an empty block
            editor.current?.clear()
          }
        })
        .catch(console.error)
      return
    }

    editor.current = new EditorJS({
      holder: id,
      autofocus: !(readOnly ?? false),
      readOnly,
      placeholder: placeHolder,
      // @ts-expect-error the type is wrong...
      tools: EditorJSTools,
      // @ts-expect-error null is fine
      data
    })
    editor.current.isReady.then(() => {
      // eslint-disable-next-line no-new
      new Undo({ editor: editor.current })
    }).catch(console.error)

    editorCreatedCb(editor.current)
  }, [data])
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
