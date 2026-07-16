import React, { createContext, useContext, useRef, useEffect, useCallback } from 'react'

type WSEventType = 'message.new' | 'friend_request.new' | 'friend_request.accepted' | 'channel.new'

interface WSEvent {
  type: WSEventType
  data: unknown
}

interface WSContextValue {
  on: (type: WSEventType, callback: (data: unknown) => void) => () => void
}

const WSContext = createContext<WSContextValue | null>(null)

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const wsRef = useRef<WebSocket | null>(null)
  const listenersRef = useRef<Map<string, Set<(data: unknown) => void>>>(new Map())
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const on = useCallback((type: WSEventType, callback: (data: unknown) => void) => {
    if (!listenersRef.current.has(type)) {
      listenersRef.current.set(type, new Set())
    }
    listenersRef.current.get(type)!.add(callback)
    return () => {
      listenersRef.current.get(type)?.delete(callback)
    }
  }, [])

  useEffect(() => {
    let reconnectDelay = 1000
    const maxDelay = 30000
    let closed = false

    const connect = () => {
      const ws = new WebSocket('/api/ws')
      wsRef.current = ws

      ws.onopen = () => {
        reconnectDelay = 1000
      }

      ws.onmessage = (e) => {
        try {
          const event: WSEvent = JSON.parse(e.data)
          const listeners = listenersRef.current.get(event.type)
          if (listeners) {
            listeners.forEach(cb => cb(event.data))
          }
        } catch {
          // ignore malformed messages
        }
      }

      ws.onclose = () => {
        if (closed) return
        reconnectTimeoutRef.current = setTimeout(() => {
          reconnectDelay = Math.min(reconnectDelay * 2, maxDelay)
          connect()
        }, reconnectDelay)
      }

      ws.onerror = () => {
        ws.close()
      }
    }

    connect()

    return () => {
      closed = true
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current)
      wsRef.current?.close()
    }
  }, [])

  return <WSContext.Provider value={{ on }}>{children}</WSContext.Provider>
}

export function useWebSocketEvent(
  type: WSEventType,
  callback: (data: unknown) => void
): void {
  const ctx = useContext(WSContext)
  const callbackRef = useRef(callback)
  callbackRef.current = callback

  useEffect(() => {
    if (!ctx) return
    return ctx.on(type, (data) => callbackRef.current(data))
  }, [ctx, type])
}
