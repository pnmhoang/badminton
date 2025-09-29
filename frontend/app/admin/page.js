'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import axios from 'axios'

export default function AdminPage() {
  const [user, setUser] = useState(null)
  const [tournaments, setTournaments] = useState([])
  const [showTournamentForm, setShowTournamentForm] = useState(false)
  const [loading, setLoading] = useState(true)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    type: 'singles',
    start_date: '',
    end_date: '',
    max_players: 16,
    max_teams: 8,
    entry_fee: 0,
    prize_pool: 0
  })
  const router = useRouter()

  const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')

    if (!token || !userData) {
      router.push('/auth/login')
      return
    }

    const user = JSON.parse(userData)
    if (user.role !== 'admin') {
      router.push('/dashboard')
      return
    }

    setUser(user)
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
        handleLogout()
      }
    } finally {
      setLoading(false)
    }
  }

  const handleCreateTournament = async (e) => {
    e.preventDefault()
    const token = localStorage.getItem('token')

    try {
      await axios.post(`${API_BASE}/api/v1/tournaments`, formData, {
        headers: { Authorization: `Bearer ${token}` }
      })
      
      setShowTournamentForm(false)
      setFormData({
        name: '',
        description: '',
        type: 'singles',
        start_date: '',
        end_date: '',
        max_players: 16,
        max_teams: 8,
        entry_fee: 0,
        prize_pool: 0
      })
      fetchTournaments(token)
    } catch (error) {
      console.error('Error creating tournament:', error)
      alert('Failed to create tournament')
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
        <h1>Admin Dashboard</h1>
        <div>
          <button 
            onClick={() => setShowTournamentForm(!showTournamentForm)}
            className="button"
            style={{ marginRight: '1rem' }}
          >
            {showTournamentForm ? 'Cancel' : 'Create Tournament'}
          </button>
          <button onClick={handleLogout} className="button" style={{ background: '#dc3545' }}>
            Logout
          </button>
        </div>
      </div>

      {user && (
        <div className="card" style={{ marginBottom: '2rem' }}>
          <h2>Welcome, {user.full_name}!</h2>
          <p><strong>Admin Panel</strong> - Manage tournaments and users</p>
        </div>
      )}

      {showTournamentForm && (
        <form className="form" onSubmit={handleCreateTournament} style={{ marginBottom: '2rem' }}>
          <h2>Create New Tournament</h2>
          
          <div className="form-group">
            <label>Tournament Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label>Description</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              rows="3"
            />
          </div>

          <div className="form-group">
            <label>Tournament Type</label>
            <select
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value })}
            >
              <option value="singles">Singles</option>
              <option value="doubles">Doubles</option>
            </select>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <div className="form-group">
              <label>Start Date</label>
              <input
                type="datetime-local"
                value={formData.start_date}
                onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                required
              />
            </div>

            <div className="form-group">
              <label>End Date</label>
              <input
                type="datetime-local"
                value={formData.end_date}
                onChange={(e) => setFormData({ ...formData, end_date: e.target.value })}
                required
              />
            </div>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <div className="form-group">
              <label>Max Players</label>
              <input
                type="number"
                value={formData.max_players}
                onChange={(e) => setFormData({ ...formData, max_players: parseInt(e.target.value) })}
                min="2"
              />
            </div>

            <div className="form-group">
              <label>Max Teams (Doubles)</label>
              <input
                type="number"
                value={formData.max_teams}
                onChange={(e) => setFormData({ ...formData, max_teams: parseInt(e.target.value) })}
                min="2"
              />
            </div>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <div className="form-group">
              <label>Entry Fee ($)</label>
              <input
                type="number"
                value={formData.entry_fee}
                onChange={(e) => setFormData({ ...formData, entry_fee: parseFloat(e.target.value) })}
                min="0"
                step="0.01"
              />
            </div>

            <div className="form-group">
              <label>Prize Pool ($)</label>
              <input
                type="number"
                value={formData.prize_pool}
                onChange={(e) => setFormData({ ...formData, prize_pool: parseFloat(e.target.value) })}
                min="0"
                step="0.01"
              />
            </div>
          </div>

          <button type="submit" className="button">Create Tournament</button>
        </form>
      )}

      <div className="card">
        <h2>Manage Tournaments</h2>
        {tournaments.length === 0 ? (
          <p>No tournaments created yet.</p>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Start Date</th>
                <th>Status</th>
                <th>Participants</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {tournaments.map((tournament) => (
                <tr key={tournament.id}>
                  <td>{tournament.name}</td>
                  <td>{tournament.type || 'singles'}</td>
                  <td>{new Date(tournament.start_date).toLocaleDateString()}</td>
                  <td>
                    <span style={{ 
                      padding: '0.25rem 0.5rem', 
                      borderRadius: '4px', 
                      fontSize: '0.875rem',
                      background: tournament.status === 'upcoming' ? '#e3f2fd' : 
                                 tournament.status === 'ongoing' ? '#fff3e0' : '#f1f8e9',
                      color: tournament.status === 'upcoming' ? '#1565c0' :
                             tournament.status === 'ongoing' ? '#ef6c00' : '#388e3c'
                    }}>
                      {tournament.status}
                    </span>
                  </td>
                  <td>{tournament.match_count || 0}</td>
                  <td>
                    <div className="actions">
                      <button className="btn-sm btn-edit">Edit</button>
                      <button className="btn-sm btn-delete">Delete</button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}