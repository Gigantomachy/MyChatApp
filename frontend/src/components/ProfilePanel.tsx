import React, { useEffect, useState } from 'react'
import { useUser } from '../context/UserContext'
import {
  getFriends,
  getFriendRequests,
  acceptFriendRequest,
  rejectFriendRequest,
  type BackendUser,
  type FriendRequestItem,
} from '../api/client'
import './ProfilePanel.css'

interface ProfilePanelProps {
  onClose: () => void
  onLogout: () => void
}

const ProfilePanel: React.FC<ProfilePanelProps> = ({ onClose, onLogout }) => {
  const user = useUser()
  const [friends, setFriends] = useState<BackendUser[]>([])
  const [requests, setRequests] = useState<FriendRequestItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [actionError, setActionError] = useState('')

  const fetchData = async () => {
    setLoading(true)
    setError('')
    try {
      const [f, r] = await Promise.all([getFriends(), getFriendRequests()])
      setFriends(f)
      setRequests(r)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  const handleAccept = async (senderId: string) => {
    setActionError('')
    try {
      await acceptFriendRequest(senderId)
      setRequests(prev => prev.filter(r => r.sender_id !== senderId))
      await fetchData()
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to accept')
    }
  }

  const handleReject = async (senderId: string) => {
    setActionError('')
    try {
      await rejectFriendRequest(senderId)
      setRequests(prev => prev.filter(r => r.sender_id !== senderId))
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to reject')
    }
  }

  return (
    <>
      <div className="profile-overlay" onClick={onClose} />
      <div className="profile-panel">
        <div className="profile-header">
          <div className="profile-avatar">
            {user.first_name[0]}{user.last_name[0]}
          </div>
          <div className="profile-user-info">
            <div className="profile-name">{user.first_name} {user.last_name}</div>
            <div className="profile-handle">@{user.username}</div>
          </div>
          <button className="profile-close" onClick={onClose} aria-label="Close">×</button>
        </div>

        {loading && <div className="profile-section-message">Loading...</div>}
        {error && <div className="profile-error">{error}</div>}
        {actionError && <div className="profile-error">{actionError}</div>}

        {!loading && !error && (
          <>
            <div className="profile-section">
              <div className="profile-section-title">
                Friends ({friends.length})
              </div>
              {friends.length === 0 ? (
                <div className="profile-section-message">No friends yet</div>
              ) : (
                <div className="profile-list">
                  {friends.map(f => (
                    <div key={f.user_id} className="profile-row">
                      <div className="profile-row-avatar">
                        {f.first_name[0]}{f.last_name[0]}
                      </div>
                      <div className="profile-row-info">
                        <div className="profile-row-name">{f.first_name} {f.last_name}</div>
                        <div className="profile-row-meta">@{f.username}</div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="profile-section">
              <div className="profile-section-title">
                Friend Requests ({requests.length})
              </div>
              {requests.length === 0 ? (
                <div className="profile-section-message">No pending requests</div>
              ) : (
                <div className="profile-list">
                  {requests.map(r => (
                    <div key={r.sender_id} className="profile-row">
                      <div className="profile-row-avatar">
                        {r.sender_first_name[0]}{r.sender_last_name[0]}
                      </div>
                      <div className="profile-row-info">
                        <div className="profile-row-name">{r.sender_first_name} {r.sender_last_name}</div>
                        <div className="profile-row-meta">@{r.sender_username}</div>
                      </div>
                      <div className="profile-row-actions">
                        <button
                          className="profile-btn profile-btn-accept"
                          onClick={() => handleAccept(r.sender_id)}
                        >
                          Accept
                        </button>
                        <button
                          className="profile-btn profile-btn-reject"
                          onClick={() => handleReject(r.sender_id)}
                        >
                          Decline
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </>
        )}

        <div className="profile-footer">
          <button className="profile-logout-btn" onClick={onLogout}>
            Sign Out
          </button>
        </div>
      </div>
    </>
  )
}

export default ProfilePanel
