import React from "react";
import { Link } from "react-router-dom";

const Layout = ({ children }) => {
  return (
    <div>
      <header className="bg-blue-500 text-white p-4">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold">ポートフォリオブログ</h1>
          <nav className="flex space-x-4">
            <Link to="/" className="hover:underline">
              投稿一覧
            </Link>
            <Link to="/create" className="hover:underline">
              新規投稿
            </Link>
          </nav>
        </div>
      </header>
      <main>{children}</main>
    </div>
  );
};

export default Layout;
