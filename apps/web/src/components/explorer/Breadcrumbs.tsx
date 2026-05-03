"use client";

import { ChevronRight, FileCode2, Home } from "lucide-react";
import { cn } from "@/lib/utils";

interface BreadcrumbsProps {
  path?: string | null;
  className?: string;
}

export function Breadcrumbs({ path, className }: BreadcrumbsProps) {
  const segments = path?.split("/").filter(Boolean) ?? [];

  return (
    <nav
      aria-label="File path"
      className={cn("flex min-w-0 items-center gap-1 overflow-hidden text-xs text-zinc-400", className)}
    >
      <Home className="h-3.5 w-3.5 shrink-0 text-zinc-500" />
      {segments.length === 0 ? (
        <span className="truncate text-zinc-500">Select a file</span>
      ) : (
        segments.map((segment, index) => {
          const isLast = index === segments.length - 1;

          return (
            <span key={`${segment}-${index}`} className="flex min-w-0 items-center gap-1">
              <ChevronRight className="h-3.5 w-3.5 shrink-0 text-blue-500/50" />
              <span
                className={cn(
                  "truncate",
                  isLast ? "font-medium text-blue-100 drop-shadow-[0_0_10px_rgba(59,130,246,0.3)]" : "text-zinc-400 transition-colors hover:text-zinc-300"
                )}
                title={segments.slice(0, index + 1).join("/")}
              >
                {isLast ? <FileCode2 className="mr-1 inline h-3.5 w-3.5 text-blue-400 drop-shadow-[0_0_8px_rgba(96,165,250,0.5)]" /> : null}
                {segment}
              </span>
            </span>
          );
        })
      )}
    </nav>
  );
}
