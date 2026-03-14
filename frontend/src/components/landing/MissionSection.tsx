"use client";

import { GraduationCap, Quote, ArrowRight, Star, Trophy, BookOpen, Sparkles } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useI18n } from "@/lib/i18n";

const content = {
  uz: {
    badge: "Missiyamiz",
    title: "Olimpiada – bu imkoniyat, biz esa yo'ldoshmiz!",
    text: "Bizning jamoamiz 2014-yildan beri xalqaro va Respublika miqyosidagi olimpiadalarda yuqori natijalarga erishgan, tajribali mutaxassislar tomonidan asos solingan. Bizning maqsadimiz — har bir o'quvchining bilim va iqtidorini yuksaltirish, uni raqobatbardosh va motivatsion muhitda tarbiyalashdir.",
    quote: "Menga qolsa har kuni olimpiada tashkil qil, chunki mening bolalarim ertaga olimpiada bor deb, bugun kitob o'qiydi. Menga aynan shu kerak.",
    afterQuote: "Aynan shu ruhda biz ilk sinfdan boshlab har bir o'quvchining salohiyatini ochish uchun imkoniyat yaratamiz.",
    cta: "Platformaga qo'shilish",
    values: [
      { icon: Star, label: "Sifat va ishonch" },
      { icon: Trophy, label: "Adolatli musobaqa" },
      { icon: BookOpen, label: "Chuqur bilim" },
      { icon: Sparkles, label: "Ilhom va motivatsiya" },
    ],
  },
  ru: {
    badge: "Наша миссия",
    title: "Олимпиада – это возможность, а мы – ваши спутники!",
    text: "Наша команда основана опытными специалистами, добившимися высоких результатов на международных и республиканских олимпиадах с 2014 года. Наша цель — развивать знания и таланты каждого ученика в конкурентной и мотивирующей среде.",
    quote: "Если бы это зависело от меня, я бы проводил олимпиады каждый день, потому что мои дети читают книги, когда знают, что завтра олимпиада. Именно это мне и нужно.",
    afterQuote: "Именно в этом духе мы создаём возможности для раскрытия потенциала каждого ученика, начиная с первого класса.",
    cta: "Присоединиться к платформе",
    values: [
      { icon: Star, label: "Качество и доверие" },
      { icon: Trophy, label: "Честное соревнование" },
      { icon: BookOpen, label: "Глубокие знания" },
      { icon: Sparkles, label: "Вдохновение и мотивация" },
    ],
  },
  en: {
    badge: "Our Mission",
    title: "Olympiad is an opportunity, and we are your companions!",
    text: "Our team was founded by experienced specialists who have achieved high results in international and national olympiads since 2014. Our goal is to elevate every student's knowledge and talent in a competitive and motivating environment.",
    quote: "If it were up to me, I'd organize olympiads every day, because my children read books when they know there's an olympiad tomorrow. That's exactly what I need.",
    afterQuote: "In this very spirit, we create opportunities to unlock every student's potential, starting from the first grade.",
    cta: "Join the Platform",
    values: [
      { icon: Star, label: "Quality & Trust" },
      { icon: Trophy, label: "Fair Competition" },
      { icon: BookOpen, label: "Deep Knowledge" },
      { icon: Sparkles, label: "Inspiration & Motivation" },
    ],
  },
};

export function MissionSection() {
  const { lang } = useI18n();
  const c = content[lang] || content.uz;

  return (
    <section className="relative py-20 overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-background via-primary/[0.03] to-background" />
      <div className="container mx-auto px-4 relative">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left: Text Content */}
          <div>
            <div className="inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/10 px-4 py-1.5 text-sm text-primary mb-6">
              <GraduationCap className="h-4 w-4" />
              {c.badge}
            </div>

            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-6 leading-tight">
              {c.title}
            </h2>

            <p className="text-muted-foreground text-lg leading-relaxed mb-8">
              {c.text}
            </p>

            {/* Quote Block */}
            <div className="relative rounded-2xl border border-primary/20 bg-primary/5 p-6 mb-6">
              <Quote className="absolute top-4 left-4 h-8 w-8 text-primary/20" />
              <p className="text-foreground font-medium italic pl-8 text-base leading-relaxed">
                &ldquo;{c.quote}&rdquo;
              </p>
            </div>

            <p className="text-muted-foreground mb-8">
              {c.afterQuote}
            </p>

            <Link href="/register">
              <Button size="lg" className="gap-2 bg-gradient-to-r from-blue-500 to-indigo-600 text-white hover:from-blue-600 hover:to-indigo-700 shadow-lg shadow-blue-500/25 border-0">
                {c.cta}
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>

          {/* Right: Values Card */}
          <div className="relative">
            <div className="rounded-3xl border border-border bg-card p-8 shadow-xl">
              <div className="text-center mb-8">
                <div className="inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500 to-indigo-600 text-white mb-4 shadow-lg shadow-blue-500/25">
                  <GraduationCap className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-bold text-foreground">
                  {lang === "ru" ? "Наши ценности" : lang === "en" ? "Our Values" : "Qadriyatlarimiz"}
                </h3>
              </div>

              <div className="grid grid-cols-2 gap-4">
                {c.values.map((value, i) => (
                  <div
                    key={i}
                    className="flex flex-col items-center gap-3 rounded-xl border border-border bg-background p-4 text-center transition-colors hover:border-primary/20 hover:bg-primary/5"
                  >
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                      <value.icon className="h-5 w-5 text-primary" />
                    </div>
                    <span className="text-sm font-medium text-foreground">{value.label}</span>
                  </div>
                ))}
              </div>

              {/* Stats inside card */}
              <div className="mt-6 pt-6 border-t border-border grid grid-cols-3 gap-4 text-center">
                <div>
                  <div className="text-2xl font-bold text-foreground">10+</div>
                  <div className="text-xs text-muted-foreground">{lang === "ru" ? "Лет опыта" : lang === "en" ? "Years exp." : "Yillik tajriba"}</div>
                </div>
                <div>
                  <div className="text-2xl font-bold text-foreground">50K+</div>
                  <div className="text-xs text-muted-foreground">{lang === "ru" ? "Учеников" : lang === "en" ? "Students" : "O'quvchilar"}</div>
                </div>
                <div>
                  <div className="text-2xl font-bold text-foreground">14</div>
                  <div className="text-xs text-muted-foreground">{lang === "ru" ? "Регионов" : lang === "en" ? "Regions" : "Viloyatlar"}</div>
                </div>
              </div>
            </div>

            {/* Decorative elements */}
            <div className="absolute -top-4 -right-4 h-24 w-24 rounded-full bg-gradient-to-br from-blue-500/10 to-indigo-500/10 blur-2xl" />
            <div className="absolute -bottom-4 -left-4 h-24 w-24 rounded-full bg-gradient-to-br from-purple-500/10 to-pink-500/10 blur-2xl" />
          </div>
        </div>
      </div>
    </section>
  );
}
