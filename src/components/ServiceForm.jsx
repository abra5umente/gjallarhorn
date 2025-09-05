import React, { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useServices } from '../context/ServiceContext'

const ServiceForm = () => {
  const navigate = useNavigate()
  const { id } = useParams()
  const { services, createService, updateService } = useServices()
  
  const [formData, setFormData] = useState({
    name: '',
    url: '',
    interval: 60
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const isEdit = Boolean(id)

  useEffect(() => {
    if (isEdit) {
      const service = services.find(s => s.id === id)
      if (service) {
        setFormData({
          name: service.name,
          url: service.url,
          interval: service.interval
        })
      }
    }
  }, [isEdit, id, services])

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: name === 'interval' ? parseInt(value) || 0 : value
    }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      if (isEdit) {
        await updateService(id, formData)
      } else {
        await createService(formData)
      }
      navigate('/')
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl mx-auto">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">
          {isEdit ? 'Edit Service' : 'Add New Service'}
        </h1>
        <p className="text-gray-600 mt-2">
          {isEdit 
            ? 'Update the service monitoring configuration.' 
            : 'Configure a new service to monitor for uptime.'
          }
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

          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
              Service Name
            </label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              className="input-field"
              placeholder="e.g., My Website"
              required
            />
          </div>

          <div>
            <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-2">
              URL
            </label>
            <input
              type="url"
              id="url"
              name="url"
              value={formData.url}
              onChange={handleChange}
              className="input-field"
              placeholder="https://example.com"
              required
            />
            <p className="text-sm text-gray-500 mt-1">
              The URL to monitor. Must include http:// or https://
            </p>
          </div>

          <div>
            <label htmlFor="interval" className="block text-sm font-medium text-gray-700 mb-2">
              Check Interval (seconds)
            </label>
            <input
              type="number"
              id="interval"
              name="interval"
              value={formData.interval}
              onChange={handleChange}
              className="input-field"
              min="30"
              max="3600"
              required
            />
            <p className="text-sm text-gray-500 mt-1">
              How often to check the service (minimum 30 seconds)
            </p>
          </div>

          <div className="flex space-x-4">
            <button
              type="submit"
              disabled={loading}
              className="btn-primary flex-1"
            >
              {loading ? 'Saving...' : (isEdit ? 'Update Service' : 'Add Service')}
            </button>
            <button
              type="button"
              onClick={() => navigate('/')}
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

export default ServiceForm
