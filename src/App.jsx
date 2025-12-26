import React from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Header from './components/Header'
import ServiceList from './components/ServiceList'
import ServiceForm from './components/ServiceForm'
import NotificationSettings from './components/NotificationSettings'
import { ServiceProvider } from './context/ServiceContext'
import { NotificationProvider } from './context/NotificationContext'

function App() {
  return (
    <ServiceProvider>
      <NotificationProvider>
        <Router>
          <div className="min-h-screen bg-gray-50">
            <Header />
            <main className="container mx-auto px-4 py-8">
              <Routes>
                <Route path="/" element={<ServiceList />} />
                <Route path="/add" element={<ServiceForm />} />
                <Route path="/edit/:id" element={<ServiceForm />} />
                <Route path="/settings" element={<NotificationSettings />} />
              </Routes>
            </main>
          </div>
        </Router>
      </NotificationProvider>
    </ServiceProvider>
  )
}

export default App
