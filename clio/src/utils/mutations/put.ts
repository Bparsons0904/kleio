/* eslint-disable @typescript-eslint/no-explicit-any */
import { apiClient } from "../api";

export const putApi = async (url: string, data: any) => {
  const response = await apiClient.put(`/${url}`, data);
  return response;
};

export const updateStylus = async (id: number, stylus: any) => {
  return await putApi(`styluses/${id}`, stylus);
};

export const updatePlayHistory = async (id: number, playHistory: any) => {
  return await putApi(`plays/${id}`, playHistory);
};

export const updateCleaningHistory = async (
  id: number,
  cleaningHistory: any,
) => {
  return await putApi(`cleanings/${id}`, cleaningHistory);
};
