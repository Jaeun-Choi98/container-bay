import React, { useState } from 'react';
import { useApi } from '../api/hooks/useApi';
import { BaseResponse } from '../api/types/response';
import { DockerImageListRequest, DcokerImageRemoveRequest } from '../api/types/request';
import { dockerService } from '../services/dockerService';
import Nav from '../components/Header/Header';

const ImageManager: React.FC = (): React.ReactElement => {
  const [host, setHost] = useState<string>('');
  const [name, setName] = useState<string>('');

  const getImgListApi = useApi<BaseResponse>((data: DockerImageListRequest) => dockerService.GetImageList(data));
  const imgRemoveApi  = useApi<BaseResponse>((data: DcokerImageRemoveRequest) => dockerService.ImageRemove(data));

  const handleGetImageList = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    getImgListApi.reset();
    await getImgListApi.execute({ host });
  };

  const handleImageRemove = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    imgRemoveApi.reset();
    await imgRemoveApi.execute({ host, name });
  };

  const renderErrorBlock = (api: typeof getImgListApi) =>
    api.result && api.result.result !== 0 ? (
      <div className="result-error-block">
        {api.result.data['stderr'].slice(1).map((item: string, idx: number) =>
          item.split(';').filter(col => col !== '').map((col, i) => (
            <div key={`${idx}-${i}`} className="result-error-line">{col}</div>
          ))
        )}
      </div>
    ) : null;

  const renderTableResult = (api: typeof getImgListApi) =>
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

  return (
    <div>
      <Nav />
      <div className="page">
        <h1 className="page-title">Image Management</h1>

        {/* Image List */}
        <section className="card">
          <h2 className="card-title">Image List</h2>
          <form className="form-row" onSubmit={handleGetImageList}>
            <input
              className="input"
              type="text"
              placeholder="Host Docker Daemon (e.g. 192.168.0.1:2375)"
              value={host}
              onChange={e => setHost(e.target.value)}
              style={{ width: '300px' }}
              required
            />
            <button className="btn btn-primary" type="submit" disabled={getImgListApi.loading}>
              {getImgListApi.loading ? 'Loading...' : 'Get Image List'}
            </button>
          </form>
          {getImgListApi.error && <div className="api-error">{getImgListApi.error}</div>}
          {renderErrorBlock(getImgListApi)}
          {renderTableResult(getImgListApi)}
        </section>

        {/* Image Remove */}
        <section className="card">
          <h2 className="card-title">Remove Image</h2>
          <form className="form-row" onSubmit={handleImageRemove}>
            <input
              className="input"
              type="text"
              placeholder="Host Docker Daemon (e.g. 192.168.0.1:2375)"
              value={host}
              onChange={e => setHost(e.target.value)}
              style={{ width: '300px' }}
              required
            />
            <input
              className="input"
              type="text"
              placeholder="Image Name"
              value={name}
              onChange={e => setName(e.target.value)}
              style={{ width: '200px' }}
              required
            />
            <button className="btn btn-primary" type="submit" disabled={imgRemoveApi.loading}>
              {imgRemoveApi.loading ? 'Removing...' : 'Remove Image'}
            </button>
          </form>
          {imgRemoveApi.error && <div className="api-error">{imgRemoveApi.error}</div>}
          {renderErrorBlock(imgRemoveApi)}
          {imgRemoveApi.result && imgRemoveApi.result.result === 0 && (
            <div className="result-success-msg">✓ Image removed successfully</div>
          )}
        </section>
      </div>
    </div>
  );
};

export default ImageManager;
