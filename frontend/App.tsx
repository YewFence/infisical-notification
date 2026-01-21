import React, { useState } from 'react';
import { Layout } from './components/Layout';
import { TaskTable } from './components/TaskTable';
import { Task } from './types';
import {
  Box,
  Search,
  Filter,
  Plus,
  ChevronDown,
} from 'lucide-react';

const INITIAL_TASKS: Task[] = [
  { 
    id: '1', 
    title: 'Mikan', 
    status: 'done', 
    priority: 'low', 
    tags: ['Dev'], 
    isFolder: true,
    createdAt: '2023-10-20',
    completedAt: '2023-10-24'
  },
  { 
    id: '2', 
    title: 'PacificYew', 
    status: 'done', 
    priority: 'high', 
    tags: ['Prod'], 
    isFolder: true,
    createdAt: '2023-10-25',
    completedAt: '2023-11-01'
  },
  { 
    id: '3', 
    title: 'Database Migration', 
    status: 'in-progress', 
    priority: 'high', 
    tags: ['Backend'], 
    isFolder: false,
    createdAt: '2023-11-10'
  },
  { 
    id: '4', 
    title: 'Update Documentation', 
    status: 'todo', 
    priority: 'medium', 
    tags: ['Docs'], 
    isFolder: false,
    createdAt: '2023-11-18'
  },
];

function App() {
  const [tasks, setTasks] = useState<Task[]>(INITIAL_TASKS);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterPriority, setFilterPriority] = useState<string | null>(null);

  const toggleStatus = (id: string) => {
    setTasks(prev => prev.map(t => {
      if (t.id !== id) return t;
      const nextStatus = t.status === 'todo' ? 'in-progress' : t.status === 'in-progress' ? 'done' : 'todo';
      
      const updates: Partial<Task> = { status: nextStatus };
      if (nextStatus === 'done') {
        updates.completedAt = new Date().toISOString().split('T')[0];
      } else {
        updates.completedAt = undefined;
      }

      return { ...t, ...updates };
    }));
  };

  const deleteTask = (id: string) => {
    setTasks(prev => prev.filter(t => t.id !== id));
  };

  const addTasks = (newTasks: Partial<Task>[]) => {
    const createdTasks: Task[] = newTasks.map(t => ({
      id: crypto.randomUUID(),
      title: t.title || 'Untitled',
      status: t.status || 'todo',
      priority: t.priority || 'medium',
      tags: t.tags || [],
      createdAt: new Date().toISOString().split('T')[0],
      isFolder: false
    }));
    setTasks(prev => [...prev, ...createdTasks]);
  };

  const filteredTasks = tasks.filter(task => {
    const matchesSearch = task.title.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesPriority = filterPriority ? task.priority === filterPriority : true;
    return matchesSearch && matchesPriority;
  });

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
             <button 
                onClick={() => setFilterPriority(prev => prev ? null : 'high')}
                className={`flex items-center gap-2 px-3 py-2 rounded bg-surfaceHighlight border ${filterPriority ? 'border-accent text-accent' : 'border-border text-textMuted hover:text-textMain'} transition-colors text-sm font-medium whitespace-nowrap`}
             >
               <Filter className="w-4 h-4" />
               Filters
             </button>
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
              onClick={() => addTasks([{ title: 'New Secret Task', priority: 'medium' }])}
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
