"use client";

import { useMemo, useState } from "react";
import {
  Braces,
  ChevronRight,
  Code2,
  Database,
  File,
  FileCode2,
  FileCog,
  FileImage,
  FileJson,
  FileText,
  Folder,
  FolderOpen,
  Hash,
  TerminalSquare,
  type LucideIcon,
} from "lucide-react";
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

function getExtension(name: string) {
  const parts = name.toLowerCase().split(".");
  return parts.length > 1 ? (parts.pop() ?? "") : "";
}

function normalizeLanguage(node: FileTreeNode) {
  const language = node.language?.toLowerCase();
  if (language) return language;
  return getExtension(node.name);
}

function getFileIcon(node: FileTreeNode): {
  icon: LucideIcon;
  className: string;
} {
  const language = normalizeLanguage(node);
  const name = node.name.toLowerCase();

  if (["package.json", "tsconfig.json", "jsconfig.json"].includes(name)) {
    return { icon: FileJson, className: "text-yellow-400" };
  }

  if (
    ["dockerfile"].includes(name) ||
    ["sh", "bash", "zsh", "fish", "ps1"].includes(language)
  ) {
    return { icon: TerminalSquare, className: "text-emerald-400" };
  }

  if (["json", "jsonc", "json5"].includes(language))
    return { icon: FileJson, className: "text-yellow-400" };
  if (["md", "mdx", "markdown", "txt", "rst"].includes(language))
    return { icon: FileText, className: "text-sky-300" };
  if (
    ["png", "jpg", "jpeg", "gif", "svg", "webp", "ico", "bmp"].includes(
      language,
    )
  )
    return { icon: FileImage, className: "text-fuchsia-300" };
  if (["sql", "prisma", "graphql", "gql"].includes(language))
    return { icon: Database, className: "text-cyan-300" };
  if (["yml", "yaml", "toml", "ini", "env", "config"].includes(language))
    return { icon: FileCog, className: "text-orange-300" };
  if (["css", "scss", "sass", "less"].includes(language))
    return { icon: Hash, className: "text-blue-300" };
  if (["html", "xml", "vue", "svelte"].includes(language))
    return { icon: Braces, className: "text-orange-400" };
  if (["ts", "tsx", "typescript"].includes(language))
    return { icon: FileCode2, className: "text-blue-400" };
  if (["js", "jsx", "javascript", "mjs", "cjs"].includes(language))
    return { icon: FileCode2, className: "text-yellow-300" };
  if (
    [
      "py",
      "python",
      "go",
      "rs",
      "rust",
      "java",
      "kt",
      "kotlin",
      "swift",
      "php",
      "rb",
      "ruby",
      "c",
      "cpp",
      "cs",
    ].includes(language)
  ) {
    return { icon: Code2, className: "text-violet-300" };
  }

  return { icon: File, className: "text-zinc-400" };
}

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
  const fileIcon = getFileIcon(node);

  const handleClick = () => {
    if (isDirectory) {
      setExpanded((current) => !current);
      return;
    }

    onFileSelect(node);
  };

  const Icon = isDirectory ? (expanded ? FolderOpen : Folder) : fileIcon.icon;

  return (
    <div className="select-none">
      <button
        type="button"
        onClick={handleClick}
        title={node.path}
        className={cn(
          "group flex h-7 w-full items-center gap-2 rounded-md pr-2 text-left text-[13px] leading-none transition-all duration-200 ease-out",
          "text-zinc-400 hover:bg-white/5 hover:text-zinc-100",
          isSelected && "bg-gradient-to-r from-blue-500/20 to-transparent border-l-2 border-blue-500 text-blue-100 hover:bg-blue-500/20 shadow-[inset_1px_0_0_0_rgba(59,130,246,1)]",
        )}
        style={{ paddingLeft: `${level * 14 + (isSelected ? 6 : 8)}px` }}
      >
        <span className="flex h-4 w-4 shrink-0 items-center justify-center">
          {isDirectory ? (
            <ChevronRight
              className={cn(
                "h-3.5 w-3.5 text-zinc-400 transition-transform duration-200 ease-out",
                expanded && "rotate-90",
              )}
            />
          ) : null}
        </span>
        <Icon
          className={cn(
            "h-4 w-4 shrink-0 transition-colors",
            isDirectory ? "text-[#dcb67a]" : fileIcon.className,
          )}
        />
        <span className="min-w-0 flex-1 truncate">{node.name}</span>
      </button>

      {isDirectory ? (
        <div
          className={cn(
            "grid transition-all duration-200 ease-out",
            expanded
              ? "grid-rows-[1fr] opacity-100"
              : "grid-rows-[0fr] opacity-0",
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
