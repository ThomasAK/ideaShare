import { type FormEvent, type ReactNode, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import Editor from '../component/Editor.tsx'
import { Button, TextField } from '@mui/material'
import type EditorJS from '@editorjs/editorjs'
import { type EditableIdea, newIdea } from '../types/idea.ts'

export default function EditIdeaPage (): ReactNode {
  const params = useParams()
  const [idea, setIdea] = useState<EditableIdea>(newIdea())
  const navigate = useNavigate()
  const ideaId = params.id ?? 'new'
  const isNew = ideaId === 'new'
  let editor: EditorJS | null = null
  function editorCreated (editorJS: EditorJS): void {
    editor = editorJS
  }
  useEffect(() => {
    if (!isNew) {
      fetch(`/api/idea/${ideaId}`)
        .then(async r => await r.json())
        .then(setIdea)
        .catch(console.error)
    }
  })
  async function handleSubmit (e: FormEvent<HTMLFormElement>): Promise<boolean> {
    e.preventDefault()
    const ideaData = await editor?.save()
    const url = isNew ? '/api/idea' : `/api/idea/${ideaId}`
    const method = isNew ? 'POST' : 'PUT'
    const resp = await fetch(url, {
      method,
      headers: {
        'content-type': 'application/json'
      },
      body: JSON.stringify({
        id: idea.id,
        title: idea.title,
        description: JSON.stringify(ideaData)
      })
    })
    if (!resp.ok) {
      console.error(`failed to ${method} idea with id ${ideaId}: ${await resp.text()}`)
      return false
    }
    const saved = await resp.json()
    navigate(`/idea/${saved.id}`, { state: { idea: saved } })
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
        <TextField label="Title" variant="outlined" fullWidth value={idea.title} onChange={(e) => { setIdea({ ...idea, title: e.target.value }) }}/>
        <Editor
          data={idea.description}
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
