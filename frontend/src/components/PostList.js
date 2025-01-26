import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Layout from "./Layout/Layout";

const PostList = () => {
  const [posts, setPosts] = useState([]);

  useEffect(() => {
    fetch(`${process.env.REACT_APP_URL_DOMAIN}/api/posts`)
      .then((response) => response.json())
      .then((data) => setPosts(data))
      .catch((error) => console.error("Error fetching posts:", error));
  }, []);

  return (
    <Layout>
      <div className="bg-gray-100 min-h-screen">
        <div className="container mx-auto px-4 py-10">
          <h1 className="text-4xl font-bold text-gray-800 text-center mb-8">
            投稿一覧
          </h1>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {posts.length > 0 ? (
              posts.map((post) => (
                <div
                  key={post.ID}
                  className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
                >
                  <h2 className="text-xl font-semibold text-blue-600 mb-4">
                    <Link to={`/posts/${post.ID}`} className="hover:underline">
                      {post.Title}
                    </Link>
                  </h2>
                  <p className="text-gray-700 mb-4">
                    {post.Content.slice(0, 100)}...
                  </p>
                  <div className="flex justify-between">
                    <Link
                      to={`/posts/${post.ID}`}
                      className="text-blue-500 hover:underline"
                    >
                      詳細を見る →
                    </Link>
                    <Link
                      to={`/edit/${post.ID}`}
                      className="text-gray-600 bg-gray-200 px-3 py-1 rounded-md hover:bg-gray-300"
                    >
                      編集
                    </Link>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-gray-600 text-center">投稿がありません。</p>
            )}
          </div>
        </div>
        <footer className="bg-gray-800 text-white py-6 text-center mt-10">
          <p className="text-sm">© 2024 ポートフォリオブログ. All Rights Reserved.</p>
        </footer>
      </div>
    </Layout>
  );
};

export default PostList;
