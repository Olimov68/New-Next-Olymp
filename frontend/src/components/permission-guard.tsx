"use client";

import { usePanelAuth } from "@/lib/panel-auth-context";
import { AccessDenied } from "@/components/access-denied";
import type { ReactNode } from "react";

interface PermissionGuardProps {
  permission?: string;
  permissions?: string[];
  module?: string;
  fallback?: ReactNode;
  showAccessDenied?: boolean;
  children: ReactNode;
}

export function PermissionGuard({
  permission,
  permissions,
  module,
  fallback = null,
  showAccessDenied = false,
  children,
}: PermissionGuardProps) {
  const { hasPermission, hasAnyPermission, hasModuleAccess, loading } = usePanelAuth();

  if (loading) return null;

  let allowed = false;

  if (permission) {
    allowed = hasPermission(permission);
  } else if (permissions && permissions.length > 0) {
    allowed = hasAnyPermission(permissions);
  } else if (module) {
    allowed = hasModuleAccess(module);
  } else {
    allowed = true;
  }

  if (!allowed) {
    if (showAccessDenied) return <AccessDenied />;
    return <>{fallback}</>;
  }

  return <>{children}</>;
}
