import {ReactNode} from 'react'

export default function EditIdeaPage({isNew}: {isNew?: boolean}): ReactNode {
  return (
    <div>
      EditPage: {isNew}
    </div>
  )
}
