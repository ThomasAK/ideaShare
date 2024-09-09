import { useEffect, useState, type ReactNode } from 'react'
import {
  List,
  ListItem,
  ListItemButton,
  IconButton,
  Divider,
  Typography, Chip
} from '@mui/material'
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder'
import InfiniteScroll from 'react-infinite-scroll-component'

interface ListIdea {
  title: string
  id: number
  status: string
  likes: number
}

export default function ListPage ({ currentUser }: { currentUser?: boolean }): ReactNode {
  const [ideas, setIdeas] = useState<ListIdea[]>([])
  const [hasMore, setHasMore] = useState(true)
  const [index, setIndex] = useState(1)

  function fetchMoreData (): void {
    fetch(`/api/idea?page=${index}&size=20`)
      .then(async response => await response.json())
      .then(json => {
        if (!Array.isArray(json)) {
          throw new Error('invalid data')
        }
        setIdeas(ideas.concat(json))
        setIndex(index + 1)
        setHasMore(json.length > 0)
      })
      .catch(error => { console.error(error) })
  }

  useEffect(() => {
    fetchMoreData()
  }, [])
  return (
        <>
            <List className="list" sx={{ width: '100%' }}>
                <InfiniteScroll
                    dataLength={ideas.length}
                    next={fetchMoreData}
                    hasMore={hasMore}
                    loader={<h4>Loading...</h4>}
                >
                {ideas.map(idea =>
                <ListItem disablePadding key={idea.id} sx={{ width: '100%' }}>
                        <ListItemButton sx={{ width: '100%', display: 'flex', justifyContent: 'space-between' }}>
                                <Typography gutterBottom variant="h6" component="div">
                                    {idea.title + idea.id}
                                </Typography>
                            <span>
                                <span>{idea.likes > 1000 ? (idea.likes / 100) + 'K' : idea.likes }</span>
                                <IconButton aria-label="Favorite">
                                    <FavoriteBorderIcon/>
                                </IconButton>
                                <Chip sx={{ marginLeft: '2rem' }} variant="outlined" label={idea.status}/>
                            </span>
                        </ListItemButton>
                </ListItem>
                )}
                <Divider/>
                </InfiniteScroll>
            </List>
        </>
  )
}
