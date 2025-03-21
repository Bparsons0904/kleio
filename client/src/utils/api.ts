import axios from 'axios';

// Create API client
export const apiClient = axios.create({
  baseURL: 'http://localhost:38080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const fetchApi = async (url: string) => {
  const data = await apiClient.get(`/${url}`);
  return data;
};
