import React, { useState } from 'react'
import { login, register } from '../api/client'
import type { BackendUser } from '../api/client'
import './LoginScreen.css'

interface LoginScreenProps {
  onLogin: (user: BackendUser) => void
}

const LoginScreen: React.FC<LoginScreenProps> = ({ onLogin }) => {
  const [tab, setTab] = useState<'login' | 'register'>('login')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [email, setEmail] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await login({ username, password })
      onLogin(res.user)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await register({
        username,
        password,
        email,
        first_name: firstName,
        last_name: lastName,
      })
      onLogin(res.user)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  const switchTab = (t: 'login' | 'register') => {
    setTab(t)
    setError('')
    setUsername('')
    setPassword('')
    setEmail('')
    setFirstName('')
    setLastName('')
  }

  return (
    <div className="login-screen">
      <div className="login-card">
        <h1 className="login-title">MyChatApp</h1>

        <div className="login-tabs">
          <button
            className={`login-tab ${tab === 'login' ? 'login-tab--active' : ''}`}
            onClick={() => switchTab('login')}
          >
            Sign In
          </button>
          <button
            className={`login-tab ${tab === 'register' ? 'login-tab--active' : ''}`}
            onClick={() => switchTab('register')}
          >
            Register
          </button>
        </div>

        {tab === 'login' ? (
          <form onSubmit={handleLogin}>
            <div className="login-form-group">
              <label className="login-label">Username</label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                className="login-input"
                placeholder="jdoe"
                autoFocus
                required
                disabled={loading}
              />
            </div>
            <div className="login-form-group">
              <label className="login-label">Password</label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                className="login-input"
                placeholder="••••••••"
                required
                disabled={loading}
              />
            </div>
            {error && <div className="login-error">{error}</div>}
            <button type="submit" className="login-button" disabled={loading}>
              {loading ? 'Signing in...' : 'Sign In'}
            </button>
          </form>
        ) : (
          <form onSubmit={handleRegister}>
            <div className="login-form-group">
              <label className="login-label">Username</label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                className="login-input"
                placeholder="jdoe"
                autoFocus
                required
                disabled={loading}
              />
            </div>
            <div className="login-form-group">
              <label className="login-label">Password</label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                className="login-input"
                placeholder="At least 6 characters"
                minLength={6}
                required
                disabled={loading}
              />
            </div>
            <div className="login-form-group">
              <label className="login-label">Email</label>
              <input
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                className="login-input"
                placeholder="jdoe@example.com"
                required
                disabled={loading}
              />
            </div>
            <div className="login-form-row">
              <div className="login-form-group">
                <label className="login-label">First Name</label>
                <input
                  type="text"
                  value={firstName}
                  onChange={e => setFirstName(e.target.value)}
                  className="login-input"
                  placeholder="John"
                  required
                  disabled={loading}
                />
              </div>
              <div className="login-form-group">
                <label className="login-label">Last Name</label>
                <input
                  type="text"
                  value={lastName}
                  onChange={e => setLastName(e.target.value)}
                  className="login-input"
                  placeholder="Doe"
                  required
                  disabled={loading}
                />
              </div>
            </div>
            {error && <div className="login-error">{error}</div>}
            <button type="submit" className="login-button" disabled={loading}>
              {loading ? 'Creating account...' : 'Create Account'}
            </button>
          </form>
        )}
      </div>
    </div>
  )
}

export default LoginScreen