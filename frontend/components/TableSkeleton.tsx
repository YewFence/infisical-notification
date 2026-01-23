import React from 'react';

const SkeletonRow = () => (
  <div className="grid grid-cols-12 gap-4 px-4 py-3 border-b border-border items-center">
    {/* Name Column */}
    <div className="col-span-5 flex items-center gap-3">
      <div className="w-4 h-4 rounded bg-white/10 animate-pulse" />
      <div className="h-4 bg-white/10 rounded w-3/4 animate-pulse" />
    </div>

    {/* Status Column */}
    <div className="col-span-2 flex justify-center">
      <div className="h-6 w-16 bg-white/10 rounded-full animate-pulse" />
    </div>

    {/* Created Date Column */}
    <div className="col-span-2 flex justify-center">
      <div className="h-4 w-20 bg-white/10 rounded animate-pulse" />
    </div>

    {/* Completed Date Column */}
    <div className="col-span-2 flex justify-center">
      <div className="h-4 w-20 bg-white/10 rounded animate-pulse" />
    </div>

    {/* Actions Column */}
    <div className="col-span-1 flex justify-end">
      <div className="w-7 h-7 bg-white/10 rounded animate-pulse" />
    </div>
  </div>
);

export const TableSkeleton: React.FC<{ rows?: number }> = ({ rows = 5 }) => {
  return (
    <div className="w-full overflow-x-auto">
      <div className="min-w-[900px] border-t border-border">
        {/* Table Header */}
        <div className="grid grid-cols-12 gap-4 px-4 py-3 bg-surfaceHighlight/50 text-xs font-medium text-textMuted border-b border-border">
          <div className="col-span-5 pl-8">NAME</div>
          <div className="col-span-2 text-center">STATUS</div>
          <div className="col-span-2 text-center">CREATED</div>
          <div className="col-span-2 text-center">COMPLETED</div>
          <div className="col-span-1"></div>
        </div>

        {/* Skeleton Rows */}
        {Array.from({ length: rows }).map((_, index) => (
          <SkeletonRow key={index} />
        ))}
      </div>

      {/* Footer Skeleton */}
      <div className="px-4 py-3 flex items-center border-b border-border">
        <div className="w-4 h-4 bg-white/10 rounded animate-pulse mr-2" />
        <div className="h-4 w-16 bg-white/10 rounded animate-pulse" />
      </div>
    </div>
  );
};
