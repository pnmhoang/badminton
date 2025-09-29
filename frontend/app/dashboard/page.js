'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import axios from 'axios'

export default function DashboardPage() {
  const [user, setUser] = useState(null)
  const [tournaments, setTournaments] = useState([])
  const [loading, setLoading] = useState(true)
  const router = useRouter()

  const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')

    if (!token || !userData) {
      router.push('/auth/login')
      return
    }

    setUser(JSON.parse(userData))
    fetchTournaments(token)
  }, [])

  const fetchTournaments = async (token) => {
    try {
      const response = await axios.get(`${API_BASE}/api/v1/tournaments`, {
        headers: { Authorization: `Bearer ${token}` }
      })
      setTournaments(response.data || [])
    } catch (error) {
      console.error('Error fetching tournaments:', error)
      if (error.response?.status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        router.push('/auth/login')
      }
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/')
  }

  if (loading) return <div>Loading...</div>

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1>Player Dashboard</h1>
        <button onClick={handleLogout} className="button" style={{ background: '#dc3545' }}>
          Logout
        </button>
      </div>

      {user && (
        <div className="card" style={{ marginBottom: '2rem' }}>
          <h2>Welcome, {user.full_name}!</h2>
          <p><strong>Username:</strong> {user.username}</p>
          <p><strong>Email:</strong> {user.email}</p>
          <p><strong>Role:</strong> {user.role}</p>
          {user.ranking > 0 && <p><strong>Ranking:</strong> {user.ranking}</p>}
        </div>
      )}

      <div className="card">
        <h2>Available Tournaments</h2>
        {tournaments.length === 0 ? (
          <p>No tournaments available at the moment.</p>
        ) : (
          <div className="grid">
            {tournaments.map((tournament) => (
              <div key={tournament.id} className="card">
                <h3>{tournament.name}</h3>
                <p>{tournament.description}</p>
                <p><strong>Start:</strong> {new Date(tournament.start_date).toLocaleDateString()}</p>
                <p><strong>Status:</strong> {tournament.status}</p>
                <p><strong>Max Players:</strong> {tournament.max_players}</p>
                <button className="button">Register</button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}