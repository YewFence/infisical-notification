export interface Task {
  id: string;
  title: string;
  status: 'todo' | 'done';
  createdAt: string;
  completedAt?: string;
}

export type SortField = 'title' | 'status' | 'createdAt' | 'completedAt';
export type SortOrder = 'asc' | 'desc';

export interface FilterState {
  search: string;
  status: string | null;
}
