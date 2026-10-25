import { useEffect, useState } from 'react';
import { api, type StatusResponse, type Commit, type BranchList } from './api';
import { GitBranch, GitCommit, FileText, Plus, Check } from 'lucide-react';
import './index.css';

export default function App() {
  const [status, setStatus] = useState<StatusResponse | null>(null);
  const [branches, setBranches] = useState<BranchList | null>(null);
  const [logs, setLogs] = useState<Commit[]>([]);
  const [commitMsg, setCommitMsg] = useState('');
  const [newBranch, setNewBranch] = useState('');

  // Fetch status, branches, and commit logs from the backend
  const loadData = async () => {
    try {
      const [s, b, l] = await Promise.all([
        api.getStatus(),
        api.getBranches(),
        api.getLog()
      ]);
      setStatus(s);
      setBranches(b);
      setLogs(l || []);
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    loadData();
    const timer = setInterval(loadData, 3000);
    return () => clearInterval(timer);
  }, []);

  // Stage all files
  const handleStageAll = async () => {
    await api.addFiles(['.']);
    loadData();
  };

  // Create commit
  const handleCommit = async () => {
    if (!commitMsg.trim()) return;
    await api.commit(commitMsg);
    setCommitMsg('');
    loadData();
  };

  // Create branch
  const handleCreateBranch = async () => {
    if (!newBranch.trim()) return;
    await api.createBranch(newBranch);
    setNewBranch('');
    loadData();
  };

  return (
    <div className="app-container">
      {/* Header */}
      <header className="app-header">
        <div className="app-title">
          <GitBranch size={16} /> mgit Source Control
        </div>
      </header>

      <div className="workspace">
        {/* Left Panel: Files and Commit */}
        <div className="panel-left">
          <div className="section-header">
            <span>Staged Changes</span>
            <button className="icon-btn" onClick={handleStageAll} title="Stage All">
              <Plus size={14} />
            </button>
          </div>

          <div className="scrollable-content" style={{ padding: '8px 0' }}>
            {status?.staged?.map(f => (
              <div key={f} className="file-item">
                <span className="file-name"><FileText size={12} style={{ marginRight: 6 }} />{f}</span>
                <span className="file-status status-S">S</span>
              </div>
            ))}
            {status?.untracked?.map(f => (
              <div key={f} className="file-item">
                <span className="file-name"><FileText size={12} style={{ marginRight: 6 }} />{f}</span>
                <span className="file-status status-U">U</span>
              </div>
            ))}
            {(!status?.staged?.length && !status?.untracked?.length) && (
              <div style={{ padding: '24px', textAlign: 'center', color: 'var(--text-muted)' }}>
                Working tree clean
              </div>
            )}
          </div>

          <div className="commit-box">
            <textarea
              className="input-field"
              placeholder="Commit message..."
              value={commitMsg}
              onChange={e => setCommitMsg(e.target.value)}
              rows={3}
            />
            <button
              className="btn"
              style={{ marginTop: 8 }}
              onClick={handleCommit}
              disabled={!commitMsg.trim()}
            >
              <Check size={14} /> Commit
            </button>
          </div>
        </div>

        {/* Right Panel: Branches and Logs */}
        <div className="panel-right">
          {/* Branch bar */}
          <div className="branches-bar">
            {branches?.branches.map(b => (
              <div
                key={b}
                className={`branch-tag ${branches.current === b ? 'active' : ''}`}
                onClick={() => api.checkout(b).then(loadData)}
                style={{ cursor: 'pointer' }}
              >
                <GitBranch size={12} /> {b}
              </div>
            ))}

            <div style={{ marginLeft: 'auto', display: 'flex', gap: 8 }}>
              <input
                className="input-field"
                style={{ width: 140 }}
                placeholder="New branch..."
                value={newBranch}
                onChange={e => setNewBranch(e.target.value)}
              />
              <button className="btn btn-secondary" onClick={handleCreateBranch}>
                Create
              </button>
            </div>
          </div>

          <div className="section-header" style={{ borderBottom: '1px solid var(--border-color)' }}>
            <span>Commit History</span>
          </div>

          {/* History table */}
          <div className="scrollable-content">
            <table className="history-table">
              <tbody>
                {logs.map(log => (
                  <tr key={log.id} className="history-row">
                    <td className="history-cell hash-cell">
                      <GitCommit size={12} style={{ marginRight: 6, verticalAlign: 'middle' }} />
                      {log.id.substring(0, 7)}
                    </td>
                    <td className="history-cell msg-cell">{log.message}</td>
                    <td className="history-cell date-cell">
                      {new Date(log.timestamp).toLocaleTimeString()}
                    </td>
                  </tr>
                ))}
                {logs.length === 0 && (
                  <tr>
                    <td colSpan={3} className="history-cell" style={{ textAlign: 'center', color: 'var(--text-muted)', padding: 32 }}>
                      No commits yet
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
