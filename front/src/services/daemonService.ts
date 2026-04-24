import { apiClient } from '../api/client';
import { BaseResponse } from '../api/types/response';
import { AddDaemonRequest } from '../api/types/request';

class DaemonService {
  async getDaemons(): Promise<BaseResponse> {
    return apiClient.get<BaseResponse>('/daemons');
  }

  async addDaemon(data: AddDaemonRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/daemons', data);
  }

  async removeDaemon(id: number): Promise<BaseResponse> {
    return apiClient.delete<BaseResponse>(`/daemons/${id}`);
  }
}

export const daemonService = new DaemonService();
