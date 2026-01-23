import { z } from 'zod';

// ============================================
// Zod Schemas（运行时验证）
// ============================================

export const TodoItemDTOSchema = z.object({
  id: z.number(),
  secretPath: z.string(),
  isCompleted: z.boolean(),
  createdAt: z.string(),
  completedAt: z.string().nullable(),
});

export const TodoItemDTOArraySchema = z.array(TodoItemDTOSchema);

export const CreateTodoRequestSchema = z.object({
  secretPath: z.string(),
});

// 通用 API 响应包装器
export const createApiResponseSchema = <T extends z.ZodType>(dataSchema: T) =>
  z.object({ data: dataSchema });

export const ApiErrorSchema = z.object({
  error: z.string(),
});

// ============================================
// TypeScript 类型（从 Schema 推断）
// ============================================

export type TodoItemDTO = z.infer<typeof TodoItemDTOSchema>;
export type CreateTodoRequest = z.infer<typeof CreateTodoRequestSchema>;
