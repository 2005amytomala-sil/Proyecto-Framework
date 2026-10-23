import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth, type Rol } from "../context/AuthContext";
import { apiFetch } from "../services/http";

const SHOW_DEMO = import.meta.env.VITE_SHOW_DEMO_CREDENTIALS === "true";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoading(true);
    apiFetch("/login", { method: "POST", body: JSON.stringify({ email, password }) })
      .then((response) => {
        if (!response.ok) throw new Error("Credenciales incorrectas");
        return response.json();
      })
      .then((data) => {
        const rol: Rol = data.rol === "admin" ? "admin" : "cliente";
        login({ email: data.email, rol, token: data.token });
        navigate(rol === "admin" ? "/" : "/tienda");
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-slate-950 px-4 py-10">
      <div className="absolute -left-24 -top-24 h-72 w-72 rounded-full bg-indigo-600/20 blur-3xl" />
      <div className="absolute -bottom-32 -right-20 h-80 w-80 rounded-full bg-violet-500/15 blur-3xl" />
      <div className="relative w-full max-w-md rounded-2xl border border-white/10 bg-white p-7 shadow-2xl shadow-indigo-950/30 sm:p-9">
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-indigo-600 text-white shadow-lg shadow-indigo-600/30">
            <svg className="h-8 w-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.8" d="M4 7.5 12 4l8 3.5v9L12 20l-8-3.5v-9Z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.8" d="M4 7.5 12 11l8-3.5M12 11v9" /></svg>
          </div>
          <h2 className="text-3xl font-extrabold tracking-tight text-slate-900">MultiCatálogo</h2>
          <p className="mt-2 text-slate-500">Ingresa a tu cuenta para continuar</p>
        </div>

        {error && <div className="mb-6 rounded-xl border border-red-200 bg-red-50 p-3 text-center text-sm text-red-600" role="alert">{error}</div>}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="mb-2 block text-sm font-semibold text-slate-700" htmlFor="email">Correo electrónico</label>
            <div className="relative">
              <svg className="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.8" d="M4 6h16v12H4z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.8" d="m4 7 8 6 8-6" /></svg>
              <input id="email" type="email" autoComplete="email" value={email} onChange={(event) => setEmail(event.target.value)} className="w-full rounded-xl border border-slate-300 bg-slate-50 py-3 pl-11 pr-4 outline-none transition focus:border-indigo-600 focus:bg-white focus:ring-2 focus:ring-indigo-600/20" placeholder="admin@upse.edu.ec" required />
            </div>
          </div>

          <div>
            <label className="mb-2 block text-sm font-semibold text-slate-700" htmlFor="password">Contraseña</label>
            <div className="relative">
              <svg className="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><rect x="5" y="10" width="14" height="10" rx="2" strokeWidth="1.8" /><path strokeLinecap="round" strokeWidth="1.8" d="M8 10V7a4 4 0 0 1 8 0v3" /></svg>
              <input id="password" type={showPassword ? "text" : "password"} autoComplete="current-password" value={password} onChange={(event) => setPassword(event.target.value)} className="w-full rounded-xl border border-slate-300 bg-slate-50 py-3 pl-11 pr-12 outline-none transition focus:border-indigo-600 focus:bg-white focus:ring-2 focus:ring-indigo-600/20" placeholder="••••••" required />
              <button type="button" onClick={() => setShowPassword((visible) => !visible)} className="absolute right-3 top-1/2 -translate-y-1/2 rounded p-1 text-slate-400 hover:text-indigo-600" aria-label={showPassword ? "Ocultar contraseña" : "Mostrar contraseña"}>{showPassword ? "◉" : "◌"}</button>
            </div>
          </div>

          <button type="submit" disabled={loading} className="w-full rounded-xl bg-indigo-600 py-3.5 font-bold text-white shadow-lg shadow-indigo-600/25 transition hover:bg-indigo-700 hover:shadow-indigo-600/40 disabled:cursor-not-allowed disabled:opacity-50">
            {loading ? "Validando..." : "Iniciar sesión"}
          </button>
        </form>

        {SHOW_DEMO && <div className="mt-6 rounded-xl border border-slate-200 bg-slate-50 p-4 text-xs text-slate-600"><p className="font-semibold text-slate-700">Cuentas de prueba</p><p className="mt-2">👑 Admin: <span className="font-mono">admin@upse.edu.ec / 123456</span></p><p className="mt-1">🛍️ Cliente: <span className="font-mono">cliente@upse.edu.ec / 123456</span></p></div>}
      </div>
    </div>
  );
};

export default Login;
