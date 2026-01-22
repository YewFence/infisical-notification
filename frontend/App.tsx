import React, { useState, useEffect } from 'react';
import { Layout } from './components/Layout';
import { TaskTable } from './components/TaskTable';
import { Task } from './types';
import { todoService } from './services/todoService';
import { ApiException } from './services/api';
import {
  Box,
  Search,
  Filter,
  Plus,
  ChevronDown,
} from 'lucide-react';

function App() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  // const [filterPriority, setFilterPriority] = useState<string | null>(null);

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

    // 简化状态切换: todo ↔ done
    const nextStatus = task.status === 'todo' ? 'done' : 'todo';

    // 原逻辑（包含 in-progress）:
    // const nextStatus = task.status === 'todo' ? 'in-progress' :
    //                    task.status === 'in-progress' ? 'done' : 'todo';

    // 前端乐观更新
    setTasks(prev => prev.map(t =>
      t.id === id ? { ...t, status: nextStatus } : t
    ));

    try {
      // 直接调用后端切换完成状态
      const updated = await todoService.toggleComplete(id);
      setTasks(prev => prev.map(t => t.id === id ? updated : t));

      // 原逻辑（仅在 done 状态时调用后端）:
      // if (nextStatus === 'done' || task.status === 'done') {
      //   const updated = await todoService.toggleComplete(id);
      //   setTasks(prev => prev.map(t => t.id === id ? updated : t));
      // }
      // 'in-progress' 状态仅在前端维护
    } catch (err) {
      // 回滚
      setTasks(prev => prev.map(t => t.id === id ? task : t));
      if (err instanceof ApiException) {
        alert(`操作失败: ${err.message}`);
      }
    }
  };

  const deleteTask = async (id: string) => {
    const taskToDelete = tasks.find(t => t.id === id);
    if (!taskToDelete) return;

    // 乐观删除
    setTasks(prev => prev.filter(t => t.id !== id));

    try {
      await todoService.delete(id);
    } catch (err) {
      // 回滚
      setTasks(prev => [...prev, taskToDelete]);
      if (err instanceof ApiException) {
        alert(`删除失败: ${err.message}`);
      }
    }
  };

  const addTasks = async (newTasks: Partial<Task>[]) => {
    for (const taskPartial of newTasks) {
      const secretPath = taskPartial.title || '/new/secret/path';

      try {
        const created = await todoService.create(secretPath);
        setTasks(prev => [...prev, created]);
      } catch (err) {
        if (err instanceof ApiException) {
          alert(`创建失败: ${err.message}`);
        }
      }
    }
  };

  const filteredTasks = tasks.filter(task => {
    const matchesSearch = task.title.toLowerCase().includes(searchQuery.toLowerCase());
    // const matchesPriority = filterPriority ? task.priority === filterPriority : true;
    return matchesSearch; // && matchesPriority;
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
      <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-4">

        {/* Left Side: Folder Button */}
        <div className="flex items-center gap-2">
            <button className="w-9 h-9 flex items-center justify-center bg-surfaceHighlight border border-border rounded hover:border-accent transition-colors group">
              <div className="w-4 h-4 text-accent group-hover:scale-110 transition-transform">
                <svg viewBox="0 0 24 24" fill="currentColor" className="w-full h-full"><path d="M20 6h-8l-2-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm0 12H4V8h16v10z"/></svg>
              </div>
            </button>
        </div>

        {/* Right Side: Filters & Actions */}
        <div className="flex flex-col md:flex-row items-center gap-3 w-full md:w-auto">

          <div className="flex items-center gap-2 w-full md:w-auto">
             {/* <button
                onClick={() => setFilterPriority(prev => prev ? null : 'high')}
                className={`flex items-center gap-2 px-3 py-2 rounded bg-surfaceHighlight border ${filterPriority ? 'border-accent text-accent' : 'border-border text-textMuted hover:text-textMain'} transition-colors text-sm font-medium whitespace-nowrap`}
             >
               <Filter className="w-4 h-4" />
               Filters
             </button> */}
          </div>

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
              onClick={() => addTasks([{ title: '/new/secret/path' }])}
              className="flex-1 md:flex-none flex items-center justify-between gap-3 px-3 py-2 bg-surfaceHighlight border border-border hover:border-textMuted rounded text-textMain text-sm font-medium transition-colors group min-w-[130px]"
            >
              <div className="flex items-center gap-2">
                <Plus className="w-4 h-4" />
                Add Secret
              </div>
              <ChevronDown className="w-3 h-3 text-textMuted group-hover:text-textMain" />
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
    </Layout>
  );
}

export default App;
