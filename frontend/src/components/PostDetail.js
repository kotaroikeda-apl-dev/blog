import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Layout from "./Layout/Layout";

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [post, setPost] = useState(null);

  useEffect(() => {
    fetch(`https://www.mynoteblog.com/api/posts/${id}`)
      .then((response) => response.json())
      .then((data) => setPost(data))
      .catch((error) => console.error("Error fetching post:", error));
  }, [id]);

  if (!post) {
    return (
      <Layout>
        <p className="text-center text-gray-600">読み込み中...</p>
      </Layout>
    );
  }

  const handleDelete = () => {
    if (window.confirm("本当にこの投稿を削除しますか？")) {
      fetch(`https://www.mynoteblog.com/api/posts/${id}`, {
        method: "DELETE",
      })
        .then((response) => {
          if (response.ok) {
            alert("投稿を削除しました");
            navigate("/");
          } else {
            alert("投稿の削除に失敗しました");
          }
        })
        .catch((error) => console.error("Error deleting post:", error));
    }
  };

  return (
    <Layout>
      <div className="bg-white p-10 rounded-lg shadow-lg max-w-4xl mx-auto">
        <h1 className="text-4xl font-extrabold text-gray-800 mb-6">{post.Title}</h1>
        <p className="text-gray-700 mb-4">{post.Content}</p>
        <p className="text-sm text-gray-500 mb-6">投稿者: {post.Author}</p>
        <div className="flex space-x-4">
          <button
            onClick={() => navigate(`/edit/${id}`)}
            className="bg-blue-500 text-white font-medium px-6 py-2 rounded-lg shadow hover:bg-blue-600 transition"
          >
            編集する
          </button>
          <button
            onClick={handleDelete}
            className="bg-red-500 text-white font-medium px-6 py-2 rounded-lg shadow hover:bg-red-600 transition"
          >
            削除する
          </button>
          <button
            onClick={() => navigate("/")}
            className="bg-gray-400 text-white font-medium px-6 py-2 rounded-lg shadow hover:bg-gray-500 transition"
          >
            投稿一覧に戻る
          </button>
        </div>
      </div>
    </Layout>
  );
};

export default PostDetail;
