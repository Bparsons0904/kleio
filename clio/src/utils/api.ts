/* eslint-disable @typescript-eslint/no-explicit-any */
import axios from "axios";

// Create API client
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:38180/api",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

export const fetchApi = async (url: string) => {
  console.log("Fetching:", url, apiClient.defaults.baseURL);
  const data = await apiClient.get(`/${url}`);
  return data;
};

// Stylus API functions
export const getStyluses = async () => {
  return await fetchApi("styluses");
};

export const getActivePrimaryStylus = async () => {
  return await fetchApi("styluses/active");
};

export const getStylusByID = async (id: number) => {
  return await fetchApi(`styluses/${id}`);
};

export const getPlayCounts = async (limit = 50) => {
  return await fetchApi(`plays/counts?limit=${limit}`);
};

export const getRecentPlays = async (limit = 10) => {
  return await fetchApi(`plays/recent?limit=${limit}`);
};

export const refreshCollection = async () => {
  console.log("Starting collection sync...");
  const response = await apiClient.post("collection/resync");
  console.log("Collection sync response:", response.data);
  return response;
};
