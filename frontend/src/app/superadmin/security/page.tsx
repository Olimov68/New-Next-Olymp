"use client";

import { useEffect, useState } from "react";
import { getSecuritySettings, updateSecuritySettings } from "@/lib/superadmin-api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Shield, Save, Loader2 } from "lucide-react";

interface SecuritySettings {
  fullscreen_required: boolean;
  tab_switch_limit: number;
  blur_limit: number;
  offline_limit: number;
  heartbeat_interval_seconds: number;
  one_device_per_session: boolean;
  reconnect_policy: string;
  risk_threshold: number;
}

const defaultSettings: SecuritySettings = {
  fullscreen_required: true,
  tab_switch_limit: 3,
  blur_limit: 5,
  offline_limit: 3,
  heartbeat_interval_seconds: 30,
  one_device_per_session: true,
  reconnect_policy: "allow",
  risk_threshold: 70,
};

export default function SecurityPage() {
  const [settings, setSettings] = useState<SecuritySettings>(defaultSettings);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const fetchSettings = async () => {
    setLoading(true);
    try {
      const res = await getSecuritySettings();
      setSettings(res.data || res);
    } catch {
      setError("Sozlamalarni yuklashda xatolik");
    }
    setLoading(false);
  };

  useEffect(() => { fetchSettings(); }, []);

  const handleSave = async () => {
    setSaving(true);
    setError("");
    setSuccess("");
    try {
      await updateSecuritySettings({
        ...settings,
        tab_switch_limit: Number(settings.tab_switch_limit),
        blur_limit: Number(settings.blur_limit),
        offline_limit: Number(settings.offline_limit),
        heartbeat_interval_seconds: Number(settings.heartbeat_interval_seconds),
        risk_threshold: Number(settings.risk_threshold),
      });
      setSuccess("Sozlamalar saqlandi");
      setTimeout(() => setSuccess(""), 3000);
    } catch (e: unknown) {
      setError((e as { response?: { data?: { message?: string } } })?.response?.data?.message || "Xatolik");
    }
    setSaving(false);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-2xl">
      <div className="flex items-center gap-3">
        <Shield className="w-6 h-6 text-orange-500" />
        <h1 className="text-2xl font-bold">Xavfsizlik sozlamalari</h1>
      </div>

      {error && <div className="bg-red-900/50 border border-red-800 text-red-300 rounded-lg p-3 text-sm">{error}</div>}
      {success && <div className="bg-green-900/50 border border-green-800 text-green-300 rounded-lg p-3 text-sm">{success}</div>}

      <div className="bg-muted border border-border rounded-lg p-6 space-y-6">
        {/* Fullscreen Required */}
        <div className="flex items-center justify-between">
          <div>
            <Label className="text-sm font-medium">To&apos;liq ekran talab qilinadi</Label>
            <p className="text-xs text-muted-foreground mt-1">Test vaqtida to&apos;liq ekran rejimi majburiy bo&apos;ladi</p>
          </div>
          <label className="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" checked={settings.fullscreen_required} onChange={(e) => setSettings({ ...settings, fullscreen_required: e.target.checked })}
              className="sr-only peer" />
            <div className="w-11 h-6 bg-gray-600 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-orange-500"></div>
          </label>
        </div>

        {/* One Device Per Session */}
        <div className="flex items-center justify-between">
          <div>
            <Label className="text-sm font-medium">Bir qurilma - bir sessiya</Label>
            <p className="text-xs text-muted-foreground mt-1">Foydalanuvchi faqat bitta qurilmadan foydalana oladi</p>
          </div>
          <label className="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" checked={settings.one_device_per_session} onChange={(e) => setSettings({ ...settings, one_device_per_session: e.target.checked })}
              className="sr-only peer" />
            <div className="w-11 h-6 bg-gray-600 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-orange-500"></div>
          </label>
        </div>

        {/* Tab Switch Limit */}
        <div>
          <Label className="text-sm font-medium">Tab almashtirish limiti</Label>
          <p className="text-xs text-muted-foreground mt-1 mb-2">Ruxsat etilgan tab almashtirishlar soni</p>
          <Input type="number" value={settings.tab_switch_limit} onChange={(e) => setSettings({ ...settings, tab_switch_limit: Number(e.target.value) })}
            className="bg-card border-border w-32" min={0} />
        </div>

        {/* Blur Limit */}
        <div>
          <Label className="text-sm font-medium">Blur limiti</Label>
          <p className="text-xs text-muted-foreground mt-1 mb-2">Oynadan chiqish (blur) limiti</p>
          <Input type="number" value={settings.blur_limit} onChange={(e) => setSettings({ ...settings, blur_limit: Number(e.target.value) })}
            className="bg-card border-border w-32" min={0} />
        </div>

        {/* Offline Limit */}
        <div>
          <Label className="text-sm font-medium">Oflayn limiti</Label>
          <p className="text-xs text-muted-foreground mt-1 mb-2">Internet uzilishi limiti</p>
          <Input type="number" value={settings.offline_limit} onChange={(e) => setSettings({ ...settings, offline_limit: Number(e.target.value) })}
            className="bg-card border-border w-32" min={0} />
        </div>

        {/* Heartbeat Interval */}
        <div>
          <Label className="text-sm font-medium">Heartbeat intervali (soniya)</Label>
          <p className="text-xs text-muted-foreground mt-1 mb-2">Server bilan aloqa tekshirish oraligi</p>
          <Input type="number" value={settings.heartbeat_interval_seconds} onChange={(e) => setSettings({ ...settings, heartbeat_interval_seconds: Number(e.target.value) })}
            className="bg-card border-border w-32" min={5} />
        </div>

        {/* Reconnect Policy */}
        <div>
          <Label className="text-sm font-medium">Qayta ulanish siyosati</Label>
          <p className="text-xs text-muted-foreground mt-1 mb-2">Internet uzilganda qayta ulanish imkoniyati</p>
          <Select value={settings.reconnect_policy} onValueChange={(v) => setSettings({ ...settings, reconnect_policy: v ?? "" })}>
            <SelectTrigger className="bg-card border-border w-48"><SelectValue /></SelectTrigger>
            <SelectContent>
              <SelectItem value="allow">Ruxsat berilgan</SelectItem>
              <SelectItem value="deny">Taqiqlangan</SelectItem>
              <SelectItem value="once">Bir marta</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Risk Threshold */}
        <div>
          <Label className="text-sm font-medium">Risk chegarasi</Label>
          <p className="text-xs text-muted-foreground mt-1 mb-2">Xavf darajasi chegarasi (0-100)</p>
          <Input type="number" value={settings.risk_threshold} onChange={(e) => setSettings({ ...settings, risk_threshold: Number(e.target.value) })}
            className="bg-card border-border w-32" min={0} max={100} />
        </div>
      </div>

      <Button onClick={handleSave} disabled={saving} className="bg-orange-500 hover:bg-orange-600">
        {saving ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <Save className="w-4 h-4 mr-2" />}
        Saqlash
      </Button>
    </div>
  );
}
