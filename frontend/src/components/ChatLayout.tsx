import React, { useState, useCallback, useEffect } from 'react'
import { useUser } from '../context/UserContext'
import {
  getMyChannels,
  createChannel,
  joinChannel,
  type ChannelMembership,
} from '../api/client'
import Sidebar from './Sidebar'
import ChatArea from './ChatArea'
import NewChatModal from './NewChatModal'
import SearchModal from './SearchModal'
import ProfilePanel from './ProfilePanel'
import CreateChannelModal from './CreateChannelModal'
import './ChatLayout.css'

interface ChatLayoutProps {
  onLogout: () => void
}

const ChatLayout: React.FC<ChatLayoutProps> = ({ onLogout }) => {
  const currentUser = useUser()
  const [selectedChannelId, setSelectedChannelId] = useState<string | null>(null)
  const [isNewChatOpen, setIsNewChatOpen] = useState(false)
  const [isSearchOpen, setIsSearchOpen] = useState(false)
  const [isProfileOpen, setIsProfileOpen] = useState(false)
  const [isCreateChannelOpen, setIsCreateChannelOpen] = useState(false)
  const [channels, setChannels] = useState<ChannelMembership[]>([])
  const [actionError, setActionError] = useState('')

  const fetchChannels = useCallback(async () => {
    try {
      const chs = await getMyChannels()
      setChannels(chs ?? [])
    } catch {
      // silent fail — sidebar just shows stale data
    }
  }, [])

  useEffect(() => {
    fetchChannels()
  }, [fetchChannels])

  const handleStartDM = async (friendId: string) => {
    setActionError('')
    try {
      const chn = await createChannel({
        channel_type: 'dm',
        members: [friendId],
      })
      await fetchChannels()
      setSelectedChannelId(chn.channel_id)
      setIsNewChatOpen(false)
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to start DM')
    }
  }

  const handleCreateChannel = async (name: string) => {
    setActionError('')
    try {
      const chn = await createChannel({
        channel_name: name,
        channel_type: 'public',
      })
      await fetchChannels()
      setSelectedChannelId(chn.channel_id)
      setIsCreateChannelOpen(false)
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to create channel')
    }
  }

  const handleJoinChannel = async (channelId: string) => {
    setActionError('')
    try {
      await joinChannel(channelId)
      await fetchChannels()
      setSelectedChannelId(channelId)
      setIsSearchOpen(false)
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to join channel')
    }
  }

  return (
    <div className="chat-layout">
      <div className="top-bar">
        <div className="top-bar-brand">MyChatApp</div>

        <div className="top-bar-user">
          <button
            className="top-bar-avatar-btn"
            onClick={() => setIsProfileOpen(prev => !prev)}
            title={`${currentUser.first_name} ${currentUser.last_name}`}
          >
            {currentUser.first_name[0]}{currentUser.last_name[0]}
          </button>
        </div>
      </div>

      {actionError && (
        <div className="chat-layout-error">{actionError}</div>
      )}

      {isProfileOpen && (
        <ProfilePanel
          onClose={() => setIsProfileOpen(false)}
          onLogout={onLogout}
        />
      )}

      <div className="chat-layout-body">
        <Sidebar
          channels={channels}
          selectedChannelId={selectedChannelId}
          onSelectChannel={setSelectedChannelId}
          onStartNewChat={() => setIsNewChatOpen(true)}
          onCreateChannel={() => setIsCreateChannelOpen(true)}
          onOpenSearch={() => setIsSearchOpen(true)}
        />
        <ChatArea channelId={selectedChannelId ?? ''} />

        <NewChatModal
          isOpen={isNewChatOpen}
          onClose={() => setIsNewChatOpen(false)}
          onStartDM={handleStartDM}
        />
        <SearchModal
          isOpen={isSearchOpen}
          onClose={() => setIsSearchOpen(false)}
          onJoinChannel={handleJoinChannel}
          myChannelIds={new Set(channels.map(c => c.channel_id))}
        />
        <CreateChannelModal
          isOpen={isCreateChannelOpen}
          onClose={() => setIsCreateChannelOpen(false)}
          onCreate={handleCreateChannel}
        />
      </div>
    </div>
  )
}

export default ChatLayout
