"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";
import {
  Trophy,
  FileText,
  Award,
  Bell,
  Wallet,
  Calendar,
  TrendingUp,
  Newspaper,
  ChevronRight,
} from "lucide-react";
import {
  getBalance,
  getMyOlympiads,
  getMyMockTests,
  listCertificates,
  getUnreadCount,
  getMyResults,
  listNews,
  type BalanceInfo,
  type UserResult,
} from "@/lib/user-api";
import type { Olympiad, MockExam } from "@/lib/api";
import type { Certificate } from "@/lib/user-api";

interface DashboardNewsItem {
  id: number;
  title: string;
  cover_image?: string;
  excerpt?: string;
  created_at: string;
}

export default function DashboardPage() {
  const { user } = useAuth();
  const [balance, setBalance] = useState<BalanceInfo | null>(null);
  const [myOlympiads, setMyOlympiads] = useState<Olympiad[]>([]);
  const [myMockTests, setMyMockTests] = useState<MockExam[]>([]);
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [results, setResults] = useState<UserResult[]>([]);
  const [news, setNews] = useState<DashboardNewsItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [bal, olymp, mocks, certs, unread, res, newsData] =
          await Promise.allSettled([
            getBalance(),
            getMyOlympiads(),
            getMyMockTests(),
            listCertificates(),
            getUnreadCount(),
            getMyResults(),
            listNews(),
          ]);
        if (bal.status === "fulfilled") setBalance(bal.value);
        if (olymp.status === "fulfilled") setMyOlympiads(Array.isArray(olymp.value) ? olymp.value : []);
        if (mocks.status === "fulfilled") setMyMockTests(Array.isArray(mocks.value) ? mocks.value : []);
        if (certs.status === "fulfilled") setCertificates(Array.isArray(certs.value) ? certs.value : []);
        if (unread.status === "fulfilled") setUnreadCount(typeof unread.value === "number" ? unread.value : 0);
        if (res.status === "fulfilled") setResults(Array.isArray(res.value) ? res.value : []);
        if (newsData.status === "fulfilled") setNews(Array.isArray(newsData.value) ? newsData.value : []);
      } catch {
        // silently handle
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const upcomingExams = [
    ...myOlympiads
      .filter((o) => o.status === "active" || o.status === "upcoming")
      .map((o) => ({ id: o.id, title: o.title, type: "olympiad" as const, date: o.start_date, subject: o.subject })),
  ].sort((a, b) => (a.date && b.date ? new Date(a.date).getTime() - new Date(b.date).getTime() : 0));

  const latestResults = results.slice(0, 5);
  const latestNews = news.slice(0, 3);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Yuklanmoqda...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Welcome */}
      <div>
        <h1 className="text-2xl font-bold text-foreground">
          Xush kelibsiz, {user?.username}!
        </h1>
        <p className="text-muted-foreground mt-1">Shaxsiy kabinetingiz</p>
      </div>

      {/* Balance Card */}
      <Link href="/dashboard/balance">
        <Card className="bg-gradient-to-r from-blue-600 to-blue-700 border-0 shadow-lg hover:shadow-xl transition-shadow cursor-pointer">
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-blue-100 text-sm font-medium">Balans</p>
                <p className="text-3xl font-bold text-white mt-1">
                  {(balance?.balance ?? 0).toLocaleString()} so&apos;m
                </p>
              </div>
              <div className="h-12 w-12 rounded-xl bg-white/20 flex items-center justify-center">
                <Wallet className="h-6 w-6 text-white" />
              </div>
            </div>
          </CardContent>
        </Card>
      </Link>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Link href="/dashboard/olympiads">
          <Card className="border-0 shadow-sm hover:shadow-md transition-shadow cursor-pointer h-full">
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Olimpiadalarim</p>
                  <p className="text-2xl font-bold text-foreground mt-1">{myOlympiads.length}</p>
                </div>
                <div className="h-10 w-10 rounded-lg bg-purple-50 flex items-center justify-center">
                  <Trophy className="h-5 w-5 text-purple-600" />
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/dashboard/mock-tests">
          <Card className="border-0 shadow-sm hover:shadow-md transition-shadow cursor-pointer h-full">
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Mock testlar</p>
                  <p className="text-2xl font-bold text-foreground mt-1">{myMockTests.length}</p>
                </div>
                <div className="h-10 w-10 rounded-lg bg-green-50 flex items-center justify-center">
                  <FileText className="h-5 w-5 text-green-600" />
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/dashboard/certificates">
          <Card className="border-0 shadow-sm hover:shadow-md transition-shadow cursor-pointer h-full">
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Sertifikatlar</p>
                  <p className="text-2xl font-bold text-foreground mt-1">{certificates.length}</p>
                </div>
                <div className="h-10 w-10 rounded-lg bg-amber-50 flex items-center justify-center">
                  <Award className="h-5 w-5 text-amber-600" />
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>

        <Link href="/dashboard/notifications">
          <Card className="border-0 shadow-sm hover:shadow-md transition-shadow cursor-pointer h-full">
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">O&apos;qilmagan xabarlar</p>
                  <p className="text-2xl font-bold text-foreground mt-1">{unreadCount}</p>
                </div>
                <div className="h-10 w-10 rounded-lg bg-red-50 flex items-center justify-center">
                  <Bell className="h-5 w-5 text-red-600" />
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Two Column Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Upcoming Exams */}
        <Card className="border-0 shadow-sm">
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-foreground flex items-center gap-2">
                <Calendar className="h-5 w-5 text-blue-600" />
                Kelgusi imtihonlar
              </h2>
            </div>
            {upcomingExams.length === 0 ? (
              <p className="text-muted-foreground text-sm py-4 text-center">
                Hozircha kelgusi imtihon yo&apos;q
              </p>
            ) : (
              <div className="space-y-3">
                {upcomingExams.slice(0, 5).map((exam) => (
                  <Link
                    key={`${exam.type}-${exam.id}`}
                    href={
                      exam.type === "olympiad"
                        ? `/dashboard/olympiads/${exam.id}`
                        : `/dashboard/mock-tests/${exam.id}`
                    }
                    className="flex items-center justify-between p-3 rounded-lg bg-muted hover:bg-accent transition-colors"
                  >
                    <div>
                      <p className="text-sm font-medium text-foreground">{exam.title}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <Badge
                          className={`text-xs border-0 ${
                            exam.type === "olympiad"
                              ? "bg-purple-100 text-purple-700"
                              : "bg-green-100 text-green-700"
                          }`}
                        >
                          {exam.type === "olympiad" ? "Olimpiada" : "Mock test"}
                        </Badge>
                        <span className="text-xs text-muted-foreground">{exam.subject}</span>
                      </div>
                    </div>
                    {exam.date && (
                      <span className="text-xs text-muted-foreground">
                        {new Date(exam.date).toLocaleDateString("uz-UZ")}
                      </span>
                    )}
                  </Link>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Latest Results */}
        <Card className="border-0 shadow-sm">
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-foreground flex items-center gap-2">
                <TrendingUp className="h-5 w-5 text-green-600" />
                So&apos;nggi natijalar
              </h2>
              <Link
                href="/dashboard/results"
                className="text-sm text-primary hover:text-primary/80 flex items-center gap-1"
              >
                Barchasi <ChevronRight className="h-4 w-4" />
              </Link>
            </div>
            {latestResults.length === 0 ? (
              <p className="text-muted-foreground text-sm py-4 text-center">
                Hozircha natija yo&apos;q
              </p>
            ) : (
              <div className="space-y-3">
                {latestResults.map((r) => (
                  <div
                    key={r.id}
                    className="flex items-center justify-between p-3 rounded-lg bg-muted"
                  >
                    <div>
                      <p className="text-sm font-medium text-foreground">{r.title}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <Badge
                          className={`text-xs border-0 ${
                            r.type === "olympiad"
                              ? "bg-purple-100 text-purple-700"
                              : "bg-green-100 text-green-700"
                          }`}
                        >
                          {r.type === "olympiad" ? "Olimpiada" : "Mock test"}
                        </Badge>
                        <span className="text-xs text-muted-foreground">{r.subject}</span>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-bold text-foreground">
                        {r.score}/{r.max_score}
                      </p>
                      <p className="text-xs text-muted-foreground">{r.percentage}%</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Latest News */}
      <Card className="border-0 shadow-sm">
        <CardContent className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-foreground flex items-center gap-2">
              <Newspaper className="h-5 w-5 text-blue-600" />
              So&apos;nggi yangiliklar
            </h2>
            <Link
              href="/dashboard/news"
              className="text-sm text-primary hover:text-primary/80 flex items-center gap-1"
            >
              Barchasi <ChevronRight className="h-4 w-4" />
            </Link>
          </div>
          {latestNews.length === 0 ? (
            <p className="text-muted-foreground text-sm py-4 text-center">
              Hozircha yangilik yo&apos;q
            </p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {latestNews.map((n) => (
                <Link key={n.id} href={`/dashboard/news`}>
                  <div className="rounded-lg border border-border bg-card p-4 hover:shadow-md transition-shadow cursor-pointer">
                    {n.cover_image && (
                      <img
                        src={n.cover_image}
                        alt={n.title}
                        className="w-full h-32 object-cover rounded-lg mb-3"
                      />
                    )}
                    <h3 className="text-sm font-semibold text-foreground line-clamp-2">
                      {n.title}
                    </h3>
                    <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                      {n.excerpt}
                    </p>
                    <p className="text-xs text-muted-foreground mt-2">
                      {new Date(n.created_at).toLocaleDateString("uz-UZ")}
                    </p>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
