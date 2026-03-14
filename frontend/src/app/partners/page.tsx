"use client";

import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { AnnouncementBar } from "@/components/AnnouncementBar";
import { useI18n } from "@/lib/i18n";
import { Handshake, ExternalLink, Globe, Mail, ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";

const partners = [
  {
    name: "UNICEF",
    category: { uz: "Xalqaro tashkilot", ru: "Международная организация", en: "International Organization" },
    desc: {
      uz: "Bolalar huquqlari va ta'lim sohalarida global hamkor. UNICEF bilan birgalikda maktab o'quvchilarining ta'lim imkoniyatlarini kengaytiramiz.",
      ru: "Глобальный партнёр в области прав детей и образования. Совместно с UNICEF расширяем образовательные возможности школьников.",
      en: "A global partner in children's rights and education. Together with UNICEF, we expand learning opportunities for students.",
    },
    color: "from-blue-500/10 to-cyan-500/10 border-blue-500/20",
    iconColor: "text-blue-600 dark:text-blue-400",
    badgeColor: "bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20",
    website: "https://unicef.org",
  },
  {
    name: "British Council",
    category: { uz: "Ta'lim tashkiloti", ru: "Образовательная организация", en: "Educational Organization" },
    desc: {
      uz: "Buyuk Britaniyaning xalqaro ta'lim va madaniyat tashkiloti. Ingliz tili va xalqaro standartlar bo'yicha hamkorlik.",
      ru: "Международная организация Великобритании в сфере образования и культуры. Сотрудничество по английскому языку и международным стандартам.",
      en: "UK's international organization for education and culture. Partnership on English language and international standards.",
    },
    color: "from-indigo-500/10 to-blue-500/10 border-indigo-500/20",
    iconColor: "text-indigo-600 dark:text-indigo-400",
    badgeColor: "bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border-indigo-500/20",
    website: "https://britishcouncil.org",
  },
  {
    name: "Samsung",
    category: { uz: "Texnologiya hamkori", ru: "Технологический партнёр", en: "Technology Partner" },
    desc: {
      uz: "Zamonaviy texnologiyalar sohasida strategik hamkor. Samsung bilan raqamli ta'lim infratuzilmasini rivojlantirmoqdamiz.",
      ru: "Стратегический партнёр в области современных технологий. Совместно с Samsung развиваем цифровую образовательную инфраструктуру.",
      en: "Strategic partner in modern technologies. Together with Samsung, we develop digital education infrastructure.",
    },
    color: "from-slate-500/10 to-gray-500/10 border-slate-500/20",
    iconColor: "text-slate-600 dark:text-slate-400",
    badgeColor: "bg-slate-500/10 text-slate-600 dark:text-slate-300 border-slate-500/20",
    website: "https://samsung.com",
  },
  {
    name: "Google Education",
    category: { uz: "Raqamli ta'lim", ru: "Цифровое образование", en: "Digital Education" },
    desc: {
      uz: "Google'ning ta'lim bo'yicha maxsus dasturlari orqali o'quvchilar va o'qituvchilarga zamonaviy raqamli vositalar taqdim etamiz.",
      ru: "Через образовательные программы Google предоставляем учащимся и преподавателям современные цифровые инструменты.",
      en: "Through Google's education programs, we provide students and teachers with modern digital tools.",
    },
    color: "from-green-500/10 to-emerald-500/10 border-green-500/20",
    iconColor: "text-green-600 dark:text-green-400",
    badgeColor: "bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20",
    website: "https://edu.google.com",
  },
  {
    name: "Khan Academy",
    category: { uz: "Onlayn ta'lim", ru: "Онлайн-образование", en: "Online Education" },
    desc: {
      uz: "Bepul va sifatli onlayn ta'lim platformasi bilan integratsiya. Khan Academy resurslari orqali o'quvchilar bilimlarini mustahkamlaydi.",
      ru: "Интеграция с платформой бесплатного качественного онлайн-образования. Ресурсы Khan Academy помогают студентам закрепить знания.",
      en: "Integration with a free, quality online education platform. Khan Academy resources help students reinforce their knowledge.",
    },
    color: "from-violet-500/10 to-purple-500/10 border-violet-500/20",
    iconColor: "text-violet-600 dark:text-violet-400",
    badgeColor: "bg-violet-500/10 text-violet-600 dark:text-violet-400 border-violet-500/20",
    website: "https://khanacademy.org",
  },
  {
    name: "UNESCO",
    category: { uz: "Xalqaro tashkilot", ru: "Международная организация", en: "International Organization" },
    desc: {
      uz: "BMTning ta'lim, fan va madaniyat tashkiloti. UNESCO dasturlari doirasida xalqaro olimpiadalar o'tkazish va sertifikatlashtirish.",
      ru: "Организация ООН по образованию, науке и культуре. В рамках программ UNESCO проводятся международные олимпиады и сертификация.",
      en: "UN's organization for education, science and culture. International olympiads and certification under UNESCO programs.",
    },
    color: "from-sky-500/10 to-blue-500/10 border-sky-500/20",
    iconColor: "text-sky-600 dark:text-sky-400",
    badgeColor: "bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/20",
    website: "https://unesco.org",
  },
  {
    name: "Coursera",
    category: { uz: "Onlayn ta'lim", ru: "Онлайн-образование", en: "Online Education" },
    desc: {
      uz: "Dunyoning yetakchi onlayn ta'lim platformasi bilan hamkorlik. Coursera kurslari va sertifikatlari bilan o'quvchilarning imkoniyatlarini oshiramiz.",
      ru: "Партнёрство с ведущей мировой платформой онлайн-образования. Курсы и сертификаты Coursera расширяют возможности студентов.",
      en: "Partnership with the world's leading online education platform. Coursera courses and certificates expand student opportunities.",
    },
    color: "from-blue-500/10 to-indigo-500/10 border-blue-500/20",
    iconColor: "text-blue-600 dark:text-blue-400",
    badgeColor: "bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20",
    website: "https://coursera.org",
  },
  {
    name: "MIT OpenCourseWare",
    category: { uz: "Akademik hamkor", ru: "Академический партнёр", en: "Academic Partner" },
    desc: {
      uz: "MIT universitetining ochiq ta'lim resurslari. Dunyoning eng yaxshi universitetidan bepul materiallar orqali bilim olish imkoniyati.",
      ru: "Открытые образовательные ресурсы MIT. Возможность получить знания через бесплатные материалы от лучшего университета мира.",
      en: "MIT's open educational resources. Opportunity to gain knowledge through free materials from the world's best university.",
    },
    color: "from-red-500/10 to-orange-500/10 border-red-500/20",
    iconColor: "text-red-600 dark:text-red-400",
    badgeColor: "bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/20",
    website: "https://ocw.mit.edu",
  },
];

export default function PartnersPage() {
  const { lang } = useI18n();

  const title   = lang === "ru" ? "Наши партнёры" : lang === "en" ? "Our Partners" : "Hamkorlarimiz";
  const badge   = lang === "ru" ? "Партнёры"      : lang === "en" ? "Partners"     : "Hamkorlar";
  const desc    = lang === "ru"
    ? "Мы сотрудничаем с ведущими международными организациями для обеспечения высочайшего качества образования"
    : lang === "en"
    ? "We collaborate with leading international organizations to ensure the highest quality of education"
    : "Biz yetakchi xalqaro tashkilotlar bilan hamkorlik qilamiz — ta'limning eng yuqori sifatini ta'minlash uchun";

  const cta     = lang === "ru" ? "Стать партнёром" : lang === "en" ? "Become a Partner" : "Hamkor bo'lish";
  const ctaDesc = lang === "ru"
    ? "Хотите сотрудничать с нами? Свяжитесь с нами — мы открыты к новым партнёрствам."
    : lang === "en"
    ? "Want to collaborate with us? Get in touch — we are open to new partnerships."
    : "Biz bilan hamkorlik qilmoqchimisiz? Bog'laning — yangi hamkorlikka doimo tayyormiz.";

  const visitLabel = lang === "ru" ? "Сайт" : lang === "en" ? "Website" : "Sayt";

  return (
    <div className="min-h-screen bg-background">
      <div className="sticky top-0 z-50">
        <AnnouncementBar />
        <Header />
      </div>

      <main>
        {/* Page header */}
        <section className="py-16 border-b border-border bg-gradient-to-b from-primary/5 to-background">
          <div className="container mx-auto px-4 text-center">
            <div className="inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/10 px-4 py-1.5 text-sm text-primary mb-5">
              <Handshake className="h-4 w-4" />
              {badge}
            </div>
            <h1 className="text-3xl md:text-4xl lg:text-5xl font-bold text-foreground mb-4">
              {title}
            </h1>
            <p className="text-muted-foreground max-w-lg mx-auto text-base leading-relaxed">
              {desc}
            </p>
          </div>
        </section>

        {/* Partners grid */}
        <section className="py-16">
          <div className="container mx-auto px-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {partners.map((p) => (
                <div
                  key={p.name}
                  className={`group flex flex-col rounded-2xl border bg-gradient-to-br ${p.color} bg-card p-6 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg`}
                >
                  {/* Logo placeholder */}
                  <div className={`flex h-14 w-14 items-center justify-center rounded-xl border ${p.color} mb-4`}>
                    <Globe className={`h-7 w-7 ${p.iconColor}`} />
                  </div>

                  {/* Badge */}
                  <span className={`inline-flex items-center self-start rounded-full border px-2.5 py-0.5 text-xs font-medium mb-3 ${p.badgeColor}`}>
                    {p.category[lang] || p.category.uz}
                  </span>

                  <h3 className="font-bold text-foreground text-base mb-2">{p.name}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed flex-1">
                    {p.desc[lang] || p.desc.uz}
                  </p>

                  {/* Link */}
                  <a
                    href={p.website}
                    target="_blank"
                    rel="noopener noreferrer"
                    className={`mt-4 inline-flex items-center gap-1.5 text-sm font-medium ${p.iconColor} hover:opacity-80 transition-opacity`}
                  >
                    <ExternalLink className="h-3.5 w-3.5" />
                    {visitLabel}
                  </a>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Become a partner CTA */}
        <section className="py-16 border-t border-border">
          <div className="container mx-auto px-4">
            <div className="max-w-2xl mx-auto rounded-3xl border border-primary/20 bg-primary/5 p-10 text-center">
              <div className="inline-flex h-14 w-14 items-center justify-center rounded-2xl bg-primary/10 mb-6">
                <Handshake className="h-7 w-7 text-primary" />
              </div>
              <h2 className="text-2xl md:text-3xl font-bold text-foreground mb-4">{cta}</h2>
              <p className="text-muted-foreground mb-8 leading-relaxed">{ctaDesc}</p>
              <div className="flex flex-col sm:flex-row gap-3 justify-center">
                <a href="mailto:info@nextolymp.uz">
                  <Button className="btn-gradient gap-2">
                    <Mail className="h-4 w-4" />
                    info@nextolymp.uz
                  </Button>
                </a>
                <Link href="/news">
                  <Button variant="outline" className="gap-2 border-primary/30 text-primary hover:bg-primary/10">
                    {lang === "ru" ? "Узнать больше" : lang === "en" ? "Learn more" : "Ko'proq bilish"}
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
