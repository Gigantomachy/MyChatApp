import React, { useState } from 'react'
import { users } from '../data/database'
import './LoginScreen.css'

interface LoginScreenProps {
  onLogin: (userId: string) => void
}

const LoginScreen: React.FC<LoginScreenProps> = ({ onLogin }) => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const user = users.find(
      u => u.username === username && u.password === password
    )
    if (user) {
      setError('')
      onLogin(user.id)
    } else {
      setError('Invalid username or password')
    }
  }

  return (
    <div className="login-screen">
      <div className="login-card">
        <h1 className="login-title">MyChatApp</h1>
        <p className="login-subtitle">Sign in to continue</p>

        <form onSubmit={handleSubmit}>
          <div className="login-form-group">
            <label className="login-label">Username</label>
            <input
              type="text"
              value={username}
              onChange={e => setUsername(e.target.value)}
              className="login-input"
              placeholder="steven"
              autoFocus
            />
          </div>

          <div className="login-form-group">
            <label className="login-label">Password</label>
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              className="login-input"
              placeholder="password123"
            />
          </div>

          {error && <div className="login-error">{error}</div>}

          <button type="submit" className="login-button">
            Sign In
          </button>
        </form>

        <p className="login-hint">Try steven / password123</p>
      </div>
    </div>
  )
}

export default LoginScreen
