// Type definitions for data returned by the backend API

export interface StatusResponse {
  staged: string[];
  modified: string[];
  untracked: string[];
  currentBranch: string;
  headCommit: string;
}

export interface Commit {
  id: string;
  message: string;
  timestamp: string;
  parents: string[];
}

export interface BranchList {
  branches: string[];
  current: string;
}

const API_BASE = '/api';

// Helper methods to make HTTP requests to the backend server
export const api = {
  // Fetch repository status (staged, modified, untracked files)
  getStatus: async (): Promise<StatusResponse> => {
    const res = await fetch(`${API_BASE}/status`);
    if (!res.ok) throw new Error(await res.text());
    return res.json();
  },

  // Stage files for commit
  addFiles: async (paths: string[]): Promise<void> => {
    const res = await fetch(`${API_BASE}/add`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ paths })
    });
    if (!res.ok) throw new Error(await res.text());
  },

  // Create a new commit with a message
  commit: async (message: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/commit`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message })
    });
    if (!res.ok) throw new Error(await res.text());
  },

  // Fetch full commit history
  getLog: async (): Promise<Commit[]> => {
    const res = await fetch(`${API_BASE}/log`);
    if (!res.ok) throw new Error(await res.text());
    return res.json();
  },

  // Get list of branches and the active branch
  getBranches: async (): Promise<BranchList> => {
    const res = await fetch(`${API_BASE}/branch`);
    if (!res.ok) throw new Error(await res.text());
    return res.json();
  },

  // Create a new branch
  createBranch: async (name: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/branch`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name })
    });
    if (!res.ok) throw new Error(await res.text());
  },

  // Switch to another branch or commit
  checkout: async (target: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/checkout`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target })
    });
    if (!res.ok) throw new Error(await res.text());
  },

  // Merge another branch into current branch
  merge: async (targetBranch: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/merge`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ targetBranch })
    });
    if (!res.ok) throw new Error(await res.text());
  },

  // Reset HEAD to an earlier commit (soft, mixed, or hard)
  reset: async (mode: string, target: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/reset`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ mode, target })
    });
    if (!res.ok) throw new Error(await res.text());
  }
};
