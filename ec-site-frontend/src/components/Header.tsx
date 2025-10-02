'use client';

import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

export default function Header() {
  // ★ Contextから user と setUser を取得します
  const { user, setUser } = useAuth(); 
  const router = useRouter();

  const handleLogout = () => {
    // ★ userとtokenの両方をクリアします
    setUser(null, null); 
    router.push('/');
  };

  return (
    <header className="bg-white shadow-md">
      <nav className="container mx-auto px-6 py-4 flex justify-between items-center">
        <Link href="/" className="text-2xl font-bold text-gray-800">
          ECサイト
        </Link>
        <div>
          {/* ★★★ userオブジェクトが存在するかどうかで表示を切り替えます ★★★ */}
          {user ? (
            // ログインしている場合
            <div className="flex items-center gap-4">
              <span className="text-gray-800">こんにちは、{user.name}さん</span>
              {/* ★ user.role を見て、管理者リンクを表示します */}
              {user.role === 'admin' && (
                <Link href="/admin" className="text-gray-600 font-bold hover:text-blue-500">
                  商品管理
                </Link>
              )}
              <button
                onClick={handleLogout}
                className="bg-red-500 text-white font-bold px-4 py-2 rounded hover:bg-red-700"
              >
                ログアウト
              </button>
            </div>
          ) : (
            // ログインしていない場合
            <>
              <Link href="/login" className="text-gray-600 hover:text-gray-800 px-3 py-2">
                ログイン
              </Link>
              <Link href="/register" className="bg-blue-500 text-white font-bold px-4 py-2 rounded hover:bg-blue-700">
                新規登録
              </Link>
            </>
          )}
        </div>
      </nav>
    </header>
  );
}