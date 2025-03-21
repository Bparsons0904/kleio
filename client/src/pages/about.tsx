import { createResource, Match, Show, Switch } from 'solid-js';
import { apiClient } from '../utils/api';

export const fetchApi = async (url: string) => {
  // const response = await fetch(`http://localhost:38080/health`);
  // console.log(response);
  // return response.json();
  const data = await apiClient.get(`/${url}`);
  return data;
};

export default function About() {
  const [data] = createResource('health', fetchApi);

  return (
    <div class="about">
      <Show when={data.loading}>Loading...</Show>
      <Show when={data.error}>Error: {data.error.message}</Show>
      <Show when={data()}>{JSON.stringify(data(), null, 2)}</Show>
      <Switch>
        <Match when={data()}>
          <div>
            <h2>System Health</h2>
          </div>
        </Match>
      </Switch>
    </div>
  );
}
