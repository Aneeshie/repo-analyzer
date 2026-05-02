"use client";

import { useEffect, useRef, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Loader2, Package, ArrowLeft, CheckCircle2, XCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

type RepoStatus = "pending" | "processing" | "completed" | "failed" | "analyzing" | "cloning" | "parsing" | string;

const normalizeStatus = (status: RepoStatus): "pending" | "processing" | "completed" | "failed" => {
  if (["analyzing", "cloning", "parsing", "processing"].includes(status as string)) {
    return "processing";
  }
  if (status === "completed" || status === "failed") {
    return status as "completed" | "failed";
  }
  return "pending";
};

interface Repo {
  id: string;
  url: string;
  status: RepoStatus;
  created_at: string;
}

interface Dependency {
  id: string;
  name: string;
  version: string;
  ecosystem: string;
  scope: string;
  source_file: string;
}

export default function AnalysisPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [repo, setRepo] = useState<Repo | null>(null);
  const [dependencies, setDependencies] = useState<Dependency[]>([]);
  const [loadingDeps, setLoadingDeps] = useState(true);
  const [depsError, setDepsError] = useState<string | null>(null);
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

          if (data.status === "completed") {
            fetchDependencies();
          }
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Error fetching status");
        if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
      }
    };

    const fetchDependencies = async () => {
      try {
        setLoadingDeps(true);
        setDepsError(null);
        const res = await fetch(`/api/v1/repos/${id}/dependencies`);
        if (!res.ok) {
          throw new Error("Failed to fetch dependencies");
        }

        const data = await res.json();
        // Handle both possible structures: array or { dependencies: [...] }, and null when empty
        const depsList = data ? (Array.isArray(data) ? data : data.dependencies || []) : [];
        setDependencies(depsList);
      } catch (err) {
        setDepsError(err instanceof Error ? err.message : "Error fetching dependencies");
        setDependencies([]);
      } finally {
        setLoadingDeps(false);
      }
    };

    pollStatus(); // initial check
    pollIntervalRef.current = setInterval(pollStatus, 2000);

    return () => {
      if (pollIntervalRef.current) clearInterval(pollIntervalRef.current);
    };
  }, [id]);

  const groupedDependencies = dependencies.reduce((acc, dep) => {
    const ecosystem = dep.ecosystem || "unknown";
    if (!acc[ecosystem]) {
      acc[ecosystem] = [];
    }
    acc[ecosystem].push(dep);
    return acc;
  }, {} as Record<string, Dependency[]>);

  if (error) {
    return (
      <div className="min-h-screen p-8 flex flex-col items-center justify-center bg-background">
        <XCircle className="w-16 h-16 text-destructive mb-4" />
        <h1 className="text-2xl font-bold text-foreground mb-2">Analysis Failed</h1>
        <p className="text-muted-foreground mb-6">{error}</p>
        <Button onClick={() => router.push("/")} variant="outline" className="h-12 px-6">
          <ArrowLeft className="mr-2 h-4 w-4" /> Go Back
        </Button>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background/50 p-4 md:p-8">
      <div className="max-w-5xl mx-auto space-y-8">
        <Button onClick={() => router.push("/")} variant="ghost" className="mb-4">
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Home
        </Button>

        <Card className="border-border/50 bg-card/50 backdrop-blur-sm shadow-lg">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-2xl font-bold">Analysis Results</CardTitle>
                <CardDescription className="text-base mt-1">
                  {repo ? repo.url : <Skeleton className="h-4 w-64 mt-2" />}
                </CardDescription>
              </div>
              <div>
                {(() => {
                  const displayStatus = normalizeStatus(repo?.status || "pending");
                  if (displayStatus === "completed") return (
                    <Badge className="bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20 px-3 py-1.5 text-sm font-medium">
                      <CheckCircle2 className="w-4 h-4 mr-1.5" /> Completed
                    </Badge>
                  );
                  if (displayStatus === "failed") return (
                    <Badge variant="destructive" className="px-3 py-1.5 text-sm font-medium">
                      <XCircle className="w-4 h-4 mr-1.5" /> Failed
                    </Badge>
                  );
                  return (
                    <Badge variant="secondary" className="px-3 py-1.5 bg-blue-500/10 text-blue-500 hover:bg-blue-500/20 text-sm font-medium">
                      <Loader2 className="w-4 h-4 mr-1.5 animate-spin" />
                      {displayStatus === "processing" ? "Processing..." : "Pending..."}
                    </Badge>
                  );
                })()}
              </div>
            </div>
          </CardHeader>
        </Card>

        {repo?.status === "completed" && (
          <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <h2 className="text-2xl font-semibold tracking-tight flex items-center">
              <Package className="w-6 h-6 mr-2 text-primary" />
              Dependencies
            </h2>

            {loadingDeps ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Skeleton className="h-[300px] rounded-xl" />
                <Skeleton className="h-[300px] rounded-xl" />
              </div>
            ) : depsError ? (
              <Card className="border-dashed border-destructive/50 bg-destructive/5 shadow-none">
                <CardContent className="flex flex-col items-center justify-center p-16 text-center">
                  <XCircle className="w-16 h-16 text-destructive mb-4 opacity-80" />
                  <p className="text-xl font-medium text-destructive mb-2">Failed to load dependencies</p>
                  <p className="text-muted-foreground max-w-sm mt-1 leading-relaxed">
                    {depsError}
                  </p>
                </CardContent>
              </Card>
            ) : Object.keys(groupedDependencies).length > 0 ? (
              <div className="grid grid-cols-1 gap-6">
                {Object.entries(groupedDependencies).map(([ecosystem, deps]) => (
                  <Card key={ecosystem} className="overflow-hidden border-border/50 shadow-md transition-all duration-300 hover:shadow-lg">
                    <CardHeader className="bg-muted/30 border-b border-border/50">
                      <CardTitle className="flex items-center text-xl">
                        <Badge variant="outline" className="mr-3 font-mono text-sm capitalize px-3 py-1">
                          {ecosystem}
                        </Badge>
                        <span className="text-muted-foreground text-sm font-normal ml-auto">
                          {deps.length} package{deps.length === 1 ? '' : 's'}
                        </span>
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="p-0">
                      <div className="max-h-[400px] overflow-auto">
                        <Table>
                          <TableHeader className="bg-muted/50 sticky top-0 backdrop-blur-sm z-10">
                            <TableRow>
                              <TableHead className="w-1/2">Package Name</TableHead>
                              <TableHead>Version</TableHead>
                              <TableHead className="text-right">Scope</TableHead>
                            </TableRow>
                          </TableHeader>
                          <TableBody>
                            {deps.map((dep, i) => (
                              <TableRow key={`${dep.name}-${i}`} className="hover:bg-muted/20 transition-colors">
                                <TableCell className="font-medium text-foreground">{dep.name}</TableCell>
                                <TableCell className="font-mono text-muted-foreground text-sm">
                                  {dep.version || "latest"}
                                </TableCell>
                                <TableCell className="text-right">
                                  {dep.scope ? (
                                    <Badge variant="secondary" className="text-xs font-normal">
                                      {dep.scope}
                                    </Badge>
                                  ) : (
                                    <span className="text-muted-foreground text-sm">-</span>
                                  )}
                                </TableCell>
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            ) : (
              <Card className="border-dashed bg-transparent shadow-none">
                <CardContent className="flex flex-col items-center justify-center p-16 text-center">
                  <Package className="w-16 h-16 text-muted-foreground mb-4 opacity-30" />
                  <p className="text-xl font-medium text-foreground mb-2">No dependencies found</p>
                  <p className="text-muted-foreground max-w-sm mt-1 leading-relaxed">
                    We couldn&apos;t detect any dependencies for this repository. It might be empty or using an unsupported package manager.
                  </p>
                </CardContent>
              </Card>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
