import React, { createContext, useContext, useState, useEffect } from 'react'
import { api } from '../services/api'

const NotificationContext = createContext()

export const useNotifications = () => {
  const context = useContext(NotificationContext)
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider')
  }
  return context
}

export const NotificationProvider = ({ children }) => {
  const [config, setConfig] = useState({
    userKey: '',
    appToken: '',
    enabled: false
  })
  const [loading, setLoading] = useState(true)

  const fetchConfig = async () => {
    try {
      setLoading(true)
      const data = await api.getNotificationConfig()
      setConfig(data)
    } catch (err) {
      console.error('Failed to fetch notification config:', err)
    } finally {
      setLoading(false)
    }
  }

  const updateConfig = async (newConfig) => {
    await api.updateNotificationConfig(newConfig)
    setConfig(newConfig)
  }

  useEffect(() => {
    fetchConfig()
  }, [])

  const value = {
    config,
    loading,
    updateConfig,
  }

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  )
}
