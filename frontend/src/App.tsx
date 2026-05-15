import React, { useState } from 'react'
import LoginScreen from './components/LoginScreen'
import ChatLayout from './components/ChatLayout'

const App: React.FC = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)

  const handleLogin = (userId: string) => {
    setCurrentUserId(userId)
    setIsLoggedIn(true)
  }

  const handleLogout = () => {
    setCurrentUserId(null)
    setIsLoggedIn(false)
  }

  if (!isLoggedIn || !currentUserId) {
    return <LoginScreen onLogin={handleLogin} />
  }

  return <ChatLayout currentUserId={currentUserId} onLogout={handleLogout} />
}

export default App
