import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import PostList from "./components/PostList";
import PostForm from "./components/PostForm";
import PostDetail from "./components/PostDetail";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<PostList />} />
        <Route path="/posts/:id" element={<PostDetail />} />
        <Route path="/edit/:id" element={<PostForm />} /> {/* 編集用 */}
        <Route path="/create" element={<PostForm />} /> {/* 新規作成用 */}
      </Routes>
    </Router>
  );
}

export default App;
