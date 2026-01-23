import { z } from 'zod';
import { api } from './api';
import { TodoItemDTOSchema, TodoItemDTOArraySchema, CreateTodoRequest } from './types';
import { Task } from '../types';
import { todoItemToTask } from './adapter';

const DeleteResponseSchema = z.string();

export const todoService = {
  async getAll(): Promise<Task[]> {
    const dtos = await api.get('', TodoItemDTOArraySchema);
    return dtos.map(todoItemToTask);
  },

  async create(secretPath: string): Promise<Task> {
    const request: CreateTodoRequest = { secretPath };
    const dto = await api.post('', TodoItemDTOSchema, request);
    return todoItemToTask(dto);
  },

  async toggleComplete(id: string): Promise<Task> {
    const dto = await api.patch(`/${id}`, TodoItemDTOSchema);
    return todoItemToTask(dto);
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/${id}`, DeleteResponseSchema);
  },
};
