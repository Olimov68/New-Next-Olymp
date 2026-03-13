"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import {
  adminFetchMockExams,
  adminCreateMockExam,
  adminUpdateMockExam,
  adminDeleteMockExam,
  type MockExam,
} from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Plus, Pencil, Trash2, Eye } from "lucide-react";

// ── Label maps ───────────────────────────────────────────────────

const STATUS_LABELS: Record<string, string> = {
  draft: "Qoralama",
  published: "Nashr etilgan",
  archived: "Arxivlangan",
};

const STATUS_COLORS: Record<string, string> = {
  draft: "bg-yellow-50 text-yellow-800 border-yellow-200 dark:bg-yellow-950/30 dark:text-yellow-400 dark:border-yellow-800",
  published: "bg-green-50 text-green-800 border-green-200 dark:bg-green-950/30 dark:text-green-400 dark:border-green-800",
  archived: "bg-muted text-muted-foreground border-border",
};

const LANGUAGE_LABELS: Record<string, string> = {
  uz: "O'zbekcha",
  en: "Inglizcha",
  ru: "Ruscha",
};

const FORMAT_LABELS: Record<string, string> = {
  standard: "Standart",
  sectioned: "Seksiyali",
  adaptive: "Adaptiv",
};

const ASSESSMENT_LABELS: Record<string, string> = {
  normal: "Oddiy",
  rasch: "Rasch",
};

// ── Form ─────────────────────────────────────────────────────────

interface MockExamForm {
  title: string;
  subject: string;
  description: string;
  language: string;
  format: string;
  status: string;
  duration_minutes: number;
  max_attempts: number;
  price: number;
}

const emptyForm: MockExamForm = {
  title: "",
  subject: "",
  description: "",
  language: "uz",
  format: "standard",
  status: "draft",
  duration_minutes: 60,
  max_attempts: 1,
  price: 0,
};

// ── Page ─────────────────────────────────────────────────────────

