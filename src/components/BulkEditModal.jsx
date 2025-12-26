import React, { useState, useEffect } from 'react'
import { useServices } from '../context/ServiceContext'

const BulkEditModal = ({ isOpen, onClose }) => {
  const { services, selectedIds, bulkUpdateServices } = useServices()
  const [interval, setInterval] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // Reset form state when modal opens/closes
  useEffect(() => {
    if (isOpen) {
      setInterval('')
      setError('')
    }
  }, [isOpen])

  if (!isOpen) return null

  const selectedServices = services.filter(s => selectedIds.has(s.id))

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!interval) {
      setError('Please enter an interval')
      return
    }

    const intervalNum = parseInt(interval)
    if (isNaN(intervalNum) || intervalNum < 30 || intervalNum > 3600) {
      setError('Interval must be between 30 and 3600 seconds')
      return
    }

    setLoading(true)
    setError('')

    try {
      const updates = selectedServices.map(s => ({
        id: s.id,
        name: s.name,
        url: s.url,
        interval: intervalNum,
      }))

      await bulkUpdateServices(updates)
      onClose()
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <h2 className="text-xl font-bold mb-4">
          Edit {selectedServices.length} Service{selectedServices.length > 1 ? 's' : ''}
        </h2>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-3 mb-4 text-red-800">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              New Check Interval (seconds)
            </label>
            <input
              type="number"
              value={interval}
              onChange={(e) => setInterval(e.target.value)}
              className="input-field"
              min="30"
              max="3600"
              placeholder="30-3600"
            />
            <p className="text-sm text-gray-500 mt-1">
              This will update the check interval for all selected services.
            </p>
          </div>

          <div className="flex space-x-3">
            <button
              type="submit"
              disabled={loading}
              className="btn-primary flex-1"
            >
              {loading ? 'Updating...' : 'Update All'}
            </button>
            <button
              type="button"
              onClick={onClose}
              className="btn-secondary flex-1"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default BulkEditModal
