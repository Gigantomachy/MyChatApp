import React, { useState } from 'react'
import { createDmOrGroupChannel } from '../data/database'
import type { BackendUser } from '../api/client'
import Sidebar from './Sidebar'
import ChatArea from './ChatArea'
import NewChatModal from './NewChatModal'
import SearchModal from './SearchModal'
import './ChatLayout.css'

interface ChatLayoutProps {
  currentUser: BackendUser
  onLogout: () => void
}

const ChatLayout: React.FC<ChatLayoutProps> = ({ currentUser }) => {
  const [selectedChannelId, setSelectedChannelId] = useState<string | null>(null)
  const [isNewChatOpen, setIsNewChatOpen] = useState(false)
  const [isSearchOpen, setIsSearchOpen] = useState(false)

  const handleStartChat = (memberIds: string[]) => {
    const channel = createDmOrGroupChannel(currentUser.user_id, memberIds)
    setSelectedChannelId(channel.id)
    setIsNewChatOpen(false)
  }

  return (
    <div className="chat-layout">
      <div className="top-bar">
        <div className="top-bar-brand">MyChatApp</div>

        <div className="top-bar-user">
          <div
            className="top-bar-avatar"
            title={`${currentUser.first_name} ${currentUser.last_name}`}
          >
            {currentUser.first_name[0]}{currentUser.last_name[0]}
          </div>
        </div>
      </div>

      <div className="chat-layout-body">
        <Sidebar
          currentUser={currentUser}
          selectedChannelId={selectedChannelId}
          onSelectChannel={setSelectedChannelId}
          onStartNewChat={() => setIsNewChatOpen(true)}
          onOpenSearch={() => setIsSearchOpen(true)}
        />
        <ChatArea
          currentUser={currentUser}
          channelId={selectedChannelId ?? ''}
        />
      </div>

      <NewChatModal
        currentUser={currentUser}
        isOpen={isNewChatOpen}
        onClose={() => setIsNewChatOpen(false)}
        onStartChat={handleStartChat}
      />
      <SearchModal
        currentUser={currentUser}
        isOpen={isSearchOpen}
        onClose={() => setIsSearchOpen(false)}
      />
    </div>
  )
}

export default ChatLayout