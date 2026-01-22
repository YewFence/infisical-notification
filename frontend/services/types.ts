export interface TodoItemDTO {
  id: number;
  secretPath: string;
  isCompleted: boolean;
  createdAt: string;
  completedAt: string | null;
}

export interface CreateTodoRequest {
  secretPath: string;
}
