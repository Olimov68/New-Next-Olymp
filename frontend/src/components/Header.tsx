"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, X, LogIn, UserPlus, LogOut, LayoutDashboard, Trophy, Newspaper, BarChart3, Home, Handshake } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger, SheetTitle } from "@/components/ui/sheet";
import { useAuth } from "@/lib/auth-context";
import { useI18n, type Lang } from "@/lib/i18n";
import { ThemeToggle } from "@/components/ThemeToggle";

const languages: { code: Lang; label: string }[] = [
  { code: "uz", label: "UZ" },
  { code: "ru", label: "RU" },
  { code: "en", label: "EN" },
];

export function Header() {
  const [open, setOpen] = useState(false);
  const { user, logout } = useAuth();
  const { t, lang, setLang } = useI18n();
  const pathname = usePathname();

  const menuItems = [
    { label: t("nav.home") || "Asosiy",        href: "/",          icon: Home },
    { label: t("nav.olympiads"),                href: "/olympiads", icon: Trophy },
    { label: t("nav.news"),                     href: "/news",      icon: Newspaper },
    { label: t("nav.results"),                  href: "/results",   icon: BarChart3 },
    { label: t("nav.partners") || "Hamkorlar",  href: "/partners",  icon: Handshake },
  ];

  const isActive = (href: string) => {
    if (href === "/") return pathname === "/";
    if (href.startsWith("/#")) return pathname === "/";
    return pathname.startsWith(href);
  };

  return (
    <header className="w-full border-b border-border bg-background/80 backdrop-blur-xl">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2 shrink-0">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 text-white font-bold text-sm shadow-lg shadow-blue-500/25">
            N
          </div>
          <span className="text-xl font-bold text-foreground">NextOly</span>
        </Link>

        {/* Desktop Nav */}
        <nav className="hidden lg:flex items-center gap-0.5">
          {menuItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-lg transition-colors ${
                isActive(item.href)
                  ? "text-primary bg-primary/10"
                  : "text-muted-foreground hover:text-foreground hover:bg-accent"
              }`}
            >
              <item.icon className="h-4 w-4" />
              {item.label}
            </Link>
          ))}
        </nav>

        {/* Right side controls */}
        <div className="flex items-center gap-2">
          {/* Language Selector */}
          <div className="hidden sm:flex items-center rounded-lg border border-border bg-muted/50 overflow-hidden">
            {languages.map((l) => (
              <button
                key={l.code}
                onClick={() => setLang(l.code)}
                className={`px-2.5 py-1.5 text-xs font-medium transition-colors ${
                  lang === l.code
                    ? "bg-primary text-primary-foreground shadow-inner"
                    : "text-muted-foreground hover:text-foreground hover:bg-accent"
                }`}
              >
                {l.label}
              </button>
            ))}
          </div>

          {/* Theme Toggle */}
          <ThemeToggle variant="compact" />

          {/* Auth buttons - desktop only */}
          {user ? (
            <div className="hidden sm:flex items-center gap-2">
              <Link href="/dashboard">
                <Button variant="outline" size="sm" className="gap-1.5">
                  <LayoutDashboard className="h-4 w-4" />
                  {t("nav.dashboard")}
                </Button>
              </Link>
              <Button
                variant="ghost"
                size="sm"
                onClick={logout}
                className="text-muted-foreground hover:text-foreground"
              >
                <LogOut className="h-4 w-4" />
              </Button>
            </div>
          ) : (
            <div className="hidden sm:flex items-center gap-2">
              <Link href="/login">
                <Button variant="ghost" size="sm" className="gap-1.5 text-muted-foreground hover:text-foreground">
                  <LogIn className="h-4 w-4" />
                  {t("nav.login")}
                </Button>
              </Link>
              <Link href="/register">
                <Button
                  size="sm"
                  className="gap-1.5 bg-gradient-to-r from-blue-500 to-indigo-600 text-white hover:from-blue-600 hover:to-indigo-700 shadow-lg shadow-blue-500/25 border-0"
                >
                  <UserPlus className="h-4 w-4" />
                  {t("nav.register")}
                </Button>
              </Link>
            </div>
          )}

          {/* Mobile hamburger */}
          <Sheet open={open} onOpenChange={setOpen}>
            <SheetTrigger className="ml-1 inline-flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground hover:bg-accent hover:text-foreground transition-colors lg:hidden">
              {open ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </SheetTrigger>
            <SheetContent side="right" className="w-72 bg-background border-l border-border">
              <SheetTitle className="sr-only">Menu</SheetTitle>
              <nav className="flex flex-col gap-1 mt-8">
                {/* Dashboard link for logged-in users */}
                {user && (
                  <Link
                    href="/dashboard"
                    onClick={() => setOpen(false)}
                    className="flex items-center gap-2 rounded-lg px-4 py-3 text-sm font-medium text-primary bg-primary/10 border border-primary/20 mb-2"
                  >
                    <LayoutDashboard className="h-4 w-4" />
                    {t("nav.dashboard")} (@{user.username})
                  </Link>
                )}

                {/* Navigation items - proper routes */}
                {menuItems.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    onClick={() => setOpen(false)}
                    className={`flex items-center gap-2 rounded-lg px-4 py-3 text-sm font-medium transition-colors ${
                      isActive(item.href)
                        ? "text-primary bg-primary/10"
                        : "text-muted-foreground hover:bg-accent hover:text-foreground"
                    }`}
                  >
                    <item.icon className="h-4 w-4" />
                    {item.label}
                  </Link>
                ))}

                {/* Mobile auth buttons */}
                {!user && (
                  <div className="flex flex-col gap-2 mt-4 px-1">
                    <Link href="/login" onClick={() => setOpen(false)}>
                      <Button variant="outline" size="sm" className="w-full gap-1.5">
                        <LogIn className="h-4 w-4" />
                        {t("nav.login")}
                      </Button>
                    </Link>
                    <Link href="/register" onClick={() => setOpen(false)}>
                      <Button
                        size="sm"
                        className="w-full gap-1.5 bg-gradient-to-r from-blue-500 to-indigo-600 text-white border-0"
                      >
                        <UserPlus className="h-4 w-4" />
                        {t("nav.register")}
                      </Button>
                    </Link>
                  </div>
                )}

                {/* Language selector */}
                <div className="flex items-center gap-1 mt-4 px-4">
                  {languages.map((l) => (
                    <button
                      key={l.code}
                      onClick={() => setLang(l.code)}
                      className={`px-3 py-1.5 text-xs font-medium rounded-md transition-colors ${
                        lang === l.code
                          ? "bg-primary text-primary-foreground"
                          : "text-muted-foreground hover:text-foreground hover:bg-accent"
                      }`}
                    >
                      {l.label}
                    </button>
                  ))}
                </div>

                {/* Mobile logout */}
                {user && (
                  <div className="mt-4 px-4">
                    <button
                      onClick={() => { logout(); setOpen(false); }}
                      className="w-full flex items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium text-destructive hover:bg-destructive/10 transition-colors"
                    >
                      <LogOut className="h-4 w-4" />
                      {t("nav.logout")}
                    </button>
                  </div>
                )}
              </nav>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </header>
  );
}
