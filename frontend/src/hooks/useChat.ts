"use client";

import { useState, useEffect, useRef, useCallback } from "react";

export interface ChatMessage {
  id: number;
  user_id: number;
  username: string;
  photo_url: string;
  content: string;
  type: string;
  role?: string;
  created_at: string;
}

export interface ChatStatus {
  is_banned: boolean;
  ban_reason: string;
  ban_type: string;
  ban_expires_at: string | null;
  is_chat_open: boolean;
  max_message_len: number;
  slow_mode: number;
  online_count: number;
}

interface UseChatConfig {
  apiUrl: string;
  token: string;
  enabled: boolean;
}

export function useChat(config: UseChatConfig) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [onlineCount, setOnlineCount] = useState(0);
  const [isChatOpen, setIsChatOpen] = useState(true);
  const [chatStatus, setChatStatus] = useState<ChatStatus | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [cooldown, setCooldown] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const cooldownTimerRef = useRef<NodeJS.Timeout | null>(null);

  // Load chat status
  const loadStatus = useCallback(async () => {
    try {
      const res = await fetch(`${config.apiUrl}/user/chat/status`, {
        headers: { Authorization: `Bearer ${config.token}` },
      });
      if (res.ok) {
        const data = await res.json();
        const status = data?.data ?? data;
        setChatStatus(status);
        setIsChatOpen(status.is_chat_open);
        setOnlineCount(status.online_count || 0);
      }
    } catch {
      // ignore
    }
  }, [config.apiUrl, config.token]);

  // Load message history
  const loadHistory = useCallback(async () => {
    try {
      const res = await fetch(`${config.apiUrl}/user/chat/messages?limit=50`, {
        headers: { Authorization: `Bearer ${config.token}` },
      });
      if (res.ok) {
        const data = await res.json();
        const msgs = data?.data?.messages || [];
        setMessages(msgs.reverse());
        setHasMore(msgs.length >= 50);
        if (data?.data?.meta?.online) {
          setOnlineCount(data.data.meta.online);
        }
      }
    } catch {
      // ignore
    }
  }, [config.apiUrl, config.token]);

  // Load older messages (cursor-based)
  const loadOlderMessages = useCallback(async () => {
    if (!hasMore || loadingMore || messages.length === 0) return;
    setLoadingMore(true);
    try {
      const oldestId = messages[0]?.id;
      const res = await fetch(
        `${config.apiUrl}/user/chat/messages?before_id=${oldestId}&limit=50`,
        { headers: { Authorization: `Bearer ${config.token}` } }
      );
      if (res.ok) {
        const data = await res.json();
        const older = data?.data?.messages || [];
        if (older.length === 0) {
          setHasMore(false);
        } else {
          setMessages((prev) => [...older.reverse(), ...prev]);
          setHasMore(data?.data?.meta?.has_more ?? older.length >= 50);
        }
      }
    } catch {
      // ignore
    } finally {
      setLoadingMore(false);
    }
  }, [config.apiUrl, config.token, hasMore, loadingMore, messages]);

  // WebSocket connect
  const connect = useCallback(() => {
    if (!config.enabled || !config.token) return;

    const wsUrl = config.apiUrl
      .replace("https://", "wss://")
      .replace("http://", "ws://");

    const ws = new WebSocket(`${wsUrl}/user/chat/ws?token=${config.token}`);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      reconnectAttemptsRef.current = 0;
      loadHistory();
      loadStatus();
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        switch (data.type) {
          case "new_message": {
            const newMsg = data.payload as ChatMessage;
            setMessages((prev) => {
              // Duplicate check — xabar allaqachon bormi
              if (prev.some((m) => m.id === newMsg.id)) return prev;
              return [...prev, newMsg];
            });
            break;
          }
          case "message_deleted":
            setMessages((prev) =>
              prev.filter((m) => m.id !== data.payload?.id)
            );
            break;
          case "online_count":
            setOnlineCount(data.payload?.count || 0);
            break;
          case "chat_status":
            setIsChatOpen(data.payload?.is_open ?? true);
            break;
          case "user_banned":
            loadStatus();
            break;
        }
      } catch {
        // parse error
      }
    };

    ws.onclose = () => {
      setIsConnected(false);
      if (config.enabled) {
        const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000);
        reconnectAttemptsRef.current += 1;
        reconnectTimeoutRef.current = setTimeout(connect, delay);
      }
    };

    ws.onerror = () => {
      ws.close();
    };
  }, [config.enabled, config.token, config.apiUrl, loadHistory, loadStatus]);

  // Send message with cooldown
  const sendMessage = useCallback(
    (content: string) => {
      if (wsRef.current?.readyState === WebSocket.OPEN && !cooldown) {
        wsRef.current.send(
          JSON.stringify({ type: "message", content })
        );
        // Set 3 second cooldown
        setCooldown(true);
        if (cooldownTimerRef.current) clearTimeout(cooldownTimerRef.current);
        cooldownTimerRef.current = setTimeout(() => {
          setCooldown(false);
        }, 3000);
      }
    },
    [cooldown]
  );

  // Connect on mount
  useEffect(() => {
    connect();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (cooldownTimerRef.current) {
        clearTimeout(cooldownTimerRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  return {
    messages,
    isConnected,
    onlineCount,
    isChatOpen,
    chatStatus,
    hasMore,
    loadingMore,
    cooldown,
    sendMessage,
    loadOlderMessages,
    loadStatus,
  };
}
