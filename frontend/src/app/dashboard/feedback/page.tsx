"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter,
} from "@/components/ui/dialog";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from "@/components/ui/select";
import { Loader2, MessageSquare, Plus, ChevronRight, CheckCircle2 } from "lucide-react";
import { toast } from "sonner";
import { createFeedback, listFeedbacks, getFeedback } from "@/lib/user-api";

interface Feedback {
  id: number;
  category: string;
  subject: string;
  message: string;
  status: "open" | "in_review" | "answered" | "closed";
  admin_reply?: string;
  replied_at?: string;
  created_at: string;
}

const statusColors: Record<string, string> = {
  open: "bg-blue-500/10 text-blue-600 border-blue-500/20",
  in_review: "bg-yellow-500/10 text-yellow-600 border-yellow-500/20",
  answered: "bg-green-500/10 text-green-600 border-green-500/20",
  closed: "bg-gray-500/10 text-gray-500 border-gray-500/20",
};

const statusLabels: Record<string, string> = {
  open: "Ochiq",
  in_review: "Ko'rib chiqilmoqda",
  answered: "Javob berildi",
  closed: "Yopildi",
};

const categories = [
  { value: "general", label: "Umumiy" },
  { value: "technical", label: "Texnik muammo" },
  { value: "payment", label: "To'lov" },
  { value: "olympiad", label: "Olimpiada" },
  { value: "account", label: "Hisob" },
  { value: "other", label: "Boshqa" },
];

