import Link from 'next/link'; // ★ next/linkをインポートします

// ... (Product型定義とgetProducts関数は同じです) ...
type Product = {
  id: number;
  name: string;
  price: number;
};
async function getProducts() {
  const res = await fetch('http://localhost:8081/products', { cache: 'no-store' });
  if (!res.ok) {
    throw new Error('商品データの取得に失敗しました');
  }
  return res.json();
}

export default async function HomePage() {
  const products: Product[] = await getProducts();

  return (
    <main className="container mx-auto p-8">
      <h1 className="text-4xl font-bold mb-10 text-center">商品一覧</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {products.map((product) => (
          // ★★★ divをLinkコンポーネントで囲みます ★★★
          <Link href={`/products/${product.id}`} key={product.id} className="block hover:opacity-80">
            <div className="border rounded-lg p-6 shadow-lg bg-white h-full cursor-pointer">
              <h2 className="text-2xl font-semibold mb-2">{product.name}</h2>
              <p className="text-lg text-gray-700">{product.price.toLocaleString()}円</p>
            </div>
          </Link>
        ))}
      </div>
    </main>
  );
}