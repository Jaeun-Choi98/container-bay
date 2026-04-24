import React, { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { useApi } from '../api/hooks/useApi';
import { dockerService } from '../services/dockerService';
import {
  DockerPsRequest, DockerRunRequest, DockerStopRequest, DockerRestartRequest,
  DockerRemoveRequest, DokcerBuildRequest, DockerImageListRequest,
  DcokerImageRemoveRequest, DockerLogsRequest
} from '../api/types/request';
import { BaseResponse } from '../api/types/response';
import Nav from '../components/Header/Header';

type Tab = 'containers' | 'images' | 'logs';

/* ── Shared output helpers ─────────────────────────── */
const ApiError = ({ api }: { api: { error: string | null } }) =>
  api.error ? <div className="api-error">{api.error}</div> : null;

const ErrorBlock = ({ api }: { api: { result: BaseResponse | null } }) =>
  api.result && api.result.result !== 0 ? (
    <div className="result-error-block">
      {(api.result.data?.['stderr'] as string[] ?? []).slice(1)
        .flatMap((item: string, idx: number) =>
          item.split(';').filter(c => c !== '').map((col, i) => (
            <div key={`${idx}-${i}`} className="result-error-line">{col}</div>
          ))
        )}
    </div>
  ) : null;

const TableResult = ({ api }: { api: { result: BaseResponse | null } }) =>
  api.result && api.result.result === 0 ? (
    <div className="result-table-wrap">
      <table className="result-table">
        <tbody>
          {(api.result.data?.['stdout'] as string[] ?? []).slice(1).map((item: string, idx: number) => (
            <tr key={idx}>
              <td>{idx + 1}</td>
              {item.split(';').filter(c => c !== '').map((col, i) => (
                <td key={i}>{col}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  ) : null;

const IdResult = ({ api, label }: { api: { result: BaseResponse | null }; label: string }) =>
  api.result && api.result.result === 0 ? (
    <div className="result-success-msg">
      ✓ {label}: {(api.result.data?.['stdout'] as string[] ?? []).slice(1).join('')}
    </div>
  ) : null;

/* ── Containers tab ────────────────────────────────── */
const ContainersTab: React.FC<{ host: string }> = ({ host }) => {
  const [containerName, setContainerName] = useState('');
  const [imageName, setImageName]         = useState('');
  const [ports, setPorts]                 = useState('');
  const [volumes, setVolumes]             = useState('');
  const [envs, setEnvs]                   = useState('');
  const [pjtName, setPjtName]             = useState('');
  const [gitUrl, setGitUrl]               = useState('');
  const [buildContext, setBuildContext]   = useState('');

  const psApi      = useApi<BaseResponse>((d: DockerPsRequest) => dockerService.ContainerPs(d));
  const runApi     = useApi<BaseResponse>((d: DockerRunRequest) => dockerService.RunContainer(d));
  const stopApi    = useApi<BaseResponse>((d: DockerStopRequest) => dockerService.StopContainer(d));
  const restartApi = useApi<BaseResponse>((d: DockerRestartRequest) => dockerService.RestartContainer(d));
  const removeApi  = useApi<BaseResponse>((d: DockerRemoveRequest) => dockerService.RemovceContainer(d));
  const buildApi   = useApi<BaseResponse>((d: DokcerBuildRequest) => dockerService.BuildImageAndPush(d));

  useEffect(() => {
    psApi.execute({ host });
  }, [host]); // eslint-disable-line react-hooks/exhaustive-deps

  const nameInput = (
    <input
      className="input"
      type="text"
      placeholder="Container Name"
      value={containerName}
      onChange={e => setContainerName(e.target.value)}
      style={{ width: '160px' }}
    />
  );

  return (
    <div>
      {/* Container list */}
      <section className="card">
        <h2 className="card-title">
          Container List
          <button
            className="btn btn-primary"
            style={{ marginLeft: '12px', padding: '3px 10px', fontSize: '0.7rem' }}
            onClick={() => psApi.execute({ host })}
            disabled={psApi.loading}
          >
            {psApi.loading ? '...' : '↻ Refresh'}
          </button>
        </h2>
        <ApiError api={psApi} />
        <ErrorBlock api={psApi} />
        <TableResult api={psApi} />
      </section>

      {/* Stop / Restart / Remove */}
      <section className="card">
        <h2 className="card-title">Container Operations</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <form className="form-row" onSubmit={async e => { e.preventDefault(); stopApi.reset(); await stopApi.execute({ host, name: containerName }); }}>
            {nameInput}
            <button className="btn btn-primary" type="submit" disabled={stopApi.loading}>
              {stopApi.loading ? '...' : 'Stop'}
            </button>
          </form>
          <ApiError api={stopApi} />
          <ErrorBlock api={stopApi} />
          <IdResult api={stopApi} label="Stopped" />

          <form className="form-row" onSubmit={async e => { e.preventDefault(); restartApi.reset(); await restartApi.execute({ host, name: containerName }); }}>
            {nameInput}
            <button className="btn btn-primary" type="submit" disabled={restartApi.loading}>
              {restartApi.loading ? '...' : 'Restart'}
            </button>
          </form>
          <ApiError api={restartApi} />
          <ErrorBlock api={restartApi} />
          <IdResult api={restartApi} label="Restarted" />

          <form className="form-row" onSubmit={async e => { e.preventDefault(); removeApi.reset(); await removeApi.execute({ host, name: containerName }); }}>
            {nameInput}
            <button className="btn btn-primary" type="submit" disabled={removeApi.loading}>
              {removeApi.loading ? '...' : 'Remove'}
            </button>
          </form>
          <ApiError api={removeApi} />
          <ErrorBlock api={removeApi} />
          <IdResult api={removeApi} label="Removed" />
        </div>
      </section>

      {/* Run container */}
      <section className="card">
        <h2 className="card-title">Run Container</h2>
        <form className="form-row" onSubmit={async e => {
          e.preventDefault();
          runApi.reset();
          await runApi.execute({
            host,
            image: imageName,
            name: containerName,
            port: ports.split(',').map(p => p.trim()).filter(Boolean),
            volume: volumes.split(',').map(v => v.trim()).filter(Boolean),
            env: envs.split(',').map(e => e.trim()).filter(Boolean),
          });
        }}>
          <input className="input" type="text" placeholder="Image Name" value={imageName} onChange={e => setImageName(e.target.value)} style={{ width: '150px' }} required />
          <input className="input" type="text" placeholder="Container Name" value={containerName} onChange={e => setContainerName(e.target.value)} style={{ width: '150px' }} required />
          <input className="input" type="text" placeholder="Ports (8080:80, ...)" value={ports} onChange={e => setPorts(e.target.value)} style={{ width: '180px' }} />
          <input className="input" type="text" placeholder="Volumes (/src:/dst, ...)" value={volumes} onChange={e => setVolumes(e.target.value)} style={{ width: '180px' }} />
          <input className="input" type="text" placeholder="Env (KEY=val, ...)" value={envs} onChange={e => setEnvs(e.target.value)} style={{ width: '160px' }} />
          <button className="btn btn-primary" type="submit" disabled={runApi.loading}>
            {runApi.loading ? 'Running...' : 'Run Container'}
          </button>
        </form>
        <ApiError api={runApi} />
        <ErrorBlock api={runApi} />
        <IdResult api={runApi} label="Container ID" />
      </section>

      {/* Build image */}
      <section className="card">
        <h2 className="card-title">Build Docker Image</h2>
        <form className="form-row" onSubmit={async e => {
          e.preventDefault();
          buildApi.reset();
          await buildApi.execute({ pjt_name: pjtName, url: gitUrl, context_path: buildContext });
        }}>
          <input className="input" type="text" placeholder="Project Name" value={pjtName} onChange={e => setPjtName(e.target.value)} style={{ width: '150px' }} required />
          <input className="input" type="text" placeholder="Git URL" value={gitUrl} onChange={e => setGitUrl(e.target.value)} style={{ width: '280px' }} required />
          <input className="input" type="text" placeholder="Dockerfile Context Path" value={buildContext} onChange={e => setBuildContext(e.target.value)} style={{ width: '190px' }} required />
          <button className="btn btn-primary" type="submit" disabled={buildApi.loading}>
            {buildApi.loading ? 'Building...' : 'Build Image'}
          </button>
        </form>
        <ApiError api={buildApi} />
        <ErrorBlock api={buildApi} />
        <TableResult api={buildApi} />
      </section>
    </div>
  );
};

/* ── Images tab ────────────────────────────────────── */
const ImagesTab: React.FC<{ host: string }> = ({ host }) => {
  const [imageName, setImageName] = useState('');

  const listApi   = useApi<BaseResponse>((d: DockerImageListRequest) => dockerService.GetImageList(d));
  const removeApi = useApi<BaseResponse>((d: DcokerImageRemoveRequest) => dockerService.ImageRemove(d));

  useEffect(() => {
    listApi.execute({ host });
  }, [host]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div>
      <section className="card">
        <h2 className="card-title">
          Image List
          <button
            className="btn btn-primary"
            style={{ marginLeft: '12px', padding: '3px 10px', fontSize: '0.7rem' }}
            onClick={() => listApi.execute({ host })}
            disabled={listApi.loading}
          >
            {listApi.loading ? '...' : '↻ Refresh'}
          </button>
        </h2>
        <ApiError api={listApi} />
        <ErrorBlock api={listApi} />
        <TableResult api={listApi} />
      </section>

      <section className="card">
        <h2 className="card-title">Remove Image</h2>
        <form className="form-row" onSubmit={async e => {
          e.preventDefault();
          removeApi.reset();
          await removeApi.execute({ host, name: imageName });
        }}>
          <input
            className="input"
            type="text"
            placeholder="Image Name"
            value={imageName}
            onChange={e => setImageName(e.target.value)}
            style={{ width: '220px' }}
            required
          />
          <button className="btn btn-primary" type="submit" disabled={removeApi.loading}>
            {removeApi.loading ? 'Removing...' : 'Remove Image'}
          </button>
        </form>
        <ApiError api={removeApi} />
        <ErrorBlock api={removeApi} />
        {removeApi.result && removeApi.result.result === 0 && (
          <div className="result-success-msg">✓ Image removed successfully</div>
        )}
      </section>
    </div>
  );
};

/* ── Logs tab ──────────────────────────────────────── */
const LogsTab: React.FC<{ host: string }> = ({ host }) => {
  const [containerName, setContainerName] = useState('');
  const [tail, setTail]                   = useState('100');

  const logsApi = useApi<BaseResponse>((d: DockerLogsRequest) => dockerService.GetContainerLogs(d));

  return (
    <div>
      <section className="card">
        <h2 className="card-title">Container Logs</h2>
        <form className="form-row" onSubmit={async e => {
          e.preventDefault();
          logsApi.reset();
          await logsApi.execute({ host, name: containerName, tail: parseInt(tail) || 100 });
        }}>
          <input
            className="input"
            type="text"
            placeholder="Container Name"
            value={containerName}
            onChange={e => setContainerName(e.target.value)}
            style={{ width: '200px' }}
            required
          />
          <input
            className="input"
            type="number"
            placeholder="Tail lines"
            value={tail}
            onChange={e => setTail(e.target.value)}
            style={{ width: '120px' }}
          />
          <button className="btn btn-primary" type="submit" disabled={logsApi.loading}>
            {logsApi.loading ? 'Loading...' : 'Get Logs'}
          </button>
        </form>
        <ApiError api={logsApi} />
        <ErrorBlock api={logsApi} />
        {logsApi.result && logsApi.result.result === 0 && (
          <div className="result-log">
            {(logsApi.result.data?.['stdout'] as string[] ?? []).slice(1).map((item: string, idx: number) => (
              <div key={idx} className="result-log-line">{item.replace(/;$/, '')}</div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
};

/* ── DaemonDetail page ─────────────────────────────── */
const DaemonDetail: React.FC = (): React.ReactElement => {
  const { host: encodedHost } = useParams<{ host: string }>();
  const host = decodeURIComponent(encodedHost ?? '');
  const [tab, setTab] = useState<Tab>('containers');

  return (
    <div>
      <Nav />
      <div className="page">
        <div style={{ marginBottom: '6px' }}>
          <Link to="/dashboard" className="back-link">← Daemons</Link>
        </div>
        <h1 className="page-title" style={{ fontFamily: 'monospace', fontSize: '1rem', letterSpacing: '0.01em' }}>
          {host}
        </h1>

        <div className="tab-bar">
          {(['containers', 'images', 'logs'] as Tab[]).map(t => (
            <button
              key={t}
              className={`tab-btn${tab === t ? ' tab-btn-active' : ''}`}
              onClick={() => setTab(t)}
            >
              {t.charAt(0).toUpperCase() + t.slice(1)}
            </button>
          ))}
        </div>

        {tab === 'containers' && <ContainersTab host={host} />}
        {tab === 'images'     && <ImagesTab host={host} />}
        {tab === 'logs'       && <LogsTab host={host} />}
      </div>
    </div>
  );
};

export default DaemonDetail;
