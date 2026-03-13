"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { AnnouncementBar } from "@/components/AnnouncementBar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select";
import {
  Trophy, Search, Clock, Users, BookOpen, Calendar, ChevronRight,
} from "lucide-react";
import { listOlympiads } from "@/lib/user-api";
import { useAuth } from "@/lib/auth-context";

interface Olympiad {
  id: number;
  title: string;
  slug: string;
  description: string;
  subject: string;
  grade: number;
  language: string;
  start_time: string;
  end_time: string;
  duration_minutes: number;
  total_questions: number;
  status: string;
  is_paid: boolean;
  price?: number;
}

const statusColors: Record<string, string> = {
  published: "bg-blue-500/10 text-blue-600 border-blue-500/20",
  active: "bg-green-500/10 text-green-600 border-green-500/20",
  ended: "bg-gray-500/10 text-gray-500 border-gray-500/20",
};
const statusLabels: Record<string, string> = {
  published: "Ro'yxatdan o'tish",
  active: "Faol",
  ended: "Tugagan",
};

export default function OlympiadsPage() {
  const router = useRouter();
  const { user } = useAuth();
  const [items, setItems] = useState<Olympiad[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [subjectFilter, setSubjectFilter] = useState("all");
  const [priceFilter, setPriceFilter] = useState("all");

  useEffect(() => {
    listOlympiads({ page: 1, page_size: 50 })
      .then(data => setItems(Array.isArray(data) ? data : (data as any)?.data || []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false));
  }, []);

  const subjects = Array.from(new Set(items.map(i => i.subject).filter(Boolean)));

  const filtered = items.filter(i => {
    const matchSearch = !search || i.title.toLowerCase().includes(search.toLowerCase()) || i.subject.toLowerCase().includes(search.toLowerCase());
    const matchSubject = subjectFilter === "all" || i.subject === subjectFilter;
    const matchPrice = priceFilter === "all" || (priceFilter === "free" ? !i.is_paid : i.is_paid);
    return matchSearch && matchSubject && matchPrice;
  });

  const handleJoin = (id: number) => {
    if (!user) { router.push("/login"); return; }
    router.push(`/dashboard/olympiads/${id}`);
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="sticky top-0 z-50"><AnnouncementBar /><Header /></div>
      <main>
        {/* Hero */}
        <section className="py-14 border-b border-border bg-gradient-to-b from-primary/5 to-background">
          <div className="container mx-auto px-4 text-center">
            <div className="inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/10 px-4 py-1.5 text-sm text-primary mb-5">
              <Trophy className="h-4 w-4" />Olimpiadalar
            </div>
            <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">Olimpiadalar va musobaqalar</h1>
            <p className="text-muted-foreground max-w-lg mx-auto">Bilimingizni sinab ko&apos;ring va g&apos;olib bo&apos;ling</p>
          </div>
        </section>

        <section className="py-12">
          <div className="container mx-auto px-4">
            {/* Filters */}
            <div className="flex flex-wrap gap-3 mb-8">
              <div className="relative flex-1 min-w-[200px]">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input placeholder="Qidirish..." value={search} onChange={e => setSearch(e.target.value)} className="pl-9" />
              </div>
              <Select value={subjectFilter} onValueChange={(v) => setSubjectFilter(v ?? "all")}>
                <SelectTrigger className="w-44"><SelectValue placeholder="Fan" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Barcha fanlar</SelectItem>
                  {subjects.map(s => <SelectItem key={s} value={s}>{s}</SelectItem>)}
                </SelectContent>
              </Select>
              <Select value={priceFilter} onValueChange={(v) => setPriceFilter(v ?? "all")}>
                <SelectTrigger className="w-36"><SelectValue placeholder="To'lov" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Barchasi</SelectItem>
                  <SelectItem value="free">Bepul</SelectItem>
                  <SelectItem value="paid">To&apos;lovli</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Cards */}
            {loading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {Array.from({ length: 6 }).map((_, i) => <Skeleton key={i} className="h-64 rounded-2xl" />)}
              </div>
            ) : filtered.length === 0 ? (
              <div className="flex flex-col items-center py-24 text-muted-foreground">
                <Trophy className="h-14 w-14 mb-4 opacity-20" />
                <p className="text-lg font-medium">Olimpiadalar topilmadi</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filtered.map(item => (
                  <div key={item.id} className="group rounded-2xl border border-border bg-card p-6 flex flex-col hover:shadow-lg transition-all duration-300">
                    <div className="flex items-start justify-between mb-4">
                      <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10">
                        <Trophy className="h-6 w-6 text-primary" />
                      </div>
                      <div className="flex gap-2 flex-wrap justify-end">
                        {item.status in statusColors && (
                          <Badge variant="outline" className={`text-xs ${statusColors[item.status]}`}>
                            {statusLabels[item.status]}
                          </Badge>
                        )}
                        <Badge variant="outline" className={`text-xs ${item.is_paid ? "bg-orange-500/10 text-orange-600 border-orange-500/20" : "bg-green-500/10 text-green-600 border-green-500/20"}`}>
                          {item.is_paid ? `${item.price?.toLocaleString()} UZS` : "Bepul"}
                        </Badge>
                      </div>
                    </div>

                    <h3 className="font-bold text-foreground text-base mb-2 line-clamp-2 group-hover:text-primary transition-colors flex-1">{item.title}</h3>

                    {item.description && <p className="text-sm text-muted-foreground line-clamp-2 mb-4">{item.description}</p>}

                    <div className="grid grid-cols-2 gap-2 text-xs text-muted-foreground mb-4">
                      <span className="flex items-center gap-1.5"><BookOpen className="h-3.5 w-3.5" />{item.subject}</span>
                      <span className="flex items-center gap-1.5"><Clock className="h-3.5 w-3.5" />{item.duration_minutes} daqiqa</span>
                      {item.grade > 0 && <span className="flex items-center gap-1.5"><Users className="h-3.5 w-3.5" />{item.grade}-sinf</span>}
                      {item.total_questions > 0 && <span className="flex items-center gap-1.5"><Trophy className="h-3.5 w-3.5" />{item.total_questions} savol</span>}
                    </div>

                    {item.end_time && (
                      <p className="text-xs text-muted-foreground flex items-center gap-1 mb-4">
                        <Calendar className="h-3.5 w-3.5" />
                        Tugash: {new Date(item.end_time).toLocaleDateString("uz-UZ")}
                      </p>
                    )}

                    <Button
                      className="w-full gap-2 mt-auto"
                      onClick={() => handleJoin(item.id)}
                      variant={item.status === "ended" ? "outline" : "default"}
                      disabled={item.status === "ended"}
                    >
                      {item.status === "ended" ? "Tugagan" : "Batafsil"}
                      {item.status !== "ended" && <ChevronRight className="h-4 w-4" />}
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </section>
      </main>
      <Footer />
    </div>
  );
}
