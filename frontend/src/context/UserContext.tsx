import { createContext, useContext } from 'react'
import type { BackendUser } from '../api/client'

const UserContext = createContext<BackendUser | null>(null)

export const UserProvider = UserContext.Provider

export const useUser = (): BackendUser => {
  const user = useContext(UserContext)
  if (!user) {
    throw new Error('useUser must be used within a UserProvider')
  }
  return user
}
