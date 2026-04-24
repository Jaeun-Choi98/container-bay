import React, { useState, useRef } from 'react';
import { useApi } from '../api/hooks/useApi';
import { BaseResponse } from '../api/types/response';
import { fileService } from '../services/fileService';
import Nav from '../components/Header/Header';

const VolumeManager: React.FC = (): React.ReactElement => {
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [path, setPath] = useState<string>('');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const uploadFileApi = useApi<BaseResponse>(
    (file: File, additionalData?: Record<string, string>) => fileService.uploadFile(file, additionalData)
  );

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    if (e.target.files && e.target.files.length > 0) {
      setUploadFile(e.target.files[e.target.files.length - 1]);
    }
  };

  const handleUploadFile = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!uploadFile) {
      alert('Please select a file.');
      return;
    }
    uploadFileApi.reset();
    await uploadFileApi.execute(uploadFile, { path });
  };

  return (
    <div>
      <Nav />
      <div className="page">
        <h1 className="page-title">Volume Directory</h1>

        <section className="card">
          <h2 className="card-title">Upload File</h2>
          <form className="form-row" onSubmit={handleUploadFile}>
            <input
              className="input"
              type="file"
              ref={fileInputRef}
              onChange={handleFileChange}
              required
            />
            <input
              className="input"
              type="text"
              value={path}
              placeholder="Destination path (e.g. /data/file.txt)"
              onChange={e => setPath(e.target.value)}
              style={{ width: '280px' }}
              required
            />
            <button className="btn btn-primary" type="submit" disabled={!uploadFile || uploadFileApi.loading}>
              {uploadFileApi.loading ? 'Uploading...' : 'Upload'}
            </button>
          </form>
          {uploadFileApi.error && <div className="api-error">{uploadFileApi.error}</div>}
          {uploadFileApi.result && uploadFileApi.result.result === 0 && (
            <div className="result-success-msg">✓ File uploaded successfully</div>
          )}
        </section>
      </div>
    </div>
  );
};

export default VolumeManager;
