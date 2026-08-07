import { requireAuth, goBack, showToast, formatPhone } from '../app.js';
import { bindMasks } from '../masks.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { ClientService } from '../../client/services/clientService.js';
import { ServiceOrderService } from '../../client/services/OrdensServicoService.js';
import { OrdemServico } from '../../client/models/OrdemServicoModel.js';

window.NotificationService = NotificationService;

if (!requireAuth()) throw new Error('redirect');

document.getElementById('btn-back').onclick = () => goBack();
bindMasks({ 'os-telefone': 'phone' });

const clientService = new ClientService();
const serviceOrderService = new ServiceOrderService();

// ─── Estado ─────────────────────────────────────────────
let modoForm = 'nova';
let ordemAtual = null;
let todasOrdens = [];
let filtroAtual = 'andamento';
let clienteAtual = null; // ← Armazena os dados do cliente

const STATUS_ANDAMENTO = ['ABERTA', 'EM_ANDAMENTO', 'AGUARDANDO_PECA'];
const STATUS_FINALIZADAS = ['CONCLUIDA', 'ENTREGUE', 'CANCELADA'];

// ─── Elementos do DOM ────────────────────────────────────
const overlay = document.getElementById('os-form-overlay');
const formTitulo = document.getElementById('form-titulo');
const grupoStatus = document.getElementById('grupo-status');
const btnSalvar = document.getElementById('btn-salvar');
const formOS = document.getElementById('form-os');

const asBackdrop = document.getElementById('action-sheet-backdrop');
const asSheet = document.getElementById('action-sheet');
const asTitle = document.getElementById('action-sheet-title');
const asEditar = document.getElementById('as-editar');
const asFinalizar = document.getElementById('as-finalizar');
const asExcluir = document.getElementById('as-excluir');
const asCancelar = document.getElementById('as-cancelar');

const canvas = document.getElementById('sig-canvas');
const ctx = canvas.getContext('2d');

let osSelecionada = null;
let drawing = false;
let hasSignature = false;

