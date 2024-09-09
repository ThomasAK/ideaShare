import { useLocation, useParams } from 'react-router-dom'
import { type ReactNode, useEffect, useState } from 'react'
import { type EditableIdea, newIdea } from '../types/idea.ts'
import Editor from '../component/Editor.tsx'

export default function IdeaPage (): ReactNode {
  const params = useParams()
  const location = useLocation()
  const [idea, setIdea] = useState<EditableIdea>(newIdea())
  useEffect(() => {
    if ((Boolean(location?.state?.idea)) && location.state?.idea?.id === params.id) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
      setIdea({ ...location.state.idea, description: JSON.parse(location.state.idea.description) })
    } else {
      fetch(`/api/idea/${params.id}`)
        .then(async r => await r.json())
        .then(i => {
          // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
          setIdea({ ...i, description: JSON.parse(i.description || '{}') })
        })
        .catch(e => {
          console.error(`failed to load idea with id ${params.id}: ${e}`)
        })
    }
  }, [])
  return (
    <div style={{ width: '100%' }}>
      <h1 style={{ width: '100%', textAlign: 'center' }}>{idea?.title}</h1>
      {((idea?.description) != null) && <Editor
        style={{ width: '100%' }}
        id={String(params.id ?? idea?.id)}
        placeHolder=""
        editorCreatedCb={() => {}}
        data={idea?.description}
        readOnly={true}
      />}
    </div>
  )
}
