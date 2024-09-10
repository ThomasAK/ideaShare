import { type ReactNode, useEffect, useState } from 'react'
import {
  Chip,
  IconButton,
  LinearProgress,
  List,
  ListItem,
  ListItemButton,
  Typography,
  useMediaQuery,
  useTheme
} from '@mui/material'
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder'
import InfiniteScroll from 'react-infinite-scroll-component'
import { useNavigate } from 'react-router-dom'
import { fetchIdeas, type ListIdea } from '../types/idea.ts'

export default function ListPage ({ currentUser }: { currentUser?: boolean }): ReactNode {
  const [ideas, setIdeas] = useState<ListIdea[]>([])
  const [hasMore, setHasMore] = useState(true)
  const [index, setIndex] = useState(1)
  const navigate = useNavigate()
  const theme = useTheme()
  const isSmall = useMediaQuery(theme.breakpoints.down('sm'))

  function fetchMoreData (): void {
    fetchIdeas(index, Boolean(currentUser))
      .then(i => {
        if (!Array.isArray(i)) {
          throw new Error('invalid data')
        }
        setIdeas(ideas.concat(i))
        setIndex(index + 1)
        setHasMore(i.length > 0)
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
                    loader={<LinearProgress />}
                >
                {ideas.map(idea =>
                <ListItem disablePadding key={idea.id} sx={{ width: '100%', borderBottom: '1px solid grey' }}>
                        <ListItemButton
                          sx={{ width: '100%', display: 'flex', justifyContent: 'space-between' }}
                          onClick={() => { navigate(`/idea/${idea.id}`) }}
                        >
                                <Typography gutterBottom variant="h6" component="div">
                                    {idea.title}
                                </Typography>
                            <span>
                                <span>{idea.likes > 1000 ? (idea.likes / 100) + 'K' : idea.likes }</span>
                                <IconButton aria-label="Favorite">
                                    <FavoriteBorderIcon/>
                                </IconButton>
                              {!isSmall && <Chip sx={{ marginLeft: '2rem' }} variant="outlined" label={idea.status}/> }
                            </span>
                        </ListItemButton>
                </ListItem>
                )}
                </InfiniteScroll>
            </List>
        </>
  )
}
