import React, { useState, useEffect, useCallback } from 'react'
import { useUser } from '../context/UserContext'
import {
  searchUsers,
  getFriends,
  getFriendRequestByID,
  sendFriendRequest,
  cancelFriendRequest,
  type BackendUser,
} from '../api/client'
import './SearchModal.css'

interface SearchModalProps {
  isOpen: boolean
  onClose: () => void
}

const SearchModal: React.FC<SearchModalProps> = ({ isOpen, onClose }) => {
  const currentUser = useUser()
  const [search, setSearch] = useState('')
  const [activeTab, setActiveTab] = useState<'channels' | 'people'>('channels')
  const [allUsers, setAllUsers] = useState<BackendUser[]>([])
  const [friendIds, setFriendIds] = useState<Set<string>>(new Set())
  const [pendingSent, setPendingSent] = useState<Set<string>>(new Set())
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [actionError, setActionError] = useState('')

  const fetchData = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const [users, friends] = await Promise.all([searchUsers(), getFriends()])
      setAllUsers(users.filter(u => u.user_id !== currentUser.user_id))
      setFriendIds(new Set(friends.map(f => f.user_id)))

      const pending = new Set<string>()
      const checks = users
        .filter(u => u.user_id !== currentUser.user_id)
        .map(async u => {
          try {
            const req = await getFriendRequestByID(u.user_id)
            if (req.status === 'PENDING') {
              pending.add(u.user_id)
            }
          } catch {
            // 404 or error means no pending request, ignore
          }
        })
      await Promise.allSettled(checks)
      setPendingSent(pending)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load users')
    } finally {
      setLoading(false)
    }
  }, [currentUser.user_id])

  useEffect(() => {
    if (isOpen) {
      fetchData()
    }
  }, [isOpen, fetchData])

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    if (isOpen) window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [isOpen, onClose])

  const filteredPeople = allUsers.filter(u => {
    const fullName = `${u.first_name} ${u.last_name}`.toLowerCase()
    const handle = u.username.toLowerCase()
    const q = search.toLowerCase()
    return fullName.includes(q) || handle.includes(q)
  })

  const handleAddFriend = async (recipientId: string) => {
    setActionError('')
    try {
      await sendFriendRequest(recipientId)
      setPendingSent(prev => new Set(prev).add(recipientId))
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to send request')
    }
  }

  const handleCancel = async (recipientId: string) => {
    setActionError('')
    try {
      await cancelFriendRequest(recipientId)
      setPendingSent(prev => {
        const next = new Set(prev)
        next.delete(recipientId)
        return next
      })
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to cancel request')
    }
  }

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
              {actionError && <div className="search-error">{actionError}</div>}
              {error && <div className="search-error">{error}</div>}
              {loading && <div className="search-empty">Loading...</div>}
              {!loading && !error && filteredPeople.length === 0 && (
                <div className="search-empty">No people found.</div>
              )}
              {!loading && filteredPeople.map(u => {
                const isFriend = friendIds.has(u.user_id)
                const isPending = pendingSent.has(u.user_id)
                return (
                <div key={u.user_id} className="search-row">
                  <div className="search-avatar">
                    {u.first_name[0]}{u.last_name[0]}
                  </div>
                  <div className="search-info">
                    <div className="search-name">{u.first_name} {u.last_name}</div>
                    <div className="search-meta">@{u.username}</div>
                  </div>
                  {isFriend ? (
                    <span className="search-friend-badge">Friends</span>
                  ) : isPending ? (
                    <button
                      className="search-action-btn search-action-btn--cancel"
                      onClick={() => handleCancel(u.user_id)}
                    >
                      Cancel
                    </button>
                  ) : (
                    <button
                      className="search-action-btn"
                      onClick={() => handleAddFriend(u.user_id)}
                    >
                      Add Friend
                    </button>
                  )}
                </div>
              )})}
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export default SearchModal
