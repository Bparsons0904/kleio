/* eslint-disable @typescript-eslint/no-explicit-any */
import { apiClient } from "../api";
import { Payload } from "../../provider/Provider";

export const postApi = async (url: string, data: any) => {
  const response = await apiClient.post(`/${url}`, data);
  return response;
};

export const createPlayHistory = async (playData: {
  releaseId: number;
  stylusId?: number | null;
  playedAt: string;
  notes?: string;
}) => {
  try {
    const response = await postApi("plays", playData);
    if (response.status !== 201) {
      throw new Error("Failed to log play");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error logging play:", error);
    return { success: false, error };
  }
};

export const createCleaningHistory = async (cleaningData: {
  releaseId: number;
  cleanedAt: string;
  notes?: string;
}) => {
  try {
    const response = await postApi("cleanings", cleaningData);
    if (response.status !== 201) {
      throw new Error("Failed to log cleaning");
    }
    return { success: true, data: response.data as Payload };
  } catch (error) {
    console.error("Error logging cleaning:", error);
    return { success: false, error };
  }
};

export const createPlayAndCleaning = async (
  playData: {
    releaseId: number;
    stylusId?: number | null;
    playedAt: string;
    notes?: string;
  },
  cleaningData: {
    releaseId: number;
    cleanedAt: string;
    notes?: string;
  },
) => {
  try {
    const cleaningResponse = await createCleaningHistory(cleaningData);
    if (!cleaningResponse.success) {
      throw new Error("Failed to log cleaning");
    }

    const playResponse = await createPlayHistory(playData);
    if (!playResponse.success) {
      throw new Error("Failed to log play");
    }

    return { success: true, data: playResponse.data };
  } catch (error) {
    console.error("Error logging both activities:", error);
    return { success: false, error };
  }
};

export const createStylus = async (stylus: any) => {
  return await postApi("styluses", stylus);
};
