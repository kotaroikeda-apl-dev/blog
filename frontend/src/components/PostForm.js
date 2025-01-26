import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Layout from "./Layout/Layout";

const PostForm = () => {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [author, setAuthor] = useState("");
  const navigate = useNavigate();
  const { id } = useParams();

  useEffect(() => {
    if (id) {
      fetch(`${process.env.REACT_APP_URL_DOMAIN}/api/posts/${id}`)
        .then((response) => response.json())
        .then((data) => {
          setTitle(data.Title || "");
          setContent(data.Content || "");
          setAuthor(data.Author || "");
        })
        .catch((error) => console.error("データ取得エラー:", error));
    }
  }, [id]);

  const handleSubmit = (e) => {
    e.preventDefault();

    const post = { title, content, author };

    const method = id ? "PUT" : "POST";
    const url = id
      ? `${process.env.REACT_APP_URL_DOMAIN}/api/posts/${id}`
      : `${process.env.REACT_APP_URL_DOMAIN}/api/posts`;

    fetch(url, {
      method,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(post),
    })
      .then((response) => {
        if (response.ok) {
          alert(id ? "投稿を更新しました！" : "新規投稿が作成されました！");
          navigate("/");
        } else {
          alert("保存に失敗しました。");
        }
      })
      .catch((error) => console.error("データ送信エラー:", error));
  };

  return (
    <Layout>
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <div className="bg-white p-10 rounded-lg shadow-lg max-w-4xl w-full">
          <h1 className="text-4xl font-bold text-gray-800 mb-8 text-center">
            {id ? "投稿を編集" : "新規投稿を作成"}
          </h1>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-gray-700 font-medium mb-2">タイトル</label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="例: Reactの使い方"
                className="w-full border border-gray-300 px-4 py-3 rounded-lg shadow-sm focus:ring-4 focus:ring-blue-300"
                required
              />
            </div>
            <div>
              <label className="block text-gray-700 font-medium mb-2">内容</label>
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="投稿内容を入力してください..."
                rows={6}
                className="w-full border border-gray-300 px-4 py-3 rounded-lg shadow-sm focus:ring-4 focus:ring-blue-300"
                required
              />
            </div>
            <div>
              <label className="block text-gray-700 font-medium mb-2">投稿者</label>
              <input
                type="text"
                value={author}
                onChange={(e) => setAuthor(e.target.value)}
                placeholder="例: 山田太郎"
                className="w-full border border-gray-300 px-4 py-3 rounded-lg shadow-sm focus:ring-4 focus:ring-blue-300"
                required
              />
            </div>
            <div className="flex justify-between">
              <button
                type="submit"
                className="bg-blue-500 text-white px-6 py-3 rounded-lg shadow hover:shadow-lg hover:bg-blue-600"
              >
                {id ? "更新する" : "作成する"}
              </button>
              <button
                type="button"
                onClick={() => navigate("/")}
                className="bg-gray-400 text-white px-6 py-3 rounded-lg shadow hover:shadow-lg hover:bg-gray-500"
              >
                投稿一覧に戻る
              </button>
            </div>
          </form>
        </div>
      </div>
    </Layout>
  );
};

export default PostForm;
