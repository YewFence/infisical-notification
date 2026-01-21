import React from 'react';

interface LayoutProps {
  children: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen bg-background text-textMain font-sans selection:bg-accent selection:text-black">
      {/* Main Content */}
      <main className="p-6 md:p-8 max-w-[1600px] mx-auto animate-fade-in">
        {children}
      </main>
    </div>
  );
};