export default function Home() {
  return (
    <div className="container">
      <h1>Welcome to Badminton Tournament Manager</h1>
      <div className="grid">
        <div className="card">
          <h2>Players</h2>
          <p>Manage player profiles and rankings</p>
          <a href="/players" className="button">View Players</a>
        </div>
        <div className="card">
          <h2>Matches</h2>
          <p>Schedule and track match results</p>
          <a href="/matches" className="button">View Matches</a>
        </div>
        <div className="card">
          <h2>Tournaments</h2>
          <p>Organize and manage tournaments</p>
          <a href="/tournaments" className="button">View Tournaments</a>
        </div>
      </div>
    </div>
  )
}