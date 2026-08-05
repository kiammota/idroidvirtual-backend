import { requireAuth, goBack, showToast, formatPhone } from '../app.js';
import { bindMasks } from '../masks.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { ClientService } from '../../client/services/clientService.js';
import { ServiceOrderService } from '../../client/services/OrdensServicoService.js'
import { API_BASE_URL } from '../../client/env.js';


window.NotificationService = NotificationService;

if (!requireAuth()) throw new Error('redirect');
document.getElementById('btn-back').onclick = () => goBack();

bindMasks({ 'os-telefone': 'phone' });

document.querySelectorAll('#os-tabs .tab').forEach(tab => {
  tab.onclick = () => {
    document.querySelectorAll('#os-tabs .tab').forEach(t => t.classList.remove('active'));
    tab.classList.add('active');
    document.getElementById('tab-lista').classList.toggle('hidden', tab.dataset.tab !== 'lista');
    document.getElementById('tab-nova').classList.toggle('hidden', tab.dataset.tab !== 'nova');
  };
});

const canvas = document.getElementById('sig-canvas');
const ctx = canvas.getContext('2d');
let drawing = false;
let hasSignature = false;

function resizeCanvas() {
  const rect = canvas.parentElement.getBoundingClientRect();
  canvas.width = rect.width;
  canvas.height = 160;
  ctx.strokeStyle = '#1a1a1a';
  ctx.lineWidth = 2.5;
  ctx.lineCap = 'round';
  ctx.lineJoin = 'round';
}
resizeCanvas();
window.addEventListener('resize', resizeCanvas);

function getPos(e) {
  const rect = canvas.getBoundingClientRect();
  const t = e.touches ? e.touches[0] : e;
  return { x: t.clientX - rect.left, y: t.clientY - rect.top };
}
function start(e) { e.preventDefault(); drawing = true; const p = getPos(e); ctx.beginPath(); ctx.moveTo(p.x, p.y); }
function move(e) {
  if (!drawing) return;
  e.preventDefault();
  const p = getPos(e);
  ctx.lineTo(p.x, p.y);
  ctx.stroke();
  hasSignature = true;
}
function end() { drawing = false; }

canvas.addEventListener('mousedown', start);
canvas.addEventListener('mousemove', move);
canvas.addEventListener('mouseup', end);
canvas.addEventListener('mouseleave', end);
canvas.addEventListener('touchstart', start, { passive: false });
canvas.addEventListener('touchmove', move, { passive: false });
canvas.addEventListener('touchend', end);

document.getElementById('sig-clear').onclick = () => {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  hasSignature = false;
  document.getElementById('os-assinatura').value = '';
  document.getElementById('sig-status').textContent = 'Aguardando assinatura…';
};
document.getElementById('sig-confirm').onclick = () => {
  if (!hasSignature) {
    showToast('Assine no campo antes de confirmar', 'error');
    return;
  }
  document.getElementById('os-assinatura').value = canvas.toDataURL('image/png');
  document.getElementById('sig-status').textContent = '✓ Assinatura confirmada';
  showToast('Assinatura salva', 'success');
};

document.getElementById('os-telefone').addEventListener('blur', () => {
  const tel = document.getElementById('os-telefone').value.replace(/\D/g, '');
  // CHAMADA DA API AQUI — verificar OS existentes por telefone
  if (tel.length >= 10) {
    document.getElementById('os-alerta').classList.toggle('hidden', !tel.endsWith('00'));
  }
});

document.addEventListener('DOMContentLoaded', () => {
  const formOS = document.getElementById('form-os');
  
  if (!formOS) {
    console.error('Formulário #form-os não foi encontrado na página.');
    return;
  }

  formOS.onsubmit = async (e) => {
    e.preventDefault();

    // 1. Coleta os dados
    const clienteNome = document.getElementById('os-cliente')?.value.trim() || '';
    const telefone = document.getElementById('os-telefone')?.value.replace(/\D/g, '') || '';
    const defeito = document.getElementById('os-defeito')?.value.trim() || '';
    const servicos = document.getElementById('os-servico')?.value.trim() || '';
    const valor = parseFloat(document.getElementById('os-valor')?.value) || 0;
    const assinatura = document.getElementById('os-assinatura')?.value || null;

    const descricaoComposta = [
      `Cliente: ${clienteNome}`,
      `Telefone: ${telefone}`,
      `Defeito: ${defeito}`,
      `Serviços: ${servicos}`,
      `Valor: ${valor}`
    ].filter(Boolean).join(' | ');

    try {
      const serviceOrderService = new ServiceOrderService();
      
      const response = await serviceOrderService.create({
        cliente_id: 1,
        descricao: descricaoComposta,
        status: 'ABERTA'
      });

      if (response && response.success) {
        showToast('Ordem de serviço salva com sucesso', 'success');
        
        formOS.reset();
        
        if (typeof ctx !== 'undefined' && typeof canvas !== 'undefined') {
          ctx.clearRect(0, 0, canvas.width, canvas.height);
          if (typeof hasSignature !== 'undefined') hasSignature = false;
          const statusElem = document.getElementById('sig-status');
          if (statusElem) statusElem.textContent = 'Aguardando assinatura…';
        }

        const abaLista = document.querySelector('[data-tab="lista"]');
        if (abaLista) abaLista.click();

        if (typeof carregarOrdensServico === 'function') {
          await carregarOrdensServico();
        }
      } else {
        showToast(response?.message || 'Erro ao salvar ordem de serviço', 'error');
      }
    } catch (err) {
      console.error('Erro ao enviar form:', err);
      showToast(err.message || 'Erro de conexão ao salvar OS', 'error');
    }
  };
});

function renderList(items) {
  const el = document.getElementById('os-list');
  if (!items || !items.length) {
    el.innerHTML = `<div class="empty-state"><div class="icon">🛠️</div><p>Nenhuma ordem de serviço ainda.</p></div>`;
    return;
  }
  el.innerHTML = items.map(os => `
    <div class="card list-item" style="flex-direction:column;align-items:stretch;gap:6px">
      <div class="flex justify-between items-center">
        <span class="fw-bold">OS #${os.numero_ordem || os.id}</span>
        <span class="badge badge-aberta">${os.status || 'ABERTA'}</span>
      </div>
      <div class="title">${os.cliente_nome || 'Cliente'}</div>
      <div class="sub">${formatPhone(os.telefone)} · ${os.defeito || ''}</div>
      ${os.valor ? `<div class="text-primary fw-bold">${Number(os.valor).toLocaleString('pt-BR',{style:'currency',currency:'BRL'})}</div>` : ''}
    </div>
  `).join('');
}

// CHAMADA DA API AQUI
renderList([]);
