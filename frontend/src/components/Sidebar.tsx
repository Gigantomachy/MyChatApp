import React from 'react'
import { channels, users } from '../data/database'
import type { BackendUser } from '../api/client'
import './Sidebar.css'

interface SidebarProps {
  currentUser: BackendUser
  selectedChannelId: string | null
  onSelectChannel: (channelId: string) => void
  onStartNewChat: () => void
  onOpenSearch: () => void
}

const Sidebar: React.FC<SidebarProps> = ({
  currentUser,
  selectedChannelId,
  onSelectChannel,
  onStartNewChat,
  onOpenSearch,
}) => {
  const userId = currentUser.user_id

  const publicChannels = channels.filter(
    ch => ch.type === 'public' && ch.memberIds.includes(userId)
  )

  const dmChannels = channels.filter(
    ch => (ch.type === 'dm' || ch.type === 'group') && ch.memberIds.includes(userId)
  )

  const getDmLabel = (ch: typeof channels[0]) => {
    const otherIds = ch.memberIds.filter(id => id !== userId)
    if (otherIds.length === 0) return ch.name
    const names = otherIds.map(id => {
      const u = users.find(user => user.id === id)
      return u ? `${u.firstName} ${u.lastName}` : null
    }).filter(Boolean) as string[]
    return names.length > 0 ? names.join(', ') : ch.name
  }

  const renderChannelItem = (ch: typeof channels[0], label: string) => {
    const isActive = selectedChannelId === ch.id
    return (
      <div
        key={ch.id}
        onClick={() => onSelectChannel(ch.id)}
        className={`sidebar-channel ${isActive ? 'sidebar-channel--active' : ''}`}
      >
        <span className="sidebar-channel-icon">
          {ch.type === 'public' ? '#' : '@'}
        </span>
        <span className="sidebar-channel-name">{label}</span>
      </div>
    )
  }

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <div className="sidebar-avatar">
          {currentUser.first_name[0]}{currentUser.last_name[0]}
        </div>
        <div className="sidebar-user-info">
          <div className="sidebar-user-name">
            {currentUser.first_name} {currentUser.last_name}
          </div>
          <div className="sidebar-user-handle">@{currentUser.username}</div>
        </div>
      </div>

      <button
        className="sidebar-search-bar"
        onClick={onOpenSearch}
        aria-label="Search channels and people"
        title="Search channels and people"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <span>Search</span>
      </button>

      <div className="sidebar-scroll">
        <div className="sidebar-section">
          <div className="sidebar-section-title">Channels</div>
          {publicChannels.map(ch => renderChannelItem(ch, ch.name))}
        </div>

        <div className="sidebar-section">
          <div className="sidebar-section-header">
            <div className="sidebar-section-title">Direct Messages</div>
            <button
              className="sidebar-new-chat-btn"
              onClick={onStartNewChat}
              aria-label="Start new chat"
              title="Start new chat"
            >
              +
            </button>
          </div>
          {dmChannels.map(ch => renderChannelItem(ch, getDmLabel(ch)))}
        </div>
      </div>
    </div>
  )
}

export default Sidebar