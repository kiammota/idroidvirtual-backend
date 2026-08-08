import { requireAuth, goBack, showToast, formatPhone } from '../app.js';
import { bindMasks } from '../masks.js';
import { Cliente } from '../../client/models/ClienteModel.js';
import { ClientService } from '../../client/services/clientService.js';
import { NotificationService } from '../../client/services/notificationService.js';
import { openInspect } from '../components/InspectModal.js';

window.NotificationService = NotificationService;

/**
 * Página de Clientes
 */
async function init() {
  if (!requireAuth()) return;

  document.getElementById('btn-back').onclick = () => goBack();
  bindMasks({ 'cli-telefone': 'phone' });

  let clientes = [];

  document.querySelectorAll('#cli-tabs .tab').forEach((tab) => {
    tab.onclick = () => {
      document.querySelectorAll('#cli-tabs .tab').forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');
      document.getElementById('tab-lista').classList.toggle('hidden', tab.dataset.tab !== 'lista');
      document.getElementById('tab-novo').classList.toggle('hidden', tab.dataset.tab !== 'novo');
    };
  });

  function abrirInspectCliente(c) {
    const criado = c.criado_em
      ? new Date(c.criado_em).toLocaleString('pt-BR')
      : '—';

    openInspect({
      title: 'Cliente',
      subtitle: 'IDroid Virtual',
      fields: [
        { label: 'ID', value: `#${c.id}` },
        { label: 'Nome', value: c.nome || '—' },
        { label: 'Telefone', value: formatPhone(c.telefone) || '—' },
        { label: 'E-mail', value: c.email || '—' },
        { label: 'Endereço', value: c.endereco || '—' },
        { label: 'Cadastrado em', value: criado },
      ],
      footer: 'IDroid Virtual · Gestão',
    });
  }

  function render(list) {
    const el = document.getElementById('cli-list');
    if (!list || !list.length) {
      el.innerHTML = `<div class="empty-state"><div class="icon">👥</div><p>Nenhum cliente encontrado.</p></div>`;
      return;
    }
    el.innerHTML = list
      .map(
        (c) => {
          const nome = c.nome || 'NÃO INFORMADO';
          const tel = formatPhone(c.telefone) || '—';
          const extra = c.email || c.endereco || 'NÃO INFORMADO';
          return `
      <div class="card list-item" data-id="${c.id}" style="cursor:pointer">
        <div class="info">
          <div class="title" title="${String(nome).replace(/"/g, '&quot;')}">${nome}</div>
          <div class="sub">${tel} · ${extra}</div>
        </div>
      </div>
    `;
        }
      )
      .join('');

    el.querySelectorAll('.list-item').forEach((card) => {
      card.onclick = () => {
        const id = Number(card.dataset.id);
        const cliente = clientes.find((x) => Number(x.id) === id);
        if (cliente) abrirInspectCliente(cliente);
      };
    });
  }

  document.getElementById('busca-cliente').oninput = (e) => {
    const q = e.target.value.toLowerCase().trim();
    if (!q) return render(clientes);
    render(
      clientes.filter(
        (c) =>
          (c.nome || '').toLowerCase().includes(q) ||
          (c.telefone || '').includes(q.replace(/\D/g, '')) ||
          (c.email || '').toLowerCase().includes(q)
      )
    );
  };

  document.getElementById('form-cliente').onsubmit = async (e) => {
    e.preventDefault();

    const nome = document.getElementById('cli-nome').value;
    const telefone = document.getElementById('cli-telefone').value;
    const endereco = document.getElementById('cli-endereco').value;
    const email = document.getElementById('cli-email').value;

    const model = new Cliente({ nome, telefone, endereco, email });
    if (!model.isValid()) {
      showToast(model.validar()[0], 'error');
      return;
    }

    try {
      const response = await ClientService.create({ nome, telefone, endereco, email });

      if (response.success) {
        showToast('Cliente cadastrado com sucesso', 'success');
        document.getElementById('form-cliente').reset();
        await carregarClientes();
        document.querySelector('[data-tab="lista"]').click();
      } else {
        showToast(response.message || 'Erro ao cadastrar cliente', 'error');
      }
    } catch (err) {
      console.error('Erro ao criar cliente:', err);
      showToast(err.message || 'Erro de conexão ao cadastrar cliente', 'error');
    }
  };

  async function carregarClientes() {
    try {
      const response = await ClientService.listAll();
      if (response.success && Array.isArray(response.data)) {
        clientes = response.data;
      } else {
        clientes = [];
      }
    } catch (err) {
      showToast('Erro ao carregar clientes', 'error');
      clientes = [];
    } finally {
      render(clientes);
    }
  }

  await carregarClientes();
}

init();
