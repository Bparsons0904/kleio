import { apiClient } from "../api";

export const exportHistory = async () => {
  try {
    const response = await apiClient.get("/export/history", {
      responseType: "blob",
    });

    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");
    link.href = url;

    const contentDisposition = response.headers["content-disposition"];
    let filename = "kleio_history_export.json";

    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename="?(.+)"?/);
      if (filenameMatch && filenameMatch[1]) {
        filename = filenameMatch[1];
      }
    } else {
      const now = new Date();
      const dateStr = now.toISOString().split("T")[0];
      filename = `kleio_history_export_${dateStr}.json`;
    }

    link.setAttribute("download", filename);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);

    return { success: true };
  } catch (error) {
    console.error("Error exporting history:", error);
    return { success: false, error };
  }
};
