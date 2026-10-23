import { useEffect, useState } from 'react';
import { apiFetch } from '../services/http';

type Summary = { total_vendido: number; pedidos: number; ticket_promedio: number; comisiones_devengadas: number };
type TopRef = { nombre: string; email: string; nivel: number; ventas: number; comisiones: number };
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
        setSummary(await summaryResponse.json()); setTop(await topResponse.json()); setProducts(await productsResponse.json());
      }).catch((reason) => setError(reason instanceof Error ? reason.message : 'No se pudieron cargar los reportes'));
  }, []);
  const money = (value: number | undefined) => `$${(value ?? 0).toFixed(2)}`;
  const cards = [
    { label: 'Ventas totales', help: 'Todos los pedidos', value: money(summary?.total_vendido), color: 'text-indigo-600' },
    { label: 'Pedidos', help: 'Compras confirmadas', value: String(summary?.pedidos ?? 0), color: 'text-indigo-600' },
    { label: 'Ticket promedio', help: 'Promedio por pedido', value: money(summary?.ticket_promedio), color: 'text-amber-600' },
    { label: 'Comisiones', help: 'Pendientes de la red', value: money(summary?.comisiones_devengadas), color: 'text-emerald-600' },
  ];
  return <div className="mx-auto max-w-7xl space-y-6">
    <div><p className="text-sm font-semibold uppercase tracking-wider text-indigo-600">Panel administrativo</p><h1 className="mt-1 text-3xl font-extrabold tracking-tight text-slate-900">Resumen del negocio</h1><p className="mt-1 text-slate-500">Consulta tus pedidos y el volumen generado por tu red de referidos.</p></div>
    {error && <p className="rounded-xl border border-red-200 bg-red-50 p-3 text-red-700">{error}</p>}
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">{cards.map((card) => <div key={card.label} className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"><p className="text-xs font-semibold uppercase tracking-wide text-slate-500">{card.label}</p><p className={`mt-2 text-3xl font-extrabold ${card.color}`}>{summary ? card.value : '…'}</p><p className="mt-1 text-xs text-slate-400">{card.help}</p></div>)}</div>
    <div className="rounded-2xl border border-indigo-100 bg-indigo-50/70 p-5"><div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"><div><h2 className="font-bold text-indigo-950">¿Cómo se calculan las ventas de la red?</h2><p className="mt-1 max-w-3xl text-sm leading-6 text-indigo-900/75">Cuando un usuario compra, el pedido cuenta como volumen para sus referentes. El comprador no recibe comisión por su propia compra; las comisiones suben por su árbol.</p></div><div className="grid grid-cols-3 gap-2 text-center text-xs"><div className="rounded-xl bg-white px-3 py-2 shadow-sm"><b className="block text-indigo-700">Nivel 1</b><span className="text-slate-500">10%</span></div><div className="rounded-xl bg-white px-3 py-2 shadow-sm"><b className="block text-indigo-700">Nivel 2</b><span className="text-slate-500">5%</span></div><div className="rounded-xl bg-white px-3 py-2 shadow-sm"><b className="block text-indigo-700">Nivel 3</b><span className="text-slate-500">2%</span></div></div></div></div>
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-2"><div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"><div className="mb-4"><h2 className="text-lg font-bold text-slate-900">Productos más vendidos</h2><p className="text-sm text-slate-500">Productos incluidos en los pedidos</p></div>{products.length === 0 ? <p className="text-slate-500">Aún no hay ventas registradas.</p> : <div className="space-y-3">{products.slice(0, 5).map((product, index) => <div key={product.producto_id} className="flex items-center justify-between border-b border-slate-100 pb-3"><span><b className="mr-2 text-indigo-600">{index + 1}.</b>{product.nombre}<small className="ml-2 text-slate-500">{product.cantidad} unidades</small></span><b>{money(product.total)}</b></div>)}</div>}</div>
      <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"><div className="mb-4"><h2 className="text-lg font-bold text-slate-900">Ventas generadas por la red</h2><p className="text-sm text-slate-500">Volumen atribuido a cada referente y comisión acumulada</p></div>{top.length === 0 ? <p className="text-slate-500">No hay referidos registrados.</p> : <div className="space-y-3">{top.slice(0, 5).map((ref, index) => <div key={ref.email} className="flex items-center justify-between gap-3 border-b border-slate-100 pb-3"><div><div className="flex flex-wrap items-center gap-2"><b className="text-indigo-600">{index + 1}.</b><span className="font-semibold text-slate-800">{ref.nombre || ref.email}</span><span className="rounded-full bg-indigo-50 px-2 py-0.5 text-[11px] font-semibold text-indigo-700">{ref.nivel === 0 ? 'Raíz' : `Nivel ${ref.nivel}`}</span></div><p className="ml-5 text-xs text-slate-500">Comisión generada: {money(ref.comisiones)}</p></div><b>{money(ref.ventas)}</b></div>)}</div>}</div></div>
  </div>;
}
