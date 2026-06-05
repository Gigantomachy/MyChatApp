import React, { useState } from 'react'
import LoginScreen from './components/LoginScreen'
import ChatLayout from './components/ChatLayout'
import type { BackendUser } from './api/client'
import { UserProvider } from './context/UserContext'

const App: React.FC = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [_token, setToken] = useState<string | null>(null)
  const [currentUser, setCurrentUser] = useState<BackendUser | null>(null)

  const handleLogin = (token: string, user: BackendUser) => {
    setToken(token)
    setCurrentUser(user)
    setIsLoggedIn(true)
  }

  const handleLogout = () => {
    setToken(null)
    setCurrentUser(null)
    setIsLoggedIn(false)
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