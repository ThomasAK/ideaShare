import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import '@fontsource/roboto/300.css'
import '@fontsource/roboto/400.css'
import '@fontsource/roboto/500.css'
import '@fontsource/roboto/700.css'

const root = document.getElementById('root')
if (root === null) {
  throw new Error('no root')
}
createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>
)
