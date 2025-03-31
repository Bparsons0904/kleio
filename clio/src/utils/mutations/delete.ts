import { apiClient } from "../api";

export const deleteApi = async (url: string) => {
  const response = await apiClient.delete(`/${url}`);
  return response;
};

export const deleteStylus = async (id: number) => {
  return await deleteApi(`styluses/${id}`);
};

export const deletePlayHistory = async (id: number) => {
  return await deleteApi(`plays/${id}`);
};

export const deleteCleaningHistory = async (id: number) => {
  return await deleteApi(`cleanings/${id}`);
};
