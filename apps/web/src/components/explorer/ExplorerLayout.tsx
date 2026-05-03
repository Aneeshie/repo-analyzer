"use client";

import { useState } from "react";
import { Group, Panel, Separator } from "react-resizable-panels";
import {
  FolderTree,
  PanelLeftClose,
  PanelLeftOpen,
  X,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { CodeViewer } from "./CodeViewer";
import { FileTree } from "./FileTree";
import type { FileTreeNode } from "./FileTreeItem";

interface Repo {
  id: string;
  url: string;
  status: string;
  created_at: string;
  entry_points?: string[];
}

interface ExplorerLayoutProps {
  repo: Repo;
  className?: string;
}

export function ExplorerLayout({ repo, className }: ExplorerLayoutProps) {
  const repoId = repo.id;
  const [selectedFile, setSelectedFile] = useState<FileTreeNode | null>(null);
  const [mobileTreeOpen, setMobileTreeOpen] = useState(false);

  const handleFileSelect = (node: FileTreeNode) => {
    setSelectedFile(node);
    setMobileTreeOpen(false);
  };

  return (
    <div className={cn("flex flex-col h-full w-full bg-[#1e1e1e] overflow-hidden", className)}>
      <div className="flex h-10 shrink-0 items-center justify-between border-b border-white/5 px-4 bg-[#1e1e1e]">
        <div className="flex items-center gap-2 text-xs font-medium text-zinc-400 uppercase tracking-widest">
          Explorer
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="border-white/10 bg-white/5 text-zinc-200 hover:bg-white/10 md:hidden h-7 text-xs px-2"
          onClick={() => setMobileTreeOpen((open) => !open)}
        >
          {mobileTreeOpen ? (
            <PanelLeftClose className="h-3 w-3 mr-1" />
          ) : (
            <PanelLeftOpen className="h-3 w-3 mr-1" />
          )}
          Files
        </Button>
      </div>

      <div className="relative flex-1 min-h-0">
        <div className="hidden h-full md:block">
          <Group orientation="horizontal" className="h-full">
            <Panel
              defaultSize="20"
              minSize="15"
              maxSize="40"
              className="min-w-[200px]"
            >
              <FileTree
                repoId={repoId}
                selectedPath={selectedFile?.path}
                onFileSelect={handleFileSelect}
              />
            </Panel>

            <Separator className="flex w-[1px] bg-white/10 hover:w-[2px] hover:bg-blue-500 transition-all cursor-col-resize" />

            <Panel minSize={45}>
              <CodeViewer repo={repo} selectedPath={selectedFile?.path} onFileSelect={handleFileSelect} />
            </Panel>
          </Group>
        </div>

        <div className="h-full md:hidden">
          <CodeViewer repo={repo} selectedPath={selectedFile?.path} onFileSelect={handleFileSelect} />

          <div
            className={cn(
              "absolute inset-0 z-10 bg-black/60 transition-opacity duration-200",
              mobileTreeOpen ? "opacity-100" : "pointer-events-none opacity-0",
            )}
            onClick={() => setMobileTreeOpen(false)}
          />

          <div
            className={cn(
              "absolute inset-y-0 left-0 z-20 w-[86%] max-w-sm border-r border-white/10 bg-zinc-950/95 shadow-2xl backdrop-blur-2xl transition-transform duration-300 ease-out",
              mobileTreeOpen ? "translate-x-0" : "-translate-x-full",
            )}
          >
            <div className="flex h-12 items-center justify-between border-b border-white/5 px-4 text-zinc-200">
              <span className="flex items-center gap-2 text-xs font-bold uppercase tracking-widest text-zinc-400">
                <FolderTree className="h-4 w-4 text-indigo-400" />
                Files
              </span>
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                className="text-zinc-400 hover:bg-white/10 hover:text-white"
                onClick={() => setMobileTreeOpen(false)}
                aria-label="Close file explorer"
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
            <div className="h-[calc(100%-3rem)]">
              <FileTree
                repoId={repoId}
                selectedPath={selectedFile?.path}
                onFileSelect={handleFileSelect}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
