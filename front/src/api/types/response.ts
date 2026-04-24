export interface BaseResponse {
  result: number | null;
  data: any | null;
}

export interface Daemon {
  id: number;
  host: string;
  label: string;
}

export interface DockerResponse {
  execute_result: string[] | null;
  stdout: string[] | null;
  stderr: string[] | null;
}
