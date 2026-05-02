"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { GitFork, Loader2 } from "lucide-react";

export default function Home() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;

    try {
      const parsedUrl = new URL(url);
      if (
        !parsedUrl.hostname.endsWith("github.com") ||
        !/^\/[^/]+\/[^/]+(\/|\.git)?$/.test(parsedUrl.pathname)
      ) {
        setError("Please enter a valid GitHub repository URL");
        return;
      }
    } catch {
      setError("Please enter a valid GitHub repository URL");
      return;
    }

    setLoading(true);
    setError("");
    try {
      const res = await fetch("/api/v1/repos", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url }),
      });

      if (!res.ok) {
        throw new Error("Failed to start analysis");
      }

      const data = await res.json();
      if (!data.id) {
        throw new Error("Invalid response from server");
      }
      router.push(`/analysis/${data.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-b from-background to-secondary/20">
      <Card className="w-full max-w-lg border-border/50 shadow-2xl backdrop-blur-sm bg-card/95">
        <CardHeader className="space-y-1 pb-8">
          <div className="flex justify-center mb-4">
            <div className="p-3 bg-primary/10 rounded-full ring-1 ring-primary/20">
              <GitFork className="w-8 h-8 text-primary" />
            </div>
          </div>

          <CardTitle className="text-3xl text-center font-bold tracking-tight">
            Repository Analyzer
          </CardTitle>
          <CardDescription className="text-center text-lg">
            Enter a GitHub repository URL to analyze its dependencies
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2 relative">
              <Input
                type="url"
                placeholder="https://github.com/username/repo"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                required
                className="h-12 text-base px-4 bg-background/50 border-border/50 focus-visible:ring-primary/50 transition-all duration-300"
              />
            </div>
            {error && (
              <div className="text-destructive text-sm font-medium animate-in fade-in slide-in-from-top-1 text-center">
                {error}
              </div>
            )}
            <Button
              type="submit"
              className="w-full h-12 text-base font-medium transition-all duration-300 hover:scale-[1.02] active:scale-[0.98]"
              disabled={loading || !url}
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                  Analyzing...
                </>
              ) : (
                "Analyze Repository"
              )}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}