import { useLocation, useNavigate, useParams } from 'react-router-dom'
import { type ReactNode, useEffect, useState } from 'react'
import { type EditableIdea, fetchIdea, newIdea } from '../types/idea.ts'
import Editor from '../component/Editor.tsx'
import { Fab } from '@mui/material'
import EditIcon from '@mui/icons-material/Edit'

export default function IdeaPage (): ReactNode {
  const params = useParams()
  const location = useLocation()
  const [idea, setIdea] = useState<EditableIdea>(newIdea())
  const navigate = useNavigate()
  const editButtonSize = '28px'
  useEffect(() => {
    if (location?.state?.idea && location.state?.idea?.id === params.id) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
      setIdea({ ...location.state.idea, description: JSON.parse(location.state.idea.description) })
    } else {
      fetchIdea(parseInt(params.id ?? '-1'))
        .then(i => { i && setIdea(i) })
        .catch(e => {
          console.error(`failed to load idea with id ${params.id}: ${e}`)
        })
    }
  }, [])
  return (
    <div style={{ width: '100%' }}>
      <Fab size="small"
           sx={{ position: 'absolute', right: '.5rem', width: editButtonSize, height: editButtonSize, minWidth: editButtonSize, minHeight: editButtonSize }}
           onClick={() => { navigate(`/idea/${idea.id}/edit`, { state: { idea } }) }}
      >
        <EditIcon/>
      </Fab>
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
