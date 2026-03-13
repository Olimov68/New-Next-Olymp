"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { AnnouncementBar } from "@/components/AnnouncementBar";
import { fetchResults } from "@/lib/api";
import { useI18n } from "@/lib/i18n";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { BarChart3, Medal, Loader2, AlertCircle, Search } from "lucide-react";
import { Input } from "@/components/ui/input";

const subjects = [
  { key: "mathematics", uz: "Matematika", ru: "Математика", en: "Mathematics" },
  { key: "physics", uz: "Fizika", ru: "Физика", en: "Physics" },
  { key: "chemistry", uz: "Kimyo", ru: "Химия", en: "Chemistry" },
  { key: "biology", uz: "Biologiya", ru: "Биология", en: "Biology" },
  { key: "informatics", uz: "Informatika", ru: "Информатика", en: "Informatics" },
];

const medalConfig: Record<string, { cls: string; icon: string }> = {
  Gold:   { cls: "bg-yellow-500/15 text-yellow-600 dark:text-yellow-400 border-yellow-500/30", icon: "🥇" },
  Silver: { cls: "bg-slate-500/15 text-slate-600 dark:text-slate-300 border-slate-500/30", icon: "🥈" },
  Bronze: { cls: "bg-orange-500/15 text-orange-600 dark:text-orange-400 border-orange-500/30", icon: "🥉" },
};

const rankStyle = (rank: number) => {
  if (rank === 1) return "text-yellow-500 font-bold text-lg";
  if (rank === 2) return "text-slate-500 font-bold text-lg";
  if (rank === 3) return "text-orange-500 font-bold text-lg";
  return "text-muted-foreground font-medium";
};

