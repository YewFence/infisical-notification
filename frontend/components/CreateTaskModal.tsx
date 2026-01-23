import React, { useState, useEffect } from 'react';
import { X } from 'lucide-react';

interface CreateTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (path: string) => void;
}

export const CreateTaskModal: React.FC<CreateTaskModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
}) => {
  const [path, setPath] = useState('');

  useEffect(() => {
    if (isOpen) {
      setPath('');
    }
  }, [isOpen]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (path.trim()) {
      onConfirm(path.trim());
      onClose();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      onClose();
    }
  };

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm animate-fade-in"
      onClick={onClose}
    >
      <div
        className="bg-surface border border-border rounded-lg shadow-2xl w-full max-w-md mx-4 p-6 relative"
        onClick={(e) => e.stopPropagation()}
        onKeyDown={handleKeyDown}
      >
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-textMain">创建新的任务</h2>
          <button
            onClick={onClose}
            className="text-textMuted hover:text-textMain transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Body */}
        <form onSubmit={handleSubmit}>
          <div className="mb-6">
            <label className="block text-sm font-medium text-textMain mb-2">
              Secret 路径
            </label>
            <input
              type="text"
              value={path}
              onChange={(e) => setPath(e.target.value)}
              placeholder="例如: /app/database/credential"
              className="w-full bg-surfaceHighlight border border-border rounded px-3 py-2 text-textMain placeholder:text-textMuted/50 focus:outline-none focus:border-accent focus:ring-1 focus:ring-accent transition-all"
              autoFocus
              autoComplete="off"
            />
            <p className="mt-2 text-xs text-textMuted leading-relaxed">
              请输入要监控的 Secret 路径
              <br/>
              <span className="opacity-70">按 <kbd className="font-mono bg-white/10 px-1 rounded text-textMain">Enter</kbd> 确认创建</span>
            </p>
          </div>

          {/* Footer */}
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 rounded-md border border-border text-sm text-textMain hover:bg-surfaceHighlight transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={!path.trim()}
              className="px-4 py-2 rounded-md bg-accent/10 border border-accent/20 text-accent text-sm font-medium hover:bg-accent/20 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              确认创建
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
