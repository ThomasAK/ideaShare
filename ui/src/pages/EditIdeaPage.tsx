import { type FormEvent, type ReactNode, useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import Editor from '../component/Editor.tsx'
import { Button, TextField } from '@mui/material'
import type EditorJS from '@editorjs/editorjs'

interface EditableIdea {
  title: string
  data: EditorJS.OutputData | null
}

export default function EditIdeaPage (): ReactNode {
  const params = useParams()
  const [idea, setIdea] = useState<EditableIdea>({
    title: '',
    data: null
  })
  const ideaId = params.id
  let editor: EditorJS | null = null
  function editorCreated (editorJS: EditorJS): void {
    editor = editorJS
  }
  useEffect(() => {
    if (ideaId && ideaId !== 'new') {
      fetch(`/api/idea/${ideaId}`)
        .then(async r => await r.json())
        .then(setIdea)
        .catch(console.error)
    }
  })
  async function handleSubmit (e: FormEvent<HTMLFormElement>): Promise<boolean> {
    e.preventDefault()
    const ideaData = await editor?.save()
    console.dir(ideaData)
    return false
  }
  // sx prop was getting overridden without this
  const styles = `
  input {
    text-align: center;
  }
  `
  return (
    <div style={{ width: '100%' }}>
      <form onSubmit={e => { handleSubmit(e).catch(console.error) }}
            style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}
      >
        <style>{styles}</style>
        <TextField label="Title" variant="outlined" fullWidth value={idea.title} />
        <Editor
          data={idea.data}
          style={{ width: '100%', minWidth: '240px', padding: '1rem' }}
          editorCreatedCb={editorCreated}
          id='edit-idea-editor'
          placeHolder='Write Your New Idea Here!'
        />
        <Button variant='contained' type='submit'>Save</Button>
      </form>
    </div>
  )
}
