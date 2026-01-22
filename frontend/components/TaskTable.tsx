import React, { useState, useEffect } from 'react';
import { Task } from '../types';
import {
  Folder,
  Trash2,
  Check,
  AlertTriangle
} from 'lucide-react';

interface TaskTableProps {
  tasks: Task[];
  onToggleStatus: (id: string) => void;
  onDelete: (id: string) => void;
}

const StatusIcon = ({ status }: { status: Task['status'] }) => {
  if (status === 'done') return <Check className="w-4 h-4 text-accent" />;
  return <div className="w-3 h-3 rounded-full border border-textMuted" />;
};

export const TaskTable: React.FC<TaskTableProps> = ({ tasks, onToggleStatus, onDelete }) => {
  const [deletingTaskId, setDeletingTaskId] = useState<string | null>(null);

  // Handle keyboard events for the confirmation modal
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!deletingTaskId) return;

      if (e.key === 'Enter') {
        e.preventDefault();
        onDelete(deletingTaskId);
        setDeletingTaskId(null);
      } else if (e.key === 'Escape') {
        e.preventDefault();
        setDeletingTaskId(null);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [deletingTaskId, onDelete]);

  const confirmDelete = () => {
    if (deletingTaskId) {
      onDelete(deletingTaskId);
      setDeletingTaskId(null);
    }
  };

  return (
    <>
      <div className="w-full overflow-x-auto">
        <div className="min-w-[900px] border-t border-border">
          {/* Table Header */}
          <div className="grid grid-cols-12 gap-4 px-4 py-3 bg-surfaceHighlight/50 text-xs font-medium text-textMuted border-b border-border">
            <div className="col-span-5 pl-8 flex items-center gap-1 cursor-pointer hover:text-textMain">
              NAME
            </div>
            <div className="col-span-2 text-center cursor-pointer hover:text-textMain">STATUS</div>
            <div className="col-span-2 text-center cursor-pointer hover:text-textMain">CREATED</div>
            <div className="col-span-2 text-center cursor-pointer hover:text-textMain">COMPLETED</div>
            <div className="col-span-1"></div>
          </div>

          {/* Table Body */}
          {tasks.length === 0 ? (
            <div className="py-12 flex flex-col items-center justify-center text-textMuted border-b border-border">
              <Folder className="w-12 h-12 mb-4 opacity-20" />
              <p>No tasks found</p>
            </div>
          ) : (
            tasks.map((task) => (
              <div
                key={task.id}
                className="group grid grid-cols-12 gap-4 px-4 py-3 border-b border-border hover:bg-surfaceHighlight/30 transition-colors items-center text-sm"
              >
                {/* Name Column */}
                <div className="col-span-5 flex items-center gap-3">
                  <input
                    type="checkbox"
                    checked={task.status === 'done'}
                    onChange={() => onToggleStatus(task.id)}
                    className="w-4 h-4 rounded border-border bg-background checked:bg-accent checked:border-accent focus:ring-1 focus:ring-accent transition-all cursor-pointer appearance-none border flex items-center justify-center after:content-['✓'] after:text-black after:text-xs after:hidden checked:after:block"
                  />
                  <span className={`truncate ${task.status === 'done' ? 'text-textMuted line-through' : 'text-textMain'}`}>
                    {task.title}
                  </span>
                </div>

                {/* Status Column */}
                <div className="col-span-2 flex justify-center">
                  <div
                    className={`flex items-center gap-2 px-2 py-1 rounded-full text-xs font-medium border
                      ${task.status === 'done' ? 'bg-accent/5 text-accent border-accent/20' :
                        'bg-white/5 text-textMuted border-white/10'}`}
                  >
                    <StatusIcon status={task.status} />
                    <span className="capitalize">{task.status}</span>
                  </div>
                </div>

                {/* Created Date Column */}
                <div className="col-span-2 flex justify-center text-textMuted text-xs font-mono">
                  {task.createdAt}
                </div>

                {/* Completed Date Column */}
                <div className="col-span-2 flex justify-center text-textMuted text-xs font-mono">
                  {task.completedAt || '--'}
                </div>

                {/* Actions Column */}
                <div className="col-span-1 flex justify-end opacity-0 group-hover:opacity-100 transition-opacity">
                  <button
                    onClick={() => setDeletingTaskId(task.id)}
                    className="p-1.5 rounded hover:bg-white/10 text-textMuted hover:text-red-400 transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Footer */}
        <div className="px-4 py-3 flex items-center border-b border-border text-xs text-textMuted">
          <Folder className="w-4 h-4 text-yellow-500 mr-2" />
          <span>{tasks.length} {tasks.length === 1 ? 'Item' : 'Items'}</span>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      {deletingTaskId && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm animate-fade-in">
          <div className="bg-surface border border-border rounded-lg shadow-2xl w-full max-w-sm p-6 relative">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center border border-red-500/20">
                <AlertTriangle className="w-5 h-5 text-red-400" />
              </div>
              <h3 className="text-lg font-semibold text-textMain">Delete Task</h3>
            </div>
            
            <p className="text-sm text-textMuted mb-6 leading-relaxed">
              Are you sure you want to delete this task? This action cannot be undone. 
              <br/>
              <span className="text-xs opacity-70 mt-2 block">Press <kbd className="font-mono bg-white/10 px-1 rounded text-textMain">Enter</kbd> to confirm</span>
            </p>

            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeletingTaskId(null)}
                className="px-4 py-2 rounded-md border border-border text-sm text-textMain hover:bg-surfaceHighlight transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={confirmDelete}
                className="px-4 py-2 rounded-md bg-red-500/10 border border-red-500/20 text-red-400 text-sm font-medium hover:bg-red-500/20 transition-colors flex items-center gap-2"
              >
                <Trash2 className="w-4 h-4" />
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};
