"use client";

import { useEffect, useMemo, useState } from "react";
import { codeToHtml } from "shiki";
import { AlertTriangle, FileCode2, Loader2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import { Breadcrumbs } from "./Breadcrumbs";

interface CodeViewerProps {
  repoId: string;
  selectedPath?: string | null;
  className?: string;
}

interface FileContentResponse {
  path: string;
  name: string;
  content: string;
  language?: string | null;
  size?: number;
}

const LANGUAGE_ALIASES: Record<string, string> = {
  bash: "shellscript",
  cjs: "javascript",
  env: "dotenv",
  js: "javascript",
  jsx: "jsx",
  md: "markdown",
  mjs: "javascript",
  py: "python",
  rb: "ruby",
  rs: "rust",
  sh: "shellscript",
  ts: "typescript",
  tsx: "tsx",
  yml: "yaml",
};

function getExtension(path: string) {
  const filename = path.split("/").pop() ?? path;
  const parts = filename.toLowerCase().split(".");
  return parts.length > 1 ? (parts.pop() ?? "") : "";
}

function normalizeLanguage(language: string | null | undefined, path: string) {
  const input = (language || getExtension(path) || "text").toLowerCase();
  return LANGUAGE_ALIASES[input] ?? input;
}

function formatBytes(size?: number) {
  if (typeof size !== "number") return null;
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / (1024 * 1024)).toFixed(1)} MB`;
}

function CodeSkeleton() {
  return (
    <div className="space-y-3 p-5">
      {Array.from({ length: 18 }).map((_, index) => (
        <div key={index} className="flex items-center gap-4">
          <Skeleton className="h-4 w-8 bg-zinc-700" />
          <Skeleton
            className={cn(
              "h-4 bg-zinc-700",
              index % 5 === 0
                ? "w-3/4"
                : index % 5 === 1
                  ? "w-1/2"
                  : index % 5 === 2
                    ? "w-5/6"
                    : "w-2/3",
            )}
          />
        </div>
      ))}
    </div>
  );
}

export function CodeViewer({
  repoId,
  selectedPath,
  className,
}: CodeViewerProps) {
  const [file, setFile] = useState<FileContentResponse | null>(null);
  const [highlightedHtml, setHighlightedHtml] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [highlighting, setHighlighting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const displayLanguage = useMemo(
    () =>
      file
        ? normalizeLanguage(file.language, file.path)
        : selectedPath
          ? normalizeLanguage(null, selectedPath)
          : null,
    [file, selectedPath],
  );

  useEffect(() => {
    const filePath = selectedPath;
    if (typeof filePath !== "string" || filePath.length === 0) return;

    const controller = new AbortController();

    async function fetchFile(path: string) {
      try {
        await Promise.resolve();
        setLoading(true);
        setHighlighting(false);
        setError(null);
        setFile(null);
        setHighlightedHtml("");

        const response = await fetch(
          `/api/v1/repos/${repoId}/file?path=${encodeURIComponent(path)}`,
          {
            signal: controller.signal,
          },
        );

        if (!response.ok) {
          let message = "File too large or unavailable";
          try {
            const body = await response.json();
            message = body?.message || body?.error || message;
          } catch {
            // Keep the generic message when the backend does not return JSON.
          }
          throw new Error(message);
        }

        const data = (await response.json()) as FileContentResponse;
        if (!controller.signal.aborted) setFile(data);
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") return;
        setError(
          err instanceof Error ? err.message : "File too large or unavailable",
        );
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }

    fetchFile(filePath);

    return () => controller.abort();
  }, [repoId, selectedPath]);

  useEffect(() => {
    const currentFile = file;
    if (!currentFile) return;

    let cancelled = false;

    async function highlight(fileToHighlight: FileContentResponse) {
      try {
        await Promise.resolve();
        setHighlighting(true);
        const lang = normalizeLanguage(
          fileToHighlight.language,
          fileToHighlight.path,
        );
        let html: string;

        try {
          html = await codeToHtml(fileToHighlight.content, {
            lang,
            theme: "dark-plus",
          });
        } catch {
          html = await codeToHtml(fileToHighlight.content, {
            lang: "text",
            theme: "dark-plus",
          });
        }

        if (!cancelled) setHighlightedHtml(html);
      } catch (err) {
        if (!cancelled)
          setError(
            err instanceof Error ? err.message : "Failed to highlight file",
          );
      } finally {
        if (!cancelled) setHighlighting(false);
      }
    }

    highlight(currentFile);

    return () => {
      cancelled = true;
    };
  }, [file]);

  return (
    <section
      className={cn(
        "flex h-full min-h-0 flex-col bg-transparent text-zinc-200",
        className,
      )}
    >
      <div className="flex h-14 shrink-0 items-center justify-between gap-3 border-b border-white/5 bg-zinc-950/30 px-5 backdrop-blur-md">
        <Breadcrumbs path={selectedPath} className="flex-1" />
        {selectedPath && file ? (
          <div className="flex shrink-0 items-center gap-2">
            {displayLanguage ? (
              <Badge
                variant="outline"
                className="border-zinc-700 bg-zinc-900/60 text-[10px] uppercase text-zinc-400"
              >
                {displayLanguage}
              </Badge>
            ) : null}
            {formatBytes(file.size) ? (
              <Badge
                variant="outline"
                className="border-zinc-700 bg-zinc-900/60 text-[10px] text-zinc-400"
              >
                {formatBytes(file.size)}
              </Badge>
            ) : null}
          </div>
        ) : null}
      </div>

      <div className="min-h-0 flex-1 overflow-hidden">
        {!selectedPath ? (
          <div className="flex h-full flex-col items-center justify-center p-8 text-center text-zinc-400">
            <div className="relative mb-6 flex h-20 w-20 items-center justify-center rounded-2xl border border-white/10 bg-gradient-to-b from-white/5 to-white/0 shadow-2xl backdrop-blur-xl">
              <div className="absolute inset-0 rounded-2xl bg-blue-500/10 blur-xl" />
              <FileCode2 className="relative z-10 h-10 w-10 text-blue-400 drop-shadow-[0_0_15px_rgba(96,165,250,0.5)]" />
            </div>
            <h3 className="text-lg font-medium text-zinc-200 tracking-tight">
              No file selected
            </h3>
            <p className="mt-2 max-w-sm text-sm leading-relaxed text-zinc-500">
              Select a file from the explorer to view its source code with
              syntax highlighting and line numbers.
            </p>
          </div>
        ) : null}

        {selectedPath && loading ? <CodeSkeleton /> : null}

        {selectedPath && !loading && error ? (
          <div className="flex h-full flex-col items-center justify-center p-8 text-center">
            <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl border border-red-500/30 bg-red-500/10">
              <AlertTriangle className="h-8 w-8 text-red-300" />
            </div>
            <h3 className="text-base font-semibold text-red-200">
              File too large or unavailable
            </h3>
            <p className="mt-2 max-w-md text-sm leading-relaxed text-zinc-400">
              {error}
            </p>
          </div>
        ) : null}

        {selectedPath && !loading && !error && file ? (
          <div className="h-full overflow-auto">
            {highlighting && !highlightedHtml ? (
              <div className="flex h-full items-center justify-center gap-2 text-sm text-zinc-400">
                <Loader2 className="h-4 w-4 animate-spin" />
                Highlighting source…
              </div>
            ) : (
              <div
                className="explorer-shiki h-full min-w-max text-[13px] leading-6"
                dangerouslySetInnerHTML={{ __html: highlightedHtml }}
              />
            )}
          </div>
        ) : null}
      </div>

      <style jsx global>{`
        .explorer-shiki pre {
          margin: 0;
          min-height: 100%;
          background: transparent !important;
          padding: 20px 0;
        }

        .explorer-shiki code {
          display: block;
          counter-reset: line;
          font-family:
            var(--font-geist-mono), ui-monospace, SFMono-Regular, Menlo, Monaco,
            Consolas, "Liberation Mono", "Courier New", monospace;
        }

        .explorer-shiki .line {
          display: block;
          min-height: 1.5rem;
          padding-right: 24px;
          white-space: pre;
        }

        .explorer-shiki .line::before {
          counter-increment: line;
          content: counter(line);
          display: inline-block;
          width: 56px;
          margin-right: 16px;
          padding-right: 16px;
          text-align: right;
          color: #6e7681;
          user-select: none;
        }

        .explorer-shiki .line:hover {
          background: rgba(255, 255, 255, 0.035);
        }
      `}</style>
    </section>
  );
}
