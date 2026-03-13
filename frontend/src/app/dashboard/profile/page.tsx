"use client";

import { useState, useEffect, useRef } from "react";
import { useMutation } from "@tanstack/react-query";
import { completeProfile, updateProfile, uploadPhoto, getMe } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import { useI18n } from "@/lib/i18n";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Save, UserCircle, Camera, Loader2 } from "lucide-react";

const regions = [
  "Toshkent shahri", "Toshkent viloyati", "Samarqand", "Buxoro",
  "Farg'ona", "Andijon", "Namangan", "Xorazm", "Qashqadaryo",
  "Surxondaryo", "Navoiy", "Jizzax", "Sirdaryo", "Qoraqalpog'iston",
];

const grades = ["5", "6", "7", "8", "9", "10", "11"];

export default function ProfilePage() {
  const { user, refreshUser } = useAuth();
  const { t } = useI18n();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [profileCompleted, setProfileCompleted] = useState(false);
  const [photoUrl, setPhotoUrl] = useState("");
  const [selectedPhoto, setSelectedPhoto] = useState<File | null>(null);
  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    birth_date: "",
    gender: "",
    region: "",
    district: "",
    school_name: "",
    grade: "",
  });
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (user) {
      getMe().then((data) => {
        setProfileCompleted(data.user.is_profile_completed);
        if (data.profile) {
          setPhotoUrl(data.profile.photo_url || "");
          setForm({
            first_name: data.profile.first_name || "",
            last_name: data.profile.last_name || "",
            birth_date: data.profile.birth_date || "",
            gender: data.profile.gender || "",
            region: data.profile.region || "",
            district: data.profile.district || "",
            school_name: data.profile.school_name || "",
            grade: data.profile.grade ? String(data.profile.grade) : "",
          });
        }
        setLoading(false);
      }).catch(() => setLoading(false));
    }
  }, [user]);

  const isNewProfile = !profileCompleted;

  const mutation = useMutation({
    mutationFn: async (formData: typeof form) => {
      if (isNewProfile) {
        const fd = new FormData();
        Object.entries(formData).forEach(([key, value]) => {
          if (key === "grade") {
            fd.append(key, String(Number(value)));
          } else {
            fd.append(key, value);
          }
        });
        if (selectedPhoto) {
          fd.append("photo", selectedPhoto);
        }
        return completeProfile(fd);
      } else {
        return updateProfile({
          ...formData,
          grade: Number(formData.grade),
        });
      }
    },
    onSuccess: async () => {
      setProfileCompleted(true);
      await refreshUser();
      setSuccess(true);
      setError("");
      setTimeout(() => setSuccess(false), 3000);
    },
    onError: (err: unknown) => {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message;
      setError(msg || "Xatolik yuz berdi");
    },
  });

  const photoMutation = useMutation({
    mutationFn: (file: File) => uploadPhoto(file),
    onSuccess: (data) => {
      setPhotoUrl(data.photo_url);
      setError("");
    },
    onError: (err: unknown) => {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message;
      setError(msg || "Rasm yuklashda xatolik");
    },
  });

  const handlePhotoChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    if (isNewProfile) {
      setSelectedPhoto(file);
      setPhotoUrl(URL.createObjectURL(file));
      setError("");
    } else {
      photoMutation.mutate(file);
    }
  };

  const update = (field: string, value: string) => {
    setForm((f) => ({ ...f, [field]: value }));
    setSuccess(false);
    setError("");
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isNewProfile && !selectedPhoto) {
      setError("Iltimos, avval rasmingizni yuklang");
      return;
    }
    mutation.mutate(form);
  };

  if (!user) return null;
  if (loading) return <div className="text-muted-foreground p-8 text-center">{t("common.loading")}</div>;

  const apiBase = process.env.NEXT_PUBLIC_API_URL?.replace("/api/v1", "") || "http://localhost:8080";

  return (
    <div className="max-w-lg mx-auto">
      <h1 className="text-2xl font-bold text-foreground mb-6">{t("profile.title")}</h1>

      <Card className="border-0 shadow-sm">
        <CardHeader>
          <div className="flex items-center gap-4">
            {/* Photo */}
            <div className="relative group">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary/10 text-primary overflow-hidden">
                {photoUrl ? (
                  <img
                    src={photoUrl.startsWith("blob:") ? photoUrl : `${apiBase}${photoUrl}`}
                    alt="Profile"
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <UserCircle className="h-8 w-8" />
                )}
              </div>
              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                className={`absolute inset-0 flex items-center justify-center rounded-full bg-black/40 transition-opacity ${
                  isNewProfile && !selectedPhoto ? "opacity-100" : "opacity-0 group-hover:opacity-100"
                }`}
                disabled={photoMutation.isPending}
              >
                {photoMutation.isPending ? (
                  <Loader2 className="h-5 w-5 text-white animate-spin" />
                ) : (
                  <Camera className="h-5 w-5 text-white" />
                )}
              </button>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/jpeg,image/png,image/webp"
                className="hidden"
                onChange={handlePhotoChange}
              />
            </div>
            <div>
              <CardTitle className="text-lg">@{user.username}</CardTitle>
              <p className="text-sm text-muted-foreground">
                {isNewProfile ? "2-qadam: Profilni to'ldiring" : t("profile.title")}
              </p>
              {isNewProfile && !selectedPhoto && (
                <p className="text-xs text-red-500 mt-1">Rasmingizni yuklang *</p>
              )}
              {isNewProfile && selectedPhoto && (
                <p className="text-xs text-green-500 mt-1">Rasm tanlandi</p>
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="rounded-lg bg-destructive/10 text-destructive text-sm p-3 text-center">{error}</div>
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>{t("auth.firstName") || "Ism"} *</Label>
                <Input value={form.first_name} onChange={(e) => update("first_name", e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label>{t("auth.lastName") || "Familiya"} *</Label>
                <Input value={form.last_name} onChange={(e) => update("last_name", e.target.value)} required />
              </div>
            </div>

            <div className="space-y-2">
              <Label>{t("auth.username")}</Label>
              <Input value={user.username} disabled className="bg-muted" />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>{"Tug'ilgan sana"} *</Label>
                <Input type="date" value={form.birth_date} onChange={(e) => update("birth_date", e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label>{"Jinsi"} *</Label>
                <Select value={form.gender} onValueChange={(v) => update("gender", v ?? "")}>
                  <SelectTrigger><SelectValue placeholder="Tanlang" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="male">Erkak</SelectItem>
                    <SelectItem value="female">Ayol</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-2">
              <Label>{t("auth.region") || "Viloyat"} *</Label>
              <Select value={form.region} onValueChange={(v) => update("region", v ?? "")}>
                <SelectTrigger><SelectValue placeholder="Tanlang" /></SelectTrigger>
                <SelectContent>
                  {regions.map((r) => (
                    <SelectItem key={r} value={r}>{r}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>{t("auth.district") || "Tuman"} *</Label>
                <Input value={form.district} onChange={(e) => update("district", e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label>{"Maktab nomi"} *</Label>
                <Input value={form.school_name} onChange={(e) => update("school_name", e.target.value)} required />
              </div>
            </div>

            <div className="space-y-2">
              <Label>{t("auth.grade") || "Sinf"} *</Label>
              <Select value={form.grade} onValueChange={(v) => update("grade", v ?? "")}>
                <SelectTrigger><SelectValue placeholder="Tanlang" /></SelectTrigger>
                <SelectContent>
                  {grades.map((g) => (
                    <SelectItem key={g} value={g}>{g}-sinf</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 gap-2" disabled={mutation.isPending}>
              {mutation.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Save className="h-4 w-4" />
              )}
              {mutation.isPending ? t("common.loading") : (isNewProfile ? "Profilni yaratish" : (t("profile.save") || "Saqlash"))}
            </Button>

            {success && (
              <div className="rounded-lg bg-green-50 text-green-600 text-sm p-3 text-center">
                {t("profile.saved") || "Saqlandi!"}
              </div>
            )}
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
