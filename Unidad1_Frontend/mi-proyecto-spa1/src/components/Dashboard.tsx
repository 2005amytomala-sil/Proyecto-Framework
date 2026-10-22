import { useEffect, useState } from 'react';
import { apiFetch } from '../services/http';

type Summary = { total_vendido: number; pedidos: number; ticket_promedio: number; comisiones_devengadas: number };
type TopRef = { nombre: string; email: string; ventas: number; comisiones: number };
type TopProduct = { producto_id: number; nombre: string; cantidad: number; total: number };

export default function Dashboard() {
  const [summary, setSummary] = useState<Summary | null>(null);
  const [top, setTop] = useState<TopRef[]>([]);
  const [products, setProducts] = useState<TopProduct[]>([]);
  const [error, setError] = useState('');
  useEffect(() => {
    Promise.all([apiFetch('/reportes/resumen'), apiFetch('/reportes/top-referidos'), apiFetch('/reportes/top-productos')])
      .then(async ([summaryResponse, topResponse, productsResponse]) => {
        if (!summaryResponse.ok || !topResponse.ok || !productsResponse.ok) throw new Error('No se pudieron cargar los reportes');
        setSummary(await summaryResponse.json());
        setTop(await topResponse.json());
        setProducts(await productsResponse.json());
      })
      .catch((reason) => setError(reason instanceof Error ? reason.message : 'No se pudieron cargar los reportes'));
  }, []);
  const money = (value: number | undefined) => `$${(value ?? 0).toFixed(2)}`;
  const cards = [
    ['Ventas totales', money(summary?.total_vendido), 'text-indigo-600'],
    ['Pedidos', String(summary?.pedidos ?? 0), 'text-indigo-600'],
    ['Ticket promedio', money(summary?.ticket_promedio), 'text-amber-600'],
    ['Comisiones', money(summary?.comisiones_devengadas), 'text-green-600'],
  ];
  return <div><h1 className="mb-6 text-2xl font-bold text-slate-800">Resumen del negocio</h1>{error && <p className="mb-4 rounded bg-red-50 p-3 text-red-700">{error}</p>}<div className="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-4">{cards.map(([label, value, color]) => <div key={label} className="rounded-lg border border-slate-100 bg-white p-6 shadow-sm"><p className="text-sm font-semibold uppercase text-slate-500">{label}</p><p className={`mt-2 text-3xl font-bold ${color}`}>{summary ? value : '…'}</p></div>)}</div><div className="grid grid-cols-1 gap-6 lg:grid-cols-2"><div className="rounded-lg border border-slate-100 bg-white p-6 shadow-sm"><h2 className="mb-4 text-lg font-bold text-slate-800">Productos más vendidos</h2>{products.length === 0 ? <p className="text-slate-500">Aún no hay ventas registradas.</p> : <div className="space-y-3">{products.slice(0, 5).map((product, index) => <div key={product.producto_id} className="flex items-center justify-between border-b border-slate-100 pb-3"><span><b className="mr-2">{index + 1}.</b>{product.nombre}<small className="ml-2 text-slate-500">{product.cantidad} unidades</small></span><b>{money(product.total)}</b></div>)}</div>}</div><div className="rounded-lg border border-slate-100 bg-white p-6 shadow-sm"><h2 className="mb-4 text-lg font-bold text-slate-800">Top de la red de ventas</h2>{top.length === 0 ? <p className="text-slate-500">No hay referidos registrados.</p> : <div className="space-y-3">{top.slice(0, 5).map((ref, index) => <div key={ref.email} className="flex items-center justify-between border-b border-slate-100 pb-3"><span><b className="mr-2">{index + 1}.</b>{ref.nombre || ref.email}</span><b>{money(ref.ventas)}</b></div>)}</div>}</div></div></div>;
}
