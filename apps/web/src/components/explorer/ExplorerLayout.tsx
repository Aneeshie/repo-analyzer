"use client";

import { useState } from "react";
import { Group, Panel, Separator } from "react-resizable-panels";
import {
  FolderTree,
  GripVertical,
  PanelLeftClose,
  PanelLeftOpen,
  X,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { CodeViewer } from "./CodeViewer";
import { FileTree } from "./FileTree";
import type { FileTreeNode } from "./FileTreeItem";

interface ExplorerLayoutProps {
  repoId: string;
  className?: string;
}

export function ExplorerLayout({ repoId, className }: ExplorerLayoutProps) {
  const [selectedFile, setSelectedFile] = useState<FileTreeNode | null>(null);
  const [mobileTreeOpen, setMobileTreeOpen] = useState(false);

  const handleFileSelect = (node: FileTreeNode) => {
    setSelectedFile(node);
    setMobileTreeOpen(false);
  };

  return (
    <div className="relative group">
      {/* Subtle background glow effect */}
      <div className="absolute -inset-1 -z-10 rounded-xl bg-gradient-to-r from-blue-500/20 via-indigo-500/20 to-purple-500/20 opacity-50 blur-xl transition-all duration-500 group-hover:opacity-70" />
      
      <Card
        className={cn(
          "overflow-hidden border border-white/10 bg-zinc-950/80 py-0 shadow-2xl backdrop-blur-xl",
          className,
        )}
      >
        <CardHeader className="h-14 border-b border-white/5 bg-transparent px-4">
          <div className="flex items-center justify-between gap-3">
            <CardTitle className="flex items-center gap-2 text-sm font-semibold text-zinc-100 tracking-wide">
              <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 shadow-lg shadow-blue-500/20">
                <FolderTree className="h-4 w-4 text-white" />
              </div>
              <span className="bg-gradient-to-r from-white to-zinc-400 bg-clip-text text-transparent">Repository Explorer</span>
            </CardTitle>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="border-white/10 bg-white/5 text-zinc-200 hover:bg-white/10 md:hidden"
              onClick={() => setMobileTreeOpen((open) => !open)}
            >
              {mobileTreeOpen ? (
                <PanelLeftClose className="h-4 w-4" />
              ) : (
                <PanelLeftOpen className="h-4 w-4" />
              )}
              Files
            </Button>
          </div>
        </CardHeader>

      <CardContent className="relative h-[70vh] min-h-[520px] p-0">
        <div className="hidden h-full md:block">
          <Group orientation="horizontal" className="h-full">
            <Panel
              defaultSize="300px"
              minSize="220px"
              maxSize="45%"
              className="min-w-[220px]"
            >
              <FileTree
                repoId={repoId}
                selectedPath={selectedFile?.path}
                onFileSelect={handleFileSelect}
              />
            </Panel>

            <Separator className="group relative flex w-[3px] items-center justify-center bg-white/5 transition-all hover:w-[4px] hover:bg-blue-500/70 data-[resize-handle-active]:w-[4px] data-[resize-handle-active]:bg-blue-500">
              <div className="absolute left-1/2 top-1/2 hidden -translate-x-1/2 -translate-y-1/2 rounded-full border border-white/10 bg-zinc-900/80 p-1 text-zinc-400 shadow-xl backdrop-blur-md transition-all group-hover:flex">
                <GripVertical className="h-4 w-4 opacity-70" />
              </div>
            </Separator>

            <Panel minSize={45}>
              <CodeViewer repoId={repoId} selectedPath={selectedFile?.path} />
            </Panel>
          </Group>
        </div>

        <div className="h-full md:hidden">
          <CodeViewer repoId={repoId} selectedPath={selectedFile?.path} />

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
      </CardContent>
    </Card>
    </div>
  );
}
