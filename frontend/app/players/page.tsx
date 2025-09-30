'use client';

import React from 'react';
import Layout from '../../components/Layout';

interface Player {
  id: string;
  name: string;
  email: string;
  skill_level: string;
  created_at: string;
}

const PlayersPage: React.FC = () => {
  return (
    <Layout>
      <div className="players-page">
        <h1>Players Management</h1>
        <p>Manage badminton players here.</p>
        
        <div className="players-list">
          {/* Players list will be implemented here */}
          <div className="placeholder">
            <p>Players list coming soon...</p>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default PlayersPage;