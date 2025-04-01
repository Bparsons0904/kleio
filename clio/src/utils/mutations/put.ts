/* eslint-disable @typescript-eslint/no-explicit-any */
import { apiClient } from "../api";
import { Payload } from "../../provider/Provider";

export const putApi = async (url: string, data: any) => {
  const response = await apiClient.put(`/${url}`, data);
  return response;
};

export const updatePlayHistory = async (
  id: number,
  playData: {
    releaseId: number;
    stylusId?: number | null;
    playedAt: string;
    notes?: string;
  },
) => {
  try {
    const response = await putApi(`plays/${id}`, playData);
    if (response.status !== 200) {
      throw new Error("Failed to update play history");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error updating play history:", error);
    return { success: false, error };
  }
};

export const updateCleaningHistory = async (
  id: number,
  cleaningData: {
    releaseId: number;
    cleanedAt: string;
    notes?: string;
  },
) => {
  try {
    const response = await putApi(`cleanings/${id}`, cleaningData);
    if (response.status !== 200) {
      throw new Error("Failed to update cleaning history");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error updating cleaning history:", error);
    return { success: false, error };
  }
};

export const updateStylus = async (id: number, stylus: any) => {
  return await putApi(`styluses/${id}`, stylus);
};
