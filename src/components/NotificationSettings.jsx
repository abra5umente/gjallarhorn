import React, { useState, useEffect } from 'react'
import { useNotifications } from '../context/NotificationContext'
import LoadingSpinner from './LoadingSpinner'

const NotificationSettings = () => {
  const { config, loading, updateConfig } = useNotifications()
  const [formData, setFormData] = useState({
    userKey: '',
    appToken: '',
    enabled: false
  })
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  useEffect(() => {
    if (config) {
      setFormData(config)
    }
  }, [config])

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)
    setError('')
    setSuccess('')

    try {
      await updateConfig(formData)
      setSuccess('Notification settings saved successfully!')
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <LoadingSpinner />
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Notification Settings</h1>
        <p className="text-gray-600 mt-2">
          Configure Pushover notifications for service downtime alerts.
        </p>
      </div>

      <div className="card">
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <div className="flex">
                <svg className="w-5 h-5 text-red-400 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p className="text-red-800">{error}</p>
              </div>
            </div>
          )}

          {success && (
            <div className="bg-green-50 border border-green-200 rounded-lg p-4">
              <div className="flex">
                <svg className="w-5 h-5 text-green-400 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p className="text-green-800">{success}</p>
              </div>
            </div>
          )}

          <div className="flex items-center">
            <input
              type="checkbox"
              id="enabled"
              name="enabled"
              checked={formData.enabled}
              onChange={handleChange}
              className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
            />
            <label htmlFor="enabled" className="ml-2 block text-sm font-medium text-gray-700">
              Enable Pushover notifications
            </label>
          </div>

          <div>
            <label htmlFor="userKey" className="block text-sm font-medium text-gray-700 mb-2">
              Pushover User Key
            </label>
            <input
              type="text"
              id="userKey"
              name="userKey"
              value={formData.userKey}
              onChange={handleChange}
              className="input-field"
              placeholder="Your Pushover user key"
            />
            <p className="text-sm text-gray-500 mt-1">
              Your unique user key from Pushover. Get it from{' '}
              <a 
                href="https://pushover.net/" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-primary-600 hover:text-primary-700"
              >
                pushover.net
              </a>
            </p>
          </div>

          <div>
            <label htmlFor="appToken" className="block text-sm font-medium text-gray-700 mb-2">
              Pushover App Token
            </label>
            <input
              type="text"
              id="appToken"
              name="appToken"
              value={formData.appToken}
              onChange={handleChange}
              className="input-field"
              placeholder="Your Pushover app token"
            />
            <p className="text-sm text-gray-500 mt-1">
              The application token for Gjallarhorn. You can use the default token or create your own app.
            </p>
          </div>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h3 className="text-sm font-medium text-blue-800 mb-2">How to get Pushover credentials:</h3>
            <ol className="text-sm text-blue-700 space-y-1 list-decimal list-inside">
              <li>Sign up for a free account at <a href="https://pushover.net/" target="_blank" rel="noopener noreferrer" className="underline">pushover.net</a></li>
              <li>Find your User Key on the main page after logging in</li>
              <li>Create a new application or use the default token</li>
              <li>Enter both credentials above and enable notifications</li>
            </ol>
          </div>

          <div className="flex space-x-4">
            <button
              type="submit"
              disabled={saving}
              className="btn-primary flex-1"
            >
              {saving ? 'Saving...' : 'Save Settings'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default NotificationSettings
