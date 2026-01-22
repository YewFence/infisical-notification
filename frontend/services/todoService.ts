import { api } from './api';
import { TodoItemDTO, CreateTodoRequest } from './types';
import { Task } from '../types';
import { todoItemToTask } from './adapter';

export const todoService = {
  async getAll(): Promise<Task[]> {
    const dtos = await api.get<TodoItemDTO[]>('');
    return dtos.map(todoItemToTask);
  },

  async create(secretPath: string): Promise<Task> {
    const request: CreateTodoRequest = { secretPath };
    const dto = await api.post<TodoItemDTO>('', request);
    return todoItemToTask(dto);
  },

  async toggleComplete(id: string): Promise<Task> {
    const dto = await api.patch<TodoItemDTO>(`/${id}`);
    return todoItemToTask(dto);
  },

  async delete(id: string): Promise<void> {
    await api.delete<string>(`/${id}`);
  },
};
