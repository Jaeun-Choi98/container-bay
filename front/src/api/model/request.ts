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