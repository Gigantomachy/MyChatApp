import React, { useState, useEffect } from 'react'
import { getFriends, type BackendUser } from '../api/client'
import './NewChatModal.css'

interface NewChatModalProps {
  isOpen: boolean
  onClose: () => void
  onStartDM: (friendId: string) => void
}

const NewChatModal: React.FC<NewChatModalProps> = ({
  isOpen,
  onClose,
  onStartDM,
}) => {
  const [search, setSearch] = useState('')
  const [friends, setFriends] = useState<BackendUser[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (!isOpen) return
    setLoading(true)
    setError('')
    getFriends()
      .then(f => { setFriends(f ?? []) })
      .catch(err => { setError(err instanceof Error ? err.message : 'Failed to load friends') })
      .finally(() => setLoading(false))
  }, [isOpen])

  const filtered = friends.filter(f => {
    const fullName = `${f.first_name} ${f.last_name}`.toLowerCase()
    const handle = f.username.toLowerCase()
    const q = search.toLowerCase()
    return fullName.includes(q) || handle.includes(q)
  })

  const handleStart = (friendId: string) => {
    setCreating(true)
    setError('')
    onStartDM(friendId)
    setCreating(false)
    setSearch('')
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
    <div className="new-chat-overlay" onClick={onClose}>
      <div className="new-chat-modal" onClick={e => e.stopPropagation()}>
        <div className="new-chat-header">
          <h3>Start a DM</h3>
          <button className="new-chat-close" onClick={onClose} aria-label="Close">
            ×
          </button>
        </div>

        <input
          className="new-chat-search"
          type="text"
          placeholder="Search friends..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          autoFocus
        />

        <div className="new-chat-list">
          {loading && <div className="new-chat-empty">Loading...</div>}
          {error && <div className="new-chat-empty">{error}</div>}
          {!loading && !error && filtered.length === 0 && (
            <div className="new-chat-empty">No friends found.</div>
          )}
          {!loading && filtered.map(f => (
            <div
              key={f.user_id}
              className="new-chat-row"
              onClick={() => handleStart(f.user_id)}
            >
              <div className="new-chat-avatar">
                {f.first_name[0]}{f.last_name[0]}
              </div>
              <div className="new-chat-info">
                <div className="new-chat-name">{f.first_name} {f.last_name}</div>
                <div className="new-chat-handle">@{f.username}</div>
              </div>
              <div className="new-chat-dm-label">DM</div>
            </div>
          ))}
        </div>

        {creating && (
          <div className="new-chat-footer">
            <span className="new-chat-empty">Starting DM...</span>
          </div>
        )}
      </div>
    </div>
  )
}

export default NewChatModal
