import React, { useState, useEffect } from 'react'
import LoginScreen from './components/LoginScreen'
import ChatLayout from './components/ChatLayout'
import type { BackendUser } from './api/client'
import { me, logout as apiLogout } from './api/client'
import { UserProvider } from './context/UserContext'

const App: React.FC = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [currentUser, setCurrentUser] = useState<BackendUser | null>(null)
  const [checkingSession, setCheckingSession] = useState(true)

  useEffect(() => {
    me()
      .then(res => {
        setCurrentUser(res.user)
        setIsLoggedIn(true)
      })
      .catch(() => {
        // not authenticated — show login screen
      })
      .finally(() => setCheckingSession(false))
  }, [])

  const handleLogin = (user: BackendUser) => {
    setCurrentUser(user)
    setIsLoggedIn(true)
  }

  const handleLogout = async () => {
    try {
      await apiLogout()
    } catch {
      // ignore — clearing state regardless
    }
    setCurrentUser(null)
    setIsLoggedIn(false)
  }

  if (checkingSession) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
        <span style={{ color: 'var(--text-muted)' }}>Loading...</span>
      </div>
    )
  }

  if (!isLoggedIn || !currentUser) {
    return <LoginScreen onLogin={handleLogin} />
  }

  return (
    <UserProvider value={currentUser}>
      <ChatLayout onLogout={handleLogout} />
    </UserProvider>
  )
}

export default App