// ─── Funções Auxiliares de Formatação ───────────────────
function formatarMoeda(valor) {
  if (!valor || isNaN(valor)) return '';
  return `R$ ${Number(valor).toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function parseMoedaParaNumero(str) {
  if (!str) return 0;
  const digitos = str.replace(/\D/g, '');
  if (!digitos) return 0;
  return parseFloat(digitos) / 100;
}

// ─── Buscar Cliente pelo ID ─────────────────────────────
async function buscarClientePorId(clienteId) {
  try {
    const res = await clientService.getById(clienteId);
    if (res.success && res.data) {
      return res.data;
    }
  } catch (err) {
    console.warn('Erro ao buscar cliente por ID', err);
  }
  return null;
}

// ─── Funções de Assinatura ──────────────────────────────
function resizeCanvas() {
  const rect = canvas.parentElement.getBoundingClientRect();
  canvas.width = rect.width || 400;
  canvas.height = 160;
  ctx.strokeStyle = '#1a1a1a';
  ctx.lineWidth = 2.5;
  ctx.lineCap = 'round';
  ctx.lineJoin = 'round';
}

function getPos(e) {
  const rect = canvas.getBoundingClientRect();
  const t = e.touches ? e.touches[0] : e;
  return { x: t.clientX - rect.left, y: t.clientY - rect.top };
}

function start(e) {
  e.preventDefault();
  drawing = true;
  const p = getPos(e);
  ctx.beginPath();
  ctx.moveTo(p.x, p.y);
}

function move(e) {
  if (!drawing) return;
  e.preventDefault();
  const p = getPos(e);
  ctx.lineTo(p.x, p.y);
  ctx.stroke();
  hasSignature = true;
}

function end() {
  drawing = false;
}

function setupAssinatura() {
  resizeCanvas();
  window.addEventListener('resize', resizeCanvas);

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
}

function limparAssinatura() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  hasSignature = false;
  document.getElementById('os-assinatura').value = '';
  document.getElementById('sig-status').textContent = 'Aguardando assinatura…';
}

function exibirAssinaturaExistente(assinaturaDados) {
  if (!assinaturaDados) return;

  document.getElementById('os-assinatura').value = assinaturaDados;
  document.getElementById('sig-status').textContent = '✓ Assinatura já registrada';

  if (typeof assinaturaDados === 'string' && assinaturaDados.startsWith('data:image')) {
    const img = new Image();
    img.onload = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
      hasSignature = true;
    };
    img.src = assinaturaDados;
  }
}

// ─── Funções do Formulário ──────────────────────────────
async function preencherFormulario(dados) {
  document.getElementById('os-id').value = dados.id || '';

  // ← BUSCAR TELEFONE DO CLIENTE PELO cliente_id
  let telefone = '';
  if (dados.cliente_id) {
    const cliente = await buscarClientePorId(dados.cliente_id);
    if (cliente) {
      clienteAtual = cliente;
      telefone = cliente.telefone || '';
    }
  }

  document.getElementById('os-telefone').value = telefone
    ? formatPhone(String(telefone).replace(/\D/g, ''))
    : '';

  const elNomeServico = document.getElementById('os-nome-servico');
  if (elNomeServico) {
    elNomeServico.value = (dados.nome_servico || '').toUpperCase();
  }

  document.getElementById('os-defeito').value = dados.defeito || '';

  let servicosTexto = dados.servicos || '';
  let valor = dados.valor != null ? Number(dados.valor) : 0;

  const matchValor = servicosTexto.match(/\|\s*Valor:\s*R\$\s*([\d.,]+)/i);
  if (matchValor) {
    valor = parseFloat(matchValor[1].replace(',', '.')) || valor;
    servicosTexto = servicosTexto.replace(/\s*\|\s*Valor:\s*R\$\s*[\d.,]+$/i, '').trim();
  }

  document.getElementById('os-servico').value = servicosTexto;
  document.getElementById('os-valor').value = valor > 0 ? formatarMoeda(valor) : '';
  document.getElementById('os-status').value = (dados.status || 'ABERTA').toUpperCase();

  limparAssinatura();
  exibirAssinaturaExistente(dados.assinatura);
}

function abrirFormulario(modo = 'nova', dados = null) {
  modoForm = modo;
  ordemAtual = dados;
  clienteAtual = null;

  formOS.reset();
  limparAssinatura();
  document.getElementById('os-alerta').classList.add('hidden');
  document.getElementById('os-id').value = '';

  if (modo === 'nova') {
    formTitulo.textContent = 'Nova Ordem de Serviço';
    grupoStatus.style.display = 'none';
    btnSalvar.textContent = 'Salvar Ordem de Serviço';
    document.getElementById('os-telefone').readOnly = false;
  } else {
    formTitulo.textContent = `Editar OS #${dados.numero_ordem || dados.id}`;
    grupoStatus.style.display = 'block';
    btnSalvar.textContent = 'Atualizar Ordem de Serviço';
    document.getElementById('os-telefone').readOnly = true;
    preencherFormulario(dados);
  }

  overlay.classList.add('open');
  setTimeout(() => {
    const first =
      modo === 'nova'
        ? document.getElementById('os-telefone')
        : document.getElementById('os-nome-servico') || document.getElementById('os-defeito');
    first?.focus();
  }, 300);
}

function fecharFormulario() {
  overlay.classList.remove('open');
  modoForm = 'nova';
  ordemAtual = null;
  clienteAtual = null;
}

// ─── Action Sheet ───────────────────────────────────────
function abrirActionSheet(os) {
  osSelecionada = os;
  const num = os.numero_ordem || os.id;
  asTitle.textContent = `OS #${num}${os.cliente_nome ? ' · ' + os.cliente_nome : ''}`;

  const status = (os.status || '').toUpperCase();
  const jaFinalizada = STATUS_FINALIZADAS.includes(status);

  asFinalizar.style.display = jaFinalizada ? 'none' : 'flex';

  asBackdrop.classList.add('open');
  asSheet.classList.add('open');
}

function fecharActionSheet() {
  asBackdrop.classList.remove('open');
  asSheet.classList.remove('open');
  osSelecionada = null;
}

// ─── Event Listeners ────────────────────────────────────
// Abas de filtro
document.querySelectorAll('#os-filter-tabs .tab').forEach((tab) => {
  tab.onclick = () => {
    document.querySelectorAll('#os-filter-tabs .tab').forEach((t) => t.classList.remove('active'));
    tab.classList.add('active');
    filtroAtual = tab.dataset.filter;
    renderListFiltrado();
  };
});

// Botões do formulário
document.getElementById('btn-fechar-form').onclick = fecharFormulario;
document.getElementById('btn-nova-os').onclick = () => abrirFormulario('nova');

// Action Sheet
asBackdrop.onclick = fecharActionSheet;
asCancelar.onclick = fecharActionSheet;

