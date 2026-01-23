import React, { useState, useEffect, useCallback, useRef } from 'react';
import { Layout } from './components/Layout';
import { TaskTable } from './components/TaskTable';
import { CreateTaskModal } from './components/CreateTaskModal';
import { Task } from './types';
import { todoService } from './services/todoService';
import { ApiException } from './services/api';
import { usePolling } from './hooks/usePolling';
import {
  Box,
  Search,
  Plus,
} from 'lucide-react';

// 从环境变量读取轮询间隔，默认 30 秒
const POLL_INTERVAL_SECONDS = parseInt(import.meta.env.VITE_POLL_INTERVAL_SECONDS || '30', 10);
const POLL_INTERVAL_MS = POLL_INTERVAL_SECONDS * 1000;

function App() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);

  // 用于追踪是否有用户操作正在进行
  const isUserActionRef = useRef(false);

  // 静默加载任务（用于轮询，不显示 loading 状态）
  const silentLoadTasks = useCallback(async () => {
    // 如果用户正在操作，跳过本次轮询
    if (isUserActionRef.current) {
      return;
    }

    try {
      const fetchedTasks = await todoService.getAll();
      setTasks(fetchedTasks);
      // 轮询成功时清除之前的错误状态
      setError(null);
    } catch (err) {
      // 轮询失败时静默处理，不更新错误状态
      console.warn('[Polling] 获取任务列表失败:', err);
    }
  }, []);

  // 使用轮询 Hook
  const { pause: pausePolling, resume: resumePolling } = usePolling(
    silentLoadTasks,
    {
      interval: POLL_INTERVAL_MS,
      enabled: !loading && !error, // 初始加载完成且无错误时才启用轮询
    }
  );

  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    try {
      setLoading(true);
      setError(null);
      const fetchedTasks = await todoService.getAll();
      setTasks(fetchedTasks);
    } catch (err) {
      if (err instanceof ApiException) {
        setError(`加载失败: ${err.message}`);
      } else {
        setError('网络错误，请检查后端服务是否启动');
      }
    } finally {
      setLoading(false);
    }
  };

  const toggleStatus = async (id: string) => {
    const task = tasks.find(t => t.id === id);
    if (!task) return;

    // 暂停轮询
    isUserActionRef.current = true;
    pausePolling();

    // 简化状态切换: todo ↔ done
    const nextStatus = task.status === 'todo' ? 'done' : 'todo';

    // 前端乐观更新
    setTasks(prev => prev.map(t =>
      t.id === id ? { ...t, status: nextStatus } : t
    ));

    try {
      // 直接调用后端切换完成状态
      const updated = await todoService.toggleComplete(id);
      setTasks(prev => prev.map(t => t.id === id ? updated : t));
    } catch (err) {
      // 回滚
      setTasks(prev => prev.map(t => t.id === id ? task : t));
      if (err instanceof ApiException) {
        alert(`操作失败: ${err.message}`);
      }
    } finally {
      // 恢复轮询
      isUserActionRef.current = false;
      resumePolling();
    }
  };

  const deleteTask = async (id: string) => {
    const taskIndex = tasks.findIndex(t => t.id === id);
    if (taskIndex === -1) return;

    const taskToDelete = tasks[taskIndex];

    // 暂停轮询
    isUserActionRef.current = true;
    pausePolling();

    // 乐观删除
    setTasks(prev => prev.filter(t => t.id !== id));

    try {
      await todoService.delete(id);
    } catch (err) {
      // 回滚到原来的位置
      setTasks(prev => {
        const newTasks = [...prev];
        newTasks.splice(taskIndex, 0, taskToDelete);
        return newTasks;
      });
      if (err instanceof ApiException) {
        alert(`删除失败: ${err.message}`);
      }
    } finally {
      // 恢复轮询
      isUserActionRef.current = false;
      resumePolling();
    }
  };

  const addTask = async (secretPath: string) => {
    // 暂停轮询
    isUserActionRef.current = true;
    pausePolling();

    try {
      const created = await todoService.create(secretPath);
      setTasks(prev => [...prev, created]);
    } catch (err) {
      if (err instanceof ApiException) {
        alert(`创建失败: ${err.message}`);
      }
    } finally {
      // 恢复轮询
      isUserActionRef.current = false;
      resumePolling();
    }
  };

  const filteredTasks = tasks.filter(task => {
    return task.title.toLowerCase().includes(searchQuery.toLowerCase());
  });

  if (loading) {
    return (
      <Layout>
        <div className="flex items-center justify-center h-64">
          <div className="text-textMuted">加载中...</div>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className="flex flex-col items-center justify-center h-64">
          <div className="text-red-400 mb-4">{error}</div>
          <button
            onClick={loadTasks}
            className="px-4 py-2 bg-accent text-black rounded hover:bg-accent/80"
          >
            重试
          </button>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      {/* Header Section */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
           <Box className="w-8 h-8 text-accent" strokeWidth={1.5} />
           <h1 className="text-2xl font-bold text-textMain tracking-tight hover:underline decoration-accent underline-offset-4 cursor-pointer">
             Project Tasks
           </h1>
        </div>
        <p className="text-textMuted text-sm max-w-2xl">
          管理你的待办事项,跟踪任务进度和优先级。
        </p>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col md:flex-row items-start md:items-center justify-end gap-4 mb-4">

        {/* Filters & Actions */}
        <div className="flex flex-col md:flex-row items-center gap-3 w-full md:w-auto">
          <div className="relative group w-full md:w-64">
             <Search className="absolute left-3 top-2.5 w-4 h-4 text-textMuted group-focus-within:text-accent transition-colors" />
             <input
               type="text"
               placeholder="Search by task or folder name..."
               value={searchQuery}
               onChange={(e) => setSearchQuery(e.target.value)}
               className="w-full bg-surfaceHighlight border border-border rounded pl-9 pr-4 py-2 text-sm text-textMain placeholder:text-textMuted/50 focus:outline-none focus:border-textMuted/50 focus:ring-1 focus:ring-border transition-all"
             />
          </div>

          <div className="flex items-center gap-2 w-full md:w-auto">
            <button
              onClick={() => setIsModalOpen(true)}
              className="flex-1 md:flex-none flex items-center gap-2 px-3 py-2 bg-surfaceHighlight border border-border hover:border-textMuted rounded text-textMain text-sm font-medium transition-colors"
            >
              <Plus className="w-4 h-4" />
              Add Task
            </button>
          </div>
        </div>
      </div>

      {/* Main Table Container */}
      <div className="bg-surface border border-border rounded-lg overflow-hidden shadow-sm">
        <TaskTable
          tasks={filteredTasks}
          onToggleStatus={toggleStatus}
          onDelete={deleteTask}
        />
      </div>

      {/* Create Task Modal */}
      <CreateTaskModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onConfirm={addTask}
      />
    </Layout>
  );
}

export default App;
