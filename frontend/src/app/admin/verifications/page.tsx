"use client";

import { useEffect, useState, useCallback } from "react";
import { getVerifications, approveVerification, rejectVerification } from "@/lib/admin-api";
import { UserCheck, UserX, Search, Loader2, CheckCircle, XCircle, Clock, Eye } from "lucide-react";
import { toast } from "sonner";

interface VerificationItem {
  id: number;
  user_id: number;
  username: string;
  first_name: string;
  last_name: string;
  photo_url: string;
  region: string;
  district: string;
  school_name: string;
  grade: number;
  method: string;
  status: string;
  note: string;
  reason: string;
  registered_at: string;
  created_at: string;
  verified_at: string | null;
}

export default function AdminVerificationsPage() {
  const [items, setItems] = useState<VerificationItem[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("pending");
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [approveDialog, setApproveDialog] = useState<VerificationItem | null>(null);
  const [rejectDialog, setRejectDialog] = useState<VerificationItem | null>(null);
  const [approveNote, setApproveNote] = useState("");
  const [rejectReason, setRejectReason] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [detailItem, setDetailItem] = useState<VerificationItem | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await getVerifications({ page, page_size: 20, status: filter, search: search || undefined });
      const d = res?.data || res;
      setItems(d?.data || []);
      setTotal(d?.total || 0);
    } catch { /* ignore */ }
    setLoading(false);
  }, [page, filter, search]);

  useEffect(() => { load(); }, [load]);

  const handleApprove = async () => {
    if (!approveDialog || submitting) return;
    setSubmitting(true);
    try {
      await approveVerification(approveDialog.id, { note: approveNote });
      toast.success(`${approveDialog.first_name || approveDialog.username} tasdiqlandi`);
      setApproveDialog(null);
      setApproveNote("");
      load();
    } catch {
      toast.error("Xatolik yuz berdi");
    }
    setSubmitting(false);
  };

  const handleReject = async () => {
    if (!rejectDialog || submitting || !rejectReason.trim()) return;
    setSubmitting(true);
    try {
      await rejectVerification(rejectDialog.id, { reason: rejectReason });
      toast.success("Rad etildi");
      setRejectDialog(null);
      setRejectReason("");
      load();
    } catch {
      toast.error("Xatolik yuz berdi");
    }
    setSubmitting(false);
  };

  const formatDate = (s: string) => {
    try { return new Date(s).toLocaleDateString("uz", { day: "numeric", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" }); }
    catch { return s; }
  };

  const statusBadge = (status: string) => {
    switch (status) {
      case "pending": return "bg-yellow-500/10 text-yellow-400";
      case "approved": return "bg-emerald-500/10 text-emerald-400";
      case "rejected": return "bg-red-500/10 text-red-400";
      default: return "bg-gray-500/10 text-gray-400";
    }
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Foydalanuvchi Tasdiqlash</h1>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="relative flex-1 max-w-xs">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1); }}
            placeholder="Qidirish..."
            className="w-full pl-9 pr-3 py-2 bg-background border border-border rounded-lg text-sm"
          />
        </div>
        <div className="flex gap-1 bg-accent/50 p-1 rounded-lg">
          {[
            { key: "pending", label: "Kutilmoqda", icon: Clock },
            { key: "approved", label: "Tasdiqlangan", icon: CheckCircle },
            { key: "rejected", label: "Rad etilgan", icon: XCircle },
            { key: "all", label: "Hammasi", icon: Eye },
          ].map((f) => (
            <button
              key={f.key}
              onClick={() => { setFilter(f.key); setPage(1); }}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs transition-colors ${
                filter === f.key ? "bg-background text-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"
              }`}
            >
              <f.icon className="w-3 h-3" /> {f.label}
            </button>
          ))}
        </div>
      </div>

      {/* Table */}
      {loading ? (
        <div className="flex justify-center py-12"><Loader2 className="w-8 h-8 animate-spin text-muted-foreground" /></div>
      ) : items.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          {filter === "pending" ? "Kutilayotgan tasdiqlar yo'q" : "Natijalar topilmadi"}
        </div>
      ) : (
        <div className="border border-border rounded-xl bg-card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-muted-foreground">
                  <th className="px-4 py-3 font-medium">Foydalanuvchi</th>
                  <th className="px-4 py-3 font-medium">Ma&apos;lumot</th>
                  <th className="px-4 py-3 font-medium">Ro&apos;yxatdan</th>
                  <th className="px-4 py-3 font-medium">Status</th>
                  <th className="px-4 py-3 font-medium text-right">Amallar</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {items.map((item) => (
                  <tr key={item.id} className="hover:bg-accent/30 transition-colors">
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        {item.photo_url ? (
                          <img src={item.photo_url} alt="" className="w-9 h-9 rounded-full object-cover" />
                        ) : (
                          <div className="w-9 h-9 rounded-full bg-muted flex items-center justify-center text-xs font-bold">
                            {(item.first_name || item.username)?.[0]?.toUpperCase() || "?"}
                          </div>
                        )}
                        <div>
                          <p className="font-medium text-foreground">{item.first_name} {item.last_name}</p>
                          <p className="text-xs text-muted-foreground">@{item.username}</p>
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <p className="text-xs text-muted-foreground">{item.region}{item.district ? `, ${item.district}` : ""}</p>
                      <p className="text-xs text-muted-foreground">{item.school_name}{item.grade ? ` (${item.grade}-sinf)` : ""}</p>
                    </td>
                    <td className="px-4 py-3">
                      <span className="text-xs text-muted-foreground">{formatDate(item.registered_at || item.created_at)}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`text-[10px] px-2 py-0.5 rounded-full font-medium ${statusBadge(item.status)}`}>
                        {item.status === "pending" ? "Kutilmoqda" : item.status === "approved" ? "Tasdiqlangan" : "Rad etilgan"}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      {item.status === "pending" ? (
                        <div className="flex items-center justify-end gap-1">
                          <button
                            onClick={() => setDetailItem(item)}
                            className="p-1.5 text-blue-400 hover:bg-blue-500/10 rounded-md transition-colors"
                            title="Batafsil"
                          >
                            <Eye className="w-4 h-4" />
                          </button>
                          <button
                            onClick={() => setApproveDialog(item)}
                            className="p-1.5 text-emerald-400 hover:bg-emerald-500/10 rounded-md transition-colors"
                            title="Tasdiqlash"
                          >
                            <UserCheck className="w-4 h-4" />
                          </button>
                          <button
                            onClick={() => setRejectDialog(item)}
                            className="p-1.5 text-red-400 hover:bg-red-500/10 rounded-md transition-colors"
                            title="Rad etish"
                          >
                            <UserX className="w-4 h-4" />
                          </button>
                        </div>
                      ) : (
                        <span className="text-xs text-muted-foreground">
                          {item.reason || item.note || "\u2014"}
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {total > 20 && (
            <div className="flex items-center justify-between px-4 py-3 border-t border-border">
              <span className="text-xs text-muted-foreground">Jami: {total}</span>
              <div className="flex gap-1">
                <button disabled={page <= 1} onClick={() => setPage(page - 1)} className="px-3 py-1 text-xs border border-border rounded disabled:opacity-30">Oldingi</button>
                <button disabled={page * 20 >= total} onClick={() => setPage(page + 1)} className="px-3 py-1 text-xs border border-border rounded disabled:opacity-30">Keyingi</button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Detail Modal */}
      {detailItem && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-xl p-6 max-w-md w-full mx-4 space-y-4">
            <h3 className="text-lg font-bold">Foydalanuvchi ma&apos;lumotlari</h3>
            <div className="flex items-center gap-3">
              {detailItem.photo_url ? (
                <img src={detailItem.photo_url} alt="" className="w-16 h-16 rounded-full object-cover" />
              ) : (
                <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center text-xl font-bold">
                  {(detailItem.first_name || "?")?.[0]?.toUpperCase()}
                </div>
              )}
              <div>
                <p className="font-bold text-lg">{detailItem.first_name} {detailItem.last_name}</p>
                <p className="text-sm text-muted-foreground">@{detailItem.username}</p>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div><span className="text-muted-foreground">Viloyat:</span> <span className="text-foreground">{detailItem.region || "\u2014"}</span></div>
              <div><span className="text-muted-foreground">Tuman:</span> <span className="text-foreground">{detailItem.district || "\u2014"}</span></div>
              <div><span className="text-muted-foreground">Maktab:</span> <span className="text-foreground">{detailItem.school_name || "\u2014"}</span></div>
              <div><span className="text-muted-foreground">Sinf:</span> <span className="text-foreground">{detailItem.grade || "\u2014"}</span></div>
              <div className="col-span-2"><span className="text-muted-foreground">Ro&apos;yxatdan:</span> <span className="text-foreground">{formatDate(detailItem.registered_at || detailItem.created_at)}</span></div>
            </div>
            <div className="flex gap-3">
              <button onClick={() => setDetailItem(null)} className="flex-1 px-4 py-2 bg-muted text-muted-foreground rounded-lg text-sm">Yopish</button>
              <button onClick={() => { setDetailItem(null); setApproveDialog(detailItem); }} className="flex-1 px-4 py-2 bg-emerald-600 text-white rounded-lg text-sm">Tasdiqlash</button>
              <button onClick={() => { setDetailItem(null); setRejectDialog(detailItem); }} className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg text-sm">Rad etish</button>
            </div>
          </div>
        </div>
      )}

      {/* Approve Dialog */}
      {approveDialog && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-xl p-6 max-w-sm w-full mx-4 space-y-4">
            <h3 className="text-lg font-bold text-emerald-400">Tasdiqlash</h3>
            <p className="text-sm text-muted-foreground">
              <strong>{approveDialog.first_name || approveDialog.username}</strong> ni tasdiqlaysizmi? Tasdiqlanganidan keyin barcha funksiyalar ochiladi.
            </p>
            <div>
              <label className="text-sm text-muted-foreground">Izoh (ixtiyoriy)</label>
              <input
                value={approveNote}
                onChange={(e) => setApproveNote(e.target.value)}
                className="w-full mt-1 bg-background border border-border rounded-lg px-3 py-2 text-sm"
                placeholder="Masalan: Telegram bot ishlamagani sabab..."
              />
            </div>
            <div className="flex gap-3">
              <button onClick={() => { setApproveDialog(null); setApproveNote(""); }} className="flex-1 px-4 py-2 bg-muted text-muted-foreground rounded-lg text-sm">Bekor</button>
              <button onClick={handleApprove} disabled={submitting} className="flex-1 px-4 py-2 bg-emerald-600 text-white rounded-lg text-sm disabled:opacity-50">
                {submitting ? "..." : "Tasdiqlash"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Reject Dialog */}
      {rejectDialog && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-xl p-6 max-w-sm w-full mx-4 space-y-4">
            <h3 className="text-lg font-bold text-red-400">Rad etish</h3>
            <p className="text-sm text-muted-foreground">
              <strong>{rejectDialog.first_name || rejectDialog.username}</strong> ni rad etasizmi?
            </p>
            <div>
              <label className="text-sm text-muted-foreground">Sabab (majburiy)</label>
              <textarea
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
                className="w-full mt-1 bg-background border border-border rounded-lg px-3 py-2 text-sm"
                rows={3}
                placeholder="Rad etish sababini yozing..."
              />
            </div>
            <div className="flex gap-3">
              <button onClick={() => { setRejectDialog(null); setRejectReason(""); }} className="flex-1 px-4 py-2 bg-muted text-muted-foreground rounded-lg text-sm">Bekor</button>
              <button onClick={handleReject} disabled={submitting || !rejectReason.trim()} className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg text-sm disabled:opacity-50">
                {submitting ? "..." : "Rad etish"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
