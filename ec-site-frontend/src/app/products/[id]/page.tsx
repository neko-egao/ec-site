// 1. Productの型を定義します
type Product = {
  id: number;
  name: string;
  price: number;
};

// 2. IDを元に、APIから特定の商品データを取得する関数です
async function getProductById(id: string) {
  const res = await fetch(`http://localhost:8081/products/${id}`, { cache: 'no-store' });
  if (!res.ok) {
    // もし商品が見つからなければ、Next.jsが自動で404ページを表示します
    throw new Error('商品が見つかりませんでした');
  }
  return res.json();
}

// 3. ページ本体です
// paramsという引数で、URLの動的な部分（今回はid）を受け取れます
export default async function ProductDetailPage({ params }: { params: { id: string } }) {
  const product: Product = await getProductById(params.id);

  return (
    <main className="container mx-auto p-8">
      <div className="border rounded-lg p-8 shadow-lg bg-white max-w-2xl mx-auto">
        <h1 className="text-4xl font-bold mb-4">{product.name}</h1>
        <p className="text-2xl text-gray-800 mb-6">{product.price.toLocaleString()}円</p>
        <p className="text-gray-600">（ここに商品の詳細な説明文が入ります）</p>
        <button className="mt-8 bg-blue-500 text-white font-bold py-2 px-4 rounded hover:bg-blue-700">
          カートに入れる
        </button>
      </div>
    </main>
  );
}