export interface Task {
  id: string;
  title: string;
  description?: string;
  status: 'todo' | /* 'in-progress' | */ 'done';
  // priority: 'low' | 'medium' | 'high';
  createdAt: string;
  completedAt?: string;
  assignee?: string;
  children?: Task[];
}

export type SortField = 'title' | 'status' | /* 'priority' | */ 'createdAt' | 'completedAt';
export type SortOrder = 'asc' | 'desc';

export interface FilterState {
  search: string;
  status: string | null;
  // priority: string | null;
}
