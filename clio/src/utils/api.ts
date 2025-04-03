/* eslint-disable @typescript-eslint/no-explicit-any */
import axios from "axios";

// Create API client
export const apiClient = axios.create({
  baseURL: "http://localhost:38080/api",
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
  return await fetchApi("collection/resync");
};
