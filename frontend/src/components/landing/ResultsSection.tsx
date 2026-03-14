"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { fetchResults } from "@/lib/api";
import { useI18n } from "@/lib/i18n";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

const subjects = [
  { key: "mathematics", label: "Mathematics" },
  { key: "physics", label: "Physics" },
  { key: "chemistry", label: "Chemistry" },
  { key: "biology", label: "Biology" },
];

const medalColors: Record<string, string> = {
  Gold: "bg-yellow-500/15 text-yellow-600 dark:text-yellow-400 border-yellow-500/30",
  Silver: "bg-slate-500/15 text-slate-600 dark:text-slate-300 border-slate-500/30",
  Bronze: "bg-orange-500/15 text-orange-600 dark:text-orange-400 border-orange-500/30",
};

export function ResultsSection() {
  const [activeSubject, setActiveSubject] = useState("mathematics");
  const { t } = useI18n();

  const { data: results, isLoading } = useQuery({
    queryKey: ["results", activeSubject],
    queryFn: () => fetchResults(activeSubject),
  });

  return (
    <section id="results" className="py-20 bg-background border-t border-border">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <div className="inline-flex items-center gap-2 rounded-full border border-amber-500/20 bg-amber-500/10 px-4 py-1.5 text-sm text-amber-600 dark:text-amber-400 mb-4">
            {t("nav.results")}
          </div>
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">{t("results.title")}</h2>
          <p className="text-muted-foreground max-w-md mx-auto">
            {t("results.desc")}
          </p>
        </div>

        <div className="flex flex-wrap justify-center gap-2 mb-8">
          {subjects.map((s) => (
            <button
              key={s.key}
              onClick={() => setActiveSubject(s.key)}
              className={`rounded-full px-6 py-2 text-sm font-medium transition-all ${
                activeSubject === s.key
                  ? "bg-gradient-to-r from-blue-500 to-indigo-500 text-foreground shadow-lg shadow-blue-500/25"
                  : "bg-card text-muted-foreground border border-border hover:bg-accent hover:text-foreground"
              }`}
            >
              {s.label}
            </button>
          ))}
        </div>

        <div className="rounded-2xl border border-border bg-card backdrop-blur-sm shadow-2xl overflow-hidden">
          {isLoading ? (
            <div className="p-8 text-center text-muted-foreground">{t("common.loading")}</div>
          ) : !results?.length ? (
            <div className="p-8 text-center text-muted-foreground">{t("results.no_results") || "Natijalar topilmadi"}</div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow className="bg-card border-b border-border">
                  <TableHead className="w-16 text-center text-muted-foreground">{t("results.rank")}</TableHead>
                  <TableHead className="text-muted-foreground">{t("results.name")}</TableHead>
                  <TableHead className="text-muted-foreground">{t("results.country")}</TableHead>
                  <TableHead className="text-center text-muted-foreground">{t("results.score")}</TableHead>
                  <TableHead className="text-center text-muted-foreground">{t("results.medal")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {results.slice(0, 10).map((r) => (
                  <TableRow key={`${r.rank}-${r.name}`} className="border-b border-border hover:bg-card transition-colors">
                    <TableCell className="text-center font-bold text-muted-foreground">{r.rank}</TableCell>
                    <TableCell className="font-medium text-foreground">{r.name}</TableCell>
                    <TableCell>
                      <span className="text-sm text-muted-foreground">{r.country}</span>
                    </TableCell>
                    <TableCell className="text-center font-semibold text-foreground">{r.score}</TableCell>
                    <TableCell className="text-center">
                      {r.medal ? (
                        <Badge className={`${medalColors[r.medal] || "bg-accent text-muted-foreground"} hover:opacity-90`}>
                          {r.medal}
                        </Badge>
                      ) : (
                        <span className="text-blue-200/30">&mdash;</span>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </div>
    </section>
  );
}
