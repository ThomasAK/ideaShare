import { ReactNode } from "react";
import { Route, Routes, useNavigate } from "react-router-dom";
import { Fab, Paper } from "@mui/material";
import { Add } from "@mui/icons-material";
import { Pages } from "../page.ts";

export default function Page({pages}: { pages: Pages }): ReactNode {
  const navigate = useNavigate()
  return (
    <div className="page">
      <Paper id="edit-idea-page" sx={{width: '100%', height: '100%', position: 'relative'}} elevation={4}>
        <Routes>
          {pages.map(p => <Route key={p.path} path={p.path} element={p.element} errorElement={p.errorElement}/>)}
        </Routes>
        <Fab color="primary" size="large" aria-label="new idea"
             sx={{position: 'absolute', bottom: 16, right: 16}}
             onClick={function handleAdd() {
               navigate("/idea/new")
             }}
        >
          <Add/>
        </Fab>
      </Paper>
    </div>
  )
}