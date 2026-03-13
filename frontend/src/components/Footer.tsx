"use client";

import { useI18n } from "@/lib/i18n";
import { Mail, Phone, MapPin, Trophy, Newspaper, BarChart3, Users, BookOpen } from "lucide-react";
import Link from "next/link";

export function Footer() {
  const { t, lang } = useI18n();

  const platformLinks = [
    { label: t("nav.olympiads"),     href: "/olympiads",  icon: Trophy },
    { label: t("nav.results"),       href: "/results",    icon: BarChart3 },
    { label: t("nav.news"),          href: "/news",       icon: Newspaper },
    { label: lang === "ru" ? "Mock testlar" : lang === "en" ? "Mock tests" : "Mock testlar", href: "/dashboard/mock-tests", icon: BookOpen },
  ];

  const infoLinks = [
    { label: t("nav.about"),    href: "/#about" },
    { label: t("nav.partners"), href: "/#partners" },
    { label: t("nav.login"),    href: "/login" },
    { label: t("nav.register"), href: "/register" },
  ];

  const footerDesc = lang === "ru"
    ? "Профессиональная платформа для проведения и участия в международных академических олимпиадах."
    : lang === "en"
    ? "A professional platform for organizing and participating in international academic olympiads."
    : "Xalqaro akademik olimpiadalarni tashkil etish va ularda ishtirok etish uchun professional platforma.";

  return (
    <footer className="bg-card border-t border-border text-muted-foreground">
      <div className="container mx-auto px-4 py-14">
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-8">
          {/* Brand */}
          <div>
            <Link href="/" className="flex items-center gap-2.5 mb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 text-white font-bold text-sm shadow-lg shadow-blue-500/25">
                NO
              </div>
              <span className="text-xl font-bold text-foreground">NextOly</span>
            </Link>
            <p className="text-sm leading-relaxed">{footerDesc}</p>
            {/* social placeholder */}
            <div className="flex gap-2 mt-4">
              <a
                href="https://t.me/nextolymp"
                target="_blank"
                rel="noopener noreferrer"
                className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary hover:bg-primary/20 transition-colors"
              >
                <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm4.917 7.014-1.797 8.466c-.134.604-.49.75-.99.467l-2.747-2.024-1.325 1.274c-.146.147-.269.27-.552.27l.197-2.793 5.088-4.598c.221-.197-.048-.307-.342-.11l-6.288 3.956-2.71-.846c-.588-.184-.6-.588.123-.87l10.577-4.077c.49-.178.92.11.766.885z" />
                </svg>
              </a>
              <a
                href="https://instagram.com/nextolymp"
                target="_blank"
                rel="noopener noreferrer"
                className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary hover:bg-primary/20 transition-colors"
              >
                <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/>
                </svg>
              </a>
            </div>
          </div>

          {/* Platform */}
          <div>
            <h3 className="text-foreground font-semibold mb-4">{t("footer.platform") || "Platforma"}</h3>
            <ul className="space-y-2.5 text-sm">
              {platformLinks.map((l) => (
                <li key={l.href}>
                  <Link
                    href={l.href}
                    className="flex items-center gap-2 hover:text-primary transition-colors"
                  >
                    <l.icon className="h-3.5 w-3.5" />
                    {l.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Info */}
          <div>
            <h3 className="text-foreground font-semibold mb-4">{t("footer.info") || "Ma'lumot"}</h3>
            <ul className="space-y-2.5 text-sm">
              {infoLinks.map((l) => (
                <li key={l.href}>
                  <Link href={l.href} className="hover:text-primary transition-colors">
                    {l.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Contact */}
          <div>
            <h3 className="text-foreground font-semibold mb-4">{t("footer.contact") || "Bog'lanish"}</h3>
            <ul className="space-y-3 text-sm">
              <li>
                <a href="mailto:info@nextolymp.uz" className="flex items-center gap-2.5 hover:text-primary transition-colors">
                  <Mail className="h-4 w-4 text-primary shrink-0" />
                  info@nextolymp.uz
                </a>
              </li>
              <li>
                <a href="tel:+998915627229" className="flex items-center gap-2.5 hover:text-primary transition-colors">
                  <Phone className="h-4 w-4 text-primary shrink-0" />
                  +998 91 562 7229
                </a>
              </li>
              <li className="flex items-center gap-2.5">
                <MapPin className="h-4 w-4 text-primary shrink-0" />
                Toshkent, O&apos;zbekiston
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="section-divider mt-10 mb-8" />
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 text-sm">
          <span>&copy; {new Date().getFullYear()} NextOly. {t("footer.rights") || "Barcha huquqlar himoyalangan."}</span>
          <div className="flex items-center gap-4">
            <Link href="/login" className="hover:text-primary transition-colors">{t("nav.login")}</Link>
            <Link href="/register" className="hover:text-primary transition-colors">{t("nav.register")}</Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
