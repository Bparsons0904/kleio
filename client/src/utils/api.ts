import axios, { AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { createSignal, createResource } from 'solid-js';

// Create a base axios instance with common configuration
const apiClient = axios.create({
  baseURL: 'http://localhost:38080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Create a useApi composable for SolidJS
export function useApi<T>() {
  const [isLoading, setIsLoading] = createSignal(false);
  const [error, setError] = createSignal<AxiosError | null>(null);

  // Generic request function
  async function request<R = T>(config: AxiosRequestConfig): Promise<AxiosResponse<R>> {
    setIsLoading(true);
    setError(null);

    try {
      const response = await apiClient(config);
      return response;
    } catch (err) {
      setError(err as AxiosError);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }

  // Convenience methods
  const get = <R = T>(url: string, config?: AxiosRequestConfig) =>
    request<R>({ ...config, method: 'get', url });

  //  eslint-disable-next-line @typescript-eslint/no-explicit-any
  const post = <R = T>(url: string, data?: any, config?: AxiosRequestConfig) =>
    request<R>({ ...config, method: 'post', url, data });

  //  eslint-disable-next-line @typescript-eslint/no-explicit-any
  const put = <R = T>(url: string, data?: any, config?: AxiosRequestConfig) =>
    request<R>({ ...config, method: 'put', url, data });

  const del = <R = T>(url: string, config?: AxiosRequestConfig) =>
    request<R>({ ...config, method: 'delete', url });

  return {
    isLoading,
    error,
    request,
    get,
    post,
    put,
    del,
  };
}

// Optional: Create a resource-based API utility
//  eslint-disable-next-line @typescript-eslint/no-explicit-any
export function createApiResource<T>(fetcher: () => Promise<T>, options?: { initialValue?: any }) {
  const [resource, { refetch, mutate }] = createResource<T>(fetcher, options);

  return {
    data: resource,
    loading: () => resource.loading,
    error: () => resource.error,
    refetch,
    mutate,
  };
}
