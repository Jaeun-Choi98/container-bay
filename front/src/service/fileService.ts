import { apiClient } from "../api/client";
import { BaseResponse } from "../api/model/response";

class FileService {
  async uploadFile(file: File, additionalData?: Record<string, string>): Promise<BaseResponse> {
    return apiClient.uploadFile<BaseResponse>('/upload-file', file, additionalData);
  }
}

export const fileService = new FileService();