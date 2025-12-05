export interface DockerPsRequest {
  host: string;
}

export interface DokcerBuildRequest {
  pjt_name: string;
  url: string;
  context_path: string;
}

export interface DockerRunRequest {
  host: string;
  image: string;
  port: string[];
  name: string;
  volume: string[];
  env: string[]
}

export interface DockerRestartRequest {
  host: string;
  name: string;
}


export interface DockerStopRequest {
  host: string;
  name: string;
}

export interface DockerRemoveRequest {
  host: string;
  name: string;
}

export interface DockerImageListRequest {
  host: string;
}
