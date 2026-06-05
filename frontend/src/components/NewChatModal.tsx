import React, { useState, useEffect } from 'react'
import { getFriendsForUser } from '../data/database'
import { useUser } from '../context/UserContext'
import './NewChatModal.css'

interface NewChatModalProps {
  isOpen: boolean
  onClose: () => void
  onStartChat: (memberIds: string[]) => void
}

const NewChatModal: React.FC<NewChatModalProps> = ({
  isOpen,
  onClose,
  onStartChat,
}) => {
  const currentUser = useUser()
  const [search, setSearch] = useState('')
  const [selectedIds, setSelectedIds] = useState<string[]>([])

  const friends = getFriendsForUser(currentUser.user_id)

  const filtered = friends.filter(f => {
    const fullName = `${f.firstName} ${f.lastName}`.toLowerCase()
    const handle = f.username.toLowerCase()
    const q = search.toLowerCase()
    return fullName.includes(q) || handle.includes(q)
  })

  const toggleSelection = (id: string) => {
    setSelectedIds(prev =>
      prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]
    )
  }

  const handleStart = () => {
    if (selectedIds.length === 0) return
    onStartChat(selectedIds)
    setSearch('')
    setSelectedIds([])
  }

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    if (isOpen) window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [isOpen, onClose])

  if (!isOpen) return null

  const buttonLabel =
    selectedIds.length === 0
      ? 'Start Chat'
      : selectedIds.length === 1
      ? 'Start Chat'
      : 'Start Group Chat'

  return (
    <div className="new-chat-overlay" onClick={onClose}>
      <div className="new-chat-modal" onClick={e => e.stopPropagation()}>
        <div className="new-chat-header">
          <h3>New Chat</h3>
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
          {filtered.length === 0 && (
            <div className="new-chat-empty">No friends found.</div>
          )}
          {filtered.map(f => {
            const isSelected = selectedIds.includes(f.id)
            return (
              <div
                key={f.id}
                className={`new-chat-row ${isSelected ? 'new-chat-row--selected' : ''}`}
                onClick={() => toggleSelection(f.id)}
              >
                <div className="new-chat-avatar">
                  {f.firstName[0]}{f.lastName[0]}
                </div>
                <div className="new-chat-info">
                  <div className="new-chat-name">{f.firstName} {f.lastName}</div>
                  <div className="new-chat-handle">@{f.username}</div>
                </div>
                {isSelected && <div className="new-chat-check">✓</div>}
              </div>
            )
          })}
        </div>

        <div className="new-chat-footer">
          <button
            className="new-chat-start"
            disabled={selectedIds.length === 0}
            onClick={handleStart}
          >
            {buttonLabel}
          </button>
        </div>
      </div>
    </div>
  )
}

export default NewChatModal