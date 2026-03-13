"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Card, CardContent } from "@/components/ui/card";
import {
  Users,
  ShieldCheck,
  Trophy,
  GraduationCap,
  MessageSquare,
  Award,
  CreditCard,
  AlertCircle,
} from "lucide-react";

interface DashboardStats {
  total_users: number;
  blocked_users: number;
  total_admins: number;
  total_olympiads: number;
  total_mock_tests: number;
  total_feedbacks: number;
  open_feedbacks: number;
  total_payments: number;
  total_certificates: number;
}

interface LatestUser {
  id: number;
  username: string;
  status: string;
  created_at: string;
}

interface LatestFeedback {
  id: number;
  username: string;
  subject: string;
  status: string;
  created_at: string;
}

export default function SuperAdminDashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [latestUsers, setLatestUsers] = useState<LatestUser[]>([]);
  const [latestFeedbacks, setLatestFeedbacks] = useState<LatestFeedback[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("panel_access_token");
    if (!token) return;

    api
      .get("/superadmin/dashboard", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((r) => {
        const d = r.data.data;
        setStats(d.stats);
        setLatestUsers(d.latest_users || []);
        setLatestFeedbacks(d.latest_feedbacks || []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64 text-muted-foreground">
        Yuklanmoqda...
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="flex items-center justify-center h-64 text-red-400 gap-2">
        <AlertCircle className="h-5 w-5" /> Ma'lumot yuklanmadi
      </div>
    );
  }

  const statCards = [
    { label: "Foydalanuvchilar", value: stats.total_users, icon: Users, color: "blue" },
    { label: "Adminlar", value: stats.total_admins, icon: ShieldCheck, color: "purple" },
    { label: "Olimpiadalar", value: stats.total_olympiads, icon: Trophy, color: "amber" },
    { label: "Mock testlar", value: stats.total_mock_tests, icon: GraduationCap, color: "green" },
    { label: "Fikrlar", value: stats.total_feedbacks, sub: `${stats.open_feedbacks} ochiq`, icon: MessageSquare, color: "rose" },
    { label: "Sertifikatlar", value: stats.total_certificates, icon: Award, color: "cyan" },
    { label: "To'lovlar", value: stats.total_payments, icon: CreditCard, color: "emerald" },
    { label: "Bloklangan", value: stats.blocked_users, icon: AlertCircle, color: "red" },
  ];

  const colorMap: Record<string, string> = {
    blue: "from-blue-500/20 to-blue-600/5 border-blue-400/20 text-blue-400",
    purple: "from-purple-500/20 to-purple-600/5 border-purple-400/20 text-purple-400",
    amber: "from-amber-500/20 to-amber-600/5 border-amber-400/20 text-amber-400",
    green: "from-green-500/20 to-green-600/5 border-green-400/20 text-green-400",
    rose: "from-rose-500/20 to-rose-600/5 border-rose-400/20 text-rose-400",
    cyan: "from-cyan-500/20 to-cyan-600/5 border-cyan-400/20 text-cyan-400",
    emerald: "from-emerald-500/20 to-emerald-600/5 border-emerald-400/20 text-emerald-400",
    red: "from-red-500/20 to-red-600/5 border-red-400/20 text-red-400",
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        {statCards.map((s) => (
          <Card key={s.label} className={`border bg-gradient-to-br ${colorMap[s.color]} shadow-none`}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-2">
                <s.icon className="h-5 w-5 opacity-70" />
                <span className="text-2xl font-bold text-foreground">{s.value}</span>
              </div>
              <p className="text-xs opacity-60">{s.label}</p>
              {s.sub && <p className="text-[10px] opacity-40 mt-0.5">{s.sub}</p>}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Latest data */}
      <div className="grid md:grid-cols-2 gap-4">
        {/* Latest Users */}
        <Card className="border border-border bg-accent/50 shadow-none">
          <CardContent className="p-4">
            <h3 className="text-sm font-semibold text-muted-foreground mb-3">Oxirgi foydalanuvchilar</h3>
            <div className="space-y-2">
              {latestUsers.length === 0 && <p className="text-xs text-muted-foreground">Hali yo'q</p>}
              {latestUsers.map((u) => (
                <div key={u.id} className="flex items-center justify-between py-1.5 border-b border-border last:border-0">
                  <div>
                    <span className="text-sm text-foreground">{u.username}</span>
                    <span className={`ml-2 text-[10px] px-1.5 py-0.5 rounded-full ${u.status === "active" ? "bg-green-500/10 text-green-400" : "bg-red-500/10 text-red-400"}`}>
                      {u.status}
                    </span>
                  </div>
                  <span className="text-[10px] text-muted-foreground">{new Date(u.created_at).toLocaleDateString()}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Latest Feedbacks */}
        <Card className="border border-border bg-accent/50 shadow-none">
          <CardContent className="p-4">
            <h3 className="text-sm font-semibold text-muted-foreground mb-3">Oxirgi fikrlar</h3>
            <div className="space-y-2">
              {latestFeedbacks.length === 0 && <p className="text-xs text-muted-foreground">Hali yo'q</p>}
              {latestFeedbacks.map((f) => (
                <div key={f.id} className="flex items-center justify-between py-1.5 border-b border-border last:border-0">
                  <div className="min-w-0">
                    <span className="text-sm text-foreground truncate block">{f.subject}</span>
                    <span className="text-[10px] text-muted-foreground">{f.username}</span>
                  </div>
                  <span className={`text-[10px] px-1.5 py-0.5 rounded-full flex-shrink-0 ${
                    f.status === "open" ? "bg-amber-500/10 text-amber-400" :
                    f.status === "answered" ? "bg-green-500/10 text-green-400" :
                    "bg-muted/50 text-muted-foreground"
                  }`}>
                    {f.status}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