asEditar.onclick = async () => {
  if (!osSelecionada) return;
  const os = osSelecionada;
  fecharActionSheet();

  try {
    const res = await serviceOrderService.getById(os.id);
    if (res.success && res.data) {
      abrirFormulario('editar', res.data);
      return;
    }
  } catch (err) {
    console.warn('Falha ao buscar OS por ID, usando dados da lista', err);
  }
  abrirFormulario('editar', os);
};

asFinalizar.onclick = async () => {
  if (!osSelecionada) return;
  const os = osSelecionada;
  fecharActionSheet();

  const confirmar = confirm(
    `Deseja finalizar a OS #${os.numero_ordem || os.id}?\nO status será alterado para CONCLUÍDA.`
  );
  if (!confirmar) return;

  try {
    const response = await serviceOrderService.update({
      id: os.id,
      status: 'CONCLUIDA',
      nome_servico: os.nome_servico,
      defeito: os.defeito,
      servicos: os.servicos,
      assinatura: os.assinatura || null,
      valor: os.valor
    });

    if (response && response.success) {
      showToast('Ordem de serviço finalizada!', 'success');
      await carregarOrdensServico();
    } else {
      showToast(response?.message || 'Erro ao finalizar OS', 'error');
    }
  } catch (err) {
    console.error(err);
    showToast(err.message || 'Erro ao finalizar ordem', 'error');
  }
};

asExcluir.onclick = async () => {
  if (!osSelecionada) return;
  const os = osSelecionada;
  fecharActionSheet();

  const confirmar = confirm(
    `Tem certeza que deseja EXCLUIR a OS #${os.numero_ordem || os.id}?\nEssa ação não pode ser desfeita.`
  );
  if (!confirmar) return;

  try {
    const response = await serviceOrderService.remove(os.id);

    if (response && response.success) {
      showToast('Ordem de serviço excluída', 'success');
      await carregarOrdensServico();
    } else {
      showToast(response?.message || 'Erro ao excluir OS', 'error');
    }
  } catch (err) {
    console.error(err);
    showToast(err.message || 'Erro ao excluir ordem', 'error');
  }
};

// Alerta de duplicidade por telefone
document.getElementById('os-telefone').addEventListener('blur', async () => {
  if (modoForm === 'editar') return;
  const tel = document.getElementById('os-telefone').value.replace(/\D/g, '');
  const alerta = document.getElementById('os-alerta');
  if (tel.length < 10) {
    alerta.classList.add('hidden');
    return;
  }
  try {
    const res = await clientService.searchByPhone(tel);
    if (res.success && Array.isArray(res.data) && res.data.length > 0) {
      const cliente = res.data[0];
      clienteAtual = cliente; // ← Armazena o cliente encontrado
      const osRes = await serviceOrderService.listByClientId(cliente.id);
      const temAberta =
        osRes.success &&
        Array.isArray(osRes.data) &&
        osRes.data.some((o) => STATUS_ANDAMENTO.includes((o.status || '').toUpperCase()));
      alerta.classList.toggle('hidden', !temAberta);
    } else {
      alerta.classList.add('hidden');
      clienteAtual = null;
    }
  } catch {
    alerta.classList.add('hidden');
    clienteAtual = null;
  }
});

