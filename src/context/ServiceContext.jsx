import React, { createContext, useContext, useState, useEffect } from 'react'
import { api } from '../services/api'

const ServiceContext = createContext()

export const useServices = () => {
  const context = useContext(ServiceContext)
  if (!context) {
    throw new Error('useServices must be used within a ServiceProvider')
  }
  return context
}

export const ServiceProvider = ({ children }) => {
  const [services, setServices] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchServices = async () => {
    try {
      setLoading(true)
      const data = await api.getServices()
      setServices(data)
      setError(null)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const createService = async (serviceData) => {
    try {
      const newService = await api.createService(serviceData)
      setServices(prev => [...prev, newService])
      return newService
    } catch (err) {
      setError(err.message)
      throw err
    }
  }

  const updateService = async (id, serviceData) => {
    try {
      const updatedService = await api.updateService(id, serviceData)
      setServices(prev => prev.map(service => 
        service.id === id ? updatedService : service
      ))
      return updatedService
    } catch (err) {
      setError(err.message)
      throw err
    }
  }

  const deleteService = async (id) => {
    try {
      await api.deleteService(id)
      setServices(prev => prev.filter(service => service.id !== id))
    } catch (err) {
      setError(err.message)
      throw err
    }
  }

  useEffect(() => {
    fetchServices()
    
    // Refresh services every 30 seconds
    const interval = setInterval(fetchServices, 30000)
    return () => clearInterval(interval)
  }, [])

  const value = {
    services,
    loading,
    error,
    fetchServices,
    createService,
    updateService,
    deleteService,
  }

  return (
    <ServiceContext.Provider value={value}>
      {children}
    </ServiceContext.Provider>
  )
}