export default function AdminMockExamsPage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data: response, isLoading } = useQuery({
    queryKey: ["admin-mock-exams"],
    queryFn: () => adminFetchMockExams(),
  });
  const mockExams = response?.data ?? [];

  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<MockExam | null>(null);
  const [form, setForm] = useState<MockExamForm>(emptyForm);

  const createMut = useMutation({
    mutationFn: adminCreateMockExam,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-mock-exams"] });
      setOpen(false);
      resetForm();
    },
  });

  const updateMut = useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<MockExam> }) =>
      adminUpdateMockExam(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-mock-exams"] });
      setOpen(false);
      resetForm();
    },
  });

  const deleteMut = useMutation({
    mutationFn: adminDeleteMockExam,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["admin-mock-exams"] }),
  });

  function resetForm() {
    setForm(emptyForm);
    setEditing(null);
  }

  function openEdit(exam: MockExam) {
    setEditing(exam);
    setForm({
      title: exam.title,
      subject: exam.subject,
      description: exam.description || "",
      language: exam.language || "uz",
      format: exam.format || "standard",
      status: exam.status || "draft",
      duration_minutes: exam.duration_minutes || 60,
      max_attempts: exam.max_attempts || 1,
      price: exam.price || 0,
    });
    setOpen(true);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (editing) {
      updateMut.mutate({ id: editing.id, data: form });
    } else {
      createMut.mutate(form);
    }
  }

  function handleDelete(exam: MockExam) {
    if (confirm(`"${exam.title}" mock imtihonini o'chirishni tasdiqlaysizmi?`)) {
      deleteMut.mutate(exam.id);
    }
  }

  const isPending = createMut.isPending || updateMut.isPending;

  return (
    <div>
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-foreground">Mock imtihonlar</h1>
        <Dialog
          open={open}
          onOpenChange={(v) => {
            setOpen(v);
            if (!v) resetForm();
          }}
        >
          <DialogTrigger className="inline-flex h-9 items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 text-sm font-medium text-white hover:bg-blue-700 transition-colors">
            <Plus className="h-4 w-4" /> Yangi imtihon
          </DialogTrigger>

          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle>
                {editing ? "Mock imtihonni tahrirlash" : "Yangi mock imtihon"}
              </DialogTitle>
            </DialogHeader>
            <form
              onSubmit={handleSubmit}
              className="space-y-4 max-h-[72vh] overflow-y-auto pr-1"
            >
              {/* Nomi */}
              <div className="space-y-2">
                <Label>Nomi</Label>
                <Input
                  value={form.title}
                  onChange={(e) => setForm({ ...form, title: e.target.value })}
                  placeholder="Mock imtihon nomi"
                  required
                />
              </div>

              {/* Fan */}
              <div className="space-y-2">
                <Label>Fan</Label>
                <Input
                  value={form.subject}
                  onChange={(e) => setForm({ ...form, subject: e.target.value })}
                  placeholder="Masalan: Matematika"
                  required
                />
              </div>

              {/* Tavsif */}
              <div className="space-y-2">
                <Label>Tavsif</Label>
                <Textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  rows={3}
                  placeholder="Mock imtihon haqida qisqacha ma'lumot..."
                />
              </div>

              {/* Til + Format */}
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Til</Label>
                  <select
                    className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={form.language}
                    onChange={(e) => setForm({ ...form, language: e.target.value })}
                  >
                    <option value="uz">O&apos;zbekcha</option>
                    <option value="ru">Ruscha</option>
                    <option value="en">Inglizcha</option>
                  </select>
                </div>
                <div className="space-y-2">
                  <Label>Format</Label>
                  <select
                    className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={form.format}
                    onChange={(e) => setForm({ ...form, format: e.target.value })}
                  >
                    <option value="standard">Standart</option>
                    <option value="sectioned">Seksiyali</option>
                    <option value="adaptive">Adaptiv</option>
                  </select>
                </div>
              </div>

              {/* Status */}
              <div className="space-y-2">
                <Label>Status</Label>
                <select
                  className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  value={form.status}
                  onChange={(e) => setForm({ ...form, status: e.target.value })}
                >
                  <option value="draft">Qoralama</option>
                  <option value="published">Nashr etilgan</option>
                  <option value="archived">Arxivlangan</option>
                </select>
              </div>

              {/* Davomiyligi + Urinishlar */}
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Davomiyligi (daqiqa)</Label>
                  <Input
                    type="number"
                    min={1}
                    value={form.duration_minutes}
                    onChange={(e) =>
                      setForm({ ...form, duration_minutes: Number(e.target.value) })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label>Maks. urinishlar</Label>
                  <Input
                    type="number"
                    min={1}
                    value={form.max_attempts}
                    onChange={(e) =>
                      setForm({ ...form, max_attempts: Number(e.target.value) })
                    }
                  />
                </div>
              </div>

              {/* Narx */}
              <div className="space-y-2">
                <Label>Narx (so&apos;m)</Label>
                <Input
                  type="number"
                  min={0}
                  value={form.price}
                  onChange={(e) => setForm({ ...form, price: Number(e.target.value) })}
                />
              </div>

              <Button
                type="submit"
                className="w-full bg-blue-600 hover:bg-blue-700"
                disabled={isPending}
              >
                {isPending ? "Saqlanmoqda..." : editing ? "Saqlash" : "Yaratish"}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Table */}
      <Card className="border border-border shadow-sm">
        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-10 text-center text-muted-foreground">Yuklanmoqda...</div>
          ) : mockExams.length === 0 ? (
            <div className="p-10 text-center text-muted-foreground">
              Hali mock imtihonlar qo&apos;shilmagan.
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-12">ID</TableHead>
                  <TableHead>Nomi</TableHead>
                  <TableHead>Fan</TableHead>
                  <TableHead>Til</TableHead>
                  <TableHead>Format</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Baholash</TableHead>
                  <TableHead className="text-center">Savollar</TableHead>
                  <TableHead className="text-center">Urinishlar</TableHead>
                  <TableHead className="text-right">Amallar</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {mockExams.map((exam: MockExam) => (
                  <TableRow key={exam.id}>
                    <TableCell className="text-muted-foreground text-sm">{exam.id}</TableCell>
                    <TableCell className="font-medium max-w-[200px] truncate">
                      {exam.title}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">{exam.subject}</TableCell>
                    <TableCell className="text-sm">
                      {LANGUAGE_LABELS[exam.language] ?? exam.language}
                    </TableCell>
                    <TableCell className="text-sm">
                      {FORMAT_LABELS[exam.format] ?? exam.format}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="outline"
                        className={`text-xs ${STATUS_COLORS[exam.status] ?? ""}`}
                      >
                        {STATUS_LABELS[exam.status] ?? exam.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-sm">
                      {ASSESSMENT_LABELS[exam.assessment_method] ?? exam.assessment_method}
                    </TableCell>
                    <TableCell className="text-center text-sm">
                      {exam.questions_count ?? 0}
                    </TableCell>
                    <TableCell className="text-center text-sm">
                      {exam.attempts_count ?? 0}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => router.push(`/admin/mock-exams/${exam.id}`)}
                          title="Ko'rish"
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => openEdit(exam)}
                          title="Tahrirlash"
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
                          onClick={() => handleDelete(exam)}
                          title="O'chirish"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