export default function DashboardFeedbackPage() {
  const [feedbacks, setFeedbacks] = useState<Feedback[]>([]);
  const [loading, setLoading] = useState(true);
  const [createOpen, setCreateOpen] = useState(false);
  const [detailItem, setDetailItem] = useState<Feedback | null>(null);
  const [form, setForm] = useState({ category: "general", subject: "", message: "" });
  const [saving, setSaving] = useState(false);

  const load = async () => {
    setLoading(true);
    try {
      const data = await listFeedbacks();
      setFeedbacks(Array.isArray(data) ? data : (data as any)?.data || []);
    } catch {
      setFeedbacks([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    if (!form.subject.trim() || !form.message.trim()) {
      toast.error("Mavzu va xabar majburiy");
      return;
    }
    if (form.message.length < 10) {
      toast.error("Xabar kamida 10 ta belgidan iborat bo'lishi kerak");
      return;
    }
    setSaving(true);
    try {
      await createFeedback(form.category, form.subject, form.message);
      toast.success("Murojaat yuborildi");
      setCreateOpen(false);
      setForm({ category: "general", subject: "", message: "" });
      load();
    } catch (e: any) {
      toast.error(e?.response?.data?.error || "Xatolik yuz berdi");
    } finally {
      setSaving(false);
    }
  };

  const openDetail = async (item: Feedback) => {
    try {
      const detail = await getFeedback(item.id);
      setDetailItem(detail);
    } catch {
      setDetailItem(item);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Murojaatlar</h1>
          <p className="text-sm text-muted-foreground mt-1">Savollar, shikoyatlar va takliflar</p>
        </div>
        <Button onClick={() => setCreateOpen(true)} className="gap-2">
          <Plus className="h-4 w-4" /> Murojaat yuborish
        </Button>
      </div>

      {/* List */}
      {loading ? (
        <div className="flex items-center justify-center py-20">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      ) : feedbacks.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-muted-foreground border border-dashed border-border rounded-2xl">
          <MessageSquare className="h-12 w-12 mb-3 opacity-20" />
          <p className="font-medium">Murojaatlar yo'q</p>
          <p className="text-sm mt-1">Birinchi murojaatingizni yuboring</p>
          <Button className="mt-4 gap-2" onClick={() => setCreateOpen(true)}>
            <Plus className="h-4 w-4" /> Murojaat yuborish
          </Button>
        </div>
      ) : (
        <div className="space-y-3">
          {feedbacks.map(item => (
            <div
              key={item.id}
              className="rounded-xl border border-border bg-card p-4 flex items-start gap-4 cursor-pointer hover:shadow-md transition-all"
              onClick={() => openDetail(item)}
            >
              <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl bg-primary/10">
                <MessageSquare className="h-5 w-5 text-primary" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1 flex-wrap">
                  <span className="font-semibold text-foreground line-clamp-1">{item.subject}</span>
                  <Badge variant="outline" className={`text-xs flex-shrink-0 ${statusColors[item.status]}`}>
                    {statusLabels[item.status]}
                  </Badge>
                </div>
                <p className="text-sm text-muted-foreground line-clamp-2">{item.message}</p>
                <div className="flex items-center gap-3 mt-2 text-xs text-muted-foreground">
                  <span>{categories.find(c => c.value === item.category)?.label || item.category}</span>
                  <span>•</span>
                  <span>{new Date(item.created_at).toLocaleDateString("uz-UZ")}</span>
                  {item.admin_reply && (
                    <>
                      <span>•</span>
                      <span className="text-green-600 flex items-center gap-1"><CheckCircle2 className="h-3 w-3" />Javob bor</span>
                    </>
                  )}
                </div>
              </div>
              <ChevronRight className="h-4 w-4 text-muted-foreground flex-shrink-0 mt-1" />
            </div>
          ))}
        </div>
      )}

      {/* Create Dialog */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Murojaat yuborish</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="space-y-1.5">
              <Label>Kategoriya</Label>
              <Select value={form.category} onValueChange={v => setForm(f => ({ ...f, category: v ?? "general" }))}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  {categories.map(c => <SelectItem key={c.value} value={c.value}>{c.label}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1.5">
              <Label>Mavzu *</Label>
              <Input placeholder="Murojaat mavzusi" value={form.subject} onChange={e => setForm(f => ({ ...f, subject: e.target.value }))} />
            </div>
            <div className="space-y-1.5">
              <Label>Xabar *</Label>
              <Textarea
                placeholder="Murojaat matnini yozing..."
                rows={5}
                value={form.message}
                onChange={e => setForm(f => ({ ...f, message: e.target.value }))}
              />
              <p className="text-xs text-muted-foreground text-right">{form.message.length} / (min 10)</p>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateOpen(false)}>Bekor qilish</Button>
            <Button onClick={handleCreate} disabled={saving}>
              {saving && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
              Yuborish
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Detail Dialog */}
      <Dialog open={detailItem !== null} onOpenChange={() => setDetailItem(null)}>
        <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
          {detailItem && (
            <>
              <DialogHeader>
                <DialogTitle>{detailItem.subject}</DialogTitle>
              </DialogHeader>
              <div className="space-y-4 py-2">
                <div className="flex items-center gap-2">
                  <Badge variant="outline" className={`text-xs ${statusColors[detailItem.status]}`}>
                    {statusLabels[detailItem.status]}
                  </Badge>
                  <span className="text-xs text-muted-foreground">
                    {categories.find(c => c.value === detailItem.category)?.label || detailItem.category}
                  </span>
                  <span className="text-xs text-muted-foreground">• {new Date(detailItem.created_at).toLocaleDateString("uz-UZ")}</span>
                </div>

                <div className="rounded-xl border border-border bg-muted/40 p-4">
                  <p className="text-sm font-medium text-muted-foreground mb-2">Sizning xabaringiz:</p>
                  <p className="text-sm text-foreground whitespace-pre-wrap">{detailItem.message}</p>
                </div>

                {detailItem.admin_reply ? (
                  <div className="rounded-xl border border-green-500/20 bg-green-500/5 p-4">
                    <p className="text-sm font-medium text-green-600 mb-2 flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4" /> Admin javobi:
                    </p>
                    <p className="text-sm text-foreground whitespace-pre-wrap">{detailItem.admin_reply}</p>
                    {detailItem.replied_at && (
                      <p className="text-xs text-muted-foreground mt-2">{new Date(detailItem.replied_at).toLocaleDateString("uz-UZ")}</p>
                    )}
                  </div>
                ) : (
                  <div className="rounded-xl border border-dashed border-border p-4 text-center text-sm text-muted-foreground">
                    Hali javob berilmagan. Tez orada javob beriladi.
                  </div>
                )}
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setDetailItem(null)}>Yopish</Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
