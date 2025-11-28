import { apiClient } from "../api/client";
import { BaseResponse } from "../api/model/response";
import { DockerPsRequest, DockerRunRequest } from "../api/model/request";

class DockerServicc {

  async ContainerPs(data: DockerPsRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/ps', data);
  }

  async BuildImageAndPush(data: DockerRunRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/build', data);
  }

  async RunContainer(data: DockerRunRequest): Promise<BaseResponse> {
    return apiClient.post<BaseResponse>('/run', data);
  }

}

export const dockerService = new DockerServicc();