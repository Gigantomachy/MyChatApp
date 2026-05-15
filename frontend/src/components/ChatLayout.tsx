import React, { useState } from 'react'
import { users, createDmOrGroupChannel } from '../data/database'
import Sidebar from './Sidebar'
import ChatArea from './ChatArea'
import NewChatModal from './NewChatModal'
import SearchModal from './SearchModal'
import './ChatLayout.css'

interface ChatLayoutProps {
  currentUserId: string
  onLogout: () => void
}

const ChatLayout: React.FC<ChatLayoutProps> = ({ currentUserId }) => {
  const [selectedChannelId, setSelectedChannelId] = useState<string | null>(null)
  const [isNewChatOpen, setIsNewChatOpen] = useState(false)
  const [isSearchOpen, setIsSearchOpen] = useState(false)
  const currentUser = users.find(u => u.id === currentUserId)!

  const handleStartChat = (memberIds: string[]) => {
    const channel = createDmOrGroupChannel(currentUserId, memberIds)
    setSelectedChannelId(channel.id)
    setIsNewChatOpen(false)
  }

  return (
    <div className="chat-layout">
      {/* Top bar */}
      <div className="top-bar">
        <div className="top-bar-brand">MyChatApp</div>

        <div className="top-bar-user">
          <div
            className="top-bar-avatar"
            title={`${currentUser.firstName} ${currentUser.lastName}`}
          >
            {currentUser.firstName[0]}{currentUser.lastName[0]}
          </div>
        </div>
      </div>

      {/* Main layout */}
      <div className="chat-layout-body">
        <Sidebar
          currentUserId={currentUserId}
          selectedChannelId={selectedChannelId}
          onSelectChannel={setSelectedChannelId}
          onStartNewChat={() => setIsNewChatOpen(true)}
          onOpenSearch={() => setIsSearchOpen(true)}
        />
        <ChatArea
          currentUserId={currentUserId}
          channelId={selectedChannelId ?? ''}
        />
      </div>

      <NewChatModal
        currentUserId={currentUserId}
        isOpen={isNewChatOpen}
        onClose={() => setIsNewChatOpen(false)}
        onStartChat={handleStartChat}
      />
      <SearchModal
        currentUserId={currentUserId}
        isOpen={isSearchOpen}
        onClose={() => setIsSearchOpen(false)}
      />
    </div>
  )
}

export default ChatLayout