export default function ResultsPage() {
  const { lang } = useI18n();
  const [activeSubject, setActiveSubject] = useState("mathematics");
  const [search, setSearch] = useState("");

  const { data: results = [], isLoading, isError } = useQuery({
    queryKey: ["results", activeSubject],
    queryFn: () => fetchResults(activeSubject),
  });

  const getLabel = (s: typeof subjects[0]) =>
    lang === "ru" ? s.ru : lang === "en" ? s.en : s.uz;

  const filtered = results.filter(
    (r) => !search || r.name?.toLowerCase().includes(search.toLowerCase())
  );

  const title = lang === "ru" ? "Результаты олимпиад" : lang === "en" ? "Olympiad Results" : "Olimpiada natijalari";
  const desc = lang === "ru"
    ? "Таблица лидеров по предметам — посмотрите результаты участников"
    : lang === "en"
    ? "Leaderboard by subject — see participant results"
    : "Fan bo'yicha reytinglar — ishtirokchilar natijalarini ko'ring";

  return (
    <div className="min-h-screen bg-background">
      <div className="sticky top-0 z-50">
        <AnnouncementBar />
        <Header />
      </div>

      <main>
        {/* Page Header */}
        <section className="py-14 border-b border-border bg-gradient-to-b from-primary/5 to-background">
          <div className="container mx-auto px-4 text-center">
            <div className="inline-flex items-center gap-2 rounded-full border border-amber-500/20 bg-amber-500/10 px-4 py-1.5 text-sm text-amber-600 dark:text-amber-400 mb-4">
              <BarChart3 className="h-4 w-4" />
              {lang === "ru" ? "Результаты" : lang === "en" ? "Results" : "Natijalar"}
            </div>
            <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">{title}</h1>
            <p className="text-muted-foreground max-w-md mx-auto">{desc}</p>
          </div>
        </section>

        {/* Subject tabs */}
        <section className="py-6 border-b border-border bg-background sticky top-[calc(4rem+var(--announcement-h,0px))] z-30 backdrop-blur-sm bg-background/80">
          <div className="container mx-auto px-4">
            <div className="flex flex-col sm:flex-row gap-3 items-stretch sm:items-center">
              {/* Search */}
              <div className="relative max-w-xs">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder={lang === "ru" ? "Поиск по имени..." : lang === "en" ? "Search by name..." : "Ism bo'yicha qidirish..."}
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="pl-9 bg-background"
                />
              </div>

              {/* Subject buttons */}
              <div className="flex flex-wrap gap-2">
                {subjects.map((s) => (
                  <button
                    key={s.key}
                    onClick={() => setActiveSubject(s.key)}
                    className={`rounded-full px-4 py-1.5 text-sm font-medium transition-all ${
                      activeSubject === s.key
                        ? "bg-primary text-primary-foreground shadow-sm"
                        : "bg-muted text-muted-foreground hover:text-foreground hover:bg-accent"
                    }`}
                  >
                    {getLabel(s)}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* Results table */}
        <section className="py-10">
          <div className="container mx-auto px-4">
            {isLoading ? (
              <div className="flex justify-center py-20">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
              </div>
            ) : isError ? (
              <div className="flex flex-col items-center py-20 gap-4">
                <AlertCircle className="h-12 w-12 text-destructive/40" />
                <p className="text-muted-foreground">
                  {lang === "ru" ? "Ошибка загрузки" : lang === "en" ? "Failed to load" : "Yuklab bo'lmadi"}
                </p>
              </div>
            ) : filtered.length === 0 ? (
              <div className="flex flex-col items-center py-20 gap-4">
                <Medal className="h-16 w-16 text-muted-foreground/20" />
                <p className="text-muted-foreground text-lg">
                  {lang === "ru" ? "Результатов нет" : lang === "en" ? "No results found" : "Natijalar topilmadi"}
                </p>
              </div>
            ) : (
              <div className="rounded-2xl border border-border bg-card overflow-hidden shadow-sm">
                <Table>
                  <TableHeader>
                    <TableRow className="border-b border-border hover:bg-transparent">
                      <TableHead className="w-16 text-center text-xs uppercase tracking-wider text-muted-foreground">
                        {lang === "ru" ? "Место" : lang === "en" ? "Rank" : "O'rin"}
                      </TableHead>
                      <TableHead className="text-xs uppercase tracking-wider text-muted-foreground">
                        {lang === "ru" ? "Участник" : lang === "en" ? "Participant" : "Ishtirokchi"}
                      </TableHead>
                      <TableHead className="text-xs uppercase tracking-wider text-muted-foreground">
                        {lang === "ru" ? "Страна/Регион" : lang === "en" ? "Country/Region" : "Viloyat"}
                      </TableHead>
                      <TableHead className="text-center text-xs uppercase tracking-wider text-muted-foreground">
                        {lang === "ru" ? "Балл" : lang === "en" ? "Score" : "Ball"}
                      </TableHead>
                      <TableHead className="text-center text-xs uppercase tracking-wider text-muted-foreground">
                        {lang === "ru" ? "Медаль" : lang === "en" ? "Medal" : "Medal"}
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filtered.map((r, idx) => (
                      <TableRow
                        key={`${r.rank}-${r.name}-${idx}`}
                        className={`border-b border-border transition-colors hover:bg-muted/50 ${r.rank <= 3 ? "bg-amber-500/[0.02]" : ""}`}
                      >
                        <TableCell className={`text-center ${rankStyle(r.rank)}`}>
                          {r.rank <= 3 ? ["🥇", "🥈", "🥉"][r.rank - 1] : `#${r.rank}`}
                        </TableCell>
                        <TableCell className="font-medium text-foreground">{r.name}</TableCell>
                        <TableCell className="text-sm text-muted-foreground">{r.country}</TableCell>
                        <TableCell className="text-center font-semibold text-foreground">{r.score}</TableCell>
                        <TableCell className="text-center">
                          {r.medal ? (
                            <Badge
                              variant="outline"
                              className={`${medalConfig[r.medal]?.cls || "bg-muted text-muted-foreground"} text-xs`}
                            >
                              {medalConfig[r.medal]?.icon} {r.medal}
                            </Badge>
                          ) : (
                            <span className="text-muted-foreground/30 text-sm">—</span>
                          )}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
                <div className="px-6 py-3 border-t border-border text-xs text-muted-foreground">
                  {lang === "ru" ? `Всего участников: ${filtered.length}` : lang === "en" ? `Total participants: ${filtered.length}` : `Jami ishtirokchilar: ${filtered.length}`}
                </div>
              </div>
            )}
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
