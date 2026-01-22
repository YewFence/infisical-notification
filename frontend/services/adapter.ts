import { Task } from '../types';
import { TodoItemDTO } from './types';

export function todoItemToTask(dto: TodoItemDTO): Task {
  const formatDate = (isoDate: string) => isoDate.split('T')[0];

  return {
    id: dto.id.toString(),
    title: dto.secretPath,
    status: dto.isCompleted ? 'done' : 'todo',
    createdAt: formatDate(dto.createdAt),
    completedAt: dto.completedAt ? formatDate(dto.completedAt) : undefined,
  };
}
