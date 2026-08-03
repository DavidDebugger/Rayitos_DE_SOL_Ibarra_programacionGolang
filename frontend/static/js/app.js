/* ============================================================
   Rayitos de Sol - Utilidades compartidas del frontend
   Depende de: Tailwind CSS (CDN) y Chart.js (CDN, solo dashboard)
   ============================================================ */

const API_BASE = '/api';

/* ---------- Fetch wrapper con manejo de errores ---------- */
async function api(path, options = {}) {
    const opts = { headers: { 'Content-Type': 'application/json' }, ...options };
    try {
        const res = await fetch(API_BASE + path, opts);
        if (!res.ok) {
            let msg = 'Error ' + res.status;
            try {
                const err = await res.json();
                if (err.error) msg = err.error;
            } catch (_) {}
            throw new Error(msg);
        }
        if (res.status === 204) return null;
        return await res.json();
    } catch (err) {
        if (err.name !== 'Error' || !err.message.startsWith('Error ')) {
            // error de red / JSON
            toast(err.message, 'error');
            throw err;
        }
        throw err;
    }
}

/* ---------- Toasts ---------- */
function toast(message, type = 'success') {
    let container = document.getElementById('toast-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-container';
        document.body.appendChild(container);
    }
    const icons = { success: '✅', error: '❌', info: 'ℹ️' };
    const el = document.createElement('div');
    el.className = 'toast toast-' + type;
    el.innerHTML = `<span>${icons[type] || ''}</span><span>${escapeHtml(message)}</span>`;
    container.appendChild(el);
    setTimeout(() => {
        el.classList.add('toast-out');
        setTimeout(() => el.remove(), 260);
    }, 3500);
}

/* ---------- Modales ---------- */
function abrirModal(id) {
    const modal = document.getElementById(id);
    if (modal) modal.classList.remove('hidden');
}
function cerrarModal(id) {
    const modal = document.getElementById(id);
    if (modal) modal.classList.add('hidden');
}
function cerrarModalConClick(event, id) {
    if (event.target === event.currentTarget) cerrarModal(id);
}

/* ---------- Confirmación ---------- */
function confirmar(mensaje) {
    return new Promise(resolve => {
        const modal = document.createElement('div');
        modal.className = 'fixed inset-0 z-[9998] modal-backdrop flex items-center justify-center p-4';
        modal.innerHTML = `
            <div class="bg-white rounded-xl p-6 w-full max-w-sm shadow-2xl modal-content">
                <div class="flex items-center gap-3 mb-4">
                    <span class="text-3xl">⚠️</span>
                    <h3 class="font-bold text-slate-800 text-lg">Confirmar</h3>
                </div>
                <p class="text-sm text-slate-600 mb-6">${escapeHtml(mensaje)}</p>
                <div class="flex justify-end gap-2">
                    <button class="btn btn-ghost" data-act="no">Cancelar</button>
                    <button class="btn btn-danger" data-act="si">Confirmar</button>
                </div>
            </div>`;
        document.body.appendChild(modal);
        modal.addEventListener('click', e => {
            const act = e.target.closest('[data-act]');
            if (!act) return;
            document.body.removeChild(modal);
            resolve(act.dataset.act === 'si');
        });
    });
}

/* ---------- Badges por estado ---------- */
function badge(estado) {
    const colores = {
        activo: 'badge-green', inactivo: 'badge-red', suspendido: 'badge-yellow',
        programada: 'badge-blue', confirmada: 'badge-green', cancelada: 'badge-red',
        completada: 'badge-green', reprogramada: 'badge-orange',
        pendiente: 'badge-yellow', pagada: 'badge-green', anulada: 'badge-red',
        entrada: 'badge-green', salida: 'badge-red',
        ingreso: 'badge-green', gasto: 'badge-red',
        ok: 'badge-green', bajo: 'badge-red',
    };
    return `<span class="badge ${colores[estado] || 'badge-gray'}">${escapeHtml(estado)}</span>`;
}

/* ---------- Utilidades ---------- */
function escapeHtml(str) {
    if (str === null || str === undefined) return '';
    return String(str)
        .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;').replace(/'/g, '&#39;');
}
function dinero(n) {
    const num = parseFloat(n);
    if (isNaN(num)) return '$0.00';
    return '$' + num.toLocaleString('es-EC', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}
function fechaCorta(f) {
    if (!f) return '';
    return String(f).slice(0, 10);
}
function horaCorta(h) {
    if (!h) return '';
    const s = String(h);
    return s.length > 10 ? s.slice(11, 16) : s.slice(0, 5);
}
function spinnerHtml() {
    return '<tr><td colspan="99"><div class="spinner"></div></td></tr>';
}
function vacioHtml(mensaje = 'No hay registros') {
    return `<tr><td colspan="99"><div class="empty-state">${escapeHtml(mensaje)}</div></td></tr>`;
}

/* ---------- Llenar selects desde API ---------- */
async function llenarSelect(path, selectId, texto = (item) => `#${item.id}`, textoVacio = '— Seleccionar —') {
    try {
        const datos = await api(path);
        const sel = document.getElementById(selectId);
        if (!sel) return;
        let html = `<option value="">${textoVacio}</option>`;
        html += datos.map(item => `<option value="${item.id}">${escapeHtml(texto(item))}</option>`).join('');
        sel.innerHTML = html;
    } catch (_) {}
}

/* ---------- Búsqueda en tabla ---------- */
function initBusqueda(inputId, filasSelector) {
    const input = document.getElementById(inputId);
    if (!input) return;
    input.addEventListener('input', () => {
        const q = input.value.toLowerCase();
        document.querySelectorAll(filasSelector).forEach(tr => {
            tr.style.display = tr.textContent.toLowerCase().includes(q) ? '' : 'none';
        });
    });
}

/* ---------- Sidebar (carga el menú compartido) ---------- */
async function cargarSidebar() {
    const cont = document.getElementById('sidebar');
    if (!cont) return;
    try {
        const res = await fetch('sidebar.html');
        if (res.ok) {
            cont.innerHTML = await res.text();
            // marcar enlace activo según la ruta actual
            const ruta = window.location.pathname;
            cont.querySelectorAll('a').forEach(a => {
                if (a.getAttribute('href') === ruta) a.classList.add('active');
            });
        }
    } catch (_) {
        cont.innerHTML = '<div class="text-slate-400 text-sm p-4">Menú no disponible</div>';
    }
}
