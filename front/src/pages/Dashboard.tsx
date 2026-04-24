import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useApi } from '../api/hooks/useApi';
import { daemonService } from '../services/daemonService';
import { AddDaemonRequest } from '../api/types/request';
import { BaseResponse, Daemon } from '../api/types/response';
import Nav from '../components/Header/Header';
import BuildImageSidebar from '../components/BuildImageSidebar/BuildImageSidebar';

const Dashboard: React.FC = (): React.ReactElement => {
  const [host, setHost] = useState<string>('');
  const [label, setLabel] = useState<string>('');
  const navigate = useNavigate();

  const listApi   = useApi<BaseResponse>(() => daemonService.getDaemons());
  const addApi    = useApi<BaseResponse>((data: AddDaemonRequest) => daemonService.addDaemon(data));
  const removeApi = useApi<BaseResponse>((id: number) => daemonService.removeDaemon(id));

  useEffect(() => {
    listApi.execute();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const daemons: Daemon[] = listApi.result?.result === 0
    ? (listApi.result.data as Daemon[]) ?? []
    : [];

  const handleAdd = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    addApi.reset();
    await addApi.execute({ host, label });
    setHost('');
    setLabel('');
    listApi.reset();
    await listApi.execute();
  };

  const handleRemove = async (e: React.MouseEvent, id: number): Promise<void> => {
    e.stopPropagation();
    await removeApi.execute(id);
    listApi.reset();
    await listApi.execute();
  };

  return (
    <div>
      <Nav />
      <div className="page">
        <h1 className="page-title">Docker Daemons</h1>
        <div className="page-with-sidebar">
          <div className="page-main">
            {/* Add daemon */}
            <section className="card">
              <h2 className="card-title">Register Daemon</h2>
              <form className="form-row" onSubmit={handleAdd}>
                <input
                  className="input"
                  type="text"
                  placeholder="Host (e.g. 192.168.0.1:2375)"
                  value={host}
                  onChange={e => setHost(e.target.value)}
                  style={{ width: '260px' }}
                  required
                />
                <input
                  className="input"
                  type="text"
                  placeholder="Label (optional)"
                  value={label}
                  onChange={e => setLabel(e.target.value)}
                  style={{ width: '180px' }}
                />
                <button className="btn btn-primary" type="submit" disabled={addApi.loading}>
                  {addApi.loading ? 'Adding...' : 'Add Daemon'}
                </button>
              </form>
              {addApi.error && <div className="api-error">{addApi.error}</div>}
              {addApi.result && addApi.result.result !== 0 && (
                <div className="api-error">Failed to add daemon</div>
              )}
            </section>

            {/* Daemon list */}
            <section className="card">
              <h2 className="card-title">
                Registered Daemons
                <button
                  className="btn btn-primary"
                  style={{ marginLeft: '12px', padding: '3px 10px', fontSize: '0.7rem' }}
                  onClick={() => listApi.execute()}
                  disabled={listApi.loading}
                >
                  {listApi.loading ? '...' : '↻ Refresh'}
                </button>
              </h2>

              {listApi.error && <div className="api-error">{listApi.error}</div>}

              {daemons.length === 0 && !listApi.loading && (
                <div style={{ color: 'var(--text-muted)', fontSize: '0.875rem', padding: '12px 0' }}>
                  No daemons registered yet. Add one above.
                </div>
              )}

              <div className="daemon-list">
                {daemons.map((daemon: Daemon) => (
                  <div
                    key={daemon.id}
                    className="daemon-row"
                    onClick={() => navigate(`/daemon/${encodeURIComponent(daemon.host)}`)}
                  >
                    <div className="daemon-row-info">
                      <span className="daemon-row-icon">▣</span>
                      <div>
                        <div className="daemon-row-host">{daemon.host}</div>
                        {daemon.label && (
                          <div className="daemon-row-label">{daemon.label}</div>
                        )}
                      </div>
                    </div>
                    <div className="daemon-row-actions">
                      <span className="daemon-row-open">Open →</span>
                      <button
                        className="btn daemon-remove-btn"
                        onClick={e => handleRemove(e, daemon.id)}
                        disabled={removeApi.loading}
                      >
                        Remove
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </section>
          </div>

          <BuildImageSidebar />
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
