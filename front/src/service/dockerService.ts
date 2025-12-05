import { apiClient } from "../api/client";
import { BaseResponse } from "../api/model/response";
import { DockerImageListRequest, DockerPsRequest, DockerRemoveRequest, DockerRestartRequest, DockerRunRequest, DockerStopRequest, DokcerBuildRequest } from "../api/model/request";

class DockerServicc {

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
}

export const dockerService = new DockerServicc();