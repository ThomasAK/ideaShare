import { type FormEvent, type ReactNode, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import Editor from '../component/Editor.tsx'
import { Button } from '@mui/material'
import type EditorJS from '@editorjs/editorjs'

export default function EditIdeaPage ({ isNew }: { isNew?: boolean }): ReactNode {
  const params = useParams()
  useEffect(() => { console.log('editIdeaPage') })
  let editor: EditorJS | null = null
  function editorCreated (editorJS: EditorJS): void {
    editor = editorJS
  }
  async function handleSubmit (e: FormEvent<HTMLFormElement>): Promise<boolean> {
    e.preventDefault()
    console.dir(e)
    const ideaData = await editor?.save()
    console.dir(ideaData)
    return false
  }
  return (
    <div style={{ width: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
      <form onSubmit={e => { handleSubmit(e).catch(console.error) }}>
        EditPage: {isNew ?? false ? 'new' : params.id},
        <Editor
          style={{ width: '80%', minWidth: '240px', padding: '1rem' }}
          editorCreatedCb={editorCreated}
          id='edit-idea-editor'
          placeHolder='Write Your New Idea Here!'
        />
        <Button variant='contained' type='submit'>Save</Button>
      </form>
    </div>
  )
}
