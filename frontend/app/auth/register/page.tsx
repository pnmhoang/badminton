'use client';

import { useState, ChangeEvent, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { authAPI } from '../../../utils/api';
import { AuthForm, FormField } from '../../../components/auth/AuthForm';
import type { RegisterRequest } from '../../../types';

interface RegisterFormData extends RegisterRequest {
  confirmPassword: string;
}

export default function RegisterPage() {
  const [formData, setFormData] = useState<RegisterFormData>({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    full_name: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  const handleInputChange = (field: keyof RegisterFormData) => 
    (e: ChangeEvent<HTMLInputElement>) => {
      setFormData(prev => ({
        ...prev,
        [field]: e.target.value
      }));
    };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
      setLoading(false);
      return;
    }

    try {
      const { confirmPassword, ...registerData } = formData;
      const response = await authAPI.register(registerData);
      
      // Store token and user info
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
      
      // Redirect to dashboard (all new users are players)
      router.push('/dashboard');
    } catch (error: any) {
      console.error('Registration error:', error);
      setError(error.response?.data?.message || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthForm
      title="Register"
      onSubmit={handleSubmit}
      loading={loading}
      error={error}
      submitText="Register"
      footerText="Already have an account?"
      footerLink={{ text: 'Login here', href: '/auth/login' }}
    >
      <FormField
        label="Full Name"
        type="text"
        value={formData.full_name}
        onChange={handleInputChange('full_name')}
        required
        disabled={loading}
      />

      <FormField
        label="Username"
        type="text"
        value={formData.username}
        onChange={handleInputChange('username')}
        required
        disabled={loading}
      />

      <FormField
        label="Email"
        type="email"
        value={formData.email}
        onChange={handleInputChange('email')}
        required
        disabled={loading}
      />

      <FormField
        label="Password"
        type="password"
        value={formData.password}
        onChange={handleInputChange('password')}
        required
        minLength={6}
        disabled={loading}
      />

      <FormField
        label="Confirm Password"
        type="password"
        value={formData.confirmPassword}
        onChange={handleInputChange('confirmPassword')}
        required
        disabled={loading}
      />
    </AuthForm>
  );
}