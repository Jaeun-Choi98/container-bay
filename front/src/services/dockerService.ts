import { apiClient } from "../api/client";
import { BaseResponse } from "../api/types/response";
import { DockerImageListRequest, DockerLogsRequest, DockerPsRequest, DockerRemoveRequest, DockerRestartRequest, DockerRunRequest, DockerStopRequest, DokcerBuildRequest } from "../api/types/request";

class DockerService {

  async ContainerPs(data: DockerPsRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/ps', data);
  }

  async BuildImageAndPush(data: DokcerBuildRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/build', data);
  }

  async RunContainer(data: DockerRunRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/run', data);
  }

  async StopContainer(data: DockerStopRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/stop', data);
  }

  async RestartContainer(data: DockerRestartRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/restart', data);
  }

  async RemovceContainer(data: DockerRemoveRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/rm', data);
  }

  async GetImageList(data: DockerImageListRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/image-ls', data);
  }

  async ImageRemove(data: DockerRemoveRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/image-rm', data);
  }

  async GetContainerLogs(data: DockerLogsRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/logs', data);
  }
}

export const dockerService = new DockerService();
