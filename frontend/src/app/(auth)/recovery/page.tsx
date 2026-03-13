"use client";

import { useState } from "react";
import Link from "next/link";
import { useI18n } from "@/lib/i18n";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowLeft, Bot, CheckCircle2, AlertCircle, Loader2, Send } from "lucide-react";

type RecoveryStep = "identify" | "bot_verify" | "reset_password" | "success";

export default function RecoveryPage() {
  const { t } = useI18n();
  const [step, setStep] = useState<RecoveryStep>("identify");
  const [identifier, setIdentifier] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [verificationCode, setVerificationCode] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [maskedContact, setMaskedContact] = useState("");

  const handleIdentify = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      // Simulate API call - in production this would call backend
      // POST /auth/recovery/identify { identifier }
      // Backend checks if user exists, sends code via Telegram bot
      await new Promise(resolve => setTimeout(resolve, 1500));
      setMaskedContact("@t***gram_user");
      setStep("bot_verify");
    } catch {
      setError(t("auth.recovery_user_not_found"));
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyCode = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (verificationCode.length < 6) {
      setError(t("auth.recovery_invalid_code"));
      return;
    }
    setLoading(true);
    try {
      // POST /auth/recovery/verify { identifier, code }
      await new Promise(resolve => setTimeout(resolve, 1500));
      setStep("reset_password");
    } catch {
      setError(t("auth.recovery_wrong_code"));
    } finally {
      setLoading(false);
    }
  };

  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (newPassword.length < 8) {
      setError(t("auth.recovery_password_short"));
      return;
    }
    if (newPassword !== confirmPassword) {
      setError(t("auth.password_mismatch"));
      return;
    }
    setLoading(true);
    try {
      // POST /auth/recovery/reset { identifier, code, new_password }
      await new Promise(resolve => setTimeout(resolve, 1500));
      setStep("success");
    } catch {
      setError(t("common.error"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background hero-mesh px-4">
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-500/10 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-indigo-500/10 rounded-full blur-3xl" />
      </div>
      <Card className="relative w-full max-w-md shadow-2xl border border-white/10 bg-white/5 backdrop-blur-xl">
        <CardHeader className="text-center pb-2">
          <Link href="/" className="flex justify-center mb-3">
            <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500 to-indigo-600 text-white font-bold text-xl shadow-lg shadow-blue-500/25">
              NO
            </div>
          </Link>
          <CardTitle className="text-2xl text-white">{t("auth.account_recovery")}</CardTitle>
          <p className="text-sm text-blue-200/60 mt-1">{t("auth.recovery_desc")}</p>
        </CardHeader>
        <CardContent>
          {/* Step indicators */}
          <div className="flex items-center justify-center gap-2 mb-6">
            {(["identify", "bot_verify", "reset_password", "success"] as RecoveryStep[]).map((s, i) => (
              <div key={s} className="flex items-center gap-2">
                <div className={`h-2 w-2 rounded-full transition-colors ${
                  step === s ? "bg-blue-400 ring-4 ring-blue-400/20" :
                  (["identify", "bot_verify", "reset_password", "success"].indexOf(step) > i ? "bg-blue-400" : "bg-white/20")
                }`} />
                {i < 3 && <div className={`w-8 h-0.5 ${
                  ["identify", "bot_verify", "reset_password", "success"].indexOf(step) > i ? "bg-blue-400" : "bg-white/10"
                }`} />}
              </div>
            ))}
          </div>

          {error && (
            <div className="rounded-xl bg-red-500/10 border border-red-500/20 text-red-300 text-sm p-3 mb-4 flex items-center gap-2">
              <AlertCircle className="h-4 w-4 shrink-0" />
              {error}
            </div>
          )}

          {/* Step 1: Identify user */}
          {step === "identify" && (
            <form onSubmit={handleIdentify} className="space-y-4">
              <div className="rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-sm text-blue-200">
                <div className="flex items-start gap-3">
                  <Bot className="h-5 w-5 text-blue-400 mt-0.5 shrink-0" />
                  <p>{t("auth.recovery_bot_info")}</p>
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-blue-100/80 text-xs">{t("auth.recovery_identifier")}</Label>
                <Input
                  className="bg-white/5 border-white/10 text-white placeholder:text-white/30 focus:border-blue-400/50 focus:ring-blue-400/20"
                  value={identifier}
                  onChange={(e) => setIdentifier(e.target.value)}
                  placeholder="username, email yoki telefon"
                  required
                />
              </div>
              <Button type="submit" className="w-full bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white shadow-lg shadow-blue-500/25 border-0 gap-2" disabled={loading}>
                {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
                {loading ? t("common.loading") : t("auth.recovery_send_code")}
              </Button>
            </form>
          )}

          {/* Step 2: Bot verification */}
          {step === "bot_verify" && (
            <form onSubmit={handleVerifyCode} className="space-y-4">
              <div className="rounded-xl bg-green-500/10 border border-green-500/20 p-4 text-sm text-green-200">
                <div className="flex items-start gap-3">
                  <CheckCircle2 className="h-5 w-5 text-green-400 mt-0.5 shrink-0" />
                  <div>
                    <p className="font-medium">{t("auth.recovery_code_sent")}</p>
                    <p className="text-green-200/70 mt-1">{t("auth.recovery_check_bot")} {maskedContact}</p>
                  </div>
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-blue-100/80 text-xs">{t("auth.recovery_enter_code")}</Label>
                <Input
                  className="bg-white/5 border-white/10 text-white placeholder:text-white/30 focus:border-blue-400/50 focus:ring-blue-400/20 text-center text-lg tracking-widest"
                  value={verificationCode}
                  onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                  placeholder="000000"
                  maxLength={6}
                  required
                />
              </div>
              <Button type="submit" className="w-full bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white shadow-lg shadow-blue-500/25 border-0" disabled={loading}>
                {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
                {loading ? t("common.loading") : t("auth.recovery_verify")}
              </Button>
              <button type="button" onClick={() => setStep("identify")} className="w-full text-center text-sm text-blue-300/50 hover:text-blue-300 transition-colors">
                {t("auth.recovery_resend")}
              </button>
            </form>
          )}

          {/* Step 3: Reset password */}
          {step === "reset_password" && (
            <form onSubmit={handleResetPassword} className="space-y-4">
              <div className="space-y-2">
                <Label className="text-blue-100/80 text-xs">{t("auth.recovery_new_password")}</Label>
                <Input
                  className="bg-white/5 border-white/10 text-white placeholder:text-white/30 focus:border-blue-400/50 focus:ring-blue-400/20"
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label className="text-blue-100/80 text-xs">{t("auth.confirm_password")}</Label>
                <Input
                  className="bg-white/5 border-white/10 text-white placeholder:text-white/30 focus:border-blue-400/50 focus:ring-blue-400/20"
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                />
              </div>
              <Button type="submit" className="w-full bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white shadow-lg shadow-blue-500/25 border-0" disabled={loading}>
                {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
                {loading ? t("common.loading") : t("auth.recovery_reset")}
              </Button>
            </form>
          )}

          {/* Step 4: Success */}
          {step === "success" && (
            <div className="text-center space-y-4">
              <div className="flex justify-center">
                <div className="h-16 w-16 rounded-full bg-green-500/10 border border-green-500/20 flex items-center justify-center">
                  <CheckCircle2 className="h-8 w-8 text-green-400" />
                </div>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white">{t("auth.recovery_success")}</h3>
                <p className="text-sm text-blue-200/60 mt-1">{t("auth.recovery_success_desc")}</p>
              </div>
              <Link href="/login">
                <Button className="w-full bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white shadow-lg shadow-blue-500/25 border-0 gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  {t("auth.login")}
                </Button>
              </Link>
            </div>
          )}

          <div className="mt-6 text-center">
            <Link href="/login" className="text-sm text-blue-300/50 hover:text-blue-300 transition-colors inline-flex items-center gap-1">
              <ArrowLeft className="h-3 w-3" />
              {t("auth.back_to_login")}
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