// Submit do formulário
formOS.onsubmit = async (e) => {
  e.preventDefault();

  const osId = document.getElementById('os-id')?.value;
  const telefone = document.getElementById('os-telefone').value.replace(/\D/g, '');
  const nomeServico = (document.getElementById('os-nome-servico')?.value || '').trim().toUpperCase();
  const defeito = document.getElementById('os-defeito').value.trim();
  const servicos = document.getElementById('os-servico').value.trim();
  const valorInputStr = document.getElementById('os-valor')?.value || '';
  const valor = parseMoedaParaNumero(valorInputStr);
  const assinatura = document.getElementById('os-assinatura').value || null;
  const statusSelecionado = document.getElementById('os-status')?.value || 'ABERTA';

  if (!telefone || !nomeServico || !defeito || !servicos) {
    showToast('Preencha todos os campos obrigatórios', 'error');
    return;
  }
  if (telefone.length < 10 || telefone.length > 11) {
    showToast('Telefone inválido (DDD + número)', 'error');
    return;
  }

  btnSalvar.disabled = true;
  btnSalvar.textContent = modoForm === 'nova' ? 'Salvando…' : 'Atualizando…';

  try {
    if (modoForm === 'nova') {
      // Tenta reaproveitar o ID se o cliente já estiver selecionado no contexto local
      let clienteId = clienteAtual?.id || null;

      const ordemPayload = {
        cliente_id: clienteId ? Number(clienteId) : undefined,
        telefone: telefone, // O backend Go vai buscar/criar o cliente com este número
        nome_servico: nomeServico,
        defeito,
        servicos,
        assinatura,
        status: 'ABERTA',
        valor
      };

      const response = await serviceOrderService.create(ordemPayload);
      if (response && response.success) {
        showToast('Ordem de serviço salva com sucesso', 'success');
        fecharFormulario();
        await carregarOrdensServico();
      } else {
        showToast(response?.message || 'Erro ao salvar ordem de serviço', 'error');
      }

    } else {
      // Edição de OS existente
      const osExistente = todasOrdens.find(o => String(o.id) === String(osId));

      const payload = {
        id: Number(osId),
        numero_ordem: osExistente?.numero_ordem || osExistente?.numeroOrdem || 0,
        telefone: telefone,
        nome_servico: nomeServico,
        defeito,
        servicos,
        assinatura,
        status: statusSelecionado,
        valor
      };

      const response = await serviceOrderService.update(payload);

      if (response && response.success) {
        showToast('Ordem de serviço atualizada com sucesso', 'success');
        fecharFormulario();
        await carregarOrdensServico();
      } else {
        showToast(response?.message || 'Erro ao atualizar ordem de serviço', 'error');
      }
    }
  } catch (err) {
    console.error('Erro ao enviar form OS:', err);
    showToast(err.message || 'Erro de conexão ao salvar OS', 'error');
  } finally {
    btnSalvar.disabled = false;
    btnSalvar.textContent =
      modoForm === 'nova' ? 'Salvar Ordem de Serviço' : 'Atualizar Ordem de Serviço';
  }
};

// ─── Renderização da Lista ──────────────────────────────
function badgeClass(status) {
  const s = (status || '').toUpperCase();
  if (['CONCLUIDA', 'ENTREGUE'].includes(s)) return 'badge-concluida';
  if (['CANCELADA'].includes(s)) return 'badge-cancelada';
  if (['EM_ANDAMENTO', 'AGUARDANDO_PECA'].includes(s)) return 'badge-pendente';
  return 'badge-aberta';
}

function renderListFiltrado() {
  const filtradas =
    filtroAtual === 'andamento'
      ? todasOrdens.filter((o) => STATUS_ANDAMENTO.includes((o.status || '').toUpperCase()))
      : todasOrdens.filter((o) => STATUS_FINALIZADAS.includes((o.status || '').toUpperCase()));

  renderList(filtradas);
}

async function renderList(items) {
  const el = document.getElementById('os-list');
  if (!items || !items.length) {
    const msg =
      filtroAtual === 'andamento'
        ? 'Nenhuma ordem em andamento.'
        : 'Nenhuma ordem finalizada.';
    el.innerHTML = `<div class="empty-state"><div class="icon">🛠️</div><p>${msg}</p></div>`;
    return;
  }

  el.innerHTML = items
    .map((os) => {
      const num = os.numero_ordem || os.id;
      const valorStr = os.valor > 0 ? ` · ${formatarMoeda(os.valor)}` : '';
      const nomeServico = os.nome_servico || 'Serviço Sem Nome';

      return `
      <div class="card list-item" data-id="${os.id}" style="flex-direction:column;align-items:stretch;gap:6px">
        <div class="flex justify-between items-center">
          <span class="fw-bold">OS #${num}</span>
          <span class="badge ${badgeClass(os.status)}">${os.status || 'ABERTA'}</span>
        </div>
        <div class="title">${nomeServico}</div>
        <div class="sub">${[os.defeito].filter(Boolean).join(' · ')}${valorStr}</div>
        ${os.servicos ? `<div class="sub">${os.servicos}</div>` : ''}
      </div>`;
    })
    .join('');

  el.querySelectorAll('.list-item').forEach((card) => {
    card.onclick = () => {
      const id = card.dataset.id;
      const os = todasOrdens.find((o) => String(o.id) === String(id));
      if (os) abrirActionSheet(os);
    };
  });
}
// ─── Carregar dados ─────────────────────────────────────
async function carregarOrdensServico() {
  try {
    const response = await serviceOrderService.listAll();
    if (response.success && Array.isArray(response.data)) {
      todasOrdens = response.data;
      renderListFiltrado();
    } else {
      todasOrdens = [];
      renderListFiltrado();
      if (response.message) showToast(response.message, 'error');
    }
  } catch (err) {
    console.error(err);
    todasOrdens = [];
    renderListFiltrado();
    showToast(err.message || 'Erro ao carregar ordens', 'error');
  }
}

// ─── Inicialização ──────────────────────────────────────
setupAssinatura();
carregarOrdensServico();