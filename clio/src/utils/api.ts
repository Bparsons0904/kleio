/* eslint-disable @typescript-eslint/no-explicit-any */
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
  console.log("Fetching:", url, apiClient.defaults.baseURL);
  const data = await apiClient.get(`/${url}`);
  return data;
};

export const postApi = async (url: string, data: any) => {
  const response = await apiClient.post(`/${url}`, data);
  return response;
};

export const putApi = async (url: string, data: any) => {
  const response = await apiClient.put(`/${url}`, data);
  return response;
};

export const deleteApi = async (url: string) => {
  const response = await apiClient.delete(`/${url}`);
  return response;
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

export const createStylus = async (stylus: any) => {
  return await postApi("styluses", stylus);
};

export const updateStylus = async (id: number, stylus: any) => {
  return await putApi(`styluses/${id}`, stylus);
};

export const deleteStylus = async (id: number) => {
  return await deleteApi(`styluses/${id}`);
};

export const getPlayCounts = async (limit = 50) => {
  return await fetchApi(`plays/counts?limit=${limit}`);
};

export const getRecentPlays = async (limit = 10) => {
  return await fetchApi(`plays/recent?limit=${limit}`);
};

export const createPlayHistory = async (playHistory: any) => {
  return await postApi("plays", playHistory);
};

export const updatePlayHistory = async (id: number, playHistory: any) => {
  return await putApi(`plays/${id}`, playHistory);
};

export const deletePlayHistory = async (id: number) => {
  return await deleteApi(`plays/${id}`);
};
export const createCleaningHistory = async (cleaningHistory: any) => {
  return await postApi("cleanings", cleaningHistory);
};

export const updateCleaningHistory = async (
  id: number,
  cleaningHistory: any,
) => {
  return await putApi(`cleanings/${id}`, cleaningHistory);
};

export const deleteCleaningHistory = async (id: number) => {
  return await deleteApi(`cleanings/${id}`);
};

export const refreshCollection = async () => {
  return await fetchApi("collection/resync");
};
