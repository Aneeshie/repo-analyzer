"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Badge } from "@/components/ui/badge";
import {
  Loader2,
  ArrowLeft,
  CheckCircle2,
  XCircle,
  Sparkles,
  ChevronDown,
  ChevronUp,
  GitBranch,
  Code2
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { ExplorerLayout } from "@/components/explorer/ExplorerLayout";

type RepoStatus = "pending" | "processing" | "completed" | "failed" | "analyzing" | "cloning" | "parsing" | string;

const normalizeStatus = (status: RepoStatus): "pending" | "processing" | "completed" | "failed" => {
  if (["analyzing", "cloning", "parsing", "processing"].includes(status as string)) return "processing";
  if (status === "completed" || status === "failed") return status as "completed" | "failed";
  return "pending";
};

interface Repo {
  id: string;
  url: string;
  status: RepoStatus;
  created_at: string;
}

export default function AnalysisPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [repo, setRepo] = useState<Repo | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [showSummary, setShowSummary] = useState(true);

  const pollIntervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const pollStatus = async () => {
      try {
        const res = await fetch(`/api/v1/repos/${id}`);
        if (!res.ok) throw new Error("Failed to fetch repository status");

        const data: Repo = await res.json();
        setRepo(data);

        if (data.status === "completed" || data.status === "failed") {
          if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Error fetching status");
        if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
      }
    };

    pollStatus();
    pollIntervalRef.current = setInterval(pollStatus, 2000);

    return () => {
      if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
    };
  }, [id]);

  if (error) {
    return (
      <div className="flex h-screen flex-col items-center justify-center bg-zinc-950 p-8">
        <XCircle className="mb-4 h-16 w-16 text-red-500" />
        <h1 className="mb-2 text-2xl font-bold text-white">Analysis Failed</h1>
        <p className="mb-6 text-zinc-400">{error}</p>
        <Button onClick={() => router.push("/")} variant="outline" className="h-12 border-white/10 bg-white/5 px-6 hover:bg-white/10">
          <ArrowLeft className="mr-2 h-4 w-4" /> Go Back
        </Button>
      </div>
    );
  }

  const displayStatus = normalizeStatus(repo?.status || "pending");

  return (
    <div className="relative flex h-screen flex-col bg-zinc-950 text-zinc-200 selection:bg-blue-500/30">
      {/* Premium ambient background */}
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(120,119,198,0.15),rgba(255,255,255,0))]" />
      
      {/* Top Navigation Bar */}
      <header className="relative z-10 flex h-14 shrink-0 items-center justify-between border-b border-white/[0.08] bg-zinc-950/40 px-4 backdrop-blur-2xl">
        <div className="flex items-center gap-4">
          <Button onClick={() => router.push("/")} variant="ghost" size="icon" className="h-8 w-8 text-zinc-400 hover:bg-white/10 hover:text-white transition-colors duration-300">
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div className="flex items-center gap-3">
            <div className="relative flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-indigo-600 shadow-lg shadow-blue-500/20 ring-1 ring-white/10">
              <div className="absolute inset-0 rounded-lg bg-white/20 mix-blend-overlay" />
              <Code2 className="h-4 w-4 text-white drop-shadow-md" />
            </div>
            <div className="flex flex-col">
              <span className="text-[13px] font-semibold text-zinc-100 tracking-wide">
                {repo ? repo.url.replace("https://github.com/", "") : "Loading workspace..."}
              </span>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-3">
          {displayStatus === "completed" && (
            <div className="relative flex items-center gap-2 rounded-full border border-emerald-500/30 bg-emerald-500/10 px-3 py-1 text-[11px] font-semibold tracking-wide text-emerald-400 shadow-[0_0_15px_rgba(16,185,129,0.15)] backdrop-blur-sm">
              <span className="relative flex h-2 w-2">
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75"></span>
                <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500"></span>
              </span>
              READY
            </div>
          )}
          {displayStatus === "failed" && (
            <Badge variant="destructive" className="px-3 py-1 text-xs font-medium border border-red-500/20 shadow-lg">
              <XCircle className="mr-1.5 h-3.5 w-3.5" /> Failed
            </Badge>
          )}
          {(displayStatus === "processing" || displayStatus === "pending") && (
            <div className="flex items-center gap-2 rounded-full border border-blue-500/30 bg-blue-500/10 px-3 py-1 text-[11px] font-semibold tracking-wide text-blue-400 backdrop-blur-sm">
              <Loader2 className="h-3 w-3 animate-spin" />
              {displayStatus === "processing" ? "ANALYZING..." : "PENDING..."}
            </div>
          )}
        </div>
      </header>

      {/* Main Content Area */}
      {displayStatus === "completed" ? (
        <main className="relative z-10 flex min-h-0 flex-1 flex-col overflow-hidden p-4 md:p-6 gap-6 animate-in fade-in zoom-in-[0.98] duration-700 ease-out">
          
          {/* AI Summary Section (Experimental) */}
          <section className="group relative shrink-0 rounded-2xl border border-white/5 bg-zinc-950/40 shadow-2xl backdrop-blur-3xl transition-all duration-500 hover:border-white/10">
            {/* Animated glowing border effect */}
            <div className="pointer-events-none absolute -inset-px -z-10 rounded-2xl bg-gradient-to-r from-blue-500/20 via-indigo-500/20 to-purple-500/20 opacity-0 blur-sm transition-opacity duration-500 group-hover:opacity-100" />
            
            <button 
              onClick={() => setShowSummary(!showSummary)}
              className="flex w-full items-center justify-between p-4 px-6 outline-none transition-colors"
            >
              <div className="flex items-center gap-3">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 shadow-lg shadow-purple-500/20 ring-1 ring-white/10">
                  <Sparkles className="h-4 w-4 text-white" />
                </div>
                <h2 className="text-[15px] font-semibold text-zinc-100 tracking-wide">Repository Architecture Summary</h2>
              </div>
              <div className="flex items-center gap-3">
                <Badge variant="outline" className="border-indigo-500/30 bg-indigo-500/10 text-indigo-300 text-[10px] uppercase tracking-wider shadow-[0_0_10px_rgba(99,102,241,0.1)]">AI Powered</Badge>
                {showSummary ? <ChevronUp className="h-4 w-4 text-zinc-500" /> : <ChevronDown className="h-4 w-4 text-zinc-500" />}
              </div>
            </button>
            
            {showSummary && (
              <div className="border-t border-white/[0.04] p-6 text-[13px] leading-relaxed text-zinc-400 animate-in slide-in-from-top-4 fade-in duration-500 ease-out">
                <p className="mb-6 max-w-4xl text-zinc-300 text-[14px]">
                  <strong className="text-zinc-100 font-semibold mr-2">Overview:</strong> 
                  This is a highly scalable full-stack web application built using <span className="text-blue-400 font-medium bg-blue-500/10 px-1 py-0.5 rounded">Next.js 15</span>, <span className="text-emerald-400 font-medium bg-emerald-500/10 px-1 py-0.5 rounded">Go 1.22</span>, and <span className="text-indigo-400 font-medium bg-indigo-500/10 px-1 py-0.5 rounded">PostgreSQL</span>. It employs a clean architecture pattern with distinct handler, service, and repository layers on the backend.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                  <div className="rounded-xl border border-white/5 bg-black/20 p-4">
                    <h3 className="mb-3 text-[11px] font-bold uppercase tracking-widest text-zinc-500 flex items-center gap-2">
                      <div className="h-1.5 w-1.5 rounded-full bg-purple-500" />
                      Architecture Highlights
                    </h3>
                    <ul className="space-y-2.5 text-zinc-300">
                      <li className="flex items-start gap-2">
                        <CheckCircle2 className="h-4 w-4 text-zinc-600 shrink-0 mt-0.5" />
                        <span>Client-side React with App Router for fast UI rendering</span>
                      </li>
                      <li className="flex items-start gap-2">
                        <CheckCircle2 className="h-4 w-4 text-zinc-600 shrink-0 mt-0.5" />
                        <span>Go microservices for high-performance processing</span>
                      </li>
                      <li className="flex items-start gap-2">
                        <CheckCircle2 className="h-4 w-4 text-zinc-600 shrink-0 mt-0.5" />
                        <span>Async worker pool implementation for heavy IO tasks</span>
                      </li>
                    </ul>
                  </div>
                  <div className="rounded-xl border border-white/5 bg-black/20 p-4">
                    <h3 className="mb-3 text-[11px] font-bold uppercase tracking-widest text-zinc-500 flex items-center gap-2">
                      <div className="h-1.5 w-1.5 rounded-full bg-blue-500" />
                      Key Entry Points
                    </h3>
                    <div className="flex flex-wrap gap-2.5">
                      <Badge variant="outline" className="border-white/10 bg-zinc-900 text-xs text-zinc-300 hover:bg-white/5 hover:text-white transition-colors cursor-pointer py-1">cmd/api/main.go</Badge>
                      <Badge variant="outline" className="border-white/10 bg-zinc-900 text-xs text-zinc-300 hover:bg-white/5 hover:text-white transition-colors cursor-pointer py-1">src/app/layout.tsx</Badge>
                      <Badge variant="outline" className="border-white/10 bg-zinc-900 text-xs text-zinc-300 hover:bg-white/5 hover:text-white transition-colors cursor-pointer py-1">docker-compose.yml</Badge>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </section>

          {/* IDE Layout */}
          <div className="min-h-0 flex-1 rounded-2xl ring-1 ring-white/10 shadow-2xl relative group overflow-hidden">
            {/* Explorer glow */}
            <div className="absolute inset-0 bg-gradient-to-b from-white/[0.02] to-transparent pointer-events-none" />
            <ExplorerLayout repoId={id} className="h-full rounded-2xl border-none shadow-none bg-black/40 backdrop-blur-sm" />
          </div>

        </main>
      ) : (
        <div className="flex flex-1 items-center justify-center animate-in fade-in duration-1000">
          <div className="flex flex-col items-center text-center">
            <div className="relative mb-8 flex h-32 w-32 items-center justify-center">
              {/* Outer spinning rings */}
              <div className="absolute inset-0 rounded-full border border-blue-500/20 animate-[spin_4s_linear_infinite]" />
              <div className="absolute inset-2 rounded-full border border-purple-500/20 animate-[spin_3s_linear_infinite_reverse]" />
              
              {/* Glowing core */}
              <div className="relative flex h-24 w-24 items-center justify-center rounded-2xl border border-white/10 bg-gradient-to-br from-blue-500/10 to-purple-500/10 shadow-2xl backdrop-blur-xl">
                <div className="absolute inset-0 rounded-2xl bg-blue-500/20 blur-2xl animate-pulse" />
                <GitBranch className="relative z-10 h-10 w-10 text-blue-300 drop-shadow-[0_0_15px_rgba(96,165,250,0.5)]" />
              </div>
            </div>
            <h2 className="text-2xl font-bold tracking-tight text-white mb-3">Analyzing Workspace</h2>
            <div className="flex items-center gap-2 text-sm font-medium text-blue-400">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span>Parsing repository architecture...</span>
            </div>
            <p className="mt-4 text-sm text-zinc-500 max-w-md leading-relaxed">
              We are cloning the repository, mapping out its file structure, and preparing the interactive code explorer.
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
