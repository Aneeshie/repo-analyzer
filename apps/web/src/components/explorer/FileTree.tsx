"use client";

import { useEffect, useMemo, useState } from "react";
import { AlertCircle, FolderTree, RefreshCw } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import { FileTreeItem, type FileTreeNode } from "./FileTreeItem";

interface FileTreeProps {
  repoId: string;
  selectedPath?: string | null;
  onFileSelect: (node: FileTreeNode) => void;
  className?: string;
}

function sortNodes(nodes: FileTreeNode[]) {
  return [...nodes].sort((a, b) => {
    if (a.type !== b.type) return a.type === "directory" ? -1 : 1;
    return a.name.localeCompare(b.name, undefined, { sensitivity: "base" });
  });
}

function countFiles(nodes: FileTreeNode[]): number {
  return nodes.reduce((total, node) => {
    if (node.type === "file") return total + 1;
    return total + countFiles(node.children ?? []);
  }, 0);
}

function FileTreeSkeleton() {
  return (
    <div className="space-y-2 p-3">
      {Array.from({ length: 12 }).map((_, index) => (
        <div
          key={index}
          className="flex items-center gap-2"
          style={{ paddingLeft: `${(index % 4) * 14}px` }}
        >
          <Skeleton className="h-4 w-4 rounded-sm bg-zinc-700" />
          <Skeleton
            className={cn(
              "h-4 bg-zinc-700",
              index % 3 === 0 ? "w-24" : index % 3 === 1 ? "w-36" : "w-44",
            )}
          />
        </div>
      ))}
    </div>
  );
}

export function FileTree({
  repoId,
  selectedPath,
  onFileSelect,
  className,
}: FileTreeProps) {
  const [tree, setTree] = useState<FileTreeNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const sortedTree = useMemo(() => sortNodes(tree), [tree]);
  const fileCount = useMemo(() => countFiles(tree), [tree]);

  useEffect(() => {
    const controller = new AbortController();

    async function fetchTree() {
      try {
        setLoading(true);
        setError(null);

        const response = await fetch(`/api/v1/repos/${repoId}/tree`, {
          signal: controller.signal,
        });

        if (!response.ok) {
          throw new Error("Failed to load repository tree");
        }

        const data = (await response.json()) as FileTreeNode[];
        setTree(Array.isArray(data) ? data : []);
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return;
        setError(
          err instanceof Error ? err.message : "Failed to load repository tree",
        );
        setTree([]);
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }

    fetchTree();

    return () => controller.abort();
  }, [repoId]);

  const retry = () => {
    setTree([]);
    setLoading(true);
    setError(null);
    fetch(`/api/v1/repos/${repoId}/tree`)
      .then((response) => {
        if (!response.ok) throw new Error("Failed to load repository tree");
        return response.json();
      })
      .then((data: FileTreeNode[]) => setTree(Array.isArray(data) ? data : []))
      .catch((err) =>
        setError(
          err instanceof Error ? err.message : "Failed to load repository tree",
        ),
      )
      .finally(() => setLoading(false));
  };

  return (
    <aside
      className={cn(
        "flex h-full min-h-0 flex-col bg-transparent text-zinc-200",
        className,
      )}
    >
      <div className="flex h-12 shrink-0 items-center justify-between border-b border-white/5 px-4 backdrop-blur-sm">
        <div className="flex items-center gap-2 text-xs font-bold uppercase tracking-widest text-zinc-400">
          <FolderTree className="h-4 w-4 text-zinc-500" />
          Explorer
        </div>
        {!loading && !error ? (
          <Badge
            variant="outline"
            className="h-5 border-blue-500/30 bg-blue-500/10 px-2 text-[10px] font-semibold text-blue-300 shadow-[0_0_10px_rgba(59,130,246,0.1)]"
          >
            {fileCount} files
          </Badge>
        ) : null}
      </div>

      <ScrollArea className="min-h-0 flex-1 py-2">
        {loading ? <FileTreeSkeleton /> : null}

        {!loading && error ? (
          <div className="flex h-full flex-col items-center justify-center p-6 text-center">
            <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-red-500/10">
              <AlertCircle className="h-6 w-6 text-red-400" />
            </div>
            <p className="text-sm font-medium text-zinc-200">
              Unable to load file tree
            </p>
            <p className="mt-1 text-xs text-zinc-500">{error}</p>
            <Button
              onClick={retry}
              size="sm"
              variant="outline"
              className="mt-4 border-white/10 bg-white/5 text-zinc-200 hover:bg-white/10"
            >
              <RefreshCw className="mr-1.5 h-3.5 w-3.5" />
              Try again
            </Button>
          </div>
        ) : null}

        {!loading && !error && sortedTree.length === 0 ? (
          <div className="flex h-full flex-col items-center justify-center p-6 text-center text-zinc-500">
            <FolderTree className="mb-3 h-8 w-8" />
            <p className="text-sm">No files found</p>
          </div>
        ) : null}

        {!loading && !error && sortedTree.length > 0 ? (
          <div className="px-1">
            {sortedTree.map((node) => (
              <FileTreeItem
                key={node.path}
                node={node}
                selectedPath={selectedPath}
                onFileSelect={onFileSelect}
                defaultExpanded={node.type === "directory"}
              />
            ))}
          </div>
        ) : null}
      </ScrollArea>
    </aside>
  );
}
