<script lang="ts">
  import { onMount } from 'svelte';
  import { GetPredefinedServers, SyncTime, StartAutoSync } from '../wailsjs/go/main/App';
  import { EventsOn } from '../wailsjs/runtime/runtime';

  let predefinedServers: string[] = [];
  let selectedServer = '';
  let customServer = '';
  let logs = '';
  let interval = '';

  function addLog(message: string) {
    const timestamp = new Date().toLocaleString();
    logs = `[${timestamp}] ${message}\n` + logs;
  }

  onMount(async () => {
    try {
      predefinedServers = await GetPredefinedServers();
    } catch (error) {
      console.error("Ошибка при получении серверов:", error);
    }

    // Подписываемся на события из Go
    EventsOn('log', (message: string) => {
      addLog(message);
    });
  });

  async function syncTime() {
    const server = customServer || selectedServer;
    if (!server) {
      alert('Пожалуйста, выберите или введите NTP-сервер');
      return;
    }

    try {
      const ntpTime = await SyncTime(server);
      addLog(`Синхронизированное время с сервера ${server}: ${ntpTime}`);
      alert('Системное время синхронизировано');
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      addLog(errorMessage);
      alert(errorMessage);
    }
  }

  async function startAutoSync() {
    const intInterval = parseInt(interval);
    if (isNaN(intInterval) || intInterval <= 0) {
      alert('Пожалуйста, введите корректный положительный интервал в секундах');
      return;
    }

    try {
      await StartAutoSync(intInterval);
      addLog('Автоматическая синхронизация запущена');
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      addLog(errorMessage);
      alert(errorMessage);
    }
  }
</script>

<style>
  /* Ваши стили */
</style>

<div>
  <h1>NTP Синхронизация времени</h1>

  <select bind:value={selectedServer}>
    <option value="">Выберите NTP-сервер</option>
    {#each predefinedServers as server}
      <option value={server}>{server}</option>
    {/each}
  </select>

  <input type="text" placeholder="Или введите адрес NTP-сервера" bind:value={customServer} />

  <button on:click={syncTime}>Синхронизировать системное время</button>

  <h2>Настройка автоматической синхронизации</h2>
  <input type="text" placeholder="Введите интервал синхронизации (в секундах)" bind:value={interval} />
  <button on:click={startAutoSync}>Запустить автоматическую синхронизацию</button>

  <h2>Логи синхронизации</h2>
  <textarea readonly rows="10" bind:value={logs}></textarea>
</div>
