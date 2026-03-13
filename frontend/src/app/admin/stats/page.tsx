"use client";

import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchStats, updateStats } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Save } from "lucide-react";

export default function AdminStatsPage() {
  const queryClient = useQueryClient();
  const { data: stats } = useQuery({ queryKey: ["stats"], queryFn: fetchStats });
  const [form, setForm] = useState({ total_users: 0, total_olympiads: 0, total_regions: 0, total_mock_tests: 0 });

  useEffect(() => {
    if (stats) {
      setForm({
        total_users: stats.total_users,
        total_olympiads: stats.total_olympiads,
        total_regions: stats.total_regions,
        total_mock_tests: stats.total_mock_tests,
      });
    }
  }, [stats]);

  const updateMut = useMutation({
    mutationFn: updateStats,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["stats"] }),
  });

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    updateMut.mutate(form);
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-foreground mb-6">Statistika</h1>

      <Card className="border-0 shadow-sm max-w-lg">
        <CardHeader>
          <CardTitle className="text-lg">Platformaning umumiy statistikasi</CardTitle>
          <p className="text-sm text-muted-foreground">Bu raqamlar bosh sahifada ko&apos;rsatiladi</p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label>Ishtirokchilar soni</Label>
              <Input type="number" value={form.total_users} onChange={(e) => setForm({ ...form, total_users: Number(e.target.value) })} />
            </div>
            <div className="space-y-2">
              <Label>Olimpiadalar soni</Label>
              <Input type="number" value={form.total_olympiads} onChange={(e) => setForm({ ...form, total_olympiads: Number(e.target.value) })} />
            </div>
            <div className="space-y-2">
              <Label>Viloyatlar soni</Label>
              <Input type="number" value={form.total_regions} onChange={(e) => setForm({ ...form, total_regions: Number(e.target.value) })} />
            </div>
            <div className="space-y-2">
              <Label>Mock testlar soni</Label>
              <Input type="number" value={form.total_mock_tests} onChange={(e) => setForm({ ...form, total_mock_tests: Number(e.target.value) })} />
            </div>
            <Button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 gap-2" disabled={updateMut.isPending}>
              <Save className="h-4 w-4" />
              {updateMut.isPending ? "Saqlanmoqda..." : "Saqlash"}
            </Button>
            {updateMut.isSuccess && (
              <div className="rounded-lg bg-green-50 text-green-600 text-sm p-3 text-center">Muvaffaqiyatli saqlandi!</div>
            )}
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
