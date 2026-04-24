import React, { useState } from 'react';
import { useApi } from '../../api/hooks/useApi';
import { dockerService } from '../../services/dockerService';
import { DokcerBuildRequest } from '../../api/types/request';
import { BaseResponse } from '../../api/types/response';

const BuildImageSidebar: React.FC = () => {
  const [pjtName, setPjtName]         = useState('');
  const [gitUrl, setGitUrl]           = useState('');
  const [buildContext, setBuildContext] = useState('');

  const buildApi = useApi<BaseResponse>((d: DokcerBuildRequest) => dockerService.BuildImageAndPush(d));

  return (
    <aside className="sidebar">
      <section className="card">
        <h2 className="card-title">Build Docker Image</h2>
        <form
          className="sidebar-form"
          onSubmit={async e => {
            e.preventDefault();
            buildApi.reset();
            await buildApi.execute({ pjt_name: pjtName, url: gitUrl, context_path: buildContext });
          }}
        >
          <input
            className="input sidebar-input"
            type="text"
            placeholder="Project Name"
            value={pjtName}
            onChange={e => setPjtName(e.target.value)}
            required
          />
          <input
            className="input sidebar-input"
            type="text"
            placeholder="Git URL"
            value={gitUrl}
            onChange={e => setGitUrl(e.target.value)}
            required
          />
          <input
            className="input sidebar-input"
            type="text"
            placeholder="Dockerfile Context Path"
            value={buildContext}
            onChange={e => setBuildContext(e.target.value)}
            required
          />
          <button className="btn btn-primary sidebar-btn" type="submit" disabled={buildApi.loading}>
            {buildApi.loading ? 'Building...' : 'Build Image'}
          </button>
        </form>

        {buildApi.error && <div className="api-error">{buildApi.error}</div>}

        {buildApi.result && buildApi.result.result !== 0 && (
          <div className="result-error-block">
            {(buildApi.result.data?.['stderr'] as string[] ?? []).slice(1)
              .flatMap((item: string, idx: number) =>
                item.split(';').filter(c => c !== '').map((col, i) => (
                  <div key={`${idx}-${i}`} className="result-error-line">{col}</div>
                ))
              )}
          </div>
        )}

        {buildApi.result && buildApi.result.result === 0 && (
          <div className="result-success-msg">✓ Build complete</div>
        )}
      </section>
    </aside>
  );
};

export default BuildImageSidebar;
