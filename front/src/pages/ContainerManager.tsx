import React, { useState } from 'react';
import { useApi } from '../api/hooks/useApi';
import { dockerService } from '../services/dockerService';
import {
  DockerPsRequest, DockerRunRequest, DockerStopRequest, DokcerBuildRequest,
  DockerRemoveRequest, DockerRestartRequest, DockerLogsRequest
} from '../api/types/request';
import { BaseResponse } from '../api/types/response';
import Nav from '../components/Header/Header';

const ContainerManage: React.FC = (): React.ReactElement => {
  const [pjtName, setPjtName] = useState<string>('');
  const [host, setHost] = useState<string>('');
  const [imageName, setImageName] = useState<string>('');
  const [volumes, setVolume] = useState<string>('');
  const [enves, setEnv] = useState<string>('');
  const [containerName, setContainerName] = useState<string>('');
  const [ports, setPorts] = useState<string>('');
  const [giturl, setGitUrl] = useState<string>('');
  const [buildContext, setBuildContext] = useState<string>('');
  const [logTail, setLogTail] = useState<string>('100');

  const psApi      = useApi<BaseResponse>((data: DockerPsRequest) => dockerService.ContainerPs(data));
  const runApi     = useApi<BaseResponse>((data: DockerRunRequest) => dockerService.RunContainer(data));
  const stopApi    = useApi<BaseResponse>((data: DockerStopRequest) => dockerService.StopContainer(data));
  const restartApi = useApi<BaseResponse>((data: DockerRestartRequest) => dockerService.RestartContainer(data));
  const removeApi  = useApi<BaseResponse>((data: DockerRemoveRequest) => dockerService.RemovceContainer(data));
  const buildApi   = useApi<BaseResponse>((data: DokcerBuildRequest) => dockerService.BuildImageAndPush(data));
  const logsApi    = useApi<BaseResponse>((data: DockerLogsRequest) => dockerService.GetContainerLogs(data));

  const handleGetContainers = async (): Promise<void> => {
    psApi.reset();
    await psApi.execute({ host });
  };

  const handleRunContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    const portArray = ports.split(',').map(p => p.trim()).filter(Boolean);
    const volumeArray = volumes.split(',').map(v => (process.env.VOLUME_DIR ?? '') + v.trim()).filter(Boolean);
    const envArray = enves.split(',').map(e => e.trim()).filter(Boolean);
    runApi.reset();
    await runApi.execute({ host, image: imageName, name: containerName, port: portArray, volume: volumeArray, env: envArray });
  };

  const handleStopContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    stopApi.reset();
    await stopApi.execute({ host, name: containerName });
  };

  const handleRestartContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    restartApi.reset();
    await restartApi.execute({ host, name: containerName });
  };

  const handleRemoveContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    removeApi.reset();
    await removeApi.execute({ host, name: containerName });
  };

  const handleBuildContainer = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    buildApi.reset();
    await buildApi.execute({ pjt_name: pjtName, url: giturl, context_path: buildContext });
  };

  const handleGetLogs = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    logsApi.reset();
    await logsApi.execute({ host, name: containerName, tail: parseInt(logTail) || 100 });
  };

  const hostInput = (
    <input
      className="input"
      type="text"
      placeholder="Host Docker Daemon (e.g. 192.168.0.1:2375)"
      value={host}
      onChange={e => setHost(e.target.value)}
      style={{ width: '300px' }}
      required
    />
  );

  const nameInput = (
    <input
      className="input"
      type="text"
      placeholder="Container Name"
      value={containerName}
      onChange={e => setContainerName(e.target.value)}
      style={{ width: '150px' }}
      required
    />
  );

  const renderApiError = (api: typeof psApi) =>
    api.error ? <div className="api-error">{api.error}</div> : null;

  const renderErrorBlock = (api: typeof psApi) =>
    api.result && api.result.result !== 0 ? (
      <div className="result-error-block">
        {api.result.data['stderr'].slice(1).map((item: string, idx: number) =>
          item.split(';').filter(col => col !== '').map((col, i) => (
            <div key={`${idx}-${i}`} className="result-error-line">{col}</div>
          ))
        )}
      </div>
    ) : null;

  const renderTableResult = (api: typeof psApi) =>
    api.result && api.result.result === 0 ? (
      <div className="result-table-wrap">
        <table className="result-table">
          <tbody>
            {api.result.data['stdout'].slice(1).map((item: string, idx: number) => (
              <tr key={idx}>
                <td>{idx + 1}</td>
                {item.split(';').filter(col => col !== '').map((col, i) => (
                  <td key={i}>{col}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    ) : null;

  const renderIdResult = (api: typeof psApi, label: string) =>
    api.result && api.result.result === 0 ? (
      <div className="result-success-msg">✓ {label}: {api.result.data['stdout'].slice(1).join('')}</div>
    ) : null;

  return (
    <div>
      <Nav />
      <div className="page">
        <h1 className="page-title">Container Management</h1>

        {/* Container List */}
        <section className="card">
          <h2 className="card-title">Container List</h2>
          <div className="form-row">
            <input
              className="input"
              type="text"
              placeholder="Host Docker Daemon (e.g. 192.168.0.1:2375)"
              value={host}
              onChange={e => setHost(e.target.value)}
              style={{ width: '300px' }}
            />
            <button className="btn btn-primary" onClick={handleGetContainers} disabled={psApi.loading}>
              {psApi.loading ? 'Loading...' : 'Get Containers'}
            </button>
          </div>
          {renderApiError(psApi)}
          {renderErrorBlock(psApi)}
          {renderTableResult(psApi)}
        </section>

        {/* Run Container */}
        <section className="card">
          <h2 className="card-title">Run Container</h2>
          <form className="form-row" onSubmit={handleRunContainer}>
            {hostInput}
            <input className="input" type="text" placeholder="Image Name" value={imageName} onChange={e => setImageName(e.target.value)} style={{ width: '140px' }} required />
            {nameInput}
            <input className="input" type="text" placeholder="Ports (8080:80, 3000:3000, ...)" value={ports} onChange={e => setPorts(e.target.value)} style={{ width: '210px' }} />
            <input className="input" type="text" placeholder="Volumes (/src:/dst, ...)" value={volumes} onChange={e => setVolume(e.target.value)} style={{ width: '180px' }} />
            <input className="input" type="text" placeholder="Env (KEY=val, ...)" value={enves} onChange={e => setEnv(e.target.value)} style={{ width: '160px' }} />
            <button className="btn btn-primary" type="submit" disabled={runApi.loading}>
              {runApi.loading ? 'Running...' : 'Run Container'}
            </button>
          </form>
          {renderApiError(runApi)}
          {renderErrorBlock(runApi)}
          {renderIdResult(runApi, 'Container ID')}
        </section>

        {/* Stop Container */}
        <section className="card">
          <h2 className="card-title">Stop Container</h2>
          <form className="form-row" onSubmit={handleStopContainer}>
            {hostInput}
            {nameInput}
            <button className="btn btn-primary" type="submit" disabled={stopApi.loading}>
              {stopApi.loading ? 'Stopping...' : 'Stop Container'}
            </button>
          </form>
          {renderApiError(stopApi)}
          {renderErrorBlock(stopApi)}
          {renderIdResult(stopApi, 'Container ID')}
        </section>

        {/* Restart Container */}
        <section className="card">
          <h2 className="card-title">Restart Container</h2>
          <form className="form-row" onSubmit={handleRestartContainer}>
            {hostInput}
            {nameInput}
            <button className="btn btn-primary" type="submit" disabled={restartApi.loading}>
              {restartApi.loading ? 'Restarting...' : 'Restart Container'}
            </button>
          </form>
          {renderApiError(restartApi)}
          {renderErrorBlock(restartApi)}
          {renderIdResult(restartApi, 'Container ID')}
        </section>

        {/* Remove Container */}
        <section className="card">
          <h2 className="card-title">Remove Container</h2>
          <form className="form-row" onSubmit={handleRemoveContainer}>
            {hostInput}
            {nameInput}
            <button className="btn btn-primary" type="submit" disabled={removeApi.loading}>
              {removeApi.loading ? 'Removing...' : 'Remove Container'}
            </button>
          </form>
          {renderApiError(removeApi)}
          {renderErrorBlock(removeApi)}
          {renderIdResult(removeApi, 'Container ID')}
        </section>

        {/* Container Logs */}
        <section className="card">
          <h2 className="card-title">Container Logs</h2>
          <form className="form-row" onSubmit={handleGetLogs}>
            {hostInput}
            {nameInput}
            <input
              className="input"
              type="number"
              placeholder="Tail lines"
              value={logTail}
              onChange={e => setLogTail(e.target.value)}
              style={{ width: '120px' }}
            />
            <button className="btn btn-primary" type="submit" disabled={logsApi.loading}>
              {logsApi.loading ? 'Loading...' : 'Get Logs'}
            </button>
          </form>
          {renderApiError(logsApi)}
          {renderErrorBlock(logsApi)}
          {logsApi.result && logsApi.result.result === 0 && (
            <div className="result-log">
              {logsApi.result.data['stdout'].slice(1).map((item: string, idx: number) => (
                <div key={idx} className="result-log-line">{item.replace(/;$/, '')}</div>
              ))}
            </div>
          )}
        </section>

        {/* Build Docker Image */}
        <section className="card">
          <h2 className="card-title">Build Docker Image</h2>
          <form className="form-row" onSubmit={handleBuildContainer}>
            <input className="input" type="text" placeholder="Project Name" value={pjtName} onChange={e => setPjtName(e.target.value)} style={{ width: '150px' }} required />
            <input className="input" type="text" placeholder="Git URL" value={giturl} onChange={e => setGitUrl(e.target.value)} style={{ width: '280px' }} required />
            <input className="input" type="text" placeholder="Dockerfile Context Path" value={buildContext} onChange={e => setBuildContext(e.target.value)} style={{ width: '190px' }} required />
            <button className="btn btn-primary" type="submit" disabled={buildApi.loading}>
              {buildApi.loading ? 'Building...' : 'Build Image'}
            </button>
          </form>
          {renderApiError(buildApi)}
          {renderErrorBlock(buildApi)}
          {buildApi.result && buildApi.result.result === 0 && (
            <div className="result-table-wrap">
              <table className="result-table">
                <tbody>
                  {buildApi.result.data['stdout'].slice(1).map((item: string, idx: number) => (
                    <tr key={idx}>
                      <td>{idx + 1}</td>
                      {item.split(';').filter(col => col !== '').map((col, i) => (
                        <td key={i}>{col}</td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </div>
  );
};

export default ContainerManage;
