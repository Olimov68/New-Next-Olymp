"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  adminFetchMockExam,
  adminUpdateMockExam,
  adminFetchMockQuestions,
  adminCreateMockQuestion,
  adminUpdateMockQuestion,
  adminDeleteMockQuestion,
  adminFetchRaschConfig,
  adminUpdateRaschConfig,
  adminUpdateRaschScales,
  type MockExam,
  type MockQuestion,
  type MockOption,
  type RaschScale,
} from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
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
} from "@/components/ui/dialog";
import {
  ArrowLeft,
  Plus,
  Pencil,
  Trash2,
  ChevronDown,
  ChevronUp,
  Save,
} from "lucide-react";

// ── Overview Tab ──────────────────────────────────────────────────

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
const STATUS_LABELS: Record<string, string> = {
  draft: "Qoralama",
  published: "Nashr etilgan",
  archived: "Arxivlangan",
};

function OverviewTab({ exam }: { exam: MockExam }) {
  const items = [
    { label: "Nomi", value: exam.title },
    { label: "Fan", value: exam.subject },
    { label: "Til", value: LANGUAGE_LABELS[exam.language] ?? exam.language },
    { label: "Format", value: FORMAT_LABELS[exam.format] ?? exam.format },
    { label: "Status", value: STATUS_LABELS[exam.status] ?? exam.status },
    {
      label: "Baholash usuli",
      value: exam.assessment_method === "rasch" ? "Rasch modeli" : "Oddiy",
    },
    {
      label: "Davomiyligi",
      value: exam.duration_minutes ? `${exam.duration_minutes} daqiqa` : "Belgilanmagan",
    },
    {
      label: "Maks. urinishlar",
      value: exam.max_attempts ? String(exam.max_attempts) : "Cheksiz",
    },
    {
      label: "Narx",
      value: exam.price === 0 ? "Bepul" : `${exam.price?.toLocaleString()} so'm`,
    },
    { label: "Savollar soni", value: String(exam.questions_count ?? 0) },
    { label: "Urinishlar soni", value: String(exam.attempts_count ?? 0) },
    {
      label: "Yaratilgan",
      value: exam.created_at ? new Date(exam.created_at).toLocaleDateString("uz-UZ") : "-",
    },
  ];

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {items.map((item) => (
          <Card key={item.label} className="border shadow-sm">
            <CardContent className="p-4">
              <p className="text-xs text-muted-foreground mb-1">{item.label}</p>
              <p className="text-sm font-medium">{item.value}</p>
            </CardContent>
          </Card>
        ))}
      </div>
      {exam.description && (
        <Card className="border shadow-sm">
          <CardContent className="p-4">
            <p className="text-xs text-muted-foreground mb-1">Tavsif</p>
            <p className="text-sm">{exam.description}</p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

// ── Questions Tab ─────────────────────────────────────────────────

interface QuestionForm {
  title: string;
  content: string;
  type: string;
  points: number;
  order_num: number;
  options: { id?: number; content: string; is_correct: boolean; order_num: number }[];
}

const defaultQuestionForm: QuestionForm = {
  title: "",
  content: "",
  type: "single_choice",
  points: 1,
  order_num: 1,
  options: [
    { content: "", is_correct: true, order_num: 0 },
    { content: "", is_correct: false, order_num: 1 },
    { content: "", is_correct: false, order_num: 2 },
    { content: "", is_correct: false, order_num: 3 },
  ],
};

const OPTION_LETTERS = "ABCDEFGHIJ";

function QuestionsTab({ examId }: { examId: string }) {
  const queryClient = useQueryClient();

  const { data: questions, isLoading } = useQuery({
    queryKey: ["mock-questions", examId],
    queryFn: () => adminFetchMockQuestions(examId),
  });

  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingQ, setEditingQ] = useState<MockQuestion | null>(null);
  const [form, setForm] = useState<QuestionForm>(defaultQuestionForm);
  const [expandedId, setExpandedId] = useState<number | null>(null);

  const createMut = useMutation({
    mutationFn: (data: QuestionForm) => adminCreateMockQuestion(examId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mock-questions", examId] });
      queryClient.invalidateQueries({ queryKey: ["admin-mock-exam", examId] });
      closeDialog();
    },
  });

  const updateMut = useMutation({
    mutationFn: ({ id, data }: { id: number; data: QuestionForm }) =>
      adminUpdateMockQuestion(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mock-questions", examId] });
      closeDialog();
    },
  });

  const deleteMut = useMutation({
    mutationFn: adminDeleteMockQuestion,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mock-questions", examId] });
      queryClient.invalidateQueries({ queryKey: ["admin-mock-exam", examId] });
    },
  });

  function resetForm() {
    setForm({
      ...defaultQuestionForm,
      order_num: (questions?.length ?? 0) + 1,
    });
    setEditingQ(null);
  }

  function closeDialog() {
    setDialogOpen(false);
    setEditingQ(null);
    setForm(defaultQuestionForm);
  }

  function openCreate() {
    resetForm();
    setDialogOpen(true);
  }

  function openEdit(q: MockQuestion) {
    setEditingQ(q);
    setForm({
      title: q.title,
      content: q.content || "",
      type: q.type || "single_choice",
      points: q.points,
      order_num: q.order_num,
      options:
        q.options && q.options.length > 0
          ? q.options.map((o) => ({
              id: o.id,
              content: o.content,
              is_correct: o.is_correct,
              order_num: o.order_num,
            }))
          : defaultQuestionForm.options,
    });
    setDialogOpen(true);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (editingQ) {
      updateMut.mutate({ id: editingQ.id, data: form });
    } else {
      createMut.mutate(form);
    }
  }

  function setCorrectOption(idx: number) {
    setForm({
      ...form,
      options: form.options.map((o, i) => ({ ...o, is_correct: i === idx })),
    });
  }

  function updateOptionContent(idx: number, content: string) {
    setForm({
      ...form,
      options: form.options.map((o, i) => (i === idx ? { ...o, content } : o)),
    });
  }

  function addOption() {
    setForm({
      ...form,
      options: [
        ...form.options,
        { content: "", is_correct: false, order_num: form.options.length },
      ],
    });
  }

  function removeOption(idx: number) {
    if (form.options.length <= 2) return;
    const next = form.options.filter((_, i) => i !== idx);
    if (!next.some((o) => o.is_correct) && next.length > 0) {
      next[0].is_correct = true;
    }
    setForm({ ...form, options: next.map((o, i) => ({ ...o, order_num: i })) });
  }

  if (isLoading) {
    return <div className="p-8 text-center text-muted-foreground">Yuklanmoqda...</div>;
  }

  const list = questions ?? [];

  return (
    <div>
      {/* Header */}
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold">Savollar ({list.length})</h3>
        <Button className="bg-blue-600 hover:bg-blue-700" onClick={openCreate}>
          <Plus className="h-4 w-4 mr-1" /> Savol qo&apos;shish
        </Button>
      </div>

      {/* Question list */}
      {list.length === 0 ? (
        <Card className="border shadow-sm">
          <CardContent className="p-10 text-center text-muted-foreground">
            Hali savollar qo&apos;shilmagan. &quot;Savol qo&apos;shish&quot; tugmasini bosing.
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {list.map((q: MockQuestion, idx: number) => (
            <Card key={q.id} className="border shadow-sm">
              <CardContent className="p-4">
                <div className="flex items-start justify-between">
                  {/* Title row (clickable to expand) */}
                  <div
                    className="flex-1 cursor-pointer"
                    onClick={() =>
                      setExpandedId(expandedId === q.id ? null : q.id)
                    }
                  >
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-sm font-bold text-blue-600">
                        #{idx + 1}
                      </span>
                      <span className="font-medium text-sm">{q.title}</span>
                      <Badge variant="outline" className="text-xs capitalize">
                        {q.type?.replace(/_/g, " ")}
                      </Badge>
                      <Badge variant="outline" className="text-xs">
                        {q.points} ball
                      </Badge>
                    </div>
                  </div>
                  {/* Actions */}
                  <div className="flex items-center gap-1 ml-3">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() =>
                        setExpandedId(expandedId === q.id ? null : q.id)
                      }
                    >
                      {expandedId === q.id ? (
                        <ChevronUp className="h-4 w-4" />
                      ) : (
                        <ChevronDown className="h-4 w-4" />
                      )}
                    </Button>
                    <Button variant="ghost" size="icon" onClick={() => openEdit(q)}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
                      onClick={() => {
                        if (confirm(`"${q.title}" savolini o'chirishni tasdiqlaysizmi?`)) {
                          deleteMut.mutate(q.id);
                        }
                      }}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>

                {/* Expanded content */}
                {expandedId === q.id && (
                  <div className="mt-3 pt-3 border-t border-border">
                    {q.content && (
                      <p className="text-sm text-muted-foreground mb-3">{q.content}</p>
                    )}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                      {q.options?.map((opt: MockOption, oi: number) => (
                        <div
                          key={opt.id ?? oi}
                          className={`flex items-center gap-2 p-2 rounded-md text-sm border ${
                            opt.is_correct
                              ? "bg-green-50 border-green-200 text-green-800 dark:bg-green-950/30 dark:border-green-800 dark:text-green-400"
                              : "bg-muted border-border text-foreground"
                          }`}
                        >
                          <span
                            className={`font-bold text-xs ${
                              opt.is_correct ? "text-green-600 dark:text-green-400" : "text-muted-foreground"
                            }`}
                          >
                            {OPTION_LETTERS[oi] ?? oi}
                          </span>
                          <span>{opt.content}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Create/Edit Dialog */}
      <Dialog
        open={dialogOpen}
        onOpenChange={(v) => {
          if (!v) closeDialog();
          else setDialogOpen(true);
        }}
      >
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              {editingQ ? "Savolni tahrirlash" : "Yangi savol"}
            </DialogTitle>
          </DialogHeader>
          <form
            onSubmit={handleSubmit}
            className="space-y-4 max-h-[72vh] overflow-y-auto pr-1"
          >
            {/* Savol nomi */}
            <div className="space-y-2">
              <Label>Savol nomi (sarlavha)</Label>
              <Input
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                placeholder="Savol sarlavhasi"
                required
              />
            </div>

            {/* Savol matni */}
            <div className="space-y-2">
              <Label>Savol matni (ixtiyoriy)</Label>
              <Textarea
                value={form.content}
                onChange={(e) => setForm({ ...form, content: e.target.value })}
                rows={3}
                placeholder="Batafsil savol matni..."
              />
            </div>

            {/* Turi + Ball + Tartib */}
            <div className="grid grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label>Savol turi</Label>
                <select
                  className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  value={form.type}
                  onChange={(e) => setForm({ ...form, type: e.target.value })}
                >
                  <option value="single_choice">Bitta javob</option>
                  <option value="multiple_choice">Ko&apos;p javob</option>
                  <option value="true_false">To&apos;g&apos;ri / Noto&apos;g&apos;ri</option>
                  <option value="open_ended">Ochiq javob</option>
                </select>
              </div>
              <div className="space-y-2">
                <Label>Ball</Label>
                <Input
                  type="number"
                  min={1}
                  value={form.points}
                  onChange={(e) =>
                    setForm({ ...form, points: Number(e.target.value) })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label>Tartib raqami</Label>
                <Input
                  type="number"
                  min={1}
                  value={form.order_num}
                  onChange={(e) =>
                    setForm({ ...form, order_num: Number(e.target.value) })
                  }
                />
              </div>
            </div>

            {/* Options */}
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <Label>Javob variantlari (to&apos;g&apos;risini bosing)</Label>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={addOption}
                >
                  <Plus className="h-3 w-3 mr-1" /> Variant
                </Button>
              </div>
              <div className="space-y-2">
                {form.options.map((opt, idx) => (
                  <div key={idx} className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => setCorrectOption(idx)}
                      className={`flex-shrink-0 w-8 h-8 rounded-full border-2 flex items-center justify-center text-sm font-bold transition-colors ${
                        opt.is_correct
                          ? "bg-green-500 border-green-500 text-white"
                          : "border-border text-muted-foreground hover:border-green-400"
                      }`}
                    >
                      {OPTION_LETTERS[idx]}
                    </button>
                    <Input
                      value={opt.content}
                      onChange={(e) => updateOptionContent(idx, e.target.value)}
                      placeholder={`${OPTION_LETTERS[idx]} varianti`}
                      className="flex-1"
                    />
                    {form.options.length > 2 && (
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="text-red-500 flex-shrink-0"
                        onClick={() => removeOption(idx)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    )}
                  </div>
                ))}
              </div>
            </div>

            <Button
              type="submit"
              className="w-full bg-blue-600 hover:bg-blue-700"
              disabled={createMut.isPending || updateMut.isPending}
            >
              {createMut.isPending || updateMut.isPending
                ? "Saqlanmoqda..."
                : editingQ
                ? "Saqlash"
                : "Qo'shish"}
            </Button>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}

// ── Rasch Settings Tab ────────────────────────────────────────────

function RaschSettingsTab({
  examId,
  currentMethod,
}: {
  examId: string;
  currentMethod: string;
}) {
  const queryClient = useQueryClient();

  const { data: raschData, isLoading } = useQuery({
    queryKey: ["rasch-config", examId],
    queryFn: () => adminFetchRaschConfig(examId),
  });

  // Config form state
  const [configForm, setConfigForm] = useState({
    min_theta: -4,
    max_theta: 4,
    scaling_method: "linear",
    pass_threshold: 0,
  });
  const [configSynced, setConfigSynced] = useState(false);

  // Scales state
  const [scales, setScales] = useState<RaschScale[]>([]);
  const [scalesSynced, setScalesSynced] = useState(false);

  // Sync from fetched data (run once after load)
  if (raschData && !configSynced) {
    if (raschData.config) {
      setConfigForm({
        min_theta: raschData.config.min_theta,
        max_theta: raschData.config.max_theta,
        scaling_method: raschData.config.scaling_method || "linear",
        pass_threshold: raschData.config.pass_threshold,
      });
    }
    setConfigSynced(true);
  }
  if (raschData && !scalesSynced) {
    setScales(raschData.scales ?? []);
    setScalesSynced(true);
  }

  // Mutations
  const updateMethodMut = useMutation({
    mutationFn: (method: string) =>
      adminUpdateMockExam(examId, { assessment_method: method }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-mock-exam", examId] });
      queryClient.invalidateQueries({ queryKey: ["rasch-config", examId] });
    },
  });

  const updateConfigMut = useMutation({
    mutationFn: () => adminUpdateRaschConfig(examId, configForm),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["rasch-config", examId] }),
  });

  const updateScalesMut = useMutation({
    mutationFn: () => adminUpdateRaschScales(examId, { scales }),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["rasch-config", examId] }),
  });

  // Scales helpers
  function addScaleRow() {
    setScales([
      ...scales,
      { raw_score: scales.length, theta: 0, scaled_score: 0, level: "" },
    ]);
  }

  function removeScaleRow(idx: number) {
    setScales(scales.filter((_, i) => i !== idx));
  }

  function updateScaleField(
    idx: number,
    field: keyof RaschScale,
    value: string
  ) {
    setScales(
      scales.map((s, i) =>
        i === idx
          ? { ...s, [field]: field === "level" ? value : Number(value) }
          : s
      )
    );
  }

  if (isLoading) {
    return <div className="p-8 text-center text-muted-foreground">Yuklanmoqda...</div>;
  }

  const assessmentMethod = raschData?.assessment_method ?? currentMethod ?? "normal";
  const isRasch = assessmentMethod === "rasch";

  return (
    <div className="space-y-6">
      {/* Assessment method toggle */}
      <Card className="border shadow-sm">
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Baholash usuli</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4">
            <button
              type="button"
              onClick={() => updateMethodMut.mutate("normal")}
              className={`px-4 py-2 rounded-lg border text-sm font-medium transition-colors ${
                !isRasch
                  ? "bg-blue-600 border-blue-600 text-white"
                  : "bg-card border-border text-foreground hover:border-blue-400"
              }`}
              disabled={updateMethodMut.isPending}
            >
              Oddiy baholash
            </button>
            <button
              type="button"
              onClick={() => updateMethodMut.mutate("rasch")}
              className={`px-4 py-2 rounded-lg border text-sm font-medium transition-colors ${
                isRasch
                  ? "bg-purple-600 border-purple-600 text-white"
                  : "bg-card border-border text-foreground hover:border-purple-400"
              }`}
              disabled={updateMethodMut.isPending}
            >
              Rasch modeli
            </button>
            {updateMethodMut.isPending && (
              <span className="text-sm text-muted-foreground">Saqlanmoqda...</span>
            )}
          </div>
          <p className="text-sm text-muted-foreground mt-3">
            {isRasch
              ? "Rasch modeli yoqlgan. Quyida Rasch konfiguratsiyasini va shkalasini sozlashingiz mumkin."
              : "Oddiy baholash usuli yoqlgan. Rasch modelini yoqish uchun yuqoridagi tugmani bosing."}
          </p>
        </CardContent>
      </Card>

      {/* Rasch config -- only when rasch is active */}
      {isRasch && (
        <>
          {/* Config form */}
          <Card className="border shadow-sm">
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Rasch konfiguratsiyasi</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Min Theta</Label>
                  <Input
                    type="number"
                    step="0.1"
                    value={configForm.min_theta}
                    onChange={(e) =>
                      setConfigForm({
                        ...configForm,
                        min_theta: Number(e.target.value),
                      })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label>Max Theta</Label>
                  <Input
                    type="number"
                    step="0.1"
                    value={configForm.max_theta}
                    onChange={(e) =>
                      setConfigForm({
                        ...configForm,
                        max_theta: Number(e.target.value),
                      })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label>Scaling usuli</Label>
                  <select
                    className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={configForm.scaling_method}
                    onChange={(e) =>
                      setConfigForm({
                        ...configForm,
                        scaling_method: e.target.value,
                      })
                    }
                  >
                    <option value="linear">Linear</option>
                    <option value="logistic">Logistic</option>
                  </select>
                </div>
                <div className="space-y-2">
                  <Label>O&apos;tish chegarasi (Pass Threshold)</Label>
                  <Input
                    type="number"
                    step="0.1"
                    value={configForm.pass_threshold}
                    onChange={(e) =>
                      setConfigForm({
                        ...configForm,
                        pass_threshold: Number(e.target.value),
                      })
                    }
                  />
                </div>
              </div>
              <Button
                className="mt-4 bg-blue-600 hover:bg-blue-700"
                onClick={() => updateConfigMut.mutate()}
                disabled={updateConfigMut.isPending}
              >
                <Save className="h-4 w-4 mr-1" />
                {updateConfigMut.isPending
                  ? "Saqlanmoqda..."
                  : "Konfiguratsiyani saqlash"}
              </Button>
              {updateConfigMut.isSuccess && (
                <p className="text-sm text-green-600 mt-2">
                  Muvaffaqiyatli saqlandi!
                </p>
              )}
            </CardContent>
          </Card>

          {/* Scales table */}
          <Card className="border shadow-sm">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-base">
                  Rasch shkalasi ({scales.length} qator)
                </CardTitle>
                <Button variant="outline" size="sm" onClick={addScaleRow}>
                  <Plus className="h-3 w-3 mr-1" /> Qator qo&apos;shish
                </Button>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Raw Score</TableHead>
                    <TableHead>Theta</TableHead>
                    <TableHead>Scaled Score</TableHead>
                    <TableHead>Level</TableHead>
                    <TableHead className="w-12"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {scales.map((s, idx) => (
                    <TableRow key={idx}>
                      <TableCell>
                        <Input
                          type="number"
                          value={s.raw_score}
                          onChange={(e) =>
                            updateScaleField(idx, "raw_score", e.target.value)
                          }
                          className="h-8 w-20"
                        />
                      </TableCell>
                      <TableCell>
                        <Input
                          type="number"
                          step="0.01"
                          value={s.theta}
                          onChange={(e) =>
                            updateScaleField(idx, "theta", e.target.value)
                          }
                          className="h-8 w-24"
                        />
                      </TableCell>
                      <TableCell>
                        <Input
                          type="number"
                          value={s.scaled_score}
                          onChange={(e) =>
                            updateScaleField(
                              idx,
                              "scaled_score",
                              e.target.value
                            )
                          }
                          className="h-8 w-24"
                        />
                      </TableCell>
                      <TableCell>
                        <Input
                          value={s.level}
                          onChange={(e) =>
                            updateScaleField(idx, "level", e.target.value)
                          }
                          className="h-8 w-32"
                          placeholder="beginner, intermediate..."
                        />
                      </TableCell>
                      <TableCell>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="text-red-500 h-8 w-8"
                          onClick={() => removeScaleRow(idx)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                  {scales.length === 0 && (
                    <TableRow>
                      <TableCell
                        colSpan={5}
                        className="text-center text-muted-foreground py-8"
                      >
                        Hali shkalalar qo&apos;shilmagan
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
              {scales.length > 0 && (
                <div className="p-4 border-t border-border">
                  <Button
                    className="bg-blue-600 hover:bg-blue-700"
                    onClick={() => updateScalesMut.mutate()}
                    disabled={updateScalesMut.isPending}
                  >
                    <Save className="h-4 w-4 mr-1" />
                    {updateScalesMut.isPending
                      ? "Saqlanmoqda..."
                      : "Barcha shkalalarni saqlash"}
                  </Button>
                  {updateScalesMut.isSuccess && (
                    <p className="text-sm text-green-600 mt-2">
                      Shkalalar muvaffaqiyatli saqlandi!
                    </p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}

// ── Settings Tab ──────────────────────────────────────────────────

function SettingsTab({ exam, examId }: { exam: MockExam; examId: string }) {
  const queryClient = useQueryClient();
  const [form, setForm] = useState({
    title: exam.title,
    subject: exam.subject,
    description: exam.description || "",
    language: exam.language || "uz",
    format: exam.format || "standard",
    status: exam.status || "draft",
    assessment_method: exam.assessment_method || "normal",
    duration_minutes: exam.duration_minutes || 60,
    max_attempts: exam.max_attempts || 1,
    price: exam.price || 0,
  });

  const updateMut = useMutation({
    mutationFn: (data: Partial<MockExam>) => adminUpdateMockExam(examId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-mock-exam", examId] });
    },
  });

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    updateMut.mutate(form);
  }

  return (
    <Card className="border shadow-sm max-w-2xl">
      <CardHeader className="pb-3">
        <CardTitle className="text-base">Mock imtihon sozlamalari</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Nomi</Label>
            <Input
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              required
            />
          </div>
          <div className="space-y-2">
            <Label>Fan</Label>
            <Input
              value={form.subject}
              onChange={(e) => setForm({ ...form, subject: e.target.value })}
              required
            />
          </div>
          <div className="space-y-2">
            <Label>Tavsif</Label>
            <Textarea
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              rows={3}
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Til</Label>
              <select
                className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm"
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
                className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm"
                value={form.format}
                onChange={(e) => setForm({ ...form, format: e.target.value })}
              >
                <option value="standard">Standart</option>
                <option value="sectioned">Seksiyali</option>
                <option value="adaptive">Adaptiv</option>
              </select>
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>Status</Label>
              <select
                className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm"
                value={form.status}
                onChange={(e) => setForm({ ...form, status: e.target.value })}
              >
                <option value="draft">Qoralama</option>
                <option value="published">Nashr etilgan</option>
                <option value="archived">Arxivlangan</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label>Baholash usuli</Label>
              <select
                className="w-full rounded-md border border-border bg-background text-foreground px-3 py-2 text-sm"
                value={form.assessment_method}
                onChange={(e) =>
                  setForm({ ...form, assessment_method: e.target.value })
                }
              >
                <option value="normal">Oddiy</option>
                <option value="rasch">Rasch modeli</option>
              </select>
            </div>
          </div>
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
            className="bg-blue-600 hover:bg-blue-700"
            disabled={updateMut.isPending}
          >
            <Save className="h-4 w-4 mr-1" />
            {updateMut.isPending ? "Saqlanmoqda..." : "Saqlash"}
          </Button>
          {updateMut.isSuccess && (
            <p className="text-sm text-green-600">Muvaffaqiyatli saqlandi!</p>
          )}
        </form>
      </CardContent>
    </Card>
  );
}

// ── Placeholder Tab ───────────────────────────────────────────────

function PlaceholderTab({ message }: { message: string }) {
  return (
    <Card className="border shadow-sm">
      <CardContent className="p-12 text-center">
        <p className="text-muted-foreground text-base">{message}</p>
      </CardContent>
    </Card>
  );
}

// ── Main Detail Page ──────────────────────────────────────────────

export default function MockExamDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const { data: exam, isLoading } = useQuery({
    queryKey: ["admin-mock-exam", id],
    queryFn: () => adminFetchMockExam(id),
    enabled: !!id,
  });

  if (isLoading) {
    return <div className="p-8 text-center text-muted-foreground">Yuklanmoqda...</div>;
  }

  if (!exam) {
    return (
      <div className="p-8 text-center text-muted-foreground">
        Mock imtihon topilmadi
      </div>
    );
  }

  return (
    <div>
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => router.push("/admin/mock-exams")}
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-foreground">{exam.title}</h1>
          <div className="flex items-center gap-2 mt-1">
            <Badge
              variant="outline"
              className={`text-xs ${
                exam.status === "published"
                  ? "bg-green-50 text-green-800 border-green-200 dark:bg-green-950/30 dark:text-green-400 dark:border-green-800"
                  : exam.status === "archived"
                  ? "bg-muted text-muted-foreground"
                  : "bg-yellow-50 text-yellow-800 border-yellow-200 dark:bg-yellow-950/30 dark:text-yellow-400 dark:border-yellow-800"
              }`}
            >
              {STATUS_LABELS[exam.status] ?? exam.status}
            </Badge>
            <Badge
              variant="outline"
              className={`text-xs ${
                exam.assessment_method === "rasch"
                  ? "bg-purple-50 text-purple-700 border-purple-200 dark:bg-purple-950/30 dark:text-purple-400 dark:border-purple-800"
                  : "bg-muted"
              }`}
            >
              {exam.assessment_method === "rasch" ? "Rasch" : "Oddiy baholash"}
            </Badge>
            <span className="text-sm text-muted-foreground">{exam.subject}</span>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="overview">
        <TabsList className="mb-4 flex-wrap h-auto gap-1">
          <TabsTrigger value="overview">Umumiy</TabsTrigger>
          <TabsTrigger value="questions">Savollar</TabsTrigger>
          <TabsTrigger value="test-structure">Test tuzilishi</TabsTrigger>
          <TabsTrigger value="rasch">Rasch sozlamalari</TabsTrigger>
          <TabsTrigger value="participants">Ishtirokchilar</TabsTrigger>
          <TabsTrigger value="results">Natijalar</TabsTrigger>
          <TabsTrigger value="certificate">Sertifikat</TabsTrigger>
          <TabsTrigger value="settings">Sozlamalar</TabsTrigger>
        </TabsList>

        <TabsContent value="overview">
          <OverviewTab exam={exam} />
        </TabsContent>

        <TabsContent value="questions">
          <QuestionsTab examId={id} />
        </TabsContent>

        <TabsContent value="test-structure">
          <PlaceholderTab message="Test tuzilmasi tez orada..." />
        </TabsContent>

        <TabsContent value="rasch">
          <RaschSettingsTab
            examId={id}
            currentMethod={exam.assessment_method}
          />
        </TabsContent>

        <TabsContent value="participants">
          <PlaceholderTab message="Ishtirokchilar tez orada..." />
        </TabsContent>

        <TabsContent value="results">
          <PlaceholderTab message="Natijalar tez orada..." />
        </TabsContent>

        <TabsContent value="certificate">
          <PlaceholderTab message="Sertifikat sozlamalari tez orada..." />
        </TabsContent>

        <TabsContent value="settings">
          <SettingsTab exam={exam} examId={id} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
