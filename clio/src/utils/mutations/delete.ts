import { apiClient } from "../api";
import { Payload } from "../../provider/Provider";

export const deleteApi = async (url: string) => {
  const response = await apiClient.delete(`/${url}`);
  return response;
};

export const deletePlayHistory = async (id: number) => {
  try {
    const response = await deleteApi(`plays/${id}`);
    if (response.status !== 200) {
      throw new Error("Failed to delete play history");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error deleting play history:", error);
    return { success: false, error };
  }
};

export const deleteCleaningHistory = async (id: number) => {
  try {
    const response = await deleteApi(`cleanings/${id}`);
    if (response.status !== 200) {
      throw new Error("Failed to delete cleaning history");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error deleting cleaning history:", error);
    return { success: false, error };
  }
};

export const deleteStylus = async (id: number) => {
  return await deleteApi(`styluses/${id}`);
};

export const deleteRelease = async (id: number) => {
  try {
    const response = await deleteApi(`releases/${id}delete`);
    if (response.status !== 200) {
      throw new Error("Failed to delete release");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error deleting release:", error);
    return { success: false, error };
  }
};

export const archiveRelease = async (id: number) => {
  try {
    const response = await deleteApi(`releases/${id}/archive`);
    if (response.status !== 200) {
      throw new Error("Failed to archive release");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error archiving release:", error);
    return { success: false, error };
  }
};
