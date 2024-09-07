import { ReactNode, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import Editor from '../component/Editor.tsx'
import { Button } from '@mui/material'

export default function EditIdeaPage ({ isNew }: { isNew?: boolean }): ReactNode {
  const params = useParams()
  useEffect(() => console.log('editIdeaPage'))
  let editor = null
  function editorCreated (editorJS): void {
    editor = editorJS
  }
  async function handleSubmit (e): Promise<boolean> {
    e.preventDefault()
    console.dir(e)
    const ideaData = await editor.save()
    console.dir(ideaData)
    return false
  }
  return (
    <div style={{ width: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
      {/* @ts-expect-error */}
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
