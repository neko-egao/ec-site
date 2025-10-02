'use client';

import React, { createContext, useState, useContext, useEffect } from 'react';

// ★ ユーザー情報の型を定義
type User = {
  id: number;
  name: string;
  email: string;
  role: string;
}

// Contextに保存するデータの型を定義
interface AuthContextType {
  user: User | null; // ★ tokenからuserオブジェクトに変更
  token: string | null;
  setUser: (user: User | null, token: string | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUserState] = useState<User | null>(null);
  const [token, setTokenState] = useState<string | null>(null);

  useEffect(() => {
    // ★ localStorageからuserとtokenの両方を読み込む
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    if (storedToken && storedUser) {
      setTokenState(storedToken);
      setUserState(JSON.parse(storedUser));
    }
  }, []);

  const setUser = (newUser: User | null, newToken: string | null) => {
    setUserState(newUser);
    setTokenState(newToken);

    // ★ localStorageにuserとtokenの両方を保存/削除
    if (newUser && newToken) {
      localStorage.setItem('user', JSON.stringify(newUser));
      localStorage.setItem('token', newToken);
    } else {
      localStorage.removeItem('user');
      localStorage.removeItem('token');
    }
  };

  return (
    <AuthContext.Provider value={{ user, token, setUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}