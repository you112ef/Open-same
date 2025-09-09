import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const register = (data: any) => {
  return apiClient.post('/register', data);
};

export const login = (data: any) => {
  return apiClient.post('/login', data);
};
