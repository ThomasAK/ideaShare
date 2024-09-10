import { type FormEvent, type ReactNode, useEffect, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import Editor from '../component/Editor.tsx'
import { Button, TextField } from '@mui/material'
import type EditorJS from '@editorjs/editorjs'
import { type EditableIdea, fetchIdea, newIdea, saveIdea } from '../types/idea.ts'
import { ApiError } from '../lib/api.ts'

export default function EditIdeaPage (): ReactNode {
  const params = useParams()
  const [idea, setIdea] = useState<EditableIdea>(newIdea())
  const location = useLocation()
  const navigate = useNavigate()
  const ideaId = parseInt(params.id ?? '-1') ?? -1
  const isNew = ideaId === -1
  // this happens when navigating from an edit page to a new idea page as the useEffect is not triggered again
  if (idea.id && ideaId !== idea.id) {
    setIdea(newIdea())
  }
  let editor: EditorJS | null = null
  function editorCreated (editorJS: EditorJS): void {
    editor = editorJS
  }
  useEffect(() => {
    if (location.state?.idea && location.state?.idea?.id === ideaId) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
      setIdea(location.state.idea)
    } else {
      if (!isNew) {
        fetchIdea(ideaId).then(i => {
          i && setIdea(i)
        })
          .catch(console.error)
      }
    }
  }, [])
  async function handleSubmit (e: FormEvent<HTMLFormElement>): Promise<boolean> {
    e.preventDefault()
    const ideaData = await editor?.save() ?? null

    try {
      const saved = await saveIdea({
        id: idea.id,
        title: idea.title,
        description: ideaData
      })
      if (!saved) {
        console.error('got empty response back from saving idea')
        return false
      }
      navigate(`/idea/${saved.id}`, { state: { idea: saved } })
    } catch (e) {
      if (e instanceof ApiError) {
        console.error(`failed to save idea with id ${ideaId}: ${e.body}`)
      }
    }

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
