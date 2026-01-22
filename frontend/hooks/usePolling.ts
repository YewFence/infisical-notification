import { useEffect, useRef, useCallback } from 'react';

interface UsePollingOptions {
  /** 轮询间隔（毫秒） */
  interval: number;
  /** 是否启用轮询 */
  enabled?: boolean;
}

/**
 * 通用轮询 Hook
 * @param callback 轮询执行的回调函数
 * @param options 轮询配置选项
 */
export function usePolling(
  callback: () => Promise<void>,
  options: UsePollingOptions
) {
  const { interval, enabled = true } = options;

  const timerRef = useRef<number | null>(null);
  const isPollingRef = useRef(false);
  const isPausedRef = useRef(false);
  const callbackRef = useRef(callback);

  // 保持 callback 引用最新
  useEffect(() => {
    callbackRef.current = callback;
  }, [callback]);

  // 执行轮询
  const poll = useCallback(async () => {
    // 如果已经在轮询中或已暂停，跳过本次
    if (isPollingRef.current || isPausedRef.current) {
      return;
    }

    isPollingRef.current = true;
    try {
      await callbackRef.current();
    } catch (error) {
      // 轮询错误静默处理，避免打断用户体验
      console.warn('[Polling] 轮询请求失败:', error);
    } finally {
      isPollingRef.current = false;
    }
  }, []);

  // 暂停轮询（用户操作时调用）
  const pause = useCallback(() => {
    isPausedRef.current = true;
  }, []);

  // 恢复轮询
  const resume = useCallback(() => {
    isPausedRef.current = false;
  }, []);

  // 设置定时器
  useEffect(() => {
    if (!enabled || interval <= 0) {
      return;
    }

    // 启动轮询定时器
    timerRef.current = window.setInterval(poll, interval);

    // 清理函数：组件卸载时清除定时器
    return () => {
      if (timerRef.current !== null) {
        window.clearInterval(timerRef.current);
        timerRef.current = null;
      }
    };
  }, [enabled, interval, poll]);

  return { pause, resume };
}
