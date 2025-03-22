import axios from "axios";

// Create API client
export const apiClient = axios.create({
  baseURL: "http://localhost:38080",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

export const fetchApi = async (url: string) => {
  const data = await apiClient.get(`/${url}`);
  return data;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const postApi = async (url: string, data: any) => {
  const response = await apiClient.post(`/${url}`, data);
  return response;
};
