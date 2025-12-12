import React, { useState } from 'react';
import { useApi } from '../api/hooks/useApi';
import { BaseResponse } from '../api/model/response';
import { DockerImageListRequest, DcokerImageRemoveRequest } from '../api/model/request';
import { dockerService } from '../service/dockerService';
import Nav from '../component/header';

const ImageManager: React.FC = (): React.ReactElement => {
  const [host, setHost] = useState<string>('');
  const [name, setName] = useState<string>('');

  const getImgListApi = useApi<BaseResponse>((data: DockerImageListRequest): Promise<BaseResponse> => {
    return dockerService.GetImageList(data);
  })

  const imgRemoveApi = useApi<BaseResponse>((data: DcokerImageRemoveRequest): Promise<BaseResponse> => {
    return dockerService.ImageRemove(data);
  })

  const handleGetImageList = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    getImgListApi.reset();
    return getImgListApi.execute({
      host: host,
    })
  }

  const handleImageRemove = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    imgRemoveApi.reset();
    return imgRemoveApi.execute(
      {
        host: host,
        name: name,
      }
    );
  }

  return (
    <div>
      <Nav></Nav>
      {/* 이미지 리스트 */}
      <h2>Image List</h2>
      <section>
        <input
          type='text'
          placeholder='Host Docker Daemon(e.g. 192.168.0.1:2375)'
          value={host}
          onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
            setHost(e.target.value)
          }
          style={{ width: '300px', padding: '2px' }}
          required
        />
        <button onClick={handleGetImageList} disabled={getImgListApi.loading}>
          {getImgListApi.loading ? 'Loading...' : 'Get Image List'}
        </button>
        {getImgListApi.error && <div style={{ color: 'red' }}>Error: {getImgListApi.error}</div>}
        {getImgListApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {getImgListApi.result?.data['stderr'].slice(1).map((item: string, index: number): React.ReactElement => {
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
              )
            })}
          </pre>
        }
        {getImgListApi.result?.result === 0 &&
          <pre style={{ color: 'green' }}>
            {getImgListApi.result?.data["stdout"].slice(1).map((item: string, index: number): React.ReactElement => {
              const colums = item.split(';');
              return (
                <tr key={index} style={{ tableLayout: 'fixed', width: '100%' }}>
                  <td style={{ width: '5%', border: '1px solid #ccc', padding: '5px' }}>{index !== 0 ? index : ""}</td>
                  {
                    colums.map((col: string, i: number): React.ReactElement | null => {
                      if (col === "") return null;
                      return (
                        <td key={i} style={{ width: '5%', border: '1px solid #ccc', padding: '5px' }}>{col} </td>
                      );
                    })
                  }
                </tr>
              );
            })}
          </pre>
        }
      </section>

      {/*이미지 삭제*/}
      <section>
        <h2>Image Remove</h2>
        <input
          type='text'
          placeholder='Host Docker Daemon(e.g. 192.168.0.1:2375)'
          value={host}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setHost(e.target.value);
          }}
          style={{ width: '300px', padding: '2px' }}
          required
        />
        <input
          type='text'
          placeholder='Image name'
          value={name}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setName(e.target.value);
          }}
          style={{ width: '200px', padding: '2px' }}
          required
        />
        <button onClick={handleImageRemove} disabled={imgRemoveApi.loading}>
          {imgRemoveApi.loading ? 'Removing...' : 'Remove Image'}
        </button>
        {imgRemoveApi.error && <div style={{ color: 'red' }}>Error: {imgRemoveApi.error}</div>}
        {imgRemoveApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {imgRemoveApi.result?.data['stderr'].slice(1).map((item: string, index: number): React.ReactElement => {
              const colums = item.split(';');
              return (
                <ul style={{ margin: 'auto', width: '50%' }}>
                  {colums.map((col: string, i: number): React.ReactElement | null => {
                    if (col === "") return null
                    return (
                      <li key={i} style={{ whiteSpace: 'pre-wrap', width: '100%', padding: '5px', listStyle: 'none' }}>{col}</li>
                    );
                  })}
                </ul>
              )
            })}
          </pre>
        }
        {imgRemoveApi.result?.result === 0 &&
          <pre style={{ color: 'green' }}>
            {imgRemoveApi.result?.data["stdout"].slice(1).map((item: string, index: number): React.ReactElement => {
              const colums = item.split(';');
              return (
                <tr key={index} style={{ tableLayout: 'fixed', width: '100%' }}>
                  <td style={{ width: '5%', border: '1px solid #ccc', padding: '5px' }}>{index !== 0 ? index : ""}</td>
                  {
                    colums.map((col: string, i: number): React.ReactElement | null => {
                      if (col === "") return null;
                      return (
                        <td key={i} style={{ width: '5%', border: '1px solid #ccc', padding: '5px' }}>{col} </td>
                      );
                    })
                  }
                </tr>
              );
            })}
          </pre>
        }
      </section>
    </div>
  )
}

export default ImageManager;