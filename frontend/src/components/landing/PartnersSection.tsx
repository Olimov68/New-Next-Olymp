"use client";

import { useI18n } from "@/lib/i18n";
import { Handshake } from "lucide-react";

const partners = [
  { name: "UNICEF", color: "from-blue-500/10 to-cyan-500/10 border-blue-500/20 text-blue-600 dark:text-blue-400" },
  { name: "British Council", color: "from-indigo-500/10 to-blue-500/10 border-indigo-500/20 text-indigo-600 dark:text-indigo-400" },
  { name: "Samsung", color: "from-slate-500/10 to-gray-500/10 border-slate-500/20 text-slate-600 dark:text-slate-400" },
  { name: "Google Education", color: "from-green-500/10 to-emerald-500/10 border-green-500/20 text-green-600 dark:text-green-400" },
  { name: "Khan Academy", color: "from-violet-500/10 to-purple-500/10 border-violet-500/20 text-violet-600 dark:text-violet-400" },
  { name: "UNESCO", color: "from-sky-500/10 to-blue-500/10 border-sky-500/20 text-sky-600 dark:text-sky-400" },
  { name: "Coursera", color: "from-blue-500/10 to-indigo-500/10 border-blue-500/20 text-blue-600 dark:text-blue-400" },
  { name: "MIT OpenCourseWare", color: "from-red-500/10 to-orange-500/10 border-red-500/20 text-red-600 dark:text-red-400" },
];

export function PartnersSection() {
  const { t, lang } = useI18n();

  const title = lang === "ru" ? "Наши партнёры" : lang === "en" ? "Our Partners" : "Hamkorlarimiz";
  const badge = lang === "ru" ? "Партнёры" : lang === "en" ? "Partners" : "Hamkorlar";
  const desc =
    lang === "ru"
      ? "Мы сотрудничаем с ведущими международными организациями для обеспечения высочайшего качества образования"
      : lang === "en"
      ? "We collaborate with leading international organizations to ensure the highest quality of education"
      : "Biz yetakchi xalqaro tashkilotlar bilan hamkorlik qilamiz — ta'limning eng yuqori sifatini ta'minlash uchun";

  return (
    <section id="partners" className="py-20 bg-background border-t border-border">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <div className="inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/10 px-4 py-1.5 text-sm text-primary mb-4">
            <Handshake className="h-4 w-4" />
            {badge}
          </div>
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">{title}</h2>
          <p className="text-muted-foreground max-w-lg mx-auto text-base">{desc}</p>
        </div>

        {/* Partners grid */}
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4 max-w-4xl mx-auto">
          {partners.map((p) => (
            <div
              key={p.name}
              className={`group flex h-24 items-center justify-center rounded-2xl border bg-gradient-to-br ${p.color} text-center px-3 font-semibold text-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-md cursor-default`}
            >
              {p.name}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
