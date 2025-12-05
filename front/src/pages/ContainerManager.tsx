import React, { useState } from 'react';
import { useApi } from '../api/hooks/useApi';
import { dockerService } from '../service/dockerService';
import {
  DockerPsRequest, DockerRunRequest, DockerStopRequest, DokcerBuildRequest, DockerRemoveRequest
  , DockerRestartRequest
} from '../api/model/request';
import { BaseResponse } from '../api/model/response';

/*
host: string;
  name: string;
*/
const ContainerManage: React.FC = (): React.ReactElement => {
  const [pjtName, setPjtName] = useState<string>('');
  const [host, setHost] = useState<string>('');
  const [imageName, setImageName] = useState<string>('');
  const [volumes, setVolume] = useState<string>();
  const [enves, setEnv] = useState<string>();
  const [containerName, setContainerName] = useState<string>('');
  const [ports, setPorts] = useState<string>('');
  const [giturl, setGitUrl] = useState<string>('');
  const [buildContext, setBuildContext] = useState<string>('');

  const psApi = useApi<BaseResponse>((data: DockerPsRequest): Promise<BaseResponse> => {
    return dockerService.ContainerPs(data);
  });

  const runApi = useApi<BaseResponse>((data: DockerRunRequest): Promise<BaseResponse> => {
    return dockerService.RunContainer(data);
  });

  const stopApi = useApi<BaseResponse>((data: DockerStopRequest): Promise<BaseResponse> => {
    return dockerService.StopContainer(data);
  });

  const restartApi = useApi<BaseResponse>((data: DockerRestartRequest): Promise<BaseResponse> => {
    return dockerService.RestartContainer(data);
  })

  const removeApi = useApi<BaseResponse>((data: DockerRemoveRequest): Promise<BaseResponse> => {
    return dockerService.RemovceContainer(data);
  })

  const buildApi = useApi<BaseResponse>((data: DokcerBuildRequest): Promise<BaseResponse> => {
    return dockerService.BuildImageAndPush(data);
  });

  const handleGetContainers = async (): Promise<void> => {
    psApi.reset();
    await psApi.execute({ host });
  };

  const handleRunContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();

    const portArray: string[] | undefined = ports?.split(',').map((p: string): string => p.trim());
    const volumeArray: string[] | undefined = volumes?.split(',').
      map((v: string): string => {
        return process.env.VOLUME_DIR + v.trim();
      });
    const envArray: string[] | undefined = enves?.split(',').map((e: string): string => e.trim());

    runApi.reset();
    await runApi.execute({
      host: host,
      image: imageName,
      name: containerName,
      port: portArray,
      volume: volumeArray,
      env: envArray,
    });
  };

  const handleStopContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    stopApi.reset();
    await stopApi.execute({
      host: host,
      name: containerName,
    });
  }

  const handleRestartContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    restartApi.reset();
    await restartApi.execute({
      host: host,
      name: containerName,
    });
  }

  const handleRemoveContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    removeApi.reset();
    await removeApi.execute({
      host: host,
      name: containerName,
    }
    );
  }

  const handleBuildContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    buildApi.reset();
    await buildApi.execute({
      pjt_name: pjtName,
      url: giturl,
      context_path: buildContext,
    });
  }

  return (
    <div>
      <h1>Container Management</h1>

      {/* 컨테이너 목록 조회 */}
      <section>
        <h2>Container List</h2>
        <input
          type="text"
          placeholder="Host"
          value={host}
          onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
            setHost(e.target.value)
          }
        />
        <button onClick={handleGetContainers} disabled={psApi.loading}>
          {psApi.loading ? 'Loading...' : 'Get Containers'}
        </button>

        {psApi.error && <div style={{ color: 'red' }}>Error: {psApi.error}</div>}
        {psApi.result && psApi.result.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {
              psApi.result.data["stderr"].slice(1).map((item: string, index: number): React.ReactElement => {
                const colums = item.split(';');
                return (
                  <ul style={{ margin: 'auto', width: '50%' }}>
                    {colums.map((col, i): React.ReactElement | null => {
                      if (col === "") return null;
                      return (
                        < li key={i} style={{ width: '100%', padding: '5px', listStyle: 'none' }}> {col}</li>
                      );
                    })}
                  </ul>
                );
              })

            }
          </pre>
        }
        {psApi.result && psApi.result.result === 0 && (
          <pre>{psApi.result.data["stdout"].slice(1).map((item: string, index: number): React.ReactElement => {
            const colums = item.split(';');
            return (
              <tr key={index} style={{ tableLayout: 'fixed', width: '100%' }}>
                <td style={{ width: '5%', border: '1px solid #ccc', padding: '5px' }} >{index !== 0 ? index : ""}</td>
                {colums.map((col, i): React.ReactElement | null => {
                  if (col === "") return null;
                  return (
                    <td key={i} style={{ width: '5%', border: '1px solid #ccc', padding: '5px' }}>{col}</td>
                  );
                }
                )}
              </tr>
            );
          }
          )}</pre>
        )}
      </section>

      {/* 컨테이너 실행 */}
      <section>
        <h2>Run Container</h2>
        <form onSubmit={handleRunContainer}>
          <input
            type="text"
            placeholder="Host Docker Daemon(e.g. 192.168.0.1:2375)"
            value={host}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setHost(e.target.value)
            }
            style={{ width: '300px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="Only Image Name"
            value={imageName}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setImageName(e.target.value)
            }
            style={{ width: '120px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="Container Name"
            value={containerName}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setContainerName(e.target.value)
            }
            style={{ width: '120px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="Ports(e.g. 8080:80,3000:3000,...)"
            value={ports}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setPorts(e.target.value)
            }
            style={{ width: '200px', padding: '2px' }}
          />
          <input
            type="text"
            placeholder="Volumes(e.g. /file:/file,/dir:/dir,...)"
            value={volumes}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setVolume(e.target.value)
            }
            style={{ width: '220px', padding: '2px' }}
          />
          <input
            type="text"
            placeholder="Enves(e.g. key1=value1,key2=value2,...)"
            value={enves}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setEnv(e.target.value)
            }
            style={{
              width: "240px", padding: "2px"
            }}
          />
          <button
            type="submit" disabled={runApi.loading}>
            {runApi.loading ? 'Running...' : 'Run Container'}
          </button>
        </form>
        {runApi.error && <div style={{ color: 'red' }}>Error: {runApi.error}</div>}
        {runApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {
              runApi.result?.data["stderr"].slice(1).map((item: string, index: number): React.ReactElement => {
                const colums = item.split(';');
                return (
                  <ul style={{ margin: 'auto', width: '50%' }}>
                    {colums.map((col, i): React.ReactElement | null => {
                      if (col === "") return null;
                      return (
                        < li key={i} style={{ width: '100%', padding: '5px', listStyle: 'none' }}> {col}</li>
                      );
                    })}
                  </ul>
                );
              })

            }
          </pre>
        }
        {
          runApi.result?.result === 0 && (
            <div style={{ color: 'green' }}>
              Container started successfully!
              <pre style={{ whiteSpace: 'pre-wrap' }}>{'Container ID: ' + runApi.result.data['stdout'].slice(1)}</pre>
            </div>
          )
        }
      </section >

      {/* 컨테이너 중지 */}
      <h2>Stop Container</h2>
      <section>
        <form onSubmit={handleStopContainer}>
          <input
            type="text"
            placeholder="Host Docker Daemon(e.g. 192.168.0.1:2375)"
            value={host}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setHost(e.target.value)
            }
            style={{ width: '300px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="Container Name"
            value={containerName}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setContainerName(e.target.value)
            }
            style={{ width: '120px', padding: '2px' }}
            required
          />
          <button type="submit" disabled={stopApi.loading}>
            {stopApi.loading ? 'Running...' : 'Stop Container'}
          </button>
        </form>
        {stopApi.error && <div style={{ color: 'red' }}>Error: {stopApi.error}</div>}
        {
          stopApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {stopApi.result?.data['stderr'].slice(1).map((item: string, index: number): React.ReactElement | null => {
              const colums: string[] = item.split(";");
              return (
                <ul style={{ margin: "auto", width: "50%" }}>
                  {colums.map((col: string, i: number): React.ReactElement | null => {
                    if (col === "") return null;
                    return (
                      <li key={i} style={{ width: "100%", padding: "5px", listStyle: "none" }}>{col}</li>
                    );
                  })}
                </ul>
              )
            })}
          </pre>
        }
        {
          stopApi.result?.result === 0 && (
            <div style={{ color: 'green' }}>
              Container started successfully!
              <pre style={{ whiteSpace: 'pre-wrap' }}>{'Container ID: ' + stopApi.result.data['stdout'].slice(1)}</pre>
            </div>
          )
        }
      </section>

      {/* 컨테이너 재시작 */}
      <h2>Restart Container</h2>
      <section>
        <form onSubmit={handleRestartContainer}>
          <input
            type="text"
            placeholder="Host Docker Daemon(e.g. 192.168.0.1:2375)"
            value={host}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setHost(e.target.value)
            }
            style={{ width: '300px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="Container Name"
            value={containerName}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setContainerName(e.target.value)
            }
            style={{ width: '120px', padding: '2px' }}
            required
          />
          <button type="submit" disabled={stopApi.loading}>
            {restartApi.loading ? 'Running...' : 'Restart Container'}
          </button>
        </form>
        {restartApi.error && <div style={{ color: 'red' }}>Error: {restartApi.error}</div>}
        {
          restartApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {restartApi.result?.data['stderr'].slice(1).map((item: string, index: number): React.ReactElement | null => {
              const colums: string[] = item.split(";");
              return (
                <ul style={{ margin: "auto", width: "50%" }}>
                  {colums.map((col: string, i: number): React.ReactElement | null => {
                    if (col === "") return null;
                    return (
                      <li key={i} style={{ width: "100%", padding: "5px", listStyle: "none" }}>{col}</li>
                    );
                  })}
                </ul>
              )
            })}
          </pre>
        }
        {
          restartApi.result?.result === 0 && (
            <div style={{ color: 'green' }}>
              Container started successfully!
              <pre style={{ whiteSpace: 'pre-wrap' }}>{'Container ID: ' + restartApi.result.data['stdout'].slice(1)}</pre>
            </div>
          )
        }
      </section>

      {/* 컨테이너 삭제 */}
      <h2>Remove Container</h2>
      <section>
        <form onSubmit={handleRemoveContainer}>
          <input
            type="text"
            placeholder="Host Docker Daemon(e.g. 192.168.0.1:2375)"
            value={host}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setHost(e.target.value)
            }
            style={{ width: '300px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="Container Name"
            value={containerName}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setContainerName(e.target.value)
            }
            style={{ width: '120px', padding: '2px' }}
            required
          />
          <button type="submit" disabled={removeApi.loading}>
            {removeApi.loading ? 'Running...' : 'Remove Container'}
          </button>
        </form>
        {removeApi.error && <div style={{ color: 'red' }}>Error: {removeApi.error}</div>}
        {
          removeApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {removeApi.result?.data['stderr'].slice(1).map((item: string, index: number): React.ReactElement | null => {
              const colums: string[] = item.split(";");
              return (
                <ul style={{ margin: "auto", width: "50%" }}>
                  {colums.map((col: string, i: number): React.ReactElement | null => {
                    if (col === "") return null;
                    return (
                      <li key={i} style={{ width: "100%", padding: "5px", listStyle: "none" }}>{col}</li>
                    );
                  })}
                </ul>
              )
            })}
          </pre>
        }
        {
          removeApi.result?.result === 0 && (
            <div style={{ color: 'green' }}>
              Container started successfully!
              <pre style={{ whiteSpace: 'pre-wrap' }}>{'Container ID: ' + removeApi.result.data['stdout'].slice(1)}</pre>
            </div>
          )
        }
      </section>

      {/* 깃 클론 및 이미지 빌드 */}
      < section >
        <h2>Build Docker Image</h2>
        <form onSubmit={handleBuildContainer}>
          <input
            type="text"
            placeholder="Project Name"
            value={pjtName}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setPjtName(e.target.value)
            }
            style={{ width: '150px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="URL"
            value={giturl}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setGitUrl(e.target.value)
            }
            style={{ width: '250px', padding: '2px' }}
            required
          />
          <input
            type="text"
            placeholder="Dockerfile Context Path"
            value={buildContext}
            onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
              setBuildContext(e.target.value)
            }
            style={{ width: '150px', padding: '2px' }}
            required
          />
          <button
            type="submit" disabled={buildApi.loading}>
            {runApi.loading ? 'Running...' : 'Build Image'}
          </button>
        </form>
        {buildApi.error && <div style={{ color: 'red' }}>Error: {buildApi.error}</div>}
        {
          buildApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {
              buildApi.result?.data['stderr'].slice(1).map((item: string, index: number): React.ReactElement => {
                const colums = item.split(';');
                return (
                  <ul style={{ margin: 'auto', width: '50%' }}>
                    {colums.map((col: string, i: number): React.ReactElement | null => {
                      if (col === "") return null
                      return (
                        <li key={i} style={{ width: '100%', padding: '5px', listStyle: 'none' }}>{col}</li>
                      );
                    })}
                  </ul>
                );
              })
            }
          </pre>
        }
        {buildApi.result?.result === 0 && (<div style={{ color: 'green' }}>
          Image Builed successfully!
          <pre style={{ whiteSpace: 'pre-wrap' }}>
            {
              buildApi.result.data['stdout'].slice(1).map((item: string, index: number): React.ReactElement => {
                const cloums: string[] = item.split(";");
                return (
                  <ul style={{ margin: 'auto', width: '50%' }}>
                    {cloums.map((col: string, i: number): React.ReactElement | null => {
                      if (col === "") return null;
                      return (
                        <li key={i} style={{ width: '100%', padding: '5px', listStyle: 'none' }} > {col} </li>
                      );
                    })}
                  </ul>
                );
              })
            }
          </pre>
        </div>)}
      </section >

      {/* 볼륨 디렉토리 */}
      < section >
        <h2>Volume Directory</h2>
        <form>

        </form>
      </section >

    </div >
  );
};

export default ContainerManage; 