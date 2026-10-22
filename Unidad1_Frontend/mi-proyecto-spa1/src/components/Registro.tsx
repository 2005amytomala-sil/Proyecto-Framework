import { useState, type FormEvent } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { apiFetch } from '../services/http';

export default function Registro() {
  const [params] = useSearchParams(); const navigate = useNavigate();
  const [form, setForm] = useState({ nombre: '', email: '', password: '', ref: params.get('ref') ?? '' });
  const [error, setError] = useState(''); const [saving, setSaving] = useState(false);
  const submit = async (event: FormEvent) => { event.preventDefault(); setSaving(true); setError(''); try { const response = await apiFetch('/register', { method:'POST', body: JSON.stringify(form) }); const data = await response.json(); if (!response.ok) throw new Error(data.message ?? data.error ?? 'No se pudo registrar'); navigate('/login'); } catch (e) { setError(e instanceof Error ? e.message : 'No se pudo registrar'); } finally { setSaving(false); } };
  return <main className="min-h-screen grid place-items-center bg-slate-100 p-4"><form onSubmit={submit} className="w-full max-w-md space-y-4 rounded-xl bg-white p-6 shadow"><h1 className="text-2xl font-bold">Crear cuenta</h1>{error && <p className="rounded bg-red-50 p-3 text-red-700">{error}</p>}{(['nombre','email','password','ref'] as const).map(field => <input key={field} required={field !== 'ref'} type={field === 'password' ? 'password' : field === 'email' ? 'email' : 'text'} placeholder={field === 'ref' ? 'Código de referido (opcional)' : field[0].toUpperCase()+field.slice(1)} value={form[field]} onChange={e=>setForm({...form,[field]:e.target.value})} className="w-full rounded border p-3" />)}<button disabled={saving} className="w-full rounded bg-indigo-600 p-3 font-semibold text-white">{saving ? 'Registrando…' : 'Registrarme'}</button><Link className="block text-center text-indigo-600" to="/login">Ya tengo una cuenta</Link></form></main>;
}
