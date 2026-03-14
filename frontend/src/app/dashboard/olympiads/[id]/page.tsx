"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Trophy, Clock, BookOpen, Users, CheckCircle2, XCircle, AlertCircle,
  ArrowLeft, ArrowRight, Flag, DollarSign, Loader2, Award,
} from "lucide-react";
import {
  getOlympiad,
  joinOlympiad,
  startOlympiad,
  submitOlympiadAnswer,
  finishOlympiad,
  getOlympiadAttemptResult,
  getBalance,
} from "@/lib/user-api";
import { toast } from "sonner";

interface Option {
  id: number;
  label: string;
  text: string;
  is_correct?: boolean;
}

interface Question {
  id: number;
  text: string;
  points: number;
  order_num: number;
  options: Option[];
}

interface AttemptResult {
  total_questions: number;
  correct_answers: number;
  score: number;
  max_score: number;
  percentage: number;
  passed: boolean;
  time_taken_seconds: number;
  answers?: { question_id: number; is_correct: boolean; correct_option_id: number; selected_option_id: number }[];
}

type Phase = "detail" | "payment" | "ready" | "exam" | "result";

function formatTime(sec: number) {
  const m = Math.floor(sec / 60);
  const s = sec % 60;
  return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
}

export default function OlympiadDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();

  const [olympiad, setOlympiad] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [phase, setPhase] = useState<Phase>("detail");

  // Attempt state
  const [attemptId, setAttemptId] = useState<number | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentIdx, setCurrentIdx] = useState(0);
  const [answers, setAnswers] = useState<Record<number, number>>({}); // questionId → optionId
  const [timeLeft, setTimeLeft] = useState(0);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const [finishing, setFinishing] = useState(false);
  const [finishConfirm, setFinishConfirm] = useState(false);

  // Result state
  const [result, setResult] = useState<AttemptResult | null>(null);

  // Payment state
  const [balance, setBalance] = useState<number>(0);

  // Stable ref to handleFinish so the timer closure doesn't capture stale state
  const handleFinishRef = useRef<(auto?: boolean) => void>(() => {});

  const handleFinish = useCallback(async (auto = false) => {
    if (!attemptId) return;
    setFinishConfirm(false);
    setFinishing(true);
    if (timerRef.current) clearInterval(timerRef.current);
    try {
      await finishOlympiad(attemptId);
      const resultData = await getOlympiadAttemptResult(attemptId);
      setResult(resultData as unknown as AttemptResult);
      setPhase("result");
    } catch (e: any) {
      toast.error(e?.response?.data?.error || "Yakunlashda xatolik yuz berdi");
    } finally {
      setFinishing(false);
    }
  }, [attemptId]);

  useEffect(() => {
    handleFinishRef.current = handleFinish;
  }, [handleFinish]);

  useEffect(() => {
    Promise.all([
      getOlympiad(Number(id)),
      getBalance().catch(() => ({ balance: 0 })),
    ]).then(([olympiadData, balanceData]) => {
      setOlympiad(olympiadData);
      setBalance((balanceData as any)?.balance || 0);
    }).catch(() => {})
      .finally(() => setLoading(false));
    return () => { if (timerRef.current) clearInterval(timerRef.current); };
  }, [id]);

  const startTimer = useCallback((durationMinutes: number) => {
    setTimeLeft(durationMinutes * 60);
    timerRef.current = setInterval(() => {
      setTimeLeft(prev => {
        if (prev <= 1) {
          clearInterval(timerRef.current!);
          handleFinishRef.current(true);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  }, []);

  const handleJoin = async () => {
    if (!olympiad) return;
    if (olympiad.is_paid) {
      setPhase("payment");
      return;
    }
    try {
      await joinOlympiad(Number(id));
      setPhase("ready");
    } catch (e: any) {
      const msg = e?.response?.data?.error || "";
      if (msg.includes("already")) {
        setPhase("ready");
      } else {
        console.error(msg || "Xatolik yuz berdi");
      }
    }
  };

  const handlePayAndJoin = async () => {
    if (!olympiad) return;
    if (balance < (olympiad.price || 0)) {
      return;
    }
    try {
      await joinOlympiad(Number(id));
      setPhase("ready");
    } catch (e: any) {
      const msg = e?.response?.data?.error || "";
      if (msg.includes("already")) setPhase("ready");
      else console.error(msg || "Xatolik yuz berdi");
    }
  };

  const handleStart = async () => {
    try {
      const attempt = await startOlympiad(Number(id));
      const aId = (attempt as any).id || (attempt as any).attempt_id;
      setAttemptId(aId);
      const qs: Question[] = (attempt as any).questions || [];
      setQuestions(qs);
      setCurrentIdx(0);
      setAnswers({});
      setPhase("exam");
      startTimer((attempt as any).duration_minutes || olympiad?.duration_minutes || 60);
    } catch (e: any) {
      console.error(e?.response?.data?.error || "Test boshlashda xatolik");
    }
  };

  const handleAnswer = async (optionId: number) => {
    const question = questions[currentIdx];
    if (!question || !attemptId) return;
    setAnswers(prev => ({ ...prev, [question.id]: optionId }));
    try {
      await submitOlympiadAnswer(attemptId, question.id, optionId);
    } catch {
      // Ignore submit errors - answer saved locally
    }
    // Auto-advance if not last
    if (currentIdx < questions.length - 1) {
      setTimeout(() => setCurrentIdx(i => i + 1), 300);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-32">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!olympiad) {
    return (
      <div className="flex flex-col items-center py-32 text-muted-foreground">
        <Trophy className="h-14 w-14 mb-4 opacity-20" />
        <p>Olimpiada topilmadi</p>
        <Button variant="outline" className="mt-4" onClick={() => router.push("/dashboard/olympiads")}>Orqaga</Button>
      </div>
    );
  }

  // === RESULT PHASE ===
  if (phase === "result" && result) {
    const pct = result.percentage || Math.round((result.correct_answers / result.total_questions) * 100);
    return (
      <div className="max-w-lg mx-auto py-10 space-y-6">
        <div className="rounded-2xl border border-border bg-card p-8 text-center space-y-4">
          <div className={`inline-flex h-20 w-20 items-center justify-center rounded-full ${result.passed !== false ? "bg-green-500/10" : "bg-red-500/10"}`}>
            {result.passed !== false
              ? <Award className="h-10 w-10 text-green-600" />
              : <XCircle className="h-10 w-10 text-red-500" />}
          </div>
          <div>
            <h2 className="text-2xl font-bold text-foreground">{result.passed !== false ? "Tabriklaymiz!" : "Keyingi safar yaxshiroq!"}</h2>
            <p className="text-muted-foreground mt-1">{olympiad.title}</p>
          </div>
          <div className="grid grid-cols-3 gap-4 py-4 border-y border-border">
            <div className="text-center">
              <p className="text-2xl font-bold text-foreground">{result.correct_answers}</p>
              <p className="text-xs text-muted-foreground">To&apos;g&apos;ri</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold text-foreground">{result.total_questions - result.correct_answers}</p>
              <p className="text-xs text-muted-foreground">Noto&apos;g&apos;ri</p>
            </div>
            <div className="text-center">
              <p className="text-2xl font-bold text-primary">{pct}%</p>
              <p className="text-xs text-muted-foreground">Natija</p>
            </div>
          </div>
          <div>
            <div className="flex justify-between text-sm mb-2">
              <span className="text-muted-foreground">Ball</span>
              <span className="font-semibold text-foreground">{result.score} / {result.max_score}</span>
            </div>
            <Progress value={pct} className="h-3" />
          </div>
          {result.time_taken_seconds > 0 && (
            <p className="text-sm text-muted-foreground">
              Sarflangan vaqt: {formatTime(result.time_taken_seconds)}
            </p>
          )}
          <div className="flex gap-3 pt-2">
            <Button variant="outline" className="flex-1" onClick={() => router.push("/dashboard/results")}>Natijalar</Button>
            <Button className="flex-1" onClick={() => router.push("/dashboard/olympiads")}>Olimpiadalar</Button>
          </div>
        </div>
      </div>
    );
  }

  // === EXAM PHASE ===
  if (phase === "exam" && questions.length > 0) {
    const question = questions[currentIdx];
    const selectedOpt = question ? answers[question.id] : undefined;
    const answered = Object.keys(answers).length;
    const pct = (answered / questions.length) * 100;
    const urgent = timeLeft < 60;

    return (
      <div className="max-w-2xl mx-auto py-6 space-y-6">
        {/* Timer + Progress */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">{answered}/{questions.length} javob berildi</span>
          </div>
          <div className={`flex items-center gap-2 rounded-full px-4 py-1.5 text-sm font-mono font-bold ${urgent ? "bg-red-500/10 text-red-600" : "bg-primary/10 text-primary"}`}>
            <Clock className="h-4 w-4" />
            {formatTime(timeLeft)}
          </div>
        </div>
        <Progress value={pct} className="h-2" />

        {/* Question */}
        {question && (
          <div className="rounded-2xl border border-border bg-card p-6 space-y-5">
            <div className="flex items-start gap-3">
              <span className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-primary/10 text-sm font-bold text-primary">
                {currentIdx + 1}
              </span>
              <p className="font-medium text-foreground text-base leading-relaxed">{question.text}</p>
            </div>

            <div className="space-y-3">
              {question.options.map(opt => {
                const isSelected = selectedOpt === opt.id;
                return (
                  <button
                    key={opt.id}
                    onClick={() => handleAnswer(opt.id)}
                    className={`w-full flex items-center gap-3 rounded-xl border p-4 text-left transition-all ${
                      isSelected
                        ? "border-primary bg-primary/5 text-foreground"
                        : "border-border bg-background hover:border-primary/30 hover:bg-primary/5 text-foreground"
                    }`}
                  >
                    <span className={`flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-xs font-bold border ${
                      isSelected ? "border-primary bg-primary text-primary-foreground" : "border-border text-muted-foreground"
                    }`}>
                      {opt.label}
                    </span>
                    <span className="text-sm leading-relaxed">{opt.text}</span>
                    {isSelected && <CheckCircle2 className="h-4 w-4 text-primary ml-auto flex-shrink-0" />}
                  </button>
                );
              })}
            </div>

            {question.points > 0 && (
              <p className="text-xs text-muted-foreground text-right">{question.points} ball</p>
            )}
          </div>
        )}

        {/* Navigation */}
        <div className="flex items-center justify-between">
          <Button variant="outline" onClick={() => setCurrentIdx(i => Math.max(0, i - 1))} disabled={currentIdx === 0}>
            <ArrowLeft className="h-4 w-4 mr-1" /> Oldingi
          </Button>

          <div className="flex gap-1 flex-wrap justify-center max-w-xs">
            {questions.map((q, i) => (
              <button
                key={q.id}
                onClick={() => setCurrentIdx(i)}
                className={`h-7 w-7 rounded-md text-xs font-medium transition-colors ${
                  i === currentIdx ? "bg-primary text-primary-foreground"
                    : answers[q.id] ? "bg-green-500/20 text-green-600"
                    : "bg-muted text-muted-foreground hover:bg-muted/80"
                }`}
              >
                {i + 1}
              </button>
            ))}
          </div>

          {currentIdx < questions.length - 1 ? (
            <Button onClick={() => setCurrentIdx(i => i + 1)}>
              Keyingi <ArrowRight className="h-4 w-4 ml-1" />
            </Button>
          ) : (
            <Button variant="destructive" onClick={() => setFinishConfirm(true)} disabled={finishing} className="gap-2">
              <Flag className="h-4 w-4" /> Yakunlash
            </Button>
          )}
        </div>

        {/* Always visible finish button when not on last question */}
        {currentIdx < questions.length - 1 && answered > 0 && (
          <div className="text-center">
            <Button variant="outline" size="sm" className="gap-2 text-destructive border-destructive/30 hover:bg-destructive/5" onClick={() => setFinishConfirm(true)}>
              <Flag className="h-3.5 w-3.5" /> Erta yakunlash
            </Button>
          </div>
        )}

        {/* Finish confirmation dialog */}
        <Dialog open={finishConfirm} onOpenChange={setFinishConfirm}>
          <DialogContent className="max-w-sm">
            <DialogHeader><DialogTitle>Testni yakunlash</DialogTitle></DialogHeader>
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">
                Jami: {questions.length} savol, {answered} ta javob berildi.
              </p>
              {answered < questions.length && (
                <p className="text-sm text-yellow-600 flex items-center gap-1.5">
                  <AlertCircle className="h-4 w-4" />
                  {questions.length - answered} ta savol javobsiz qolmoqda
                </p>
              )}
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setFinishConfirm(false)}>Davom etish</Button>
              <Button variant="destructive" onClick={() => handleFinish()} disabled={finishing}>
                {finishing && <Loader2 className="h-4 w-4 animate-spin mr-2" />}Yakunlash
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    );
  }

  // === PAYMENT PHASE ===
  if (phase === "payment") {
    const canPay = balance >= (olympiad.price || 0);
    return (
      <div className="max-w-md mx-auto py-10">
        <Button variant="ghost" size="sm" className="mb-6 gap-2" onClick={() => setPhase("detail")}>
          <ArrowLeft className="h-4 w-4" /> Orqaga
        </Button>
        <div className="rounded-2xl border border-border bg-card p-8 space-y-6">
          <div className="text-center">
            <div className="inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-orange-500/10 mb-4">
              <DollarSign className="h-8 w-8 text-orange-600" />
            </div>
            <h2 className="text-xl font-bold text-foreground">To&apos;lovli olimpiada</h2>
            <p className="text-sm text-muted-foreground mt-2">{olympiad.title}</p>
          </div>

          <div className="rounded-xl border border-border bg-muted/40 p-4 space-y-3">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Olimpiada narxi</span>
              <span className="font-semibold text-foreground">{olympiad.price?.toLocaleString()} UZS</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Balansingiz</span>
              <span className={`font-semibold ${canPay ? "text-green-600" : "text-red-500"}`}>{balance.toLocaleString()} UZS</span>
            </div>
            {!canPay && (
              <div className="flex justify-between text-sm border-t border-border pt-2">
                <span className="text-muted-foreground">Yetishmaydi</span>
                <span className="font-semibold text-red-500">{((olympiad.price || 0) - balance).toLocaleString()} UZS</span>
              </div>
            )}
          </div>

          {canPay ? (
            <Button className="w-full gap-2" onClick={handlePayAndJoin}>
              <DollarSign className="h-4 w-4" />
              To&apos;lov qilish va qatnashish
            </Button>
          ) : (
            <div className="space-y-3">
              <p className="text-sm text-center text-muted-foreground">Balansingizni to&apos;ldiring va qatnashing</p>
              <Button className="w-full" onClick={() => router.push("/dashboard/balance")}>
                Balansni to&apos;ldirish
              </Button>
            </div>
          )}
        </div>
      </div>
    );
  }

  // === READY PHASE (joined, ready to start) ===
  if (phase === "ready") {
    return (
      <div className="max-w-md mx-auto py-10">
        <div className="rounded-2xl border border-border bg-card p-8 space-y-6 text-center">
          <div className="inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-green-500/10">
            <CheckCircle2 className="h-8 w-8 text-green-600" />
          </div>
          <div>
            <h2 className="text-xl font-bold text-foreground">Tayyor!</h2>
            <p className="text-sm text-muted-foreground mt-2">{olympiad.title}</p>
          </div>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div className="rounded-xl border border-border p-3 text-center">
              <Clock className="h-5 w-5 mx-auto mb-1 text-muted-foreground" />
              <p className="font-semibold text-foreground">{olympiad.duration_minutes} daqiqa</p>
            </div>
            <div className="rounded-xl border border-border p-3 text-center">
              <Trophy className="h-5 w-5 mx-auto mb-1 text-muted-foreground" />
              <p className="font-semibold text-foreground">{olympiad.total_questions} savol</p>
            </div>
          </div>
          <div className="rounded-xl border border-yellow-500/20 bg-yellow-500/5 p-4 text-left space-y-1.5 text-sm">
            <p className="flex items-start gap-2 text-foreground"><AlertCircle className="h-4 w-4 text-yellow-600 flex-shrink-0 mt-0.5" />Test boshlanganidan keyin vaqt hisoblanadi</p>
            <p className="flex items-start gap-2 text-foreground"><AlertCircle className="h-4 w-4 text-yellow-600 flex-shrink-0 mt-0.5" />Vaqt tugaganda test avtomatik yakunlanadi</p>
            {olympiad.rules && <p className="flex items-start gap-2 text-foreground"><AlertCircle className="h-4 w-4 text-yellow-600 flex-shrink-0 mt-0.5" />{olympiad.rules}</p>}
          </div>
          <Button className="w-full gap-2" onClick={handleStart}>
            <Trophy className="h-4 w-4" /> Testni boshlash
          </Button>
        </div>
      </div>
    );
  }

  // === DETAIL PHASE (default) ===
  return (
    <div className="max-w-2xl mx-auto py-6 space-y-6">
      <Button variant="ghost" size="sm" className="gap-2" onClick={() => router.push("/dashboard/olympiads")}>
        <ArrowLeft className="h-4 w-4" /> Olimpiadalar
      </Button>

      <div className="rounded-2xl border border-border bg-card p-8 space-y-6">
        {/* Header */}
        <div className="flex items-start gap-4">
          <div className="flex h-14 w-14 flex-shrink-0 items-center justify-center rounded-2xl bg-primary/10">
            <Trophy className="h-7 w-7 text-primary" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex flex-wrap gap-2 mb-2">
              <Badge variant="outline" className={`text-xs ${olympiad.is_paid ? "bg-orange-500/10 text-orange-600 border-orange-500/20" : "bg-green-500/10 text-green-600 border-green-500/20"}`}>
                {olympiad.is_paid ? `${olympiad.price?.toLocaleString()} UZS` : "Bepul"}
              </Badge>
              {olympiad.subject && <Badge variant="outline" className="text-xs">{olympiad.subject}</Badge>}
              {olympiad.grade > 0 && <Badge variant="outline" className="text-xs">{olympiad.grade}-sinf</Badge>}
            </div>
            <h1 className="text-xl font-bold text-foreground">{olympiad.title}</h1>
          </div>
        </div>

        {olympiad.description && (
          <p className="text-sm text-muted-foreground leading-relaxed">{olympiad.description}</p>
        )}

        {/* Info grid */}
        <div className="grid grid-cols-2 gap-4">
          <div className="rounded-xl border border-border p-4 text-center">
            <Clock className="h-5 w-5 mx-auto mb-2 text-muted-foreground" />
            <p className="font-semibold text-foreground">{olympiad.duration_minutes} daqiqa</p>
            <p className="text-xs text-muted-foreground">Davomiyligi</p>
          </div>
          <div className="rounded-xl border border-border p-4 text-center">
            <Trophy className="h-5 w-5 mx-auto mb-2 text-muted-foreground" />
            <p className="font-semibold text-foreground">{olympiad.total_questions ?? olympiad.questions_count ?? "—"} ta</p>
            <p className="text-xs text-muted-foreground">Savollar soni</p>
          </div>
        </div>

        {/* Dates */}
        {(olympiad.start_time || olympiad.start_date) && (
          <div className="grid grid-cols-2 gap-4 text-sm">
            {(olympiad.start_time || olympiad.start_date) && (
              <div className="rounded-xl border border-border p-3">
                <p className="text-xs text-muted-foreground mb-1">Boshlanish</p>
                <p className="font-medium text-foreground">
                  {new Date(olympiad.start_time || olympiad.start_date).toLocaleDateString("uz-UZ")}
                </p>
              </div>
            )}
            {(olympiad.end_time || olympiad.end_date) && (
              <div className="rounded-xl border border-border p-3">
                <p className="text-xs text-muted-foreground mb-1">Tugash</p>
                <p className="font-medium text-foreground">
                  {new Date(olympiad.end_time || olympiad.end_date).toLocaleDateString("uz-UZ")}
                </p>
              </div>
            )}
          </div>
        )}

        {/* Rules */}
        {olympiad.rules && (
          <div className="rounded-xl border border-border p-4">
            <p className="text-sm font-semibold text-foreground mb-2">Qoidalar</p>
            <p className="text-sm text-muted-foreground whitespace-pre-wrap">{olympiad.rules}</p>
          </div>
        )}

        {/* CTA */}
        {olympiad.status === "ended" || olympiad.status === "completed" ? (
          <Button disabled className="w-full">Olimpiada tugagan</Button>
        ) : (
          <Button className="w-full gap-2" onClick={handleJoin} size="lg">
            <Trophy className="h-4 w-4" />
            {olympiad.is_paid ? `To'lov qilish (${olympiad.price?.toLocaleString()} UZS)` : "Qatnashish"}
          </Button>
        )}
      </div>
    </div>
  );
}
