"use client";

import { ShieldX } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export function AccessDenied() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-4">
      <div className="h-20 w-20 rounded-full bg-red-50 flex items-center justify-center mb-6">
        <ShieldX className="h-10 w-10 text-red-500" />
      </div>
      <h1 className="text-2xl font-bold text-gray-900 mb-2">Ruxsat berilmagan</h1>
      <p className="text-gray-500 mb-6 max-w-md">
        Sizda bu sahifaga kirish uchun ruxsat yo&apos;q. Administrator bilan bog&apos;laning yoki bosh sahifaga qayting.
      </p>
      <Link href="/admin">
        <Button variant="outline">Dashboard ga qaytish</Button>
      </Link>
    </div>
  );
}
