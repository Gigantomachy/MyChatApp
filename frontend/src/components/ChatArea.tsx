import React, { useState, useRef, useEffect } from 'react'
import { channels, users, messages as allMessages } from '../data/database'
import { useUser } from '../context/UserContext'
import type { Message } from '../data/database'
import './ChatArea.css'

interface ChatAreaProps {
  channelId: string
}

const ChatArea: React.FC<ChatAreaProps> = ({ channelId }) => {
  const currentUser = useUser()
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const channel = channels.find(c => c.id === channelId)

  useEffect(() => {
    const channelMessages = allMessages
      .filter(m => m.channelId === channelId)
      .sort((a, b) => a.timestamp - b.timestamp)
    setMessages(channelMessages)
  }, [channelId])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = () => {
    if (!input.trim()) return
    const newMsg: Message = {
      id: `local-${Date.now()}`,
      channelId,
      authorId: currentUser.user_id,
      content: input.trim(),
      timestamp: Date.now(),
    }
    setMessages(prev => [...prev, newMsg])
    setInput('')
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (ts: number) => {
    const d = new Date(ts)
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  if (!channel) {
    return (
      <div className="chat-area-empty">
        Select a conversation to start chatting
      </div>
    )
  }

  const getAuthorName = (authorId: string) => {
    if (authorId === currentUser.user_id) {
      return `${currentUser.first_name} ${currentUser.last_name}`
    }
    const u = users.find(user => user.id === authorId)
    return u ? `${u.firstName} ${u.lastName}` : 'Unknown'
  }

  const isCurrentUser = (authorId: string) => authorId === currentUser.user_id

  return (
    <div className="chat-area">
      <div className="chat-area-header">
        <span className="chat-area-header-icon">
          {channel.type === 'public' ? '#' : channel.type === 'group' ? '◆' : '@'}
        </span>
        <span className="chat-area-header-name">{channel.name}</span>
        <span className="chat-area-header-meta">
          {channel.memberIds.length} members
        </span>
      </div>

      <div className="chat-messages">
        {messages.map(msg => (
          <div
            key={msg.id}
            className={`message ${isCurrentUser(msg.authorId) ? 'message--self' : 'message--other'}`}
          >
            <div
              className={`message-avatar ${isCurrentUser(msg.authorId) ? 'message-avatar--self' : 'message-avatar--other'}`}
            >
              {getAuthorName(msg.authorId)
                .split(' ')
                .map(n => n[0])
                .join('')}
            </div>

            <div
              className={`message-bubble ${isCurrentUser(msg.authorId) ? 'message-bubble--self' : 'message-bubble--other'}`}
            >
              <div className="message-author">
                {getAuthorName(msg.authorId)}
                <span className="message-time">{formatTime(msg.timestamp)}</span>
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
          placeholder={`Message ${channel.type === 'public' ? '#' : ''}${channel.name}`}
          className="chat-input"
        />
        <button
          onClick={handleSend}
          disabled={!input.trim()}
          className="chat-send-button"
        >
          Send
        </button>
      </div>
    </div>
  )
}

export default ChatArea