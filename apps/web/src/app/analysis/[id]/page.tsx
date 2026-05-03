"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  Loader2,
  ArrowLeft,
  XCircle,
  GitBranch
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
  entry_points?: string[];
}

export default function AnalysisPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [repo, setRepo] = useState<Repo | null>(null);
  const [error, setError] = useState<string | null>(null);

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
    <div className="flex h-screen w-screen flex-col bg-[#1e1e1e] text-zinc-200 overflow-hidden">
      {displayStatus === "completed" ? (
        <main className="flex-1 overflow-hidden bg-[#1e1e1e]">
          <ExplorerLayout repo={repo!} className="h-full w-full border-none shadow-none rounded-none bg-transparent" />
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
