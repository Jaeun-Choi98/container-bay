import React, { useState } from 'react';
import { useApi } from '../api/hooks/useApi';
import { BaseResponse } from '../api/model/response';
import { DockerImageListRequest } from '../api/model/request';
import { dockerService } from '../service/dockerService';

const ImageManager: React.FC = (): React.ReactElement => {
  const [host, setHost] = useState<string>('');

  const GetImgListApi = useApi<BaseResponse>((data: DockerImageListRequest): Promise<BaseResponse> => {
    return dockerService.GetImageList(data);
  })

  const handleGetImageList = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    GetImgListApi.reset();
    return GetImgListApi.execute({
      host: host,
    })
  }

  return (
    <div>
      <h1>Image List</h1>
      {/* 이미지 리스트 */}
      <section>
        <input
          type='text'
          placeholder='host'
          value={host}
          onChange={(e: React.ChangeEvent<HTMLInputElement>): void =>
            setHost(e.target.value)
          }
        />
        <button onClick={handleGetImageList} disabled={GetImgListApi.loading}>
          {GetImgListApi.loading ? 'Loading...' : 'Get Image List'}
        </button>
        {GetImgListApi.error && <div style={{ color: 'red' }}>Error: {GetImgListApi.error}</div>}
        {GetImgListApi.result?.result !== 0 &&
          <pre style={{ color: 'red' }}>
            {GetImgListApi.result?.data['stderr'].slice(1).map((item: string, index: number): React.ReactElement => {
              const colums = item.split(';');
              return (
                <ul style={{ margin: 'auto', }}>
                  { }
                </ul>
              )
            })}
          </pre>
        }
      </section>
    </div>
  )
}