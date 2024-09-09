import { type ReactNode } from 'react'
import { Box, Chip, Stack } from '@mui/material'
import { LineChart } from '@mui/x-charts'
export default function MetricsPage (): ReactNode {
  const statuses = [
    {
      name: 'New Idea',
      count: '312'
    },
    {
      name: 'Approved',
      count: '423'
    },
    {
      name: 'In Review',
      count: '223'
    },
    {
      name: 'Infeasible',
      count: '14'
    },
    {
      name: 'Implemented',
      count: '154'
    },
    {
      name: 'Rejected',
      count: '255'
    }
  ]
  return (
      <>
          <Stack direction={{ lg: 'row', sm: 'column' }} sx={{ justifyContent: 'space-evenly', margin: 'auto', width: '100%' }}>
        <Box sx={{ width: '80%', boxShadow: 2, margin: '2rem' }}>
            <Stack sx={{ marginLeft: '30%', marginRight: '30%', textAlign: 'center' }}>
                <h2> Ideas Submitted </h2>
                {statuses.map(status =>
                    <Stack direction="row" sx={{ justifyContent: 'space-between' }}>
                        <Chip label={status.name}/>
                        <p>{status.count}</p>
                    </Stack>
                )}
            </Stack>
        </Box>
        <Box sx={{ width: '80%', boxShadow: 2, margin: '2rem' }}>
            <Box sx={{ height: '90%' }}>
            <LineChart
              xAxis={[{ data: [1, 2, 3, 5, 8, 10] }]}
              series={[
                {
                  id: 'Ideas',
                  label: 'Ideas Per Month',
                  data: [2, -5.5, 2, -7.5, 1.5, 6],
                  area: true,
                  baseline: 'min'
                }
              ]}
          />
            </Box>
        </Box>
          </Stack>
      </>
  )
}
