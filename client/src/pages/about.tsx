import { createEffect, createSignal } from 'solid-js';
import { useApi } from '../utils/api';

export default function About() {
  const api = useApi<Record<string, string>>();

  // State to store the health data
  const [healthData, setHealthData] = createSignal<Record<string, string> | null>(null);

  // Function to fetch health data
  const checkHealth = async () => {
    try {
      const response = await api.get('/health');
      setHealthData(response.data);
    } catch (error) {
      console.error('Health check failed:', error);
      // You could set some error state here if needed
    }
  };

  // Call the health endpoint when the component mounts
  createEffect(() => {
    checkHealth();
  });

  return (
    <div>
      <h2>System Health</h2>
      {api.isLoading() && <p>Loading health information...</p>}
      {api.error() && <p class="error">Error: {api.error()?.message}</p>}

      {healthData() && (
        <div class="health-stats">
          <p>
            Status:{' '}
            <span class={healthData()?.status === 'up' ? 'status-up' : 'status-down'}>
              {healthData()?.status}
            </span>
          </p>

          {healthData()?.message && <p>Message: {healthData()?.message}</p>}

          {/* Display other health metrics if they exist */}
          {healthData()?.open_connections && (
            <div class="metrics">
              <h3>Database Metrics</h3>
              <p>Open Connections: {healthData()?.open_connections}</p>
              <p>In Use: {healthData()?.in_use}</p>
              <p>Idle: {healthData()?.idle}</p>
              <p>Wait Count: {healthData()?.wait_count}</p>
              <p>Wait Duration: {healthData()?.wait_duration}</p>
            </div>
          )}
        </div>
      )}

      <button onClick={checkHealth} disabled={api.isLoading()}>
        Refresh Health Status
      </button>
    </div>
  );
}
