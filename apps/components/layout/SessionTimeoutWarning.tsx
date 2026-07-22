'use client';

import * as React from 'react';
import { useState, useEffect } from 'react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { useSessionTimeoutWarning } from '@/hooks/use-session-persistence';
import { useAuth } from '@/context/AuthContext';

export interface SessionTimeoutWarningProps {
  warningThreshold?: number; // ms before warning (default: 5 minutes)
  criticalThreshold?: number; // ms before critical (default: 2 minutes)
  showWarning?: boolean;
  showCritical?: boolean;
}

export function SessionTimeoutWarning({
  warningThreshold = 5 * 60 * 1000,
  criticalThreshold = 2 * 60 * 1000,
  showWarning = true,
  showCritical = true,
}: SessionTimeoutWarningProps) {
  const [showDialog, setShowDialog] = useState(false);
  const [isCritical, setIsCritical] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const { refresh, logout } = useAuth();

  const warningType = useSessionTimeoutWarning(
    warningThreshold,
    criticalThreshold,
    () => {
      if (showWarning) {
        setIsCritical(false);
        setShowDialog(true);
      }
    },
    () => {
      if (showCritical) {
        setIsCritical(true);
        setShowDialog(true);
      }
    },
    () => {
      // Session expired
      setShowDialog(false);
    }
  );

  // Calculate countdown
  useEffect(() => {
    if (!showDialog) {
      setCountdown(0);
      return;
    }

    const interval = setInterval(() => {
      const remaining = isCritical ? criticalThreshold : warningThreshold;
      const elapsed = Date.now() % remaining;
      const timeLeft = Math.ceil((remaining - elapsed) / 1000);
      setCountdown(Math.max(0, timeLeft));
    }, 1000);

    return () => clearInterval(interval);
  }, [showDialog, isCritical, warningThreshold, criticalThreshold]);

  const handleExtendSession = async () => {
    try {
      await refresh();
      setShowDialog(false);
      setIsCritical(false);
    } catch {
      // Refresh failed, logout
      await logout();
    }
  };

  const handleLogout = async () => {
    await logout();
    setShowDialog(false);
  };

  if (!showDialog) return null;

  return (
    <Alert
      variant={isCritical ? 'destructive' : 'default'}
      className="fixed bottom-4 right-4 max-w-sm shadow-lg z-50"
    >
      <AlertTitle className="text-sm">
        {isCritical ? 'Session expiring soon!' : 'Session warning'}
      </AlertTitle>
      <AlertDescription className="text-sm">
        {isCritical
          ? `Your session will expire in ${countdown} seconds. `
          : `Your session will expire in ${countdown} seconds. `}
        {isCritical ? 'Save your work!' : 'Would you like to extend your session?'}
      </AlertDescription>
      <div className="mt-3 flex gap-2">
        <Button
          variant={isCritical ? 'default' : 'outline'}
          size="sm"
          onClick={handleExtendSession}
        >
          Extend Session
        </Button>
        <Button
          variant={isCritical ? 'outline' : 'ghost'}
          size="sm"
          onClick={handleLogout}
        >
          Logout
        </Button>
      </div>
    </Alert>
  );
}

export default SessionTimeoutWarning;
