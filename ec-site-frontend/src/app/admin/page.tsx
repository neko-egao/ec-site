'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { jwtDecode } from 'jwt-decode';

// 商品の型定義
type Product = {
  id: number;
  name: string;
  price: number;
};

// JWTペイロードの型定義
type DecodedToken = {
  user_id: number;
  role: string;
  exp: number;
};

export default function AdminPage() {
  const { token } = useAuth();
  const router = useRouter();
  const [products, setProducts] = useState<Product[]>([]);
  const [isAdmin, setIsAdmin] = useState(false);

  // ★★★ 商品追加フォーム用の状態を追加 ★★★
  const [newProductName, setNewProductName] = useState('');
  const [newProductPrice, setNewProductPrice] = useState('');
  const [addError, setAddError] = useState('');

  // ★★★ 編集モーダル用の状態を追加 ★★★
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  // ★★★★★★★★★★★★★★★★★★★★★

  useEffect(() => {
    if (token) {
      try {
        const decoded = jwtDecode<DecodedToken>(token);
        if (decoded.role === 'admin') {
          setIsAdmin(true);
          fetchProducts();
        } else {
          router.push('/');
        }
      } catch (error) {
        router.push('/login');
      }
    } else {
      router.push('/login');
    }
  }, [token, router]);

  const fetchProducts = async () => {
    const res = await fetch('http://localhost:8081/products');
    const data = await res.json();
    setProducts(data);
  };
  
  const handleDelete = async (id: number) => {
    if (confirm('本当にこの商品を削除しますか？') && token) {
      try {
        const res = await fetch(`http://localhost:8081/products/${id}`, {
          method: 'DELETE',
          headers: { 'Authorization': `Bearer ${token}` },
        });

        if (!res.ok) {
          const errorData = await res.json();
          throw new Error(errorData.error || '商品の削除に失敗しました。');
        }
        fetchProducts(); // 成功した場合のみリストを再取得
      } catch (error) {
        console.error('削除エラー:', error);
        alert(error instanceof Error ? error.message : '予期せぬエラーが発生しました。');
      }
    }
  };

  // ★★★ 商品追加のロジックを追加 ★★★
  const handleCreate = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setAddError('');
    if (!token) return;

    try {
      const res = await fetch(`http://localhost:8081/products`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          name: newProductName,
          price: parseInt(newProductPrice, 10),
        }),
      });

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || '商品の追加に失敗しました。');
      }

      fetchProducts(); // 成功したらリストを再取得
      setNewProductName(''); // フォームをリセット
      setNewProductPrice('');
    } catch (error) {
      setAddError(error instanceof Error ? error.message : '予期せぬエラーが発生しました。');
    }
  };
  
  // ★★★ 商品更新のロジックを追加 ★★★
  const handleUpdate = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!editingProduct || !token) return;

    try {
      const res = await fetch(`http://localhost:8081/products/${editingProduct.id}`, {
        method: 'PUT',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}` 
        },
        body: JSON.stringify({
          name: editingProduct.name,
          price: editingProduct.price,
        }),
      });

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || '商品の更新に失敗しました。');
      }
      
      fetchProducts(); // 成功した場合のみリストを再取得
      setIsModalOpen(false); // モーダルを閉じる
      setEditingProduct(null);
    } catch (error) {
      console.error('更新エラー:', error);
      alert(error instanceof Error ? error.message : '予期せぬエラーが発生しました。');
    }
  };

  const openEditModal = (product: Product) => {
    setEditingProduct({ ...product });
    setIsModalOpen(true);
  };
  // ★★★★★★★★★★★★★★★★★★★★★

  if (!isAdmin) {
    return null; 
  }

  return (
    <main className="container mx-auto p-8">
      <h1 className="text-4xl font-bold mb-10">管理者用 商品管理</h1>

      {/* ★★★ 商品追加フォームを追加 ★★★ */}
      <div className="mb-10">
        <h2 className="text-2xl font-bold mb-4">商品を追加する</h2>
        <form onSubmit={handleCreate} className="bg-white shadow-md rounded px-8 pt-6 pb-8">
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="new-product-name">
              商品名
            </label>
            <input
              id="new-product-name"
              type="text"
              value={newProductName}
              onChange={(e) => setNewProductName(e.target.value)}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              required
            />
          </div>
          <div className="mb-6">
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="new-product-price">
              価格
            </label>
            <input
              id="new-product-price"
              type="number"
              value={newProductPrice}
              onChange={(e) => setNewProductPrice(e.target.value)}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              required
            />
          </div>
          {addError && <p className="text-red-500 text-xs italic mb-4">{addError}</p>}
          <div className="flex items-center">
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
              type="submit"
            >
              商品を追加
            </button>
          </div>
        </form>
      </div>

      <div className="bg-white shadow-md rounded my-6">
        <table className="min-w-full table-auto">
          {/* ... theadは変更なし ... */}
          <tbody className="text-gray-600 text-sm font-light">
            {products.map((product) => (
              <tr key={product.id} className="border-b border-gray-200 hover:bg-gray-100">
                <td className="py-3 px-6 text-left whitespace-nowrap">{product.name}</td>
                <td className="py-3 px-6 text-left">{product.price.toLocaleString()}円</td>
                <td className="py-3 px-6 text-center">
                  <div className="flex item-center justify-center">
                    {/* ★ 編集ボタンにonClickイベントを追加 */}
                    <button onClick={() => openEditModal(product)} className="w-4 mr-2 transform hover:text-purple-500 hover:scale-110">
                      ✏️
                    </button>
                    <button onClick={() => handleDelete(product.id)} className="w-4 mr-2 transform hover:text-red-500 hover:scale-110">
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* ★★★ 編集用モーダルを追加 ★★★ */}
      {isModalOpen && editingProduct && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
          <div className="bg-white p-8 rounded-lg shadow-2xl w-full max-w-md">
            <h2 className="text-2xl font-bold mb-4">商品を編集</h2>
            <form onSubmit={handleUpdate}>
              <div className="mb-4">
                <label className="block text-gray-700 text-sm font-bold mb-2">商品名</label>
                <input
                  type="text"
                  value={editingProduct.name}
                  onChange={(e) => setEditingProduct({ ...editingProduct, name: e.target.value })}
                  className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                />
              </div>
              <div className="mb-6">
                <label className="block text-gray-700 text-sm font-bold mb-2">価格</label>
                <input
                  type="number"
                  value={editingProduct.price}
                  onChange={(e) => setEditingProduct({ ...editingProduct, price: parseInt(e.target.value, 10) || 0 })}
                  className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700"
                />
              </div>
              <div className="flex items-center justify-end gap-4">
                <button type="button" onClick={() => setIsModalOpen(false)} className="text-gray-600">
                  キャンセル
                </button>
                <button type="submit" className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                  更新する
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </main>
  );
}