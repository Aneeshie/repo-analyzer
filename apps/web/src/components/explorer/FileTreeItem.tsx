"use client";

import { useMemo, useState } from "react";
import { cn } from "@/lib/utils";

export type FileTreeNodeType = "file" | "directory";

export interface FileTreeNode {
  name: string;
  path: string;
  type: FileTreeNodeType;
  children?: FileTreeNode[];
  size?: number;
  language?: string | null;
}

interface FileTreeItemProps {
  node: FileTreeNode;
  level?: number;
  selectedPath?: string | null;
  onFileSelect: (node: FileTreeNode) => void;
  defaultExpanded?: boolean;
}

// Native-looking SVG icons
const FolderIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
    <path d="M14.5 4H7.71l-1.5-1.5H1.5v11h13V4zM2.5 3.5h3.29l1.5 1.5H13.5v.5h-11v-2z" />
  </svg>
);

const FolderOpenIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
    <path d="M14.5 4H7.71l-1.5-1.5H1.5v11h13V4zM2.5 3.5h3.29l1.5 1.5H13.5v.5h-11v-2zM2.5 6.5h11v7h-11v-7z" />
  </svg>
);

const FileIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
    <path d="M13.85 4.44l-3.28-3.3-.35-.14H2.5v14h11V4.8l-.15-.36zM10 2.41l2.59 2.59H10V2.41zM12.5 14h-9V2h5.5v4h4v8z" />
  </svg>
);

const ChevronRightIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
    <path fillRule="evenodd" clipRule="evenodd" d="M10.072 8l-3.536 3.536.708.707L11.485 8 7.244 3.757l-.708.707L10.072 8z" />
  </svg>
);

function sortNodes(nodes: FileTreeNode[]) {
  return [...nodes].sort((a, b) => {
    if (a.type !== b.type) return a.type === "directory" ? -1 : 1;
    return a.name.localeCompare(b.name, undefined, { sensitivity: "base" });
  });
}

export function FileTreeItem({
  node,
  level = 0,
  selectedPath,
  onFileSelect,
  defaultExpanded = false,
}: FileTreeItemProps) {
  const isDirectory = node.type === "directory";
  const isSelected = node.type === "file" && selectedPath === node.path;
  const [expanded, setExpanded] = useState(defaultExpanded);
  const children = useMemo(
    () => sortNodes(node.children ?? []),
    [node.children],
  );

  const handleClick = () => {
    if (isDirectory) {
      setExpanded((current) => !current);
      return;
    }
    onFileSelect(node);
  };

  const Icon = isDirectory ? (expanded ? FolderOpenIcon : FolderIcon) : FileIcon;

  return (
    <div className="select-none">
      <button
        type="button"
        onClick={handleClick}
        title={node.path}
        className={cn(
          "group flex h-[24px] w-full items-center gap-1.5 pr-2 text-left text-[13px] leading-none transition-none",
          "text-zinc-400 hover:bg-white/5 hover:text-zinc-200",
          isSelected && "bg-white/10 text-white hover:bg-white/10",
        )}
        style={{ paddingLeft: `${level * 12 + 12}px` }}
      >
        <div className="flex h-4 w-4 shrink-0 items-center justify-center">
          {isDirectory ? (
            <ChevronRightIcon
              className={cn(
                "h-3.5 w-3.5 text-zinc-500 transition-transform duration-100 ease-out",
                expanded && "rotate-90",
              )}
            />
          ) : null}
        </div>
        <Icon
          className={cn(
            "h-4 w-4 shrink-0",
            isDirectory ? "text-zinc-300" : "text-zinc-400"
          )}
        />
        <span className="min-w-0 flex-1 truncate">{node.name}</span>
      </button>

      {isDirectory ? (
        <div
          className={cn(
            "grid",
            expanded
              ? "grid-rows-[1fr] opacity-100"
              : "grid-rows-[0fr] opacity-0 overflow-hidden",
          )}
        >
          <div className="overflow-hidden">
            {children.map((child) => (
              <FileTreeItem
                key={child.path}
                node={child}
                level={level + 1}
                selectedPath={selectedPath}
                onFileSelect={onFileSelect}
              />
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}
