"use client";

import Link from "next/link";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { MaintenanceBanner } from "@/components/MaintenanceBanner";
import { useSettings } from "@/lib/settings-context";
import { Trophy, ClipboardCheck, BarChart3, ArrowRight, Zap, Target, Users } from "lucide-react";

export default function Home() {
  const settings = useSettings();

  if (settings.maintenance_mode) {
    return <MaintenanceBanner />;
  }

  return (
    <div className="min-h-screen bg-background">
      <Header />

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-background to-primary/10" />
        <div className="relative max-w-6xl mx-auto px-4 py-24 md:py-36 text-center">
          <div className="inline-flex items-center gap-2 rounded-full bg-primary/10 border border-primary/20 px-4 py-1.5 text-sm text-primary mb-6">
            <Zap className="h-4 w-4" />
            Onlayn test va olimpiada platformasi
          </div>
          <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold text-foreground tracking-tight mb-6">
            Next<span className="text-primary">Olymp</span>
          </h1>
          <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-10">
            Olimpiada va imtihonlarga tayyorlov onlayn test platformasi.
            O&apos;z bilimingizni sinab ko&apos;ring va reytingda o&apos;rningizni bilib oling.
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link
              href="/dashboard/mock-tests"
              className="inline-flex items-center gap-2 rounded-xl bg-primary text-primary-foreground px-8 py-3.5 text-base font-semibold hover:bg-primary/90 transition-colors shadow-lg shadow-primary/25"
            >
              <ClipboardCheck className="h-5 w-5" />
              Testni boshlash
            </Link>
            <Link
              href="/dashboard"
              className="inline-flex items-center gap-2 rounded-xl border border-border bg-card text-foreground px-8 py-3.5 text-base font-semibold hover:bg-accent transition-colors"
            >
              <Trophy className="h-5 w-5" />
              Olimpiadalarni ko&apos;rish
            </Link>
          </div>
        </div>
      </section>

      {/* Stats */}
      <section className="border-y border-border bg-card/50">
        <div className="max-w-6xl mx-auto px-4 py-12 grid grid-cols-1 sm:grid-cols-3 gap-8 text-center">
          <div>
            <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center mx-auto mb-3">
              <Target className="h-6 w-6 text-primary" />
            </div>
            <p className="text-2xl font-bold text-foreground">Mock testlar</p>
            <p className="text-sm text-muted-foreground mt-1">Istalgan vaqtda ishlash mumkin</p>
          </div>
          <div>
            <div className="h-12 w-12 rounded-xl bg-amber-500/10 flex items-center justify-center mx-auto mb-3">
              <Trophy className="h-6 w-6 text-amber-500" />
            </div>
            <p className="text-2xl font-bold text-foreground">Olimpiadalar</p>
            <p className="text-sm text-muted-foreground mt-1">Fan olimpiadalarida qatnashing</p>
          </div>
          <div>
            <div className="h-12 w-12 rounded-xl bg-green-500/10 flex items-center justify-center mx-auto mb-3">
              <Users className="h-6 w-6 text-green-500" />
            </div>
            <p className="text-2xl font-bold text-foreground">Reyting</p>
            <p className="text-sm text-muted-foreground mt-1">Respublika bo&apos;yicha o&apos;rningizni bilib oling</p>
          </div>
        </div>
      </section>

      {/* Preview Blocks */}
      <section className="max-w-6xl mx-auto px-4 py-20">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Mock test */}
          <div className="group rounded-2xl border border-border bg-card p-6 hover:border-primary/30 hover:shadow-lg transition-all">
            <div className="h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center mb-4">
              <ClipboardCheck className="h-6 w-6 text-primary" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">Mock testlar</h3>
            <p className="text-sm text-muted-foreground mb-4">
              Fan bo&apos;yicha testlarni istalgan vaqtda ishlang. Timer, progress bar va natija — hammasi bir joyda.
            </p>
            <Link href="/dashboard/mock-tests" className="inline-flex items-center gap-1 text-sm text-primary font-medium group-hover:gap-2 transition-all">
              Boshqalar <ArrowRight className="h-4 w-4" />
            </Link>
          </div>

          {/* Olimpiada */}
          <div className="group rounded-2xl border border-border bg-card p-6 hover:border-amber-500/30 hover:shadow-lg transition-all">
            <div className="h-12 w-12 rounded-xl bg-amber-500/10 flex items-center justify-center mb-4">
              <Trophy className="h-6 w-6 text-amber-500" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">Olimpiadalar</h3>
            <p className="text-sm text-muted-foreground mb-4">
              Vaqtga bog&apos;liq test olimpiadalarida qatnashing. Live ranking va umumiy reyting.
            </p>
            <Link href="/dashboard" className="inline-flex items-center gap-1 text-sm text-amber-500 font-medium group-hover:gap-2 transition-all">
              Ko&apos;rish <ArrowRight className="h-4 w-4" />
            </Link>
          </div>

          {/* Reyting */}
          <div className="group rounded-2xl border border-border bg-card p-6 hover:border-green-500/30 hover:shadow-lg transition-all">
            <div className="h-12 w-12 rounded-xl bg-green-500/10 flex items-center justify-center mb-4">
              <BarChart3 className="h-6 w-6 text-green-500" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">Reyting</h3>
            <p className="text-sm text-muted-foreground mb-4">
              Eng yuqori ball to&apos;plaganlar. Fan va sinf bo&apos;yicha filter. Haftalik va oylik reyting.
            </p>
            <Link href="/dashboard/leaderboard" className="inline-flex items-center gap-1 text-sm text-green-500 font-medium group-hover:gap-2 transition-all">
              Reytingni ko&apos;rish <ArrowRight className="h-4 w-4" />
            </Link>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="bg-gradient-to-r from-primary/10 to-primary/5 border-t border-border">
        <div className="max-w-4xl mx-auto px-4 py-16 text-center">
          <h2 className="text-2xl md:text-3xl font-bold text-foreground mb-4">
            Bilimingizni sinab ko&apos;rishga tayyormisiz?
          </h2>
          <p className="text-muted-foreground mb-8">
            Ro&apos;yxatdan o&apos;ting va birinchi testingizni hoziroq boshlang
          </p>
          <Link
            href="/register"
            className="inline-flex items-center gap-2 rounded-xl bg-primary text-primary-foreground px-8 py-3.5 text-base font-semibold hover:bg-primary/90 transition-colors shadow-lg shadow-primary/25"
          >
            Ro&apos;yxatdan o&apos;tish
            <ArrowRight className="h-5 w-5" />
          </Link>
        </div>
      </section>

      <Footer />
    </div>
  );
}
