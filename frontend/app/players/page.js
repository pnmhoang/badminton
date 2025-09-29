'use client'

import { useState, useEffect } from 'react'
import axios from 'axios'

export default function PlayersPage() {
  const [players, setPlayers] = useState([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    ranking: 0
  })

  const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  useEffect(() => {
    fetchPlayers()
  }, [])

  const fetchPlayers = async () => {
    try {
      const response = await axios.get(`${API_BASE}/api/v1/players`)
      setPlayers(response.data || [])
    } catch (error) {
      console.error('Error fetching players:', error)
      setPlayers([])
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await axios.post(`${API_BASE}/api/v1/players`, formData)
      setFormData({ name: '', email: '', ranking: 0 })
      setShowForm(false)
      fetchPlayers()
    } catch (error) {
      console.error('Error creating player:', error)
      alert('Error creating player')
    }
  }

  const handleDelete = async (id) => {
    if (confirm('Are you sure you want to delete this player?')) {
      try {
        await axios.delete(`${API_BASE}/api/v1/players/${id}`)
        fetchPlayers()
      } catch (error) {
        console.error('Error deleting player:', error)
        alert('Error deleting player')
      }
    }
  }

  if (loading) return <div>Loading...</div>

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1>Players</h1>
        <button 
          className="button" 
          onClick={() => setShowForm(!showForm)}
        >
          {showForm ? 'Cancel' : 'Add Player'}
        </button>
      </div>

      {showForm && (
        <form className="form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
          </div>
          <div className="form-group">
            <label>Email</label>
            <input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              required
            />
          </div>
          <div className="form-group">
            <label>Ranking</label>
            <input
              type="number"
              value={formData.ranking}
              onChange={(e) => setFormData({ ...formData, ranking: parseInt(e.target.value) })}
            />
          </div>
          <button type="submit" className="button">Create Player</button>
        </form>
      )}

      <table className="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Ranking</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {players.map((player) => (
            <tr key={player.id}>
              <td>{player.id}</td>
              <td>{player.name}</td>
              <td>{player.email}</td>
              <td>{player.ranking}</td>
              <td className="actions">
                <button 
                  className="btn-sm btn-delete"
                  onClick={() => handleDelete(player.id)}
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
          {players.length === 0 && (
            <tr>
              <td colSpan="5" style={{ textAlign: 'center', padding: '2rem' }}>
                No players found
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}