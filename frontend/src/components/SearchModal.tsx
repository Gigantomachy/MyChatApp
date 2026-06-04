import React, { useState, useEffect } from 'react'
import { allUsers, getFriendsForUser } from '../data/database'
import type { BackendUser } from '../api/client'
import type { User } from '../data/database'
import './SearchModal.css'

interface SearchModalProps {
  currentUser: BackendUser
  isOpen: boolean
  onClose: () => void
}

const isFriend = (userId: string, friends: User[]) =>
  friends.some(f => f.id === userId)

const SearchModal: React.FC<SearchModalProps> = ({ currentUser, isOpen, onClose }) => {
  const [search, setSearch] = useState('')
  const [activeTab, setActiveTab] = useState<'channels' | 'people'>('channels')

  const friends = getFriendsForUser(currentUser.user_id)

  const nonFriends = allUsers.filter(
    u => u.id !== currentUser.user_id && !isFriend(u.id, friends)
  )

  const filteredPeople = nonFriends.filter(u => {
    const fullName = `${u.firstName} ${u.lastName}`.toLowerCase()
    const handle = u.username.toLowerCase()
    const q = search.toLowerCase()
    return fullName.includes(q) || handle.includes(q)
  })

  const handleAddFriend = (userId: string) => {
    console.log(`[dummy] Send friend request to ${userId}`)
  }

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
            <div className="search-empty">Channel search not yet available.</div>
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