"use client";

import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Award, Calendar, Shield, FileCheck, Hash } from "lucide-react";
import { listCertificates, type Certificate } from "@/lib/user-api";

export default function CertificatesPage() {
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listCertificates()
      .then((data) => setCertificates(Array.isArray(data) ? data : []))
      .catch(() => setCertificates([]))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Yuklanmoqda...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Sertifikatlar</h1>
        <p className="text-muted-foreground mt-1">
          Olimpiada va mock testlardagi yutuqlaringiz uchun sertifikatlar
        </p>
      </div>

      {certificates.length === 0 ? (
        <Card className="border-0 shadow-sm">
          <CardContent className="p-12 text-center">
            <Award className="h-12 w-12 text-muted-foreground/50 mx-auto mb-4" />
            <p className="text-muted-foreground mb-2">Hozircha sertifikat yo&apos;q</p>
            <p className="text-sm text-muted-foreground">
              Olimpiada va mock testlarda qatnashib, sertifikat oling
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {certificates.map((cert) => (
            <Card
              key={cert.id}
              className="border-0 shadow-sm hover:shadow-md transition-shadow"
            >
              <CardContent className="p-6">
                {/* Header */}
                <div className="flex items-start justify-between mb-4">
                  <div className="h-12 w-12 rounded-xl bg-amber-50 flex items-center justify-center">
                    <Award className="h-6 w-6 text-amber-600" />
                  </div>
                  <Badge
                    className={`border-0 ${
                      cert.source_type === "olympiad"
                        ? "bg-purple-100 text-purple-700"
                        : "bg-green-100 text-green-700"
                    }`}
                  >
                    {cert.source_type === "olympiad" ? "Olimpiada" : "Mock test"}
                  </Badge>
                </div>

                {/* Title */}
                <h3 className="font-semibold text-foreground mb-1">{cert.title}</h3>
                {cert.source_title && (
                  <p className="text-sm text-muted-foreground mb-3">{cert.source_title}</p>
                )}

                {/* Score */}
                {cert.score > 0 && (
                  <div className="mb-3 p-2 bg-muted rounded-lg">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">Ball:</span>
                      <span className="font-semibold text-foreground">
                        {cert.score} ({cert.percentage}%)
                      </span>
                    </div>
                  </div>
                )}

                {/* Details */}
                <div className="space-y-2 text-sm">
                  <div className="flex items-center gap-2 text-muted-foreground">
                    <Hash className="h-3.5 w-3.5" />
                    <span className="font-mono text-xs">
                      {cert.certificate_number || cert.verification_code}
                    </span>
                  </div>
                  {cert.verification_code && cert.certificate_number && (
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <Shield className="h-3.5 w-3.5" />
                      <span className="font-mono text-xs">{cert.verification_code}</span>
                    </div>
                  )}
                  <div className="flex items-center gap-2 text-muted-foreground">
                    <Calendar className="h-3.5 w-3.5" />
                    {new Date(cert.issued_at || cert.created_at).toLocaleDateString("uz-UZ")}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
