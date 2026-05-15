import React, { useState, useEffect } from 'react'
import { allUsers, getFriendsForUser, discoverableChannels } from '../data/database'
import type { User, Channel } from '../data/database'
import './SearchModal.css'

interface SearchModalProps {
  currentUserId: string
  isOpen: boolean
  onClose: () => void
}

// Helper: returns true if current user is already in the channel
const isChannelMember = (ch: Channel, userId: string) => ch.memberIds.includes(userId)

// Helper: returns true if current user is already friends with the person
const isFriend = (userId: string, friends: User[]) =>
  friends.some(f => f.id === userId)

const SearchModal: React.FC<SearchModalProps> = ({ currentUserId, isOpen, onClose }) => {
  const [search, setSearch] = useState('')
  const [activeTab, setActiveTab] = useState<'channels' | 'people'>('channels')

  const friends = getFriendsForUser(currentUserId)

  // Filter discoverable channels (not already joined)
  const joinableChannels = discoverableChannels.filter(
    ch => !isChannelMember(ch, currentUserId)
  )

  // Filter users that are not already friends
  const nonFriends = allUsers.filter(
    u => u.id !== currentUserId && !isFriend(u.id, friends)
  )

  const filteredChannels = joinableChannels.filter(ch => {
    const q = search.toLowerCase()
    return ch.name.toLowerCase().includes(q)
  })

  const filteredPeople = nonFriends.filter(u => {
    const fullName = `${u.firstName} ${u.lastName}`.toLowerCase()
    const handle = u.username.toLowerCase()
    const q = search.toLowerCase()
    return fullName.includes(q) || handle.includes(q)
  })

  const handleJoinChannel = (channelId: string) => {
    // Dummy action — does nothing for now
    console.log(`[dummy] Join channel ${channelId}`)
  }

  const handleAddFriend = (userId: string) => {
    // Dummy action — does nothing for now
    console.log(`[dummy] Send friend request to ${userId}`)
  }

  // Close on Escape
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    if (isOpen) window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [isOpen, onClose])

  if (!isOpen) return null

  return (
    <div className="search-overlay" onClick={onClose}>
      <div className="search-modal" onClick={e => e.stopPropagation()}>
        <div className="search-header">
          <h3>Search</h3>
          <button className="search-close" onClick={onClose} aria-label="Close">
            ×
          </button>
        </div>

        <input
          className="search-input"
          type="text"
          placeholder="Search channels and people..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          autoFocus
        />

        <div className="search-tabs">
          <button
            className={`search-tab ${activeTab === 'channels' ? 'search-tab--active' : ''}`}
            onClick={() => setActiveTab('channels')}
          >
            Channels
          </button>
          <button
            className={`search-tab ${activeTab === 'people' ? 'search-tab--active' : ''}`}
            onClick={() => setActiveTab('people')}
          >
            People
          </button>
        </div>

        <div className="search-list">
          {activeTab === 'channels' && (
            <>
              {filteredChannels.length === 0 && (
                <div className="search-empty">No channels found.</div>
              )}
              {filteredChannels.map(ch => (
                <div key={ch.id} className="search-row">
                  <div className="search-row-icon">#</div>
                  <div className="search-info">
                    <div className="search-name">{ch.name}</div>
                    <div className="search-meta">Public channel</div>
                  </div>
                  <button
                    className="search-action-btn"
                    onClick={() => handleJoinChannel(ch.id)}
                  >
                    Join
                  </button>
                </div>
              ))}
            </>
          )}

          {activeTab === 'people' && (
            <>
              {filteredPeople.length === 0 && (
                <div className="search-empty">No people found.</div>
              )}
              {filteredPeople.map(u => (
                <div key={u.id} className="search-row">
                  <div className="search-avatar">
                    {u.firstName[0]}{u.lastName[0]}
                  </div>
                  <div className="search-info">
                    <div className="search-name">{u.firstName} {u.lastName}</div>
                    <div className="search-meta">@{u.username}</div>
                  </div>
                  <button
                    className="search-action-btn"
                    onClick={() => handleAddFriend(u.id)}
                  >
                    Add Friend
                  </button>
                </div>
              ))}
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default SearchModal
