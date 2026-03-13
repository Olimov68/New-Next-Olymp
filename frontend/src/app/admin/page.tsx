"use client";

import { useEffect, useState } from "react";
import { getAdminDashboard } from "@/lib/admin-api";
import { usePanelAuth } from "@/lib/panel-auth-context";
import { Card, CardContent } from "@/components/ui/card";
import {
  Users,
  Trophy,
  GraduationCap,
  Newspaper,
  Award,
  BarChart3,
  AlertCircle,
} from "lucide-react";

interface DashboardStats {
  total_olympiads: number;
  total_mock_tests: number;
  total_news: number;
  total_users: number;
  total_results: number;
  total_certificates: number;
}

export default function AdminDashboard() {
  const { staff, hasModuleAccess, loading: authLoading } = usePanelAuth();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    if (authLoading) return;
    getAdminDashboard()
      .then((res) => {
        const d = res.data || res;
        setStats({
          total_olympiads: d.total_olympiads || d.stats?.total_olympiads || 0,
          total_mock_tests: d.total_mock_tests || d.stats?.total_mock_tests || 0,
          total_news: d.total_news || d.stats?.total_news || 0,
          total_users: d.total_users || d.stats?.total_users || 0,
          total_results: d.total_results || d.stats?.total_results || 0,
          total_certificates: d.total_certificates || d.stats?.total_certificates || 0,
        });
      })
      .catch(() => setError(true))
      .finally(() => setLoading(false));
  }, [authLoading]);

  if (loading || authLoading) {
    return (
      <div className="space-y-6">
        <div className="h-8 w-48 bg-muted rounded animate-pulse" />
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="h-28 bg-muted rounded-lg animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  if (error || !stats) {
    return (
      <div className="flex items-center justify-center h-64 text-red-500 gap-2">
        <AlertCircle className="h-5 w-5" /> Ma&apos;lumot yuklanmadi
      </div>
    );
  }

  const statCards = [
    { label: "Olimpiadalar", value: stats.total_olympiads, icon: Trophy, color: "amber", module: "olympiads" },
    { label: "Mock testlar", value: stats.total_mock_tests, icon: GraduationCap, color: "green", module: "mock_tests" },
    { label: "Yangiliklar", value: stats.total_news, icon: Newspaper, color: "blue", module: "news" },
    { label: "Foydalanuvchilar", value: stats.total_users, icon: Users, color: "purple", module: "users" },
    { label: "Natijalar", value: stats.total_results, icon: BarChart3, color: "rose", module: "results" },
    { label: "Sertifikatlar", value: stats.total_certificates, icon: Award, color: "cyan", module: "certificates" },
  ];

  const colorMap: Record<string, { bg: string; icon: string; border: string }> = {
    amber: { bg: "bg-amber-50 dark:bg-amber-950/30", icon: "text-amber-600 bg-amber-100 dark:bg-amber-900/50", border: "border-amber-100 dark:border-amber-900" },
    green: { bg: "bg-green-50 dark:bg-green-950/30", icon: "text-green-600 bg-green-100 dark:bg-green-900/50", border: "border-green-100 dark:border-green-900" },
    blue: { bg: "bg-blue-50 dark:bg-blue-950/30", icon: "text-blue-600 bg-blue-100 dark:bg-blue-900/50", border: "border-blue-100 dark:border-blue-900" },
    purple: { bg: "bg-purple-50 dark:bg-purple-950/30", icon: "text-purple-600 bg-purple-100 dark:bg-purple-900/50", border: "border-purple-100 dark:border-purple-900" },
    rose: { bg: "bg-rose-50 dark:bg-rose-950/30", icon: "text-rose-600 bg-rose-100 dark:bg-rose-900/50", border: "border-rose-100 dark:border-rose-900" },
    cyan: { bg: "bg-cyan-50 dark:bg-cyan-950/30", icon: "text-cyan-600 bg-cyan-100 dark:bg-cyan-900/50", border: "border-cyan-100 dark:border-cyan-900" },
  };

  const visibleCards = statCards.filter((s) => hasModuleAccess(s.module));

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Dashboard</h1>
        {staff && (
          <p className="text-muted-foreground mt-1">
            Xush kelibsiz, {staff.full_name || staff.username}!
          </p>
        )}
      </div>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        {visibleCards.map((s) => {
          const colors = colorMap[s.color];
          return (
            <Card key={s.label} className={`border ${colors.border} ${colors.bg} shadow-none`}>
              <CardContent className="p-4">
                <div className="flex items-center justify-between mb-3">
                  <div className={`flex h-10 w-10 items-center justify-center rounded-lg ${colors.icon}`}>
                    <s.icon className="h-5 w-5" />
                  </div>
                  <span className="text-3xl font-bold text-foreground">{s.value}</span>
                </div>
                <p className="text-sm text-muted-foreground font-medium">{s.label}</p>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {visibleCards.length === 0 && (
        <div className="text-center py-12 text-muted-foreground">
          Sizga hech qanday modul ruxsati berilmagan
        </div>
      )}
    </div>
  );
}
