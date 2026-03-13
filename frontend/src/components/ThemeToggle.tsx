"use client";

import { Moon, Sun } from "lucide-react";
import { useTheme } from "@/lib/theme-context";

interface ThemeToggleProps {
  className?: string;
  variant?: "default" | "compact" | "sidebar";
}

export function ThemeToggle({ className = "", variant = "default" }: ThemeToggleProps) {
  const { theme, toggleTheme } = useTheme();

  if (variant === "compact") {
    return (
      <button
        onClick={toggleTheme}
        className={`inline-flex items-center justify-center rounded-lg p-2 transition-colors hover:bg-accent text-muted-foreground hover:text-foreground ${className}`}
        title={theme === "dark" ? "Light mode" : "Dark mode"}
      >
        {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
      </button>
    );
  }

  if (variant === "sidebar") {
    return (
      <button
        onClick={toggleTheme}
        className={`flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors text-muted-foreground hover:bg-accent hover:text-foreground w-full ${className}`}
      >
        {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        {theme === "dark" ? "Yorug' rejim" : "Qorong'u rejim"}
      </button>
    );
  }

  return (
    <button
      onClick={toggleTheme}
      className={`inline-flex items-center justify-center rounded-lg border border-border bg-background p-2 transition-colors hover:bg-accent text-muted-foreground hover:text-foreground ${className}`}
      title={theme === "dark" ? "Light mode" : "Dark mode"}
    >
      {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
    </button>
  );
}
