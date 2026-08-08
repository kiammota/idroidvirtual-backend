import { requireAuth, goBack, showToast } from '../app.js';
import { AuditoriaService } from '../../client/services/AuditoriaService.js';

let auditorias = [];
let periodoAtual = 'dia';

async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();

  const modal = document.getElementById('modal-auditoria');
  const btnClose = document.getElementById('btn-close-modal');

  btnClose.onclick = () => modal.classList.add('hidden');
  modal.onclick = (e) => {
    if (e.target === modal) modal.classList.add('hidden');
  };

  // Botões de filtro
  document.getElementById('btn-dia').onclick = function () {
    this.className = 'btn btn-primary btn-sm';
    document.getElementById('btn-mes').className = 'btn btn-outline btn-sm';
    periodoAtual = 'dia';
    carregarAuditorias();
  };

  document.getElementById('btn-mes').onclick = function () {
    this.className = 'btn btn-primary btn-sm';
    document.getElementById('btn-dia').className = 'btn btn-outline btn-sm';
    periodoAtual = 'mes';
    carregarAuditorias();
  };

  async function carregarAuditorias() {
    try {
      const response = await AuditoriaService.listar({ periodo: periodoAtual, pagina: 1, itens_por_pagina: 50 });

      if (response.success && Array.isArray(response.data)) {
        auditorias = response.data;
      } else if (Array.isArray(response)) {
        auditorias = response;
      } else {
        auditorias = [];
      }
    } catch (err) {
      showToast('Erro ao carregar registros de auditoria', 'error');
    }
    renderLista();
  }

  function renderLista() {
    const listContainer = document.getElementById('auditoria-list');

    if (!auditorias || auditorias.length === 0) {
      listContainer.innerHTML = `<p class="text-muted" style="font-size:0.85rem; padding: 8px 0;">Nenhum registro encontrado.</p>`;
      return;
    }

    listContainer.innerHTML = auditorias
      .map((item) => {
        const dataFormatada = item.getHorarioFormatado ? item.getHorarioFormatado() : new Date(item.horario).toLocaleString('pt-BR');
        
        return `
          <div class="list-item" data-id="${item.id}" style="cursor: pointer; padding: 10px 0; border-bottom: 1px solid rgba(255,255,255,0.08); display: flex; justify-content: space-between; align-items: center;">
            <div>
              <div style="font-weight: 600;">${item.acao} - <span style="opacity: 0.8;">${item.tabela}</span></div>
              <div style="font-size: 0.78rem; opacity: 0.6;">Usuário #${item.usuario_id} • ${dataFormatada}</div>
            </div>
            <div style="font-size: 0.85rem; opacity: 0.5;">➜</div>
          </div>
        `;
      })
      .join('');

    // Evento de clique para abrir o modal
    listContainer.querySelectorAll('.list-item').forEach((el) => {
      el.onclick = () => {
        const id = el.dataset.id;
        const item = auditorias.find((a) => String(a.id) === String(id));
        if (item) abrirModal(item);
      };
    });
  }

  function abrirModal(item) {
    const content = document.getElementById('modal-content');
    const dataFormatada = item.getHorarioFormatado ? item.getHorarioFormatado() : new Date(item.horario).toLocaleString('pt-BR');

    content.innerHTML = `
      <div><strong>ID do Log:</strong> ${item.id}</div>
      <div><strong>Ação:</strong> ${item.acao}</div>
      <div><strong>Tabela Afetada:</strong> ${item.tabela} (ID: ${item.registro_id})</div>
      <div><strong>Usuário ID:</strong> ${item.usuario_id}</div>
      <div><strong>Data/Hora:</strong> ${dataFormatada}</div>

      <div style="margin-top: 10px;">
        <strong>Valores Antigos:</strong>
        <pre style="background: rgba(0,0,0,0.3); padding: 8px; border-radius: 4px; overflow-x: auto; font-size: 0.78rem;">${JSON.stringify(item.valores_antigos || {}, null, 2)}</pre>
      </div>

      <div>
        <strong>Valores Novos:</strong>
        <pre style="background: rgba(0,0,0,0.3); padding: 8px; border-radius: 4px; overflow-x: auto; font-size: 0.78rem;">${JSON.stringify(item.valores_novos || {}, null, 2)}</pre>
      </div>
    `;

    modal.classList.remove('hidden');
  }

  await carregarAuditorias();
}

init();