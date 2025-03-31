/* eslint-disable @typescript-eslint/no-explicit-any */
import { apiClient } from "../api";

export const postApi = async (url: string, data: any) => {
  const response = await apiClient.post(`/${url}`, data);
  return response;
};

export const createStylus = async (stylus: any) => {
  return await postApi("styluses", stylus);
};

export const createPlayHistory = async (playHistory: any) => {
  return await postApi("plays", playHistory);
};

export const createCleaningHistory = async (cleaningHistory: any) => {
  return await postApi("cleanings", cleaningHistory);
};
