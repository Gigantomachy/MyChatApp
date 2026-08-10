import React, { useState, useEffect } from 'react'
import './CreateChannelModal.css'

interface CreateChannelModalProps {
  isOpen: boolean
  onClose: () => void
  onCreate: (name: string) => void
}

const CreateChannelModal: React.FC<CreateChannelModalProps> = ({
  isOpen,
  onClose,
  onCreate,
}) => {
  const [name, setName] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    if (isOpen) {
      setName('')
      setError('')
    }
  }, [isOpen])

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    if (isOpen) window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [isOpen, onClose])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = name.trim()
    if (!trimmed) {
      setError('Channel name is required')
      return
    }
    const channelName = trimmed.replace(/^#+/, '').trim()
    onCreate(channelName)
  }

  if (!isOpen) return null

  return (
    <div className="create-channel-overlay" onClick={onClose}>
      <div className="create-channel-modal" onClick={e => e.stopPropagation()}>
        <div className="create-channel-header">
          <h3>Create Channel</h3>
          <button className="create-channel-close" onClick={onClose} aria-label="Close">
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className="create-channel-form">
          <label className="create-channel-label">Channel Name</label>
          <div className="create-channel-input-wrapper">
            <span className="create-channel-prefix">#</span>
            <input
              className="create-channel-input"
              type="text"
              placeholder="new-channel"
              value={name.startsWith('#') ? name.slice(1) : name}
              onChange={e => setName(e.target.value)}
              autoFocus
            />
          </div>
          {error && <div className="create-channel-error">{error}</div>}
          <button type="submit" className="create-channel-btn">
            Create Channel
          </button>
        </form>
      </div>
    </div>
  )
}

export default CreateChannelModal
