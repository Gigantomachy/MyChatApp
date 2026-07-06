import React from 'react'
import { useUser } from '../context/UserContext'
import type { ChannelMembership } from '../api/client'
import './Sidebar.css'

interface SidebarProps {
  channels: ChannelMembership[]
  selectedChannelId: string | null
  onSelectChannel: (channelId: string) => void
  onStartNewChat: () => void
  onCreateChannel: () => void
  onOpenSearch: () => void
}

const Sidebar: React.FC<SidebarProps> = ({
  channels,
  selectedChannelId,
  onSelectChannel,
  onStartNewChat,
  onCreateChannel,
  onOpenSearch,
}) => {
  const currentUser = useUser()

  const publicChannels = channels.filter(ch => ch.channel_type === 'public')
  const dmChannels = channels.filter(ch => ch.channel_type === 'dm' || ch.channel_type === 'group')

  const renderChannelItem = (ch: ChannelMembership, label: string) => {
    const isActive = selectedChannelId === ch.channel_id
    return (
      <div
        key={ch.channel_id}
        onClick={() => onSelectChannel(ch.channel_id)}
        className={`sidebar-channel ${isActive ? 'sidebar-channel--active' : ''}`}
      >
        <span className="sidebar-channel-icon">
          {ch.channel_type === 'public' ? '#' : '@'}
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
          <div className="sidebar-section-header">
            <div className="sidebar-section-title">Channels</div>
            <button
              className="sidebar-new-chat-btn"
              onClick={onCreateChannel}
              aria-label="Create channel"
              title="Create channel"
            >
              +
            </button>
          </div>
          {publicChannels.length === 0 ? (
            <div className="sidebar-empty">No channels joined</div>
          ) : (
            publicChannels.map(ch => renderChannelItem(ch, ch.channel_name))
          )}
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
          {dmChannels.length === 0 ? (
            <div className="sidebar-empty">No direct messages</div>
          ) : (
            dmChannels.map(ch => renderChannelItem(ch, ch.channel_name))
          )}
        </div>
      </div>
    </div>
  )
}

export default Sidebar
