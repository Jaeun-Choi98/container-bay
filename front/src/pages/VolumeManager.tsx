import React, { useState, useRef } from 'react';
import { useApi } from '../api/hooks/useApi';
import { BaseResponse } from '../api/model/response';
import { dockerService } from '../service/dockerService';
import { fileService } from '../service/fileService';

import Nav from '../component/header';

const VolumeManager: React.FC = (): React.ReactElement => {

  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [path, setPath] = useState<string>('');
  const fileInputRef = useRef<HTMLInputElement>(null);


  const uploadFileApi = useApi<BaseResponse>((file: File, additionalData?: Record<string, string>): Promise<BaseResponse> => {
    return fileService.uploadFile(file, additionalData);
  });

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    if (e.target.files && e.target.files.length > 0) {
      setUploadFile(e.target.files[e.target.files.length - 1]);
    }
  };

  const handleUploadFile = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!uploadFile) {
      alert("파일을 선택해주세요.")
      return;
    }
    uploadFileApi.reset();
    return uploadFileApi.execute(uploadFile, { "path": path })
  };

  return (
    <div>
      <Nav></Nav>
      {/* 파일 업로드 */}
      <h2>Upload File</h2>
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        style={{ width: '300px', padding: '2px' }}
        required />
      {/*uploadFile && <p>선택된 파일: {uploadFile.name}</p>*/}
      <input
        type="text"
        value={path}
        placeholder='File Path(e.g. {env.volume-dir}/path.../file)'
        onChange={(e: React.ChangeEvent<HTMLInputElement>): void => {
          setPath(e.target.value);
        }}
        style={{ width: '300px', padding: '2px' }}
        required />
      <button onClick={handleUploadFile} disabled={!uploadFile || uploadFileApi.loading}>
        {uploadFileApi.loading ? 'Uploading...' : 'Upload'}
      </button>
    </div>
  );
}

export default VolumeManager;