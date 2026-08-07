/**
 * InspectModal — modal reutilizável estilo "comprovante / inspect"
 * Use em qualquer lista do sistema para inspecionar um item.
 *
 * Exemplo:
 *   import { openInspect } from '../components/InspectModal.js';
 *   import { formatBRL } from '../app.js';
 *
 *   openInspect({
 *     title: 'Comprovante de Venda',
 *     subtitle: 'IDroid Virtual',
 *     status: { label: 'CONCLUÍDA', color: '#27ae60' },
 *     fields: [
 *       { label: 'Venda', value: '#42' },
 *       { label: 'Produto', value: 'iPhone 13' },
 *       { label: 'Total', value: formatBRL(2500), highlight: true },
 *     ],
 *     footer: 'IDroid Virtual · Gestão',
 *   });
 */

let overlay = null;
let contentEl = null;
let actionsEl = null;

function ensureDom() {
  if (overlay) return;

  overlay = document.createElement('div');
  overlay.id = 'inspect-modal';
  overlay.className = 'modal-overlay center';
  overlay.innerHTML = `
    <div class="modal inspect-modal-box" style="max-width:380px;width:100%">
      <div class="modal-handle"></div>
      <div id="inspect-content"></div>
      <div class="modal-actions" id="inspect-actions">
        <button type="button" class="btn btn-ghost" id="inspect-close">Fechar</button>
        <button type="button" class="btn btn-primary" id="inspect-print">📄 Tirar print</button>
      </div>
    </div>
  `;
  document.body.appendChild(overlay);

  contentEl = overlay.querySelector('#inspect-content');
  actionsEl = overlay.querySelector('#inspect-actions');

  overlay.querySelector('#inspect-close').onclick = closeInspect;

  overlay.querySelector('#inspect-print').onclick = () => {
    actionsEl.style.display = 'none';
    // Pequeno delay para o browser aplicar o hide antes do print
    setTimeout(() => {
      window.print();
      setTimeout(() => {
        actionsEl.style.display = '';
      }, 400);
    }, 50);
  };

  overlay.addEventListener('click', (e) => {
    if (e.target === overlay) closeInspect();
  });

  // Esc fecha
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && overlay.classList.contains('open')) {
      closeInspect();
    }
  });
}

/**
 * Abre o modal de inspeção.
 * @param {Object} options
 * @param {string}  [options.title='Detalhes']
 * @param {string}  [options.subtitle='']
 * @param {{label:string, color?:string}} [options.status]
 * @param {Array<{label:string, value:string|number, highlight?:boolean}>} options.fields
 * @param {string}  [options.footer='']
 * @param {boolean} [options.showPrint=true]
 */
export function openInspect(options = {}) {
  ensureDom();

  const {
    title = 'Detalhes',
    subtitle = '',
    status = null,
    fields = [],
    footer = '',
    showPrint = true,
  } = options;

  const statusHtml = status
    ? `<div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px">
         <span style="font-size:0.85rem;color:#666">Status</span>
         <span style="font-weight:600;color:${status.color || 'var(--text)'}">${status.label}</span>
       </div>`
    : '';

  const fieldsHtml = fields
    .map((f) => {
      const isHighlight = f.highlight;
      return `
        <div style="display:flex;justify-content:space-between;gap:12px;padding:6px 0;${isHighlight ? 'font-size:1.05rem;border-top:1px dashed #eee;margin-top:6px;padding-top:10px' : ''}">
          <span style="color:${isHighlight ? 'var(--text)' : '#666'};font-weight:${isHighlight ? 700 : 400}">${f.label}</span>
          <span style="text-align:right;max-width:62%;word-break:break-word;font-weight:${isHighlight ? 700 : 500};color:${isHighlight ? '#f06522' : 'inherit'}">${f.value ?? '—'}</span>
        </div>`;
    })
    .join('');

  contentEl.innerHTML = `
    <div class="inspect-printable">
      <div style="text-align:center;border-bottom:1px dashed #ccc;padding-bottom:12px;margin-bottom:12px">
        <div style="font-size:1.1rem;font-weight:700;color:#f06522">${escapeHtml(title)}</div>
        ${subtitle ? `<div style="font-size:0.8rem;color:#666;margin-top:2px">${escapeHtml(subtitle)}</div>` : ''}
      </div>

      ${statusHtml}

      <div style="font-size:0.85rem;line-height:1.5">
        ${fieldsHtml}
      </div>

      ${
        footer
          ? `<div style="text-align:center;margin-top:16px;font-size:0.75rem;color:#999">${escapeHtml(footer)}</div>`
          : ''
      }
    </div>
  `;

  // Mostra/esconde botão de print
  const printBtn = overlay.querySelector('#inspect-print');
  if (printBtn) printBtn.style.display = showPrint ? '' : 'none';

  overlay.classList.add('open');
}

export function closeInspect() {
  if (!overlay) return;
  overlay.classList.remove('open');
}

function escapeHtml(str) {
  return String(str ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

/** Atalho: inspecionar um objeto genérico (mostra todas as chaves) */
export function openInspectObject(obj, title = 'Detalhes') {
  if (!obj || typeof obj !== 'object') return;

  const fields = Object.entries(obj)
    .filter(([k]) => !k.startsWith('_'))
    .map(([label, value]) => ({
      label: formatKey(label),
      value: formatValue(value),
    }));

  openInspect({
    title,
    subtitle: 'IDroid Virtual',
    fields,
    footer: 'IDroid Virtual · Gestão',
  });
}

function formatKey(key) {
  return String(key)
    .replace(/_/g, ' ')
    .replace(/\b\w/g, (c) => c.toUpperCase());
}

function formatValue(val) {
  if (val == null) return '—';
  if (typeof val === 'boolean') return val ? 'Sim' : 'Não';
  if (typeof val === 'object') return JSON.stringify(val);
  return String(val);
}
