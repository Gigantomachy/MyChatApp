import React, { useState, useRef, useEffect, useCallback } from 'react'
import { useUser } from '../context/UserContext'
import {
  getMessages,
  sendMessage as apiSendMessage,
  type ChannelMembership,
  type MessageItem,
} from '../api/client'
import './ChatArea.css'

interface ChatAreaProps {
  channel: ChannelMembership | null
}

const ChatArea: React.FC<ChatAreaProps> = ({ channel }) => {
  const currentUser = useUser()
  const [messages, setMessages] = useState<MessageItem[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [sending, setSending] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const fetchMessages = useCallback(async (channelId: string) => {
    setLoading(true)
    setError('')
    try {
      const msgs = await getMessages(channelId)
      setMessages((msgs ?? []).reverse())
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load messages')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (channel) {
      fetchMessages(channel.channel_id)
    } else {
      setMessages([])
    }
  }, [channel, fetchMessages])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = async () => {
    if (!input.trim() || !channel || sending) return
    setSending(true)
    try {
      const msg = await apiSendMessage(channel.channel_id, input.trim())
      setMessages(prev => [...prev, msg])
      setInput('')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to send message')
    } finally {
      setSending(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (iso: string) => {
    const d = new Date(iso)
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  const getAuthorName = (msg: MessageItem) => {
    return `${msg.author_first_name} ${msg.author_last_name}`
  }

  const getInitials = (msg: MessageItem) => {
    return `${msg.author_first_name[0] ?? ''}${msg.author_last_name[0] ?? ''}`
  }

  const isCurrentUser = (authorId: string) => authorId === currentUser.user_id

  if (!channel) {
    return (
      <div className="chat-area-empty">
        Select a conversation to start chatting
      </div>
    )
  }

  return (
    <div className="chat-area">
      <div className="chat-area-header">
        <span className="chat-area-header-icon">
          {channel.channel_type === 'public' ? '#' : '@'}
        </span>
        <span className="chat-area-header-name">{channel.channel_name}</span>
      </div>

      <div className="chat-messages">
        {loading && (
          <div className="chat-messages-info">Loading messages...</div>
        )}
        {error && (
          <div className="chat-messages-error">{error}</div>
        )}
        {!loading && !error && messages.length === 0 && (
          <div className="chat-messages-info">No messages yet. Say hello!</div>
        )}
        {!loading && messages.map(msg => (
          <div
            key={msg.message_id}
            className={`message ${isCurrentUser(msg.author_id) ? 'message--self' : 'message--other'}`}
          >
            <div
              className={`message-avatar ${isCurrentUser(msg.author_id) ? 'message-avatar--self' : 'message-avatar--other'}`}
            >
              {getInitials(msg)}
            </div>

            <div
              className={`message-bubble ${isCurrentUser(msg.author_id) ? 'message-bubble--self' : 'message-bubble--other'}`}
            >
              <div className="message-author">
                {getAuthorName(msg)}
                <span className="message-time">{formatTime(msg.created_at)}</span>
              </div>
              <div>{msg.content}</div>
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <div className="chat-input-area">
        <input
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={`Message ${channel.channel_type === 'public' ? '#' : ''}${channel.channel_name}`}
          className="chat-input"
          disabled={sending}
        />
        <button
          onClick={handleSend}
          disabled={!input.trim() || sending}
          className="chat-send-button"
        >
          {sending ? '...' : 'Send'}
        </button>
      </div>
    </div>
  )
}

export default ChatArea
